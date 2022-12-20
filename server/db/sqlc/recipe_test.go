package db

import (
	"context"
	"testing"

	"github.com/brkss/vanillefraise2/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func RandomRecipeCategory(t *testing.T) RecipeCategory {

	arg := CreateRecipeCategoryParams{
		ID:    uuid.New().String(),
		Title: utils.RandomName(),
		Image: utils.RandomName(),
	}

	category, err := testQueries.CreateRecipeCategory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, category)

	require.Equal(t, category.ID, arg.ID)
	require.Equal(t, category.Title, arg.Title)
	require.Equal(t, category.Image, arg.Image)
	require.True(t, category.Active.Bool)

	return category

}

func RandomRecipe(t *testing.T) Recipe {
	arg := CreateRecipeParams{
		ID:          uuid.New().String(),
		Name:        utils.RandomName(),
		Description: utils.RandomString(100),
		Image:       utils.RandomName(),
		Url:         utils.RandomString(33),
		Time:        utils.RandomString(5),
		Servings:    10,
		//Description: utils.RandomString(100),
	}
	recipe, err := testQueries.CreateRecipe(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, recipe)

	require.Equal(t, recipe.ID, arg.ID)
	require.Equal(t, recipe.Name, arg.Name)
	require.Equal(t, recipe.Description, arg.Description)
	require.Equal(t, recipe.Image, arg.Image)
	require.Equal(t, recipe.Servings, arg.Servings)
	require.Equal(t, recipe.Url, arg.Url)
	require.Equal(t, recipe.Time, arg.Time)

	return recipe
}

func AssignRecipeToCategory(t *testing.T, recipe Recipe, category RecipeCategory) {

	arg := AssignRecipeToCategoryParams{
		ID:               uuid.New().String(),
		RecipeID:         recipe.ID,
		RecipeCategoryID: category.ID,
	}
	recipeToCategory, err := testQueries.AssignRecipeToCategory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, recipeToCategory)

	require.Equal(t, recipeToCategory.ID, arg.ID)
	require.Equal(t, recipeToCategory.RecipeID, arg.RecipeID)
	require.Equal(t, recipeToCategory.RecipeCategoryID, arg.RecipeCategoryID)
}

func TestCreateRecipeCategory(t *testing.T) {
	RandomRecipeCategory(t)
}

func TestCreateRecipe(t *testing.T) {
	RandomRecipe(t)
}

func TestRecipeToCategory(t *testing.T) {
	recipe := RandomRecipe(t)
	category := RandomRecipeCategory(t)
	AssignRecipeToCategory(t, recipe, category)
}

func TestGetRecipes(t *testing.T) {
	n := 10

	for i := 0; i < n; i++ {
		RandomRecipe(t)
	}
	arg := GetRecipesParams{
		Limit:  5,
		Offset: 5,
	}
	recipes, err := testQueries.GetRecipes(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, recipes, 5)

	for _, recipe := range recipes {
		require.NotEmpty(t, recipe)
	}
}

func TestGetRecipeCategories(t *testing.T) {

	n := 10
	for i := 0; i < n; i++ {
		RandomRecipeCategory(t)
	}

	arg := GetCategoriesParams{
		Limit:  5,
		Offset: 5,
	}
	categories, err := testQueries.GetCategories(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, categories, 5)

	for _, category := range categories {
		require.NotEmpty(t, category)
	}
}

func TestGetRecipesByCategory(t *testing.T) {
	recipe := RandomRecipe(t)
	category := RandomRecipeCategory(t)
	AssignRecipeToCategory(t, recipe, category)

	recipes, err := testQueries.GetRecipeByCategory(context.Background(), category.ID)
	require.NoError(t, err)
	require.Len(t, recipes, 1)

	for _, recipe := range recipes {
		require.NotEmpty(t, recipe)
	}
}
