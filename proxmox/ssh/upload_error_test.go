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
	"encoding/binary"
	"io"
	"net"
	"os"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// startExecServer starts a minimal in-process SSH server that accepts a single
// session, drains the exec command's stdin, writes the given stdout/stderr, and
// exits with exitStatus. It is used to reproduce remote command failures.
func startExecServer(t *testing.T, exitStatus uint32, stdout, stderr string) string {
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
			return &ssh.Permissions{}, nil
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
		conn, aerr := ln.Accept()
		if aerr != nil {
			return
		}

		_, chans, reqs, serr := ssh.NewServerConn(conn, cfg)
		if serr != nil {
			return
		}

		go ssh.DiscardRequests(reqs)

		for newCh := range chans {
			if newCh.ChannelType() != "session" {
				ignoreErr(newCh.Reject(ssh.UnknownChannelType, "only session channels are supported"))
				continue
			}

			ch, chReqs, aerr := newCh.Accept()
			if aerr != nil {
				return
			}

			go handleSession(ch, chReqs, exitStatus, stdout, stderr)
		}
	}()

	return ln.Addr().String()
}

func handleSession(ch ssh.Channel, reqs <-chan *ssh.Request, exitStatus uint32, stdout, stderr string) {
	for req := range reqs {
		if req.Type != "exec" {
			if req.WantReply {
				ignoreErr(req.Reply(false, nil))
			}

			continue
		}

		if req.WantReply {
			ignoreErr(req.Reply(true, nil))
		}

		// Drain stdin so the client's stdin copy completes cleanly.
		_, err := io.Copy(io.Discard, ch)
		ignoreErr(err)

		if stdout != "" {
			_, err = io.WriteString(ch, stdout)
			ignoreErr(err)
		}

		if stderr != "" {
			_, err = io.WriteString(ch.Stderr(), stderr)
			ignoreErr(err)
		}

		status := make([]byte, 4)
		binary.BigEndian.PutUint32(status, exitStatus)

		_, err = ch.SendRequest("exit-status", false, status)
		ignoreErr(err)
		ignoreErr(ch.Close())
	}
}

// ignoreErr is a best-effort discard for the throwaway test SSH server, where
// connection teardown errors are expected and irrelevant to the assertion.
func ignoreErr(error) {}

// TestUploadFileSurfacesUnderlyingError verifies that when the remote transfer
// command fails, uploadFile reports the underlying cause instead of an empty
// "error transferring file:" message — both when the remote command produces no
// output (the reported bug) and when it does.
func TestUploadFileSurfacesUnderlyingError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		stderr string
	}{
		{name: "no remote output", stderr: ""},
		{name: "with remote output", stderr: "tee: /var/lib/vz/snippets/user-data.yml: No such file or directory"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			addr := startExecServer(t, 1, "", tt.stderr)

			tmp, err := os.CreateTemp(t.TempDir(), "snippet-*.yml")
			if err != nil {
				t.Fatalf("temp file: %v", err)
			}

			if _, err = tmp.WriteString("#cloud-config\n"); err != nil {
				t.Fatalf("write temp file: %v", err)
			}

			if _, err = tmp.Seek(0, 0); err != nil {
				t.Fatalf("seek temp file: %v", err)
			}

			sshClient, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
				User:            "root",
				Auth:            []ssh.AuthMethod{ssh.Password("irrelevant")},
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			})
			if err != nil {
				t.Fatalf("dial: %v", err)
			}

			defer func() { _ = sshClient.Close() }()

			c := &client{}
			req := &api.FileUploadRequest{FileName: "user-data.yml", File: tmp}

			err = c.uploadFile(context.Background(), sshClient, req, "/var/lib/vz/snippets/user-data.yml", "")
			if err == nil {
				t.Fatal("expected an error from the non-zero remote exit, got nil")
			}

			detail := strings.TrimSpace(strings.TrimPrefix(err.Error(), "error transferring file:"))
			if detail == "" {
				t.Fatalf("error message carries no detail beyond the prefix: %q", err.Error())
			}

			if tt.stderr != "" && !strings.Contains(err.Error(), tt.stderr) {
				t.Fatalf("error message %q does not include the remote output %q", err.Error(), tt.stderr)
			}
		})
	}
}
