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
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"

	haresources "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/ha/resources"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/containers"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
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
				// Plan-only steps don't persist state, so the implicit destroy still sees the original node.
				PreConfig: func() {
					migrateContainerOutOfBand(t, te, containerID, te.Node2Name)
				},
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccDownloadContainerTemplateOnNode(t *testing.T, te *Environment, imageFileName, nodeName string) {
	t.Helper()

	nodeClient := &nodes.Client{Client: te.Client(), NodeName: nodeName}
	storageClient := &storage.Client{Client: nodeClient, StorageName: te.DatastoreID}

	err := storageClient.DownloadFileByURL(context.Background(), &storage.DownloadURLPostRequestBody{
		Content:  new("vztmpl"),
		FileName: new(imageFileName),
		Node:     new(nodeName),
		Storage:  new(te.DatastoreID),
		URL:      new(getTemplateURL(t, te.ContainerImagesServer)),
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		e := storageClient.DeleteDatastoreFile(context.Background(), fmt.Sprintf("vztmpl/%s", imageFileName)).Err()
		require.NoError(t, e)
	})
}

// TestAccResourceContainerMigrateDisabledRecreates checks that without `migrate` a node change still replaces.
func TestAccResourceContainerMigrateDisabledRecreates(t *testing.T) {
	te := InitEnvironment(t)

	if te.Node2Name == "" {
		t.Skip("PROXMOX_VE_ACC_NODE_2_NAME must be set")
	}

	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())
	testAccDownloadContainerTemplate(t, te, imageFileName)
	testAccDownloadContainerTemplateOnNode(t, te, imageFileName, te.Node2Name)

	containerID := 100000 + rand.Intn(99999)

	te.AddTemplateVars(map[string]any{
		"ImageFileName":   imageFileName,
		"TestContainerID": containerID,
	})

	resourceName := "proxmox_virtual_environment_container.test_recreate"

	containerConfig := func(node string) string {
		return te.RenderConfig(fmt.Sprintf(`
		resource "proxmox_virtual_environment_container" "test_recreate" {
			node_name    = "%s"
			vm_id        = {{.TestContainerID}}
			unprivileged = true
			started      = false

			disk {
				datastore_id  = "{{.ContainerDatastoreID}}"
				size          = 4
				mount_options = []
			}

			initialization {
				hostname = "test-recreate"
			}

			operating_system {
				template_file_id = "local:vztmpl/{{.ImageFileName}}"
				type             = "alpine"
			}
		}`, node))
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: containerConfig(te.NodeName),
				Check:  resource.TestCheckResourceAttr(resourceName, "node_name", te.NodeName),
			},
			{
				Config: containerConfig(te.Node2Name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.TestCheckResourceAttr(resourceName, "node_name", te.Node2Name),
			},
		},
	})
}

// containerClientOnNode returns a container client for any node; te.NodeClient() covers only the primary one.
func containerClientOnNode(te *Environment, node string, vmID int) *containers.Client {
	nodeClient := &nodes.Client{Client: te.Client(), NodeName: node}
	return nodeClient.Container(vmID)
}

// requireContainerOnNode asserts the container is readable on the given node.
func requireContainerOnNode(t *testing.T, te *Environment, vmID int, node string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	_, err := containerClientOnNode(te, node, vmID).GetContainer(ctx)
	require.NoError(t, err, "container %d must be readable on node %s", vmID, node)
}

// containerMigrateConfig renders a migratable container on the given node, started as requested.
func containerMigrateConfig(te *Environment, node, hostname string, started bool) string {
	return te.RenderConfig(fmt.Sprintf(`
	resource "proxmox_virtual_environment_container" "test_migrate" {
		node_name    = "%s"
		vm_id        = {{.TestContainerID}}
		unprivileged = true
		migrate      = true
		started      = %t

		disk {
			datastore_id  = "{{.ContainerDatastoreID}}"
			size          = 4
			mount_options = []
		}

		initialization {
			hostname = "%s"
		}

		operating_system {
			template_file_id = "local:vztmpl/{{.ImageFileName}}"
			type             = "alpine"
		}

		network_interface {
			name = "vmbr0"
		}
	}`, node, started, hostname))
}

