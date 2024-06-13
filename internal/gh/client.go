package gh

import (
	"context"
	"errors"

	"github.com/google/go-github/v62/github"
	"github.com/umed/girmes/internal/logging"
)

var ErrGithubApiCallFailed = errors.New("github error")

type Client struct {
	githubClient *github.Client
}

func NewClient(token string) *Client {
	if len(token) == 0 {
		return &Client{
			githubClient: github.NewClient(nil),
		}
	}
	return &Client{
		githubClient: github.NewClient(nil).WithAuthToken(token),
	}
}

func (client *Client) FetchLogins(ctx context.Context, orgName string) ([]string, error) {
	logger := logging.L(ctx).With().Str("org", orgName).Logger()
	logger.Debug().Msg("fetching logins")
	users, _, err := client.githubClient.Organizations.ListMembers(ctx, orgName, nil)
	if err != nil {
		logging.L(ctx).Err(err).Str("org", orgName).Msg("failed to list members")
		return nil, ErrGithubApiCallFailed
	}
	logins := make([]string, len(users))
	for i := range users {
		logins[i] = users[i].GetLogin()
	}
	logger.Debug().Strs("logins", logins).Msg("completed fetching logins")
	return logins, nil
}

type User struct {
	Login string
	Keys  []string
}

func (client *Client) FetchUsers(ctx context.Context, orgName string) ([]User, error) {
	logger := logging.L(ctx).With().Str("org", orgName).Logger()
	logger.Debug().Msg("fetching users")
	logins, err := client.FetchLogins(ctx, orgName)
	if err != nil {
		return nil, err
	}
	users := make([]User, len(logins))
	fails := 0
	for i, login := range logins {
		users[i].Login = login
		users[i].Keys, err = client.FetchKeys(ctx, login)
		if err != nil {
			logger.Warn().Str("user", login).Msg("failed to acquire keys, will be skipped")
			fails++
			continue
		}
	}
	if fails == len(logins) {
		logger.Error().Msg("failed to get keys for any user")
		return nil, ErrGithubApiCallFailed
	}
	logger.Debug().Msg("completed fetching users")
	return users, nil
}

func (client *Client) FetchKeys(ctx context.Context, user string) ([]string, error) {
	logger := logging.L(ctx).With().Str("username", user).Logger()
	logger.Debug().Msg("fetching keys")
	keys, _, err := client.githubClient.Users.ListKeys(ctx, user, nil)
	if err != nil {
		logging.L(ctx).Err(err)
		return nil, ErrGithubApiCallFailed
	}
	keysArray := make([]string, len(keys))
	for i := range keys {
		keysArray[i] = keys[i].GetKey()
	}
	logger.Debug().Int("number", len(keys)).Msg("completed fetching keys")
	return keysArray, nil
}
