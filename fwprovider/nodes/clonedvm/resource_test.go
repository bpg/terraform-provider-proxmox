//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package clonedvm_test

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

func TestAccResourceClonedVM_InheritAndDelete(t *testing.T) {
	if utils.GetAnyStringEnv("TF_ACC") == "" {
		t.Skip("Acceptance tests are disabled")
	}

	t.Parallel()

	te := test.InitEnvironment(t)

	config := te.RenderConfig(`
		resource "proxmox_virtual_environment_download_file" "cloud_image" {
			content_type = "iso"
			datastore_id = "local"
			node_name    = "{{.NodeName}}"
			url          = "{{.CloudImagesServer}}/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-amd64.img"
			overwrite_unmanaged = true
		}

		resource "proxmox_virtual_environment_vm" "template_vm" {
			node_name = "{{.NodeName}}"
			started   = false
			template  = true

			disk {
				datastore_id = "local-lvm"
				file_id      = proxmox_virtual_environment_download_file.cloud_image.id
				interface    = "virtio0"
				size         = 16
			}

			cpu {
				cores = 1
			}

			memory {
				dedicated = 1024
			}

			network_device {
				model  = "virtio"
				bridge = "vmbr0"
			}

			network_device {
				model  = "virtio"
				bridge = "vmbr0"
			}
		}

		resource "proxmox_virtual_environment_cloned_vm" "keep_inherited" {
			node_name = "{{.NodeName}}"
			name      = "fwk-cloned-keep"

			clone = {
				source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
			}

			network = {
				net0 = {
					bridge = "vmbr0"
					model  = "virtio"
				}
			}
		}

		resource "proxmox_virtual_environment_cloned_vm" "delete_inherited" {
			node_name = "{{.NodeName}}"
			name      = "fwk-cloned-delete"

			clone = {
				source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
			}

			network = {
				net0 = {
					bridge = "vmbr0"
					model  = "virtio"
				}
			}

			delete = {
				network = ["net1"]
			}
		}
	`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkNetworkSlot(te, "proxmox_virtual_environment_cloned_vm.keep_inherited", "net1", true),
					checkNetworkSlot(te, "proxmox_virtual_environment_cloned_vm.delete_inherited", "net1", false),
				),
			},
		},
	})
}

func TestAccResourceClonedVM_StopManagingDoesNotDelete(t *testing.T) {
	if utils.GetAnyStringEnv("TF_ACC") == "" {
		t.Skip("Acceptance tests are disabled")
	}

	t.Parallel()

	te := test.InitEnvironment(t)

	baseConfig := `
		resource "proxmox_virtual_environment_download_file" "cloud_image" {
			content_type = "iso"
			datastore_id = "local"
			node_name    = "{{.NodeName}}"
			url          = "{{.CloudImagesServer}}/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-amd64.img"
			overwrite_unmanaged = true
		}

		resource "proxmox_virtual_environment_vm" "template_vm" {
			node_name = "{{.NodeName}}"
			started   = false
			template  = true

			disk {
				datastore_id = "local-lvm"
				file_id      = proxmox_virtual_environment_download_file.cloud_image.id
				interface    = "virtio0"
				size         = 16
			}

			cpu {
				cores = 1
			}

			memory {
				dedicated = 1024
			}

			network_device {
				model  = "virtio"
				bridge = "vmbr0"
			}
		}
		`

	withManaged := te.RenderConfig(baseConfig + `
		resource "proxmox_virtual_environment_cloned_vm" "unmanaged" {
			node_name = "{{.NodeName}}"
			name      = "fwk-cloned-unmanage"

			clone = {
				source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
			}

			network = {
				net0 = {
					bridge = "vmbr0"
					model  = "virtio"
				}
			}
		}
	`)

	withoutManaged := te.RenderConfig(baseConfig + `
		resource "proxmox_virtual_environment_cloned_vm" "unmanaged" {
			node_name = "{{.NodeName}}"
			name      = "fwk-cloned-unmanage"

			clone = {
				source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
			}
		}
	`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: withManaged,
				Check: resource.ComposeTestCheckFunc(
					checkNetworkSlot(te, "proxmox_virtual_environment_cloned_vm.unmanaged", "net0", true),
				),
			},
			{
				Config: withoutManaged,
				Check: resource.ComposeTestCheckFunc(
					checkNetworkSlot(te, "proxmox_virtual_environment_cloned_vm.unmanaged", "net0", true),
				),
			},
		},
	})
}

