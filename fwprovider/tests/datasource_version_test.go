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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_VersionDatasource(t *testing.T) {
	t.Parallel()

	accProviders := AccMuxProviders(context.Background(), t)

	datasourceName := "data.proxmox_virtual_environment_version.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: accProviders,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: ProviderConfig + `data "proxmox_virtual_environment_version" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "release", "8.0"),
					resource.TestCheckResourceAttrSet(datasourceName, "repository_id"),
					resource.TestCheckResourceAttrWith(datasourceName, "version", func(value string) error {
						if strings.HasPrefix(value, "8.0") {
							return nil
						}
						return fmt.Errorf("version %s does not start with 8.0", value)
					}),
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
				),
			},
		},
	})
}
