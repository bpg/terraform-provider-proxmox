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
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

const (
	// PVE tickets are valid for 2 hours.
	rootTicketTTL = time.Hour

	// q() does not interpolate, so q(root\@pam) would mint a ticket for the user "root\@pam".
	mintTicketCommand = `perl -MPVE::AccessControl -e ` +
		`'printf("` + ticketPrefix + `%s\n` + csrfPrefix + `%s\n", ` +
		`PVE::AccessControl::assemble_ticket(q(root@pam)), ` +
		`PVE::AccessControl::assemble_csrf_prevention_token(q(root@pam)))'`

	// minting a ticket needs a root shell; sudo on perl would be root-equivalent anyway.
	rootShellUsername = "root"

	ticketPrefix = "PVE_TICKET="
	csrfPrefix   = "PVE_CSRF="
)

// NodeCommandExecutor is the part of ssh.Client this package needs. Declared here because
// proxmox/ssh imports this package.
type NodeCommandExecutor interface {
	Username() string
	ExecuteNodeCommands(ctx context.Context, nodeName string, commands []string) ([]byte, error)
}

type rootTicket struct {
	auth     Authenticator
	mintedAt time.Time
}

// rootTicketAuthenticator authenticates requests under /nodes/{node}/ with a root@pam ticket
// minted on that node over SSH, and everything else with the provider's own credentials.
type rootTicketAuthenticator struct {
	base     Authenticator
	executor NodeCommandExecutor
	clock    func() time.Time

	mu      sync.RWMutex
	tickets map[string]rootTicket
	minting singleflight.Group
}

func newRootTicketAuthenticator(base Authenticator, executor NodeCommandExecutor) *rootTicketAuthenticator {
	return &rootTicketAuthenticator{
		base:     base,
		executor: executor,
		clock:    time.Now,
		tickets:  map[string]rootTicket{},
	}
}

func (t *rootTicketAuthenticator) IsRoot(_ context.Context) bool {
	return true
}

func (t *rootTicketAuthenticator) IsRootTicket(_ context.Context) bool {
	return true
}

// AuthenticateRequest adds authentication data to a new request.
func (t *rootTicketAuthenticator) AuthenticateRequest(ctx context.Context, req *http.Request) error {
	node, ok := nodeFromPath(req.URL.Path)
	if !ok {
		return t.base.AuthenticateRequest(ctx, req)
	}

	auth, err := t.ticketAuth(ctx, node)
	if err != nil {
		return err
	}

	return auth.AuthenticateRequest(ctx, req)
}

// nodeFromPath returns the node name of a /nodes/{node}/... API path.
func nodeFromPath(path string) (string, bool) {
	path = strings.TrimPrefix(strings.TrimPrefix(path, "/"), basePathJSONAPI+"/")

	parts := strings.Split(path, "/")
	if len(parts) < 3 || parts[0] != "nodes" || parts[1] == "" || parts[2] == "" {
		return "", false
	}

	return parts[1], true
}

func (t *rootTicketAuthenticator) ticketAuth(ctx context.Context, node string) (Authenticator, error) {
	if auth, ok := t.cached(node); ok {
		return auth, nil
	}

	// WithoutCancel: one caller cancelling must not fail the shared mint for the other waiters
	ch := t.minting.DoChan(node, func() (any, error) {
		if auth, ok := t.cached(node); ok {
			return auth, nil
		}

		return t.mint(context.WithoutCancel(ctx), node)
	})

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("minting a root@pam ticket on node %q was cancelled: %w", node, ctx.Err())
	case res := <-ch:
		if res.Err != nil {
			return nil, res.Err
		}

		auth, ok := res.Val.(Authenticator)
		if !ok {
			return nil, fmt.Errorf("unexpected result while minting a root@pam ticket on node %q", node)
		}

		return auth, nil
	}
}

func (t *rootTicketAuthenticator) cached(node string) (Authenticator, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	tic, ok := t.tickets[node]
	if !ok || t.clock().Sub(tic.mintedAt) >= rootTicketTTL {
		return nil, false
	}

	return tic.auth, true
}

func (t *rootTicketAuthenticator) mint(ctx context.Context, node string) (Authenticator, error) {
	output, err := t.executor.ExecuteNodeCommands(ctx, node, []string{mintTicketCommand})
	if err != nil {
		return nil, fmt.Errorf("failed to mint a root@pam ticket on node %q over SSH: %s",
			node, scrubTicket(err.Error()))
	}

	ticket, csrf, err := parseTicketOutput(string(output))
	if err != nil {
		return nil, fmt.Errorf("failed to mint a root@pam ticket on node %q over SSH: %w", node, err)
	}

	auth, err := NewTicketAuthenticator(TicketCredentials{AuthTicket: ticket, CSRFPreventionToken: csrf})
	if err != nil {
		return nil, fmt.Errorf("the root@pam ticket minted on node %q is not usable: %w", node, err)
	}

	t.mu.Lock()
	t.tickets[node] = rootTicket{auth: auth, mintedAt: t.clock()}
	t.mu.Unlock()

	return auth, nil
}

// parseTicketOutput picks the values out of the combined stdout and stderr of the remote shell.
func parseTicketOutput(output string) (string, string, error) {
	var ticket, csrf string

	for line := range strings.SplitSeq(output, "\n") {
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, ticketPrefix):
			ticket = strings.TrimPrefix(line, ticketPrefix)
		case strings.HasPrefix(line, csrfPrefix):
			csrf = strings.TrimPrefix(line, csrfPrefix)
		}
	}

	if ticket == "" || csrf == "" {
		return "", "", errors.New("the node did not return a ticket and a CSRF prevention token; " +
			"check that the SSH user can run perl with PVE::AccessControl")
	}

	return ticket, csrf, nil
}

// scrubTicket keeps a partially successful command from leaking the ticket into an error message.
func scrubTicket(s string) string {
	lines := strings.Split(s, "\n")
	kept := make([]string, 0, len(lines))

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, ticketPrefix) || strings.HasPrefix(trimmed, csrfPrefix) ||
			strings.Contains(trimmed, "PVE:root@pam:") {
			continue
		}

		kept = append(kept, line)
	}

	return strings.TrimSpace(strings.Join(kept, "\n"))
}
