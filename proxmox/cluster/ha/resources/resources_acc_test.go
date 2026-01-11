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

// skipIfNotMultiNode skips the test if the cluster has fewer than 2 nodes.
func skipIfNotMultiNode(t *testing.T, nodes []string) {
	t.Helper()

	if len(nodes) < 2 {
		t.Skip("Test requires at least 2 nodes in the cluster for migration testing")
	}
}

// skipIfNoHA skips the test if HA is not available (checks if HA manager is responding).
func skipIfNoHA(t *testing.T, client *cluster.Client) {
	t.Helper()

	_, err := client.HA().Resources().List(context.Background(), nil)
	if err != nil {
		t.Skipf("Test requires HA-capable cluster (HA manager not responding: %v)", err)
	}
}

// TestAccResourceVMHAMigrationStopped tests that stopped HA-managed VMs can be migrated to another node.
// This uses the workaround: remove from HA -> standard migrate -> re-add to HA.
//
// Requirements:
//   - Multi-node Proxmox cluster (at least 2 nodes)
//   - HA enabled and operational (corosync quorum)
//   - local-lvm storage available on all nodes (or shared storage)
//
// Note: Migration tests do NOT run in parallel to avoid VM ID and HA resource conflicts.
func TestAccResourceVMHAMigrationStopped(t *testing.T) {
	te := test.InitEnvironment(t)

	nodes := getClusterNodes(t, te.ClusterClient())
	skipIfNotMultiNode(t, nodes)
	skipIfNoHA(t, te.ClusterClient())

	sourceNode := nodes[0]
	targetNode := nodes[1]

	te.AddTemplateVars(map[string]any{
		"SourceNode": sourceNode,
		"TargetNode": targetNode,
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			// Step 1: Create stopped VM with HA resource on source node
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_ha_stopped" {
						node_name = "{{.SourceNode}}"
						started   = false
						name      = "test-ha-migration-stopped"
						
						cpu {
							cores = 1
						}
						memory {
							dedicated = 512
						}
					}

					resource "proxmox_virtual_environment_haresource" "test_ha_stopped" {
						depends_on = [
							proxmox_virtual_environment_vm.test_ha_stopped
						]
						resource_id = "vm:${proxmox_virtual_environment_vm.test_ha_stopped.vm_id}"
						state       = "stopped"
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_virtual_environment_vm.test_ha_stopped", "node_name", sourceNode),
				),
			},
			// Step 2: Migrate stopped HA VM to target node
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_ha_stopped" {
						node_name = "{{.TargetNode}}"
						started   = false
						migrate   = true
						name      = "test-ha-migration-stopped"
						
						cpu {
							cores = 1
						}
						memory {
							dedicated = 512
						}
					}

					resource "proxmox_virtual_environment_haresource" "test_ha_stopped" {
						depends_on = [
							proxmox_virtual_environment_vm.test_ha_stopped
						]
						resource_id = "vm:${proxmox_virtual_environment_vm.test_ha_stopped.vm_id}"
						state       = "stopped"
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_virtual_environment_vm.test_ha_stopped", "node_name", targetNode),
				),
			},
		},
	})
}

// TestAccResourceVMHAMigrationRunning tests that running HA-managed VMs can be migrated to another node.
// This uses the HA manager's native migration which performs live migration.
//
// Requirements:
//   - Multi-node Proxmox cluster (at least 2 nodes)
//   - HA enabled and operational (corosync quorum)
//   - local-lvm storage available on all nodes (or shared storage)
//
// Note: Migration tests do NOT run in parallel to avoid VM ID and HA resource conflicts.
func TestAccResourceVMHAMigrationRunning(t *testing.T) {
	te := test.InitEnvironment(t)

	nodes := getClusterNodes(t, te.ClusterClient())
	skipIfNotMultiNode(t, nodes)
	skipIfNoHA(t, te.ClusterClient())

	sourceNode := nodes[0]
	targetNode := nodes[1]

	te.AddTemplateVars(map[string]any{
		"SourceNode": sourceNode,
		"TargetNode": targetNode,
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			// Step 1: Create running VM with HA resource on source node
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_ha_running" {
						node_name       = "{{.SourceNode}}"
						started         = true
						stop_on_destroy = true
						name            = "test-ha-migration-running"
						
						cpu {
							cores = 1
						}
						memory {
							dedicated = 512
						}
						disk {
							datastore_id = "local-lvm"
							file_format  = "raw"
							interface    = "scsi0"
							size         = 1
						}
					}

					resource "proxmox_virtual_environment_haresource" "test_ha_running" {
						depends_on = [
							proxmox_virtual_environment_vm.test_ha_running
						]
						resource_id = "vm:${proxmox_virtual_environment_vm.test_ha_running.vm_id}"
						state       = "started"
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_virtual_environment_vm.test_ha_running", "node_name", sourceNode),
					resource.TestCheckResourceAttr("proxmox_virtual_environment_vm.test_ha_running", "started", "true"),
				),
			},
			// Step 2: Migrate running HA VM to target node (uses HA manager live migration)
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_ha_running" {
						node_name       = "{{.TargetNode}}"
						started         = true
						migrate         = true
						stop_on_destroy = true
						name            = "test-ha-migration-running"
						
						cpu {
							cores = 1
						}
						memory {
							dedicated = 512
						}
						disk {
							datastore_id = "local-lvm"
							file_format  = "raw"
							interface    = "scsi0"
							size         = 1
						}
					}

					resource "proxmox_virtual_environment_haresource" "test_ha_running" {
						depends_on = [
							proxmox_virtual_environment_vm.test_ha_running
						]
						resource_id = "vm:${proxmox_virtual_environment_vm.test_ha_running.vm_id}"
						state       = "started"
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_virtual_environment_vm.test_ha_running", "node_name", targetNode),
					resource.TestCheckResourceAttr("proxmox_virtual_environment_vm.test_ha_running", "started", "true"),
				),
			},
		},
	})
}

