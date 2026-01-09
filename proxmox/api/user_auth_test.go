/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package api

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserAuthenticator_TFARequired(t *testing.T) {
	t.Parallel()

	// TFA response mock based on actual Proxmox API behavior
	tfaResponse := `{
		"data": {
			"NeedTFA": 1,
			"ticket": "PVE:!tfa!{totp:true}:12345678::testticket",
			"CSRFPreventionToken": "test-csrf-token",
			"username": "root@pam"
		}
	}`

	conn := &Connection{
		endpoint: "http://localhost",
		httpClient: newTestClient(func(_ *http.Request) *http.Response {
			return &http.Response{
				Status:     "200 OK",
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(tfaResponse)),
			}
		}),
	}

	auth := NewUserAuthenticator(UserCredentials{
		Username: "root@pam",
		Password: "test",
	}, conn)

	err := auth.AuthenticateRequest(t.Context(), &http.Request{
		Method: http.MethodGet,
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "two-factor authentication is required")
	require.Contains(t, err.Error(), "otp")
}

func TestUserAuthenticator_Success(t *testing.T) {
	t.Parallel()

	successResponse := `{
		"data": {
			"ticket": "PVE:root@pam:12345678::validticket",
			"CSRFPreventionToken": "test-csrf-token",
			"username": "root@pam"
		}
	}`

	conn := &Connection{
		endpoint: "http://localhost",
		httpClient: newTestClient(func(_ *http.Request) *http.Response {
			return &http.Response{
				Status:     "200 OK",
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(successResponse)),
			}
		}),
	}

	auth := NewUserAuthenticator(UserCredentials{
		Username: "root@pam",
		Password: "test",
	}, conn)

	req := &http.Request{
		Method: http.MethodGet,
		Header: make(http.Header),
	}
	err := auth.AuthenticateRequest(t.Context(), req)

	require.NoError(t, err)
}
