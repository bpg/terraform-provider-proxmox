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
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls.
func newTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

type dummyAuthenticator struct{}

func (dummyAuthenticator) IsRoot(_ context.Context) bool {
	return false
}

func (dummyAuthenticator) IsRootTicket(context.Context) bool {
	return false
}

func (dummyAuthenticator) AuthenticateRequest(_ context.Context, _ *http.Request) error {
	return nil
}

func TestClientDoRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		status  string
		wantErr error
	}{
		{name: "no error", status: "200 OK", wantErr: nil},
		{name: "not exists - 404 status", status: "404 missing", wantErr: ErrResourceDoesNotExist},
		{name: "not exists - 500 status", status: "500 This thing does not exist", wantErr: ErrResourceDoesNotExist},
		{name: "500 status", status: "500 Internal Server Error", wantErr: &HTTPError{
			Code:    500,
			Message: "Internal Server Error",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := client{
				conn: &Connection{
					endpoint: "http://localhost",
					httpClient: newTestClient(func(_ *http.Request) *http.Response {
						sc, err := strconv.Atoi(strings.Fields(tt.status)[0])
						require.NoError(t, err)
						return &http.Response{
							Status:     tt.status,
							StatusCode: sc,
							Body:       nil,
						}
					}),
				},
				auth: dummyAuthenticator{},
			}

			err := c.DoRequest(t.Context(), "POST", "any", nil, nil)
			fail := false

			switch {
			case err == nil && tt.wantErr == nil:
				return
			case err != nil && tt.wantErr == nil:
				fallthrough
			case err == nil && tt.wantErr != nil:
				fail = true
			default:
				var he, we *HTTPError
				if errors.As(err, &he) && errors.As(tt.wantErr, &we) {
					fail = !reflect.DeepEqual(he, we)
				} else {
					fail = !errors.Is(err, tt.wantErr)
				}
			}

			if fail {
				t.Errorf("DoRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
