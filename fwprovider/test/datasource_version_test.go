//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceVersion(t *testing.T) {
	te := InitEnvironment(t)

	datasourceName := "data.proxmox_virtual_environment_version.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: `data "proxmox_virtual_environment_version" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "release", "8.4"),
					resource.TestCheckResourceAttrSet(datasourceName, "repository_id"),
					resource.TestCheckResourceAttrWith(datasourceName, "version", func(value string) error {
						if strings.HasPrefix(value, "8.4") {
							return nil
						}
						return fmt.Errorf("version %s does not start with 8.4", value)
					}),
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
				),
			},
		},
	})
}
