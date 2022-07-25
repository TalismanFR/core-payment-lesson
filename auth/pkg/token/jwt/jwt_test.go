package jwt

import (
	"github.com/coocood/freecache"
	"testing"
	"time"
)

var (
	validAuthority JWTokenAuthority = JWTokenAuthority{
		accessSigningKey:  "123",
		refreshSigningKey: "456",
		accessTtl:         5 * time.Second,
		refreshTtl:        10 * time.Second,
	}
)

func TestJWTokenAuthority_DoubleRefresh(t *testing.T) {
	j := &JWTokenAuthority{
		accessSigningKey:  "123",
		refreshSigningKey: "456",
		accessTtl:         3 * time.Second,
		refreshTtl:        6 * time.Second,
		refreshIDs:        freecache.NewCache(1),
	}

	a1, r1, err := j.Issue("tim", "user")
	t.Log(err)
	a2, r2, err := j.Refresh(r1)
	t.Log(err)
	a3, r3, err := j.Refresh(r1)
	t.Log(err)
	if err == nil {
		t.Fatal("after initial refresh call, father calls with the same token should return err")
	}
	_, _, _, _, _, _ = a1, a2, a3, r1, r2, r3
}

// After revoking, access token is still functioning
func TestJWTokenAuthority_RevokeVerifyRefresh(t *testing.T) {
	j := &JWTokenAuthority{
		accessSigningKey:  "123",
		refreshSigningKey: "456",
		accessTtl:         3 * time.Second,
		refreshTtl:        6 * time.Second,
		refreshIDs:        freecache.NewCache(1),
	}

	a, r, _ := j.Issue("tim", "user")
	err := j.Revoke(a)
	t.Log(err)
	err = j.Verify(a)
	t.Log(err)
	if err != nil {
		t.Fatalf("revoked access token should NOT fail to verify")
	}
	_, _, err = j.Refresh(r)
	t.Log(err)
	if err == nil {
		t.Fatalf("revoked refresh token should fail to refresh")
	}
}

// After revoking access token, refresh token is invalidated
func TestJWTokenAuthority_RevokeRefresh(t *testing.T) {
	j := &JWTokenAuthority{
		accessSigningKey:  "123",
		refreshSigningKey: "456",
		accessTtl:         3 * time.Second,
		refreshTtl:        6 * time.Second,
		refreshIDs:        freecache.NewCache(1),
	}

	a, r, _ := j.Issue("tim", "user")
	err := j.Revoke(a)
	t.Log(err)
	if err != nil {
		t.Fatalf("revoke error: %s", err)
	}
	_, _, err = j.Refresh(r)
	t.Log(err)
	if err == nil {
		t.Fatalf("revoked refresh token should fail to refresh")
	}
}

// After revoking, access token is still functioning
func TestJWTokenAuthority_Refresh(t *testing.T) {
	j := &JWTokenAuthority{
		accessSigningKey:  "123",
		refreshSigningKey: "456",
		accessTtl:         3 * time.Second,
		refreshTtl:        6 * time.Second,
		refreshIDs:        freecache.NewCache(1),
	}

	a, r, err := j.Issue("tim", "user")
	t.Log(err)
	a, _, err = j.Refresh(a)
	t.Log(err)
	_, _, err = j.Refresh(r)
	t.Log(err)
}

// After revoking, access token is still functioning
func TestJWTokenAuthority_VerifyExpired(t *testing.T) {

	ttl := 1 * time.Second

	j := &JWTokenAuthority{
		accessSigningKey:  "123",
		refreshSigningKey: "456",
		accessTtl:         ttl,
		refreshTtl:        5 * time.Second,
		refreshIDs:        freecache.NewCache(1),
	}

	a, _, err := j.Issue("tim", "user")
	t.Log(err)
	time.Sleep(ttl)
	err = j.Verify(a)
	t.Log(err)
	if err == nil {
		t.Fatalf("expired access token should fail to verify")
	}
}

