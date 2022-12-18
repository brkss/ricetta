package api

import (
	"net/http"

	db "github.com/brkss/vanillefraise2/db/sqlc"
	"github.com/brkss/vanillefraise2/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// CreateUserRequest is wjat we're expecring from register user post request
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// createUser handler for reguster post request
func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		ID:       uuid.New().String(),
		Name:     req.Name,
		Email:    req.Email,
		Username: req.Username,
		Password: hash,
	}
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, user)
}
