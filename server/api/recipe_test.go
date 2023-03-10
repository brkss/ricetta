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
	"github.com/stretchr/testify/require"
)

type eqCreateRecipePramsMatcher struct {
	arg db.CreateRecipeParams
}

func (e eqCreateRecipePramsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateRecipeParams)
	if !ok {
		return false
	}
	e.arg.ID = arg.ID
	return true
}
func (e eqCreateRecipePramsMatcher) String() string {
	return fmt.Sprintf("id matches %v", e.arg.ID)
}

func EqCreateRecipeParams(arg db.CreateRecipeParams) gomock.Matcher {
	return &eqCreateRecipePramsMatcher{arg: arg}
}

func RandomRecipe() db.Recipe {
	recipe := db.Recipe{
		ID:          uuid.New().String(),
		Name:        utils.RandomName(),
		Description: utils.RandomString(100),
		Image:       utils.RandomName(),
		Active:      sql.NullBool{Bool: true, Valid: true},
		Time:        utils.RandomName(),
		Url:         utils.RandomString(50),
		Servings:    10,
		CreatedAt:   time.Now().Local(),
	}
	return recipe
}

func RandomRecipeCategoryAssingment(recipe db.Recipe, category db.RecipeCategory) db.RecipeCategoryRecipe {
	return db.RecipeCategoryRecipe{
		ID:               uuid.New().String(),
		RecipeID:         recipe.ID,
		RecipeCategoryID: category.ID,
	}
}

