package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	database "github.com/Hajdudev/invoice-flow/internal/adapters/postgresql"
	"github.com/jackc/pgx/v5"
)

type Server struct {
	port int
	db   database.Service
	conn *pgx.Conn
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	database := database.New()
	dbConn, err := database.GetDB(context.Background())

	defer dbConn.Close(context.Background())

	if err != nil {
		fmt.Errorf("error %v", err)
	}
	NewServer := &Server{
		port: port,
		db:   database,
		conn: dbConn,
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
