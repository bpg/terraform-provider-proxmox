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

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

const (
	accTestContainerName      = "proxmox_virtual_environment_container.test_container"
	accTestContainerCloneName = "proxmox_virtual_environment_container.test_container_clone"
)

func TestAccResourceContainer(t *testing.T) {
	t.Parallel()

	accProviders := testAccMuxProviders(context.Background(), t)
	client := newClient(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: accProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceContainerCreateConfig(false),
				Check:  testAccResourceContainerCreateCheck(t, client),
			},
			{
				Config: testAccResourceContainerCreateConfig(true) + testAccResourceContainerCreateCloneConfig(),
				Check:  testAccResourceContainerCreateCloneCheck(t, client),
			},
		},
	})
}

func testAccResourceContainerCreateConfig(isTemplate bool) string {
	return fmt.Sprintf(`
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
    template_file_id = "local:vztmpl/ubuntu-23.04-standard_23.04-1_amd64.tar.zst"
    type             = "ubuntu"
  }
}
`, accTestNodeName, isTemplate)
}

func testAccResourceContainerCreateCheck(t *testing.T, client *nodes.Client) resource.TestCheckFunc {
	t.Helper()

	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(accTestContainerName, "description", "my\ndescription\nvalue\n"),
		func(*terraform.State) error {
			err := client.Container(1100).WaitForContainerStatus(context.Background(), "running", 10, 1)
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

func testAccResourceContainerCreateCloneCheck(t *testing.T, client *nodes.Client) resource.TestCheckFunc {
	t.Helper()

	return resource.ComposeTestCheckFunc(
		func(*terraform.State) error {
			err := client.Container(1101).WaitForContainerStatus(context.Background(), "running", 10, 1)
			require.NoError(t, err, "container did not start")
			return nil
		},
	)
}

func newClient(t *testing.T) *nodes.Client {
	t.Helper()

	username := utils.GetAnyStringEnv("PROXMOX_VE_USERNAME")
	password := utils.GetAnyStringEnv("PROXMOX_VE_PASSWORD")
	endpoint := utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT")
	apiToken := utils.GetAnyStringEnv("PROXMOX_VE_API_TOKEN")

	creds, err := api.NewCredentials(username, password, "", apiToken)
	require.NoError(t, err)

	conn, err := api.NewConnection(endpoint, true)
	require.NoError(t, err)

	client, err := api.NewClient(creds, conn)
	require.NoError(t, err)

	return &nodes.Client{Client: client, NodeName: accTestNodeName}
}
