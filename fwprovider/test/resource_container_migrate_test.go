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

// TestAccResourceContainerNodeDriftDoesNotPlanBareCreate checks that read follows a container moved out of band
// instead of treating it as deleted.
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
			datastore_id  = "{{.ContainerDatastoreID}}"
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
				Config: config,
			},
		},
	})
}

// TestAccResourceContainerNodeDriftRecreatesOnDeclaredNode checks that a moved container is replaced on the
// declared node, not planned as a bare create that collides on the VMID.
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
			datastore_id  = "{{.ContainerDatastoreID}}"
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

// TestAccResourceContainerDestroyAfterNodeDrift checks that delete finds a container moved since the last persisted
// refresh: plan-only steps don't save state and the post-test destroy runs with -refresh=false.
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
			datastore_id  = "{{.ContainerDatastoreID}}"
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
				// Move it out of band in a plan-only step: the plan refreshes but never persists state, so state
				// still names the original node when the framework's implicit destroy (-refresh=false) runs.
				PreConfig: func() {
					migrateContainerOutOfBand(t, te, containerID, te.Node2Name)
				},
				Config:   config,
				PlanOnly: true,
				// The refresh finds the container on its new node and ignore_changes hides the node_name diff.
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
