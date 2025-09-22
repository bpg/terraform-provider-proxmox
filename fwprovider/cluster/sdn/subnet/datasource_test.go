//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package subnet_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccDataSourceSDNSubnet(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_simple" "datasource_subnet_zone" {
					id    = "dsrcz2"
					nodes = ["{{.NodeName}}"]
				}

				resource "proxmox_virtual_environment_sdn_vnet" "datasource_subnet_vnet" {
					id     = "dsrcv"
					zone   = proxmox_virtual_environment_sdn_zone_simple.datasource_subnet_zone.id
				}

				resource "proxmox_virtual_environment_sdn_subnet" "datasource_subnet" {
					cidr            = "10.60.0.0/24"
					vnet            = proxmox_virtual_environment_sdn_vnet.datasource_subnet_vnet.id
					gateway         = "10.60.0.1"
					dhcp_dns_server = "10.60.0.53"
					snat            = true
				}

				data "proxmox_virtual_environment_sdn_subnet" "datasource_subnet" {
					cidr = "10.60.0.0/24"
					vnet = proxmox_virtual_environment_sdn_subnet.datasource_subnet.vnet
					depends_on = [proxmox_virtual_environment_sdn_subnet.datasource_subnet]
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("data.proxmox_virtual_environment_sdn_subnet.datasource_subnet", map[string]string{
						"cidr":            "10.60.0.0/24",
						"vnet":            "dsrcv",
						"gateway":         "10.60.0.1",
						"dhcp_dns_server": "10.60.0.53",
						"snat":            "true",
					}),
				),
			},
		},
	})
}
