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

// TestWaitForTask_WithFailOnWarningsTreatsWarningsAsErrors verifies that when
// WithFailOnWarnings is used, a task completing with warnings returns an error.
func TestWaitForTask_WithFailOnWarningsTreatsWarningsAsErrors(t *testing.T) {
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

	// Task log: return warning lines (used by taskFailedResult when failOnWarnings is active).
	mux.HandleFunc("GET /api2/json/nodes/pve/tasks/"+testUPID+"/log", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		writeJSON(w, map[string]any{
			"data": []map[string]any{
				{"n": 1, "t": "WARN: some warning"},
				{"n": 2, "t": "TASK WARNINGS: 1"},
			},
		})
	})

	server := httptest.NewTLSServer(mux)
	defer server.Close()

	client := newTestClient(t, server.URL)

	result := client.WaitForTask(t.Context(), testUPID, WithFailOnWarnings())
	require.Error(t, result.Err(), "WaitForTask with WithFailOnWarnings should return error on warnings")
	assert.Contains(t, result.Err().Error(), "WARNINGS: 1", "error should include the exit code")

	// The WARN: line must also be surfaced as a separate warning so callers can emit it as a diagnostic.
	assert.True(t, result.HasWarnings(), "result should carry the warning lines")
	require.Len(t, result.Warnings(), 1)
	assert.Contains(t, result.Warnings()[0], "some warning")
}

// TestWaitForTask_FailedTaskWithWarningsDoesNotDuplicateWarnLines verifies that when a task
// fails and the log contains both error lines and WARN: lines, the WARN: text appears only in
// the separate warnings — not duplicated inside the error message.
func TestWaitForTask_FailedTaskWithWarningsDoesNotDuplicateWarnLines(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()

	// Task status: failed with an error exit code.
	mux.HandleFunc("GET /api2/json/nodes/pve/tasks/"+testUPID+"/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		writeJSON(w, map[string]any{
			"data": map[string]any{
				"status":     "stopped",
				"exitstatus": "ERROR: allocation failed",
			},
		})
	})

	// Task log: mixed error and warning lines.
	mux.HandleFunc("GET /api2/json/nodes/pve/tasks/"+testUPID+"/log", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		writeJSON(w, map[string]any{
			"data": []map[string]any{
				{"n": 1, "t": "starting task"},
				{"n": 2, "t": "WARN: disk is nearly full"},
				{"n": 3, "t": "ERROR: allocation failed"},
				{"n": 4, "t": "TASK ERROR: allocation failed"},
			},
		})
	})

	server := httptest.NewTLSServer(mux)
	defer server.Close()

	client := newTestClient(t, server.URL)

	result := client.WaitForTask(t.Context(), testUPID)
	require.Error(t, result.Err())

	// The error message should contain the non-warning log lines.
	assert.Contains(t, result.Err().Error(), "allocation failed")
	assert.Contains(t, result.Err().Error(), "starting task")

	// The WARN: line must NOT appear in the error detail — it is surfaced separately as a warning.
	assert.NotContains(t, result.Err().Error(), "disk is nearly full",
		"WARN: lines should be excluded from the error detail to avoid duplication")

	// The warning should be available as a separate diagnostic.
	assert.True(t, result.HasWarnings(), "result should carry warning lines from the task log")
	require.Len(t, result.Warnings(), 1)
	assert.Contains(t, result.Warnings()[0], "disk is nearly full")
}
