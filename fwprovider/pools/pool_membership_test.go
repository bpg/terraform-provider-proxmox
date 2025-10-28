package pools_test

import (
	"context"
	"fmt"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
	"github.com/bpg/terraform-provider-proxmox/proxmox/pools"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestAccPoolMembershipContainer(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	imageFileName := gofakeit.Word() + "-ubuntu-24.04-standard_24.04-2_amd64.tar.zst"
	accTestContainerID := 100000 + rand.Intn(99999)
	accTestContainerID2 := 100000 + rand.Intn(99999)

	te.AddTemplateVars(map[string]interface{}{
		"ImageFileName":    imageFileName,
		"TestContainerID":  accTestContainerID,
		"TestContainerID2": accTestContainerID2,
	})

	err := te.NodeStorageClient().DownloadFileByURL(context.Background(), &storage.DownloadURLPostRequestBody{
		Content:  ptr.Ptr("vztmpl"),
		FileName: &imageFileName,
		Node:     &te.NodeName,
		Storage:  &te.DatastoreID,
		URL:      ptr.Ptr(fmt.Sprintf("%s/images/system/ubuntu-24.04-standard_24.04-2_amd64.tar.zst", te.ContainerImagesServer)),
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		e := te.NodeStorageClient().DeleteDatastoreFile(context.Background(), fmt.Sprintf("vztmpl/%s", imageFileName))
		require.NoError(t, e)
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_pool" "test_pool" {
					pool_id = "test-pool"
				}

				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID}}
					disk {
						datastore_id = "local-lvm"
						size         = 4
						mount_options = []
					}
					operating_system {
						template_file_id = "local:vztmpl/{{ .ImageFileName }}"
						type             = "ubuntu"
					}

					initialization {
						hostname = "test"
					}
					started = false
				}
				
				resource "proxmox_virtual_environment_pool_membership" "test_ct_test_pool_membership" {
					pool_id = proxmox_virtual_environment_pool.test_pool.id
					vm_id = proxmox_virtual_environment_container.test_container.id
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes(
						"proxmox_virtual_environment_pool_membership.test_ct_test_pool_membership",
						map[string]string{
							"pool_id": "test-pool",
							"vm_id":   strconv.Itoa(accTestContainerID),
						},
					),
					func(state *terraform.State) error {
						ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
						defer cancel()
						exists, poolCheckErr := testAccCheckPoolContainsMember(
							ctx,
							te.PoolsClient(),
							"test-pool",
							strconv.Itoa(accTestContainerID),
							"lxc",
						)
						require.NoError(t, poolCheckErr, "couldn't get the pool")
						require.True(t, exists)
						return nil
					},
				),
			},
			{
				// destroy "test-pool" and create "test-pool-2"
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_pool" "test_pool-2" {
					pool_id = "test-pool-2"
				}

				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID}}
					disk {
						datastore_id = "local-lvm"
						size         = 4
						mount_options = []
					}
					operating_system {
						template_file_id = "local:vztmpl/{{ .ImageFileName }}"
						type             = "ubuntu"
					}

					initialization {
						hostname = "test"
					}
					started = false
				}
				
				resource "proxmox_virtual_environment_pool_membership" "test_ct_test_pool_membership" {
					pool_id = proxmox_virtual_environment_pool.test_pool-2.id
					vm_id = proxmox_virtual_environment_container.test_container.id
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes(
						"proxmox_virtual_environment_pool_membership.test_ct_test_pool_membership",
						map[string]string{
							"pool_id": "test-pool-2",
							"vm_id":   strconv.Itoa(accTestContainerID),
						},
					),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"proxmox_virtual_environment_pool_membership.test_ct_test_pool_membership",
							plancheck.ResourceActionDestroyBeforeCreate,
						),
					},
				},
			},
			{
				// destroy "test-container" and create "test-container-2"
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_pool" "test_pool-2" {
					pool_id = "test-pool-2"
				}

				resource "proxmox_virtual_environment_container" "test_container-2" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID2}}
					disk {
						datastore_id = "local-lvm"
						size         = 4
						mount_options = []
					}
					operating_system {
						template_file_id = "local:vztmpl/{{ .ImageFileName }}"
						type             = "ubuntu"
					}

					initialization {
						hostname = "test"
					}
					started = false
				}
				
				resource "proxmox_virtual_environment_pool_membership" "test_ct_test_pool_membership" {
					pool_id = proxmox_virtual_environment_pool.test_pool-2.id
					vm_id = proxmox_virtual_environment_container.test_container-2.id
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes(
						"proxmox_virtual_environment_pool_membership.test_ct_test_pool_membership",
						map[string]string{
							"pool_id": "test-pool-2",
							"vm_id":   strconv.Itoa(accTestContainerID2),
						},
					),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"proxmox_virtual_environment_pool_membership.test_ct_test_pool_membership",
							plancheck.ResourceActionDestroyBeforeCreate,
						),
					},
				},
			},
			{
				// delete membership
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_pool" "test_pool-2" {
					pool_id = "test-pool-2"
				}

				resource "proxmox_virtual_environment_container" "test_container-2" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID2}}
					disk {
						datastore_id = "local-lvm"
						size         = 4
						mount_options = []
					}
					operating_system {
						template_file_id = "local:vztmpl/{{ .ImageFileName }}"
						type             = "ubuntu"
					}

					initialization {
						hostname = "test"
					}
					started = false
				}`),
				Check: resource.ComposeTestCheckFunc(
					func(state *terraform.State) error {
						ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
						defer cancel()
						exists, poolCheckErr := testAccCheckPoolContainsMember(
							ctx,
							te.PoolsClient(),
							"test-pool-2",
							strconv.Itoa(accTestContainerID2),
							"lxc",
						)
						require.NoError(t, poolCheckErr, "couldn't get the pool")
						require.False(t, exists)
						return nil
					},
				),
			},
		},
	})
}

func testAccCheckPoolContainsMember(ctx context.Context, client *pools.Client, poolId, memberId, memberType string) (bool, error) {
	pool, err := client.GetPool(ctx, poolId)
	if err != nil {
		return false, err
	}
	exists := false
	for _, member := range pool.Members {
		switch memberType {
		case "lxc", "qemu":
			if member.VMID != nil && member.Type == memberType && strconv.Itoa(*member.VMID) == memberId {
				exists = true
				break
			}
		case "storage":
			if member.DatastoreID != nil && member.Type == memberType && *member.DatastoreID == memberId {
				exists = true
				break
			}
		}
	}
	return exists, err
}
