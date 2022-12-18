package api

import (
	"net/http"

	db "github.com/brkss/vanillefraise2/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

// testing !
func (server *Server) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "pong")
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// add routes here !

	router.POST("/register", server.createUser)
	router.GET("/ping", server.Ping)
	server.router = router
	return server
}

// Start new HTTP server and listen for requests
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
