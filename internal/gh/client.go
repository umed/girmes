package gh

import (
	"context"
	"errors"
	"log/slog"

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

func (client *Client) FetchMembers(ctx context.Context, orgName string) ([]string, error) {
	logger := logging.L(ctx).With(slog.String("org", orgName))
	logger.Debug("fetching org's members")
	users, _, err := client.githubClient.Organizations.ListMembers(ctx, orgName, nil)
	if err != nil {
		logger.Error("failed to list members", slog.String("org", orgName), logging.Err(err))
		return nil, ErrGithubApiCallFailed
	}
	logins := make([]string, len(users))
	for i := range users {
		logins[i] = users[i].GetLogin()
	}
	logger.Debug("completed fetching org's members", slog.Any("logins", logins))
	return logins, nil
}

type User struct {
	Login string
	Keys  []string
}

func (client *Client) FetchUsers(ctx context.Context, orgName string) ([]User, error) {
	logger := logging.L(ctx).With(slog.String("org", orgName))
	logger.Debug("fetching users")
	logins, err := client.FetchMembers(ctx, orgName)
	if err != nil {
		return nil, err
	}
	users := make([]User, len(logins))
	fails := 0
	for i, login := range logins {
		users[i].Login = login
		users[i].Keys, err = client.FetchKeys(ctx, login)
		if err != nil {
			logger.Warn("failed to acquire keys, will be skipped", slog.String("user", login))
			fails++
			continue
		}
	}
	if fails == len(logins) {
		logger.Error("failed to get keys for any user")
		return nil, ErrGithubApiCallFailed
	}
	logger.Debug("completed fetching users")
	return users, nil
}

func (client *Client) FetchKeys(ctx context.Context, user string) ([]string, error) {
	logger := logging.L(ctx).With(slog.String("username", user))
	logger.Debug("fetching keys")
	keys, _, err := client.githubClient.Users.ListKeys(ctx, user, nil)
	if err != nil {
		logger.Error("failed to list user's keys", logging.Err(err))
		return nil, ErrGithubApiCallFailed
	}
	keysArray := make([]string, len(keys))
	for i := range keys {
		keysArray[i] = keys[i].GetKey()
	}
	logger.Debug("completed fetching keys", slog.Int("number", len(keys)))
	return keysArray, nil
}
