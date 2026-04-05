//go:build acceptance || all

//testacc:tier=light
//testacc:resource=vm

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceSDKVMNotFound(t *testing.T) {
	te := InitEnvironment(t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					data "proxmox_virtual_environment_vm" "test" {
						node_name = "{{.NodeName}}"
						vm_id     = 999999
					}
				`),
				ExpectError: regexp.MustCompile(`(?i)not found|does not exist`),
			},
		},
	})
}
