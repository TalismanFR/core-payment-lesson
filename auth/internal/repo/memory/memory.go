package memory

import (
	"auth/internal/model/principal"
	"auth/internal/repo"
	"context"
	"fmt"
)

var _ repo.PrincipalRepo = (*RepoInMemory)(nil)

type RepoInMemory struct {
	m map[string]principal.Principal
}

func New() *RepoInMemory {
	return &RepoInMemory{
		map[string]principal.Principal{},
	}
}

func (r *RepoInMemory) CreatePrincipal(ctx context.Context, p principal.Principal) error {
	if _, ok := r.m[p.Email]; ok {
		return fmt.Errorf("email %s already exists", p.Email)
	}
	r.m[p.Email] = p
	return nil
}

func (r *RepoInMemory) GetPrincipal(ctx context.Context, email string) (*principal.Principal, error) {
	v, ok := r.m[email]
	if !ok {
		return nil, fmt.Errorf("email %s doesn't exists", email)
	}
	return &v, nil
}

func (r *RepoInMemory) UpdatePrincipal(ctx context.Context, p principal.Principal) error {
	_, ok := r.m[p.Email]
	if !ok {
		return fmt.Errorf("email %s doesn't exists", p.Email)
	}
	r.m[p.Email] = p
	return nil
}
