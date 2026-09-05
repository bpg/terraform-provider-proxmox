/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vnets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	proxmoxfirewall "github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type capturedFirewallRequest struct {
	method string
	path   string
	form   url.Values
}

func TestFirewallClientUsesVNetRuleEndpoints(t *testing.T) {
	t.Parallel()

	requests := make([]capturedFirewallRequest, 0)

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		requests = append(requests, capturedFirewallRequest{
			method: r.Method,
			path:   r.URL.Path,
			form:   r.Form,
		})

		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api2/json/cluster/sdn/vnets/testv/firewall/rules":
			writeFirewallTestJSON(w, map[string]any{
				"data": []map[string]any{{"pos": 0}},
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api2/json/cluster/sdn/vnets/testv/firewall/rules/0":
			writeFirewallTestJSON(w, map[string]any{
				"data": map[string]any{
					"action": "ACCEPT",
					"enable": 1,
					"pos":    "0",
					"type":   "forward",
				},
			})
		default:
			writeFirewallTestJSON(w, map[string]any{"data": nil})
		}
	}))
	defer server.Close()

	client := newFirewallTestClient(t, server.URL)
	rulesAPI := client.Firewall()
	enabled := proxmoxtypes.CustomBool(true)

	require.NoError(t, rulesAPI.CreateRule(t.Context(), &proxmoxfirewall.RuleCreateRequestBody{
		BaseRule: proxmoxfirewall.BaseRule{
			Enable:   &enabled,
			ICMPType: new("echo-request"),
		},
		Action: "ACCEPT",
		Type:   "forward",
	}))

	positions, err := rulesAPI.ListRules(t.Context())
	require.NoError(t, err)
	require.Equal(t, 0, positions[0].Pos)

	rule, err := rulesAPI.GetRule(t.Context(), 0)
	require.NoError(t, err)
	require.Equal(t, "ACCEPT", rule.Action)

	moveTo := 1
	require.NoError(t, rulesAPI.UpdateRule(t.Context(), 0, &proxmoxfirewall.RuleUpdateRequestBody{MoveTo: &moveTo}))
	require.NoError(t, rulesAPI.DeleteRule(t.Context(), 0))

	require.Len(t, requests, 5)
	require.Equal(t, http.MethodPost, requests[0].method)
	require.Equal(t, "/api2/json/cluster/sdn/vnets/testv/firewall/rules", requests[0].path)
	require.Equal(t, "ACCEPT", requests[0].form.Get("action"))
	require.Equal(t, "1", requests[0].form.Get("enable"))
	require.Equal(t, "echo-request", requests[0].form.Get("icmp-type"))
	require.Equal(t, "forward", requests[0].form.Get("type"))
	require.Equal(t, http.MethodGet, requests[1].method)
	require.Equal(t, "/api2/json/cluster/sdn/vnets/testv/firewall/rules", requests[1].path)
	require.Equal(t, http.MethodGet, requests[2].method)
	require.Equal(t, "/api2/json/cluster/sdn/vnets/testv/firewall/rules/0", requests[2].path)
	require.Equal(t, http.MethodPut, requests[3].method)
	require.Equal(t, "1", requests[3].form.Get("moveto"))
	require.Equal(t, http.MethodDelete, requests[4].method)
	require.Equal(t, "/api2/json/cluster/sdn/vnets/testv/firewall/rules/0", requests[4].path)
}

func newFirewallTestClient(t *testing.T, endpoint string) *Client {
	t.Helper()

	connection, err := api.NewConnection(endpoint, true, "")
	require.NoError(t, err)

	credentials, err := api.NewCredentials("", "", "", "user@pve!token=test", "", "")
	require.NoError(t, err)

	apiClient, err := api.NewClient(credentials, connection)
	require.NoError(t, err)

	return &Client{
		Client: firewallTestClusterClient{Client: apiClient},
		ID:     "testv",
	}
}

type firewallTestClusterClient struct {
	api.Client
}

func (c firewallTestClusterClient) ExpandPath(path string) string {
	return fmt.Sprintf("cluster/%s", path)
}

func writeFirewallTestJSON(w http.ResponseWriter, value any) {
	if err := json.NewEncoder(w).Encode(value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
