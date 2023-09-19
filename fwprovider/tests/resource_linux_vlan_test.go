/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	accTestLinuxVLANName = "proxmox_virtual_environment_network_linux_vlan.test"
)

func TestAccResourceLinuxVLAN(t *testing.T) {
	t.Parallel()

	accProviders := testAccMuxProviders(context.Background(), t)

	iface := "eno0"
	vlan1 := gofakeit.Number(10, 4094)
	customName := fmt.Sprintf("iface_%s", gofakeit.Word())
	vlan2 := gofakeit.Number(10, 4094)
	ipV4cidr := fmt.Sprintf("%s/24", gofakeit.IPv4Address())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: accProviders,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccResourceLinuxVLANCreatedConfig(iface, vlan1),
				Check:  testAccResourceLinuxVLANCreatedCheck(iface, vlan1),
			},
			// ImportState testing
			{
				ResourceName:      accTestLinuxVLANName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Create and Read with a custom name
			{
				Config: testAccResourceLinuxVLANCustomNameCreatedConfig(customName, iface, vlan2),
				Check:  testAccResourceLinuxVLANCustomNameCreatedCheck(customName, iface, vlan2),
				// PVE API is unreliable. Sometimes it returns a wrong VLAN ID for this second interface.
				SkipFunc: func() (bool, error) {
					return true, nil
				},
			},
			// Update testing
			{
				Config: testAccResourceLinuxVLANUpdatedConfig(iface, vlan1, ipV4cidr),
				Check:  testAccResourceLinuxVLANUpdatedCheck(iface, vlan1, ipV4cidr),
			},
		},
	})
}

func testAccResourceLinuxVLANCreatedConfig(iface string, vlan int) string {
	return fmt.Sprintf(`
	resource "proxmox_virtual_environment_network_linux_vlan" "test" {
		comment = "created by terraform"
		mtu = 1499
		name = "%s.%d"
		node_name = "%s"
	}
	`, iface, vlan, accTestNodeName)
}

func testAccResourceLinuxVLANCreatedCheck(iface string, vlan int) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(accTestLinuxVLANName, "comment", "created by terraform"),
		resource.TestCheckResourceAttr(accTestLinuxVLANName, "interface", iface),
		resource.TestCheckResourceAttr(accTestLinuxVLANName, "name", fmt.Sprintf("%s.%d", iface, vlan)),
		resource.TestCheckResourceAttr(accTestLinuxVLANName, "vlan", strconv.Itoa(vlan)),
		resource.TestCheckResourceAttrSet(accTestLinuxVLANName, "id"),
	)
}

func testAccResourceLinuxVLANCustomNameCreatedConfig(name string, iface string, vlan int) string {
	return fmt.Sprintf(`
	resource "proxmox_virtual_environment_network_linux_vlan" "%s" {
		comment = "created by terraform"
		interface = "%s"
		mtu = 1499
		name = "%s"
		node_name = "%s"
		vlan = %d
	}
	`, name, iface, name, accTestNodeName, vlan)
}

func testAccResourceLinuxVLANCustomNameCreatedCheck(name string, iface string, vlan int) resource.TestCheckFunc {
	resourceName := fmt.Sprintf("proxmox_virtual_environment_network_linux_vlan.%s", name)

	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "comment", "created by terraform"),
		resource.TestCheckResourceAttr(resourceName, "interface", iface),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "vlan", strconv.Itoa(vlan)),
		resource.TestCheckResourceAttrSet(resourceName, "id"),
	)
}

func testAccResourceLinuxVLANUpdatedConfig(iface string, vlan int, ipV4cidr string) string {
	return fmt.Sprintf(`
	resource "proxmox_virtual_environment_network_linux_vlan" "test" {
		address = "%s"
		address6 = "FE80:0000:0000:0000:0202:B3FF:FE1E:8329/64"
		comment = "updated by terraform"
		name = "%s.%d"
		node_name = "%s"
	}
	`, ipV4cidr, iface, vlan, accTestNodeName)
}

func testAccResourceLinuxVLANUpdatedCheck(iface string, vlan int, ipV4cidr string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(accTestLinuxVLANName, "address", ipV4cidr),
		resource.TestCheckResourceAttr(accTestLinuxVLANName, "address6", "FE80:0000:0000:0000:0202:B3FF:FE1E:8329/64"),
		resource.TestCheckResourceAttr(accTestLinuxVLANName, "comment", "updated by terraform"),
		resource.TestCheckResourceAttr(accTestLinuxVLANName, "interface", iface),
		resource.TestCheckResourceAttr(accTestLinuxVLANName, "name", fmt.Sprintf("%s.%d", iface, vlan)),
		resource.TestCheckResourceAttr(accTestLinuxVLANName, "vlan", strconv.Itoa(vlan)),
		resource.TestCheckNoResourceAttr(accTestLinuxVLANName, "mtu"),
		resource.TestCheckResourceAttrSet(accTestLinuxVLANName, "id"),
	)
}
