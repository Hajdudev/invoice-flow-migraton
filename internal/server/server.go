package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	database "github.com/Hajdudev/invoice-flow/internal/adapters/postgresql"
	"github.com/Hajdudev/invoice-flow/internal/auth"
)

type Server struct {
	port    int
	db      database.Service
	jwtAuth auth.JWTAuthenticator
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	database := database.New()
	jwtAuth := auth.NewJWTAuthenticator()
	auth.NewOauth()

	NewServer := &Server{
		port:    port,
		db:      database,
		jwtAuth: *jwtAuth,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
