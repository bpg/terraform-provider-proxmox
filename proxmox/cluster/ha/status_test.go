/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ha

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

func newStatusTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()

	server := httptest.NewTLSServer(handler)
	t.Cleanup(server.Close)

	conn, err := api.NewConnection(server.URL, true, "", nil)
	require.NoError(t, err)

	creds, err := api.NewCredentials("", "", "", "user@pve!token=test", "", "")
	require.NoError(t, err)

	c, err := api.NewClient(creds, conn)
	require.NoError(t, err)

	return &Client{Client: c}
}

func testCTID() types.HAResourceID {
	return types.HAResourceID{Type: types.HAResourceTypeContainer, Name: "100"}
}

func writeStatus(w http.ResponseWriter, body string) {
	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write([]byte(body)); err != nil {
		panic(err)
	}
}

func TestWaitForServiceUnmanagedReturnsOnceServiceDropped(t *testing.T) {
	t.Parallel()

	var calls atomic.Int32

	client := newStatusTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api2/json/cluster/ha/status/manager_status", r.URL.Path)

		if calls.Add(1) == 1 {
			writeStatus(w, `{"data":{"manager_status":{"service_status":{"ct:100":{"state":"started"}}}}}`)
			return
		}

		writeStatus(w, `{"data":{"manager_status":{"service_status":{}}}}`)
	})

	require.NoError(t, client.WaitForServiceUnmanaged(t.Context(), testCTID()))
	assert.Equal(t, int32(2), calls.Load())
}

func TestWaitForServiceUnmanagedSkipsWhileDisarmed(t *testing.T) {
	t.Parallel()

	var calls atomic.Int32

	client := newStatusTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		writeStatus(w, `{"data":{"manager_status":{"disarm":{"mode":"freeze","state":"disarmed"},`+
			`"service_status":{"ct:100":{"state":"freeze"}}}}}`)
	})

	require.NoError(t, client.WaitForServiceUnmanaged(t.Context(), testCTID()))
	assert.Equal(t, int32(1), calls.Load())
}

func TestWaitForServiceUnmanagedStopsOnCallerDeadline(t *testing.T) {
	t.Parallel()

	client := newStatusTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		writeStatus(w, `{"data":{"manager_status":{"service_status":{"ct:100":{"state":"started"}}}}}`)
	})

	ctx, cancel := context.WithTimeout(t.Context(), 1500*time.Millisecond)
	defer cancel()

	err := client.WaitForServiceUnmanaged(ctx, testCTID())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error waiting for HA manager to release ct:100")
}
