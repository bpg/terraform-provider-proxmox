/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package proxmox

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	// DefaultRootAccount contains the default username and realm for the root account.
	DefaultRootAccount = "root@pam"
)

// Authenticate authenticates against the specified endpoint.
func (c *VirtualEnvironmentClient) Authenticate(ctx context.Context, reset bool) error {
	if c.authenticationData != nil && !reset {
		return nil
	}

	var reqBody *bytes.Buffer

	if c.OTP != nil {
		reqBody = bytes.NewBufferString(fmt.Sprintf(
			"username=%s&password=%s&otp=%s",
			url.QueryEscape(c.Username),
			url.QueryEscape(c.Password),
			url.QueryEscape(*c.OTP),
		))
	} else {
		reqBody = bytes.NewBufferString(fmt.Sprintf(
			"username=%s&password=%s",
			url.QueryEscape(c.Username),
			url.QueryEscape(c.Password),
		))
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/%s/access/ticket", c.Endpoint, basePathJSONAPI),
		reqBody,
	)
	if err != nil {
		return errors.New("failed to create authentication request")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	tflog.Debug(ctx, "sending authentication request", map[string]interface{}{
		"path": req.URL.Path,
	})
	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to retrieve authentication response: %w", err)
	}

	defer CloseOrLogError(ctx)(res.Body)

	err = c.ValidateResponseCode(res)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	resBody := AuthenticationResponseBody{}
	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		return fmt.Errorf("failed to decode authentication response, %w", err)
	}

	if resBody.Data == nil {
		return errors.New("the server did not include a data object in the authentication response")
	}

	if resBody.Data.CSRFPreventionToken == nil {
		return errors.New(
			"the server did not include a CSRF prevention token in the authentication response",
		)
	}

	if resBody.Data.Ticket == nil {
		return errors.New("the server did not include a ticket in the authentication response")
	}

	if resBody.Data.Username == "" {
		return errors.New("the server did not include the username in the authentication response")
	}

	c.authenticationData = resBody.Data

	return nil
}

// AuthenticateRequest adds authentication data to a new request.
func (c *VirtualEnvironmentClient) AuthenticateRequest(ctx context.Context, req *http.Request) error {
	err := c.Authenticate(ctx, false)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	req.AddCookie(&http.Cookie{
		Name:  "PVEAuthCookie",
		Value: *c.authenticationData.Ticket,
	})

	if req.Method != http.MethodGet {
		req.Header.Add("CSRFPreventionToken", *c.authenticationData.CSRFPreventionToken)
	}

	return nil
}
