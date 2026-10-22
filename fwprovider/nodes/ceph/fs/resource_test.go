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

// The test creates its own MDS, so the cluster needs none, and removes the filesystem by hand the way the admin guide
// describes, since the resource leaves that to the user. PROXMOX_VE_ACC_CEPH_CRUSH_RULE, when set, places the pools
// on a CRUSH rule that can hold two replicas.
func TestAccResourceCephFS(t *testing.T) {
	te := test.InitEnvironment(t)
	te.RequireCeph()
	te.RequireRootSSH()
	te.CreateCephMDS()

	name := test.SafeResourceName("ceph-fs-res")
	rule := utils.GetAnyStringEnv("PROXMOX_VE_ACC_CEPH_CRUSH_RULE")

	// In case a step fails before the filesystem is removed by hand.
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		if err := te.RemoveCephFS(ctx, name); err != nil {
			t.Errorf("removing test CephFS: %v", err)
		}
	})

	te.AddTemplateVars(map[string]any{"FSName": name, "FSCrushRule": rule})

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		CheckDestroy:             checkFSRemoved(te, name),
		Steps: []resource.TestStep{
			{
				Config: fsConfig(te, 2, 1),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_ceph_fs.test", map[string]string{
						"id":                     "^" + name + "$",
						"name":                   "^" + name + "$",
						"node_name":              "^" + te.NodeName + "$",
						"pg_num":                 "^8$",
						"metadata_pool.name":     "^" + name + "_metadata$",
						"metadata_pool.size":     "^2$",
						"metadata_pool.min_size": "^1$",
						"data_pool.name":         "^" + name + "_data$",
						"data_pool.size":         "^2$",
						"data_pool.min_size":     "^1$",
					}),
					test.ResourceAttributesSet("proxmox_ceph_fs.test", []string{
						"metadata_pool.crush_rule",
						"metadata_pool.pool_id",
						"data_pool.crush_rule",
						"data_pool.pool_id",
					}),
					checkFSExists(te, name),
					checkPool(te, name+"_metadata", 2, 1, rule),
					checkPool(te, name+"_data", 2, 1, rule),
				),
			},
			{
				// Size goes first, so the min_size sent with it isn't lost (Proxmox bug 8193).
				Config: fsConfig(te, 3, 2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_ceph_fs.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_ceph_fs.test", map[string]string{
						"data_pool.size":     "^3$",
						"data_pool.min_size": "^2$",
						"metadata_pool.size": "^2$",
					}),
					checkPool(te, name+"_data", 3, 2, rule),
					checkPool(te, name+"_metadata", 2, 1, rule),
				),
			},
			{
				ResourceName:      "proxmox_ceph_fs.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     te.NodeName + "/" + name,
				// Only used at creation, so an import leaves it unset.
				ImportStateVerifyIgnore: []string{"pg_num"},
			},
			{
				// Removed from the configuration while it still exists: the resource won't destroy it.
				Config:      te.RenderConfig(""),
				ExpectError: regexp.MustCompile(`(?s)Unable to Delete CephFS.*manual procedure.*pveceph fs destroy`),
			},
			{
				// Once it has been removed by hand, the destroy goes through.
				PreConfig: func() {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
					defer cancel()

					require.NoError(t, te.RemoveCephFS(ctx, name))
				},
				Config: te.RenderConfig(""),
				Check:  checkFSRemoved(te, name),
			},
		},
	})
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

// fsConfig renders the test filesystem with the given data pool size and min_size; the metadata pool stays at 2/1.
func fsConfig(te *test.Environment, dataSize, dataMinSize int) string {
	return te.RenderConfig(fmt.Sprintf(`
		resource "proxmox_ceph_fs" "test" {
			node_name = "{{.NodeName}}"
			name      = "{{.FSName}}"
			pg_num    = 8

			metadata_pool = {
				size     = 2
				min_size = 1
				{{with .FSCrushRule}}crush_rule = "{{.}}"{{end}}
			}

			data_pool = {
				size     = %d
				min_size = %d
				{{with .FSCrushRule}}crush_rule = "{{.}}"{{end}}
			}
		}`, dataSize, dataMinSize))
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
