//go:build acceptance || all

//testacc:tier=heavy
//testacc:resource=container

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/containers"
)

// migrateContainerOutOfBand moves a container with the API directly, simulating an HA rebalance.
func migrateContainerOutOfBand(t *testing.T, te *Environment, vmID int, targetNode string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	body := &containers.MigrateRequestBody{TargetNode: targetNode}

	err := te.NodeClient().Container(vmID).MigrateContainer(ctx, body).Err()
	require.NoError(t, err, "out-of-band migration to %s must succeed", targetNode)
}

// TestAccResourceContainerNodeDriftDoesNotPlanBareCreate proves that a container moved out of band is
// found on its new node during read, rather than being treated as deleted. Without the fix, the read
// queries the stale node, gets a 404, blanks the ID, and Terraform plans a create that can never
// succeed because the VMID is still taken cluster-wide.
func TestAccResourceContainerNodeDriftDoesNotPlanBareCreate(t *testing.T) {
	te := InitEnvironment(t)

	if te.Node2Name == "" {
		t.Skip("PROXMOX_VE_ACC_NODE_2_NAME must be set")
	}

	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())
	testAccDownloadContainerTemplate(t, te, imageFileName)

	containerID := 100000 + rand.Intn(99999)

	te.AddTemplateVars(map[string]any{
		"ImageFileName":   imageFileName,
		"TestContainerID": containerID,
	})

	config := te.RenderConfig(`
	resource "proxmox_virtual_environment_container" "test_drift" {
		node_name    = "{{.NodeName}}"
		vm_id        = {{.TestContainerID}}
		unprivileged = true
		started      = false

		disk {
			datastore_id  = "local-lvm"
			size          = 4
			mount_options = []
		}

		initialization {
			hostname = "test-drift"
		}

		operating_system {
			template_file_id = "local:vztmpl/{{.ImageFileName}}"
			type             = "alpine"
		}

		lifecycle {
			ignore_changes = [node_name]
		}
	}`)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				PreConfig: func() {
					migrateContainerOutOfBand(t, te, containerID, te.Node2Name)
				},
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				// A real apply, so the persisted state picks up the corrected node_name (ignore_changes only
				// suppresses the diff, it doesn't skip refresh) and the test's final destroy targets pve-b.
				Config: config,
			},
		},
	})
}

// TestAccResourceContainerNodeDriftRecreatesOnDeclaredNode proves the recovery path for a container
// that HA moved: with no ignore_changes, Terraform must plan a replacement and the destroy half must
// target the node the container actually sits on. Without the read fix the container reads as deleted,
// Terraform plans a bare create, and that create fails because the VMID is still taken cluster-wide.
func TestAccResourceContainerNodeDriftRecreatesOnDeclaredNode(t *testing.T) {
	te := InitEnvironment(t)

	if te.Node2Name == "" {
		t.Skip("PROXMOX_VE_ACC_NODE_2_NAME must be set")
	}

	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())
	testAccDownloadContainerTemplate(t, te, imageFileName)

	containerID := 100000 + rand.Intn(99999)

	te.AddTemplateVars(map[string]any{
		"ImageFileName":   imageFileName,
		"TestContainerID": containerID,
	})

	resourceName := "proxmox_virtual_environment_container.test_recreate_on_drift"

	config := te.RenderConfig(`
	resource "proxmox_virtual_environment_container" "test_recreate_on_drift" {
		node_name    = "{{.NodeName}}"
		vm_id        = {{.TestContainerID}}
		unprivileged = true
		started      = false

		disk {
			datastore_id  = "local-lvm"
			size          = 4
			mount_options = []
		}

		initialization {
			hostname = "test-recreate-on-drift"
		}

		operating_system {
			template_file_id = "local:vztmpl/{{.ImageFileName}}"
			type             = "alpine"
		}
	}`)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				PreConfig: func() {
					migrateContainerOutOfBand(t, te, containerID, te.Node2Name)
				},
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Not ResourceActionCreate: the container still exists on the other node, so a bare
						// create would collide on the VMID and the apply could never succeed.
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "node_name", te.NodeName),
					func(*terraform.State) error {
						ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
						defer cancel()

						_, err := te.NodeClient().Container(containerID).GetContainer(ctx)
						require.NoError(t, err, "container %d must exist on the declared node after replacement", containerID)

						return nil
					},
				),
			},
		},
	})
}

// TestAccResourceContainerDestroyAfterNodeDrift proves the destroy path tolerates a container that
// moved after the last refresh. Terraform destroys from state, so a container HA relocated since the
// last read would otherwise be deleted against the node it left, failing with
// "configuration file does not exist" and leaving the container behind with no way to remove it.
func TestAccResourceContainerDestroyAfterNodeDrift(t *testing.T) {
	te := InitEnvironment(t)

	if te.Node2Name == "" {
		t.Skip("PROXMOX_VE_ACC_NODE_2_NAME must be set")
	}

	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())
	testAccDownloadContainerTemplate(t, te, imageFileName)

	containerID := 100000 + rand.Intn(99999)

	te.AddTemplateVars(map[string]any{
		"ImageFileName":   imageFileName,
		"TestContainerID": containerID,
	})

	config := te.RenderConfig(`
	resource "proxmox_virtual_environment_container" "test_destroy_drift" {
		node_name    = "{{.NodeName}}"
		vm_id        = {{.TestContainerID}}
		unprivileged = true
		started      = false

		disk {
			datastore_id  = "local-lvm"
			size          = 4
			mount_options = []
		}

		initialization {
			hostname = "test-destroy-drift"
		}

		operating_system {
			template_file_id = "local:vztmpl/{{.ImageFileName}}"
			type             = "alpine"
		}

		lifecycle {
			ignore_changes = [node_name]
		}
	}`)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				// Move it out of band and do not refresh, so state still names the original node when the
				// framework's implicit destroy runs at the end of the test.
				PreConfig: func() {
					migrateContainerOutOfBand(t, te, containerID, te.Node2Name)
				},
				Config:   config,
				PlanOnly: true,
				// The plan is computed from the pre-move state, so nothing should be proposed.
				ExpectNonEmptyPlan: false,
				RefreshState:       false,
			},
		},
	})
}