func TestAccResourceClonedVM_MapKeyStability(t *testing.T) {
	if utils.GetAnyStringEnv("TF_ACC") == "" {
		t.Skip("Acceptance tests are disabled")
	}

	t.Parallel()

	te := test.InitEnvironment(t)

	baseConfig := `
		resource "proxmox_virtual_environment_download_file" "cloud_image" {
			content_type = "iso"
			datastore_id = "local"
			node_name    = "{{.NodeName}}"
			url          = "{{.CloudImagesServer}}/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-amd64.img"
			overwrite_unmanaged = true
		}

		resource "proxmox_virtual_environment_vm" "template_vm" {
			node_name = "{{.NodeName}}"
			started   = false
			template  = true

			disk {
				datastore_id = "local-lvm"
				file_id      = proxmox_virtual_environment_download_file.cloud_image.id
				interface    = "virtio0"
				size         = 16
			}

			cpu {
				cores = 1
			}

			memory {
				dedicated = 1024
			}

			network_device {
				model  = "virtio"
				bridge = "vmbr0"
			}

			network_device {
				model  = "virtio"
				bridge = "vmbr0"
			}
		}
		`

	initialConfig := te.RenderConfig(baseConfig + `
		resource "proxmox_virtual_environment_cloned_vm" "stability" {
			node_name = "{{.NodeName}}"
			name      = "fwk-cloned-stability"

			clone = {
				source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
			}

			network = {
				net0 = {
					bridge = "vmbr0"
					model  = "virtio"
				}
			}
		}
	`)

	updatedConfig := te.RenderConfig(baseConfig + `
		resource "proxmox_virtual_environment_cloned_vm" "stability" {
			node_name = "{{.NodeName}}"
			name      = "fwk-cloned-stability"

			clone = {
				source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
			}

			network = {
				net0 = {
					bridge = "vmbr0"
					model  = "e1000"
					tag    = 100
				}
			}
		}
	`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					checkNetworkSlot(te, "proxmox_virtual_environment_cloned_vm.stability", "net0", true),
					checkNetworkSlot(te, "proxmox_virtual_environment_cloned_vm.stability", "net1", true),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					checkNetworkSlot(te, "proxmox_virtual_environment_cloned_vm.stability", "net0", true),
					checkNetworkSlot(te, "proxmox_virtual_environment_cloned_vm.stability", "net1", true),
				),
			},
		},
	})
}

func checkNetworkSlot(te *test.Environment, resourceName, slot string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}

		nodeName := rs.Primary.Attributes["node_name"]
		idStr := rs.Primary.Attributes["id"]

		if nodeName == "" || idStr == "" {
			return fmt.Errorf("resource %s missing node_name or id in state", resourceName)
		}

		vmid, err := strconv.Atoi(idStr)
		if err != nil {
			return err
		}

		ctx := context.Background()
		config, err := te.NodeClient().VM(vmid).GetVM(ctx)
		if err != nil {
			return err
		}

		found := networkSlotPresent(config, slot)

		if expected && !found {
			return fmt.Errorf("expected slot %s to exist for %s", slot, resourceName)
		}

		if !expected && found {
			return fmt.Errorf("expected slot %s to be absent for %s", slot, resourceName)
		}

		return nil
	}
}

func networkSlotPresent(config *vms.GetResponseData, slot string) bool {
	idx, ok := slotIndex(slot, "net")
	if !ok {
		return false
	}

	val := reflect.ValueOf(config)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	field := val.FieldByName(fmt.Sprintf("NetworkDevice%d", idx))
	if !field.IsValid() || field.IsNil() {
		return false
	}

	return true
}

