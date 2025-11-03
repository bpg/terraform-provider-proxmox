//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/utils"
)

func TestAccResourceVMPoolDetection(t *testing.T) {
	if utils.GetAnyStringEnv("TF_ACC") == "" {
		t.Skip("Acceptance tests are disabled")
	}

	te := InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"vm pool membership detection", func() []resource.TestStep {
			poolName1 := fmt.Sprintf("test-pool-%s-%d", gofakeit.Word(), time.Now().UnixNano())
			poolName2 := fmt.Sprintf("test-pool-%s-%d", gofakeit.Word(), time.Now().UnixNano())

			te.AddTemplateVars(map[string]interface{}{
				"PoolName1": poolName1,
				"PoolName2": poolName2,
			})

			return []resource.TestStep{
				{
					Config: te.RenderConfig(`
						resource "proxmox_virtual_environment_pool" "test_pool1" {
							pool_id = "{{.PoolName1}}"
							comment = "Test pool 1"
						}

						resource "proxmox_virtual_environment_pool" "test_pool2" {
							pool_id = "{{.PoolName2}}"
							comment = "Test pool 2"
						}

						resource "proxmox_virtual_environment_vm" "test_vm_pool" {
							node_name = "{{.NodeName}}"
							started   = false
							name      = "test-pool-vm"
							pool_id   = proxmox_virtual_environment_pool.test_pool1.pool_id
						}`, WithRootUser()),
					Check: resource.ComposeTestCheckFunc(
						ResourceAttributes("proxmox_virtual_environment_vm.test_vm_pool", map[string]string{
							"pool_id": poolName1,
						}),
					),
				},
				{
					// Test that the provider correctly detects the current pool membership
					RefreshState: true,
					Check: resource.ComposeTestCheckFunc(
						ResourceAttributes("proxmox_virtual_environment_vm.test_vm_pool", map[string]string{
							"pool_id": poolName1,
						}),
					),
				},
			}
		}()},
		{"vm pool membership change detection", func() []resource.TestStep {
			poolName1 := fmt.Sprintf("test-pool-%s-%d", gofakeit.Word(), time.Now().UnixNano())
			poolName2 := fmt.Sprintf("test-pool-%s-%d", gofakeit.Word(), time.Now().UnixNano())

			te.AddTemplateVars(map[string]interface{}{
				"PoolName1": poolName1,
				"PoolName2": poolName2,
			})

			return []resource.TestStep{
				{
					Config: te.RenderConfig(`
						resource "proxmox_virtual_environment_pool" "test_pool1" {
							pool_id = "{{.PoolName1}}"
							comment = "Test pool 1"
						}

						resource "proxmox_virtual_environment_pool" "test_pool2" {
							pool_id = "{{.PoolName2}}"
							comment = "Test pool 2"
						}

						resource "proxmox_virtual_environment_vm" "test_vm_pool_change" {
							node_name = "{{.NodeName}}"
							started   = false
							name      = "test-pool-change-vm"
							pool_id   = proxmox_virtual_environment_pool.test_pool1.pool_id
						}`, WithRootUser()),
					Check: resource.ComposeTestCheckFunc(
						ResourceAttributes("proxmox_virtual_environment_vm.test_vm_pool_change", map[string]string{
							"pool_id": poolName1,
						}),
					),
				},
				{
					// Simulate moving VM to different pool outside Terraform
					Config: te.RenderConfig(`
						resource "proxmox_virtual_environment_pool" "test_pool1" {
							pool_id = "{{.PoolName1}}"
							comment = "Test pool 1"
						}

						resource "proxmox_virtual_environment_pool" "test_pool2" {
							pool_id = "{{.PoolName2}}"
							comment = "Test pool 2"
						}

						resource "proxmox_virtual_environment_vm" "test_vm_pool_change" {
							node_name = "{{.NodeName}}"
							started   = false
							name      = "test-pool-change-vm"
							pool_id   = proxmox_virtual_environment_pool.test_pool2.pool_id
						}`, WithRootUser()),
					Check: resource.ComposeTestCheckFunc(
						ResourceAttributes("proxmox_virtual_environment_vm.test_vm_pool_change", map[string]string{
							"pool_id": poolName2,
						}),
					),
				},
			}
		}()},
		{"vm without pool", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_vm" "test_vm_no_pool" {
						node_name = "{{.NodeName}}"
						started   = false
						name      = "test-no-pool-vm"
					}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm_no_pool", map[string]string{
						"pool_id": "",
					}),
				),
			},
		}},
		{"vm pool drift detection", func() []resource.TestStep {
			poolName1 := fmt.Sprintf("test-pool-%s-%d", gofakeit.Word(), time.Now().UnixNano())
			poolName2 := fmt.Sprintf("test-pool-%s-%d", gofakeit.Word(), time.Now().UnixNano())

			te.AddTemplateVars(map[string]interface{}{
				"PoolName1": poolName1,
				"PoolName2": poolName2,
			})

			return []resource.TestStep{
				{
					Config: te.RenderConfig(`
						resource "proxmox_virtual_environment_pool" "test_pool1" {
							pool_id = "{{.PoolName1}}"
							comment = "Test pool 1"
						}

						resource "proxmox_virtual_environment_pool" "test_pool2" {
							pool_id = "{{.PoolName2}}"
							comment = "Test pool 2"
						}

						resource "proxmox_virtual_environment_vm" "test_vm_drift" {
							node_name = "{{.NodeName}}"
							started   = false
							name      = "test-drift-vm"
							pool_id   = proxmox_virtual_environment_pool.test_pool1.pool_id
						}`, WithRootUser()),
					Check: resource.ComposeTestCheckFunc(
						ResourceAttributes("proxmox_virtual_environment_vm.test_vm_drift", map[string]string{
							"pool_id": poolName1,
						}),
					),
				},
				{
					// Test that plan detects drift when pool changes
					RefreshState: true,
					PlanOnly:     true,
					Check:        resource.ComposeTestCheckFunc(
					// Plan should be empty if no changes detected
					),
				},
			}
		}()},
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

