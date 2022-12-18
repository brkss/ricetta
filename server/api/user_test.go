package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/brkss/vanillefraise2/db/mock"
	db "github.com/brkss/vanillefraise2/db/sqlc"
	"github.com/brkss/vanillefraise2/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// handle testing hashed password !
type eqCreateUserParamsMarcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMarcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	err := utils.VerifyPassword(e.password, arg.Password)
	if err != nil {
		return false
	}
	e.arg.Password = arg.Password
	return true
}

func (e eqCreateUserParamsMarcher) String() string {
	return fmt.Sprintf("matches arg [%v] and password [%v]", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMarcher{arg, password}
}

func RandomUser(t *testing.T) (db.User, string) {
	password := utils.RandomString(10)
	hashed, err := utils.HashPassword(password)
	require.NoError(t, err)
	user := db.User{
		ID:        uuid.New().String(),
		Name:      utils.RandomName(),
		Email:     utils.RandomEmail(),
		Username:  utils.RandomName(),
		Password:  hashed,
		CreatedAt: time.Now(),
	}
	return user, password
}

func TestCreateUserAPI(t *testing.T) {
	user, password := RandomUser(t)
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"name":     user.Name,
				"username": user.Username,
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					ID:       user.ID,
					Name:     user.Name,
					Email:    user.Email,
					Username: user.Username,
					Password: user.Password,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatch(t, recorder.Body, user)
			},
		},

		{
			name: "Internal Error",
			body: gin.H{
				"name":     user.Name,
				"username": user.Username,
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "duplicated username",
			body: gin.H{
				"name":     user.Name,
				"password": user.Password,
				"username": user.Username,
				"email":    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
	}
	fmt.Println("test @")
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(store)
			recorder := httptest.NewRecorder()

			// Marshall Body Data
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/register"
			// create http request
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func requireBodyMatch(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	require.Equal(t, user.ID, gotUser.ID)
	require.Equal(t, user.Name, gotUser.Name)
	require.Equal(t, user.Email, gotUser.Email)
	require.Equal(t, user.Username, gotUser.Username)
	require.WithinDuration(t, user.CreatedAt, gotUser.CreatedAt, time.Second)
}
