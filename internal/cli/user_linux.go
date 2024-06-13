//go:build linux

package cli

import (
	"context"
	"os/exec"
	"path"

	"github.com/umed/girmes/internal/logging"
)

func UserHomeDir(login string) string {
	return path.Join("/home", login)
}

func AddUser(ctx context.Context, login string, group string) error {
	logger := logging.L(ctx).With().Str("username", login).Logger()
	if UserExists(ctx, login) {
		logger.Debug().Msg("user exists")
		return nil
	}
	cmd := exec.Command("adduser",
		"-D",
		"-s", "/bin/sh",
		"-G", group,
		"-h", UserHomeDir(login),
		login)

	if err := cmd.Run(); err != nil {
		logger.Err(err).Msg("failed to create user")
		return ErrFailedToCreateUser
	}
	logger.Debug().Msg("created user")
	return nil
}

func UserExists(ctx context.Context, user string) bool {
	cmd := exec.Command("id", "-u", user)
	if err := cmd.Run(); err != nil {
		logging.L(ctx).Debug().Err(err).Str("username", user).Msg("user does not exists")
		return false
	}
	return true
}

func AddUsers(ctx context.Context, users []string, group string) error {
	for _, username := range users {
		if err := AddUser(ctx, username, group); err != nil {
			return err
		}
	}
	return nil
}

func AddGroup(ctx context.Context, group string) error {
	logger := logging.L(ctx).With().Str("group", group).Logger()
	if GroupExists(ctx, group) {
		logger.Debug().Msg("group exists")
		return nil
	}
	cmd := exec.Command("addgroup", group)
	out, err := cmd.Output()
	if err != nil {
		logger.Err(err).Msg("failed to create group")
		return ErrFailedToCreateGroup
	}
	logger.Debug().Str("output", string(out)).Msg("completed group creation")
	return nil
}

func GroupExists(ctx context.Context, group string) bool {
	logger := logging.L(ctx).With().Str("group", group).Logger()
	cmd := exec.Command("getent", "group", group)
	if err := cmd.Run(); err != nil {
		logger.Debug().Err(err).Msg("group not exists")
		return false
	}
	return true
}
