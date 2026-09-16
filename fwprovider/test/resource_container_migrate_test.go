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

// testAccDownloadContainerTemplateOnNode downloads a template to a specific node's "local" storage.
// Unlike testAccDownloadContainerTemplate (always te.NodeName), this is needed when a test recreates
// the container on a second node — "local" storage is per-node, not shared.
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

	resource.Test(t, resource.TestCase{
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

// TestAccResourceContainerMigrateDisabledRecreates proves the default (migrate = false) still forces
// replacement on a node change, and that the replacement is a destroy-then-create rather than a bare
// create that would collide on VMID.
func TestAccResourceContainerMigrateDisabledRecreates(t *testing.T) {
	te := InitEnvironment(t)

	if te.Node2Name == "" {
		t.Skip("PROXMOX_VE_ACC_NODE_2_NAME must be set")
	}

	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())
	testAccDownloadContainerTemplate(t, te, imageFileName)
	// The replacement container is created on Node2Name, so it needs its own copy of the template too.
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

	resource.Test(t, resource.TestCase{
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

// containerClientOnNode returns a container client pinned to an arbitrary node. te.NodeClient() is
// pinned to the primary node, and te.ClusterClient() is a cluster client with no per-node accessor,
// so the nodes client is constructed directly.
func containerClientOnNode(te *Environment, node string, vmID int) *containers.Client {
	nodeClient := &nodes.Client{Client: te.Client(), NodeName: node}
	return nodeClient.Container(vmID)
}

// requireContainerOnNode asserts the container is readable on the given node, which the Terraform
// state alone cannot prove.
func requireContainerOnNode(t *testing.T, te *Environment, vmID int, node string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	_, err := containerClientOnNode(te, node, vmID).GetContainer(ctx)
	require.NoError(t, err, "container %d must be readable on node %s", vmID, node)
}

// containerMigrateConfig renders a migratable container on the given node, started as requested.
func containerMigrateConfig(te *Environment, node string, started bool) string {
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
			hostname = "test-migrate"
		}

		operating_system {
			template_file_id = "local:vztmpl/{{.ImageFileName}}"
			type             = "alpine"
		}

		network_interface {
			name = "vmbr0"
		}
	}`, node, started))
}

// TestAccResourceContainerMigrateStopped proves a stopped container is actually relocated to the
// target node — not merely recreated — when `migrate = true` and `node_name` changes.
func TestAccResourceContainerMigrateStopped(t *testing.T) {
	te := InitEnvironment(t)

	if te.Node2Name == "" {
		t.Skip("PROXMOX_VE_ACC_NODE_2_NAME must be set")
	}

	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())
	testAccDownloadContainerTemplate(t, te, imageFileName)
	// "local" storage is node-local: the container must find the template on whichever node it lands on.
	testAccDownloadContainerTemplateOnNode(t, te, imageFileName, te.Node2Name)

	containerID := 100000 + rand.Intn(99999)

	te.AddTemplateVars(map[string]any{
		"ImageFileName":   imageFileName,
		"TestContainerID": containerID,
	})

	resourceName := "proxmox_virtual_environment_container.test_migrate"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: containerMigrateConfig(te, te.NodeName, false),
				Check:  resource.TestCheckResourceAttr(resourceName, "node_name", te.NodeName),
			},
			{
				Config: containerMigrateConfig(te, te.Node2Name, false),
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

// TestAccResourceContainerMigrateRunning proves a running container is migrated with restart, and
// ends up running again on the target node.
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: containerMigrateConfig(te, te.NodeName, true),
			},
			{
				Config: containerMigrateConfig(te, te.Node2Name, true),
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

// containerHAResourceID returns the HA resource identifier ("ct:<vmid>") for a container.
func containerHAResourceID(vmID int) proxmoxtypes.HAResourceID {
	return proxmoxtypes.HAResourceID{Type: proxmoxtypes.HAResourceTypeContainer, Name: strconv.Itoa(vmID)}
}

// registerContainerHA puts the container under HA management and returns a function that removes it
// again. auto-rebalance is disabled because this cluster runs CRS with ha-auto-rebalance enabled, and
// a test that asserts "the container is on node X" cannot tolerate the scheduler relocating it on its
// own. Deregistration is also wired into t.Cleanup so a failed run leaves no HA config behind.
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
			State:         proxmoxtypes.HAResourceStateStarted,
			AutoRebalance: new(proxmoxtypes.CustomBool(false)),
			Comment:       new("terraform provider acceptance test"),
		},
		ID: haResourceID,
	})
	require.NoError(t, err, "registering container %d in HA", vmID)

	t.Cleanup(deregister)

	return deregister
}

// waitForContainerHASettled blocks until the container has been reported on node with hastate
// "started" for several consecutive reads.
//
// One "started" read is not enough. Verified live on this cluster: a freshly registered resource goes
// ” -> started -> request_start_balance -> started over roughly 20 seconds while the CRM adopts it and
// CRS considers placement. Migrating during that window races the manager, so the test waits for the
// state to hold before it asks Terraform to move anything.
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

// TestAccResourceContainerMigrateHA proves the provider waits for the HA manager to finish relocating
// a container before the apply returns.
//
// This is the case the plain migrate tests cannot cover: for an HA-managed container PVE converts the
// migrate call into an HA request whose task reports TASK OK within ~2s, while the container is still
// running on the source node. Measured on this cluster, it reaches the target node at ~18s (still
// hastate "migrate", still stopped there) and is running at ~33s. Without the wait, the apply returns
// during that window and the post-migration checks below fail — the container is not yet readable on
// the target node at all.
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

	var deregisterHA func()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: containerMigrateConfig(te, te.NodeName, true),
			},
			{
				PreConfig: func() {
					deregisterHA = registerContainerHA(t, te, containerID)
					waitForContainerHASettled(t, te, containerID, te.NodeName, 3*time.Minute)
				},
				Config: containerMigrateConfig(te, te.Node2Name, true),
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

						// Read from the target node, never from /cluster/resources: its status field lags a
						// migration by tens of seconds and reports the container running while the node that
						// now owns it still has it stopped.
						status, err := containerClientOnNode(te, te.Node2Name, containerID).GetContainerStatus(ctx)
						require.NoError(t, err)
						require.Equal(t, "running", status.Status,
							"container must already be running on the target node when the apply returns, "+
								"not merely queued for relocation by the HA manager")

						return nil
					},
				),
			},
			{
				// Hand the container back before the framework destroys it: HA would restart a container
				// stopped for deletion, so the destroy has to run against an unmanaged guest.
				PreConfig: func() { deregisterHA() },
				Config:    containerMigrateConfig(te, te.Node2Name, true),
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

// containerTasksSince returns the tasks PVE recorded for a container on one node at or after since.
// Task history is per-node, so a migration's tasks are split between the source and the target.
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

// TestAccResourceContainerMigrateWithConfigChange covers a single apply that changes node_name and
// configuration at the same time.
//
// A restart migration starts the container on the target from the configuration as it stands when the
// migration runs, so the pending hostname would only take effect after a second bounce. The provider
// must instead shut the container down, migrate it cold, write the configuration on the target and
// start it — one restart in total.
//
// Counting the boots is what makes this test meaningful: an assertion that the container is simply
// running on the target with the new hostname passes under both designs.
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

	// The hostname is a configuration change that the provider marks as requiring a restart, so under
	// the old ordering it produced a second bounce after the restart migration.
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

	resource.Test(t, resource.TestCase{
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

						// Read from the node itself: /cluster/resources reports stale status right after
						// a migration.
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

						// The offline migration must not have started the container: that is the restart
						// migration this path exists to avoid.
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
				// Migrating back while also growing the root filesystem. The resize is a separate API
				// call, not part of the config PUT, so it has to run only once the container is on its
				// final node: resizing first would copy the extra space across the wire, and would fail
				// whenever only the target has room for it.
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
