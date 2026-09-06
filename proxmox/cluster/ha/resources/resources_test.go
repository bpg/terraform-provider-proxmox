/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resources

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

func newTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	return httptest.NewTLSServer(mux)
}

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

// writeJSON writes a JSON response in test handlers. Panics on error since
// we're in a test context and can't use require (which would fail in goroutine).
func writeJSON(w http.ResponseWriter, v any) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		panic(err)
	}
}

func TestClient_Exists_ReturnsTrue_WhenHAResourceExists(t *testing.T) {
	t.Parallel()

	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/cluster/ha/resources/vm:100")

		resp := HAResourceGetResponseBody{
			Data: &HAResourceGetResponseData{
				ID:   types.HAResourceID{Type: types.HAResourceTypeVM, Name: "100"},
				Type: types.HAResourceTypeVM,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		writeJSON(w, resp)
	})
	defer server.Close()

	client := newTestClient(t, server.URL)

	haResourceID := types.HAResourceID{Type: types.HAResourceTypeVM, Name: "100"}

	exists, err := client.Exists(t.Context(), haResourceID)

	require.NoError(t, err)
	assert.True(t, exists)
}

func TestClient_Exists_ReturnsFalse_WhenHAResourceDoesNotExist(t *testing.T) {
	t.Parallel()

	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/cluster/ha/resources/vm:100")

		w.Header().Set("Content-Type", "application/json")
		// return 404 status which is recognized as "resource does not exist"
		w.WriteHeader(http.StatusNotFound)

		resp := map[string]any{
			"data": nil,
		}
		writeJSON(w, resp)
	})
	defer server.Close()

	client := newTestClient(t, server.URL)

	haResourceID := types.HAResourceID{Type: types.HAResourceTypeVM, Name: "100"}

	exists, err := client.Exists(t.Context(), haResourceID)

	require.NoError(t, err)
	assert.False(t, exists)
}

func TestClient_Migrate_SendsCorrectRequest(t *testing.T) {
	t.Parallel()

	var capturedPath string

	var capturedMethod string

	var capturedBody string

	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		capturedPath = r.URL.Path

		if r.Body != nil {
			if err := r.ParseForm(); err != nil {
				panic(err)
			}

			capturedBody = r.Form.Encode()
		}

		w.Header().Set("Content-Type", "application/json")

		// PVE 9.x response format
		resp := HAResourceMigrateResponseBody{
			Data: &HAResourceMigrateResponseData{
				SID:           "vm:100",
				RequestedNode: "pve2",
			},
		}
		writeJSON(w, resp)
	})
	defer server.Close()

	client := newTestClient(t, server.URL)

	haResourceID := types.HAResourceID{Type: types.HAResourceTypeVM, Name: "100"}

	resp, err := client.Migrate(t.Context(), haResourceID, "pve2")

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "vm:100", resp.SID)
	assert.Equal(t, "pve2", resp.RequestedNode)
	assert.Equal(t, http.MethodPost, capturedMethod)
	assert.Contains(t, capturedPath, "/cluster/ha/resources/vm:100/migrate")
	assert.Contains(t, capturedBody, "node=pve2")
}
