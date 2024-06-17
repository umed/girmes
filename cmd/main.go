package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

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
	logger.Debug("initialized", slog.Any("config", cfg))

	if len(*orgName) == 0 {
		logger.Fatal("org name is not provided")
	}

	ctx := context.Background()
	ctx = logging.Ctx(ctx, logger)

	client := gh.NewClient(cfg.GitHubAccessToken)
	users := util.Must(client.FetchUsers(ctx, *orgName))

	if err := cli.AddGroup(ctx, *orgName); err != nil {
		logger.Error("failed to create group", logging.Err(err))
		os.Exit(1)
	}

	for _, user := range users {
		if err := cli.AddUser(ctx, user.Login, *orgName); err != nil {
			logger.Fatal("failed to create user", logging.Err(err))
		}
		if err := cli.AddAuthorizedKeys(ctx, user.Login, user.Keys); err != nil {
			logger.Fatal("failed to add keys", logging.Err(err))
		}
	}
}
