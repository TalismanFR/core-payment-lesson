package jwt

import (
	"auth/internal/config"
	"auth/pkg/token"
	"errors"
	"fmt"
	"github.com/coocood/freecache"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

var _ token.Authority = (*JWTokenAuthority)(nil)

var (
	ErrRevokedOrUsedRefreshToken = errors.New("token was revoked or refreshed before")
)

type Claims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

type JWTokenAuthority struct {
	accessTtl  time.Duration
	refreshTtl time.Duration

	// accessSigningKey and refreshSigningKey should not be the same
	accessSigningKey  string
	refreshSigningKey string

	refreshIDs *freecache.Cache
}

func NewJWTokenAuthority(cfg *config.JWTConfig) (*JWTokenAuthority, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}
	return &JWTokenAuthority{
		accessTtl:         cfg.AccessTokenTTL,
		refreshTtl:        cfg.RefreshTokenTTL,
		accessSigningKey:  cfg.AccessSigningKey,
		refreshSigningKey: cfg.RefreshSigningKey,
		refreshIDs:        freecache.NewCache(1),
	}, nil
}

func (j *JWTokenAuthority) Issue(subject, role string) (access, refresh string, err error) {
	return j.issue(subject, role)
}

// Refresh returns a new pair of tokens or an error
func (j *JWTokenAuthority) Refresh(refresh string) (string, string, error) {

	claims, err := j.parse(refresh, j.refreshSigningKey)
	if err != nil {
		return "", "", err
	}

	if _, err = j.refreshIDs.Get([]byte(claims.ID)); err != nil {
		return "", "", ErrRevokedOrUsedRefreshToken
	}

	a, r, err := j.issue(claims.Subject, claims.Role)
	if err != nil {
		return "", "", err
	}

	j.refreshIDs.Del([]byte(claims.ID))
	return a, r, nil
}

// Verify access token
func (j *JWTokenAuthority) Verify(access string) error {
	_, err := j.parse(access, j.accessSigningKey)
	return err
}

// Revoke refresh token associated with access token
func (j *JWTokenAuthority) Revoke(access string) error {
	claims, err := j.parse(access, j.accessSigningKey)
	if err != nil {
		return err
	}
	j.refreshIDs.Del([]byte(claims.ID))
	return nil
}

func (j *JWTokenAuthority) issue(subject, role string) (access, refresh string, err error) {
	id := uuid.New().String()

	if access, refresh, err = j.newTokenPair(subject, role, id); err != nil {
		return
	}
	if err = j.refreshIDs.Set([]byte(id), []byte{}, int(j.refreshTtl.Seconds())); err != nil {
		return
	}

	return
}

// newToken issues one jwt token with subject and time-to-live and sign it with key
func (j *JWTokenAuthority) newToken(subject, role, id string, ttl time.Duration, key string) (string, error) {

	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			ID:        id,
		},
		role,
	},
	)

	return tkn.SignedString([]byte(key))
}

// newTokenPair issues access and refresh tokens with subject
func (j *JWTokenAuthority) newTokenPair(subject, role, id string) (access, refresh string, err error) {

	if access, err = j.newToken(subject, role, id, j.accessTtl, j.accessSigningKey); err != nil {
		return
	}
	if refresh, err = j.newToken(subject, role, id, j.refreshTtl, j.refreshSigningKey); err != nil {
		return
	}

	return
}

// parse token with key, check validity and return parsed subject and error
func (j *JWTokenAuthority) parse(token string, key string) (*Claims, error) {

	if token == "" {
		return nil, fmt.Errorf("empty token")
	}

	var claims Claims

	parsedToken, err :=
		jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()})).
			ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) { return []byte(key), nil })

	if err != nil || !parsedToken.Valid {
		return nil, err
	}

	return &claims, nil
}

func validateConfig(cfg *config.JWTConfig) error {
	if cfg.AccessSigningKey == "" {
		return fmt.Errorf("empty access token")
	}
	if cfg.RefreshSigningKey == "" {
		return fmt.Errorf("empty refresh token")
	}
	if cfg.AccessTokenTTL == cfg.RefreshTokenTTL {
		return fmt.Errorf("access and refresh signing key cannot be the same")
	}
	return nil
}
