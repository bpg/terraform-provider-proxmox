//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resources_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
)

// TestAccResourceVMHAMigration tests that HA-managed VMs can be migrated to another node.
// This test requires:
// - A multi-node Proxmox cluster (at least 2 nodes)
// - HA capability enabled
// The test will be skipped if these requirements are not met.
func TestAccResourceVMHAMigration(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	// Get available nodes in the cluster
	nodes := getClusterNodes(t, te.ClusterClient())
	if len(nodes) < 2 {
		t.Skip("Test requires at least 2 nodes in the cluster for HA migration testing")
	}

	sourceNode := nodes[0]
	targetNode := nodes[1]

	te.AddTemplateVars(map[string]any{
		"SourceNode": sourceNode,
		"TargetNode": targetNode,
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			// Step 1: Create VM with HA resource on source node
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_hagroup" "test_ha_group" {
						group = "test-ha-migration-group"
						nodes = {
							"{{.SourceNode}}" = null
							"{{.TargetNode}}" = null
						}
					}

					resource "proxmox_virtual_environment_vm" "test_ha_vm" {
						node_name = "{{.SourceNode}}"
						started   = false
						name      = "test-ha-migration-vm"
						
						# minimal VM config
						cpu {
							cores = 1
						}
						memory {
							dedicated = 512
						}
					}

					resource "proxmox_virtual_environment_haresource" "test_ha_resource" {
						depends_on = [
							proxmox_virtual_environment_hagroup.test_ha_group,
							proxmox_virtual_environment_vm.test_ha_vm
						]
						resource_id = "vm:${proxmox_virtual_environment_vm.test_ha_vm.vm_id}"
						state       = "stopped"
						group       = proxmox_virtual_environment_hagroup.test_ha_group.group
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_virtual_environment_vm.test_ha_vm", "node_name", sourceNode),
				),
			},
			// Step 2: Migrate VM to target node
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_hagroup" "test_ha_group" {
						group = "test-ha-migration-group"
						nodes = {
							"{{.SourceNode}}" = null
							"{{.TargetNode}}" = null
						}
					}

					resource "proxmox_virtual_environment_vm" "test_ha_vm" {
						node_name = "{{.TargetNode}}"
						started   = false
						migrate   = true
						name      = "test-ha-migration-vm"
						
						# minimal VM config
						cpu {
							cores = 1
						}
						memory {
							dedicated = 512
						}
					}

					resource "proxmox_virtual_environment_haresource" "test_ha_resource" {
						depends_on = [
							proxmox_virtual_environment_hagroup.test_ha_group,
							proxmox_virtual_environment_vm.test_ha_vm
						]
						resource_id = "vm:${proxmox_virtual_environment_vm.test_ha_vm.vm_id}"
						state       = "stopped"
						group       = proxmox_virtual_environment_hagroup.test_ha_group.group
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_virtual_environment_vm.test_ha_vm", "node_name", targetNode),
				),
			},
		},
	})
}

// getClusterNodes returns a list of node names in the cluster.
func getClusterNodes(t *testing.T, client *cluster.Client) []string {
	t.Helper()

	resources, err := client.GetClusterResources(context.Background(), "node")
	if err != nil {
		t.Fatalf("Failed to get cluster nodes: %v", err)
	}

	var nodes []string
	for _, r := range resources {
		// for node resources, the name is in the Name field
		if r.Type == "node" && r.Name != "" {
			nodes = append(nodes, r.Name)
		}
	}

	return nodes
}