func TestAccResourceClonedVM_MemoryConfiguration(t *testing.T) {
	if utils.GetAnyStringEnv("TF_ACC") == "" {
		t.Skip("Acceptance tests are disabled")
	}

	t.Parallel()

	te := test.InitEnvironment(t)

	baseConfig := `
		resource "proxmox_virtual_environment_download_file" "cloud_image" {
			content_type = "iso"
			datastore_id = "local"
			node_name    = "{{.NodeName}}"
			url          = "{{.CloudImagesServer}}/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-amd64.img"
			overwrite_unmanaged = true
		}

		resource "proxmox_virtual_environment_vm" "template_vm" {
			node_name = "{{.NodeName}}"
			started   = false
			template  = true

			disk {
				datastore_id = "local-lvm"
				file_id      = proxmox_virtual_environment_download_file.cloud_image.id
				interface    = "virtio0"
				size         = 16
			}

			cpu {
				cores = 1
			}

			memory {
				dedicated = 512
			}

			network_device {
				model  = "virtio"
				bridge = "vmbr0"
			}
		}
		`

	initialConfig := te.RenderConfig(baseConfig + `
		resource "proxmox_virtual_environment_cloned_vm" "memory_test" {
			node_name = "{{.NodeName}}"
			name      = "fwk-cloned-memory"

			clone = {
				source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
			}

			memory = {
				maximum = 2048
				minimum = 1024
				shares  = 1500
			}
		}
	`)

	updatedConfig := te.RenderConfig(baseConfig + `
		resource "proxmox_virtual_environment_cloned_vm" "memory_test" {
			node_name = "{{.NodeName}}"
			name      = "fwk-cloned-memory"

			clone = {
				source_vm_id = proxmox_virtual_environment_vm.template_vm.vm_id
			}

			memory = {
				maximum = 4096
				minimum = 2048
				shares  = 2000
			}
		}
	`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					checkMemoryConfig(te, "proxmox_virtual_environment_cloned_vm.memory_test", 2048, 1024, 1500),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					checkMemoryConfig(te, "proxmox_virtual_environment_cloned_vm.memory_test", 4096, 2048, 2000),
				),
			},
		},
	})
}

func checkMemoryConfig(te *test.Environment, resourceName string, expectedMaximum, expectedMinimum, expectedShares int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}

		nodeName := rs.Primary.Attributes["node_name"]
		idStr := rs.Primary.Attributes["id"]

		if nodeName == "" || idStr == "" {
			return fmt.Errorf("resource %s missing node_name or id in state", resourceName)
		}

		vmid, err := strconv.Atoi(idStr)
		if err != nil {
			return err
		}

		ctx := context.Background()
		config, err := te.NodeClient().VM(vmid).GetVM(ctx)
		if err != nil {
			return err
		}

		// Check maximum memory (Proxmox API: 'memory', our naming: 'maximum', SDK: 'dedicated')
		if config.DedicatedMemory == nil {
			return fmt.Errorf("DedicatedMemory (maximum) is nil for %s", resourceName)
		}
		if int64(*config.DedicatedMemory) != int64(expectedMaximum) {
			return fmt.Errorf("expected maximum memory %d, got %d for %s", expectedMaximum, *config.DedicatedMemory, resourceName)
		}

		// Check minimum memory (Proxmox API: 'balloon', our naming: 'minimum', SDK: 'floating')
		if config.FloatingMemory == nil {
			return fmt.Errorf("FloatingMemory (minimum) is nil for %s", resourceName)
		}
		if int64(*config.FloatingMemory) != int64(expectedMinimum) {
			return fmt.Errorf("expected minimum memory %d, got %d for %s", expectedMinimum, *config.FloatingMemory, resourceName)
		}

		// Check shares (Proxmox API: 'shares', our naming: 'shares', SDK: 'shared')
		if config.FloatingMemoryShares == nil {
			return fmt.Errorf("FloatingMemoryShares (shares) is nil for %s", resourceName)
		}
		if *config.FloatingMemoryShares != expectedShares {
			return fmt.Errorf("expected shares %d, got %d for %s", expectedShares, *config.FloatingMemoryShares, resourceName)
		}

		return nil
	}
}

func slotIndex(slot string, prefix string) (int, bool) {
	if !strings.HasPrefix(slot, prefix) {
		return 0, false
	}

	idx, err := strconv.Atoi(strings.TrimPrefix(slot, prefix))
	if err != nil || idx < 0 {
		return 0, false
	}

	return idx, true
}