// TestAccResourceContainerMigrateStopped checks that a stopped container is migrated, not recreated.
func TestAccResourceContainerMigrateStopped(t *testing.T) {
	te := InitEnvironment(t)

	if te.Node2Name == "" {
		t.Skip("PROXMOX_VE_ACC_NODE_2_NAME must be set")
	}

	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())
	testAccDownloadContainerTemplate(t, te, imageFileName)
	testAccDownloadContainerTemplateOnNode(t, te, imageFileName, te.Node2Name)

	containerID := 100000 + rand.Intn(99999)

	te.AddTemplateVars(map[string]any{
		"ImageFileName":   imageFileName,
		"TestContainerID": containerID,
	})

	resourceName := "proxmox_virtual_environment_container.test_migrate"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: containerMigrateConfig(te, te.NodeName, "test-migrate", false),
				Check:  resource.TestCheckResourceAttr(resourceName, "node_name", te.NodeName),
			},
			{
				Config: containerMigrateConfig(te, te.Node2Name, "test-migrate", false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "node_name", te.Node2Name),
					func(*terraform.State) error {
						requireContainerOnNode(t, te, containerID, te.Node2Name)
						return nil
					},
				),
			},
		},
	})
}

// TestAccResourceContainerMigrateRunning checks that a running container is migrated with restart and keeps running.
func TestAccResourceContainerMigrateRunning(t *testing.T) {
	te := InitEnvironment(t)

	if te.Node2Name == "" {
		t.Skip("PROXMOX_VE_ACC_NODE_2_NAME must be set")
	}

	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())
	testAccDownloadContainerTemplate(t, te, imageFileName)
	testAccDownloadContainerTemplateOnNode(t, te, imageFileName, te.Node2Name)

	containerID := 100000 + rand.Intn(99999)

	te.AddTemplateVars(map[string]any{
		"ImageFileName":   imageFileName,
		"TestContainerID": containerID,
	})

	resourceName := "proxmox_virtual_environment_container.test_migrate"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: containerMigrateConfig(te, te.NodeName, "test-migrate", true),
			},
			{
				Config: containerMigrateConfig(te, te.Node2Name, "test-migrate", true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "node_name", te.Node2Name),
					func(*terraform.State) error {
						requireContainerOnNode(t, te, containerID, te.Node2Name)

						ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
						defer cancel()

						status, err := containerClientOnNode(te, te.Node2Name, containerID).GetContainerStatus(ctx)
						require.NoError(t, err)
						require.Equal(t, "running", status.Status, "container must be running after restart migration")

						return nil
					},
				),
			},
		},
	})
}

// containerHAResourceID returns the HA resource ID ("ct:<vmid>") of a container.
func containerHAResourceID(vmID int) proxmoxtypes.HAResourceID {
	return proxmoxtypes.HAResourceID{Type: proxmoxtypes.HAResourceTypeContainer, Name: strconv.Itoa(vmID)}
}

// requireContainerHAState asserts the requested state in the HA configuration, not /cluster/resources' hastate.
func requireContainerHAState(t *testing.T, te *Environment, vmID int, want proxmoxtypes.HAResourceState) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	haResource, err := te.ClusterClient().HA().Resources().Get(ctx, containerHAResourceID(vmID))
	require.NoError(t, err, "reading HA configuration of container %d", vmID)
	require.Equal(t, want, haResource.State, "unexpected HA state for container %d", vmID)
}

