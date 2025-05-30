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

func TestAccDatasourceSDNZone(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{
			"read sdn zone attributes",
			[]resource.TestStep{{
				Config: te.RenderConfig(`
					data "proxmox_virtual_environment_sdn_zone" "zone_ex" {
						name = "{{ .ZoneName }}"
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributesSet("data.proxmox_virtual_environment_sdn_zone.zone_ex", []string{
						"id",
						"name",
						"type",
						"ipam",
					}),
				),
			}},
		},
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
