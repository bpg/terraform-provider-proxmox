//go:build acceptance || all

//testacc:tier=light
//testacc:resource=ceph

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	fsapi "github.com/bpg/terraform-provider-proxmox/proxmox/nodes/ceph/fs"
)

// The storages point at a CephFS the test creates. Every one is disabled, so Proxmox never mounts it: an enabled
// storage would leave its mount behind after destroy. "ro" keeps a mount read-only should one happen anyway.
func TestAccResourceStorageCephFS(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)
	te.RequireCeph()

	fsName := newTestCephFS(t, te)

	basicID := test.SafeResourceName("ceph-fs-basic")
	replaceID := test.SafeResourceName("ceph-fs-replace")

	tests := []struct {
		name  string
		id    string
		steps []resource.TestStep
	}{
		{"create, update and import", basicID, []resource.TestStep{
			{
				Config: cephFSStorageConfig(te, basicID, fsName, ""),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_storage_cephfs.test", map[string]string{
						"id":               basicID,
						"fs_name":          fsName,
						"options":          "ro",
						"disable":          "true",
						"fuse":             "false",
						"shared":           "true",
						"create_base_path": "true",
						"create_subdirs":   "true",
						"content.#":        "1",
					}),
					resource.TestCheckTypeSetElemAttr("proxmox_storage_cephfs.test", "content.*", "backup"),
					checkCephFSConfig(te, basicID, map[string]any{
						"type":    "cephfs",
						"fs-name": fsName,
						"path":    "/mnt/pve/" + basicID,
						"options": "ro",
						"disable": 1,
						"content": "backup",
						"subdir":  nil,
					}),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("proxmox_storage_cephfs.test", tfjsonpath.New("subdir"), knownvalue.Null()),
				},
			},
			{
				Config: cephFSStorageConfig(te, basicID, fsName, `
					nodes   = ["{{.NodeName}}"]
					content = ["backup", "import", "iso", "snippets", "vztmpl"]
					options = "ro,noatime"

					backups {
						keep_last = 2
					}`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_storage_cephfs.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_storage_cephfs.test", map[string]string{
						"nodes.#":           "1",
						"content.#":         "5",
						"options":           "ro,noatime",
						"backups.keep_last": "2",
						"backups.keep_all":  "false",
						"fs_name":           fsName,
						"create_subdirs":    "true",
						"create_base_path":  "true",
						"disable":           "true",
						"shared":            "true",
						"fuse":              "false",
					}),
					resource.TestCheckTypeSetElemAttr("proxmox_storage_cephfs.test", "content.*", "snippets"),
					resource.TestCheckTypeSetElemAttr("proxmox_storage_cephfs.test", "content.*", "import"),
					checkCephFSConfig(te, basicID, map[string]any{
						"nodes":         te.NodeName,
						"options":       "ro,noatime",
						"prune-backups": "keep-last=2",
					}),
					checkCephFSContent(te, basicID, "backup", "import", "iso", "snippets", "vztmpl"),
				),
			},
			{
				ResourceName:      "proxmox_storage_cephfs.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     basicID,
				// As on every storage resource, the base reads `backups` back only once the configuration has the block.
				ImportStateVerifyIgnore: []string{
					"backups",
					"backups.%",
					"backups.keep_all",
					"backups.keep_last",
				},
			},
		}},
		{"fs_name, subdir and fuse force replacement", replaceID, []resource.TestStep{
			{
				// Proxmox doesn't check the subdirectory until it mounts the storage.
				Config: cephFSStorageConfig(te, replaceID, fsName, `subdir = "/tfacc-a"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_storage_cephfs.test", "subdir", "/tfacc-a"),
					checkCephFSConfig(te, replaceID, map[string]any{"subdir": "/tfacc-a", "fs-name": fsName}),
				),
			},
			{
				Config: cephFSStorageConfig(te, replaceID, fsName, `subdir = "/tfacc-b"`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_storage_cephfs.test", plancheck.ResourceActionReplace),
					},
				},
				Check: checkCephFSConfig(te, replaceID, map[string]any{"subdir": "/tfacc-b"}),
			},
			{
				Config: cephFSStorageConfig(te, replaceID, fsName, `
					subdir = "/tfacc-b"
					fuse   = true`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_storage_cephfs.test", plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_storage_cephfs.test", "fuse", "true"),
					checkCephFSConfig(te, replaceID, map[string]any{"fuse": 1, "subdir": "/tfacc-b"}),
				),
			},
			{
				// Proxmox accepts a filesystem name it can't find; it only fails once it mounts the storage.
				Config: cephFSStorageConfig(te, replaceID, "tfaccmissing", `
					subdir = "/tfacc-b"
					fuse   = true`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_storage_cephfs.test", plancheck.ResourceActionReplace),
					},
				},
				Check: checkCephFSConfig(te, replaceID, map[string]any{"fs-name": "tfaccmissing"}),
			},
			{
				ResourceName:      "proxmox_storage_cephfs.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     replaceID,
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				CheckDestroy:             checkCephFSDestroyed(te, tt.id),
				Steps:                    tt.steps,
			})
		})
	}
}

func TestUnitResourceStorageCephFS_Validators(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		extra string
		error string
	}{
		{"content not allowed on CephFS", `content = ["images"]`, `(?s)content.*one of`},
		{"container root directories", `content = ["rootdir"]`, `(?s)content.*one of`},
		{"empty subdir", `subdir = ""`, `(?s)subdir.*at least 1`},
		{"empty fs_name", `fs_name = ""`, `(?s)fs_name.*at least 1`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps: []resource.TestStep{{
					Config: te.RenderConfig(fmt.Sprintf(`
						resource "proxmox_storage_cephfs" "test" {
							id      = "ceph-fs-invalid"
							disable = true
							%s
						}`, tt.extra)),
					PlanOnly:    true,
					ExpectError: regexp.MustCompile(tt.error),
				}},
			})
		})
	}
}

// cephFSStorageConfig renders a disabled, read-only test storage on the filesystem, with extra attributes appended.
// An `options` in extra replaces the default `ro`.
func cephFSStorageConfig(te *test.Environment, id, fsName, extra string) string {
	options := `options = "ro"`
	if strings.Contains(extra, "options") {
		options = ""
	}

	return te.RenderConfig(fmt.Sprintf(`
		resource "proxmox_storage_cephfs" "test" {
			id      = %q
			fs_name = %q
			disable = true
			%s
			%s
		}`, id, fsName, options, extra))
}

// newTestCephFS creates a CephFS for the test, with its own MDS so the cluster needs none, and removes both once the
// test and its subtests have finished: the filesystem by hand as the admin guide describes, then the MDS. The
// filesystem is never mounted, so its pools' placement doesn't matter.
func newTestCephFS(t *testing.T, te *test.Environment) string {
	t.Helper()

	te.RequireRootSSH()
	te.CreateCephMDS()

	name := test.SafeResourceName("ceph-fs-test")

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		if err := te.RemoveCephFS(ctx, name); err != nil {
			t.Errorf("removing test CephFS: %v", err)
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	pgNum := int64(8)
	if err := te.NodeClient().Ceph().FS().Create(ctx, name, &fsapi.CreateRequestBody{PGNum: &pgNum}).Err(); err != nil {
		t.Fatalf("creating test CephFS %q: %v", name, err)
	}

	return name
}

// checkCephFSConfig compares a storage's configuration as Proxmox returns it with the wanted values; a nil value
// means the key must be absent.
func checkCephFSConfig(te *test.Environment, storageID string, want map[string]any) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		resp := &struct {
			Data map[string]any `json:"data"`
		}{}
		if err := te.Client().DoRequest(ctx, http.MethodGet, "storage/"+storageID, nil, resp); err != nil {
			return fmt.Errorf("reading storage %q: %w", storageID, err)
		}

		for key, wantValue := range want {
			gotValue, ok := resp.Data[key]

			switch {
			case wantValue == nil && ok:
				return fmt.Errorf("storage %q: expected no %q, got %v", storageID, key, gotValue)
			case wantValue == nil:
				continue
			case !ok:
				return fmt.Errorf("storage %q: expected %q = %v, but it is not set", storageID, key, wantValue)
			case fmt.Sprint(gotValue) != fmt.Sprint(wantValue):
				return fmt.Errorf("storage %q: expected %q = %v, got %v", storageID, key, wantValue, gotValue)
			}
		}

		return nil
	}
}

// checkCephFSContent checks with the API that the storage allows exactly the given content types.
func checkCephFSContent(te *test.Environment, storageID string, want ...string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		resp := &struct {
			Data struct {
				Content string `json:"content"`
			} `json:"data"`
		}{}
		if err := te.Client().DoRequest(ctx, http.MethodGet, "storage/"+storageID, nil, resp); err != nil {
			return fmt.Errorf("reading storage %q: %w", storageID, err)
		}

		got := strings.Split(resp.Data.Content, ",")
		if len(got) != len(want) {
			return fmt.Errorf("storage %q: expected content %v, got %q", storageID, want, resp.Data.Content)
		}

		for _, w := range want {
			if !strings.Contains(","+resp.Data.Content+",", ","+w+",") {
				return fmt.Errorf("storage %q: expected content %v, got %q", storageID, want, resp.Data.Content)
			}
		}

		return nil
	}
}

// checkCephFSDestroyed fails if the storage still exists.
func checkCephFSDestroyed(te *test.Environment, storageID string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		err := te.Client().DoRequest(ctx, http.MethodGet, "storage/"+storageID, nil, &struct{}{})
		if err == nil {
			return fmt.Errorf("expected storage %q to be destroyed, but it still exists", storageID)
		}

		if !errors.Is(err, api.ErrResourceDoesNotExist) {
			return fmt.Errorf("reading storage %q: %w", storageID, err)
		}

		return nil
	}
}
