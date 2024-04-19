package api

import (
	"context"
	"errors"
	"net/http"
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

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func newTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

type dummyAuthenticator struct{}

func (dummyAuthenticator) IsRoot() bool {
	return false
}

func (dummyAuthenticator) IsRootTicket() bool {
	return false
}

func (dummyAuthenticator) AuthenticateRequest(ctx context.Context, req *http.Request) error {
	return nil
}

//func Test_client_DoRequest(t *testing.T) {
//
//
//	c := client{
//		conn: &Connection{
//			endpoint: "http://localhost",
//			httpClient: newTestClient(func(req *http.Request) *http.Response {
//				return &http.Response{
//					StatusCode: 200,
//					Body:       nil,
//				}
//			}),
//		},
//		auth: dummyAuthenticator{},
//	}
//
//	c.DoRequest(context.Background(), http.MethodGet, "/test", nil, nil)
//}

func Test_client_DoRequest(t *testing.T) {
	tests := []struct {
		name    string
		status  string
		wantErr error
	}{
		{name: "no error", status: "200 OK", wantErr: nil},
		{name: "not exists - 404 status", status: "404 missing", wantErr: ErrResourceDoesNotExist},
		{name: "not exists - 500 status", status: "500 This thing does not exist", wantErr: ErrResourceDoesNotExist},
		//{name: "500 status", status: "500 Internal Server Error", wantErr: HTTPError{}},
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

			err := c.DoRequest(context.Background(), "POST", "any", nil, nil)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DoRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
