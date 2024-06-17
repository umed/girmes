//go:build linux

package cli

import (
	"context"
	"log/slog"
	"os/exec"
	"path"

	"github.com/umed/girmes/internal/logging"
)

func UserHomeDir(login string) string {
	return path.Join("/home", login)
}

func AddUser(ctx context.Context, login string, group string) error {
	logger := logging.L(ctx).With(slog.String("username", login))
	if UserExists(ctx, login) {
		logger.Debug("user exists")
		return nil
	}
	cmd := exec.Command("adduser",
		"-D",
		"-s", "/bin/sh",
		"-G", group,
		"-h", UserHomeDir(login),
		login)

	if err := cmd.Run(); err != nil {
		logger.Error("failed to execute command", logging.Err(err))
		return ErrFailedToCreateUser
	}
	logger.Debug("created user")
	return nil
}

func UserExists(ctx context.Context, user string) bool {
	cmd := exec.Command("id", "-u", user)
	if err := cmd.Run(); err != nil {
		logging.L(ctx).Debug("user does not exists", logging.Err(err), slog.String("username", user))
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
	logger := logging.L(ctx).With(slog.String("group", group))
	if GroupExists(ctx, group) {
		logger.Debug("group exists")
		return nil
	}
	cmd := exec.Command("addgroup", group)
	out, err := cmd.Output()
	if err != nil {
		logger.Error("failed to create group", logging.Err(err))
		return ErrFailedToCreateGroup
	}
	logger.Debug("completed group creation", slog.String("output", string(out)))
	return nil
}

func GroupExists(ctx context.Context, group string) bool {
	logger := logging.L(ctx).With(slog.String("group", group))
	cmd := exec.Command("getent", "group", group)
	if err := cmd.Run(); err != nil {
		logger.Debug("group not exists", logging.Err(err))
		return false
	}
	return true
}
