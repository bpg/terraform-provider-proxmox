//go:build acceptance || all

//testacc:tier=light
//testacc:resource=misc

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"net/http"
	"net/http/httptrace"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestAccAPIClientDropsIdleConnectionsBeforePveproxy verifies that the API client doesn't reuse a connection that
// has been idle long enough for pveproxy (5 s) to be closing it: a PUT written to it would fail with a bare EOF.
func TestAccAPIClientDropsIdleConnectionsBeforePveproxy(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)
	client := te.Client()

	getVersion := func(t *testing.T) bool {
		t.Helper()

		var reused bool

		ctx := httptrace.WithClientTrace(t.Context(), &httptrace.ClientTrace{
			GotConn: func(info httptrace.GotConnInfo) { reused = info.Reused },
		})

		var res map[string]any
		require.NoError(t, client.DoRequest(ctx, http.MethodGet, "version", nil, &res))

		return reused
	}

	getVersion(t)
	time.Sleep(100 * time.Millisecond)
	require.True(t, getVersion(t), "a connection idle for 100 ms should be reused")

	time.Sleep(4 * time.Second)
	require.False(t, getVersion(t), "a connection idle for 4 s should not be reused")
}
