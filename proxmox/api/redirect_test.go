/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

const redirectTestToken = "root@pam!probe=00000000-0000-0000-0000-000000000000"

func newRedirectTestClient(t *testing.T, endpoint string) Client {
	t.Helper()

	conn, err := NewConnection(endpoint, true, "", nil)
	require.NoError(t, err)

	creds, err := NewCredentials("", "", "", redirectTestToken, "", "")
	require.NoError(t, err)

	client, err := NewClient(creds, conn)
	require.NoError(t, err)

	return client
}

func TestClientRefusesEndpointRedirectToHTTP(t *testing.T) {
	t.Parallel()

	var plainHits atomic.Int32

	plain := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		plainHits.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(plain.Close)

	endpoint := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, plain.URL+r.URL.Path, http.StatusFound)
	}))
	t.Cleanup(endpoint.Close)

	client := newRedirectTestClient(t, endpoint.URL)

	err := client.DoRequest(context.Background(), http.MethodGet, "version", nil, nil)
	require.ErrorContains(t, err, "non-HTTPS")
	require.Zero(t, plainHits.Load())
}

func TestClientFollowsEndpointRedirectOverHTTPS(t *testing.T) {
	t.Parallel()

	var gotAuth atomic.Value

	endpoint := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/moved" {
			http.Redirect(w, r, "/moved", http.StatusFound)

			return
		}

		gotAuth.Store(r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(endpoint.Close)

	client := newRedirectTestClient(t, endpoint.URL)

	require.NoError(t, client.DoRequest(context.Background(), http.MethodGet, "version", nil, nil))
	require.Equal(t, "PVEAPIToken="+redirectTestToken, gotAuth.Load())
}

func TestClientFollowsNonEndpointRedirectToHTTP(t *testing.T) {
	t.Parallel()

	plain := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(plain.Close)

	mirror := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, plain.URL, http.StatusFound)
	}))
	t.Cleanup(mirror.Close)

	client := newRedirectTestClient(t, "https://pve.example.com:8006")

	req, err := http.NewRequestWithContext(context.Background(), http.MethodHead, mirror.URL, nil)
	require.NoError(t, err)

	res, err := client.HTTP().Do(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, http.StatusOK, res.StatusCode)
}
