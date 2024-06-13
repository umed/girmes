//go:build darwin

package cli

import (
	"context"
	"path"
)

func AddUser(ctx context.Context, login string, group string) error {
	panic("not implemented")
}

func UserExists(ctx context.Context, user string) bool {
	panic("not implemented")
}

func AddUsers(ctx context.Context, users []string, group string) error {
	panic("not implemented")
}

func AddGroup(ctx context.Context, group string) error {
	panic("not implemented")
}

func UserHomeDir(login string) string {
	return path.Join("/Users", login)
}
