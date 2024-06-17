package cli

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path"

	"github.com/umed/girmes/internal/logging"
)

var (
	ErrFailedToCreateUser  = errors.New("failed to create user")
	ErrFailedToCreateGroup = errors.New("failed to create group")
	ErrUserNotExists       = errors.New("user not exists")
	ErrFailedToAddKeys     = errors.New("failed to add keys")
)

func AddAuthorizedKeys(ctx context.Context, login string, keys []string) error {
	logger := logging.L(ctx).With(slog.String("username", login))
	if !UserExists(ctx, login) {
		logger.Warn("user not exists")
		return ErrUserNotExists
	}
	sshDir := path.Join(UserHomeDir(login), ".ssh")
	if err := os.MkdirAll(sshDir, os.ModePerm); err != nil {
		logger.Error("failed to create ssh directory", logging.Err(err), slog.String("dir", sshDir))
		return ErrFailedToAddKeys
	}
	authorizedKeysFile := path.Join(sshDir, "authorized_keys")
	f, err := os.OpenFile(authorizedKeysFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		logger.Error("failed to open authorized keys file", logging.Err(err), slog.String("keys_file", authorizedKeysFile))
		return ErrFailedToAddKeys
	}
	defer f.Close()

	for _, key := range keys {
		if _, err = f.WriteString(key + "\n"); err != nil {
			logger.Error("failed to append keys", logging.Err(err), slog.String("keys_file", authorizedKeysFile))
			return ErrFailedToAddKeys
		}
	}
	logger.Debug("successfully added keys")
	return nil
}
