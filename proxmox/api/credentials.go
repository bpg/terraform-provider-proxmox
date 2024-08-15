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

// Credentials is a struct that holds the credentials for the Proxmox Virtual
// Environment API.
type Credentials struct {
	AuthTicket          string
	CSRFPreventionToken string
	APIToken            *string
	OTP                 *string
	Username            string
	Password            string
}

// NewCredentials creates a new Credentials struct.
//
//nolint:lll
func NewCredentials(username, password, otp, apiToken string, authTicket string, csrfPreventionToken string) (*Credentials, error) {
	errorAuthTicketCommonMsg := "the Proxmox Virtual Environment API requires auth params; "

	switch {
	case authTicket != "" && csrfPreventionToken != "":
		switch {
		case authTicket == "":
			return nil, errors.New(errorAuthTicketCommonMsg + "detected csrfPreventionToken, but authTicket is unset")
		case csrfPreventionToken == "":
			return nil, errors.New(errorAuthTicketCommonMsg + "detected authTicket, but csrfPreventionToken is unset")
		}
	case apiToken != "":
		re := regexp.MustCompile(`^\S+@\S+!\S+=([a-zA-Z0-9-]+)$`)
		if !re.MatchString(apiToken) {
			return nil, errors.New("must be a valid API token, e.g. 'USER@REALM!TOKENID=UUID'")
		}

		return &Credentials{
			APIToken: &apiToken,
		}, nil
	case (username != "" && password != "") || (username != "" || password != ""):
		switch {
		case username == "":
			return nil, errors.New(errorAuthTicketCommonMsg + "detected password, but username is unset")
		case password == "":
			return nil, errors.New(errorAuthTicketCommonMsg + "detected username, but password is unset")
		}
	default:
		return nil, errors.New(errorAuthTicketCommonMsg +
			"choose either: authTicket + csrfPreventionToken, apiToken; username + password")
	}

	if username != "" && !strings.Contains(username, "@") {
		return nil, errors.New(
			"make sure the username for the Proxmox Virtual Environment API ends in '@pve or @pam'",
		)
	}

	c := &Credentials{
		AuthTicket:          authTicket,
		CSRFPreventionToken: csrfPreventionToken,
		Username:            username,
		Password:            password,
	}

	if otp != "" {
		c.OTP = &otp
	}

	return c, nil
}
