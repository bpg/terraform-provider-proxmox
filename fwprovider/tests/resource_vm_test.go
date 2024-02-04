/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

const (
	accTestVMName      = "proxmox_virtual_environment_vm.test_vm"
	accTestVMCloneName = "proxmox_virtual_environment_vm.test_vm_clone"
)

func TestAccResourceVM(t *testing.T) {
	t.Parallel()

	accProviders := testAccMuxProviders(context.Background(), t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: accProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVMCreateConfig(false),
				Check:  testAccResourceVMCreateCheck(t),
			},
			{
				Config: testAccResourceVMCreateConfig(true) + testAccResourceVMCreateCloneConfig(),
				Check:  testAccResourceVMCreateCloneCheck(t),
			},
		},
	})
}

func testAccResourceVMCreateConfig(isTemplate bool) string {
	return fmt.Sprintf(`
resource "proxmox_virtual_environment_vm" "test_vm" {
  node_name = "%s"
  vm_id     = 2100
  template  = %t
  started = false

  disk {
	file_format= "raw"
	datastore_id = "local-lvm"
	interface    = "virtio0"
    size         = 8
  }
  
}
`, accTestNodeName, isTemplate)
}

func testAccResourceVMCreateCheck(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return resource.ComposeTestCheckFunc(
		func(*terraform.State) error {
			err := getNodesClient().VM(2100).WaitForVMStatus(context.Background(), "stopped", 10, 1)
			require.NoError(t, err, "vm did not start")
			return nil
		},
	)
}

func testAccResourceVMCreateCloneConfig() string {
	return fmt.Sprintf(`
resource "proxmox_virtual_environment_vm" "test_vm_clone" {
  depends_on = [proxmox_virtual_environment_vm.test_vm]

  node_name = "%s"
  vm_id     = 2101
  started = false

  clone {
	vm_id = 2100
  }
}
`, accTestNodeName)
}

func testAccResourceVMCreateCloneCheck(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return resource.ComposeTestCheckFunc(
		func(*terraform.State) error {
			err := getNodesClient().VM(2101).WaitForVMStatus(context.Background(), "stopped", 20, 1)
			require.NoError(t, err, "vm did not start")
			return nil
		},
	)
}
