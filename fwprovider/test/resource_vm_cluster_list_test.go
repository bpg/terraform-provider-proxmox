//go:build acceptance || all

//testacc:tier=medium
//testacc:resource=vm

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"math/rand"
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

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

// newEmptyClusterListProxy fronts the real PVE endpoint and answers every /cluster/resources request with
// an empty list, which is what a multi-node cluster returns until the guest list has synced to the node
// serving the request.
func newEmptyClusterListProxy(t *testing.T) string {
	t.Helper()

	target, err := url.Parse(utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT"))
	require.NoError(t, err)

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
			// keep the upstream response uncompressed so the body below can be swapped in
			r.Out.Header.Del("Accept-Encoding")
		},
		Transport: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		ModifyResponse: func(resp *http.Response) error {
			if !strings.HasSuffix(resp.Request.URL.Path, "/cluster/resources") {
				return nil
			}

			_ = resp.Body.Close()

			body := `{"data":[]}`
			resp.Body = io.NopCloser(strings.NewReader(body))
			resp.ContentLength = int64(len(body))
			resp.Header.Set("Content-Length", strconv.Itoa(len(body)))

			return nil
		},
	}

	server := httptest.NewTLSServer(proxy)
	t.Cleanup(server.Close)

	return server.URL
}

// TestAccResourceVMReadSurvivesClusterListLag verifies that a VM missing from the cluster-wide resource
// list, as happens on a multi-node cluster right after creation, is still read through its node's config
// endpoint instead of being dropped from state.
func TestAccResourceVMReadSurvivesClusterListLag(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)
	vmID := 100000 + rand.Intn(99999)
	te.AddTemplateVars(map[string]any{"TestVMID": vmID})

	t.Cleanup(func() {
		err := te.NodeClient().VM(vmID).DeleteVM(context.Background(), true, true).Err()
		if err != nil && !errors.Is(err, api.ErrResourceDoesNotExist) {
			t.Logf("cleanup: delete VM %d: %v", vmID, err)
		}
	})

	endpoint := newEmptyClusterListProxy(t)
	resourceName := "proxmox_virtual_environment_vm.test_cluster_list_lag"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_cluster_list_lag" {
						node_name = "{{.NodeName}}"
						vm_id     = {{.TestVMID}}
						name      = "test-cluster-list-lag"
						started   = false
					}`, WithEndpoint(endpoint)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", strconv.Itoa(vmID)),
					func(*terraform.State) error {
						_, err := te.NodeClient().VM(vmID).GetVM(context.Background())

						return err
					},
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("%s/%d", te.NodeName, vmID),
			},
		},
	})
}