// registerContainerHA puts the container under HA and returns a deregister function, also run on cleanup.
func registerContainerHA(t *testing.T, te *Environment, vmID int) func() {
	t.Helper()

	haClient := te.ClusterClient().HA().Resources()
	haResourceID := containerHAResourceID(vmID)

	deregister := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		if _, err := haClient.Get(ctx, haResourceID); err != nil {
			return
		}

		require.NoError(t, haClient.Delete(ctx, haResourceID), "removing container %d from HA", vmID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	err := haClient.Create(ctx, &haresources.HAResourceCreateRequestBody{
		HAResourceDataBase: haresources.HAResourceDataBase{
			State:   proxmoxtypes.HAResourceStateStarted,
			Comment: new("terraform provider acceptance test"),
		},
		ID: haResourceID,
	})
	require.NoError(t, err, "registering container %d in HA", vmID)

	t.Cleanup(deregister)

	return deregister
}

// waitForContainerHASettled waits for hastate "started" on node across consecutive reads: a newly registered
// resource cycles through request_start_balance for ~20s while the CRM adopts it.
func waitForContainerHASettled(t *testing.T, te *Environment, vmID int, node string, timeout time.Duration) {
	t.Helper()

	const (
		pollInterval = 2 * time.Second
		stableReads  = 10
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	stable := 0
	last := "<none>"

	for {
		select {
		case <-ctx.Done():
			t.Fatalf("container %d did not hold HA state \"started\" on node %s within %s (last seen %s)",
				vmID, node, timeout, last)
		case <-ticker.C:
			entry, err := te.ClusterClient().GetContainerResource(ctx, vmID)
			if err != nil {
				stable = 0
				continue
			}

			last = fmt.Sprintf("node=%s hastate=%q", entry.NodeName, entry.HaState)

			if entry.NodeName != node || entry.HaState != "started" {
				stable = 0
				continue
			}

			if stable++; stable >= stableReads {
				return
			}
		}
	}
}

// TestAccResourceContainerMigrateHA checks that an HA-managed container is migrated outside HA and handed back.
// Under HA, hamigrate reports TASK OK before the container moves (arrived ~18s, running ~33s), and with
// ha-rebalance-on-start CRS relocated it off the declared node ~15s after an offline-path start.
func TestAccResourceContainerMigrateHA(t *testing.T) {
	te := InitEnvironment(t)

	if te.Node2Name == "" {
		t.Skip("PROXMOX_VE_ACC_NODE_2_NAME must be set")
	}

	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())
	testAccDownloadContainerTemplate(t, te, imageFileName)
	testAccDownloadContainerTemplateOnNode(t, te, imageFileName, te.Node2Name)

	containerID := 100000 + rand.Intn(99999)

	te.AddTemplateVars(map[string]any{
		"ImageFileName":   imageFileName,
		"TestContainerID": containerID,
	})

	resourceName := "proxmox_virtual_environment_container.test_migrate"

	var (
		deregisterHA   func()
		migrationStart int64
	)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: containerMigrateConfig(te, te.NodeName, "test-migrate", true),
			},
			{
				PreConfig: func() {
					deregisterHA = registerContainerHA(t, te, containerID)
					waitForContainerHASettled(t, te, containerID, te.NodeName, 3*time.Minute)
				},
				Config: containerMigrateConfig(te, te.Node2Name, "test-migrate", true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "node_name", te.Node2Name),
					func(*terraform.State) error {
						requireContainerOnNode(t, te, containerID, te.Node2Name)

						ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
						defer cancel()

						// /cluster/resources status lags a migration by tens of seconds.
						status, err := containerClientOnNode(te, te.Node2Name, containerID).GetContainerStatus(ctx)
						require.NoError(t, err)
						require.Equal(t, "running", status.Status,
							"container must already be running on the target node when the apply returns, "+
								"not merely queued for relocation by the HA manager")

						requireContainerHAState(t, te, containerID, proxmoxtypes.HAResourceStateStarted)

						return nil
					},
				),
			},
			{
				// node_name and hostname together take the offline path.
				PreConfig: func() {
					waitForContainerHASettled(t, te, containerID, te.Node2Name, 3*time.Minute)
					migrationStart = time.Now().Add(-2 * time.Second).Unix()
				},
				Config: containerMigrateConfig(te, te.NodeName, "test-migrate-ha", true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "node_name", te.NodeName),
					resource.TestCheckResourceAttr(resourceName, "initialization.0.hostname", "test-migrate-ha"),
					func(*terraform.State) error {
						ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
						defer cancel()

						declaredAPI := containerClientOnNode(te, te.NodeName, containerID)

						status, err := declaredAPI.GetContainerStatus(ctx)
						require.NoError(t, err)
						require.Equal(t, "running", status.Status,
							"container must be running on the node Terraform declared when the apply returns")

						containerConfig, err := declaredAPI.GetContainer(ctx)
						require.NoError(t, err)
						require.NotNil(t, containerConfig.Hostname)
						require.Equal(t, "test-migrate-ha", *containerConfig.Hostname,
							"the new configuration must be in effect on the node the container ended up on")

						counts := containerTaskTypeCounts(t, te, containerID, migrationStart)
						t.Logf("HA compound migration task types by node: %v", counts)

						byType := map[string]int{}
						for node, perNode := range counts {
							for taskType, n := range perNode {
								byType[taskType] += n

								require.NotContains(t, []string{"vzreboot", "vzstop"}, taskType,
									"the container must not be rebooted or hard-stopped (%s on %s): %v", taskType, node, counts)
							}
						}

						for _, haTask := range []string{"hastop", "hastart", "hamigrate"} {
							require.Zero(t, byType[haTask],
								"HA must not intercept the migration (%s seen): %v", haTask, counts)
						}

						require.Equal(t, 1, byType["vzmigrate"], "exactly one plain migrate is expected: %v", counts)
						require.Equal(t, 1, counts[te.Node2Name]["vzshutdown"],
							"the shutdown belongs on the source node, before the migration: %v", counts)
						require.Equal(t, 1, byType["vzstart"],
							"the container must be booted exactly once, after the config was written: %v", counts)
						require.Equal(t, 1, counts[te.NodeName]["vzstart"],
							"the boot must happen on the declared node, not wherever CRS moved the container: %v", counts)
						require.Zero(t, counts[te.NodeName]["vzshutdown"],
							"the shutdown belongs on the source node, before the migration: %v", counts)

						requireContainerHAState(t, te, containerID, proxmoxtypes.HAResourceStateStarted)

						// The container must stay where Terraform put it once HA has it back.
						waitForContainerHASettled(t, te, containerID, te.NodeName, 3*time.Minute)

						return nil
					},
				),
			},
			{
				// HA must be handed back as "stopped", or it boots the container again.
				Config: containerMigrateConfig(te, te.Node2Name, "test-migrate-ha", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "node_name", te.Node2Name),
					func(*terraform.State) error {
						requireContainerHAState(t, te, containerID, proxmoxtypes.HAResourceStateStopped)

						ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
						defer cancel()

						status, err := containerClientOnNode(te, te.Node2Name, containerID).GetContainerStatus(ctx)
						require.NoError(t, err)
						require.Equal(t, "stopped", status.Status, "container must stay stopped on the target node")

						return nil
					},
				),
			},
			{
				// Deregister before destroy, or HA restarts the container stopped for deletion.
				PreConfig: func() { deregisterHA() },
				Config:    containerMigrateConfig(te, te.Node2Name, "test-migrate-ha", false),
			},
		},
	})
}

