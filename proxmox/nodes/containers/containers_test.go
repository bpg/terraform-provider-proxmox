/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package containers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/retry"
)

// testUPID is a valid Proxmox UPID for use in tests.
// Format: UPID:<node>:<pid_hex>:<pstart_hex>:<starttime_hex>:<type>:<id>:<user@realm>:.
const testUPID = "UPID:pve:00001234:00005678:AABBCCDD:vzcreate:100:root@pam:"

type requestCaptures struct {
	mu  sync.Mutex
	req []capturedRequest
}

type capturedRequest struct {
	Method string
	Path   string
}

func (c *requestCaptures) add(method, path string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.req = append(c.req, capturedRequest{Method: method, Path: path})
}

func (c *requestCaptures) countPOST(pathSuffix string) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	n := 0

	for _, r := range c.req {
		if r.Method == http.MethodPost && strings.HasSuffix(r.Path, pathSuffix) {
			n++
		}
	}

	return n
}

func newTestServer(t *testing.T, handler http.Handler) *httptest.Server {
	t.Helper()

	return httptest.NewTLSServer(handler)
}

func newTestClient(t *testing.T, endpoint string) *Client {
	t.Helper()

	conn, err := api.NewConnection(endpoint, true, "")
	require.NoError(t, err)

	creds, err := api.NewCredentials("", "", "", "user@pve!token=test", "", "")
	require.NoError(t, err)

	c, err := api.NewClient(creds, conn)
	require.NoError(t, err)

	return &Client{Client: c, VMID: 100}
}

// writeJSON writes a JSON response in test handlers. Panics on error since
// we're in a test context and can't use require (which would fail in goroutine).
func writeJSON(w http.ResponseWriter, v any) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		panic(err)
	}
}

// taskCompletedHandler returns a handler that responds with a completed task status.
func taskCompletedHandler(captures *requestCaptures) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		captures.add(r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		writeJSON(w, map[string]any{
			"data": map[string]any{
				"status":     "stopped",
				"exitstatus": "OK",
			},
		})
	}
}

func TestDeleteContainerWaitsForTask(t *testing.T) {
	t.Parallel()

	captures := &requestCaptures{}

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api2/json/lxc/100", func(w http.ResponseWriter, r *http.Request) {
		captures.add(r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		writeJSON(w, map[string]any{"data": testUPID})
	})
	mux.HandleFunc("GET /api2/json/nodes/", taskCompletedHandler(captures))

	server := newTestServer(t, mux)
	defer server.Close()

	client := newTestClient(t, server.URL)

	err := client.DeleteContainer(t.Context())
	require.NoError(t, err)
}

func TestResizeContainerDiskWaitsForTask(t *testing.T) {
	t.Parallel()

	captures := &requestCaptures{}

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api2/json/lxc/100/resize", func(w http.ResponseWriter, r *http.Request) {
		captures.add(r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		writeJSON(w, map[string]any{"data": testUPID})
	})
	mux.HandleFunc("GET /api2/json/nodes/", taskCompletedHandler(captures))

	server := newTestServer(t, mux)
	defer server.Close()

	client := newTestClient(t, server.URL)

	err := client.ResizeContainerDisk(t.Context(), &ResizeRequestBody{
		Disk: "rootfs",
		Size: "+1G",
	})
	require.NoError(t, err)
}

