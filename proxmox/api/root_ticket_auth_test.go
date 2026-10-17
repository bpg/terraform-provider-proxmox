/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const testTicket = "PVE:root@pam:68F1A2B3::c29tZS1zaWduYXR1cmU="

type fakeExecutor struct {
	username string
	calls    atomic.Int32
	output   func(node string) string
	err      error
}

func (e *fakeExecutor) Username() string {
	if e.username == "" {
		return "root"
	}

	return e.username
}

func (e *fakeExecutor) ExecuteNodeCommands(_ context.Context, node string, _ []string) ([]byte, error) {
	e.calls.Add(1)

	if e.err != nil {
		return nil, e.err
	}

	if e.output != nil {
		return []byte(e.output(node)), nil
	}

	return fmt.Appendf(nil, "%s%s\n%s%s\n", ticketPrefix, testTicket, csrfPrefix, "68F1A2B3:dG9rZW4="), nil
}

type recordingAuthenticator struct {
	calls atomic.Int32
}

func (a *recordingAuthenticator) IsRoot(context.Context) bool       { return false }
func (a *recordingAuthenticator) IsRootTicket(context.Context) bool { return false }

func (a *recordingAuthenticator) AuthenticateRequest(_ context.Context, req *http.Request) error {
	a.calls.Add(1)
	req.Header.Set("Authorization", "PVEAPIToken=user@pve!id=secret")

	return nil
}

func newRequest(t *testing.T, path string) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://pve:8006"+path, nil)
	require.NoError(t, err)

	return req
}

func TestNodeFromPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		path string
		node string
		want bool
	}{
		{path: "/api2/json/nodes/pve1/lxc/100/config", node: "pve1", want: true},
		{path: "/api2/json/nodes/pve1/qemu", node: "pve1", want: true},
		{path: "api2/json/nodes/pve1/network", node: "pve1", want: true},
		{path: "/nodes/pve1/storage", node: "pve1", want: true},
		{path: "/api2/json/nodes/pve1", want: false},
		{path: "/api2/json/nodes/pve1/", want: false},
		{path: "/api2/json/nodes", want: false},
		{path: "/api2/json/nodes//lxc", want: false},
		{path: "/api2/json/cluster/resources", want: false},
		{path: "/api2/json/access/ticket", want: false},
		{path: "/api2/json/version", want: false},
		{path: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()

			node, ok := nodeFromPath(tt.path)
			require.Equal(t, tt.want, ok)
			require.Equal(t, tt.node, node)
		})
	}
}

func TestRootTicketAuthenticatorRouting(t *testing.T) {
	t.Parallel()

	base := &recordingAuthenticator{}
	exec := &fakeExecutor{}
	auth := newRootTicketAuthenticator(base, exec)

	nodeReq := newRequest(t, "/api2/json/nodes/pve1/lxc/100/config")
	nodeReq.Method = http.MethodPut
	require.NoError(t, auth.AuthenticateRequest(t.Context(), nodeReq))

	cookie, err := nodeReq.Cookie("PVEAuthCookie")
	require.NoError(t, err)
	require.Equal(t, testTicket, cookie.Value)
	require.Equal(t, "68F1A2B3:dG9rZW4=", nodeReq.Header.Get("CSRFPreventionToken"))
	require.Empty(t, nodeReq.Header.Get("Authorization"))
	require.Zero(t, base.calls.Load())

	clusterReq := newRequest(t, "/api2/json/cluster/resources")
	require.NoError(t, auth.AuthenticateRequest(t.Context(), clusterReq))

	_, err = clusterReq.Cookie("PVEAuthCookie")
	require.Error(t, err)
	require.NotEmpty(t, clusterReq.Header.Get("Authorization"))
	require.Equal(t, int32(1), base.calls.Load())
}

func TestRootTicketAuthenticatorCaching(t *testing.T) {
	t.Parallel()

	now := time.Now()
	exec := &fakeExecutor{}
	auth := newRootTicketAuthenticator(&recordingAuthenticator{}, exec)
	auth.clock = func() time.Time { return now }

	for range 3 {
		require.NoError(t, auth.AuthenticateRequest(t.Context(), newRequest(t, "/api2/json/nodes/pve1/lxc/100/config")))
	}

	require.Equal(t, int32(1), exec.calls.Load(), "the ticket should be minted once per node")

	// a second node needs its own ticket
	require.NoError(t, auth.AuthenticateRequest(t.Context(), newRequest(t, "/api2/json/nodes/pve2/lxc/100/config")))
	require.Equal(t, int32(2), exec.calls.Load())

	// still fresh just before the TTL, re-minted once past it
	now = now.Add(rootTicketTTL - time.Second)

	require.NoError(t, auth.AuthenticateRequest(t.Context(), newRequest(t, "/api2/json/nodes/pve1/lxc/100/config")))
	require.Equal(t, int32(2), exec.calls.Load())

	now = now.Add(2 * time.Second)

	require.NoError(t, auth.AuthenticateRequest(t.Context(), newRequest(t, "/api2/json/nodes/pve1/lxc/100/config")))
	require.Equal(t, int32(3), exec.calls.Load())
}

