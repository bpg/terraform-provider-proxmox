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
func NewTokenAuthenticator(token string) (Authenticator, error) {
	return &tokenAuthenticator{
		username: strings.Split(token, "!")[0],
		token:    token,
	}, nil
}

func (t *tokenAuthenticator) IsRoot() bool {
	return t.username == rootUsername
}

func (t *tokenAuthenticator) AuthenticateRequest(ctx context.Context, req *http.Request) error {
	req.Header.Set("Authorization", "PVEAPIToken="+t.token)
	return nil
}
