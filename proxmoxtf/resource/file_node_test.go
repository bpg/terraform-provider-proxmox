/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

func TestFileResolveNode(t *testing.T) {
	t.Parallel()

	const shared = `{"data":[{"storage":"shared","shared":1,"active":1,"enabled":1,"content":"snippets"}]}`
	const nodes = `{"data":[{"node":"b","status":"online"},{"node":"a","status":"online"}]}`
	tests := []struct {
		name        string
		preferred   string
		nodes       string
		config      string
		storage     string
		firstStatus int
		nodesStatus int
		want        string
		wantError   string
	}{
		{name: "existing preferred node", preferred: "b", want: "b"},
		{name: "existing offline node", preferred: "b", nodes: `{"data":[{"node":"b","status":"offline"}]}`, want: "b"},
		{name: "deterministic shared fallback", want: "a"},
		{name: "node restriction", config: `{"data":{"nodes":"b"}}`, want: "b"},
		{name: "restricted to removed node", config: `{"data":{"nodes":"removed"}}`, wantError: "no online node"},
		{name: "disabled datastore", config: `{"data":{"disable":1}}`, wantError: "disabled"},
		{name: "no nodes", nodes: `{"data":[]}`, wantError: "no online node"},
		{name: "offline candidate", nodes: `{"data":[{"node":"a","status":"offline"}]}`, wantError: "no online node"},
		{name: "local datastore", storage: `{"data":[{"storage":"shared","shared":0,"active":1,"enabled":1,"content":"snippets"}]}`, wantError: "no online node"},
		{name: "unknown sharing", storage: `{"data":[{"storage":"shared","active":1,"enabled":1,"content":"snippets"}]}`, wantError: "no online node"},
		{name: "inactive datastore", storage: `{"data":[{"storage":"shared","shared":1,"active":0,"enabled":1,"content":"snippets"}]}`, wantError: "no online node"},
		{name: "disabled on node", storage: `{"data":[{"storage":"shared","shared":1,"active":1,"enabled":0,"content":"snippets"}]}`, wantError: "no online node"},
		{name: "unsupported content", storage: `{"data":[{"storage":"shared","shared":1,"active":1,"enabled":1,"content":"iso"}]}`, wantError: "no online node"},
		{name: "datastore missing", storage: `{"data":[]}`, wantError: "no online node"},
		{name: "candidate error tries next node", firstStatus: http.StatusForbidden, want: "b"},
		{name: "candidate error retained", firstStatus: http.StatusForbidden, nodes: `{"data":[{"node":"a","status":"online"}]}`, wantError: `node "a"`},
		{name: "discovery permission error", nodesStatus: http.StatusForbidden, wantError: "403"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")

				body := tt.storage
				if body == "" {
					body = shared
				}

				switch r.URL.Path {
				case "/api2/json/nodes":
					if tt.nodesStatus != 0 {
						w.WriteHeader(tt.nodesStatus)
					}

					body = tt.nodes
					if body == "" {
						body = nodes
					}
				case "/api2/json/storage/shared":
					assert.Empty(t, tt.preferred, "existing nodes must not trigger datastore discovery")

					body = tt.config
					if body == "" {
						body = `{"data":{}}`
					}
				case "/api2/json/nodes/a/storage":
					assert.Equal(t, "shared", r.URL.Query().Get("storage"))

					if tt.firstStatus != 0 {
						w.WriteHeader(tt.firstStatus)
					}
				case "/api2/json/nodes/b/storage":
					assert.Equal(t, "shared", r.URL.Query().Get("storage"))
				default:
					t.Errorf("unexpected request: %s", r.URL.Path)
					w.WriteHeader(http.StatusNotFound)
				}

				_, err := fmt.Fprint(w, body)
				assert.NoError(t, err)
			}))
			t.Cleanup(server.Close)

			conn, err := api.NewConnection(server.URL, true, "", nil)
			require.NoError(t, err)
			creds, err := api.NewCredentials("", "", "", "user@pve!token=abcd", "", "")
			require.NoError(t, err)
			apiClient, err := api.NewClient(creds, conn)
			require.NoError(t, err)

			preferred := tt.preferred
			if preferred == "" {
				preferred = "removed"
			}

			contentType := "snippets"

			got, err := fileResolveNode(context.Background(), proxmox.NewClient(apiClient, nil, ""), preferred, "shared", &contentType)
			if tt.wantError != "" {
				require.ErrorContains(t, err, tt.wantError)
				assert.Empty(t, got)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
