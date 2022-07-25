package repo

import (
	"auth/internal/model/principal"
	"context"
)

//go:generate mockgen -source=repo.go -destination=../mocks/repo_pg_mock.go -package=mocks
type PrincipalRepo interface {
	CreatePrincipal(ctx context.Context, p principal.Principal) error
	GetPrincipal(ctx context.Context, email string) (*principal.Principal, error)
	UpdatePrincipal(ctx context.Context, p principal.Principal) error
}
