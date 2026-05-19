//go:build acceptance || all

//testacc:tier=light
//testacc:resource=ceph

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package pool_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceCephPool(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)
	te.RequireCeph()

	basicName := test.SafeResourceName("ceph-pool-basic")
	fullName := test.SafeResourceName("ceph-pool-full")
	updateName := test.SafeResourceName("ceph-pool-update")
	cephfsName := test.SafeResourceName("ceph-pool-cephfs")
	invalidName := test.SafeResourceName("ceph-pool-invalid")

	tests := []struct {
		name string
		step []resource.TestStep
	}{
		{"basic create and import", []resource.TestStep{
			{
				Config: te.RenderConfig(fmt.Sprintf(`
					resource "proxmox_ceph_pool" "basic" {
						node_name = "{{.NodeName}}"
						name      = %q
					}`, basicName)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_ceph_pool.basic", "name", basicName),
					resource.TestCheckResourceAttr("proxmox_ceph_pool.basic", "node_name", te.NodeName),
					test.ResourceAttributesSet("proxmox_ceph_pool.basic", []string{
						"application",
						"size",
						"min_size",
						"pg_num",
						"pg_autoscale_mode",
					}),
				),
			},
			{
				ResourceName:      "proxmox_ceph_pool.basic",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importIDFunc("proxmox_ceph_pool.basic"),
				ImportStateVerifyIgnore: []string{
					// Create-only side effect; not round-tripped through Read.
					"add_storages",
					// Delete-only flags; not round-tripped through Read.
					"force_destroy",
					"remove_storages",
					"remove_ecprofile",
				},
			},
		}},
		{"create with full settings", []resource.TestStep{
			{
				Config: te.RenderConfig(fmt.Sprintf(`
					resource "proxmox_ceph_pool" "full" {
						node_name           = "{{.NodeName}}"
						name                = %q
						application         = "rbd"
						size                = 2
						min_size            = 1
						pg_num              = 32
						pg_autoscale_mode   = "on"
						target_size_ratio   = 0.1
					}`, fullName)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_ceph_pool.full", map[string]string{
						"name":              fullName,
						"application":       "rbd",
						"size":              "2",
						"min_size":          "1",
						"pg_num":            "32",
						"pg_autoscale_mode": "on",
						"target_size_ratio": "0.1",
					}),
				),
			},
		}},
		{"update path", []resource.TestStep{
			{
				Config: te.RenderConfig(fmt.Sprintf(`
					resource "proxmox_ceph_pool" "update" {
						node_name         = "{{.NodeName}}"
						name              = %q
						size              = 3
						min_size          = 2
						pg_num            = 32
						pg_autoscale_mode = "warn"
					}`, updateName)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_ceph_pool.update", map[string]string{
						"size":              "3",
						"min_size":          "2",
						"pg_autoscale_mode": "warn",
					}),
				),
			},
			{
				Config: te.RenderConfig(fmt.Sprintf(`
					resource "proxmox_ceph_pool" "update" {
						node_name         = "{{.NodeName}}"
						name              = %q
						size              = 2
						min_size          = 1
						pg_num            = 32
						pg_autoscale_mode = "on"
						target_size_ratio = 0.05
					}`, updateName)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_ceph_pool.update", map[string]string{
						"size":              "2",
						"min_size":          "1",
						"pg_autoscale_mode": "on",
						"target_size_ratio": "0.05",
					}),
				),
			},
		}},
		{"create with cephfs application", []resource.TestStep{
			{
				Config: te.RenderConfig(fmt.Sprintf(`
					resource "proxmox_ceph_pool" "cephfs" {
						node_name   = "{{.NodeName}}"
						name        = %q
						application = "cephfs"
						pg_num      = 8
						size        = 2
						min_size    = 1
					}`, cephfsName)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_ceph_pool.cephfs", map[string]string{
						"name":        cephfsName,
						"application": "cephfs",
					}),
				),
			},
		}},
		{"invalid pg_autoscale_mode", []resource.TestStep{
			{
				Config: te.RenderConfig(fmt.Sprintf(`
					resource "proxmox_ceph_pool" "invalid" {
						node_name         = "{{.NodeName}}"
						name              = %q
						pg_autoscale_mode = "bogus"
					}`, invalidName)),
				ExpectError: regexp.MustCompile(`(?s)pg_autoscale_mode.*one of`),
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.step,
			})
		})
	}
}

// importIDFunc builds the composite import id `node_name/pool_name` from the resource state.
func importIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource %s not found in state", resourceName)
		}

		nodeName := rs.Primary.Attributes["node_name"]
		name := rs.Primary.Attributes["name"]

		return fmt.Sprintf("%s/%s", nodeName, name), nil
	}
}
