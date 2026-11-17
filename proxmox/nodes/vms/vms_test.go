/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vms

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// TestCreateVMRetry checks which failed creates are dispatched again, with the response bodies PVE sends.
func TestCreateVMRetry(t *testing.T) {
	t.Parallel()

	type response struct {
		code    int
		message string
	}

	alreadyExists := response{http.StatusInternalServerError, "unable to create VM 100 - VM 100 already exists on node 'pve'\n"}

	tests := []struct {
		name      string
		responses []response
		wantErr   string
		wantPosts int32
	}{
		{
			// PVE has written the VM's config with "lock: create" and killed the worker before it ran
			name:      "no worker upid is not retried",
			responses: []response{{http.StatusInternalServerError, "got no worker upid - start worker failed\n"}, alreadyExists},
			wantErr:   "got no worker upid",
			wantPosts: 1,
		},
		{
			// the first create may have gone through, so "already exists" on the retry counts as done
			name:      "proxy timeout is retried",
			responses: []response{{596, ""}, alreadyExists},
			wantPosts: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var posts atomic.Int32

			mux := http.NewServeMux()
			mux.HandleFunc("POST /api2/json/qemu", func(w http.ResponseWriter, _ *http.Request) {
				res := tt.responses[min(int(posts.Add(1)), len(tt.responses))-1]

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(res.code)

				if res.message != "" {
					assert.NoError(t, json.NewEncoder(w).Encode(map[string]any{"data": nil, "message": res.message}))
				}
			})

			server := httptest.NewTLSServer(mux)
			defer server.Close()

			conn, err := api.NewConnection(server.URL, true, "", nil)
			require.NoError(t, err)

			creds, err := api.NewCredentials("", "", "", "user@pve!token=test", "", "")
			require.NoError(t, err)

			apiClient, err := api.NewClient(creds, conn)
			require.NoError(t, err)

			client := &Client{Client: apiClient}

			err = client.CreateVM(t.Context(), &CreateRequestBody{VMID: 100}).Err()

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantPosts, posts.Load())
		})
	}
}
