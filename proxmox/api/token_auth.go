/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package api

import (
	"context"
	"net/http"
	"strings"
)

type tokenAuthenticator struct {
	username string
	token    string
}

// NewTokenAuthenticator creates a new authenticator that uses a PVE API Token
// for authentication.
func NewTokenAuthenticator(toc TokenCredentials) (Authenticator, error) {
	return &tokenAuthenticator{
		username: strings.Split(toc.APIToken, "!")[0],
		token:    toc.APIToken,
	}, nil
}

func (t *tokenAuthenticator) IsRoot(_ context.Context) bool {
	return t.username == rootUsername
}

func (t *tokenAuthenticator) IsRootTicket(_ context.Context) bool {
	// Logged using a token, therefore not a ticket login
	return false
}

func (t *tokenAuthenticator) AuthenticateRequest(_ context.Context, req *http.Request) error {
	req.Header.Set("Authorization", "PVEAPIToken="+t.token)
	return nil
}
