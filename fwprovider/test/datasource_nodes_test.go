//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceNodes(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"read nodes attributes", []resource.TestStep{{
			Config: `data "proxmox_virtual_environment_nodes" "test" {}`,
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributesSet("data.proxmox_virtual_environment_nodes.test", []string{
					"cpu_count.0",
					"cpu_utilization.0",
					"memory_available.0",
					"memory_used.0",
					"names.0",
					"online.0",
					"uptime.0",
				}),
			),
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
