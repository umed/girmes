package main

import (
	"context"
	"flag"

	"github.com/umed/girmes/config"
	"github.com/umed/girmes/internal/cli"
	"github.com/umed/girmes/internal/gh"
	"github.com/umed/girmes/internal/logging"
	"github.com/umed/girmes/internal/util"
)

var (
	orgName = flag.String("org", "", "organization to fetch users from")
	shell   = flag.String("shell", "/bin/sh", "default user shell")
)

func main() {
	flag.Parse()
	cfg := config.NewConfig()

	logger := logging.NewLogger(cfg.LogLevel)
	logger.Debug().Any("config", cfg).Msg("initialized")

	if orgName == nil {
		logger.Error().Msg("org name must be explicitly specified")
	}

	ctx := context.Background()
	ctx = logging.WithLogger(ctx, logger)

	client := gh.NewClient(cfg.GitHubAccessToken)
	users := util.Must(client.FetchUsers(ctx, *orgName))

	if err := cli.AddGroup(ctx, *orgName); err != nil {
		logger.Fatal().Err(err).Msg("failed to create group")
	}

	for _, user := range users {
		if err := cli.AddUser(ctx, user.Login, *orgName); err != nil {
			logger.Fatal().Err(err).Msg("failed to create user")
		}
		if err := cli.AddAuthorizedKeys(ctx, user.Login, user.Keys); err != nil {
			logger.Fatal().Err(err).Msg("failed to add keys")
		}
	}
}
