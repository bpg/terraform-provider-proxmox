//go:build acceptance || all

//testacc:tier=heavy
//testacc:resource=pool

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/bpg/terraform-provider-proxmox/utils"
)

func TestAccResourcePool(t *testing.T) {
	if utils.GetAnyStringEnv("TF_ACC") == "" {
		t.Skip("Acceptance tests are disabled")
	}

	te := InitEnvironment(t)

	tests := []struct {
		name string
		step []resource.TestStep
	}{
		{"create pool", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_pool" "test_01" {
						comment = "Managed by Terraform"
						pool_id = "test-01"
					}
				`),
				Check: ResourceAttributes("proxmox_virtual_environment_pool.test_01", map[string]string{
					"pool_id": "test-01",
					"comment": "Managed by Terraform",
				}),
			},
		}},
		{"create nested pool", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_pool" "test_02" {
  						comment = "Managed by Terraform"
  						pool_id = "test-02"
					}
					resource "proxmox_virtual_environment_pool" "test_02_nested" {
						depends_on = [proxmox_virtual_environment_pool.test_02]
  						comment = "Managed by Terraform"
  						pool_id = "test-02/test-02-01"
					}
				`),
				Check: ResourceAttributes("proxmox_virtual_environment_pool.test_02_nested", map[string]string{
					"pool_id": "test-02/test-02-01",
					"comment": "Managed by Terraform",
				}),
			},
		}},
		{"change pool description", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_pool" "test_03" {
  						comment = "Managed by Terraform"
  						pool_id = "test-03"
					}
				`),
				Check: ResourceAttributes("proxmox_virtual_environment_pool.test_03", map[string]string{
					"pool_id": "test-03",
					"comment": "Managed by Terraform",
				}),
			},
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_pool" "test_03" {
  						comment = "Still managed by Terraform"
  						pool_id = "test-03"
					}
				`),
				Check: ResourceAttributes("proxmox_virtual_environment_pool.test_03", map[string]string{
					"pool_id": "test-03",
					"comment": "Still managed by Terraform",
				}),
			},
		}},
		{"changing pool name forces replacement", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_pool" "test_04" {
  						comment = "Managed by Terraform"
  						pool_id = "test-04"
					}
				`),
				Check: ResourceAttributes("proxmox_virtual_environment_pool.test_04", map[string]string{
					"pool_id": "test-04",
					"comment": "Managed by Terraform",
				}),
			},
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_pool" "test_04" {
  						comment = "Managed by Terraform"
  						pool_id = "test-04-new"
					}
				`),
				Check: ResourceAttributes("proxmox_virtual_environment_pool.test_04", map[string]string{
					"pool_id": "test-04-new",
					"comment": "Managed by Terraform",
				}),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_virtual_environment_pool.test_04", plancheck.ResourceActionReplace),
					},
				},
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.step,
			})
		})
	}
}