// After revoking, access token is still functioning
func TestJWTokenAuthority_RefreshExpired(t *testing.T) {

	ttl := 1 * time.Second

	j := &JWTokenAuthority{
		accessSigningKey:  "123",
		refreshSigningKey: "456",
		accessTtl:         ttl,
		refreshTtl:        ttl,
		refreshIDs:        freecache.NewCache(1),
	}

	_, r, err := j.Issue("tim", "user")
	t.Log(err)
	time.Sleep(ttl)
	_, _, err = j.Refresh(r)
	t.Log(err)
	if err == nil {
		t.Fatalf("expired refresh token should fail to verify")
	}
}

func TestJWTokenAuthority_Issue(t *testing.T) {
	tests := []struct {
		name    string
		subject string
		role    string
		j       func() *JWTokenAuthority
		isErr   bool
	}{
		{
			name:    "valid1",
			subject: "user1",
			role:    "role1",
			j: func() *JWTokenAuthority {
				return &JWTokenAuthority{
					accessSigningKey:  "123",
					refreshSigningKey: "456",
					accessTtl:         3 * time.Second,
					refreshTtl:        6 * time.Second,
				}
			},
			isErr: false,
		},
		{
			name:    "valid2",
			subject: "",
			role:    "",
			j: func() *JWTokenAuthority {
				return &JWTokenAuthority{
					accessSigningKey:  "",
					refreshSigningKey: "",
					accessTtl:         0 * time.Second,
					refreshTtl:        0 * time.Second,
				}
			},
			isErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a, r, err := tc.j().Issue(tc.subject, tc.role)
			t.Logf("access: %q, refresh: %q", a, r)
			if !tc.isErr && err != nil {
				t.Fatalf("error isn't nil:\nER: nil\nAR: %s", err)
			}
			if tc.isErr && err == nil {
				t.Fatalf("error is nil:\nER: not nil error\nAR: nil")
			}
		})
	}
}

func TestJWTokenAuthority_VerifyTableTests(t *testing.T) {
	tests := []struct {
		name        string
		j           func() *JWTokenAuthority
		accessToken func() string
		isErr       bool
	}{
		{
			name: "invalidNumberOfSegments",
			j: func() *JWTokenAuthority {
				return &JWTokenAuthority{
					accessSigningKey:  "123",
					refreshSigningKey: "456",
					accessTtl:         3 * time.Second,
					refreshTtl:        6 * time.Second,
				}
			},
			accessToken: func() string {
				return "123"
			},
			isErr: true,
		},
		{
			name: "expiredToken",
			j: func() *JWTokenAuthority {
				return &JWTokenAuthority{
					accessSigningKey: "123",
				}
			},
			accessToken: func() string {
				return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE1MTYyMzkwMjN9.YS3BRO8gPMNR3OMCyv0YFv_oow-CLQU3XTFB96j00gw"
			},
			isErr: true,
		},
		{
			name: "valid_expiresIn2255",
			j: func() *JWTokenAuthority {
				return &JWTokenAuthority{
					accessSigningKey: "123",
				}
			},
			accessToken: func() string {
				return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjkwMTY1MzkwMjN9.8o-Wh3lJdXgdhrej8mRmUCbUcz-4m9lDqwJVblob_n0"
			},
			isErr: false,
		},
		{
			name: "invalidAccessTokenSignature",
			j: func() *JWTokenAuthority {
				return &JWTokenAuthority{
					accessSigningKey: "123",
				}
			},
			accessToken: func() string {
				return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE1MTYyMzkwMjR9.AFHL3WoZhnGOxnu7-DlBQMcj7alAiVDI0FfE1eyFI-o"
			},

			isErr: true,
		},
		{
			name: "tokenChanged",
			j: func() *JWTokenAuthority {
				return &JWTokenAuthority{
					accessSigningKey: "123",
				}
			},
			accessToken: func() string {
				return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjkwMTYyMzkwMjR9.OOEi_yPjM_17867xFHUNj8VQgyD6h3xjY8C1XtMzHw0"
			},

			isErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.j().Verify(tc.accessToken())
			if !tc.isErr && err != nil {
				t.Fatalf("error isn't nil:\nER: nil\nAR: %s", err)
			}
			if tc.isErr && err == nil {
				t.Fatalf("error is nil:\nER: not nil error\nAR: nil")
			}
		})
	}
}
