//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datastores_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccDatasourceDatastores(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"read datastores attributes", []resource.TestStep{{
			Config: te.RenderConfig(`data "proxmox_virtual_environment_datastores" "test" {
				node_name = "{{.NodeName}}"
				filters = {
					content_types = ["iso"]
				}
			}`),

			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributesSet("data.proxmox_virtual_environment_datastores.test", []string{
					"node_name",
				}),
				test.ResourceAttributes("data.proxmox_virtual_environment_datastores.test", map[string]string{
					"datastores.#":        "1",
					"datastores.0.active": "true",
					"datastores.0.id":     "local",
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
