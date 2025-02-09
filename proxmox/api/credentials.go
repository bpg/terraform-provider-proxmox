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

// Package level error declarations.
var (
	ErrMissingAPIToken          = errors.New("no API token provided")
	ErrMissingTicketCredentials = errors.New("no authTicket and csrfPreventionToken pair provided")
	ErrMissingUserCredentials   = errors.New("no username and password provided")
	ErrInvalidAPIToken          = errors.New("the API token must be in the format 'USER@REALM!TOKENID=UUID'")
	ErrInvalidUsernameFormat    = errors.New("the username must end with '@pve' or '@pam'")
)

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
	if tok, err := newTokenCredentials(apiToken); err == nil {
		return Credentials{TokenCredentials: &tok}, nil
	} else if errors.Is(err, ErrInvalidAPIToken) {
		return Credentials{}, err
	}

	if tic, err := newTicketCredentials(authTicket, csrfPreventionToken); err == nil {
		return Credentials{TicketCredentials: &tic}, nil
	}

	if usr, err := newUserCredentials(username, password, otp); err == nil {
		return Credentials{UserCredentials: &usr}, nil
	} else if errors.Is(err, ErrInvalidUsernameFormat) {
		return Credentials{}, err
	}

	return Credentials{}, errors.New("must provide either username and password, an API token, or a ticket")
}

func newUserCredentials(username, password, otp string) (UserCredentials, error) {
	if username == "" || password == "" {
		return UserCredentials{}, ErrMissingUserCredentials
	}

	if !strings.Contains(username, "@") {
		return UserCredentials{}, ErrInvalidUsernameFormat
	}

	return UserCredentials{
		Username: username,
		Password: password,
		OTP:      otp,
	}, nil
}

func newTokenCredentials(apiToken string) (TokenCredentials, error) {
	if apiToken == "" {
		return TokenCredentials{}, ErrMissingAPIToken
	}

	re := regexp.MustCompile(`^\S+@\S+!\S+=([a-zA-Z0-9-]+)$`)
	if !re.MatchString(apiToken) {
		return TokenCredentials{}, ErrInvalidAPIToken
	}

	return TokenCredentials{
		APIToken: apiToken,
	}, nil
}

func newTicketCredentials(authTicket, csrfPreventionToken string) (TicketCredentials, error) {
	if authTicket == "" || csrfPreventionToken == "" {
		return TicketCredentials{}, ErrMissingTicketCredentials
	}

	return TicketCredentials{
		AuthTicket:          authTicket,
		CSRFPreventionToken: csrfPreventionToken,
	}, nil
}
