/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ssh

import (
	"context"
	"errors"
	"strings"
	"syscall"
	"testing"
)

func TestOpenNodeShellErrorPointsToDocs(t *testing.T) {
	// t.Setenv keeps known_hosts out of the real home directory; it rules out t.Parallel.
	t.Setenv("HOME", t.TempDir())
	t.Setenv("USERPROFILE", t.TempDir())

	c := &client{username: "root", password: "irrelevant"}

	// nothing listens on port 1, so the dial is refused
	sshClient, err := c.openNodeShell(context.Background(), ProxmoxNode{Address: "127.0.0.1", Port: 1})
	if err == nil {
		_ = sshClient.Close()

		t.Fatal("expected an error, got nil")
	}

	if !strings.Contains(err.Error(), docsURL) {
		t.Fatalf("error %q does not mention %s", err.Error(), docsURL)
	}

	if !errors.Is(err, syscall.ECONNREFUSED) {
		t.Fatalf("error %q no longer wraps %v", err.Error(), syscall.ECONNREFUSED)
	}
}

func TestNewErrUserHasNoPermissionPointsToDocs(t *testing.T) {
	t.Parallel()

	err := NewErrUserHasNoPermission("terraform")
	if !strings.Contains(err.Error(), "'terraform'") || !strings.Contains(err.Error(), docsURL) {
		t.Fatalf("unexpected error message: %q", err.Error())
	}
}
