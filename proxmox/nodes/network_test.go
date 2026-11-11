/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

func newTestNodeClient(t *testing.T, handler http.Handler) *Client {
	t.Helper()

	server := httptest.NewTLSServer(handler)
	t.Cleanup(server.Close)

	conn, err := api.NewConnection(server.URL, true, "", nil)
	require.NoError(t, err)

	creds, err := api.NewCredentials("", "", "", "user@pve!token=test", "", "")
	require.NoError(t, err)

	c, err := api.NewClient(creds, conn)
	require.NoError(t, err)

	return &Client{Client: c, NodeName: "pve"}
}

func TestDeleteNetworkInterfaceNotFound(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		status   int
		body     string
		notFound bool
	}{
		{
			name:     "missing interface is a 400 parameter error",
			status:   http.StatusBadRequest,
			body:     `{"data":null,"errors":{"iface":"interface does not exist\n"},"message":"Parameter verification failed.\n"}`,
			notFound: true,
		},
		{
			name:     "other 400 parameter errors are not not-found",
			status:   http.StatusBadRequest,
			body:     `{"data":null,"errors":{"iface":"invalid format - value does not look like a valid interface name\n"}}`,
			notFound: false,
		},
		{
			name:     "500 without a not-found message is not not-found",
			status:   http.StatusInternalServerError,
			body:     `{"data":null,"message":"unable to write file\n"}`,
			notFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := newTestNodeClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete || r.URL.Path != "/api2/json/nodes/pve/network/vmbr9" {
					http.NotFound(w, r)

					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.status)

				if _, err := w.Write([]byte(tt.body)); err != nil {
					panic(err)
				}
			}))

			err := client.DeleteNetworkInterface(context.Background(), "vmbr9")
			require.Error(t, err)
			require.Equal(t, tt.notFound, errors.Is(err, api.ErrResourceDoesNotExist))
		})
	}
}