func TestIsTransientAPIError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		err     error
		want    bool
		comment string
	}{
		{
			name:    "HTTP 500 Internal Server Error",
			err:     &api.HTTPError{Code: http.StatusInternalServerError, Message: "Internal Server Error"},
			want:    true,
			comment: "HTTP 500 should be retried (API overload)",
		},
		{
			name:    "HTTP 503 Service Unavailable",
			err:     &api.HTTPError{Code: http.StatusServiceUnavailable, Message: "Service Unavailable"},
			want:    true,
			comment: "HTTP 503 should be retried (API overload)",
		},
		{
			name:    "HTTP 400 Bad Request",
			err:     &api.HTTPError{Code: http.StatusBadRequest, Message: "Bad Request"},
			want:    false,
			comment: "HTTP 400 should not be retried (client error)",
		},
		{
			name:    "HTTP 403 Forbidden",
			err:     &api.HTTPError{Code: http.StatusForbidden, Message: "Forbidden"},
			want:    false,
			comment: "HTTP 403 should not be retried (auth error)",
		},
		{
			name:    "HTTP 404 Not Found",
			err:     &api.HTTPError{Code: http.StatusNotFound, Message: "Not Found"},
			want:    false,
			comment: "HTTP 404 should not be retried (resource does not exist)",
		},
		{
			name:    "got no worker upid error",
			err:     fmt.Errorf("got no worker upid - start worker failed"),
			want:    true,
			comment: "PVE worker start failure should be retried",
		},
		{
			name:    "got timeout error",
			err:     fmt.Errorf("got timeout"),
			want:    true,
			comment: "timeout errors should be retried",
		},
		{
			name:    "wrapped got no worker upid error",
			err:     fmt.Errorf("error creating container: %w", fmt.Errorf("got no worker upid")),
			want:    true,
			comment: "wrapped PVE worker start failure should be retried",
		},
		{
			name:    "wrapped got timeout error",
			err:     fmt.Errorf("error waiting for task: %w", fmt.Errorf("got timeout")),
			want:    true,
			comment: "wrapped timeout errors should be retried",
		},
		{
			name:    "generic error",
			err:     errors.New("something went wrong"),
			want:    false,
			comment: "generic errors should not be retried",
		},
		{
			name:    "nil error",
			err:     nil,
			want:    false,
			comment: "nil error should not be retried",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := retry.IsTransientAPIError(tt.err)
			assert.Equal(t, tt.want, got, tt.comment)
		})
	}
}

// TestCreateContainerRetries verifies that CreateContainer retries on HTTP 500
// errors and eventually succeeds. The mock server returns 500 on the first POST
// to the create endpoint, then succeeds on the second attempt with a valid UPID,
// and returns a completed task status on the status poll.
func TestCreateContainerRetries(t *testing.T) {
	t.Parallel()

	captures := &requestCaptures{}
	var createCount int

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api2/json/lxc", func(w http.ResponseWriter, r *http.Request) {
		captures.add(r.Method, r.URL.Path)

		createCount++

		w.Header().Set("Content-Type", "application/json")

		if createCount == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			writeJSON(w, map[string]any{"data": nil})

			return
		}

		writeJSON(w, map[string]any{"data": testUPID})
	})
	mux.HandleFunc("GET /api2/json/nodes/", taskCompletedHandler(captures))

	server := newTestServer(t, mux)
	defer server.Close()

	client := newTestClient(t, server.URL)

	err := client.CreateContainer(t.Context(), &CreateRequestBody{})
	require.NoError(t, err)

	assert.Equal(t, 2, captures.countPOST("/lxc"),
		"expected exactly 2 POST calls (1 failure + 1 success), proving retry occurred")
}

// TestCreateContainerNoRetryOn400 verifies that CreateContainer does NOT retry
// on HTTP 400 errors. The mock server always returns 400, and we assert that
// only 1 call was made (no retry).
func TestCreateContainerNoRetryOn400(t *testing.T) {
	t.Parallel()

	captures := &requestCaptures{}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api2/json/lxc", func(w http.ResponseWriter, r *http.Request) {
		captures.add(r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]any{
			"errors": map[string]string{"ostemplate": "value does not exist"},
		})
	})

	server := newTestServer(t, mux)
	defer server.Close()

	client := newTestClient(t, server.URL)

	err := client.CreateContainer(t.Context(), &CreateRequestBody{})
	require.Error(t, err)

	assert.Equal(t, 1, captures.countPOST("/lxc"),
		"expected exactly 1 POST call (no retry on 400)")
}

