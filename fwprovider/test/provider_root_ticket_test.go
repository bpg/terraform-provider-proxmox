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

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

// TestAccProviderRootTicket verifies that `ssh { root_ticket = true }` lets an API token perform operations PVE allows
// only for a logged-in root@pam user, through both the SDK and the Framework provider. The first step of each case is
// the control: without the setting the same config must fail.
func TestAccProviderRootTicket(t *testing.T) {
	te := InitEnvironment(t)

	tests := []struct {
		name        string
		address     string
		vmIDAttr    string
		config      string
		expectError *regexp.Regexp
		check       func(vm *vms.GetResponseData) error
	}{
		{
			name:     "sdk vm serial device",
			address:  "proxmox_virtual_environment_vm.test_vm",
			vmIDAttr: "vm_id",
			config: `
				resource "proxmox_virtual_environment_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name      = "test-root-ticket-sdk"
					started   = false

					serial_device {
						device = "/dev/ttyS0"
					}
				}`,
			expectError: regexp.MustCompile(`only root can\s+set\s+.serial0.\s+config\s+for\s+real\s+devices`),
			check: func(vm *vms.GetResponseData) error {
				if vm.SerialDevice0 == nil || *vm.SerialDevice0 != "/dev/ttyS0" {
					return fmt.Errorf("expected serial0 to be /dev/ttyS0, got %v", vm.SerialDevice0)
				}

				return nil
			},
		},
		{
			name:     "framework vm cpu architecture",
			address:  "proxmox_vm.test_vm",
			vmIDAttr: "id",
			config: `
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name      = "test-root-ticket-fwk"
					cpu = {
						architecture = "x86_64"
					}
				}`,
			expectError: regexp.MustCompile(`only root can\s+set\s+.arch.\s+config`),
			check: func(vm *vms.GetResponseData) error {
				if vm.CPUArchitecture == nil || *vm.CPUArchitecture != "x86_64" {
					return fmt.Errorf("expected arch to be x86_64, got %v", vm.CPUArchitecture)
				}

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps: []resource.TestStep{
					{
						Config:      te.RenderConfig(tt.config, WithAPIToken()),
						ExpectError: tt.expectError,
					},
					{
						Config: te.RenderConfig(tt.config, WithAPIToken(), WithRootTicket()),
						// read from the API, not state, which only echoes the plan
						Check: func(s *terraform.State) error {
							rs, ok := s.RootModule().Resources[tt.address]
							if !ok {
								return fmt.Errorf("resource %s not found in state", tt.address)
							}

							vmID, err := strconv.Atoi(rs.Primary.Attributes[tt.vmIDAttr])
							require.NoError(t, err)

							vm, err := te.NodeClient().VM(vmID).GetVM(t.Context())
							require.NoError(t, err, "failed to get VM")

							return tt.check(vm)
						},
					},
				},
			})
		})
	}
}
