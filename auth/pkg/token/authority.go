package token

//go:generate mockgen -source=authority.go -destination=../mocks/authority_mock.go -package=mocks
type Authority interface {

	// Issue issues access and refresh token for the subject.
	Issue(subject, role string) (access, refresh string, err error)

	// Refresh checks refresh token validity and issues new access and refresh tokens.
	//
	// If refresh token is invalid, returns invalidation error.
	Refresh(oldRefresh string) (access, refresh string, err error)

	// Verify checks access token validity.
	Verify(access string) error

	// Revoke puts access to black list.
	// Sequential calls to Verify with access always return non-nil error.
	// Sequential calls to Refresh with refresh, that contains refreshID from access, always return non-nil error.
	Revoke(access string) error
}
