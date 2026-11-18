/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cluster

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// fakeNextIDCluster mimics PVE: nextid reserves nothing.
func fakeNextIDCluster(t *testing.T, existing map[int]bool) *Client {
	t.Helper()

	var mu sync.Mutex

	lowestFree := func() int {
		id := 100
		for existing[id] {
			id++
		}

		return id
	}

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api2/json/cluster/nextid":
			if v := r.URL.Query().Get("vmid"); v != "" {
				id, err := strconv.Atoi(v)
				if err != nil {
					t.Errorf("bad vmid %q: %v", v, err)
					w.WriteHeader(http.StatusBadRequest)

					return
				}

				if existing[id] {
					w.WriteHeader(http.StatusBadRequest)
					_, _ = fmt.Fprintf(w, `{"data":null,"errors":{"vmid":"VM %d already exists"}}`, id)

					return
				}

				_, _ = fmt.Fprintf(w, `{"data":"%d"}`, id)

				return
			}

			_, _ = fmt.Fprintf(w, `{"data":"%d"}`, lowestFree())
		case "/api2/json/cluster/resources":
			data := []map[string]any{}
			for id := range existing {
				data = append(data, map[string]any{"type": "lxc", "id": fmt.Sprintf("lxc/%d", id), "vmid": id})
			}

			if err := json.NewEncoder(w).Encode(map[string]any{"data": data}); err != nil {
				t.Errorf("encode resources: %v", err)
			}
		default:
			t.Errorf("unexpected request %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	conn, err := api.NewConnection(server.URL, true, "", nil)
	require.NoError(t, err)

	creds, err := api.NewCredentials("", "", "", "user@pve!token=test", "", "")
	require.NoError(t, err)

	c, err := api.NewClient(creds, conn)
	require.NoError(t, err)

	return &Client{Client: c}
}

func TestIDGeneratorSequentialSkipsTakenIDsWithoutReuse(t *testing.T) {
	t.Setenv("TMPDIR", t.TempDir())

	existing := map[int]bool{128: true, 129: true}
	for id := 100; id <= 126; id++ {
		existing[id] = true
	}

	gen := NewIDGenerator(fakeNextIDCluster(t, existing), IDGeneratorConfig{RandomIDs: false})

	got := make([]int, 0, 4)

	for range 4 {
		id, err := gen.NextID(t.Context())
		require.NoError(t, err)

		got = append(got, id)
	}

	require.Equal(t, []int{127, 130, 131, 132}, got)
}
