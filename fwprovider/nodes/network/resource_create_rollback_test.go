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
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

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

// cleanupStagedInterface removes the interface through the real endpoint if a failed create left it staged,
// then discards the pending network file so the node is left as it was.
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

func TestAccResourceLinuxBridgeCreateRollsBackOnReadMiss(t *testing.T) {
	te := test.InitEnvironment(t)

	iface := fmt.Sprintf("vmbr%d", gofakeit.Number(10, 9999))
	cleanupStagedInterface(t, te, iface)

	endpoint := newInterfaceListDropProxy(t, iface)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_bridge" "test_rollback" {
					name           = "%s"
					node_name      = "{{.NodeName}}"
					timeout_reload = 60
				}
				`, iface), test.WithInsecureEndpoint(endpoint)),
				ExpectError: regexp.MustCompile(`not found after creation`),
			},
		},
	})

	requireInterfaceNotStaged(t, te, iface)
}

func TestAccResourceLinuxVLANCreateRollsBackOnReadMiss(t *testing.T) {
	te := test.InitEnvironment(t)

	parent := os.Getenv("PROXMOX_VE_ACC_IFACE_NAME")
	if parent == "" {
		parent = "ens18"
	}

	iface := fmt.Sprintf("%s.%d", parent, gofakeit.Number(10, 4094))
	cleanupStagedInterface(t, te, iface)

	endpoint := newInterfaceListDropProxy(t, iface)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_vlan" "test_rollback" {
					name           = "%s"
					node_name      = "{{.NodeName}}"
					timeout_reload = 60
				}
				`, iface), test.WithInsecureEndpoint(endpoint)),
				ExpectError: regexp.MustCompile(`not found after creation`),
			},
		},
	})

	requireInterfaceNotStaged(t, te, iface)
}

func TestAccResourceLinuxBondCreateRollsBackOnReadMiss(t *testing.T) {
	te := test.InitEnvironment(t)

	slave1 := os.Getenv("PROXMOX_VE_ACC_BOND_SLAVE1")
	slave2 := os.Getenv("PROXMOX_VE_ACC_BOND_SLAVE2")

	if slave1 == "" || slave2 == "" {
		t.Skip("skipping: PROXMOX_VE_ACC_BOND_SLAVE1 and PROXMOX_VE_ACC_BOND_SLAVE2 must be set to eth-type interfaces")
	}

	iface := fmt.Sprintf("bond%d", gofakeit.Number(10, 9999))
	cleanupStagedInterface(t, te, iface)

	endpoint := newInterfaceListDropProxy(t, iface)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_network_linux_bond" "test_rollback" {
					name           = "%s"
					node_name      = "{{.NodeName}}"
					slaves         = ["%s", "%s"]
					timeout_reload = 60
				}
				`, iface, slave1, slave2), test.WithInsecureEndpoint(endpoint)),
				ExpectError: regexp.MustCompile(`Unable to Read Linux Bond After Creation`),
			},
		},
	})

	requireInterfaceNotStaged(t, te, iface)
}