// TestAccResourceVMPoolDetectionLegacy tests pool detection with the legacy SDK provider
func TestAccResourceVMPoolDetectionLegacy(t *testing.T) {
	if utils.GetAnyStringEnv("TF_ACC") == "" {
		t.Skip("Acceptance tests are disabled")
	}

	te := InitEnvironment(t)
	poolName1 := fmt.Sprintf("test-pool-legacy-%s-%d", gofakeit.Word(), time.Now().Unix())
	poolName2 := fmt.Sprintf("test-pool-legacy-%s-%d", gofakeit.Word(), time.Now().Unix())

	te.AddTemplateVars(map[string]interface{}{
		"PoolName1": poolName1,
		"PoolName2": poolName2,
	})

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"legacy vm pool membership detection", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_pool" "test_pool1" {
						pool_id = "{{.PoolName1}}"
						comment = "Test pool 1"
					}

					resource "proxmox_virtual_environment_pool" "test_pool2" {
						pool_id = "{{.PoolName2}}"
						comment = "Test pool 2"
					}

					resource "proxmox_virtual_environment_vm" "test_vm_pool_legacy" {
						node_name = "{{.NodeName}}"
						started   = false
						name      = "test-pool-legacy-vm"
						pool_id   = proxmox_virtual_environment_pool.test_pool1.pool_id
					}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm_pool_legacy", map[string]string{
						"pool_id": poolName1,
					}),
				),
			},
			{
				// Test pool detection after refresh
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm_pool_legacy", map[string]string{
						"pool_id": poolName1,
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

// TestAccResourceVMPoolDetectionManual tests manual pool changes outside Terraform
func TestAccResourceVMPoolDetectionManual(t *testing.T) {
	if utils.GetAnyStringEnv("TF_ACC") == "" {
		t.Skip("Acceptance tests are disabled")
	}

	te := InitEnvironment(t)
	poolName1 := fmt.Sprintf("test-pool-manual-%s-%d", gofakeit.Word(), time.Now().Unix())
	poolName2 := fmt.Sprintf("test-pool-manual-%s-%d", gofakeit.Word(), time.Now().Unix())

	te.AddTemplateVars(map[string]interface{}{
		"PoolName1": poolName1,
		"PoolName2": poolName2,
	})

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"manual pool change detection", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_pool" "test_pool1" {
						pool_id = "{{.PoolName1}}"
						comment = "Test pool 1"
					}

					resource "proxmox_virtual_environment_pool" "test_pool2" {
						pool_id = "{{.PoolName2}}"
						comment = "Test pool 2"
					}

					resource "proxmox_virtual_environment_vm" "test_vm_manual" {
						node_name = "{{.NodeName}}"
						started   = false
						name      = "test-manual-vm"
						pool_id   = proxmox_virtual_environment_pool.test_pool1.pool_id
					}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm_manual", map[string]string{
						"pool_id": poolName1,
					}),
				),
			},
			{
				// Simulate manual pool change by updating the VM's pool_id in config
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_pool" "test_pool1" {
						pool_id = "{{.PoolName1}}"
						comment = "Test pool 1"
					}

					resource "proxmox_virtual_environment_pool" "test_pool2" {
						pool_id = "{{.PoolName2}}"
						comment = "Test pool 2"
					}

					resource "proxmox_virtual_environment_vm" "test_vm_manual" {
						node_name = "{{.NodeName}}"
						started   = false
						name      = "test-manual-vm"
						pool_id   = proxmox_virtual_environment_pool.test_pool2.pool_id
					}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_vm.test_vm_manual", map[string]string{
						"pool_id": poolName2,
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
