package gh_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/umed/girmes/internal/gh"
	"github.com/umed/girmes/internal/logging"
)

func TestFetchLogins(t *testing.T) {
	ctx := logging.WithLogger(context.TODO(), logging.NewLogger("debug"))
	client := gh.NewClient("")
	logins, err := client.FetchLogins(ctx, "girmes")

	require.NoError(t, err, "failed to fetch users")

	require.Equal(t, []string{"umed"}, logins)
}

func TestFetchUsers(t *testing.T) {
	ctx := logging.WithLogger(context.TODO(), logging.NewLogger("debug"))
	client := gh.NewClient("")
	users, err := client.FetchUsers(ctx, "girmes")

	require.NoError(t, err, "failed to fetch users")
	require.Equal(t, 1, len(users))
	require.Equal(t, "umed", users[0].Login)
	require.True(t, len(users[0].Keys) > 0)
}
