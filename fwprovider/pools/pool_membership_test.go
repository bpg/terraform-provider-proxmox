//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package pools_test

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/bpg/terraform-provider-proxmox/proxmox/pools"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
)

func TestAccPoolMembershipContainer(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	imageFileName := gofakeit.Word() + "-ubuntu-24.04-standard_24.04-2_amd64.tar.zst"
	accTestContainerID := 100000 + rand.Intn(99999)
	accTestContainerID2 := 100000 + rand.Intn(99999)
	accTestPoolName := gofakeit.Word()

	accTestPoolName2 := gofakeit.Word()
	for accTestPoolName == accTestPoolName2 {
		accTestPoolName2 = gofakeit.Word()
	}

	te.AddTemplateVars(map[string]interface{}{
		"ImageFileName":    imageFileName,
		"TestContainerID":  accTestContainerID,
		"TestContainerID2": accTestContainerID2,
		"TestPoolName":     accTestPoolName,
		"TestPoolName2":    accTestPoolName2,
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
					pool_id = "{{ .TestPoolName }}"
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
				
				resource "proxmox_virtual_environment_pool_membership" "pool_membership" {
					pool_id = proxmox_virtual_environment_pool.test_pool.id
					vm_id = proxmox_virtual_environment_container.test_container.id
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes(
						"proxmox_virtual_environment_pool_membership.pool_membership",
						map[string]string{
							"id":      fmt.Sprintf("%s/vm/%d", accTestPoolName, accTestContainerID),
							"pool_id": accTestPoolName,
							"vm_id":   strconv.Itoa(accTestContainerID),
							"type":    "vm",
						},
					),
					testAccCheckPoolContainMember(t, te, accTestPoolName, strconv.Itoa(accTestContainerID), "lxc", true),
				),
			},
			{
				ResourceName:      "proxmox_virtual_environment_pool_membership.pool_membership",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// change pool
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_pool" "test_pool" {
					pool_id = "{{ .TestPoolName }}"
				}

				resource "proxmox_virtual_environment_pool" "test_pool-2" {
					pool_id = "{{ .TestPoolName2 }}"
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
				
				resource "proxmox_virtual_environment_pool_membership" "pool_membership" {
					pool_id = proxmox_virtual_environment_pool.test_pool-2.id
					vm_id = proxmox_virtual_environment_container.test_container.id
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes(
						"proxmox_virtual_environment_pool_membership.pool_membership",
						map[string]string{
							"id":      fmt.Sprintf("%s/vm/%d", accTestPoolName2, accTestContainerID),
							"pool_id": accTestPoolName2,
							"vm_id":   strconv.Itoa(accTestContainerID),
							"type":    "vm",
						},
					),
					testAccCheckPoolContainMember(t, te, accTestPoolName, strconv.Itoa(accTestContainerID), "lxc", false),
					testAccCheckPoolContainMember(t, te, accTestPoolName2, strconv.Itoa(accTestContainerID), "lxc", true),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"proxmox_virtual_environment_pool_membership.pool_membership",
							plancheck.ResourceActionDestroyBeforeCreate,
						),
					},
				},
			},
			{
				// destroy "test-container" and create "test-container-2"
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_pool" "test_pool-2" {
					pool_id = "{{ .TestPoolName2 }}"
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
				
				resource "proxmox_virtual_environment_pool_membership" "pool_membership" {
					pool_id = proxmox_virtual_environment_pool.test_pool-2.id
					vm_id = proxmox_virtual_environment_container.test_container-2.id
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes(
						"proxmox_virtual_environment_pool_membership.pool_membership",
						map[string]string{
							"id":      fmt.Sprintf("%s/vm/%d", accTestPoolName2, accTestContainerID2),
							"pool_id": accTestPoolName2,
							"vm_id":   strconv.Itoa(accTestContainerID2),
							"type":    "vm",
						},
					),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"proxmox_virtual_environment_pool_membership.pool_membership",
							plancheck.ResourceActionDestroyBeforeCreate,
						),
					},
				},
			},
			{
				// delete membership
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_pool" "test_pool-2" {
					pool_id = "{{ .TestPoolName2 }}"
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
				Check: testAccCheckPoolContainMember(t, te, accTestPoolName2, strconv.Itoa(accTestContainerID2), "lxc", false),
			},
		},
	})
}

