//go:build acceptance || all

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
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-rng"
			}`),
			Check: test.NoResourceAttributesSet("proxmox_virtual_environment_vm2.test_vm", []string{
				"rng.source",
				"rng.max_bytes",
				"rng.period",
			}),
		}}},
		{"create VM with some rng params", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-rng"
				rng = {
					source = "/dev/urandom"
				}
			}`, test.WithRootUser()),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
					"rng.source": "/dev/urandom",
				}),
				test.NoResourceAttributesSet("proxmox_virtual_environment_vm2.test_vm", []string{
					"rng.max_bytes",
					"rng.period",
				}),
			),
		}}},
		{"create VM with RNG params and then update them", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-rng"
					rng = {
						source = "/dev/urandom"
						max_bytes = 1024
					}
				}`, test.WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
						"rng.source":    "/dev/urandom",
						"rng.max_bytes": "1024",
					}),
					test.NoResourceAttributesSet("proxmox_virtual_environment_vm2.test_vm", []string{
						"rng.period",
					}),
				),
			},
			{ // now update the rng params and check if they are updated
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-rng"
					rng = {
						source = "/dev/random"
						period = 1000
					}
				}`, test.WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
						"rng.source": "/dev/random",
						"rng.period": "1000",
					}),
					test.NoResourceAttributesSet("proxmox_virtual_environment_vm2.test_vm", []string{
						"rng.max_bytes",
					}),
				),
			},
			{
				RefreshState: true,
			},
		}},
		{"clone VM with some rng params", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_vm2" "template_vm" {
				node_name = "{{.NodeName}}"
				name = "template-rng"
				rng = {
					source = "/dev/urandom"
					max_bytes = 1024
				}
			}
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-rng"
				clone = {
					id = proxmox_virtual_environment_vm2.template_vm.id
				}
			}`, test.WithRootUser()),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
					"rng.source":    "/dev/urandom",
					"rng.max_bytes": "1024",
				}),
			),
		}}},
		{"clone VM with some rng params and updating them in the clone", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_vm2" "template_vm" {
				node_name = "{{.NodeName}}"
				name = "template-rng"
				rng = {
					source = "/dev/urandom"
					max_bytes = 1024
				}
			}
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-rng"
				clone = {
					id = proxmox_virtual_environment_vm2.template_vm.id
				}
				rng = {
					source = "/dev/random"
					period = 2000
				}
			}`, test.WithRootUser()),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
					"rng.source":    "/dev/random",
					"rng.period":    "2000",
					"rng.max_bytes": "1024",
				}),
			),
		}}},
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
