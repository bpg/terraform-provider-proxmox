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

const accTestClusterOptionsName = "proxmox_virtual_environment_cluster_options.test_options"

func TestAccResourceClusterOptions(t *testing.T) {
	t.Parallel()

	accProviders := testAccMuxProviders(context.Background(), t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: accProviders,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccResourceClusterOptionsCreatedConfig(),
				Check:  testAccResourceClusterOptionsCreatedCheck(),
			},
			// ImportState testing
			{
				ResourceName:      accTestClusterOptionsName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update testing
			{
				Config: testAccResourceClusterOptionsUpdatedConfig(),
				Check:  testAccResourceClusterOptionsUpdatedCheck(),
			},
		},
	})
}

func testAccResourceClusterOptionsCreatedConfig() string {
	return `
	resource "proxmox_virtual_environment_cluster_options" "test_options" {
		bandwidth_limit_default   = 666666
		bandwidth_limit_migration = 555554
		crs_ha                    = "static"
		email_from                = "example@example.com"
		ha_shutdown_policy        = "freeze"
		http_proxy                = "http://example.com"
		keyboard                  = "pl"
		language                  = "en"
		max_workers               = 5
		migration_cidr            = "10.0.0.0/8"
		migration_type            = "secure"
        bandwidth_limit_restore   = 777777
	}
	`
}

func testAccResourceClusterOptionsCreatedCheck() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "bandwidth_limit_default", "666666"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "bandwidth_limit_migration", "555554"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "bandwidth_limit_restore", "777777"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "crs_ha", "static"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "email_from", "example@example.com"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "ha_shutdown_policy", "freeze"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "http_proxy", "http://example.com"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "id", "cluster"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "keyboard", "pl"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "language", "en"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "max_workers", "5"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "migration_cidr", "10.0.0.0/8"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "migration_type", "secure"),
		resource.TestCheckNoResourceAttr(accTestClusterOptionsName, "bandwidth_limit_move"),
	)
}

func testAccResourceClusterOptionsUpdatedConfig() string {
	return `
	resource "proxmox_virtual_environment_cluster_options" "test_options" {
		bandwidth_limit_default   = 333333
		bandwidth_limit_migration = 111111
		email_from                = "ged@gont.earthsea"
		language                  = "en"
		max_workers               = 6
		migration_cidr            = "10.0.0.1/8"
		migration_type            = "secure"
	}
	`
}

func testAccResourceClusterOptionsUpdatedCheck() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "bandwidth_limit_default", "333333"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "bandwidth_limit_migration", "111111"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "email_from", "ged@gont.earthsea"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "id", "cluster"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "language", "en"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "max_workers", "6"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "migration_cidr", "10.0.0.1/8"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "migration_type", "secure"),
		resource.TestCheckNoResourceAttr(accTestClusterOptionsName, "bandwidth_limit_move"),
		resource.TestCheckNoResourceAttr(accTestClusterOptionsName, "crs_ha"),
		resource.TestCheckNoResourceAttr(accTestClusterOptionsName, "ha_shutdown_policy"),
		resource.TestCheckNoResourceAttr(accTestClusterOptionsName, "http_proxy"),
		resource.TestCheckNoResourceAttr(accTestClusterOptionsName, "keyboard"),
	)
}
