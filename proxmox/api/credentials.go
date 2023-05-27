/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package api

import (
	"errors"
	"strings"
)

const rootUsername = "root@pam"

type Credentials struct {
	Username string
	Password string
	OTP      *string
	APIToken *string
}

func NewCredentials(username, password, otp, apiToken string) (*Credentials, error) {
	if apiToken != "" {
		return &Credentials{
			APIToken: &apiToken,
		}, nil
	}

	if password == "" {
		return nil, errors.New(
			"you must specify a password for the Proxmox Virtual Environment API",
		)
	}

	if username == "" {
		return nil, errors.New(
			"you must specify a username for the Proxmox Virtual Environment API",
		)
	}

	if !strings.Contains(username, "@") {
		return nil, errors.New(
			"make sure the username for the Proxmox Virtual Environment API ends in '@pve or @pam'",
		)
	}

	c := &Credentials{
		Username: username,
		Password: password,
	}

	if otp != "" {
		c.OTP = &otp
	}

	return c, nil
}
