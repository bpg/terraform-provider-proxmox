//go:build acceptance || all

//testacc:tier=light
//testacc:resource=ceph

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fs_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

// Each case creates its own MDS, since a CephFS needs one in standby, and removes its filesystem by hand the way the
// admin guide describes, since the resource leaves that to the user. PROXMOX_VE_ACC_CEPH_CRUSH_RULE, when set, places
// the pools on a CRUSH rule that can hold their replicas; the pool-defaults case moves its pools there from the
// cluster's default rule.
func TestAccResourceCephFS(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)
	te.RequireCeph()

	rule := utils.GetAnyStringEnv("PROXMOX_VE_ACC_CEPH_CRUSH_RULE")
	updateName := test.SafeResourceName("ceph-fs-update")
	defaultsName := test.SafeResourceName("ceph-fs-defaults")

	// Every case ends by removing its filesystem by hand and destroying the resource.
	tests := []struct {
		name  string
		fs    string
		steps []resource.TestStep
	}{
		{"create, update and import", updateName, []resource.TestStep{
			{
				Config: fsConfig(te, updateName, poolBlock("metadata_pool", 2, 1)+poolBlock("data_pool", 2, 1)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_ceph_fs.test", map[string]string{
						"id":                     "^" + updateName + "$",
						"name":                   "^" + updateName + "$",
						"node_name":              "^" + te.NodeName + "$",
						"pg_num":                 "^8$",
						"metadata_pool.name":     "^" + updateName + "_metadata$",
						"metadata_pool.size":     "^2$",
						"metadata_pool.min_size": "^1$",
						"data_pool.name":         "^" + updateName + "_data$",
						"data_pool.size":         "^2$",
						"data_pool.min_size":     "^1$",
					}),
					test.ResourceAttributesSet("proxmox_ceph_fs.test", []string{
						"metadata_pool.crush_rule",
						"metadata_pool.pool_id",
						"data_pool.crush_rule",
						"data_pool.pool_id",
					}),
					checkFSExists(te, updateName),
					checkPool(te, updateName+"_metadata", 2, 1, rule),
					checkPool(te, updateName+"_data", 2, 1, rule),
				),
			},
			{
				// Raising size makes Ceph reset min_size to 2, and Proxmox would skip the min_size of 1 sent with it as
				// unchanged (Proxmox bug 8193), so size has to go first.
				Config: fsConfig(te, updateName, poolBlock("metadata_pool", 2, 1)+poolBlock("data_pool", 3, 1)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_ceph_fs.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_ceph_fs.test", map[string]string{
						"data_pool.size":     "^3$",
						"data_pool.min_size": "^1$",
					}),
					checkPool(te, updateName+"_data", 3, 1, rule),
					checkPool(te, updateName+"_metadata", 2, 1, rule),
				),
			},
			{
				Config: fsConfig(te, updateName, poolBlock("metadata_pool", 2, 2)+poolBlock("data_pool", 3, 1)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_ceph_fs.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_ceph_fs.test", map[string]string{
						"metadata_pool.size":     "^2$",
						"metadata_pool.min_size": "^2$",
					}),
					checkPool(te, updateName+"_metadata", 2, 2, rule),
					checkPool(te, updateName+"_data", 3, 1, rule),
				),
			},
			{
				ResourceName:      "proxmox_ceph_fs.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     te.NodeName + "/" + updateName,
				// Only used at creation, so an import leaves it unset.
				ImportStateVerifyIgnore: []string{"pg_num"},
			},
			{
				// A new pg_num replaces the filesystem, and the resource won't destroy the old one.
				Config: fsConfig(te, updateName, `pg_num = 16`+poolBlock("metadata_pool", 2, 2)+
					poolBlock("data_pool", 3, 1)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_ceph_fs.test", plancheck.ResourceActionReplace),
					},
				},
				ExpectError: regexp.MustCompile(`(?s)Unable to Delete CephFS.*manual procedure`),
			},
			{
				// Removed from the configuration while it still exists: the resource won't destroy it.
				Config:      te.RenderConfig(""),
				ExpectError: regexp.MustCompile(`(?s)Unable to Delete CephFS.*manual procedure.*pveceph fs destroy`),
			},
		}},
		{"pool blocks left out, then set", defaultsName, []resource.TestStep{
			{
				// Proxmox creates the pools under the cluster's defaults, which are read back.
				Config: fsConfig(te, defaultsName, ""),
				Check: resource.ComposeTestCheckFunc(
					checkPoolReadBack(te, "metadata_pool", defaultsName+"_metadata"),
					checkPoolReadBack(te, "data_pool", defaultsName+"_data"),
				),
			},
			{
				// With a test rule set, this moves both pools from the default rule onto it.
				Config: fsConfig(te, defaultsName, poolBlock("metadata_pool", 2, 1)+poolBlock("data_pool", 2, 1)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_ceph_fs.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					checkPool(te, defaultsName+"_metadata", 2, 1, rule),
					checkPool(te, defaultsName+"_data", 2, 1, rule),
					checkPoolReadBack(te, "metadata_pool", defaultsName+"_metadata"),
					checkPoolReadBack(te, "data_pool", defaultsName+"_data"),
				),
			},
		}},
	}

	for range tests {
		te.CreateCephMDS()
	}

	for _, tt := range tests {
		// In case a step fails before the filesystem is removed by hand.
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			if err := te.RemoveCephFS(ctx, tt.fs); err != nil {
				t.Errorf("removing test CephFS %q: %v", tt.fs, err)
			}
		})

		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				CheckDestroy:             checkFSRemoved(te, tt.fs),
				Steps:                    append(tt.steps, removedByHandStep(t, te, tt.fs)),
			})
		})
	}
}

