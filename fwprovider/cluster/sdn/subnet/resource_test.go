//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package subnet_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceSDNSubnet(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"basic subnet create", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_sdn_zone_simple" "subnet_zone" {
				id    = "subnetz"
				nodes = ["{{.NodeName}}"]
			}

			resource "proxmox_virtual_environment_sdn_vnet" "subnet_vnet" {
				id     = "subnetv"
				zone     = proxmox_virtual_environment_sdn_zone_simple.subnet_zone.id
				depends_on = [
					proxmox_virtual_environment_sdn_applier.finalizer
				]
			}

			// use variables to test resource validation with interpolated values
			// they are unknown until the plan is applied, so the resource validation should be skipped.
			variable "cidr" {
				type = string
				default = "10.10.0.0/24"
			}

			variable "gateway" {
				type = string
				default = "10.10.0.1"
			}

			resource "proxmox_virtual_environment_sdn_subnet" "test_subnet" {
				cidr  = var.cidr
				vnet    = proxmox_virtual_environment_sdn_vnet.subnet_vnet.id
				gateway = var.gateway
				depends_on = [
					proxmox_virtual_environment_sdn_applier.finalizer
				]
			}

			resource "proxmox_virtual_environment_sdn_applier" "test_applier" {
				depends_on = [
					proxmox_virtual_environment_sdn_zone_simple.subnet_zone,
					proxmox_virtual_environment_sdn_vnet.subnet_vnet,
					proxmox_virtual_environment_sdn_subnet.test_subnet
				]
			}

			resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_subnet.test_subnet", map[string]string{
					"cidr":    "10.10.0.0/24",
					"vnet":    "subnetv",
					"gateway": "10.10.0.1",
				}),
				test.ResourceAttributesSet("proxmox_virtual_environment_sdn_subnet.test_subnet", []string{
					"id",
				}),
			),
		}}},
		{"subnet with dhcp configuration", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_sdn_zone_simple" "dhcp_zone" {
				id    = "dhcpz2"
				nodes = ["{{.NodeName}}"]
			}

			resource "proxmox_virtual_environment_sdn_vnet" "dhcp_vnet" {
				id     = "dhcpv"
				zone     = proxmox_virtual_environment_sdn_zone_simple.dhcp_zone.id
				depends_on = [
					proxmox_virtual_environment_sdn_applier.finalizer
				]
			}

			resource "proxmox_virtual_environment_sdn_subnet" "dhcp_subnet" {
				cidr            = "192.168.1.0/24"
				vnet              = proxmox_virtual_environment_sdn_vnet.dhcp_vnet.id
				gateway           = "192.168.1.1"
				dhcp_dns_server   = "192.168.1.53"
				snat              = true
				dhcp_range = {
					start_address = "192.168.1.10"
					end_address   = "192.168.1.100"
				}
				depends_on = [
					proxmox_virtual_environment_sdn_applier.finalizer
				]
			}

			resource "proxmox_virtual_environment_sdn_applier" "dhcp_applier" {
				depends_on = [
					proxmox_virtual_environment_sdn_zone_simple.dhcp_zone,
					proxmox_virtual_environment_sdn_vnet.dhcp_vnet,
					proxmox_virtual_environment_sdn_subnet.dhcp_subnet
				]
			}

			resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_subnet.dhcp_subnet", map[string]string{
					"cidr":                     "192.168.1.0/24",
					"vnet":                     "dhcpv",
					"gateway":                  "192.168.1.1",
					"dhcp_dns_server":          "192.168.1.53",
					"snat":                     "true",
					"dhcp_range.start_address": "192.168.1.10",
					"dhcp_range.end_address":   "192.168.1.100",
				}),
			),
		}}},
		{"subnet with dhcp range", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_sdn_zone_simple" "multi_dhcp_zone" {
				id    = "multidh2"
				nodes = ["{{.NodeName}}"]
			}

			resource "proxmox_virtual_environment_sdn_vnet" "multi_dhcp_vnet" {
				id     = "multidhv"
				zone   = proxmox_virtual_environment_sdn_zone_simple.multi_dhcp_zone.id
				depends_on = [
					proxmox_virtual_environment_sdn_applier.finalizer
				]
			}

			resource "proxmox_virtual_environment_sdn_subnet" "multi_dhcp_subnet" {
				cidr    = "172.16.0.0/24"
				vnet    = proxmox_virtual_environment_sdn_vnet.multi_dhcp_vnet.id
				gateway = "172.16.0.1"
				dhcp_range = {
					start_address = "172.16.0.10"
					end_address   = "172.16.0.50"
				}
				depends_on = [
					proxmox_virtual_environment_sdn_applier.finalizer
				]
			}

			resource "proxmox_virtual_environment_sdn_applier" "multi_dhcp_applier" {
				depends_on = [
					proxmox_virtual_environment_sdn_zone_simple.multi_dhcp_zone,
					proxmox_virtual_environment_sdn_vnet.multi_dhcp_vnet,
					proxmox_virtual_environment_sdn_subnet.multi_dhcp_subnet
				]
			}

			resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_subnet.multi_dhcp_subnet", map[string]string{
					"cidr":                     "172.16.0.0/24",
					"vnet":                     "multidhv",
					"gateway":                  "172.16.0.1",
					"dhcp_range.start_address": "172.16.0.10",
					"dhcp_range.end_address":   "172.16.0.50",
				}),
			),
		}}},
		{"subnet update", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_simple" "update_subnet_zone" {
					id    = "updatez2"
					nodes = ["{{.NodeName}}"]
				}

				resource "proxmox_virtual_environment_sdn_vnet" "update_subnet_vnet" {
					id     = "updatev"
					zone   = proxmox_virtual_environment_sdn_zone_simple.update_subnet_zone.id
					depends_on = [
						proxmox_virtual_environment_sdn_applier.finalizer
					]
				}

				resource "proxmox_virtual_environment_sdn_subnet" "update_subnet" {
					cidr    = "10.20.0.0/24"
					vnet    = proxmox_virtual_environment_sdn_vnet.update_subnet_vnet.id
					gateway = "10.20.0.1"
					snat    = false
					depends_on = [
						proxmox_virtual_environment_sdn_applier.finalizer
					]
				}

				resource "proxmox_virtual_environment_sdn_applier" "update_applier" {
					depends_on = [
						proxmox_virtual_environment_sdn_zone_simple.update_subnet_zone,
						proxmox_virtual_environment_sdn_vnet.update_subnet_vnet,
						proxmox_virtual_environment_sdn_subnet.update_subnet
					]
				}

				resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_sdn_subnet.update_subnet", map[string]string{
						"cidr":    "10.20.0.0/24",
						"vnet":    "updatev",
						"gateway": "10.20.0.1",
						"snat":    "false",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_simple" "update_subnet_zone" {
					id    = "updatez2"
					nodes = ["{{.NodeName}}"]
				}

				resource "proxmox_virtual_environment_sdn_vnet" "update_subnet_vnet" {
					id     = "updatev"
					zone   = proxmox_virtual_environment_sdn_zone_simple.update_subnet_zone.id
					depends_on = [
						proxmox_virtual_environment_sdn_applier.finalizer
					]
				}

				resource "proxmox_virtual_environment_sdn_subnet" "update_subnet" {
					cidr            = "10.20.0.0/24"
					vnet            = proxmox_virtual_environment_sdn_vnet.update_subnet_vnet.id
					gateway         = "10.20.0.1"
					snat            = true
					dhcp_dns_server = "10.20.0.53"
					depends_on = [
						proxmox_virtual_environment_sdn_applier.finalizer
					]
				}

				resource "proxmox_virtual_environment_sdn_applier" "update_applier" {
					depends_on = [
						proxmox_virtual_environment_sdn_zone_simple.update_subnet_zone,
						proxmox_virtual_environment_sdn_vnet.update_subnet_vnet,
						proxmox_virtual_environment_sdn_subnet.update_subnet
					]
				}

				resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_sdn_subnet.update_subnet", map[string]string{
						"cidr":            "10.20.0.0/24",
						"vnet":            "updatev",
						"gateway":         "10.20.0.1",
						"snat":            "true",
						"dhcp_dns_server": "10.20.0.53",
					}),
				),
			},
		}},
		{"minimal subnet create", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_sdn_zone_simple" "minimal_zone" {
				id    = "minimalz"
				nodes = ["{{.NodeName}}"]
			}

			resource "proxmox_virtual_environment_sdn_vnet" "minimal_vnet" {
				id     = "minimalv"
				zone   = proxmox_virtual_environment_sdn_zone_simple.minimal_zone.id
				depends_on = [
					proxmox_virtual_environment_sdn_applier.finalizer
				]
			}

			resource "proxmox_virtual_environment_sdn_subnet" "minimal_subnet" {
				cidr = "172.20.0.0/24"
				vnet = proxmox_virtual_environment_sdn_vnet.minimal_vnet.id
				depends_on = [
					proxmox_virtual_environment_sdn_applier.finalizer
				]
			}

			resource "proxmox_virtual_environment_sdn_applier" "minimal_applier" {
				depends_on = [
					proxmox_virtual_environment_sdn_zone_simple.minimal_zone,
					proxmox_virtual_environment_sdn_vnet.minimal_vnet,
					proxmox_virtual_environment_sdn_subnet.minimal_subnet
				]
			}

			resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_subnet.minimal_subnet", map[string]string{
					"cidr": "172.20.0.0/24",
					"vnet": "minimalv",
				}),
				test.ResourceAttributesSet("proxmox_virtual_environment_sdn_subnet.minimal_subnet", []string{
					"id",
				}),
			),
		}}},
		{"subnet with all attributes", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_sdn_zone_simple" "all_zone" {
				id    = "allzone"
				nodes = ["{{.NodeName}}"]
			}

			resource "proxmox_virtual_environment_sdn_vnet" "all_vnet" {
				id     = "allvnet"
				zone   = proxmox_virtual_environment_sdn_zone_simple.all_zone.id
				depends_on = [
					proxmox_virtual_environment_sdn_applier.finalizer
				]
			}

			resource "proxmox_virtual_environment_sdn_subnet" "all_subnet" {
				cidr              = "172.30.0.0/24"
				vnet              = proxmox_virtual_environment_sdn_vnet.all_vnet.id
				gateway           = "172.30.0.1"
				dhcp_dns_server   = "172.30.0.53"
				dns_zone_prefix   = "example.com"
				snat              = true
				dhcp_range = {
					start_address = "172.30.0.10"
					end_address   = "172.30.0.50"
				}
				depends_on = [
					proxmox_virtual_environment_sdn_applier.finalizer
				]
			}

			resource "proxmox_virtual_environment_sdn_applier" "all_applier" {
				depends_on = [
					proxmox_virtual_environment_sdn_zone_simple.all_zone,
					proxmox_virtual_environment_sdn_vnet.all_vnet,
					proxmox_virtual_environment_sdn_subnet.all_subnet
				]
			}

			resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_sdn_subnet.all_subnet", map[string]string{
					"cidr":                     "172.30.0.0/24",
					"vnet":                     "allvnet",
					"gateway":                  "172.30.0.1",
					"dhcp_dns_server":          "172.30.0.53",
					"dns_zone_prefix":          "example.com",
					"snat":                     "true",
					"dhcp_range.start_address": "172.30.0.10",
					"dhcp_range.end_address":   "172.30.0.50",
				}),
			),
		}}},
		{"subnet with dhcp range updates", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_simple" "dhcp_update_zone" {
					id    = "dhcpupz"
					nodes = ["{{.NodeName}}"]
				}

				resource "proxmox_virtual_environment_sdn_vnet" "dhcp_update_vnet" {
					id     = "dhcpupv"
					zone   = proxmox_virtual_environment_sdn_zone_simple.dhcp_update_zone.id
					depends_on = [
						proxmox_virtual_environment_sdn_applier.finalizer
					]
				}

				resource "proxmox_virtual_environment_sdn_subnet" "dhcp_update_subnet" {
					cidr    = "172.40.0.0/24"
					vnet    = proxmox_virtual_environment_sdn_vnet.dhcp_update_vnet.id
					gateway = "172.40.0.1"
					depends_on = [
						proxmox_virtual_environment_sdn_applier.finalizer
					]
				}

				resource "proxmox_virtual_environment_sdn_applier" "dhcp_update_applier" {
					depends_on = [
						proxmox_virtual_environment_sdn_zone_simple.dhcp_update_zone,
						proxmox_virtual_environment_sdn_vnet.dhcp_update_vnet,
						proxmox_virtual_environment_sdn_subnet.dhcp_update_subnet
					]
				}

				resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_sdn_subnet.dhcp_update_subnet", map[string]string{
						"cidr":    "172.40.0.0/24",
						"vnet":    "dhcpupv",
						"gateway": "172.40.0.1",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_simple" "dhcp_update_zone" {
					id    = "dhcpupz"
					nodes = ["{{.NodeName}}"]
				}

				resource "proxmox_virtual_environment_sdn_vnet" "dhcp_update_vnet" {
					id     = "dhcpupv"
					zone   = proxmox_virtual_environment_sdn_zone_simple.dhcp_update_zone.id
					depends_on = [
						proxmox_virtual_environment_sdn_applier.finalizer
					]
				}

				resource "proxmox_virtual_environment_sdn_subnet" "dhcp_update_subnet" {
					cidr    = "172.40.0.0/24"
					vnet    = proxmox_virtual_environment_sdn_vnet.dhcp_update_vnet.id
					gateway = "172.40.0.1"
					dhcp_range = {
						start_address = "172.40.0.10"
						end_address   = "172.40.0.50"
					}
					depends_on = [
						proxmox_virtual_environment_sdn_applier.finalizer
					]
				}

				resource "proxmox_virtual_environment_sdn_applier" "dhcp_update_applier" {
					depends_on = [
						proxmox_virtual_environment_sdn_zone_simple.dhcp_update_zone,
						proxmox_virtual_environment_sdn_vnet.dhcp_update_vnet,
						proxmox_virtual_environment_sdn_subnet.dhcp_update_subnet
					]
				}

				resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_sdn_subnet.dhcp_update_subnet", map[string]string{
						"cidr":                     "172.40.0.0/24",
						"vnet":                     "dhcpupv",
						"gateway":                  "172.40.0.1",
						"dhcp_range.start_address": "172.40.0.10",
						"dhcp_range.end_address":   "172.40.0.50",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_simple" "dhcp_update_zone" {
					id    = "dhcpupz"
					nodes = ["{{.NodeName}}"]
				}

				resource "proxmox_virtual_environment_sdn_vnet" "dhcp_update_vnet" {
					id     = "dhcpupv"
					zone   = proxmox_virtual_environment_sdn_zone_simple.dhcp_update_zone.id
					depends_on = [
						proxmox_virtual_environment_sdn_applier.finalizer
					]
				}

				resource "proxmox_virtual_environment_sdn_subnet" "dhcp_update_subnet" {
					cidr    = "172.40.0.0/24"
					vnet    = proxmox_virtual_environment_sdn_vnet.dhcp_update_vnet.id
					gateway = "172.40.0.1"
					depends_on = [
						proxmox_virtual_environment_sdn_applier.finalizer
					]
				}

				resource "proxmox_virtual_environment_sdn_applier" "dhcp_update_applier" {
					depends_on = [
						proxmox_virtual_environment_sdn_zone_simple.dhcp_update_zone,
						proxmox_virtual_environment_sdn_vnet.dhcp_update_vnet,
						proxmox_virtual_environment_sdn_subnet.dhcp_update_subnet
					]
				}

				resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_sdn_subnet.dhcp_update_subnet", map[string]string{
						"cidr":    "172.40.0.0/24",
						"vnet":    "dhcpupv",
						"gateway": "172.40.0.1",
					}),
				),
			},
		}},
		{"subnet import", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_sdn_zone_simple" "import_subnet_zone" {
					id    = "importz2"
					nodes = ["{{.NodeName}}"]
				}

				resource "proxmox_virtual_environment_sdn_vnet" "import_subnet_vnet" {
					id     = "importv"
					zone   = proxmox_virtual_environment_sdn_zone_simple.import_subnet_zone.id
					depends_on = [
						proxmox_virtual_environment_sdn_applier.finalizer
					]
				}

				resource "proxmox_virtual_environment_sdn_subnet" "import_subnet" {
					cidr    = "10.30.0.0/24"
					vnet    = proxmox_virtual_environment_sdn_vnet.import_subnet_vnet.id
					gateway = "10.30.0.1"
					depends_on = [
						proxmox_virtual_environment_sdn_applier.finalizer
					]
				}

				resource "proxmox_virtual_environment_sdn_applier" "import_applier" {
					depends_on = [
						proxmox_virtual_environment_sdn_zone_simple.import_subnet_zone,
						proxmox_virtual_environment_sdn_vnet.import_subnet_vnet,
						proxmox_virtual_environment_sdn_subnet.import_subnet
					]
				}

				resource "proxmox_virtual_environment_sdn_applier" "finalizer" {}`),
			},
			{
				ResourceName:      "proxmox_virtual_environment_sdn_subnet.import_subnet",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "importv/importz2-10.30.0.0-24",
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}

func TestAccResourceSDNSubnetValidation(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name   string
		config string
		error  string
	}{
		{
			"invalid subnet cidr",
			`
			resource "proxmox_virtual_environment_sdn_zone_simple" "validation_zone" {
				id    = "validz3"
				nodes = ["{{.NodeName}}"]
			}

			resource "proxmox_virtual_environment_sdn_vnet" "validation_vnet" {
				id     = "validv3"
				zone   = proxmox_virtual_environment_sdn_zone_simple.validation_zone.id
			}

			resource "proxmox_virtual_environment_sdn_subnet" "validation_subnet" {
				cidr = "invalid-cidr"
				vnet = proxmox_virtual_environment_sdn_vnet.validation_vnet.id
			}`,
			"invalid CIDR address",
		},
		{
			"gateway outside subnet",
			`
			resource "proxmox_virtual_environment_sdn_zone_simple" "gateway_zone" {
				id    = "gatewz3"
				nodes = ["{{.NodeName}}"]
			}

			resource "proxmox_virtual_environment_sdn_vnet" "gateway_vnet" {
				id     = "gatewv3"
				zone   = proxmox_virtual_environment_sdn_zone_simple.gateway_zone.id
			}

			resource "proxmox_virtual_environment_sdn_subnet" "gateway_subnet" {
				cidr    = "10.40.0.0/24"
				vnet    = proxmox_virtual_environment_sdn_vnet.gateway_vnet.id
				gateway = "192.168.1.1"
			}`,
			"192.168.1.1 must be within the subnet 10.40.0.0/24",
		},
		{
			"dhcp range outside subnet",
			`
			resource "proxmox_virtual_environment_sdn_zone_simple" "dhcp_range_zone" {
				id    = "dhcprng3"
				nodes = ["{{.NodeName}}"]
			}

			resource "proxmox_virtual_environment_sdn_vnet" "dhcp_range_vnet" {
				id     = "dhcprng3"
				zone   = proxmox_virtual_environment_sdn_zone_simple.dhcp_range_zone.id
			}

			resource "proxmox_virtual_environment_sdn_subnet" "dhcp_range_subnet" {
				cidr = "10.50.0.0/24"
				vnet = proxmox_virtual_environment_sdn_vnet.dhcp_range_vnet.id
				dhcp_range = {
					start_address = "192.168.1.10"
					end_address   = "192.168.1.20"
				}
			}`,
			"End address 192.168.1.20 must be within the subnet 10.50.0.0/24",
		},
		{
			"dhcp dns server outside subnet",
			`
			resource "proxmox_virtual_environment_sdn_zone_simple" "dhcp_dns_zone" {
				id    = "dhcpdnsz"
				nodes = ["{{.NodeName}}"]
			}

			resource "proxmox_virtual_environment_sdn_vnet" "dhcp_dns_vnet" {
				id     = "dhcpdnsv"
				zone   = proxmox_virtual_environment_sdn_zone_simple.dhcp_dns_zone.id
			}

			resource "proxmox_virtual_environment_sdn_subnet" "dhcp_dns_subnet" {
				cidr            = "10.60.0.0/24"
				vnet            = proxmox_virtual_environment_sdn_vnet.dhcp_dns_vnet.id
				dhcp_dns_server = "192.168.1.53"
			}`,
			"192.168.1.53 must be within the subnet 10.60.0.0/24",
		},
		{
			"dhcp range start after end",
			`
			resource "proxmox_virtual_environment_sdn_zone_simple" "dhcp_order_zone" {
				id    = "dhcpordz"
				nodes = ["{{.NodeName}}"]
			}

			resource "proxmox_virtual_environment_sdn_vnet" "dhcp_order_vnet" {
				id     = "dhcpordv"
				zone   = proxmox_virtual_environment_sdn_zone_simple.dhcp_order_zone.id
			}

			resource "proxmox_virtual_environment_sdn_subnet" "dhcp_order_subnet" {
				cidr = "10.70.0.0/24"
				vnet = proxmox_virtual_environment_sdn_vnet.dhcp_order_vnet.id
				dhcp_range = {
					start_address = "10.70.0.50"
					end_address   = "10.70.0.10"
				}
			}`,
			"10.70.0.50 must be less than or equal to end address 10.70.0.10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps: []resource.TestStep{
					{
						Config:      te.RenderConfig(tt.config),
						ExpectError: regexp.MustCompile(tt.error),
					},
				},
			})
		})
	}
}
