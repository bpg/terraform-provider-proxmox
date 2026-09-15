//go:build acceptance || all

//testacc:tier=heavy
//testacc:resource=network

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package network_test

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

// TestAccResourceNetworkApplier covers the core contract: an interface staged with
// `reload = false` stays inactive until an applier reloads the node.
func TestAccResourceNetworkApplier(t *testing.T) {
	te := test.InitEnvironment(t)

	iface := fmt.Sprintf("vmbr%d", gofakeit.Number(10, 9999))
	ipV4cidr := fmt.Sprintf("%s/24", gofakeit.IPv4Address())

	bridge := fmt.Sprintf(`
	resource "proxmox_network_linux_bridge" "test" {
		node_name = "{{.NodeName}}"
		name      = "%s"
		address   = "%s"
		reload    = false
	}`, iface, ipV4cidr)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		// No finalizer: the dependent applier is destroyed before the bridge, so its delete stays staged.
		CheckDestroy: checkStagedDestroy(te),
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(bridge),
				Check: resource.ComposeTestCheckFunc(
					checkInterfaceActive(te, iface, false),
				),
			},
			{
				Config: te.RenderConfig(bridge + `
				resource "proxmox_network_applier" "test" {
					node_name  = "{{.NodeName}}"
					on_create  = false
					depends_on = [proxmox_network_linux_bridge.test]
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_network_applier.test", map[string]string{
						"on_create":  "false",
						"on_destroy": "true",
					}),
					test.ResourceAttributesSet("proxmox_network_applier.test", []string{"id"}),
					checkInterfaceActive(te, iface, false),
				),
			},
			{
				Config: te.RenderConfig(bridge + `
				resource "proxmox_network_applier" "test" {
					node_name  = "{{.NodeName}}"
					depends_on = [proxmox_network_linux_bridge.test]
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_network_applier.test", map[string]string{
						"on_create": "true",
					}),
					checkInterfaceActive(te, iface, true),
				),
			},
		},
	})
}

// TestAccResourceNetworkApplierTriggers covers the explicit re-apply escape hatch.
func TestAccResourceNetworkApplierTriggers(t *testing.T) {
	te := test.InitEnvironment(t)

	var firstID string

	config := func(v string) string {
		return fmt.Sprintf(`
		resource "proxmox_network_applier" "test" {
			node_name = "{{.NodeName}}"
			triggers = {
				version = "%s"
			}
		}`, v)
	}

	recordID := func(into *string) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			rs, ok := s.RootModule().Resources["proxmox_network_applier.test"]
			if !ok {
				return fmt.Errorf("resource not found in state")
			}

			*into = rs.Primary.Attributes["id"]

			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(config("one")),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_network_applier.test", map[string]string{
						"triggers.version": "one",
					}),
					recordID(&firstID),
				),
			},
			{
				Config: te.RenderConfig(config("two")),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_network_applier.test", map[string]string{
						"triggers.version": "two",
					}),
					func(s *terraform.State) error {
						rs := s.RootModule().Resources["proxmox_network_applier.test"]
						if rs.Primary.Attributes["id"] == firstID {
							return fmt.Errorf("id %q unchanged after trigger change; no new apply happened", firstID)
						}

						return nil
					},
				),
			},
		},
	})
}

// applierFinalizerConfig pairs a dependency-free finalizer with a dependent applier around a bridge staged
// with `reload = false`. Terraform destroys the dependent applier before the bridge, leaving only the
// finalizer to activate the staged delete.
func applierFinalizerConfig(iface, cidr string, finalizerOnDestroy bool) string {
	return fmt.Sprintf(`
	resource "proxmox_network_applier" "finalizer" {
		node_name  = "{{.NodeName}}"
		on_create  = false
		on_destroy = %t
	}

	resource "proxmox_network_linux_bridge" "test" {
		node_name  = "{{.NodeName}}"
		name       = "%s"
		address    = "%s"
		reload     = false
		depends_on = [proxmox_network_applier.finalizer]
	}

	resource "proxmox_network_applier" "apply" {
		node_name  = "{{.NodeName}}"
		depends_on = [proxmox_network_linux_bridge.test]
	}`, finalizerOnDestroy, iface, cidr)
}

func TestAccResourceNetworkApplierFinalizer(t *testing.T) {
	te := test.InitEnvironment(t)

	iface := fmt.Sprintf("vmbr%d", gofakeit.Number(10, 9999))
	ipV4cidr := fmt.Sprintf("%s/24", gofakeit.IPv4Address())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		CheckDestroy:             checkNodePendingChanges(te, false),
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(applierFinalizerConfig(iface, ipV4cidr, true)),
				Check: resource.ComposeTestCheckFunc(
					checkInterfaceActive(te, iface, true),
					checkNodePendingChanges(te, false),
				),
			},
		},
	})
}

func TestAccResourceNetworkApplierFinalizerOnDestroyDisabled(t *testing.T) {
	te := test.InitEnvironment(t)

	iface := fmt.Sprintf("vmbr%d", gofakeit.Number(10, 9999))
	ipV4cidr := fmt.Sprintf("%s/24", gofakeit.IPv4Address())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		CheckDestroy:             checkStagedDestroy(te),
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(applierFinalizerConfig(iface, ipV4cidr, false)),
				Check: resource.ComposeTestCheckFunc(
					checkInterfaceActive(te, iface, true),
					checkNodePendingChanges(te, false),
				),
			},
		},
	})
}
