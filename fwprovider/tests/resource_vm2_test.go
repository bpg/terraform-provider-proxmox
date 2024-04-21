/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceVM2(t *testing.T) {
	t.Parallel()

	te := initTestEnvironment(t)
	vmID := gofakeit.IntRange(90000, 100000)
	te.addTemplateVars(map[string]any{
		"VMID": vmID,
	})

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create minimal VM", []resource.TestStep{{
			Config: te.renderConfig(`
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
			}`),
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
					"node_name": te.nodeName,
				}),
				testResourceAttributesSet("proxmox_virtual_environment_vm2.test_vm", []string{
					"id",
				}),
			),
		}}},
		{"create minimal VM with ID", []resource.TestStep{{
			Config: te.renderConfig(`
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				
				id = {{.VMID}}
			}`),
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
					"node_name": te.nodeName,
					"id":        strconv.Itoa(vmID),
				}),
			),
		}}},
		{"set an invalid VM name", []resource.TestStep{{
			Config: te.renderConfig(`
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				
				name = "not a valid DNS name"

			}`),
			ExpectError: regexp.MustCompile(`name must be a valid DNS name`),
		}}},
		{"set, update, import with primitive fields", []resource.TestStep{
			{
				Config: te.renderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					
					name = "test-vm"
					description = "test description"
				}`),
				Check: testResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
					"name":        "test-vm",
					"description": "test description",
				}),
			},
			{
				Config: te.renderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					
					name = "test-vm"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
						"name": "test-vm",
					}),
					testNoResourceAttributesSet("proxmox_virtual_environment_vm2.test_vm", []string{
						"description",
					}),
				),
			},
			{
				ResourceName:        "proxmox_virtual_environment_vm2.test_vm",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: te.nodeName + "/",
			},
		}},
		{"set, update, import with tags", []resource.TestStep{
			{
				Config: te.renderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					
					name = "test-tags"
					tags = ["tag2", "tag1"]
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttr("proxmox_virtual_environment_vm2.test_vm", "tags.*", "tag1"),
					resource.TestCheckTypeSetElemAttr("proxmox_virtual_environment_vm2.test_vm", "tags.*", "tag2"),
				),
			},
			{
				Config: te.renderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					
					name = "test-tags"
					tags = ["tag1"]
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_virtual_environment_vm2.test_vm", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("proxmox_virtual_environment_vm2.test_vm", "tags.*", "tag1"),
				),
			},
			{
				Config: te.renderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					
					name = "test-tags"
				}`),
				Check: testNoResourceAttributesSet("proxmox_virtual_environment_vm2.test_vm", []string{
					"tags",
				}),
			},
		}},
		{"a VM can't have empty tags set", []resource.TestStep{{
			Config: te.renderConfig(`
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				
				tags = []
			}`),
			ExpectError: regexp.MustCompile(`tags set must contain at least 1 elements`),
		}}},
		{"a VM can't have empty tags", []resource.TestStep{{
			Config: te.renderConfig(`
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				
				tags = ["", "tag1"]
			}`),
			ExpectError: regexp.MustCompile(`string length must be at least 1, got: 0`),
		}}},
		{"a VM can't have empty tags", []resource.TestStep{{
			Config: te.renderConfig(`
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				
				tags = [" ", "tag1"]
			}`),
			ExpectError: regexp.MustCompile(`must be a non-empty and non-whitespace string`),
		}}},
		{"multiline description", []resource.TestStep{{
			Config: te.renderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					
					description = trimspace(<<-EOT
						my
						description
						value
					EOT
					)
				}`),
			Check: resource.ComposeTestCheckFunc(
				testResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
					"description": "my\ndescription\nvalue",
				}),
			),
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.accProviders,
				Steps:                    tt.steps,
			})
		})
	}
}
