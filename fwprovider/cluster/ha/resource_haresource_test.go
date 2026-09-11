//go:build acceptance || all

//testacc:tier=medium
//testacc:resource=ha

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ha_test

import (
	"context"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

// skipIfNoHA skips the test if the HA manager is not responding, e.g. because the test cluster does not have HA
// configured.
func skipIfNoHA(t *testing.T, te *test.Environment) {
	t.Helper()

	_, err := te.ClusterClient().HA().Resources().List(context.Background(), nil)
	if err != nil {
		t.Skipf("Test requires HA-capable cluster (HA manager not responding: %v)", err)
	}
}

// TestAccResourceHAResourceAutoRebalance verifies the lifecycle of the `auto_rebalance` attribute on the
// `proxmox_haresource` resource: setting it on create, changing it on update, and unsetting it so the resource
// reverts to the cluster default.
func TestAccResourceHAResourceAutoRebalance(t *testing.T) {
	te := test.InitEnvironment(t)

	skipIfNoHA(t, te)

	vmID := 100000 + rand.Intn(99999)

	te.AddTemplateVars(map[string]any{
		"TestVMID": vmID,
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			// Step 1: create with auto_rebalance = true
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_ha_auto_rebalance" {
						node_name = "{{.NodeName}}"
						vm_id     = {{.TestVMID}}
						started   = false
						name      = "test-ha-auto-rebalance"
					}

					resource "proxmox_haresource" "test_auto_rebalance" {
						depends_on = [
							proxmox_virtual_environment_vm.test_ha_auto_rebalance
						]
						resource_id    = "vm:{{.TestVMID}}"
						state          = "stopped"
						auto_rebalance = true
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_haresource.test_auto_rebalance", "auto_rebalance", "true"),
				),
			},
			// Step 2: import
			{
				ResourceName:      "proxmox_haresource.test_auto_rebalance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Step 3: update auto_rebalance = false
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_ha_auto_rebalance" {
						node_name = "{{.NodeName}}"
						vm_id     = {{.TestVMID}}
						started   = false
						name      = "test-ha-auto-rebalance"
					}

					resource "proxmox_haresource" "test_auto_rebalance" {
						depends_on = [
							proxmox_virtual_environment_vm.test_ha_auto_rebalance
						]
						resource_id    = "vm:{{.TestVMID}}"
						state          = "stopped"
						auto_rebalance = false
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_haresource.test_auto_rebalance", "auto_rebalance", "false"),
				),
			},
			// Step 4: unset auto_rebalance, reverting to the cluster default
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_ha_auto_rebalance" {
						node_name = "{{.NodeName}}"
						vm_id     = {{.TestVMID}}
						started   = false
						name      = "test-ha-auto-rebalance"
					}

					resource "proxmox_haresource" "test_auto_rebalance" {
						depends_on = [
							proxmox_virtual_environment_vm.test_ha_auto_rebalance
						]
						resource_id = "vm:{{.TestVMID}}"
						state       = "stopped"
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("proxmox_haresource.test_auto_rebalance", "auto_rebalance"),
				),
			},
		},
	})
}
