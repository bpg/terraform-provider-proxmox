//go:build acceptance || all

//testacc:tier=heavy
//testacc:resource=network

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package network_test

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

// checkInterfaceActive asserts whether PVE reports an interface as active on the node.
//
// A staged change is listed but not up in the kernel, so `active` stays false until a reload.
// Terraform state alone would pass whether or not `reload` was honoured.
func checkInterfaceActive(te *test.Environment, iface string, want bool) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		ifaces, err := te.NodeClient().ListNetworkInterfaces(context.Background())
		if err != nil {
			return fmt.Errorf("listing network interfaces: %w", err)
		}

		for _, i := range ifaces {
			if i.Iface != iface {
				continue
			}

			got := i.Active != nil && bool(*i.Active)
			if got != want {
				return fmt.Errorf("interface %q: active = %t, want %t", iface, got, want)
			}

			return nil
		}

		return fmt.Errorf("interface %q not found on node", iface)
	}
}

// reloadNode applies whatever is staged, so a test that leaves a pending change on purpose
// does not hand the next test a dirty node.
func reloadNode(te *test.Environment) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		return te.NodeClient().ReloadNetworkConfiguration(context.Background())
	}
}

// checkNodePendingChanges asserts whether the node has staged network changes.
//
// A staged delete drops out of the interface list immediately, so `active` cannot tell it apart
// from an applied one. PVE only sends `changes` while something is pending.
func checkNodePendingChanges(te *test.Environment, want bool) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		got, err := nodeHasPendingChanges(te)
		if err != nil {
			return err
		}

		if got != want {
			return fmt.Errorf("node %q: pending changes = %t, want %t", te.NodeName, got, want)
		}

		return nil
	}
}

func nodeHasPendingChanges(te *test.Environment) (bool, error) {
	var body struct {
		Changes *string `json:"changes"`
	}

	path := fmt.Sprintf("nodes/%s/network", te.NodeName)
	if err := te.Client().DoRequest(context.Background(), http.MethodGet, path, nil, &body); err != nil {
		return false, fmt.Errorf("reading node network: %w", err)
	}

	return body.Changes != nil && *body.Changes != "", nil
}

// checkStagedDestroy asserts that a destroy with `reload = false` left its delete staged rather
// than applied, then clears the pending change off the shared node.
func checkStagedDestroy(te *test.Environment) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		checkNodePendingChanges(te, true),
		reloadNode(te),
	)
}

// newInterfaceListDropProxy fronts the real PVE endpoint and removes the named interface from every node
// interface list response, which is what a privilege-separated token without SDN.Audit sees for bridges.
func newInterfaceListDropProxy(t *testing.T, iface string) string {
	t.Helper()

	target, err := url.Parse(utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT"))
	require.NoError(t, err)

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
			r.Out.Header.Del("Accept-Encoding")
		},
		Transport: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		ModifyResponse: func(resp *http.Response) error {
			if resp.Request.Method != http.MethodGet || !strings.HasSuffix(resp.Request.URL.Path, "/network") {
				return nil
			}

			raw, err := io.ReadAll(resp.Body)
			_ = resp.Body.Close()

			if err != nil {
				return err
			}

			var payload struct {
				Data []map[string]any `json:"data"`
			}

			if err := json.Unmarshal(raw, &payload); err != nil {
				return err
			}

			kept := make([]map[string]any, 0, len(payload.Data))

			for _, entry := range payload.Data {
				if entry["iface"] != iface {
					kept = append(kept, entry)
				}
			}

			body, err := json.Marshal(map[string]any{"data": kept})
			if err != nil {
				return err
			}

			resp.Body = io.NopCloser(strings.NewReader(string(body)))
			resp.ContentLength = int64(len(body))
			resp.Header.Set("Content-Length", strconv.Itoa(len(body)))

			return nil
		},
	}

	server := httptest.NewTLSServer(proxy)
	t.Cleanup(server.Close)

	return server.URL
}

// cleanupStagedInterface removes the interface through the real endpoint if a failed create left it staged, then
// clears the shadow file. The revert is node-wide, so it only runs when the shadow carries no diff and there is
// nothing of anyone else's to discard.
func cleanupStagedInterface(t *testing.T, te *test.Environment, iface string) {
	t.Helper()

	t.Cleanup(func() {
		ctx := context.Background()

		ifaces, err := te.NodeClient().ListNetworkInterfaces(ctx)
		if err != nil {
			t.Logf("cleanup: list interfaces: %v", err)

			return
		}

		for _, i := range ifaces {
			if i.Iface != iface {
				continue
			}

			if err := te.NodeClient().DeleteNetworkInterface(ctx, iface); err != nil {
				t.Logf("cleanup: delete interface %s: %v", iface, err)
			}
		}

		pending, err := nodeHasPendingChanges(te)
		if err != nil {
			t.Logf("cleanup: %v", err)

			return
		}

		if pending {
			t.Logf("cleanup: node %s has unrelated pending network changes, leaving them", te.NodeName)

			return
		}

		if err := te.NodeClient().RevertNetworkConfiguration(ctx); err != nil {
			t.Logf("cleanup: revert network configuration: %v", err)
		}
	})
}

// requireInterfaceNotStaged lists interfaces through the real endpoint, which includes pending entries. It runs
// after resource.Test because the harness skips CheckDestroy when the final state is empty.
func requireInterfaceNotStaged(t *testing.T, te *test.Environment, iface string) {
	t.Helper()

	ifaces, err := te.NodeClient().ListNetworkInterfaces(context.Background())
	require.NoError(t, err)

	for _, i := range ifaces {
		require.NotEqual(t, iface, i.Iface, "interface %s is still staged on node %s after the failed create", iface, te.NodeName)
	}
}
