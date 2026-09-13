/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// newStorageTestClient spins up a mock PVE API server and returns a storage
// client pointed at it.
func newStorageTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()

	server := httptest.NewTLSServer(handler)
	t.Cleanup(server.Close)

	conn, err := api.NewConnection(server.URL, true, "", nil)
	require.NoError(t, err)

	apiClient, err := api.NewClient(
		api.Credentials{TokenCredentials: &api.TokenCredentials{APIToken: "test@pve!test=secret"}},
		conn,
	)
	require.NoError(t, err)

	return &Client{Client: apiClient, StorageName: "test"}
}

// writeJSON writes a JSON response in test handlers. Panics on error since
// we're in a test context and can't use require (which would fail in goroutine).
func writeJSON(w http.ResponseWriter, v any) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		panic(err)
	}
}

func TestDeleteDatastoreFile_WaitsForTask(t *testing.T) {
	t.Parallel()

	const upid = "UPID:test:00000001:00000000:00000001:imgdel:test:test@pve!test:"
	const volid = "test:snippets/foo.yaml"

	deleteHandled := false
	statusPolls := 0

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodDelete && r.URL.Path == "/api2/json/storage/test/content/"+volid:
			// PVE returns the UPID of the async imgdel task immediately.
			w.Header().Set("Content-Type", "application/json;charset=UTF-8")
			writeJSON(w, map[string]any{"data": upid})

			deleteHandled = true
		case r.Method == http.MethodGet && r.URL.Path == "/api2/json/nodes/test/tasks/"+upid+"/status":
			statusPolls++

			switch statusPolls {
			case 1:
				// First poll: task still running; the file would still be
				// present in a listing at this point.
				writeJSON(w, map[string]any{
					"data": map[string]any{"status": "running"},
				})
			default:
				// Task finished — from now on the file is really gone.
				writeJSON(w, map[string]any{
					"data": map[string]any{"status": "stopped", "exitstatus": "OK"},
				})
			}
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	})

	c := newStorageTestClient(t, handler)

	err := c.DeleteDatastoreFile(context.Background(), volid)
	require.NoError(t, err)

	assert.True(t, deleteHandled, "DELETE should have been issued")
	assert.GreaterOrEqual(t, statusPolls, 2, "client should have polled the task status until completion")
}

func TestDeleteDatastoreFile_NoTaskInResponse(t *testing.T) {
	t.Parallel()

	const volid = "test:snippets/foo.yaml"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodDelete && r.URL.Path == "/api2/json/storage/test/content/"+volid:
			// Some PVE versions / paths return no task; deletion must still
			// succeed without polling.
			w.Header().Set("Content-Type", "application/json;charset=UTF-8")
			writeJSON(w, map[string]any{"data": nil})
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	})

	c := newStorageTestClient(t, handler)

	err := c.DeleteDatastoreFile(context.Background(), volid)
	require.NoError(t, err)
}
