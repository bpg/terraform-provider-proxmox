/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package api

import (
	"errors"
	"regexp"
	"strings"
)

const rootUsername = "root@pam"

// Credentials contains the credentials for authenticating with the Proxmox VE API.
type Credentials struct {
	UserCredentials   *UserCredentials
	TokenCredentials  *TokenCredentials
	TicketCredentials *TicketCredentials
}

// UserCredentials contains the username, password, and OTP for authenticating with the Proxmox VE API.
type UserCredentials struct {
	Username string
	Password string
	OTP      string
}

// TokenCredentials contains the API token for authenticating with the Proxmox VE API.
type TokenCredentials struct {
	APIToken string
}

// TicketCredentials contains the auth ticket and CSRF prevention token for authenticating with the Proxmox VE API.
type TicketCredentials struct {
	AuthTicket          string
	CSRFPreventionToken string
}

// NewCredentials creates a new set of credentials for authenticating with the Proxmox VE API.
// The order of precedence is:
// 1. API token
// 2. Ticket
// 3. User credentials.
func NewCredentials(username, password, otp, apiToken, authTicket, csrfPreventionToken string) (Credentials, error) {
	tok, err := newTokenCredentials(apiToken)
	if err == nil {
		return Credentials{TokenCredentials: &tok}, nil
	}

	tic, err := newTicketCredentials(authTicket, csrfPreventionToken)
	if err == nil {
		return Credentials{TicketCredentials: &tic}, nil
	}

	usr, err := newUserCredentials(username, password, otp)
	if err == nil {
		return Credentials{UserCredentials: &usr}, nil
	}

	return Credentials{}, errors.New("must provide either user credentials, an API token, or a ticket")
}

func newUserCredentials(username, password, otp string) (UserCredentials, error) {
	if username == "" || password == "" {
		return UserCredentials{}, errors.New("both username and password are required")
	}

	if !strings.Contains(username, "@") {
		return UserCredentials{}, errors.New("username must end with '@pve' or '@pam'")
	}

	return UserCredentials{
		Username: username,
		Password: password,
		OTP:      otp,
	}, nil
}

func newTokenCredentials(apiToken string) (TokenCredentials, error) {
	re := regexp.MustCompile(`^\S+@\S+!\S+=([a-zA-Z0-9-]+)$`)
	if !re.MatchString(apiToken) {
		return TokenCredentials{}, errors.New("must be a valid API token, e.g. 'USER@REALM!TOKENID=UUID'")
	}

	return TokenCredentials{
		APIToken: apiToken,
	}, nil
}

func newTicketCredentials(authTicket, csrfPreventionToken string) (TicketCredentials, error) {
	if authTicket == "" || csrfPreventionToken == "" {
		return TicketCredentials{}, errors.New("both authTicket and csrfPreventionToken are required")
	}

	return TicketCredentials{
		AuthTicket:          authTicket,
		CSRFPreventionToken: csrfPreventionToken,
	}, nil
}
