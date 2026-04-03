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

// TestWaitForTask_FailedTaskIncludesLog verifies that when a task fails,
// the error message includes the task log output so users can see what went wrong.
func TestWaitForTask_FailedTaskIncludesLog(t *testing.T) {
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

	// Task log: return lines including the warning.
	mux.HandleFunc("GET /api2/json/nodes/pve/tasks/"+testUPID+"/log", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		writeJSON(w, map[string]any{
			"data": []map[string]any{
				{"n": 1, "t": "extracting archive '/var/lib/vz/template/cache/fedora.tar.xz'"},
				{"n": 2, "t": "Total bytes read: 594892800"},
				{"n": 3, "t": "Creating SSH host key 'ssh_host_ed25519_key'"},
				{"n": 4, "t": "WARN: Systemd 258 detected. You may need to enable nesting."},
				{"n": 5, "t": "TASK WARNINGS: 1"},
			},
		})
	})

	server := httptest.NewTLSServer(mux)
	defer server.Close()

	client := newTestClient(t, server.URL)

	err := client.WaitForTask(t.Context(), testUPID)
	require.Error(t, err, "WaitForTask should return error for non-OK exit code")

	// The error message should include the task log lines.
	assert.Contains(t, err.Error(), "Systemd 258 detected")
	assert.Contains(t, err.Error(), "TASK WARNINGS: 1")
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

	err := client.WaitForTask(t.Context(), testUPID)
	require.Error(t, err)

	// Should still have the basic error info even if log fetch failed.
	assert.Contains(t, err.Error(), "failed to complete with exit code: ERROR")
}

// TestWaitForTask_IgnoredWarningsIncludesLogInContext verifies that when
// WithIgnoreWarnings is set and the task has warnings, WaitForTask succeeds
// but the warning text is available (not silently swallowed).
func TestWaitForTask_IgnoredWarningsIncludesLogInContext(t *testing.T) {
	t.Parallel()

	var logRequested bool

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
		logRequested = true

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

	err := client.WaitForTask(t.Context(), testUPID, WithIgnoreWarnings())
	require.NoError(t, err, "WaitForTask should succeed when ignoring warnings")

	assert.True(t, logRequested, "task log should be fetched to surface warnings even when ignoring them")
}
