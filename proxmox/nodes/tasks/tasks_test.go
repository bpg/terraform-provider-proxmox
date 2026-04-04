/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tasks

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// testUPID is a valid Proxmox UPID for use in tests.
const testUPID = "UPID:pve:00001234:00005678:AABBCCDD:vzcreate:100:root@pam:"

func newTestClient(t *testing.T, endpoint string) *Client {
	t.Helper()

	conn, err := api.NewConnection(endpoint, true, "")
	require.NoError(t, err)

	creds, err := api.NewCredentials("", "", "", "user@pve!token=test", "", "")
	require.NoError(t, err)

	c, err := api.NewClient(creds, conn)
	require.NoError(t, err)

	return &Client{Client: c}
}

func writeJSON(w http.ResponseWriter, v any) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		panic(err)
	}
}

// TestWaitForTask_FailedTaskIncludesLog verifies that when a task fails with an
// actual error (not warnings), the error message includes the task log output.
func TestWaitForTask_FailedTaskIncludesLog(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()

	// Task status: completed with an error exit code.
	mux.HandleFunc("GET /api2/json/nodes/pve/tasks/"+testUPID+"/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		writeJSON(w, map[string]any{
			"data": map[string]any{
				"status":     "stopped",
				"exitstatus": "ERROR: command failed",
			},
		})
	})

	// Task log: return lines including the error.
	mux.HandleFunc("GET /api2/json/nodes/pve/tasks/"+testUPID+"/log", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		writeJSON(w, map[string]any{
			"data": []map[string]any{
				{"n": 1, "t": "starting task"},
				{"n": 2, "t": "TASK ERROR: command failed"},
			},
		})
	})

	server := httptest.NewTLSServer(mux)
	defer server.Close()

	client := newTestClient(t, server.URL)

	result := client.WaitForTask(t.Context(), testUPID)
	require.Error(t, result.Err(), "WaitForTask should return error for non-OK exit code")

	assert.Contains(t, result.Err().Error(), "command failed")
}

// TestWaitForTask_FailedTaskWithLogFetchError verifies that when a task fails
// and fetching the log also fails, the original error is still returned.
func TestWaitForTask_FailedTaskWithLogFetchError(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()

	// Task status: completed with error.
	mux.HandleFunc("GET /api2/json/nodes/pve/tasks/"+testUPID+"/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		writeJSON(w, map[string]any{
			"data": map[string]any{
				"status":     "stopped",
				"exitstatus": "ERROR",
			},
		})
	})

	// Task log: return error (log endpoint unavailable).
	mux.HandleFunc("GET /api2/json/nodes/pve/tasks/"+testUPID+"/log", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	server := httptest.NewTLSServer(mux)
	defer server.Close()

	client := newTestClient(t, server.URL)

	result := client.WaitForTask(t.Context(), testUPID)
	require.Error(t, result.Err())

	// Should still have the basic error info even if log fetch failed.
	assert.Contains(t, result.Err().Error(), "failed to complete with exit code: ERROR")
}

// TestWaitForTask_WarningsAreNonFatalByDefault verifies that when a task completes
// with warnings, WaitForTask succeeds and the warning text is available.
func TestWaitForTask_WarningsAreNonFatalByDefault(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()

	// Task status: completed with warnings.
	mux.HandleFunc("GET /api2/json/nodes/pve/tasks/"+testUPID+"/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		writeJSON(w, map[string]any{
			"data": map[string]any{
				"status":     "stopped",
				"exitstatus": "WARNINGS: 1",
			},
		})
	})

	// Task log: verify it gets fetched even when warnings are ignored.
	mux.HandleFunc("GET /api2/json/nodes/pve/tasks/"+testUPID+"/log", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		writeJSON(w, map[string]any{
			"data": []map[string]any{
				{"n": 1, "t": "WARN: Systemd 258 detected. You may need to enable nesting."},
				{"n": 2, "t": "TASK WARNINGS: 1"},
			},
		})
	})

	server := httptest.NewTLSServer(mux)
	defer server.Close()

	client := newTestClient(t, server.URL)

	result := client.WaitForTask(t.Context(), testUPID)
	require.NoError(t, result.Err(), "WaitForTask should succeed — warnings are non-fatal by default")

	assert.True(t, result.HasWarnings(), "result should carry warning lines")
	assert.Len(t, result.Warnings(), 1, "TASK WARNINGS summary line should be filtered out")
	assert.Contains(t, result.Warnings()[0], "Systemd 258")
}
