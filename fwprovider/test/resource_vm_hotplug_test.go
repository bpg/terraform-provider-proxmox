//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceVMHotplug(t *testing.T) {
	te := InitEnvironment(t)
	imageFileID := te.DownloadCloudImage()

	t.Run("add disk to running VM", func(t *testing.T) {
		te := InitEnvironment(t)
		te.AddTemplateVars(map[string]any{"ImageFileID": imageFileID})

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: te.AccProviders,
			Steps: []resource.TestStep{
				{
					Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_hotplug" {
						node_name       = "{{.NodeName}}"
						started         = true
						stop_on_destroy = true
						name            = "test-hotplug-disk"

						cpu {
							cores = 2
						}
						memory {
							dedicated = 2048
						}
						disk {
							datastore_id = "local-lvm"
							file_id      = "{{.ImageFileID}}"
							interface    = "scsi0"
							size         = 20
						}
						initialization {
							ip_config {
								ipv4 {
									address = "dhcp"
								}
							}
						}
						network_device {
							bridge = "vmbr0"
						}
					}`),
					Check: resource.ComposeTestCheckFunc(
						ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
							"disk.#": "1",
						}),
					),
				},
				{
					Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_hotplug" {
						node_name       = "{{.NodeName}}"
						started         = true
						stop_on_destroy = true
						name            = "test-hotplug-disk"

						cpu {
							cores = 2
						}
						memory {
							dedicated = 2048
						}
						disk {
							datastore_id = "local-lvm"
							file_id      = "{{.ImageFileID}}"
							interface    = "scsi0"
							size         = 20
						}
						disk {
							datastore_id = "local-lvm"
							interface    = "scsi1"
							size         = 4
						}
						initialization {
							ip_config {
								ipv4 {
									address = "dhcp"
								}
							}
						}
						network_device {
							bridge = "vmbr0"
						}
					}`),
					Check: resource.ComposeTestCheckFunc(
						ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
							"disk.#": "2",
						}),
					),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction("proxmox_virtual_environment_vm.test_hotplug", plancheck.ResourceActionUpdate),
						},
					},
				},
			},
		})
	})

	t.Run("add network device to running VM", func(t *testing.T) {
		te := InitEnvironment(t)
		te.AddTemplateVars(map[string]any{"ImageFileID": imageFileID})

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: te.AccProviders,
			Steps: []resource.TestStep{
				{
					Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_hotplug" {
						node_name       = "{{.NodeName}}"
						started         = true
						stop_on_destroy = true
						name            = "test-hotplug-network"

						cpu {
							cores = 2
						}
						memory {
							dedicated = 2048
						}
						disk {
							datastore_id = "local-lvm"
							file_id      = "{{.ImageFileID}}"
							interface    = "scsi0"
							size         = 20
						}
						initialization {
							ip_config {
								ipv4 {
									address = "dhcp"
								}
							}
						}
						network_device {
							bridge = "vmbr0"
						}
					}`),
					Check: resource.ComposeTestCheckFunc(
						ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
							"network_device.#": "1",
						}),
					),
				},
				{
					Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_hotplug" {
						node_name       = "{{.NodeName}}"
						started         = true
						stop_on_destroy = true
						name            = "test-hotplug-network"

						cpu {
							cores = 2
						}
						memory {
							dedicated = 2048
						}
						disk {
							datastore_id = "local-lvm"
							file_id      = "{{.ImageFileID}}"
							interface    = "scsi0"
							size         = 20
						}
						initialization {
							ip_config {
								ipv4 {
									address = "dhcp"
								}
							}
						}
						network_device {
							bridge = "vmbr0"
						}
						network_device {
							bridge = "vmbr0"
							model  = "virtio"
						}
					}`),
					Check: resource.ComposeTestCheckFunc(
						ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
							"network_device.#": "2",
						}),
					),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction("proxmox_virtual_environment_vm.test_hotplug", plancheck.ResourceActionUpdate),
						},
					},
				},
			},
		})
	})

	t.Run("increase memory on running VM", func(t *testing.T) {
		te := InitEnvironment(t)
		te.AddTemplateVars(map[string]any{"ImageFileID": imageFileID})

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: te.AccProviders,
			Steps: []resource.TestStep{
				{
					Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_hotplug" {
						node_name       = "{{.NodeName}}"
						started         = true
						stop_on_destroy = true
						name            = "test-hotplug-memory"

						cpu {
							cores = 2
						}
						memory {
							dedicated = 2048
						}
						disk {
							datastore_id = "local-lvm"
							file_id      = "{{.ImageFileID}}"
							interface    = "scsi0"
							size         = 20
						}
						initialization {
							ip_config {
								ipv4 {
									address = "dhcp"
								}
							}
						}
						network_device {
							bridge = "vmbr0"
						}
					}`),
					Check: resource.ComposeTestCheckFunc(
						ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
							"memory.0.dedicated": "2048",
						}),
					),
				},
				{
					Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_hotplug" {
						node_name       = "{{.NodeName}}"
						started         = true
						stop_on_destroy = true
						name            = "test-hotplug-memory"

						cpu {
							cores = 2
						}
						memory {
							dedicated = 4096
						}
						disk {
							datastore_id = "local-lvm"
							file_id      = "{{.ImageFileID}}"
							interface    = "scsi0"
							size         = 20
						}
						initialization {
							ip_config {
								ipv4 {
									address = "dhcp"
								}
							}
						}
						network_device {
							bridge = "vmbr0"
						}
					}`),
					Check: resource.ComposeTestCheckFunc(
						ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
							"memory.0.dedicated": "4096",
						}),
					),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction("proxmox_virtual_environment_vm.test_hotplug", plancheck.ResourceActionUpdate),
						},
					},
				},
			},
		})
	})

	t.Run("memory increase without hotplug requires reboot", func(t *testing.T) {
		te := InitEnvironment(t)
		te.AddTemplateVars(map[string]any{"ImageFileID": imageFileID})

		var capturedUptime int

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: te.AccProviders,
			Steps: []resource.TestStep{
				{
					Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_hotplug" {
						node_name           = "{{.NodeName}}"
						started             = true
						stop_on_destroy     = true
						reboot_after_update = true
						name                = "test-reboot-memory"

						cpu {
							cores = 2
						}
						memory {
							dedicated = 2048
						}
						disk {
							datastore_id = "local-lvm"
							file_id      = "{{.ImageFileID}}"
							interface    = "scsi0"
							size         = 20
						}
						initialization {
							ip_config {
								ipv4 {
									address = "dhcp"
								}
							}
						}
						network_device {
							bridge = "vmbr0"
						}
					}`),
					Check: resource.ComposeTestCheckFunc(
						ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
							"memory.0.dedicated": "2048",
						}),
						func(s *terraform.State) error {
							rs, ok := s.RootModule().Resources["proxmox_virtual_environment_vm.test_hotplug"]
							if !ok {
								return fmt.Errorf("resource not found")
							}

							vmID, err := strconv.Atoi(rs.Primary.Attributes["vm_id"])
							if err != nil {
								return fmt.Errorf("failed to parse vm_id: %w", err)
							}

							// wait a bit for uptime to accumulate
							time.Sleep(5 * time.Second)

							ctx := context.Background()

							status, err := te.NodeClient().VM(vmID).GetVMStatus(ctx)
							if err != nil {
								return fmt.Errorf("failed to get VM status: %w", err)
							}

							if status.Uptime == nil || *status.Uptime < 3 {
								return fmt.Errorf("VM uptime too low, expected >= 3 seconds, got %v", status.Uptime)
							}

							capturedUptime = *status.Uptime

							return nil
						},
					),
				},
				{
					Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_hotplug" {
						node_name           = "{{.NodeName}}"
						started             = true
						stop_on_destroy     = true
						reboot_after_update = true
						name                = "test-reboot-memory"

						cpu {
							cores = 2
						}
						memory {
							dedicated = 4096
						}
						disk {
							datastore_id = "local-lvm"
							file_id      = "{{.ImageFileID}}"
							interface    = "scsi0"
							size         = 20
						}
						initialization {
							ip_config {
								ipv4 {
									address = "dhcp"
								}
							}
						}
						network_device {
							bridge = "vmbr0"
						}
					}`),
					Check: resource.ComposeTestCheckFunc(
						ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
							"memory.0.dedicated": "4096",
						}),
						func(s *terraform.State) error {
							rs, ok := s.RootModule().Resources["proxmox_virtual_environment_vm.test_hotplug"]
							if !ok {
								return fmt.Errorf("resource not found")
							}

							vmID, err := strconv.Atoi(rs.Primary.Attributes["vm_id"])
							if err != nil {
								return fmt.Errorf("failed to parse vm_id: %w", err)
							}

							ctx := context.Background()

							status, err := te.NodeClient().VM(vmID).GetVMStatus(ctx)
							if err != nil {
								return fmt.Errorf("failed to get VM status: %w", err)
							}

							if status.Uptime == nil {
								return fmt.Errorf("VM uptime is nil")
							}

							// uptime should be reset (reboot happened) - new uptime should be less than captured
							if *status.Uptime >= capturedUptime {
								return fmt.Errorf("VM was NOT rebooted: uptime before=%d, after=%d (expected reboot)", capturedUptime, *status.Uptime)
							}

							return nil
						},
					),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction("proxmox_virtual_environment_vm.test_hotplug", plancheck.ResourceActionUpdate),
						},
					},
				},
			},
		})
	})

	t.Run("memory increase with hotplug enabled skips reboot", func(t *testing.T) {
		te := InitEnvironment(t)
		te.AddTemplateVars(map[string]any{"ImageFileID": imageFileID})

		var capturedUptime int

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: te.AccProviders,
			Steps: []resource.TestStep{
				{
					Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_hotplug" {
						node_name           = "{{.NodeName}}"
						started             = true
						stop_on_destroy     = true
						reboot_after_update = true
						hotplug             = "memory"
						name                = "test-hotplug-memory-enabled"

						cpu {
							cores = 2
							numa  = true
						}
						memory {
							dedicated = 2048
						}
						disk {
							datastore_id = "local-lvm"
							file_id      = "{{.ImageFileID}}"
							interface    = "scsi0"
							size         = 20
						}
						initialization {
							ip_config {
								ipv4 {
									address = "dhcp"
								}
							}
						}
						network_device {
							bridge = "vmbr0"
						}
					}`),
					Check: resource.ComposeTestCheckFunc(
						ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
							"memory.0.dedicated": "2048",
						}),
						func(s *terraform.State) error {
							rs, ok := s.RootModule().Resources["proxmox_virtual_environment_vm.test_hotplug"]
							if !ok {
								return fmt.Errorf("resource not found")
							}

							vmID, err := strconv.Atoi(rs.Primary.Attributes["vm_id"])
							if err != nil {
								return fmt.Errorf("failed to parse vm_id: %w", err)
							}

							// wait a bit for uptime to accumulate
							time.Sleep(5 * time.Second)

							ctx := context.Background()

							status, err := te.NodeClient().VM(vmID).GetVMStatus(ctx)
							if err != nil {
								return fmt.Errorf("failed to get VM status: %w", err)
							}

							if status.Uptime == nil || *status.Uptime < 3 {
								return fmt.Errorf("VM uptime too low, expected >= 3 seconds, got %v", status.Uptime)
							}

							capturedUptime = *status.Uptime

							return nil
						},
					),
				},
				{
					Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_hotplug" {
						node_name           = "{{.NodeName}}"
						started             = true
						stop_on_destroy     = true
						reboot_after_update = true
						hotplug             = "memory"
						name                = "test-hotplug-memory-enabled"

						cpu {
							cores = 2
							numa  = true
						}
						memory {
							dedicated = 4096
						}
						disk {
							datastore_id = "local-lvm"
							file_id      = "{{.ImageFileID}}"
							interface    = "scsi0"
							size         = 20
						}
						initialization {
							ip_config {
								ipv4 {
									address = "dhcp"
								}
							}
						}
						network_device {
							bridge = "vmbr0"
						}
					}`),
					Check: resource.ComposeTestCheckFunc(
						ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
							"memory.0.dedicated": "4096",
						}),
						func(s *terraform.State) error {
							rs, ok := s.RootModule().Resources["proxmox_virtual_environment_vm.test_hotplug"]
							if !ok {
								return fmt.Errorf("resource not found")
							}

							vmID, err := strconv.Atoi(rs.Primary.Attributes["vm_id"])
							if err != nil {
								return fmt.Errorf("failed to parse vm_id: %w", err)
							}

							ctx := context.Background()

							status, err := te.NodeClient().VM(vmID).GetVMStatus(ctx)
							if err != nil {
								return fmt.Errorf("failed to get VM status: %w", err)
							}

							if status.Uptime == nil {
								return fmt.Errorf("VM uptime is nil")
							}

							// uptime should have increased (no reboot), allow some tolerance
							if *status.Uptime < capturedUptime {
								return fmt.Errorf("VM was rebooted: uptime before=%d, after=%d (expected no reboot)", capturedUptime, *status.Uptime)
							}

							return nil
						},
					),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction("proxmox_virtual_environment_vm.test_hotplug", plancheck.ResourceActionUpdate),
						},
					},
				},
			},
		})
	})

	t.Run("change CPU cores requires reboot", func(t *testing.T) {
		te := InitEnvironment(t)
		te.AddTemplateVars(map[string]any{"ImageFileID": imageFileID})

		var capturedUptime int

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: te.AccProviders,
			Steps: []resource.TestStep{
				{
					Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_hotplug" {
						node_name           = "{{.NodeName}}"
						started             = true
						stop_on_destroy     = true
						reboot_after_update = true
						name                = "test-reboot-cores"

						cpu {
							cores = 2
						}
						memory {
							dedicated = 2048
						}
						disk {
							datastore_id = "local-lvm"
							file_id      = "{{.ImageFileID}}"
							interface    = "scsi0"
							size         = 20
						}
						initialization {
							ip_config {
								ipv4 {
									address = "dhcp"
								}
							}
						}
						network_device {
							bridge = "vmbr0"
						}
					}`),
					Check: resource.ComposeTestCheckFunc(
						ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
							"cpu.0.cores": "2",
						}),
						func(s *terraform.State) error {
							rs, ok := s.RootModule().Resources["proxmox_virtual_environment_vm.test_hotplug"]
							if !ok {
								return fmt.Errorf("resource not found")
							}

							vmID, err := strconv.Atoi(rs.Primary.Attributes["vm_id"])
							if err != nil {
								return fmt.Errorf("failed to parse vm_id: %w", err)
							}

							// wait a bit for uptime to accumulate
							time.Sleep(5 * time.Second)

							ctx := context.Background()

							status, err := te.NodeClient().VM(vmID).GetVMStatus(ctx)
							if err != nil {
								return fmt.Errorf("failed to get VM status: %w", err)
							}

							if status.Uptime == nil || *status.Uptime < 3 {
								return fmt.Errorf("VM uptime too low, expected >= 3 seconds, got %v", status.Uptime)
							}

							capturedUptime = *status.Uptime

							return nil
						},
					),
				},
				{
					Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_hotplug" {
						node_name           = "{{.NodeName}}"
						started             = true
						stop_on_destroy     = true
						reboot_after_update = true
						name                = "test-reboot-cores"

						cpu {
							cores = 4
						}
						memory {
							dedicated = 2048
						}
						disk {
							datastore_id = "local-lvm"
							file_id      = "{{.ImageFileID}}"
							interface    = "scsi0"
							size         = 20
						}
						initialization {
							ip_config {
								ipv4 {
									address = "dhcp"
								}
							}
						}
						network_device {
							bridge = "vmbr0"
						}
					}`),
					Check: resource.ComposeTestCheckFunc(
						ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
							"cpu.0.cores": "4",
						}),
						func(s *terraform.State) error {
							rs, ok := s.RootModule().Resources["proxmox_virtual_environment_vm.test_hotplug"]
							if !ok {
								return fmt.Errorf("resource not found")
							}

							vmID, err := strconv.Atoi(rs.Primary.Attributes["vm_id"])
							if err != nil {
								return fmt.Errorf("failed to parse vm_id: %w", err)
							}

							ctx := context.Background()

							status, err := te.NodeClient().VM(vmID).GetVMStatus(ctx)
							if err != nil {
								return fmt.Errorf("failed to get VM status: %w", err)
							}

							if status.Uptime == nil {
								return fmt.Errorf("VM uptime is nil")
							}

							// uptime should be reset (reboot happened) - new uptime should be less than captured
							if *status.Uptime >= capturedUptime {
								return fmt.Errorf("VM was NOT rebooted: uptime before=%d, after=%d (expected reboot)", capturedUptime, *status.Uptime)
							}

							return nil
						},
					),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction("proxmox_virtual_environment_vm.test_hotplug", plancheck.ResourceActionUpdate),
						},
					},
				},
			},
		})
	})

	t.Run("change CPU hotplugged vcpus without hotplug requires reboot", func(t *testing.T) {
		te := InitEnvironment(t)
		te.AddTemplateVars(map[string]any{"ImageFileID": imageFileID})

		var capturedUptime int

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: te.AccProviders,
			Steps: []resource.TestStep{
				{
					Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_hotplug" {
						node_name           = "{{.NodeName}}"
						started             = true
						stop_on_destroy     = true
						reboot_after_update = true
						name                = "test-reboot-vcpus"

						cpu {
							cores      = 4
							hotplugged = 2
						}
						memory {
							dedicated = 2048
						}
						disk {
							datastore_id = "local-lvm"
							file_id      = "{{.ImageFileID}}"
							interface    = "scsi0"
							size         = 20
						}
						initialization {
							ip_config {
								ipv4 {
									address = "dhcp"
								}
							}
						}
						network_device {
							bridge = "vmbr0"
						}
					}`),
					Check: resource.ComposeTestCheckFunc(
						ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
							"cpu.0.cores":      "4",
							"cpu.0.hotplugged": "2",
						}),
						func(s *terraform.State) error {
							rs, ok := s.RootModule().Resources["proxmox_virtual_environment_vm.test_hotplug"]
							if !ok {
								return fmt.Errorf("resource not found")
							}

							vmID, err := strconv.Atoi(rs.Primary.Attributes["vm_id"])
							if err != nil {
								return fmt.Errorf("failed to parse vm_id: %w", err)
							}

							// wait a bit for uptime to accumulate
							time.Sleep(5 * time.Second)

							ctx := context.Background()

							status, err := te.NodeClient().VM(vmID).GetVMStatus(ctx)
							if err != nil {
								return fmt.Errorf("failed to get VM status: %w", err)
							}

							if status.Uptime == nil || *status.Uptime < 3 {
								return fmt.Errorf("VM uptime too low, expected >= 3 seconds, got %v", status.Uptime)
							}

							capturedUptime = *status.Uptime

							return nil
						},
					),
				},
				{
					Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_hotplug" {
						node_name           = "{{.NodeName}}"
						started             = true
						stop_on_destroy     = true
						reboot_after_update = true
						name                = "test-reboot-vcpus"

						cpu {
							cores      = 4
							hotplugged = 3
						}
						memory {
							dedicated = 2048
						}
						disk {
							datastore_id = "local-lvm"
							file_id      = "{{.ImageFileID}}"
							interface    = "scsi0"
							size         = 20
						}
						initialization {
							ip_config {
								ipv4 {
									address = "dhcp"
								}
							}
						}
						network_device {
							bridge = "vmbr0"
						}
					}`),
					Check: resource.ComposeTestCheckFunc(
						ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
							"cpu.0.cores":      "4",
							"cpu.0.hotplugged": "3",
						}),
						func(s *terraform.State) error {
							rs, ok := s.RootModule().Resources["proxmox_virtual_environment_vm.test_hotplug"]
							if !ok {
								return fmt.Errorf("resource not found")
							}

							vmID, err := strconv.Atoi(rs.Primary.Attributes["vm_id"])
							if err != nil {
								return fmt.Errorf("failed to parse vm_id: %w", err)
							}

							ctx := context.Background()

							status, err := te.NodeClient().VM(vmID).GetVMStatus(ctx)
							if err != nil {
								return fmt.Errorf("failed to get VM status: %w", err)
							}

							if status.Uptime == nil {
								return fmt.Errorf("VM uptime is nil")
							}

							// uptime should be reset (reboot happened) - new uptime should be less than captured
							if *status.Uptime >= capturedUptime {
								return fmt.Errorf("VM was NOT rebooted: uptime before=%d, after=%d (expected reboot)", capturedUptime, *status.Uptime)
							}

							return nil
						},
					),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction("proxmox_virtual_environment_vm.test_hotplug", plancheck.ResourceActionUpdate),
						},
					},
				},
			},
		})
	})

	t.Run("change CPU hotplugged vcpus with hotplug enabled skips reboot", func(t *testing.T) {
		te := InitEnvironment(t)
		te.AddTemplateVars(map[string]any{"ImageFileID": imageFileID})

		var capturedUptime int

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: te.AccProviders,
			Steps: []resource.TestStep{
				{
					Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_hotplug" {
						node_name           = "{{.NodeName}}"
						started             = true
						stop_on_destroy     = true
						reboot_after_update = true
						hotplug             = "cpu"
						name                = "test-hotplug-vcpus-enabled"

						cpu {
							cores      = 4
							hotplugged = 2
						}
						memory {
							dedicated = 2048
						}
						disk {
							datastore_id = "local-lvm"
							file_id      = "{{.ImageFileID}}"
							interface    = "scsi0"
							size         = 20
						}
						initialization {
							ip_config {
								ipv4 {
									address = "dhcp"
								}
							}
						}
						network_device {
							bridge = "vmbr0"
						}
					}`),
					Check: resource.ComposeTestCheckFunc(
						ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
							"cpu.0.cores":      "4",
							"cpu.0.hotplugged": "2",
						}),
						func(s *terraform.State) error {
							rs, ok := s.RootModule().Resources["proxmox_virtual_environment_vm.test_hotplug"]
							if !ok {
								return fmt.Errorf("resource not found")
							}

							vmID, err := strconv.Atoi(rs.Primary.Attributes["vm_id"])
							if err != nil {
								return fmt.Errorf("failed to parse vm_id: %w", err)
							}

							// wait a bit for uptime to accumulate
							time.Sleep(5 * time.Second)

							ctx := context.Background()

							status, err := te.NodeClient().VM(vmID).GetVMStatus(ctx)
							if err != nil {
								return fmt.Errorf("failed to get VM status: %w", err)
							}

							if status.Uptime == nil || *status.Uptime < 3 {
								return fmt.Errorf("VM uptime too low, expected >= 3 seconds, got %v", status.Uptime)
							}

							capturedUptime = *status.Uptime

							return nil
						},
					),
				},
				{
					Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_hotplug" {
						node_name           = "{{.NodeName}}"
						started             = true
						stop_on_destroy     = true
						reboot_after_update = true
						hotplug             = "cpu"
						name                = "test-hotplug-vcpus-enabled"

						cpu {
							cores      = 4
							hotplugged = 3
						}
						memory {
							dedicated = 2048
						}
						disk {
							datastore_id = "local-lvm"
							file_id      = "{{.ImageFileID}}"
							interface    = "scsi0"
							size         = 20
						}
						initialization {
							ip_config {
								ipv4 {
									address = "dhcp"
								}
							}
						}
						network_device {
							bridge = "vmbr0"
						}
					}`),
					Check: resource.ComposeTestCheckFunc(
						ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
							"cpu.0.cores":      "4",
							"cpu.0.hotplugged": "3",
						}),
						func(s *terraform.State) error {
							rs, ok := s.RootModule().Resources["proxmox_virtual_environment_vm.test_hotplug"]
							if !ok {
								return fmt.Errorf("resource not found")
							}

							vmID, err := strconv.Atoi(rs.Primary.Attributes["vm_id"])
							if err != nil {
								return fmt.Errorf("failed to parse vm_id: %w", err)
							}

							ctx := context.Background()

							status, err := te.NodeClient().VM(vmID).GetVMStatus(ctx)
							if err != nil {
								return fmt.Errorf("failed to get VM status: %w", err)
							}

							if status.Uptime == nil {
								return fmt.Errorf("VM uptime is nil")
							}

							// uptime should have increased (no reboot), allow some tolerance
							if *status.Uptime < capturedUptime {
								return fmt.Errorf("VM was rebooted: uptime before=%d, after=%d (expected no reboot)", capturedUptime, *status.Uptime)
							}

							return nil
						},
					),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction("proxmox_virtual_environment_vm.test_hotplug", plancheck.ResourceActionUpdate),
						},
					},
				},
			},
		})
	})

	t.Run("change disk properties on running VM", func(t *testing.T) {
		te := InitEnvironment(t)
		te.AddTemplateVars(map[string]any{"ImageFileID": imageFileID})

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: te.AccProviders,
			Steps: []resource.TestStep{
				{
					Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_hotplug" {
						node_name       = "{{.NodeName}}"
						started         = true
						stop_on_destroy = true
						name            = "test-hotplug-disk-props"

						cpu {
							cores = 2
						}
						memory {
							dedicated = 2048
						}
						disk {
							datastore_id = "local-lvm"
							file_id      = "{{.ImageFileID}}"
							interface    = "scsi0"
							size         = 20
							cache        = "none"
						}
						initialization {
							ip_config {
								ipv4 {
									address = "dhcp"
								}
							}
						}
						network_device {
							bridge = "vmbr0"
						}
					}`),
					Check: resource.ComposeTestCheckFunc(
						ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
							"disk.0.cache": "none",
						}),
					),
				},
				{
					Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_hotplug" {
						node_name       = "{{.NodeName}}"
						started         = true
						stop_on_destroy = true
						name            = "test-hotplug-disk-props"

						cpu {
							cores = 2
						}
						memory {
							dedicated = 2048
						}
						disk {
							datastore_id = "local-lvm"
							file_id      = "{{.ImageFileID}}"
							interface    = "scsi0"
							size         = 20
							cache        = "writeback"
							discard      = "on"
						}
						initialization {
							ip_config {
								ipv4 {
									address = "dhcp"
								}
							}
						}
						network_device {
							bridge = "vmbr0"
						}
					}`),
					Check: resource.ComposeTestCheckFunc(
						ResourceAttributes("proxmox_virtual_environment_vm.test_hotplug", map[string]string{
							"disk.0.cache":   "writeback",
							"disk.0.discard": "on",
						}),
					),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction("proxmox_virtual_environment_vm.test_hotplug", plancheck.ResourceActionUpdate),
						},
					},
				},
			},
		})
	})
}
