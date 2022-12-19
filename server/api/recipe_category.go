package api

import (
	"net/http"

	db "github.com/brkss/vanillefraise2/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRecipeCategoryRequest struct {
	Title string `json:"title" binding:"required"`
	Image string `json:"image" binding:"required"`
}

type GetCategorieRequest struct {
	Limit  int32 `json:"limit" binding:"required,min=1"`
	Offset int32 `json:"offset" binding:"required,min=1"`
}

func (server *Server) CreateRecipeCategoryAPI(ctx *gin.Context) {
	var req CreateRecipeCategoryRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreateRecipeCategoryParams{
		ID:    uuid.New().String(),
		Title: req.Title,
		Image: req.Image,
	}
	category, err := server.store.CreateRecipeCategory(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, category)
}

func (server *Server) GetCategoriesAPI(ctx *gin.Context) {
	var req GetCategorieRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.GetCategoriesParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	}
	categories, err := server.store.GetCategories(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, categories)
}
