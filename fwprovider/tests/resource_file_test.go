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

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	accTestFileName = "proxmox_virtual_environment_file.test"
)

func TestAccResourceFile(t *testing.T) {
	t.Parallel()

	accProviders := testAccMuxProviders(context.Background(), t)

	snippet := fmt.Sprintf("snippet-%s.txt", gofakeit.Word())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: accProviders,
		Steps: []resource.TestStep{
			// Upload a snippet file from a raw source
			{
				Config: testAccResourceFileCreatedConfig(snippet),
				Check:  testAccResourceFileCreatedCheck(snippet),
			},
			// ImportState testing
			// {
			// 	ResourceName:        accTestFileName,
			// 	ImportState:         true,
			// 	ImportStateVerify:   true,
			// 	ImportStateIdPrefix: "local:snippets/",
			// 	ImportStateId:       fmt.Sprintf("local:snippets/%s", snippet),
			// },
			// // Update testing
			// {
			// 	Config: testAccResourceLinuxVLANUpdatedConfig(iface, vlan1, ipV4cidr),
			// 	Check:  testAccResourceLinuxVLANUpdatedCheck(iface, vlan1, ipV4cidr),
			// },
		},
	})
}

func testAccResourceFileCreatedConfig(fname string) string {
	return fmt.Sprintf(`
resource "proxmox_virtual_environment_file" "test" {
  content_type = "snippets"
  datastore_id = "local"
  node_name    = "%s"

  source_raw {
    data = <<EOF
test snippet
    EOF

    file_name = "%s"
  }
}
	`, accTestNodeName, fname)
}

func testAccResourceFileCreatedCheck(fname string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(accTestFileName, "content_type", "snippets"),
		// resource.TestCheckResourceAttr(accTestFileName, "file_name", fname),
		resource.TestCheckResourceAttr(accTestFileName, "source_raw.0.file_name", fname),
		resource.TestCheckResourceAttr(accTestFileName, "id", fmt.Sprintf("local:snippets/%s", fname)),
	)
}
