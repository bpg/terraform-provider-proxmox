//go:build acceptance || all

//testacc:tier=medium
//testacc:resource=vm

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vm_test

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceVMShort(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)
	te.AddTemplateVars(map[string]interface{}{
		"TestVMID": 100000 + rand.Intn(99999),
	})

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create minimal VM", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
					"node_name": te.NodeName,
				}),
				test.ResourceAttributesSet("proxmox_vm.test_vm", []string{
					"id",
				}),
			),
		}}},
		{"create minimal VM with ID", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				id = {{.TestVMID}}
			}`),
		}}},
		{"set an invalid VM name", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "not a valid DNS name"
			}`),
			ExpectError: regexp.MustCompile(`name must be a valid DNS name`),
		}}},
		{"set a FQDN VM name", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "vm.example.com"
			}`),
			Check: test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
				"name": "vm.example.com",
			}),
		}}},
		{"set, update, import with primitive fields", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-vm"
					description = "test description"
				}`),
				Check: test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
					"name":        "test-vm",
					"description": "test description",
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-vm"
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
						"name": "test-vm",
					}),
					test.NoResourceAttributesSet("proxmox_vm.test_vm", []string{
						"description",
					}),
				),
			},
			{
				ResourceName:        "proxmox_vm.test_vm",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: te.NodeName + "/",
			},
		}},
		{"set, update, import with tags", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-tags"
					tags = ["tag2", "tag1"]
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttr("proxmox_vm.test_vm", "tags.*", "tag1"),
					resource.TestCheckTypeSetElemAttr("proxmox_vm.test_vm", "tags.*", "tag2"),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-tags"
					tags = ["tag1"]
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_vm.test_vm", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("proxmox_vm.test_vm", "tags.*", "tag1"),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-tags"
					// no tags
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_vm.test_vm", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("proxmox_vm.test_vm", "tags.*", "tag1"),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-tags"
					tags = []
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_vm.test_vm", "tags.#", "0"),
				),
			},
		}},
		{"a VM can't have empty tags", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				tags = ["", "tag1"]
			}`),
			ExpectError: regexp.MustCompile(`string length must be at least 1, got: 0`),
		}}},
		{"a VM can't have empty tags", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				tags = [" ", "tag1"]
			}`),
			ExpectError: regexp.MustCompile(`must be a non-empty and non-whitespace string`),
		}}},
		{"multiline description", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_vm" "test_vm" {
					node_name = "{{.NodeName}}"
					description = trimspace(<<-EOT
						my
						description
						value
					EOT
					)
				}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_vm.test_vm", map[string]string{
					"description": "my\ndescription\nvalue",
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

func TestAccResourceVM2AgentFreezeFS(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)
	resourceName := "proxmox_vm.test_agent"

	checkFreezeFS := func(wantDisabled bool) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			rs, ok := state.RootModule().Resources[resourceName]
			if !ok {
				return fmt.Errorf("resource %q not found", resourceName)
			}

			vmID, err := strconv.Atoi(rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("failed to parse VM ID %q: %w", rs.Primary.ID, err)
			}

			config, err := te.NodeClient().VM(vmID).GetVM(context.Background())
			if err != nil {
				return fmt.Errorf("failed to get VM %d config: %w", vmID, err)
			}

			if wantDisabled {
				if config.Agent == nil || config.Agent.FreezeFS == nil {
					return fmt.Errorf("VM %d has no freeze_fs API value", vmID)
				}
				if *config.Agent.FreezeFS {
					return fmt.Errorf("VM %d freeze_fs is enabled, want disabled", vmID)
				}
			} else if config.Agent != nil && config.Agent.FreezeFS != nil && !*config.Agent.FreezeFS {
				return fmt.Errorf("VM %d freeze_fs is disabled, want enabled", vmID)
			}

			return nil
		}
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_vm" "test_agent" {
						node_name = "{{.NodeName}}"
						name      = "test-agent-freeze-fs"
						agent = {
							enabled   = true
							trim      = false
							type      = "virtio"
							freeze_fs = false
						}
					}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes(resourceName, map[string]string{
						"agent.enabled":   "true",
						"agent.trim":      "false",
						"agent.type":      "virtio",
						"agent.freeze_fs": "false",
					}),
					checkFreezeFS(true),
				),
			},
			{
				Config: te.RenderConfig(`
					resource "proxmox_vm" "test_agent" {
						node_name = "{{.NodeName}}"
						name      = "test-agent-freeze-fs"
						agent = {
							enabled   = true
							trim      = true
							type      = "virtio"
							freeze_fs = true
						}
					}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes(resourceName, map[string]string{
						"agent.enabled":   "true",
						"agent.trim":      "true",
						"agent.type":      "virtio",
						"agent.freeze_fs": "true",
					}),
					checkFreezeFS(false),
				),
			},
		},
	})
}
