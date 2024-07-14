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
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bpg/terraform-provider-proxmox/utils"
)

type ticketAuthenticator struct {
	conn        *Connection
	authRequest string
	authData    *AuthenticationResponseData

	mu sync.Mutex
}

// custom bool to handle bools of: int, text, bool // ref: https://stackoverflow.com/questions/4278430/convert-string-to-integer-type-in-go
func (bit *cBool) UnmarshalJSON(inBits []byte) error {
	inString := strings.ToLower(string(inBits))

	inInt, err := strconv.Atoi(inString)
	if err != nil {
		return fmt.Errorf("failed injesting bool: %s", inString)
	}

	switch {
	case inInt > 0 || inString == "true":
		*bit = true
	case inInt == 0 || inString == "false":
		*bit = false
	default:
		return fmt.Errorf("cannot unmarshal boolean: %s", inString)
	}

	return nil
}

// NewTicketAuthenticator returns a new ticket authenticator.
// precedence:  AuthPayload > AuthTicket & CSRFPreventionToken > OTP > user+pass // aka:  pre-auth > OTP > user+pass
func NewTicketAuthenticator(conn *Connection, creds *Credentials) (Authenticator, error) {
	if creds.AuthPayload != "" {
		resBody := AuthenticationResponseBody{}

		err := json.Unmarshal([]byte(creds.AuthPayload), &resBody)
		if err != nil {
			return nil, fmt.Errorf("failed to decode auth_payload json, %w", err)
		}

		//nolint:lll
		if resBody.Data == nil || (resBody.Data.Ticket == nil || *resBody.Data.Ticket == "") || (resBody.Data.CSRFPreventionToken == nil || *resBody.Data.CSRFPreventionToken == "") {
			// not enforcing Username, nor capabilities in authPayload
			return nil, fmt.Errorf("external auth_payload is invalid; eg: '{\"data\":{\"username\":\"USER@REALM\",\"ticket\":\"PVE:USER@REALM:UNIQUEID::BASE64BLOB\",\"CSRFPreventionToken\":\"UNIQUEID:BLOB\",\"cap\":{}}}}'")
		}

		if resBody.Data.NeedTFA != nil && *(resBody.Data.NeedTFA) {
			return nil, fmt.Errorf("external auth_payload contains key NeedTFA and requires two-factor authentication")
		}

		creds.Username = resBody.Data.Username                              // todo: exclude line? - is creds.Username needed/used anywhere other than new-auth reqs?
		if creds.Username != "" && !strings.Contains(creds.Username, "@") { // todo: improve this vs copy-pasta from credentials.go
			return nil, errors.New(
				"make sure the username for the Proxmox Virtual Environment API ends in '@pve or @pam'",
			)
		}

		return &ticketAuthenticator{
			conn:     conn,
			authData: resBody.Data,
		}, nil
	}

	if creds.AuthTicket != "" && creds.CSRFPreventionToken != "" {
		ard := &AuthenticationResponseData{}
		ard.Ticket = &(creds.AuthTicket)
		ard.CSRFPreventionToken = &(creds.CSRFPreventionToken)

		authTicketSplits := strings.Split(creds.AuthTicket, ":")

		if len(authTicketSplits) > 3 {
			creds.Username = strings.Split(creds.AuthTicket, ":")[1] // todo: exclude line? - is creds.Username needed/used anywhere other than new-auth reqs?
		} else {
			return nil, errors.New("auth_ticket is set to an invalid value")
		}

		if creds.Username != "" && !strings.Contains(creds.Username, "@") { // todo: improve this vs copy-pasta from credentials.go
			return nil, errors.New(
				"make sure the username for the Proxmox Virtual Environment API ends in '@pve or @pam'",
			)
		}

		return &ticketAuthenticator{
			conn:     conn,
			authData: ard,
		}, nil
	}

	authRequest := fmt.Sprintf(
		"username=%s&password=%s",
		url.QueryEscape(creds.Username),
		url.QueryEscape(creds.Password),
	)

	// OTP is optional, and probably doesn't make much sense for most provider users.
	//   TOTP uses 2x requests; one with payloads `username=` and `password=`,
	//     (this returns a payload including: 'NeedTFA=1')
	//   followed by a 2nd request with payloads:
	//     `username=`, `tfa-challenge=<firsts response ticket>`, `password=totp:######`,
	//   and header: `CSRFPreventionToken: <first response CSRF>`
	//   Ticket generated lasts for ~2hours (to verify)
	if creds.OTP != nil {
		authRequest = fmt.Sprintf("%s&otp=%s", authRequest, url.QueryEscape(*creds.OTP))
	}

	return &ticketAuthenticator{
		conn:        conn,
		authRequest: authRequest,
	}, nil
}

func (t *ticketAuthenticator) authenticate(ctx context.Context) (*AuthenticationResponseData, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

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
		HttpOnly: true,
		Name:     "PVEAuthCookie",
		Secure:   true,
		Value:    *a.Ticket,
	})

	if req.Method != http.MethodGet {
		req.Header.Add("CSRFPreventionToken", *a.CSRFPreventionToken)
	}

	return nil
}
