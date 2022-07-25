package simple

import (
	"auth/internal/auth"
	"auth/internal/cache"
	"auth/internal/model/principal"
	"auth/internal/repo"
	"auth/pkg/logger"
	"auth/pkg/token"
	"context"
)

var _ auth.AuthenticationService = (*Auth)(nil)

type Auth struct {
	authority token.Authority
	repo      repo.PrincipalRepo
	cache     cache.Cache
	logger    logger.Logger
}

func New(authority token.Authority, repo repo.PrincipalRepo, cache cache.Cache, logger logger.Logger) *Auth {
	return &Auth{
		authority: authority,
		repo:      repo,
		cache:     cache,
		logger:    logger,
	}
}

func (a *Auth) Signup(ctx context.Context, request auth.SignUpRequest) (*auth.JWTokens, error) {

	a.logger.Debug("Auth.Signup: email: %s", request.Email)

	p := principal.Principal{
		Email:    request.Email,
		Password: request.Password,
		Role:     request.Role,
	}

	if err := p.Validate(); err != nil {
		a.logger.Error("Principal.Validate: email: %q, err: %s", request.Email, err)
		return nil, err
	}
	if err := p.HashPassword(); err != nil {
		a.logger.Error("Principal.HashPassword: email: %q, err: %s", request.Email, err)
		return nil, err
	}

	p.Sanitize()

	if err := a.repo.CreatePrincipal(ctx, p); err != nil {
		a.logger.Error("Principal.CreatePrincipal: email: %q, err: %s", request.Email, err)
		return nil, err
	}

	a.logger.Info("email: %q, database record created", p.Email)

	if err := a.cache.Set(p.Email, p); err != nil {
		a.logger.Error("Cache.Set: email: %q, err(ignored): %s", request.Email, err)
	}

	select {
	default:
	case <-ctx.Done():
		a.logger.Info("ctx.Done: err: %s", ctx.Err())
		return nil, ctx.Err()
	}

	acc, ref, err := a.authority.Issue(p.Email, p.Role)
	if err != nil {
		a.logger.Error("Authority.Issue: email: %s, role: %s, err: %s", p.Email, p.Role, err)
		return nil, auth.ErrInternal
	}

	a.logger.Debug("Auth.Signup: singup succeeded, tokens are issued, email: %s", request.Email)

	return &auth.JWTokens{Access: acc, Refresh: ref}, nil
}

func (a *Auth) Login(ctx context.Context, request auth.LoginRequest) (*auth.JWTokens, error) {

	a.logger.Debug("Auth.Login: email: %s", request.Email)

	if err := (&principal.Principal{
		Email:    request.Email,
		Password: request.Password,
	}).Validate(); err != nil {
		a.logger.Error("Principal.Validate: err: %s", err)
		return nil, auth.ErrIncorrectLoginOrPassword
	}

	p, err := a.cache.Get(request.Email)
	if err != nil {
		a.logger.Error("Cache.Get: key: %s, err: %s", request.Email, err)

		rp, err := a.repo.GetPrincipal(ctx, request.Email)
		if err != nil {
			a.logger.Error("Repo.GetPrincipal: email: %s, err: %s", request.Email, err)
			return nil, auth.ErrIncorrectLoginOrPassword
		}

		a.logger.Debug("email %s is found in repo", request.Email)

		if err := a.cache.Set(request.Email, p); err != nil {
			a.logger.Error("Cache.Set: key %s, err(ignored): %s", request.Email, err)
		}

		p = *rp

	} else {
		a.logger.Debug("email %s is found in cache", request.Email)
	}

	if !p.PasswordMatches(request.Password) {
		a.logger.Error("Principal.PasswordMatches: email: %s, password doesn't match", p.Email)
		return nil, auth.ErrIncorrectLoginOrPassword
	}

	select {
	default:
	case <-ctx.Done():
		a.logger.Info("ctx.Done: err: %s", ctx.Err())
		return nil, ctx.Err()
	}

	acc, ref, err := a.authority.Issue(p.Email, p.Role)
	if err != nil {
		a.logger.Error("Authority.Issue: email: %s, role: %s, err: %s", p.Email, p.Role, err)
		return nil, auth.ErrInternal
	}

	a.logger.Debug("Auth.Login: login succeeded, tokens are issued, email: %q", p.Email)

	return &auth.JWTokens{Access: acc, Refresh: ref}, nil
}

func (a *Auth) Verify(ctx context.Context, accessToken string) error {

	a.logger.Debug("Auth.Verify: accessToken: %s", accessToken)

	select {
	default:
	case <-ctx.Done():
		return ctx.Err()
	}

	if err := a.authority.Verify(accessToken); err != nil {
		a.logger.Error("Authority.Revoke: err: %s, access token: %s", err, accessToken)
		return err
	}

	a.logger.Debug("Auth.Verify: token verified: %s", accessToken)

	return nil
}

func (a *Auth) Refresh(ctx context.Context, refreshToken string) (*auth.JWTokens, error) {

	a.logger.Debug("Auth.Refresh: refreshToken: %s", refreshToken)

	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	at, rt, err := a.authority.Refresh(refreshToken)
	if err != nil {
		a.logger.Error("Authority.Refresh: err: %s, refreshToken: %s", err, refreshToken)
		return nil, err
	}

	a.logger.Debug("Auth.Refresh: tokens refreshed: access: %s, refresh: %s", at, rt)

	return &auth.JWTokens{Access: at, Refresh: rt}, nil
}

func (a *Auth) Revoke(ctx context.Context, accessToken string) error {

	a.logger.Debug("Auth.Revoke: accessToken: %s", accessToken)

	select {
	default:
	case <-ctx.Done():
		return ctx.Err()
	}

	err := a.authority.Revoke(accessToken)
	if err != nil {
		a.logger.Error("Authority.Revoke: err: %s, access token: %s", err, accessToken)
		return err
	}

	a.logger.Debug("Auth.Revoke: token revoked: accessToken: %s", accessToken)

	return nil
}