// TestAccResourceVMMigrationStopped tests that stopped non-HA VMs can be migrated to another node.
// Requires a multi-node Proxmox cluster (at least 2 nodes).
// TestAccResourceVMMigrationStopped tests that stopped non-HA VMs can be migrated to another node.
//
// Requirements:
//   - Multi-node Proxmox cluster (at least 2 nodes)
//   - local-lvm storage available on all nodes (or shared storage)
//
// Note: Migration tests do NOT run in parallel to avoid VM ID conflicts.
func TestAccResourceVMMigrationStopped(t *testing.T) {
	te := test.InitEnvironment(t)

	nodes := getClusterNodes(t, te.ClusterClient())
	skipIfNotMultiNode(t, nodes)

	sourceNode := nodes[0]
	targetNode := nodes[1]

	te.AddTemplateVars(map[string]any{
		"SourceNode": sourceNode,
		"TargetNode": targetNode,
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			// Step 1: Create stopped VM on source node
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_stopped" {
						node_name = "{{.SourceNode}}"
						started   = false
						name      = "test-migration-stopped"
						
						cpu {
							cores = 1
						}
						memory {
							dedicated = 512
						}
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_virtual_environment_vm.test_stopped", "node_name", sourceNode),
				),
			},
			// Step 2: Migrate stopped VM to target node
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_stopped" {
						node_name = "{{.TargetNode}}"
						started   = false
						migrate   = true
						name      = "test-migration-stopped"
						
						cpu {
							cores = 1
						}
						memory {
							dedicated = 512
						}
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_virtual_environment_vm.test_stopped", "node_name", targetNode),
				),
			},
		},
	})
}

// TestAccResourceVMMigrationRunning tests that running non-HA VMs can be migrated to another node.
// This uses standard online migration with local disk support.
// Requires a multi-node Proxmox cluster (at least 2 nodes).
// TestAccResourceVMMigrationRunning tests that running non-HA VMs can be migrated to another node.
// This uses standard online migration with local disk support.
//
// Requirements:
//   - Multi-node Proxmox cluster (at least 2 nodes)
//   - local-lvm storage available on all nodes (or shared storage)
//
// Note: Migration tests do NOT run in parallel to avoid VM ID conflicts.
func TestAccResourceVMMigrationRunning(t *testing.T) {
	te := test.InitEnvironment(t)

	nodes := getClusterNodes(t, te.ClusterClient())
	skipIfNotMultiNode(t, nodes)

	sourceNode := nodes[0]
	targetNode := nodes[1]

	te.AddTemplateVars(map[string]any{
		"SourceNode": sourceNode,
		"TargetNode": targetNode,
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			// Step 1: Create running VM on source node
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_running" {
						node_name       = "{{.SourceNode}}"
						started         = true
						stop_on_destroy = true
						name            = "test-migration-running"
						
						cpu {
							cores = 1
						}
						memory {
							dedicated = 512
						}
						disk {
							datastore_id = "local-lvm"
							file_format  = "raw"
							interface    = "scsi0"
							size         = 1
						}
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_virtual_environment_vm.test_running", "node_name", sourceNode),
					resource.TestCheckResourceAttr("proxmox_virtual_environment_vm.test_running", "started", "true"),
				),
			},
			// Step 2: Migrate running VM to target node (online migration)
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_running" {
						node_name       = "{{.TargetNode}}"
						started         = true
						migrate         = true
						stop_on_destroy = true
						name            = "test-migration-running"
						
						cpu {
							cores = 1
						}
						memory {
							dedicated = 512
						}
						disk {
							datastore_id = "local-lvm"
							file_format  = "raw"
							interface    = "scsi0"
							size         = 1
						}
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_virtual_environment_vm.test_running", "node_name", targetNode),
					resource.TestCheckResourceAttr("proxmox_virtual_environment_vm.test_running", "started", "true"),
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
		// for node resources, the name is in the NodeName field (json:"node")
		if r.Type == "node" && r.NodeName != "" {
			nodes = append(nodes, r.NodeName)
		}
	}

	return nodes
}
