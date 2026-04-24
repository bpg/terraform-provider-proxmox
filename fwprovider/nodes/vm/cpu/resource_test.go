//go:build acceptance || all

//testacc:tier=medium
//testacc:resource=vm

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cpu_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceVM2CPU(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		// After dropping Optional+Computed at the block level, a VM with no cpu block has a null
		// cpu attribute in state — no phantom cpu.cores=1 / cpu.sockets=1 / cpu.type=kvm64.
		{"create VM with no cpu params — block stays null", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-cpu"
			}`),
			Check: test.NoResourceAttributesSet("proxmox_vm.test_vm", []string{
				"cpu.cores",
				"cpu.numa",
				"cpu.sockets",
				"cpu.type",
				"cpu.units",
				"cpu.vcpus",
			}),
		}}},
		{"create VM with some cpu params", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-cpu"
				cpu = {
					cores = 2
					sockets = 2
					type = "host"
					flags = ["+aes"]
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
					"cpu.cores":   "2",
					"cpu.sockets": "2",
					"cpu.type":    "host",
					"cpu.flags.#": "1",
					"cpu.flags.0": `\+aes`,
				}),
			),
		}}},
		{"create VM with cpu params and then update them", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-cpu"
					cpu = {
						# affinity = "0-1"          only root can set affinity
						# architecture = "x86_64"   only root can set architecture
						cores = 2
						limit = 63.5
						numa = false
						sockets = 2
						type = "host"
						units = 1024
						vcpus = 2
						flags = ["+aes"]
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
						"cpu.cores":   "2",
						"cpu.limit":   "63.5",
						"cpu.numa":    "false",
						"cpu.sockets": "2",
						"cpu.type":    "host",
						"cpu.units":   "1024",
						"cpu.vcpus":   "2",
					}),
				),
			},
			{ // now update the cpu params and check if they are updated
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-cpu"
					cpu = {
						cores = 4
						# limit, numa, sockets, type, flags, vcpus all omitted — Optional-only means
						# each one emits delete=<key> on the wire so PVE drops the prior value.
						units = 2048
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
						"cpu.cores": "4",
						"cpu.units": "2048",
					}),
					test.NoResourceAttributesSet("proxmox_vm.test_vm", []string{
						// Legacy `sockets=1` / `type=kvm64` sentinels are gone by design
						// (audit F20-F22). PVE only returns keys the user wrote.
						"cpu.limit",
						"cpu.numa",
						"cpu.sockets",
						"cpu.type",
						"cpu.flags",
						"cpu.vcpus",
					}),
				),
			},
			{
				RefreshState: true,
			},
		}},
		{"create VM with cpu units = 1", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-cpu-units-1"
				cpu = {
					cores = 1
					sockets = 1
					type = "kvm64"
					units = 1
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
					"cpu.cores":   "1",
					"cpu.sockets": "1",
					"cpu.type":    "kvm64",
					"cpu.units":   "1",
				}),
			),
		}}},
		{"create VM with integer cpu limit and verify no drift", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-cpu-limit-int"
					cpu = {
						cores = 1
						limit = 64
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
						"cpu.limit": "64",
					}),
				),
			},
			{
				RefreshState: true,
			},
		}},
		// F24 regression — legacy FillUpdateBody line 190 had `IsDefined(plan.Sockets)` where
		// it should have checked `plan.Limit`, so a plan that kept `limit` while leaving
		// `sockets` null silently dropped the limit from the wire. The new per-field
		// CheckDeleteBody + toAPI pipeline has no such cascade, so this scenario just has to
		// round-trip a limit change without drift to prove the bug can't regress.
		{"update cpu.limit without setting sockets (F24)", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-cpu-limit-no-sockets"
					cpu = {
						cores = 1
						limit = 2.5
					}
				}`),
				Check: test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
					"cpu.limit": "2.5",
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-cpu-limit-no-sockets"
					cpu = {
						cores = 1
						limit = 5
					}
				}`),
				Check: test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
					"cpu.limit": "5",
				}),
			},
			{
				RefreshState: true,
			},
		}},
		// Verify block-level deletion: setting cpu then removing the whole block must emit
		// per-field `delete=affinity|cores|cpulimit|sockets|cpuunits|cpu` on the PUT body. The
		// resulting state has no cpu attributes set (block is null).
		{"add cpu block then remove it", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-cpu-remove"
					cpu = {
						cores = 2
						limit = 3
					}
				}`),
				Check: test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
					"cpu.cores": "2",
					"cpu.limit": "3",
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-cpu-remove"
				}`),
				Check: test.NoResourceAttributesSet("proxmox_vm.test_vm", []string{
					"cpu.cores",
					"cpu.limit",
					"cpu.sockets",
					"cpu.type",
				}),
			},
		}},
		// regression test for https://github.com/bpg/terraform-provider-proxmox/issues/2353
		{"create VM without cpu.units and verify no drift", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-cpu-units-default"
				}`),
			},
			{
				RefreshState: true,
			},
		}},
		// regression test for https://github.com/bpg/terraform-provider-proxmox/issues/2301
		{"create VM with x86-64-v4 CPU type and verify no format drift", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-cpu-x86-64-v4"
					cpu = {
						cores = 1
						sockets = 1
						type = "x86-64-v4"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
						"cpu.cores":   "1",
						"cpu.sockets": "1",
						"cpu.type":    "x86-64-v4",
					}),
				),
			},
			{
				RefreshState: true,
			},
		}},
		// regression test: CPU without type should not set CPUEmulation — and since PVE only
		// surfaces the config lines that were explicitly written (audit Section 4), a cpu
		// block with only `cores` reads back with everything else null: no sockets, no type.
		{"create VM with CPU cores only (no type) and verify no CPUEmulation error", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-cpu-no-type"
					cpu = {
						cores = 2
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
						"cpu.cores": "2",
					}),
					test.NoResourceAttributesSet("proxmox_vm.test_vm", []string{
						"cpu.sockets",
						"cpu.type",
					}),
				),
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
