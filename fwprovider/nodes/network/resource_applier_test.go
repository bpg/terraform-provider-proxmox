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

	// on_destroy = false keeps a replacement's destroy half from reloading, so any activation comes from Create.
	applier := func(onCreate bool, timeout int, version string) string {
		return fmt.Sprintf(`
		resource "proxmox_network_applier" "test" {
			node_name      = "{{.NodeName}}"
			on_create      = %t
			on_destroy     = false
			timeout_reload = %d
			triggers = {
				version = "%s"
			}
			depends_on = [proxmox_network_linux_bridge.test]
		}`, onCreate, timeout, version)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		CheckDestroy:             checkStagedDestroy(te),
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(bridge),
				Check:  checkInterfaceActive(te, iface, false),
			},
			{
				Config: te.RenderConfig(bridge + applier(false, 100, "one")),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_network_applier.test", map[string]string{
						"on_create":  "false",
						"on_destroy": "false",
					}),
					test.ResourceAttributesSet("proxmox_network_applier.test", []string{"id"}),
					checkInterfaceActive(te, iface, false),
				),
			},
			{
				Config: te.RenderConfig(bridge + applier(false, 60, "one")),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_network_applier.test", map[string]string{
						"timeout_reload": "60",
					}),
					checkInterfaceActive(te, iface, false),
				),
			},
			{
				Config: te.RenderConfig(bridge + applier(true, 60, "two")),
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

	iface := fmt.Sprintf("vmbr%d", gofakeit.Number(10, 9999))
	ipV4cidr := fmt.Sprintf("%s/24", gofakeit.IPv4Address())

	applier := func(version string) string {
		return fmt.Sprintf(`
		resource "proxmox_network_applier" "test" {
			node_name  = "{{.NodeName}}"
			on_destroy = false
			triggers = {
				version = "%s"
			}
		}`, version)
	}

	bridge := fmt.Sprintf(`
	resource "proxmox_network_linux_bridge" "test" {
		node_name = "{{.NodeName}}"
		name      = "%s"
		address   = "%s"
		reload    = false
	}`, iface, ipV4cidr)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		CheckDestroy:             checkStagedDestroy(te),
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(applier("one")),
				Check: test.ResourceAttributes("proxmox_network_applier.test", map[string]string{
					"triggers.version": "one",
				}),
			},
			{
				Config: te.RenderConfig(applier("one") + bridge),
				Check:  checkInterfaceActive(te, iface, false),
			},
			{
				Config: te.RenderConfig(applier("two") + bridge),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_network_applier.test", map[string]string{
						"triggers.version": "two",
					}),
					checkInterfaceActive(te, iface, true),
				),
			},
		},
	})
}

// applierFinalizerConfig pairs a dependency-free finalizer with a dependent applier around a bridge staged
// with `reload = false`, mirroring the bundled example. Terraform destroys the dependent applier before the
// bridge, leaving only the finalizer to activate the staged delete.
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
		on_destroy = false
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
