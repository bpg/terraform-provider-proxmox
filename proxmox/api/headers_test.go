/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	cfIDHeader     = "CF-Access-Client-Id"
	cfSecretHeader = "CF-Access-Client-Secret" // #nosec G101
	cfSecretValue  = "top-secret"
)

func testHeaders() map[string]string {
	return map[string]string{
		cfIDHeader:     "abc.access",
		cfSecretHeader: cfSecretValue,
	}
}

// newConnectionForServer points a Connection at a test server. NewConnection insists on https,
// so the http test server URL is injected afterwards, keeping the transport wiring under test.
func newConnectionForServer(t *testing.T, serverURL string, headers map[string]string) *Connection {
	t.Helper()

	u, err := url.Parse(serverURL)
	require.NoError(t, err)

	conn, err := NewConnection("https://"+u.Host, true, "", headers)
	require.NoError(t, err)

	ht, ok := conn.httpClient.Transport.(*headerTransport)
	require.True(t, ok, "expected the header transport to be the outermost one")

	ht.scheme = u.Scheme
	conn.endpoint = serverURL

	return conn
}

func TestHeaderTransportAddsHeadersToAPIRequests(t *testing.T) {
	t.Parallel()

	var got http.Header

	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		got = r.Header.Clone()
	}))
	defer server.Close()

	conn := newConnectionForServer(t, server.URL, testHeaders())

	c, err := NewClient(Credentials{TokenCredentials: &TokenCredentials{APIToken: "user@pve!tok=secret"}}, conn)
	require.NoError(t, err)

	require.NoError(t, c.DoRequest(t.Context(), http.MethodGet, "version", nil, nil))

	require.Equal(t, "abc.access", got.Get(cfIDHeader))
	require.Equal(t, cfSecretValue, got.Get(cfSecretHeader))
	// the provider's own authentication must survive the injection
	require.Equal(t, "PVEAPIToken=user@pve!tok=secret", got.Get("Authorization"))
}

// The http.Client is shared with callers that fetch arbitrary user-supplied URLs (see readURL in
// proxmoxtf/resource/file.go), so the credentials must never reach an unrelated host.
func TestHeaderTransportSkipsOtherHosts(t *testing.T) {
	t.Parallel()

	var got http.Header

	other := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		got = r.Header.Clone()
	}))
	defer other.Close()

	endpoint := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	defer endpoint.Close()

	conn := newConnectionForServer(t, endpoint.URL, testHeaders())

	req, err := http.NewRequestWithContext(t.Context(), http.MethodHead, other.URL, nil)
	require.NoError(t, err)

	res, err := conn.httpClient.Do(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())

	require.Empty(t, got.Get(cfIDHeader))
	require.Empty(t, got.Get(cfSecretHeader))
}

// Go populates a redirected request from a snapshot of the original request taken before any
// RoundTrip runs, so cloning inside RoundTrip keeps the headers out of a cross-host redirect.
func TestHeaderTransportDoesNotFollowHeadersAcrossRedirect(t *testing.T) {
	t.Parallel()

	var got http.Header

	other := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		got = r.Header.Clone()
	}))
	defer other.Close()

	endpoint := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, other.URL, http.StatusFound)
	}))
	defer endpoint.Close()

	conn := newConnectionForServer(t, endpoint.URL, testHeaders())

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, endpoint.URL, nil)
	require.NoError(t, err)

	res, err := conn.httpClient.Do(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())

	require.Empty(t, got.Get(cfIDHeader))
	require.Empty(t, got.Get(cfSecretHeader))
}

// http.Client follows an https -> http redirect on the same host, and Go's redirect filter only
// looks at the domain. Matching the scheme keeps the credentials off a plaintext connection.
func TestHeaderTransportSkipsSchemeDowngrade(t *testing.T) {
	t.Parallel()

	var got http.Header

	plaintext := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		got = r.Header.Clone()
	}))
	defer plaintext.Close()

	u, err := url.Parse(plaintext.URL)
	require.NoError(t, err)

	// the transport is configured for https on the very host the plaintext server listens on
	conn, err := NewConnection("https://"+u.Host, true, "", testHeaders())
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, plaintext.URL, nil)
	require.NoError(t, err)

	res, err := conn.httpClient.Do(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())

	require.Empty(t, got.Get(cfIDHeader))
	require.Empty(t, got.Get(cfSecretHeader))
}

func TestHeaderTransportDoesNotMutateCallerRequest(t *testing.T) {
	t.Parallel()

	transport := &headerTransport{
		next: RoundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}
		}),
		scheme:  "https",
		host:    "pve.example.com",
		headers: testHeaders(),
	}

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://pve.example.com/api2/json/version", nil)
	require.NoError(t, err)

	res, err := transport.RoundTrip(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())

	require.Empty(t, req.Header.Get(cfIDHeader))
	require.Empty(t, req.Header.Get(cfSecretHeader))
}

func TestNewConnectionRejectsInvalidHeaders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		headers map[string]string
		wantErr string
	}{
		{"provider managed", map[string]string{"Authorization": "Bearer x"}, "managed by the provider"},
		{"provider managed lowercase", map[string]string{"authorization": "Bearer x"}, "managed by the provider"},
		{"provider managed cookie", map[string]string{"Cookie": "a=b"}, "managed by the provider"},
		{"provider managed csrf", map[string]string{"CSRFPreventionToken": "x"}, "managed by the provider"},
		{"provider managed content type", map[string]string{"Content-Type": "text/plain"}, "managed by the provider"},
		{"client managed encoding", map[string]string{"Accept-Encoding": "gzip"}, "managed by the provider"},
		{"client managed host", map[string]string{"Host": "other"}, "managed by the provider"},
		{"hop by hop", map[string]string{"Connection": "close"}, "managed by the provider"},
		{"hop by hop upgrade", map[string]string{"Upgrade": "h2c"}, "managed by the provider"},
		{"proxy auth", map[string]string{"Proxy-Authorization": "Basic x"}, "managed by the provider"},
		{"empty name", map[string]string{"": "x"}, "must not be empty"},
		{"invalid name", map[string]string{"X Foo": "x"}, "invalid characters"},
		{"invalid value", map[string]string{"X-Foo": "a\r\nInjected: 1"}, "invalid characters"},
		{"case duplicate", map[string]string{"X-Foo": "a", "x-foo": "b"}, "differ only in case"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewConnection("https://pve.example.com", true, "", tt.headers)
			require.ErrorContains(t, err, tt.wantErr)
			// header values are credentials and must never be echoed back
			require.NotContains(t, err.Error(), "Bearer x")
		})
	}
}

func TestNewConnectionAcceptsCustomHeaders(t *testing.T) {
	t.Parallel()

	conn, err := NewConnection("https://pve.example.com", true, "", testHeaders())
	require.NoError(t, err)
	require.IsType(t, &headerTransport{}, conn.httpClient.Transport)
}

func TestNewConnectionWithoutHeadersLeavesTransportUnwrapped(t *testing.T) {
	t.Parallel()

	conn, err := NewConnection("https://pve.example.com", true, "", nil)
	require.NoError(t, err)

	_, wrapped := conn.httpClient.Transport.(*headerTransport)
	require.False(t, wrapped)
}
