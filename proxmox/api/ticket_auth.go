/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package api

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type ticketAuthenticator struct {
	authData *AuthenticationResponseData
}

// NewTicketAuthenticator returns a new ticket authenticator.
func NewTicketAuthenticator(creds TicketCredentials) (Authenticator, error) {
	ard := &AuthenticationResponseData{}
	ard.Ticket = &(creds.AuthTicket)
	ard.CSRFPreventionToken = &(creds.CSRFPreventionToken)

	authTicketSplits := strings.Split(creds.AuthTicket, ":")

	if len(authTicketSplits) > 3 {
		ard.Username = strings.Split(creds.AuthTicket, ":")[1]
	} else {
		return nil, errors.New("AuthTicket must include a valid username")
	}

	if !strings.Contains(ard.Username, "@") {
		return nil, errors.New("username must end with '@pve' or '@pam'")
	}

	return &ticketAuthenticator{
		authData: ard,
	}, nil
}

func (t *ticketAuthenticator) IsRoot(_ context.Context) bool {
	return t.authData != nil && t.authData.Username == rootUsername
}

func (t *ticketAuthenticator) IsRootTicket(ctx context.Context) bool {
	return t.IsRoot(ctx)
}

// AuthenticateRequest adds authentication data to a new request.
func (t *ticketAuthenticator) AuthenticateRequest(_ context.Context, req *http.Request) error {
	req.AddCookie(&http.Cookie{
		HttpOnly: true,
		Name:     "PVEAuthCookie",
		Secure:   true,
		Value:    *t.authData.Ticket,
	})

	if req.Method != http.MethodGet {
		req.Header.Add("CSRFPreventionToken", *t.authData.CSRFPreventionToken)
	}

	return nil
}