// TestCloneContainerRetries verifies that CloneContainer retries on HTTP 500
// errors and eventually succeeds.
func TestCloneContainerRetries(t *testing.T) {
	t.Parallel()

	captures := &requestCaptures{}
	var cloneCount int

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api2/json/lxc/100/clone", func(w http.ResponseWriter, r *http.Request) {
		captures.add(r.Method, r.URL.Path)

		cloneCount++

		w.Header().Set("Content-Type", "application/json")

		if cloneCount == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			writeJSON(w, map[string]any{"data": nil})

			return
		}

		writeJSON(w, map[string]any{"data": testUPID})
	})
	mux.HandleFunc("GET /api2/json/nodes/", taskCompletedHandler(captures))

	server := newTestServer(t, mux)
	defer server.Close()

	client := newTestClient(t, server.URL)

	err := client.CloneContainer(t.Context(), &CloneRequestBody{})
	require.NoError(t, err)

	assert.Equal(t, 2, captures.countPOST("/clone"),
		"expected exactly 2 POST calls (1 failure + 1 success), proving retry occurred")
}

// containerStatusHandler returns a handler that responds with the given container status.
func containerStatusHandler(captures *requestCaptures, status string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		captures.add(r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		writeJSON(w, map[string]any{
			"data": map[string]any{
				"status": status,
			},
		})
	}
}

// TestStartContainerAlreadyRunningOnFirstAttempt verifies the TOCTOU race
// condition: the status pre-check says "stopped", but by the time the start
// API call is made, the container is already running. StartContainer should
// treat this as success (not an error).
func TestStartContainerAlreadyRunningOnFirstAttempt(t *testing.T) {
	t.Parallel()

	captures := &requestCaptures{}
	var statusCallCount int

	mux := http.NewServeMux()

	// Status endpoint: first call returns "stopped" (pre-check), subsequent calls return "running".
	mux.HandleFunc("GET /api2/json/lxc/100/status/current", func(w http.ResponseWriter, r *http.Request) {
		captures.add(r.Method, r.URL.Path)

		statusCallCount++

		w.Header().Set("Content-Type", "application/json")

		status := "running"
		if statusCallCount == 1 {
			status = "stopped"
		}

		writeJSON(w, map[string]any{
			"data": map[string]any{"status": status},
		})
	})

	// Start endpoint: returns "already running" error (race condition).
	mux.HandleFunc("POST /api2/json/lxc/100/status/start", func(w http.ResponseWriter, r *http.Request) {
		captures.add(r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]any{
			"errors": map[string]string{"status": "container is already running"},
		})
	})

	server := newTestServer(t, mux)
	defer server.Close()

	client := newTestClient(t, server.URL)

	err := client.StartContainer(t.Context())
	require.NoError(t, err, "StartContainer should succeed when API returns 'already running'")

	assert.Equal(t, 1, captures.countPOST("/status/start"),
		"expected exactly 1 start attempt")
}

// TestStartContainerAlreadyRunningDetectedByPreCheck verifies that
// StartContainer short-circuits when the status pre-check already shows
// "running" — no start API call should be made.
func TestStartContainerAlreadyRunningDetectedByPreCheck(t *testing.T) {
	t.Parallel()

	captures := &requestCaptures{}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api2/json/lxc/100/status/current",
		containerStatusHandler(captures, "running"))
	mux.HandleFunc("POST /api2/json/lxc/100/status/start", func(w http.ResponseWriter, r *http.Request) {
		captures.add(r.Method, r.URL.Path)
		t.Fatal("start API should not be called when container is already running")
	})

	server := newTestServer(t, mux)
	defer server.Close()

	client := newTestClient(t, server.URL)

	err := client.StartContainer(t.Context())
	require.NoError(t, err)

	assert.Equal(t, 0, captures.countPOST("/status/start"),
		"expected no start calls when container is already running")
}

