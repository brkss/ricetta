package api

import (
	"net/http"

	db "github.com/brkss/vanillefraise2/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRecipeRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Image       string `json:"image" binding:"required"`
	Url         string `json:"url" binding:"required"`
	Time        string `json:"time" binding:"required"`
	Servings    int32  `json:"servings" binding:"required,min=1"`
	CategoryId  string `json:"category_id" binding:"required"`
}

type GetRecipeRequest struct {
	Limit  int32 `json:"limit" binding:"required,min=1"`
	Offset int32 `json:"offset" binding:"required,min=1"`
}

type GetRecipesByCategoryRequest struct {
	CategoryId string `uri:"catid" binding:"required"`
}

func (server *Server) CreateRecipeAPI(ctx *gin.Context) {
	var req CreateRecipeRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.CreateRecipeParams{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		Url:         req.Url,
		Time:        req.Time,
		Servings:    req.Servings,
	}

	recipe, err := server.store.CreateRecipe(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	assingRecipeCatgoryArg := db.AssignRecipeToCategoryParams{
		ID:               uuid.New().String(),
		RecipeID:         recipe.ID,
		RecipeCategoryID: req.CategoryId,
	}

	_, err = server.store.AssignRecipeToCategory(ctx, assingRecipeCatgoryArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, recipe)
}

func (server *Server) GetRecipes(ctx *gin.Context) {
	var req GetRecipeRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetRecipesParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	}
	recipes, err := server.store.GetRecipes(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, recipes)
	return
}

func (server *Server) GetRecipesByCategory(ctx *gin.Context) {
	var req GetRecipesByCategoryRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	recipes, err := server.store.GetRecipeByCategory(ctx, req.CategoryId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, recipes)
}
