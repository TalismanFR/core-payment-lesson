package auth

import (
	"context"
	"errors"
)

var (
	ErrEmailInvalidOrTaken      = errors.New("email invalid or taken")
	ErrWeakPassword             = errors.New("password is too weak")
	ErrIncorrectLoginOrPassword = errors.New("incorrect login or password")
	ErrInternal                 = errors.New("internal error")
)

//go:generate mockgen -source=auth.go -destination=../mocks/auth_mock.go -package=mocks
type (
	AuthenticationService interface {
		Signup(ctx context.Context, r SignUpRequest) (*JWTokens, error)
		Login(ctx context.Context, r LoginRequest) (*JWTokens, error)
		Refresh(ctx context.Context, refreshToken string) (*JWTokens, error)
		Verify(ctx context.Context, accessToken string) error
		Revoke(ctx context.Context, accessToken string) error
	}

	JWTokens struct {
		Access  string
		Refresh string
	}

	SignUpRequest struct {
		Email    string
		Password string
		Role     string
	}

	LoginRequest struct {
		Email    string
		Password string
	}
)
