/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ssh

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"net"
	"strconv"
	"strings"
	"syscall"
	"testing"

	"golang.org/x/crypto/ssh"
)

// startRejectingServer starts an in-process SSH server that refuses every password.
func startRejectingServer(t *testing.T) string {
	t.Helper()

	hostKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate host key: %v", err)
	}

	signer, err := ssh.NewSignerFromKey(hostKey)
	if err != nil {
		t.Fatalf("signer: %v", err)
	}

	cfg := &ssh.ServerConfig{
		PasswordCallback: func(ssh.ConnMetadata, []byte) (*ssh.Permissions, error) {
			return nil, errors.New("denied")
		},
	}
	cfg.AddHostKey(signer)

	var lc net.ListenConfig

	ln, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	t.Cleanup(func() { _ = ln.Close() })

	go func() {
		for {
			conn, aerr := ln.Accept()
			if aerr != nil {
				return
			}

			go func() {
				_, _, _, serr := ssh.NewServerConn(conn, cfg)
				ignoreErr(serr)
				ignoreErr(conn.Close())
			}()
		}
	}()

	return ln.Addr().String()
}

func nodeFromAddr(t *testing.T, addr string) ProxmoxNode {
	t.Helper()

	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatalf("split %q: %v", addr, err)
	}

	port, err := strconv.ParseInt(portStr, 10, 32)
	if err != nil {
		t.Fatalf("parse port %q: %v", portStr, err)
	}

	return ProxmoxNode{Address: host, Port: int32(port)}
}

func TestOpenNodeShellErrorsPointToDocs(t *testing.T) {
	// t.Setenv keeps known_hosts out of the real home directory; it rules out t.Parallel.
	t.Setenv("HOME", t.TempDir())
	t.Setenv("USERPROFILE", t.TempDir())

	closedAddr := func(t *testing.T) string {
		t.Helper()

		var lc net.ListenConfig

		ln, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("listen: %v", err)
		}

		addr := ln.Addr().String()
		_ = ln.Close()

		return addr
	}

	tests := []struct {
		name    string
		addr    func(t *testing.T) string
		wantErr error
	}{
		{name: "connection refused", addr: closedAddr, wantErr: syscall.ECONNREFUSED},
		{name: "authentication failure", addr: startRejectingServer},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &client{username: "root", password: "wrong"}

			sshClient, err := c.openNodeShell(context.Background(), nodeFromAddr(t, tt.addr(t)))
			if err == nil {
				_ = sshClient.Close()

				t.Fatal("expected an error, got nil")
			}

			if sshClient != nil {
				t.Fatal("expected a nil client on error")
			}

			if !strings.Contains(err.Error(), docsURL) {
				t.Fatalf("error %q does not mention %s", err.Error(), docsURL)
			}

			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("error %q no longer wraps %v", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestOpenNodeShellSuccessHasNoDocsHint(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("USERPROFILE", t.TempDir())

	addr := startExecServer(t, 0, "", "")
	c := &client{username: "root", password: "irrelevant"}

	sshClient, err := c.openNodeShell(context.Background(), nodeFromAddr(t, addr))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_ = sshClient.Close()
}

func TestNewErrUserHasNoPermissionPointsToDocs(t *testing.T) {
	t.Parallel()

	err := NewErrUserHasNoPermission("terraform")
	if !strings.Contains(err.Error(), "'terraform'") || !strings.Contains(err.Error(), docsURL) {
		t.Fatalf("unexpected error message: %q", err.Error())
	}
}
