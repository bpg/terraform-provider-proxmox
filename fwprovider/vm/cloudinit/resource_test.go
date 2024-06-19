/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cloudinit_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

const resourceName = "proxmox_virtual_environment_vm2.test_vm"

func TestAccResourceVM2CloudInit(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create VM with cloud-init", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				id = {{.RandomVMID}}
				name = "test-cloudinit"
				initialization = {
					dns = {
						domain = "example.com"
			        }
				}
			}`),
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}
