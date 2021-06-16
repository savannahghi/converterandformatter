package converterandformatter

import (
	"context"
	"testing"

	"firebase.google.com/go/auth"
	"github.com/savannahghi/firebasetools"
	"github.com/stretchr/testify/assert"
)

// ContextKey is used as a type for the UID key for the Firebase *auth.Token on context.Context.
// It is a custom type in order to minimize context key collissions on the context
// (.and to shut up golint).
type ContextKey string

// GetAuthenticatedContext returns a logged in context, useful for test purposes
func GetAuthenticatedContext(t *testing.T) context.Context {
	ctx := context.Background()
	authToken := getAuthToken(ctx, t)
	authenticatedContext := context.WithValue(ctx, AuthTokenContextKey, authToken)
	return authenticatedContext
}

// GetAuthenticatedContextAndToken returns a logged in context and ID token.
// It is useful for test purposes
func GetAuthenticatedContextAndToken(t *testing.T) (context.Context, *auth.Token) {
	ctx := context.Background()
	authToken := getAuthToken(ctx, t)
	authenticatedContext := context.WithValue(ctx, AuthTokenContextKey, authToken)
	return authenticatedContext, authToken
}

// GetAuthenticatedContextAndBearerToken returns a logged in context and bearer token.
// It is useful for test purposes
func GetAuthenticatedContextAndBearerToken(t *testing.T) (context.Context, string) {
	ctx := context.Background()
	authToken, bearerToken := getAuthTokenAndBearerToken(ctx, t)
	authenticatedContext := context.WithValue(ctx, AuthTokenContextKey, authToken)
	return authenticatedContext, bearerToken
}

func getAuthToken(ctx context.Context, t *testing.T) *auth.Token {
	authToken, _ := getAuthTokenAndBearerToken(ctx, t)
	return authToken
}

func getAuthTokenAndBearerToken(ctx context.Context, t *testing.T) (*auth.Token, string) {
	user, userErr := firebasetools.GetOrCreateFirebaseUser(ctx, TestUserEmail)
	assert.Nil(t, userErr)
	assert.NotNil(t, user)

	customToken, tokenErr := firebasetools.CreateFirebaseCustomToken(ctx, user.UID)
	assert.Nil(t, tokenErr)
	assert.NotNil(t, customToken)

	idTokens, idErr := firebasetools.AuthenticateCustomFirebaseToken(customToken)
	assert.Nil(t, idErr)
	assert.NotNil(t, idTokens)

	bearerToken := idTokens.IDToken
	authToken, err := firebasetools.ValidateBearerToken(ctx, bearerToken)
	assert.Nil(t, err)
	assert.NotNil(t, authToken)

	return authToken, bearerToken
}
