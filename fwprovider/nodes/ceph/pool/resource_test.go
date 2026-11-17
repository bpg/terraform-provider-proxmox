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
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

func TestAccResourceCephPool(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)
	te.RequireCeph()

	basicName := test.SafeResourceName("ceph-pool-basic")
	fullName := test.SafeResourceName("ceph-pool-full")
	updateName := test.SafeResourceName("ceph-pool-update")
	cephfsName := test.SafeResourceName("ceph-pool-cephfs")
	recreateName := test.SafeResourceName("ceph-pool-recreate")
	invalidName := test.SafeResourceName("ceph-pool-invalid")

	recreatePoolIDs := statecheck.CompareValue(compare.ValuesDiffer())

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
						{{with .CephCrushRule}}crush_rule = "{{.}}"{{end}}
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
						"pool_id",
					}),
					// The pool is created without Proxmox's storage bridge.
					checkStorageAbsent(te, basicName),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("proxmox_ceph_pool.basic", tfjsonpath.New("data_pool"), knownvalue.Null()),
					statecheck.ExpectKnownValue("proxmox_ceph_pool.basic", tfjsonpath.New("metadata_pool_name"), knownvalue.Null()),
				},
			},
			{
				ResourceName:      "proxmox_ceph_pool.basic",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importIDFunc("proxmox_ceph_pool.basic"),
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
						pg_num_min          = 4
						pg_autoscale_mode   = "on"
						target_size_ratio   = 0.1
						{{with .CephCrushRule}}crush_rule = "{{.}}"{{end}}
					}`, fullName)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_ceph_pool.full", map[string]string{
						"name":              fullName,
						"application":       "rbd",
						"size":              "2",
						"min_size":          "1",
						"pg_num":            "32",
						"pg_num_min":        "4",
						"pg_autoscale_mode": "on",
						"target_size_ratio": "0.1",
					}),
				),
			},
			{
				// Import round-trip verifies pg_num_min is read back from /status.
				// Under the prior list-based read, fromAPI left pg_num_min unset and
				// ImportStateVerify would mismatch — this step is the behavioral proof.
				ResourceName:      "proxmox_ceph_pool.full",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importIDFunc("proxmox_ceph_pool.full"),
				ImportStateVerifyIgnore: []string{
					// Write-only — see schema MarkdownDescription.
					"target_size_ratio",
				},
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
						{{with .CephCrushRule}}crush_rule = "{{.}}"{{end}}
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
						{{with .CephCrushRule}}crush_rule = "{{.}}"{{end}}
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
			{
				// Changing size alone must keep min_size: Ceph resets it when size changes.
				Config: te.RenderConfig(fmt.Sprintf(`
					resource "proxmox_ceph_pool" "update" {
						node_name         = "{{.NodeName}}"
						name              = %q
						size              = 3
						min_size          = 1
						pg_num            = 32
						pg_autoscale_mode = "on"
						target_size_ratio = 0.05
						{{with .CephCrushRule}}crush_rule = "{{.}}"{{end}}
					}`, updateName)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_ceph_pool.update", map[string]string{
						"size":     "3",
						"min_size": "1",
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
						{{with .CephCrushRule}}crush_rule = "{{.}}"{{end}}
					}`, cephfsName)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_ceph_pool.cephfs", map[string]string{
						"name":        cephfsName,
						"application": "cephfs",
					}),
				),
			},
		}},
		{"recreation assigns a new pool_id", []resource.TestStep{
			{
				Config: te.RenderConfig(cephPoolConfig("recreate", recreateName, "")),
				ConfigStateChecks: []statecheck.StateCheck{
					recreatePoolIDs.AddStateValue("proxmox_ceph_pool.recreate", tfjsonpath.New("pool_id")),
				},
			},
			{
				Taint:  []string{"proxmox_ceph_pool.recreate"},
				Config: te.RenderConfig(cephPoolConfig("recreate", recreateName, "")),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectUnknownValue("proxmox_ceph_pool.recreate", tfjsonpath.New("pool_id")),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					recreatePoolIDs.AddStateValue("proxmox_ceph_pool.recreate", tfjsonpath.New("pool_id")),
				},
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
		{"invalid erasure_coding profile with failure_domain", []resource.TestStep{
			{
				Config: te.RenderConfig(fmt.Sprintf(`
					resource "proxmox_ceph_pool" "invalid" {
						node_name      = "{{.NodeName}}"
						name           = %q
						erasure_coding = {
							k              = 2
							m              = 2
							failure_domain = "osd"
							profile        = "existing"
						}
					}`, invalidName)),
				ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Combination`),
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

func TestAccResourceCephPoolErasureCoded(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)
	te.RequireCeph()

	ecName := test.SafeResourceName("ceph-pool-ec")
	ownerName := test.SafeResourceName("ceph-pool-ec-owner")
	sharerName := test.SafeResourceName("ceph-pool-ec-sharer")

	ecDataPoolIDs := statecheck.CompareValue(compare.ValuesSame())

	tests := []struct {
		name         string
		checkDestroy resource.TestCheckFunc
		step         []resource.TestStep
	}{
		{"create, update and import", checkPoolsAbsent(te, ecName+"-data", ecName+"-metadata"), []resource.TestStep{
			{
				Config: te.RenderConfig(cephECPoolConfig("ec", ecName, 2, `pg_autoscale_mode = "warn"`)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_ceph_pool.ec", map[string]string{
						"name": ecName,
						// size, min_size and crush_rule describe the metadata pool.
						"size":                 "2",
						"min_size":             "1",
						"metadata_pool_name":   ecName + "-metadata",
						"pg_num":               "8",
						"pg_autoscale_mode":    "warn",
						"data_pool.name":       ecName + "-data",
						"data_pool.size":       "4",
						"data_pool.min_size":   "3",
						"data_pool.crush_rule": ecName + "-data",
					}),
					test.ResourceAttributesSet("proxmox_ceph_pool.ec", []string{"metadata_pool_id", "data_pool.pool_id"}),
					resource.TestCheckNoResourceAttr("proxmox_ceph_pool.ec", "pool_id"),
					// Proxmox adds a storage to erasure-coded pools unless told not to.
					checkStorageAbsent(te, ecName),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					ecDataPoolIDs.AddStateValue("proxmox_ceph_pool.ec", tfjsonpath.New("data_pool").AtMapKey("pool_id")),
				},
			},
			{
				// size goes to the metadata pool and pg_autoscale_mode to the data pool, in place.
				Config: te.RenderConfig(cephECPoolConfig("ec", ecName, 3, `pg_autoscale_mode = "off"`)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_ceph_pool.ec", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_ceph_pool.ec", map[string]string{
						"size":              "3",
						"pg_autoscale_mode": "off",
						"data_pool.size":    "4",
					}),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					ecDataPoolIDs.AddStateValue("proxmox_ceph_pool.ec", tfjsonpath.New("data_pool").AtMapKey("pool_id")),
				},
			},
			{
				// erasure_coding can't be read back: the import records it in place rather than replacing the pool.
				ResourceName:       "proxmox_ceph_pool.ec",
				ImportState:        true,
				ImportStateKind:    resource.ImportBlockWithID,
				ImportStateIdFunc:  importIDFunc("proxmox_ceph_pool.ec"),
				ExpectNonEmptyPlan: true,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_ceph_pool.ec", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "proxmox_ceph_pool.ec",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       importIDFunc("proxmox_ceph_pool.ec"),
				ImportStateVerifyIgnore: []string{"erasure_coding"},
			},
		}},
		{"existing profile", checkPoolsAbsent(te, ownerName+"-data", sharerName+"-data"), []resource.TestStep{
			{
				// The sharer uses the profile Proxmox created for the owner. It depends on the owner, so it is
				// destroyed first and must keep the profile; the owner then removes it.
				Config: te.RenderConfig(cephECPoolConfig("owner", ownerName, 2, "") + fmt.Sprintf(`
					resource "proxmox_ceph_pool" "sharer" {
						node_name = "{{.NodeName}}"
						name      = %q
						size      = 2
						min_size  = 1
						pg_num    = 8
						{{with .CephCrushRule}}crush_rule = "{{.}}"{{end}}
						erasure_coding = {
							k       = 2
							m       = 2
							profile = "pve_ec_${proxmox_ceph_pool.owner.name}"
						}
					}`, sharerName)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_ceph_pool.sharer", "data_pool.size", "4"),
					resource.TestCheckResourceAttrPair(
						"proxmox_ceph_pool.sharer", "data_pool.crush_rule",
						"proxmox_ceph_pool.sharer", "data_pool.name",
					),
				),
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				// import blocks need Terraform 1.5.
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.SkipBelow(tfversion.Version1_5_0),
				},
				CheckDestroy: tt.checkDestroy,
				Steps:        tt.step,
			})
		})
	}
}

func TestAccResourceCephPoolImport(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)
	te.RequireCeph()

	plainName := test.SafeResourceName("ceph-pool-import-plain")
	forceName := test.SafeResourceName("ceph-pool-import-force")

	forceConfig := te.RenderConfig(cephPoolConfig("force", forceName, "force_destroy = true"))
	forcePoolIDs := statecheck.CompareValue(compare.ValuesSame())

	tests := []struct {
		name string
		step []resource.TestStep
	}{
		{"import block plans no changes", []resource.TestStep{
			{
				Config: te.RenderConfig(cephPoolConfig("plain", plainName, "")),
			},
			{
				ResourceName:      "proxmox_ceph_pool.plain",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateIdFunc: importIDFunc("proxmox_ceph_pool.plain"),
			},
		}},
		{"import with a state-only setting updates in place without touching the pool", []resource.TestStep{
			{
				Config: forceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					forcePoolIDs.AddStateValue("proxmox_ceph_pool.force", tfjsonpath.New("pool_id")),
				},
			},
			{
				// force_destroy can't be read back, so the import plans recording it, in place.
				ResourceName:       "proxmox_ceph_pool.force",
				ImportState:        true,
				ImportStateKind:    resource.ImportBlockWithID,
				ImportStateIdFunc:  importIDFunc("proxmox_ceph_pool.force"),
				ExpectNonEmptyPlan: true,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_ceph_pool.force", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				// Drop the pool from state without destroying it, then import it for real.
				Config: `
					removed {
						from = proxmox_ceph_pool.force
						lifecycle {
							destroy = false
						}
					}`,
			},
			{
				Config:             forceConfig,
				ResourceName:       "proxmox_ceph_pool.force",
				ImportState:        true,
				ImportStateId:      fmt.Sprintf("%s/%s", te.NodeName, forceName),
				ImportStatePersist: true,
			},
			{
				// Recording force_destroy must only write state, not update the pool.
				Config: forceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_ceph_pool.force", "force_destroy", "true"),
					checkNoPoolUpdateTask(te, forceName),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					forcePoolIDs.AddStateValue("proxmox_ceph_pool.force", tfjsonpath.New("pool_id")),
				},
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				// removed blocks need Terraform 1.7.
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.SkipBelow(tfversion.Version1_7_0),
				},
				Steps: tt.step,
			})
		})
	}
}

// cephPoolConfig renders a small replicated pool with extra attributes appended.
func cephPoolConfig(label, name, extra string) string {
	return fmt.Sprintf(`
		resource "proxmox_ceph_pool" %q {
			node_name         = "{{.NodeName}}"
			name              = %q
			size              = 2
			min_size          = 1
			pg_num            = 8
			pg_autoscale_mode = "warn"
			{{with .CephCrushRule}}crush_rule = "{{.}}"{{end}}
			%s
		}`, label, name, extra)
}

// cephECPoolConfig renders a small erasure-coded pool (k=2, m=2 over OSDs) whose metadata pool has the given size.
func cephECPoolConfig(label, name string, metadataSize int, extra string) string {
	return fmt.Sprintf(`
		resource "proxmox_ceph_pool" %q {
			node_name = "{{.NodeName}}"
			name      = %q
			size      = %d
			min_size  = 1
			pg_num    = 8
			{{with .CephCrushRule}}crush_rule = "{{.}}"{{end}}
			erasure_coding = {
				k              = 2
				m              = 2
				failure_domain = "osd"
				{{with .CephECDeviceClass}}device_class = "{{.}}"{{end}}
			}
			%s
		}`, label, name, metadataSize, extra)
}

// checkStorageAbsent fails if a storage with the given ID exists.
func checkStorageAbsent(te *test.Environment, storageID string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		err := te.Client().DoRequest(ctx, http.MethodGet, "storage/"+storageID, nil, &struct{}{})
		if err == nil {
			return fmt.Errorf("expected no storage %q, but it exists", storageID)
		}

		if !errors.Is(err, api.ErrResourceDoesNotExist) {
			return fmt.Errorf("reading storage %q: %w", storageID, err)
		}

		return nil
	}
}

// checkPoolsAbsent fails if any of the Ceph pools, or a CRUSH rule of the same name, still exists.
func checkPoolsAbsent(te *test.Environment, names ...string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		pools := &struct {
			Data []struct {
				Name string `json:"pool_name"`
			} `json:"data"`
		}{}
		if err := te.Client().DoRequest(ctx, http.MethodGet, fmt.Sprintf("nodes/%s/ceph/pool", te.NodeName), nil, pools); err != nil {
			return fmt.Errorf("listing Ceph pools: %w", err)
		}

		rules := &struct {
			Data []struct {
				Name string `json:"name"`
			} `json:"data"`
		}{}
		if err := te.Client().DoRequest(ctx, http.MethodGet, fmt.Sprintf("nodes/%s/ceph/rules", te.NodeName), nil, rules); err != nil {
			return fmt.Errorf("listing CRUSH rules: %w", err)
		}

		for _, name := range names {
			for _, p := range pools.Data {
				if p.Name == name {
					return fmt.Errorf("expected Ceph pool %q to be destroyed, but it still exists", name)
				}
			}

			for _, r := range rules.Data {
				if r.Name == name {
					return fmt.Errorf("expected the CRUSH rule %q to be removed with its pool, but it still exists", name)
				}
			}
		}

		return nil
	}
}

// checkNoPoolUpdateTask fails if the node ran a pool update (a cephsetpool task) for the pool.
func checkNoPoolUpdateTask(te *test.Environment, poolName string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		resp := &struct {
			Data []struct {
				UPID string `json:"upid"`
				ID   string `json:"id"`
			} `json:"data"`
		}{}

		path := fmt.Sprintf("nodes/%s/tasks?typefilter=cephsetpool&source=all&limit=500", te.NodeName)
		if err := te.Client().DoRequest(ctx, http.MethodGet, path, nil, resp); err != nil {
			return fmt.Errorf("listing cephsetpool tasks: %w", err)
		}

		for _, task := range resp.Data {
			if task.ID == poolName {
				return fmt.Errorf("expected no update of pool %q, found task %s", poolName, task.UPID)
			}
		}

		return nil
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
