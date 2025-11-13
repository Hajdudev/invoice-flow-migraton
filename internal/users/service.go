package users

import (
	"context"
	"fmt"

	repo "github.com/Hajdudev/invoice-flow/internal/adapters/postgresql/sqlc"
)

type svc struct {
	repo *repo.Queries
}

func NewService(repo *repo.Queries) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) RegisterUser(ctx context.Context) error {
	fmt.Println("runned")
	return nil
}
