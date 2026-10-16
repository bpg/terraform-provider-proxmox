//go:build acceptance || all

//testacc:tier=medium
//testacc:resource=misc

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/utils"
)

// TestAccProviderRootTicket verifies that `ssh { root_ticket = true }` lets an API token set a
// serial device to a real host device, which PVE allows only for a logged-in root@pam user.
// The first step is the control: without the setting the same config must fail.
func TestAccProviderRootTicket(t *testing.T) {
	if utils.GetAnyStringEnv("PROXMOX_VE_SSH_USERNAME") != "root" {
		t.Skip("PROXMOX_VE_SSH_USERNAME is not 'root'")
	}

	te := InitEnvironment(t)

	const config = `
		resource "proxmox_virtual_environment_vm" "test_vm" {
			node_name = "{{.NodeName}}"
			name      = "test-root-ticket"
			started   = false

			serial_device {
				device = "/dev/ttyS0"
			}
		}`

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config:      te.RenderConfig(config, WithAPIToken()),
				ExpectError: regexp.MustCompile(`only root can\s+set .serial0. config for real devices`),
			},
			{
				Config: te.RenderConfig(config, WithAPIToken(), WithRootTicket()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm", map[string]string{
						"serial_device.0.device": "/dev/ttyS0",
					}),
					// read from the API, not state, which only echoes the plan
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["proxmox_virtual_environment_vm.test_vm"]
						if !ok {
							return fmt.Errorf("resource not found in state")
						}

						vmID, err := strconv.Atoi(rs.Primary.Attributes["vm_id"])
						require.NoError(te.t, err)

						vmConfig, err := te.NodeClient().VM(vmID).GetVM(t.Context())
						require.NoError(te.t, err, "failed to get VM")
						require.NotNil(te.t, vmConfig.SerialDevice0, "the VM has no serial0 device")
						require.Equal(te.t, "/dev/ttyS0", *vmConfig.SerialDevice0)

						return nil
					},
				),
			},
		},
	})
}