func TestCreateRecipe(t *testing.T) {
	recipe := RandomRecipe()
	category := RandomCategory()
	assigment := RandomRecipeCategoryAssingment(recipe, category)
	testCases := []struct {
		name          string
		body          gin.H
		buildStabs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"name":        recipe.Name,
				"description": recipe.Description,
				"image":       recipe.Image,
				"time":        recipe.Time,
				"url":         recipe.Url,
				"servings":    recipe.Servings,
				"category_id": category.ID,
			},
			buildStabs: func(store *mockdb.MockStore) {
				arg := db.CreateRecipeParams{
					ID:          category.ID,
					Name:        recipe.Name,
					Image:       recipe.Image,
					Description: recipe.Description,
					Url:         recipe.Url,
					Time:        recipe.Time,
					Servings:    recipe.Servings,
				}
				store.EXPECT().
					CreateRecipe(gomock.Any(), EqCreateRecipeParams(arg)).
					Times(1).
					Return(recipe, nil)
				store.EXPECT().
					AssignRecipeToCategory(gomock.Any(), gomock.Any()).
					Times(1).
					Return(assigment, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"name":        recipe.Name,
				"description": recipe.Description,
				"image":       recipe.Image,
				"time":        recipe.Time,
				"url":         recipe.Url,
				"servings":    recipe.Servings,
				"category_id": category.ID,
			},
			buildStabs: func(store *mockdb.MockStore) {
				arg := db.CreateRecipeParams{
					ID:          category.ID,
					Name:        recipe.Name,
					Image:       recipe.Image,
					Description: recipe.Description,
					Url:         recipe.Url,
					Time:        recipe.Time,
					Servings:    recipe.Servings,
				}
				store.EXPECT().
					CreateRecipe(gomock.Any(), EqCreateRecipeParams(arg)).
					Times(1).
					Return(db.Recipe{}, sql.ErrConnDone)
				store.EXPECT().
					AssignRecipeToCategory(gomock.Any(), gomock.Any()).
					Times(1).
					Return(assigment, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InternalError2",
			body: gin.H{
				"name":        recipe.Name,
				"description": recipe.Description,
				"image":       recipe.Image,
				"time":        recipe.Time,
				"url":         recipe.Url,
				"servings":    recipe.Servings,
				"category_id": category.ID,
			},
			buildStabs: func(store *mockdb.MockStore) {
				arg := db.CreateRecipeParams{
					ID:          category.ID,
					Name:        recipe.Name,
					Image:       recipe.Image,
					Description: recipe.Description,
					Url:         recipe.Url,
					Time:        recipe.Time,
					Servings:    recipe.Servings,
				}
				store.EXPECT().
					CreateRecipe(gomock.Any(), EqCreateRecipeParams(arg)).
					Times(1).
					Return(recipe, nil)
				store.EXPECT().
					AssignRecipeToCategory(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.RecipeCategoryRecipe{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{
				"name":        recipe.Name,
				"description": recipe.Description,
				"image":       recipe.Image,
				"time":        recipe.Time,
				"url":         recipe.Url,
				"servings":    recipe.Servings,
			},
			buildStabs: func(store *mockdb.MockStore) {
				arg := db.CreateRecipeParams{
					ID:          category.ID,
					Name:        recipe.Name,
					Image:       recipe.Image,
					Description: recipe.Description,
					Url:         recipe.Url,
					Time:        recipe.Time,
					Servings:    recipe.Servings,
				}
				store.EXPECT().
					CreateRecipe(gomock.Any(), EqCreateRecipeParams(arg)).
					Times(1).
					Return(db.Recipe{}, sql.ErrConnDone)
				store.EXPECT().
					AssignRecipeToCategory(gomock.Any(), gomock.Any()).
					Times(1).
					Return(assigment, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(t.Name(), func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStabs(store)

			server := newTestServer(store)
			recorder := httptest.NewRecorder()
			//

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/create-recipe"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetRecipes(t *testing.T) {
	recipes := make([]db.Recipe, 5)
	for i := 0; i < 5; i++ {
		recipes[i] = RandomRecipe()
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
				arg := db.GetRecipesParams{
					Limit:  5,
					Offset: 5,
				}
				store.EXPECT().GetRecipes(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(recipes, nil)
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
				arg := db.GetRecipesParams{
					Limit:  5,
					Offset: 5,
				}
				store.EXPECT().GetRecipes(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.Recipe{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{
				"limit":  -5,
				"offset": -5,
			},
			buildStabs: func(store *mockdb.MockStore) {
				arg := db.GetRecipesParams{
					Limit:  5,
					Offset: 5,
				}
				store.EXPECT().GetRecipes(gomock.Any(), gomock.Eq(arg)).
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

			url := "/recipes"
			request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func RandomRecipeByCategoryRow() (db.GetRecipeByCategoryRow, db.Recipe, db.RecipeCategory) {
	recipe := RandomRecipe()
	category := RandomCategory()

	return db.GetRecipeByCategoryRow{
		ID:               recipe.ID,
		Name:             recipe.Name,
		Description:      recipe.Description,
		Image:            recipe.Image,
		Active:           recipe.Active,
		Time:             recipe.Time,
		Url:              recipe.Url,
		Servings:         recipe.Servings,
		CreatedAt:        recipe.CreatedAt,
		ID_2:             sql.NullString{String: uuid.New().String(), Valid: true},
		RecipeID:         sql.NullString{String: recipe.ID, Valid: true},
		RecipeCategoryID: sql.NullString{String: category.ID, Valid: true},
	}, recipe, category
}

func TestGetRecipesByCategory(t *testing.T) {

	row, recipe, category := RandomRecipeByCategoryRow()
	testCases := []struct {
		name          string
		catid         string
		buildStabs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			catid: category.ID,
			buildStabs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetRecipeByCategory(gomock.Any(), category.ID).
					Times(1).
					Return([]db.GetRecipeByCategoryRow{row}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				checkBodyMatch(t, recorder.Body, recipe)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStabs(store)

			server := newTestServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/recipes/%s", tc.catid)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func checkBodyMatch(t *testing.T, body *bytes.Buffer, recipe db.Recipe) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotRecipes []db.GetRecipeByCategoryRow
	err = json.Unmarshal(data, &gotRecipes)
	require.NoError(t, err)

	for _, recipe := range gotRecipes {
		require.Equal(t, recipe.ID, recipe.ID)
		require.Equal(t, recipe.Name, recipe.Name)
		require.Equal(t, recipe.Description, recipe.Description)
		require.Equal(t, recipe.Servings, recipe.Servings)
		require.Equal(t, recipe.Time, recipe.Time)
	}
}