// containerTaskEntry is one row of GET /nodes/{node}/tasks.
type containerTaskEntry struct {
	UPID      string `json:"upid"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	StartTime int64  `json:"starttime"`
}

// containerTasksSince returns the container's tasks on one node started at or after since.
func containerTasksSince(t *testing.T, te *Environment, node string, vmID int, since int64) []containerTaskEntry {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	resp := &struct {
		Data []containerTaskEntry `json:"data"`
	}{}

	path := fmt.Sprintf("nodes/%s/tasks?vmid=%d&limit=500&source=all", node, vmID)
	require.NoError(t, te.Client().DoRequest(ctx, http.MethodGet, path, nil, resp))

	var recent []containerTaskEntry

	for _, task := range resp.Data {
		if task.StartTime >= since {
			recent = append(recent, task)
		}
	}

	return recent
}

// containerTaskTypeCounts counts the container's exact task types per node, across all nodes since HA may run
// work on a third one.
func containerTaskTypeCounts(t *testing.T, te *Environment, vmID int, since int64) map[string]map[string]int {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	nodeList, err := (&nodes.Client{Client: te.Client()}).ListNodes(ctx)
	require.NoError(t, err)

	counts := map[string]map[string]int{}

	for _, node := range nodeList {
		perNode := map[string]int{}

		for _, task := range containerTasksSince(t, te, node.Name, vmID, since) {
			require.NotContains(t, task.Status, "migration problems",
				"the apply must not leave a failed migration in the task history (%s)", task.UPID)

			perNode[task.Type]++
		}

		if len(perNode) > 0 {
			counts[node.Name] = perNode
		}
	}

	return counts
}

// TestAccResourceContainerMigrateStopOnMove checks that node_name plus `started = false` shuts down once, on the
// source.
func TestAccResourceContainerMigrateStopOnMove(t *testing.T) {
	te := InitEnvironment(t)

	if te.Node2Name == "" {
		t.Skip("PROXMOX_VE_ACC_NODE_2_NAME must be set")
	}

	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())
	testAccDownloadContainerTemplate(t, te, imageFileName)
	testAccDownloadContainerTemplateOnNode(t, te, imageFileName, te.Node2Name)

	containerID := 100000 + rand.Intn(99999)

	te.AddTemplateVars(map[string]any{
		"ImageFileName":   imageFileName,
		"TestContainerID": containerID,
	})

	resourceName := "proxmox_virtual_environment_container.test_migrate"

	var migrationStart int64

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: containerMigrateConfig(te, te.NodeName, "test-migrate", true),
			},
			{
				PreConfig: func() { migrationStart = time.Now().Add(-2 * time.Second).Unix() },
				Config:    containerMigrateConfig(te, te.Node2Name, "test-migrate", false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "node_name", te.Node2Name),
					func(*terraform.State) error {
						requireContainerOnNode(t, te, containerID, te.Node2Name)

						ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
						defer cancel()

						status, err := containerClientOnNode(te, te.Node2Name, containerID).GetContainerStatus(ctx)
						require.NoError(t, err)
						require.Equal(t, "stopped", status.Status, "container must be stopped on the target node")

						counts := containerTaskTypeCounts(t, te, containerID, migrationStart)
						t.Logf("stop-on-move task types by node: %v", counts)

						require.Equal(t, 1, counts[te.NodeName]["vzshutdown"],
							"the provider must shut the container down once, on the source node: %v", counts)
						require.Zero(t, counts[te.Node2Name]["vzstart"],
							"a restart migration booted the container on the target it was asked to stop on: %v", counts)
						require.Zero(t, counts[te.Node2Name]["vzshutdown"]+counts[te.Node2Name]["vzstop"],
							"a second shutdown on the target means the container was migrated running: %v", counts)

						return nil
					},
				),
			},
		},
	})
}

// TestAccResourceContainerMigrateWithConfigChange checks that node_name plus a config change boots the container
// once, after the config is written.
func TestAccResourceContainerMigrateWithConfigChange(t *testing.T) {
	te := InitEnvironment(t)

	if te.Node2Name == "" {
		t.Skip("PROXMOX_VE_ACC_NODE_2_NAME must be set")
	}

	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())
	testAccDownloadContainerTemplate(t, te, imageFileName)
	testAccDownloadContainerTemplateOnNode(t, te, imageFileName, te.Node2Name)

	containerID := 100000 + rand.Intn(99999)

	te.AddTemplateVars(map[string]any{
		"ImageFileName":   imageFileName,
		"TestContainerID": containerID,
	})

	resourceName := "proxmox_virtual_environment_container.test_compound"

	compoundConfig := func(node, hostname string, diskSize int) string {
		return te.RenderConfig(fmt.Sprintf(`
		resource "proxmox_virtual_environment_container" "test_compound" {
			node_name    = "%s"
			vm_id        = {{.TestContainerID}}
			unprivileged = true
			migrate      = true
			started      = true

			disk {
				datastore_id  = "{{.ContainerDatastoreID}}"
				size          = %d
				mount_options = []
			}

			initialization {
				hostname = "%s"
			}

			operating_system {
				template_file_id = "local:vztmpl/{{.ImageFileName}}"
				type             = "alpine"
			}

			network_interface {
				name = "vmbr0"
			}
		}`, node, diskSize, hostname))
	}

	var migrationStart int64

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: compoundConfig(te.NodeName, "test-compound-a", 4),
				Check:  resource.TestCheckResourceAttr(resourceName, "node_name", te.NodeName),
			},
			{
				PreConfig: func() { migrationStart = time.Now().Add(-2 * time.Second).Unix() },
				Config:    compoundConfig(te.Node2Name, "test-compound-b", 4),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "node_name", te.Node2Name),
					resource.TestCheckResourceAttr(resourceName, "initialization.0.hostname", "test-compound-b"),
					func(*terraform.State) error {
						requireContainerOnNode(t, te, containerID, te.Node2Name)

						ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
						defer cancel()

						targetAPI := containerClientOnNode(te, te.Node2Name, containerID)

						status, err := targetAPI.GetContainerStatus(ctx)
						require.NoError(t, err)
						require.Equal(t, "running", status.Status, "container must be running on the target node")

						containerConfig, err := targetAPI.GetContainer(ctx)
						require.NoError(t, err)
						require.NotNil(t, containerConfig.Hostname)
						require.Equal(t, "test-compound-b", *containerConfig.Hostname,
							"the new configuration must be in effect on the target node")

						tasks := containerTasksSince(t, te, te.NodeName, containerID, migrationStart)
						tasks = append(tasks, containerTasksSince(t, te, te.Node2Name, containerID, migrationStart)...)

						starts, reboots, shutdowns := 0, 0, 0

						for _, task := range tasks {
							require.NotContains(t, task.Status, "migration problems",
								"the apply must not leave a failed restart migration in the task history (%s)", task.UPID)

							switch {
							case strings.Contains(task.Type, "reboot"):
								reboots++
							case strings.Contains(task.Type, "start"):
								starts++
							case strings.Contains(task.Type, "shutdown") || strings.Contains(task.Type, "stop"):
								shutdowns++
							}
						}

						require.Zero(t, reboots,
							"a reboot after the migration means the container was started with the stale configuration first: %v", tasks)
						require.Equal(t, 1, starts, "the container must be started exactly once: %v", tasks)
						require.Equal(t, 1, shutdowns, "the provider must shut the container down itself before migrating: %v", tasks)

						for _, task := range tasks {
							if !strings.Contains(task.Type, "migrate") {
								continue
							}

							sourceNode := &nodes.Client{Client: te.Client(), NodeName: te.NodeName}

							log, logErr := sourceNode.Tasks().GetTaskLog(ctx, task.UPID)
							require.NoError(t, logErr)
							require.NotContains(t, strings.Join(log, "\n"), "start container on target node",
								"the migration must not start the container; the provider starts it after writing the config")
						}

						return nil
					},
				),
			},
			{
				// The resize must run on the target, after the migration.
				PreConfig: func() { migrationStart = time.Now().Add(-2 * time.Second).Unix() },
				Config:    compoundConfig(te.NodeName, "test-compound-b", 5),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "node_name", te.NodeName),
					resource.TestCheckResourceAttr(resourceName, "disk.0.size", "5"),
					func(*terraform.State) error {
						ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
						defer cancel()

						containerConfig, err := containerClientOnNode(te, te.NodeName, containerID).GetContainer(ctx)
						require.NoError(t, err)
						require.NotNil(t, containerConfig.RootFS)
						require.NotNil(t, containerConfig.RootFS.Size)
						require.Equal(t, "5G", containerConfig.RootFS.Size.String(),
							"the root filesystem must have grown on the target node")

						for _, task := range containerTasksSince(t, te, te.Node2Name, containerID, migrationStart) {
							require.NotContains(t, task.Type, "resize",
								"the resize must not run on the source node before the migration (%s)", task.UPID)
						}

						resizes := 0

						for _, task := range containerTasksSince(t, te, te.NodeName, containerID, migrationStart) {
							if strings.Contains(task.Type, "resize") {
								resizes++
							}
						}

						require.Equal(t, 1, resizes, "the resize must run once, on the node the container ended up on")

						return nil
					},
				),
			},
		},
	})
}

// TestAccResourceContainerMigrateWithIDMap checks that node_name plus an idmap change boots the container once, with
// the idmap active, and no reboot.
func TestAccResourceContainerMigrateWithIDMap(t *testing.T) {
	te := InitEnvironment(t)

	if te.Node2Name == "" {
		t.Skip("PROXMOX_VE_ACC_NODE_2_NAME must be set")
	}

	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())
	testAccDownloadContainerTemplate(t, te, imageFileName)
	testAccDownloadContainerTemplateOnNode(t, te, imageFileName, te.Node2Name)

	containerID := 100000 + rand.Intn(99999)

	te.AddTemplateVars(map[string]any{
		"ImageFileName":   imageFileName,
		"TestContainerID": containerID,
	})

	resourceName := "proxmox_virtual_environment_container.test_migrate_idmap"

	idmapConfig := func(node, idmap string) string {
		return te.RenderConfig(fmt.Sprintf(`
		resource "proxmox_virtual_environment_container" "test_migrate_idmap" {
			node_name      = "%s"
			vm_id          = {{.TestContainerID}}
			timeout_delete = 300
			unprivileged   = true
			migrate        = true
			started        = true

			disk {
				datastore_id  = "{{.ContainerDatastoreID}}"
				size          = 4
				mount_options = []
			}
			%s
			initialization {
				hostname = "test-migrate-idmap"
			}

			operating_system {
				template_file_id = "local:vztmpl/{{.ImageFileName}}"
				type             = "alpine"
			}

			network_interface {
				name = "vmbr0"
			}
		}`, node, idmap), WithRootUser())
	}

	idmaps := `
			idmap {
				type         = "uid"
				container_id = 0
				host_id      = 100000
				size         = 1000
			}
			idmap {
				type         = "uid"
				container_id = 1000
				host_id      = 101000
				size         = 64536
			}
			idmap {
				type         = "gid"
				container_id = 0
				host_id      = 100000
				size         = 1000
			}
			idmap {
				type         = "gid"
				container_id = 1000
				host_id      = 101000
				size         = 64536
			}
`

	var migrationStart int64

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: idmapConfig(te.Node2Name, ""),
			},
			{
				// Migrate onto the primary node, where assertContainerIDMapActive can exec into the container.
				PreConfig: func() { migrationStart = time.Now().Add(-2 * time.Second).Unix() },
				Config:    idmapConfig(te.NodeName, idmaps),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "node_name", te.NodeName),
					assertContainerIDMapActive(te, containerID, "/proc/self/uid_map"),
					assertContainerIDMapActive(te, containerID, "/proc/self/gid_map"),
					func(*terraform.State) error {
						counts := containerTaskTypeCounts(t, te, containerID, migrationStart)
						t.Logf("idmap migration task types by node: %v", counts)

						for node, perNode := range counts {
							require.Zero(t, perNode["vzreboot"], "the container must not be rebooted (on %s): %v", node, counts)
						}

						require.Equal(t, 1, counts[te.Node2Name]["vzshutdown"],
							"the container must be shut down once, on the source node: %v", counts)
						require.Equal(t, 1, counts[te.NodeName]["vzstart"],
							"the container must be booted once, on the target node: %v", counts)

						return nil
					},
				),
			},
		},
	})
}
