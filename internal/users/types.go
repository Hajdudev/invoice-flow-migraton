package users

import "context"

type Service interface {
	RegisterUser(ctx context.Context) error
}
