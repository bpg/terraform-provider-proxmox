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

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
)

const (
	accTestContainerName      = "proxmox_virtual_environment_container.test_container"
	accTestContainerCloneName = "proxmox_virtual_environment_container.test_container_clone"
)

func TestAccResourceContainer(t *testing.T) {
	accProviders := testAccMuxProviders(context.Background(), t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: accProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceContainerCreateConfig(false),
				Check:  testAccResourceContainerCreateCheck(t),
			},
			{
				Config: testAccResourceContainerCreateConfig(true) + testAccResourceContainerCreateCloneConfig(),
				Check:  testAccResourceContainerCreateCloneCheck(t),
			},
		},
	})
}

func testAccResourceContainerCreateConfig(isTemplate bool) string {
	return fmt.Sprintf(`
resource "proxmox_virtual_environment_download_file" "ubuntu_container_template" {
	content_type = "vztmpl"
	datastore_id = "local"
	node_name    = "pve"
	url = "http://download.proxmox.com/images/system/ubuntu-23.04-standard_23.04-1_amd64.tar.zst"
    overwrite_unmanaged = true
}
resource "proxmox_virtual_environment_container" "test_container" {
  node_name = "%s"
  vm_id     = 1100
  template  = %t

  disk {
    datastore_id = "local-lvm"
    size         = 8
  }

  description = <<-EOT
    my
    description
    value
  EOT

  initialization {
    hostname = "test"

    ip_config {
      ipv4 {
        address = "dhcp"
      }
    }
  }

  network_interface {
    name = "vmbr0"
  }

  operating_system {
	template_file_id = proxmox_virtual_environment_download_file.ubuntu_container_template.id
    type             = "ubuntu"
  }
}
`, accTestNodeName, isTemplate)
}

func testAccResourceContainerCreateCheck(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(accTestContainerName, "description", "my\ndescription\nvalue\n"),
		func(*terraform.State) error {
			err := getNodesClient().Container(1100).WaitForContainerStatus(context.Background(), "running", 10, 1)
			require.NoError(t, err, "container did not start")
			return nil
		},
	)
}

func testAccResourceContainerCreateCloneConfig() string {
	return fmt.Sprintf(`
resource "proxmox_virtual_environment_container" "test_container_clone" {
  depends_on = [proxmox_virtual_environment_container.test_container]

  node_name = "%s"
  vm_id     = 1101

  clone {
	vm_id = 1100
  }

  initialization {
    hostname = "test-clone"
  }
}
`, accTestNodeName)
}

func testAccResourceContainerCreateCloneCheck(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return resource.ComposeTestCheckFunc(
		func(*terraform.State) error {
			err := getNodesClient().Container(1101).WaitForContainerStatus(context.Background(), "running", 10, 1)
			require.NoError(t, err, "container did not start")
			return nil
		},
	)
}
