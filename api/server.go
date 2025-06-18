package api

import (
	db "github.com/dohoanggiahuy317/ACH-transactions-Microservice-app/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for the banking service.
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(store *db.Store) *Server {

	// Initialize the server with the provided store and a new Gin router.
	server := &Server{store:  store}
	router := gin.Default()

	// Define the route for creating an account.
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	server.router = router
	return server
}

// Start runs the HTTP server on the specified address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse formats an error message as a JSON response.
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}