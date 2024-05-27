package vga_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceVM2VGA(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create VM with no vga params", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-vga"
			}`),
			Check: test.NoResourceAttributesSet("proxmox_virtual_environment_vm2.test_vm", []string{
				// PVE does not set / return anything by default
				"vga.type",
			}),
		}}},
		{"create VM with some vga params", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-vga"
				vga = {
					type = "std"
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
					"vga.type": "std",
				}),
				test.NoResourceAttributesSet("proxmox_virtual_environment_vm2.test_vm", []string{
					"vga.clipboard",
					"vga.memory",
				}),
			),
		}}},
		{"create VM with VGA params and then update them", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-vga"
					vga = {
						type = "std"
						memory = 16
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
						"vga.type":   "std",
						"vga.memory": "16",
					}),
				),
			},
			{ // now update the vga params and check if they are updated
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm2" "test_vm" {
					node_name = "{{.NodeName}}"
					name = "test-cpu"
					vga = {
						type = "qxl"
						clipboard = "vnc"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
						"vga.type":      "qxl",
						"vga.clipboard": "vnc",
					}),
					test.NoResourceAttributesSet("proxmox_virtual_environment_vm2.test_vm", []string{
						"vga.memory",
					}),
				),
			},
			{
				RefreshState: true,
			},
		}},
		{"clone VM with some vga params", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_vm2" "template_vm" {
				node_name = "{{.NodeName}}"
				name = "template-vga"
				vga = {
					type = "qxl"
					clipboard = "vnc"
				}
			}
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-vga"
				clone = {
					id = proxmox_virtual_environment_vm2.template_vm.id
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
					"vga.type":      "qxl",
					"vga.clipboard": "vnc",
				}),
			),
		}}},
		{"clone VM with some vga params and updating them in the clone", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_vm2" "template_vm" {
				node_name = "{{.NodeName}}"
				name = "template-vga"
				vga = {
					type = "qxl"
					clipboard = "vnc"
				}
			}
			resource "proxmox_virtual_environment_vm2" "test_vm" {
				node_name = "{{.NodeName}}"
				name = "test-cpu"
				clone = {
					id = proxmox_virtual_environment_vm2.template_vm.id
				}
				vga = {
					type = "std"
					memory = 16
				}
			}`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_virtual_environment_vm2.test_vm", map[string]string{
					"vga.type":      "std",
					"vga.memory":    "16",
					"vga.clipboard": "vnc",
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
