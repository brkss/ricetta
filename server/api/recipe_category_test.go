package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/brkss/vanillefraise2/db/mock"
	db "github.com/brkss/vanillefraise2/db/sqlc"
	"github.com/brkss/vanillefraise2/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type eqCreateRecipeCategoryParamsMatcher struct {
	arg db.CreateRecipeCategoryParams
}

func (e eqCreateRecipeCategoryParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateRecipeCategoryParams)
	if !ok {
		return false
	}
	e.arg.ID = arg.ID
	return true
}

func (e eqCreateRecipeCategoryParamsMatcher) String() string {
	return fmt.Sprintf("id %v matches: ", e.arg.ID)
}

func EqCreateRecipeCategoryParams(arg db.CreateRecipeCategoryParams) gomock.Matcher {
	return eqCreateRecipeCategoryParamsMatcher{arg}
}

func RandomCategory() db.RecipeCategory {
	category := db.RecipeCategory{
		ID:    uuid.New().String(),
		Title: utils.RandomName(),
		Image: utils.RandomName(),
		///Active: true,
	}
	return category
}

func TestCreateCategoryAPI(t *testing.T) {
	category := RandomCategory()
	testCases := []struct {
		name          string
		body          gin.H
		buildStabs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"title": category.Title,
				"image": category.Image,
			},
			buildStabs: func(store *mockdb.MockStore) {
				arg := db.CreateRecipeCategoryParams{
					ID:    category.ID,
					Title: category.Title,
					Image: category.Image,
				}
				store.EXPECT().CreateRecipeCategory(gomock.Any(), EqCreateRecipeCategoryParams(arg)).
					Times(1).
					Return(category, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusOK)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{
				"title": category.Title,
			},
			buildStabs: func(store *mockdb.MockStore) {

				store.EXPECT().CreateRecipeCategory(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"title": category.Title,
				"image": category.Image,
			},
			buildStabs: func(store *mockdb.MockStore) {
				arg := db.CreateRecipeCategoryParams{
					ID:    category.ID,
					Title: category.Title,
					Image: category.Image,
				}
				store.EXPECT().CreateRecipeCategory(gomock.Any(), EqCreateRecipeCategoryParams(arg)).
					Times(1).
					Return(db.RecipeCategory{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusInternalServerError)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStabs(store)

			server := newTestServer(store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/create-category"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}

func TestGettingCategories(t *testing.T) {
	categories := make([]db.RecipeCategory, 5)
	for i := 0; i < 5; i++ {
		categories[i] = RandomCategory()
	}

	testCases := []struct {
		name          string
		body          gin.H
		buildStabs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"limit":  5,
				"offset": 5,
			},
			buildStabs: func(store *mockdb.MockStore) {
				arg := db.GetCategoriesParams{
					Limit:  5,
					Offset: 5,
				}
				store.EXPECT().
					GetCategories(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(categories, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"limit":  5,
				"offset": 5,
			},
			buildStabs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCategories(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.RecipeCategory{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{
				"limit":  -1,
				"offset": -1,
			},
			buildStabs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCategories(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStabs(store)

			server := newTestServer(store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/categories"
			request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
