/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bpg/terraform-provider-proxmox/utils"
)

type ticketAuthenticator struct {
	conn        *Connection
	authRequest string
	authData    *AuthenticationResponseData
}

// NewTicketAuthenticator returns a new ticket authenticator.
func NewTicketAuthenticator(conn *Connection, creds *Credentials) (Authenticator, error) {
	authRequest := fmt.Sprintf(
		"username=%s&password=%s",
		url.QueryEscape(creds.Username),
		url.QueryEscape(creds.Password),
	)

	// OTP is optional, and probably doesn't make much sense for most provider users.
	if creds.OTP != nil {
		authRequest = fmt.Sprintf("%s&otp=%s", authRequest, url.QueryEscape(*creds.OTP))
	}

	return &ticketAuthenticator{
		conn:        conn,
		authRequest: authRequest,
	}, nil
}

func (t *ticketAuthenticator) authenticate(ctx context.Context) (*AuthenticationResponseData, error) {
	if t.authData != nil {
		return t.authData, nil
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/%s/access/ticket", t.conn.endpoint, basePathJSONAPI),
		bytes.NewBufferString(t.authRequest),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create authentication request: %w", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	tflog.Debug(ctx, "Sending authentication request", map[string]interface{}{
		"path": req.URL.Path,
	})

	//nolint:bodyclose
	res, err := t.conn.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve authentication response: %w", err)
	}

	defer utils.CloseOrLogError(ctx)(res.Body)

	err = validateResponseCode(res)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}

	resBody := AuthenticationResponseBody{}

	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		return nil, fmt.Errorf("failed to decode authentication response, %w", err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the authentication response")
	}

	if resBody.Data.CSRFPreventionToken == nil {
		return nil, errors.New(
			"the server did not include a CSRF prevention token in the authentication response",
		)
	}

	if resBody.Data.Ticket == nil {
		return nil, errors.New("the server did not include a ticket in the authentication response")
	}

	if resBody.Data.Username == "" {
		return nil, errors.New("the server did not include the username in the authentication response")
	}

	t.authData = resBody.Data

	return resBody.Data, nil
}

func (t *ticketAuthenticator) IsRoot() bool {
	return t.authData != nil && t.authData.Username == rootUsername
}

func (t *ticketAuthenticator) IsRootTicket() bool {
	return t.IsRoot()
}

// AuthenticateRequest adds authentication data to a new request.
func (t *ticketAuthenticator) AuthenticateRequest(ctx context.Context, req *http.Request) error {
	a, err := t.authenticate(ctx)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	req.AddCookie(&http.Cookie{
		Name:  "PVEAuthCookie",
		Value: *a.Ticket,
	})

	if req.Method != http.MethodGet {
		req.Header.Add("CSRFPreventionToken", *a.CSRFPreventionToken)
	}

	return nil
}
