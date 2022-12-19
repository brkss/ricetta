package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
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