func TestAccPoolMembershipVm(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)
	accTestVmId := 100000 + rand.Intn(99999)
	accTestVmId2 := 100000 + rand.Intn(99999)
	accTestPoolName := gofakeit.Word()

	accTestPoolName2 := gofakeit.Word()
	for accTestPoolName == accTestPoolName2 {
		accTestPoolName2 = gofakeit.Word()
	}

	te.AddTemplateVars(map[string]interface{}{
		"TestVMID":      accTestVmId,
		"TestVMID2":     accTestVmId2,
		"TestPoolName":  accTestPoolName,
		"TestPoolName2": accTestPoolName2,
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "test_vm1" {
					vm_id = {{.TestVMID}}
					node_name = "{{.NodeName}}"
					started   = false
				}
				
				resource "proxmox_virtual_environment_pool" "test_pool" {
					pool_id = "{{.TestPoolName}}"
				}
				
				resource "proxmox_virtual_environment_pool_membership" "pool_membership" {
					pool_id = proxmox_virtual_environment_pool.test_pool.id
					vm_id = proxmox_virtual_environment_vm.test_vm1.id
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes(
						"proxmox_virtual_environment_pool_membership.pool_membership",
						map[string]string{
							"id":      fmt.Sprintf("%s/vm/%d", accTestPoolName, accTestVmId),
							"pool_id": accTestPoolName,
							"vm_id":   strconv.Itoa(accTestVmId),
							"type":    "vm",
						},
					),
					testAccCheckPoolContainMember(t, te, accTestPoolName, strconv.Itoa(accTestVmId), "qemu", true),
				),
			},
			{
				ResourceName:      "proxmox_virtual_environment_pool_membership.pool_membership",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{ // change pool
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_pool" "test_pool" {
					pool_id = "{{ .TestPoolName }}"
				}

				resource "proxmox_virtual_environment_pool" "test_pool-2" {
					pool_id = "{{ .TestPoolName2 }}"
				}

				resource "proxmox_virtual_environment_vm" "test_vm1" {
					vm_id = {{.TestVMID}}
					node_name = "{{.NodeName}}"
					started   = false
				}
				
				resource "proxmox_virtual_environment_pool_membership" "pool_membership" {
					pool_id = proxmox_virtual_environment_pool.test_pool-2.id
					vm_id = proxmox_virtual_environment_vm.test_vm1.id
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes(
						"proxmox_virtual_environment_pool_membership.pool_membership",
						map[string]string{
							"id":      fmt.Sprintf("%s/vm/%d", accTestPoolName2, accTestVmId),
							"pool_id": accTestPoolName2,
							"vm_id":   strconv.Itoa(accTestVmId),
							"type":    "vm",
						},
					),
					testAccCheckPoolContainMember(t, te, accTestPoolName, strconv.Itoa(accTestVmId), "qemu", false),
					testAccCheckPoolContainMember(t, te, accTestPoolName2, strconv.Itoa(accTestVmId), "qemu", true),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"proxmox_virtual_environment_pool_membership.pool_membership",
							plancheck.ResourceActionDestroyBeforeCreate,
						),
					},
				},
			},
			{ // destroy "test-vm1" and create "test-vm2"
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_pool" "test_pool-2" {
					pool_id = "{{ .TestPoolName2 }}"
				}

				resource "proxmox_virtual_environment_vm" "test_vm2" {
					vm_id = {{.TestVMID2}}
					node_name = "{{.NodeName}}"
					started   = false
				}
				
				resource "proxmox_virtual_environment_pool_membership" "pool_membership" {
					pool_id = proxmox_virtual_environment_pool.test_pool-2.id
					vm_id = proxmox_virtual_environment_vm.test_vm2.id
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes(
						"proxmox_virtual_environment_pool_membership.pool_membership",
						map[string]string{
							"id":      fmt.Sprintf("%s/vm/%d", accTestPoolName2, accTestVmId2),
							"pool_id": accTestPoolName2,
							"vm_id":   strconv.Itoa(accTestVmId2),
							"type":    "vm",
						},
					),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"proxmox_virtual_environment_pool_membership.pool_membership",
							plancheck.ResourceActionDestroyBeforeCreate,
						),
					},
				},
			},
			{ // delete membership
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_pool" "test_pool-2" {
					pool_id = "{{ .TestPoolName2 }}"
				}

				resource "proxmox_virtual_environment_vm" "test_vm2" {
					vm_id = {{.TestVMID2}}
					node_name = "{{.NodeName}}"
					started   = false
				}`),
				Check: testAccCheckPoolContainMember(t, te, accTestPoolName2, strconv.Itoa(accTestVmId2), "qemu", false),
			},
		},
	})
}

func TestAccPoolMembershipStorage(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)
	accTestPoolName := gofakeit.Word()

	accTestPoolName2 := gofakeit.Word()
	for accTestPoolName == accTestPoolName2 {
		accTestPoolName2 = gofakeit.Word()
	}

	te.AddTemplateVars(map[string]interface{}{
		"TestPoolName":  accTestPoolName,
		"TestPoolName2": accTestPoolName2,
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_pool" "test_pool" {
					pool_id = "{{.TestPoolName}}"
				}
				
				resource "proxmox_virtual_environment_pool_membership" "pool_membership" {
					pool_id = proxmox_virtual_environment_pool.test_pool.id
					storage_id = "{{ .DatastoreID }}"
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes(
						"proxmox_virtual_environment_pool_membership.pool_membership",
						map[string]string{
							"id":         fmt.Sprintf("%s/storage/%s", accTestPoolName, "local"),
							"pool_id":    accTestPoolName,
							"storage_id": "local",
							"type":       "storage",
						},
					),
					testAccCheckPoolContainMember(t, te, accTestPoolName, "local", "storage", true),
				),
			},
			{
				ResourceName:      "proxmox_virtual_environment_pool_membership.pool_membership",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{ // change pool
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_pool" "test_pool" {
					pool_id = "{{ .TestPoolName }}"
				}
				resource "proxmox_virtual_environment_pool" "test_pool-2" {
					pool_id = "{{ .TestPoolName2 }}"
				}
				
				resource "proxmox_virtual_environment_pool_membership" "pool_membership" {
					pool_id = proxmox_virtual_environment_pool.test_pool-2.id
					storage_id = "{{ .DatastoreID }}"
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes(
						"proxmox_virtual_environment_pool_membership.pool_membership",
						map[string]string{
							"id":         fmt.Sprintf("%s/storage/%s", accTestPoolName2, "local"),
							"pool_id":    accTestPoolName2,
							"storage_id": "local",
							"type":       "storage",
						},
					),
					testAccCheckPoolContainMember(t, te, accTestPoolName2, "local", "storage", true),
					testAccCheckPoolContainMember(t, te, accTestPoolName, "local", "storage", false),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"proxmox_virtual_environment_pool_membership.pool_membership",
							plancheck.ResourceActionDestroyBeforeCreate,
						),
					},
				},
			},
			{ // delete membership
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_pool" "test_pool-2" {
					pool_id = "{{ .TestPoolName2 }}"
				}`),
				Check: testAccCheckPoolContainMember(t, te, accTestPoolName2, "local", "storage", false),
			},
		},
	})
}

