/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
)

const (
	accTestContainerName      = "proxmox_virtual_environment_container.test_container"
	accTestContainerCloneName = "proxmox_virtual_environment_container.test_container_clone"
)

//nolint:gochecknoglobals
var (
	accTestContainerID  = 100000 + rand.Intn(99999) //nolint:gosec
	accCloneContainerID = 200000 + rand.Intn(99999) //nolint:gosec
)

func TestAccResourceContainer(t *testing.T) { //nolint:wsl
	// download fails with 404 or "exit code 8" if run in parallel
	// t.Parallel()

	te := initTestEnvironment(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.accProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceContainerCreateConfig(te, false),
				Check:  testAccResourceContainerCreateCheck(te),
			},
			{
				Config: testAccResourceContainerCreateConfig(te, true) + testAccResourceContainerCreateCloneConfig(te),
				Check:  testAccResourceContainerCreateCloneCheck(te),
			},
		},
	})
}

func testAccResourceContainerCreateConfig(te *testEnvironment, isTemplate bool) string {
	te.t.Helper()

	return fmt.Sprintf(`
resource "proxmox_virtual_environment_download_file" "ubuntu_container_template" {
	content_type = "vztmpl"
	datastore_id = "local"
	node_name    = "%[1]s"
	url = "http://download.proxmox.com/images/system/ubuntu-23.04-standard_23.04-1_amd64.tar.zst"
    overwrite_unmanaged = true
}
resource "proxmox_virtual_environment_container" "test_container" {
  node_name = "%[1]s"
  vm_id     = %d
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
`, te.nodeName, accTestContainerID, isTemplate)
}

func testAccResourceContainerCreateCheck(te *testEnvironment) resource.TestCheckFunc {
	te.t.Helper()

	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(accTestContainerName, "description", "my\ndescription\nvalue\n"),
		func(*terraform.State) error {
			err := te.nodeClient().Container(accTestContainerID).WaitForContainerStatus(context.Background(), "running", 10, 1)
			require.NoError(te.t, err, "container did not start")

			return nil
		},
	)
}

func testAccResourceContainerCreateCloneConfig(te *testEnvironment) string {
	te.t.Helper()

	return fmt.Sprintf(`
resource "proxmox_virtual_environment_container" "test_container_clone" {
  depends_on = [proxmox_virtual_environment_container.test_container]

  node_name = "%s"
  vm_id     = %d

  clone {
	vm_id = proxmox_virtual_environment_container.test_container.id
  }

  initialization {
    hostname = "test-clone"
  }
}
`, te.nodeName, accCloneContainerID)
}

func testAccResourceContainerCreateCloneCheck(te *testEnvironment) resource.TestCheckFunc {
	te.t.Helper()

	return resource.ComposeTestCheckFunc(
		func(*terraform.State) error {
			err := te.nodeClient().Container(accCloneContainerID).WaitForContainerStatus(context.Background(), "running", 10, 1)
			require.NoError(te.t, err, "container did not start")

			return nil
		},
	)
}
