package users

import (
	"context"
	"fmt"

	repo "github.com/Hajdudev/invoice-flow/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

type svc struct {
	repo *repo.Queries
	db   *pgx.Conn
}

func NewService(repo *repo.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

func (s *svc) RegisterUser(ctx context.Context) error {
	fmt.Println("runned")
	return nil
}