// TestStartContainerRetriesOnNoWorkerUpid verifies that StartContainer retries
// when the API returns "got no worker upid" and eventually succeeds.
func TestStartContainerRetriesOnNoWorkerUpid(t *testing.T) {
	t.Parallel()

	captures := &requestCaptures{}
	var startCount int

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api2/json/lxc/100/status/current", func(w http.ResponseWriter, r *http.Request) {
		captures.add(r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")

		// Return "stopped" on pre-check, "running" on WaitForContainerStatus polls.
		status := "running"
		if startCount == 0 {
			status = "stopped"
		}

		writeJSON(w, map[string]any{
			"data": map[string]any{"status": status},
		})
	})

	mux.HandleFunc("POST /api2/json/lxc/100/status/start", func(w http.ResponseWriter, r *http.Request) {
		captures.add(r.Method, r.URL.Path)

		startCount++

		w.Header().Set("Content-Type", "application/json")

		if startCount == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			writeJSON(w, map[string]any{
				"errors": map[string]string{"status": "got no worker upid - start worker failed"},
			})

			return
		}

		writeJSON(w, map[string]any{"data": testUPID})
	})
	mux.HandleFunc("GET /api2/json/nodes/", taskCompletedHandler(captures))

	server := newTestServer(t, mux)
	defer server.Close()

	client := newTestClient(t, server.URL)

	err := client.StartContainer(t.Context())
	require.NoError(t, err)

	assert.Equal(t, 2, captures.countPOST("/status/start"),
		"expected exactly 2 start calls (1 failure + 1 success), proving retry occurred")
}

// TestStartContainerAlreadyRunningOnRetry verifies the scenario where the first
// start attempt gets "got no worker upid" (retried), and the second attempt gets
// "already running" because the first attempt actually succeeded. This should be
// treated as success.
func TestStartContainerAlreadyRunningOnRetry(t *testing.T) {
	t.Parallel()

	captures := &requestCaptures{}
	var startCount int

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api2/json/lxc/100/status/current",
		containerStatusHandler(captures, "stopped"))
	mux.HandleFunc("POST /api2/json/lxc/100/status/start", func(w http.ResponseWriter, r *http.Request) {
		captures.add(r.Method, r.URL.Path)

		startCount++

		w.Header().Set("Content-Type", "application/json")

		if startCount == 1 {
			// First attempt: transient error (retryable).
			w.WriteHeader(http.StatusInternalServerError)
			writeJSON(w, map[string]any{
				"errors": map[string]string{"status": "got no worker upid - start worker failed"},
			})

			return
		}

		// Second attempt: "already running" because first actually started it.
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]any{
			"errors": map[string]string{"status": "container is already running"},
		})
	})

	server := newTestServer(t, mux)
	defer server.Close()

	client := newTestClient(t, server.URL)

	err := client.StartContainer(t.Context())
	require.NoError(t, err, "should succeed when retry finds container already running")

	assert.Equal(t, 2, captures.countPOST("/status/start"),
		"expected exactly 2 start calls (1 transient failure + 1 already running)")
}

// TestCloneContainerNoRetryOn400 verifies that CloneContainer does NOT retry
// on HTTP 400 errors.
func TestCloneContainerNoRetryOn400(t *testing.T) {
	t.Parallel()

	captures := &requestCaptures{}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api2/json/lxc/100/clone", func(w http.ResponseWriter, r *http.Request) {
		captures.add(r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]any{
			"errors": map[string]string{"vmid": "VM 999 does not exist"},
		})
	})

	server := newTestServer(t, mux)
	defer server.Close()

	client := newTestClient(t, server.URL)

	err := client.CloneContainer(t.Context(), &CloneRequestBody{})
	require.Error(t, err)

	assert.Equal(t, 1, captures.countPOST("/clone"),
		"expected exactly 1 POST call (no retry on 400)")
}
