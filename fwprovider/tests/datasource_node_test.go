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
)

func TestAccDatasourceNode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"read node attributes", []resource.TestStep{{
			Config: fmt.Sprintf(`data "proxmox_virtual_environment_node" "test" { node_name = "%s" }`, accTestNodeName),
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributesSet("data.proxmox_virtual_environment_node.test", []string{
					"cpu_count",
					"cpu_sockets",
					"cpu_model",
					"memory_available",
					"memory_used",
					"memory_total",
					"uptime",
				}),
			),
		}}},
	}

	accProviders := testAccMuxProviders(context.Background(), t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: accProviders,
				Steps:                    tt.steps,
			})
		})
	}
}
