/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestClusterOptionsResource(t *testing.T) {
	t.Parallel()

	accProviders := testAccMuxProviders(context.Background(), t)

	resourceName := "proxmox_virtual_environment_cluster_options.test_options"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: accProviders,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `
resource "proxmox_virtual_environment_cluster_options" "test_options" {
	language                  = "en"
	keyboard                  = "pl"
	email_from                = "example@example.com"
	bandwidth_limit_migration = 555554
	bandwidth_limit_default   = 666666
	max_workers               = 5
	crs_ha                    = "static"
	ha_shutdown_policy        = "freeze"
	migration_cidr            = "10.0.0.0/8"
	migration_type            = "secure"
}
 `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "language", "en"),
					resource.TestCheckResourceAttr(resourceName, "keyboard", "pl"),
					resource.TestCheckResourceAttr(resourceName, "email_from", "example@example.com"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_limit_migration", "555554"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_limit_default", "666666"),
					resource.TestCheckResourceAttr(resourceName, "max_workers", "5"),
					resource.TestCheckResourceAttr(resourceName, "crs_ha", "static"),
					resource.TestCheckResourceAttr(resourceName, "ha_shutdown_policy", "freeze"),
					resource.TestCheckResourceAttr(resourceName, "migration_cidr", "10.0.0.0/8"),
					resource.TestCheckResourceAttr(resourceName, "migration_type", "secure"),
					resource.TestCheckResourceAttr(resourceName, "id", "cluster"),
					resource.TestCheckNoResourceAttr(resourceName, "bandwidth_limit_restore"),
					resource.TestCheckNoResourceAttr(resourceName, "bandwidth_limit_move"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update testing
			{
				Config: `
resource "proxmox_virtual_environment_cluster_options" "test_options" {
	language                  = "en"
	keyboard                  = "pl"
	email_from                = "ged@gont.earthsea"
	bandwidth_limit_migration = 111111
	bandwidth_limit_default   = 666666
	max_workers               = 6
	migration_cidr            = "10.0.0.0/8"
	migration_type            = "secure"
}
 `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "language", "en"),
					resource.TestCheckResourceAttr(resourceName, "keyboard", "pl"),
					resource.TestCheckResourceAttr(resourceName, "email_from", "ged@gont.earthsea"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_limit_migration", "111111"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_limit_default", "666666"),
					resource.TestCheckResourceAttr(resourceName, "max_workers", "6"),
					resource.TestCheckResourceAttr(resourceName, "migration_cidr", "10.0.0.0/8"),
					resource.TestCheckResourceAttr(resourceName, "migration_type", "secure"),
					resource.TestCheckResourceAttr(resourceName, "id", "cluster"),
					resource.TestCheckNoResourceAttr(resourceName, "bandwidth_limit_restore"),
					resource.TestCheckNoResourceAttr(resourceName, "bandwidth_limit_move"),
					resource.TestCheckNoResourceAttr(resourceName, "crs_ha"),
					resource.TestCheckNoResourceAttr(resourceName, "ha_shutdown_policy"),
				),
			},
		},
	})
}