func TestRootTicketAuthenticatorConcurrentMint(t *testing.T) {
	t.Parallel()

	exec := &fakeExecutor{}
	auth := newRootTicketAuthenticator(&recordingAuthenticator{}, exec)

	var wg sync.WaitGroup

	errs := make([]error, 16)

	for i := range 16 {
		wg.Go(func() {
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet,
				"https://pve:8006/api2/json/nodes/pve1/lxc/100/config", nil)
			if err != nil {
				errs[i] = err

				return
			}

			errs[i] = auth.AuthenticateRequest(context.Background(), req)
		})
	}

	wg.Wait()
	require.NoError(t, errors.Join(errs...))
	require.Equal(t, int32(1), exec.calls.Load())
}

func TestRootTicketAuthenticatorMintErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		exec   *fakeExecutor
		errMsg string
	}{
		{
			name:   "command fails",
			exec:   &fakeExecutor{err: errors.New("bash: perl: command not found")},
			errMsg: "bash: perl: command not found",
		},
		{
			name:   "no ticket in the output",
			exec:   &fakeExecutor{output: func(string) string { return "Can't locate PVE/AccessControl.pm\n" }},
			errMsg: "did not return a ticket",
		},
		{
			name: "ticket without a username",
			exec: &fakeExecutor{output: func(string) string {
				return ticketPrefix + "garbage\n" + csrfPrefix + "68F1A2B3:dG9rZW4=\n"
			}},
			errMsg: "not usable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			auth := newRootTicketAuthenticator(&recordingAuthenticator{}, tt.exec)

			err := auth.AuthenticateRequest(t.Context(), newRequest(t, "/api2/json/nodes/pve1/lxc/100/config"))
			require.ErrorContains(t, err, tt.errMsg)
		})
	}
}

func TestRootTicketMintErrorDoesNotLeakTheTicket(t *testing.T) {
	t.Parallel()

	exec := &fakeExecutor{err: fmt.Errorf("%s%s\nssh: session closed", ticketPrefix, testTicket)}
	auth := newRootTicketAuthenticator(&recordingAuthenticator{}, exec)

	err := auth.AuthenticateRequest(t.Context(), newRequest(t, "/api2/json/nodes/pve1/lxc/100/config"))
	require.Error(t, err)
	require.NotContains(t, err.Error(), testTicket)
	require.Contains(t, err.Error(), "ssh: session closed")
}

func TestParseTicketOutput(t *testing.T) {
	t.Parallel()

	// the SSH client returns the combined output, so unrelated lines can precede the values
	output := "Warning: Permanently added 'pve1' to the list of known hosts.\n" +
		ticketPrefix + testTicket + "\n" + csrfPrefix + "68F1A2B3:dG9rZW4=\n"

	ticket, csrf, err := parseTicketOutput(output)
	require.NoError(t, err)
	require.Equal(t, testTicket, ticket)
	require.Equal(t, "68F1A2B3:dG9rZW4=", csrf)
}

func TestNewRootTicketClient(t *testing.T) {
	t.Parallel()

	conn, err := NewConnection("https://pve:8006", true, "", nil)
	require.NoError(t, err)

	inner, err := NewClient(Credentials{TokenCredentials: &TokenCredentials{
		APIToken: "user@pve!id=11111111-1111-1111-1111-111111111111",
	}}, conn)
	require.NoError(t, err)

	t.Run("root SSH user", func(t *testing.T) {
		t.Parallel()

		c, err := NewRootTicketClient(inner, &fakeExecutor{})
		require.NoError(t, err)
		require.True(t, c.IsRoot(t.Context()))
		require.True(t, c.IsRootTicket(t.Context()))
		require.False(t, inner.IsRootTicket(t.Context()), "the wrapped client must not be modified")
	})

	t.Run("non-root SSH user is refused", func(t *testing.T) {
		t.Parallel()

		_, err := NewRootTicketClient(inner, &fakeExecutor{username: "terraform"})
		require.ErrorContains(t, err, `the SSH user is "terraform"`)
	})
}
