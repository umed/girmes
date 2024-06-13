package cli

import (
	"context"
	"errors"
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
	logger := logging.L(ctx).With().Str("username", login).Logger()
	if !UserExists(ctx, login) {
		logger.Warn().Msg("user not exists")
		return ErrUserNotExists
	}
	sshDir := path.Join(UserHomeDir(login), ".ssh")
	if err := os.MkdirAll(sshDir, os.ModePerm); err != nil {
		logger.Err(err).Str("dir", sshDir).Msg("failed to create ssh directory")
		return ErrFailedToAddKeys
	}
	authorizedKeysFile := path.Join(sshDir, "authorized_keys")
	f, err := os.OpenFile(authorizedKeysFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		logger.Err(err).Str("keys_file", authorizedKeysFile).Msg("failed to open authorized keys file")
		return ErrFailedToAddKeys
	}
	defer f.Close()

	for _, key := range keys {
		if _, err = f.WriteString(key + "\n"); err != nil {
			logger.Err(err).Str("keys_file", authorizedKeysFile).Msg("failed to append keys")
			return ErrFailedToAddKeys
		}
	}
	logger.Debug().Msg("successfully added keys")
	return nil
}