func TestAccPoolMembership_Validators(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				PlanOnly: true,
				Config: `resource "proxmox_virtual_environment_pool_membership" "test" {
					pool_id = "test_pool"
					storage_id = "local"
					vm_id = 1234
				}`,
				ExpectError: regexp.MustCompile(`.*Error: Invalid Attribute Combination`),
			},
			{
				PlanOnly: true,
				Config: `resource "proxmox_virtual_environment_pool_membership" "test" {
					pool_id = "test_pool"
				}`,
				ExpectError: regexp.MustCompile(`.*Error: Missing Attribute Configuration`),
			},
			{
				PlanOnly: true,
				Config: `resource "proxmox_virtual_environment_pool_membership" "test" {
					pool_id = "test_pool"
					vm_id = 1234
				}`,
				ExpectNonEmptyPlan: true,
			},
			{
				PlanOnly: true,
				Config: `resource "proxmox_virtual_environment_pool_membership" "test" {
					pool_id = "test_pool"
					storage_id = "local"
				}`,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func checkPoolContainsMember(ctx context.Context, client *pools.Client, poolId, memberId, memberType string) (bool, error) {
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

func testAccCheckPoolContainMember(t *testing.T, te *test.Environment, poolName, memberId, memberType string, shouldExist bool) resource.TestCheckFunc {
	t.Helper()

	return func(state *terraform.State) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		exists, poolCheckErr := checkPoolContainsMember(
			ctx,
			te.PoolsClient(),
			poolName,
			memberId,
			memberType,
		)
		require.NoError(t, poolCheckErr, "couldn't get the pool")
		require.Equal(t, shouldExist, exists)

		return nil
	}
}
