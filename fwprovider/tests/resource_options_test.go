/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider"
)

const accTestClusterOptionsName = "proxmox_virtual_environment_cluster_options.test_options"

func TestAccResourceClusterOptions(t *testing.T) {
	t.Parallel()

	te := initTestEnvironment(t)

	resource.Test(
		t, resource.TestCase{
			ProtoV6ProviderFactories: te.accProviders,
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
		},
	)
}

func testAccResourceClusterOptionsCreatedConfig() string {
	return fmt.Sprintf(
		`
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
		next_id = {
		  lower = %d
			upper = %d
		}
		notify = {
      ha_fencing_mode            = "never"
      ha_fencing_target          = "default-matcher"
      package_updates            = "always"
      package_updates_target     = "default-matcher"
      replication        = "always"
      replication_target = "default-matcher"
    }
	}
	`,
		fwprovider.ClusterOptionsNextIDLowerMinimum,
		fwprovider.ClusterOptionsNextIDLowerMaximum,
	)
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
		resource.TestCheckResourceAttr(
			accTestClusterOptionsName,
			"next_id.lower",
			fmt.Sprintf("%d", fwprovider.ClusterOptionsNextIDLowerMinimum),
		),
		resource.TestCheckResourceAttr(
			accTestClusterOptionsName,
			"next_id.upper",
			fmt.Sprintf("%d", fwprovider.ClusterOptionsNextIDLowerMaximum),
		),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "notify.ha_fencing_mode", "never"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "notify.ha_fencing_target", "default-matcher"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "notify.package_updates", "always"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "notify.package_updates_target", "default-matcher"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "notify.replication", "always"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "notify.replication_target", "default-matcher"),
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
		next_id = {
		  lower = 555
			upper = 666
		}
    notify = {
      ha_fencing_mode        = "always"
      ha_fencing_target      = "custom-matcher"
      package_updates        = "auto"
      package_updates_target = "custom-matcher"
      replication            = "never"
      replication_target     = "custom-matcher"
    }
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
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "next_id.lower", "555"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "next_id.upper", "666"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "notify.ha_fencing_mode", "always"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "notify.ha_fencing_target", "custom-matcher"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "notify.package_updates", "auto"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "notify.package_updates_target", "custom-matcher"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "notify.replication", "never"),
		resource.TestCheckResourceAttr(accTestClusterOptionsName, "notify.replication_target", "custom-matcher"),
		resource.TestCheckNoResourceAttr(accTestClusterOptionsName, "bandwidth_limit_move"),
		resource.TestCheckNoResourceAttr(accTestClusterOptionsName, "crs_ha"),
		resource.TestCheckNoResourceAttr(accTestClusterOptionsName, "ha_shutdown_policy"),
		resource.TestCheckNoResourceAttr(accTestClusterOptionsName, "http_proxy"),
		resource.TestCheckNoResourceAttr(accTestClusterOptionsName, "keyboard"),
	)
}
