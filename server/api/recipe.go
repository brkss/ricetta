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
