package simple

import (
	"auth/internal/auth"
	"auth/internal/mocks"
	"auth/internal/model/principal"
	"auth/pkg/logger"
	pkgmocks "auth/pkg/mocks"
	"bytes"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

type fakeHasher struct {
	generate func([]byte) ([]byte, error)
	compare  func([]byte, []byte) error
}

// simpleFakeHasher returns password as hash
var simpleFakeHasher = fakeHasher{
	generate: func(password []byte) ([]byte, error) {
		return password, nil
	},
	compare: func(hash []byte, password []byte) error {
		if bytes.Equal(hash, password) {
			return nil
		}
		return errors.New("hash error")
	},
}

func (s fakeHasher) GenerateFromPassword(password []byte) ([]byte, error) {
	return s.generate(password)
}

func (s fakeHasher) CompareHashAndPassword(hash []byte, password []byte) error {
	return s.compare(hash, password)
}

func TestAuth_SignupLogin(t *testing.T) {
	defer temporarySetHasher(simpleFakeHasher)()

	ctx := context.TODO()

	expectedSignupTokens := newTokens("signup")
	expectedLoginTokens := newTokens("login")

	p := validPrincipal()
	signupRequest := auth.SignUpRequest{
		Email:    p.Email,
		Password: p.Password,
		Role:     p.Role,
	}

	loginRequest := auth.LoginRequest{
		Email:    p.Email,
		Password: p.Password,
	}
	p.Password = ""

	ctrl := gomock.NewController(t)

	c := mocks.NewMockCache(ctrl)
	c.EXPECT().Set(p.Email, p).Return(nil).Times(1)
	c.EXPECT().Get(p.Email).Return(p, nil).Times(1)

	r := mocks.NewMockPrincipalRepo(ctrl)
	r.EXPECT().CreatePrincipal(ctx, p).Return(nil).Times(1)

	a := pkgmocks.NewMockAuthority(ctrl)
	gomock.InOrder(
		a.EXPECT().Issue(p.Email, p.Role).Return(expectedSignupTokens.Access, expectedSignupTokens.Refresh, nil).Times(1),
		a.EXPECT().Issue(p.Email, p.Role).Return(expectedLoginTokens.Access, expectedLoginTokens.Refresh, nil).Times(1),
	)

	authService := New(a, r, c, logger.New(logger.Debug))

	actualSignUpTokens, err := authService.Signup(ctx, signupRequest)
	require.Nil(t, err)
	require.Equal(t, expectedSignupTokens, actualSignUpTokens)

	actualLoginTokens, err := authService.Login(ctx, loginRequest)
	require.Nil(t, err)
	require.Equal(t, expectedLoginTokens, actualLoginTokens)
}

func TestAuth_SignupVerify(t *testing.T) {
	defer temporarySetHasher(simpleFakeHasher)()

	ctx := context.TODO()
	expectedSignupTokens := newTokens("1")
	p := validPrincipal()
	signupRequest := auth.SignUpRequest{
		Email:    p.Email,
		Password: p.Password,
		Role:     p.Role,
	}
	p.Password = ""

	ctrl := gomock.NewController(t)

	c := mocks.NewMockCache(ctrl)
	c.EXPECT().Set(p.Email, p).Return(nil).Times(1)

	r := mocks.NewMockPrincipalRepo(ctrl)
	r.EXPECT().CreatePrincipal(ctx, p).Return(nil).Times(1)

	a := pkgmocks.NewMockAuthority(ctrl)
	gomock.InOrder(
		a.EXPECT().Issue(p.Email, p.Role).Return(expectedSignupTokens.Access, expectedSignupTokens.Refresh, nil).Times(1),
		a.EXPECT().Verify(expectedSignupTokens.Access).Return(nil).Times(1),
	)

	authService := New(a, r, c, logger.New(logger.Debug))

	actualSignUpTokens, err := authService.Signup(ctx, signupRequest)
	require.Nil(t, err)
	require.Equal(t, expectedSignupTokens, actualSignUpTokens)

	err = authService.Verify(ctx, actualSignUpTokens.Access)
	require.Nil(t, err)
}

func TestAuth_SingleSignup(t *testing.T) {

	defer temporarySetHasher(simpleFakeHasher)()

	ctx := context.TODO()
	expectedTokens := newTokens("1")
	p := validPrincipal()
	signupRequest := auth.SignUpRequest{
		Email:    p.Email,
		Password: p.Password,
		Role:     p.Role,
	}
	p.Password = ""

	ctrl := gomock.NewController(t)

	c := mocks.NewMockCache(ctrl)
	c.EXPECT().Set(p.Email, p).Return(nil).Times(1)

	r := mocks.NewMockPrincipalRepo(ctrl)
	r.EXPECT().CreatePrincipal(ctx, p).Return(nil).Times(1)

	a := pkgmocks.NewMockAuthority(ctrl)
	a.EXPECT().Issue(p.Email, p.Role).Return(expectedTokens.Access, expectedTokens.Refresh, nil).Times(1)

	authService := New(a, r, c, logger.New(logger.Debug))
	actualTokens, err := authService.Signup(ctx, signupRequest)

	require.Nil(t, err)
	require.Equal(t, expectedTokens, actualTokens)
}

func TestAuth_SignupWithTheSameCredentials(t *testing.T) {
	defer temporarySetHasher(simpleFakeHasher)()

	ctx := context.TODO()
	expectedTokens := newTokens("1")
	p := validPrincipal()
	signupRequest := auth.SignUpRequest{
		Email:    p.Email,
		Password: p.Password,
		Role:     p.Role,
	}
	p.Password = ""

	ctrl := gomock.NewController(t)

	c := mocks.NewMockCache(ctrl)
	c.EXPECT().Set(p.Email, p).Return(nil).Times(1)

	r := mocks.NewMockPrincipalRepo(ctrl)
	gomock.InOrder(
		r.EXPECT().CreatePrincipal(ctx, p).Return(nil).Times(1),
		r.EXPECT().CreatePrincipal(ctx, p).Return(errors.New("email taken")).Times(1),
	)

	a := pkgmocks.NewMockAuthority(ctrl)
	a.EXPECT().Issue(p.Email, p.Role).Return(expectedTokens.Access, expectedTokens.Refresh, nil).Times(1)

	authService := New(a, r, c, logger.New(logger.Debug))

	actualTokens, err := authService.Signup(ctx, signupRequest)
	require.Nil(t, err)
	require.Equal(t, expectedTokens, actualTokens)

	actualTokens, err = authService.Signup(ctx, signupRequest)
	require.NotNil(t, err, "second signup request with the same credentials should return non-nil error, got nil")
}

func validPrincipal() principal.Principal {
	email := "tim@ya.ru"
	role := "user"
	password := "password123"
	hashedPassword := []byte(password)

	return principal.Principal{
		Email:          email,
		Password:       password,
		Role:           role,
		HashedPassword: hashedPassword,
	}

}

func newTokens(suffix string) *auth.JWTokens {
	return &auth.JWTokens{
		Access:  "access_token_" + suffix,
		Refresh: "refresh_token_" + suffix,
	}
}

func temporarySetHasher(h principal.Hasher) func() {
	oldHasher := principal.GetHasher()
	principal.SetHasher(h)
	return func() {
		principal.SetHasher(oldHasher)
	}
}