func TestUnitResourceCephFS_Validators(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		attrs string
		error string
	}{
		{"name with a slash", `name = "bad/name"`, `must not contain`},
		{"name with whitespace", `name = "bad name"`, `must not contain`},
		{"pg_num too small", "name = \"tfacc\"\npg_num = 4", `(?s)pg_num.*between 8 and 32768`},
		{"pool size zero", "name = \"tfacc\"\ndata_pool = { size = 0 }", `(?s)size.*between 1 and 7`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps: []resource.TestStep{{
					Config: te.RenderConfig(fmt.Sprintf(`
						resource "proxmox_ceph_fs" "test" {
							node_name = "{{.NodeName}}"
							%s
						}`, tt.attrs)),
					PlanOnly:    true,
					ExpectError: regexp.MustCompile(tt.error),
				}},
			})
		})
	}
}

// fsConfig renders a test filesystem with 8 placement groups, or the pg_num in extra, and extra's pool blocks.
func fsConfig(te *test.Environment, name, extra string) string {
	pgNum := "pg_num = 8"
	if strings.Contains(extra, "pg_num") {
		pgNum = ""
	}

	return te.RenderConfig(fmt.Sprintf(`
		resource "proxmox_ceph_fs" "test" {
			node_name = "{{.NodeName}}"
			name      = %q
			%s
			%s
		}`, name, pgNum, extra))
}

// poolBlock renders a pool block with the given size and min_size, on the test CRUSH rule when one is set.
func poolBlock(attr string, size, minSize int) string {
	return fmt.Sprintf(`
		%s = {
			size     = %d
			min_size = %d
			{{with .CephCrushRule}}crush_rule = "{{.}}"{{end}}
		}`, attr, size, minSize)
}

// removedByHandStep removes the filesystem by hand, then checks that destroying the resource goes through.
func removedByHandStep(t *testing.T, te *test.Environment, name string) resource.TestStep {
	t.Helper()

	return resource.TestStep{
		PreConfig: func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			require.NoError(t, te.RemoveCephFS(ctx, name))
		},
		Config: te.RenderConfig(""),
		Check:  checkFSRemoved(te, name),
	}
}

// checkFSExists checks with the API that the filesystem is listed.
func checkFSExists(te *test.Environment, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		if _, err := te.NodeClient().Ceph().FS().Get(ctx, name); err != nil {
			return fmt.Errorf("CephFS %q: %w", name, err)
		}

		return nil
	}
}

// checkFSRemoved checks with the API that neither the filesystem nor its pools exist.
func checkFSRemoved(te *test.Environment, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		_, err := te.NodeClient().Ceph().FS().Get(ctx, name)
		if err == nil {
			return fmt.Errorf("expected CephFS %q to be gone, but it still exists", name)
		}

		if !errors.Is(err, api.ErrResourceDoesNotExist) {
			return fmt.Errorf("checking CephFS %q: %w", name, err)
		}

		for _, pool := range []string{name + "_metadata", name + "_data"} {
			_, err := te.NodeClient().Ceph().Pool().Get(ctx, pool)
			if err == nil {
				return fmt.Errorf("expected Ceph pool %q to be gone, but it still exists", pool)
			}

			if !errors.Is(err, api.ErrResourceDoesNotExist) {
				return fmt.Errorf("checking Ceph pool %q: %w", pool, err)
			}
		}

		return nil
	}
}

// checkPoolReadBack checks that the state holds the pool's settings as the API reports them.
func checkPoolReadBack(te *test.Environment, attr, pool string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		status, err := te.NodeClient().Ceph().Pool().Get(ctx, pool)
		if err != nil {
			return fmt.Errorf("reading Ceph pool %q: %w", pool, err)
		}

		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("proxmox_ceph_fs.test", attr+".name", pool),
			resource.TestCheckResourceAttr("proxmox_ceph_fs.test", attr+".crush_rule", status.CrushRule),
			resource.TestCheckResourceAttr("proxmox_ceph_fs.test", attr+".size", strconv.FormatInt(status.Size, 10)),
			resource.TestCheckResourceAttr("proxmox_ceph_fs.test", attr+".min_size", strconv.FormatInt(status.MinSize, 10)),
		)(s)
	}
}

// checkPool checks a pool's replication settings with the API; an empty rule isn't checked.
func checkPool(te *test.Environment, pool string, size, minSize int64, rule string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		status, err := te.NodeClient().Ceph().Pool().Get(ctx, pool)
		if err != nil {
			return fmt.Errorf("reading Ceph pool %q: %w", pool, err)
		}

		if status.Size != size || status.MinSize != minSize {
			return fmt.Errorf("Ceph pool %q: expected size %d/min_size %d, got %d/%d", pool, size, minSize,
				status.Size, status.MinSize)
		}

		if rule != "" && status.CrushRule != rule {
			return fmt.Errorf("Ceph pool %q: expected CRUSH rule %q, got %q", pool, rule, status.CrushRule)
		}

		return nil
	}
}
