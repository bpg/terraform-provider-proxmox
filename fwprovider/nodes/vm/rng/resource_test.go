//go:build acceptance || all

//testacc:tier=medium
//testacc:resource=vm

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package rng_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceVM2RNG(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create VM with no rng params", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-rng"
			}`),
			Check: test.NoResourceAttributesSet("proxmox_vm.test_vm", []string{
				"rng.source",
				"rng.max_bytes",
				"rng.period",
			}),
		}}},
		{"create VM with some rng params", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-rng"
				rng = {
					source = "/dev/urandom"
				}
			}`, test.WithRootUser()),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
					"rng.source": "/dev/urandom",
				}),
				test.NoResourceAttributesSet("proxmox_vm.test_vm", []string{
					"rng.max_bytes",
					"rng.period",
				}),
			),
		}}},
		{"create VM with RNG params and then update them", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-rng"
					rng = {
						source = "/dev/urandom"
						max_bytes = 1024
					}
				}`, test.WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
						"rng.source":    "/dev/urandom",
						"rng.max_bytes": "1024",
					}),
					test.NoResourceAttributesSet("proxmox_vm.test_vm", []string{
						"rng.period",
					}),
				),
			},
			{ // now update the rng params and check if they are updated
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-rng"
					rng = {
						source = "/dev/random"
						period = 1000
					}
				}`, test.WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
						"rng.source": "/dev/random",
						"rng.period": "1000",
					}),
					test.NoResourceAttributesSet("proxmox_vm.test_vm", []string{
						"rng.max_bytes",
					}),
				),
			},
			{
				RefreshState: true,
			},
		}},
		// F37 fix: user-set `max_bytes = 0` must reach PVE on the wire. The previous
		// `ValueInt64() != 0` guard silently dropped it despite the docs claiming 0 disables
		// limiting; the new attribute.Int64PtrFromValue-based path sends 0 through.
		{"create VM with RNG max_bytes=0", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-rng"
				rng = {
					source    = "/dev/urandom"
					max_bytes = 0
				}
			}`, test.WithRootUser()),
			Check: test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
				"rng.source":    "/dev/urandom",
				"rng.max_bytes": "0",
			}),
		}}},
		// Verifies block-level deletion: removing the whole `rng` block after it was set must
		// emit `delete=rng0` on the PUT wire (compound property — same pattern as vga).
		{"add RNG then remove the block entirely", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-rng"
					rng = {
						source = "/dev/urandom"
					}
				}`, test.WithRootUser()),
				Check: test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
					"rng.source": "/dev/urandom",
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-rng"
				}`, test.WithRootUser()),
				Check: test.NoResourceAttributesSet("proxmox_vm.test_vm", []string{
					"rng.source",
					"rng.max_bytes",
					"rng.period",
				}),
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}
