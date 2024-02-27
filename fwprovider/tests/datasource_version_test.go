/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceVersion(t *testing.T) {
	t.Parallel()

	accProviders := testAccMuxProviders(context.Background(), t)

	datasourceName := "data.proxmox_virtual_environment_version.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: accProviders,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: `data "proxmox_virtual_environment_version" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "release", "8.1"),
					resource.TestCheckResourceAttrSet(datasourceName, "repository_id"),
					resource.TestCheckResourceAttrWith(datasourceName, "version", func(value string) error {
						if strings.HasPrefix(value, "8.1") {
							return nil
						}
						return fmt.Errorf("version %s does not start with 8.1", value)
					}),
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
				),
			},
		},
	})
}
