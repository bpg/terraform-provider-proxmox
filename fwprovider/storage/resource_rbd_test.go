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
	"slices"
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
)

func TestAccResourceStorageRBD(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)
	te.RequireCeph()

	basicName := test.SafeResourceName("ceph-rbd-basic")
	namespaceName := test.SafeResourceName("ceph-rbd-ns")
	ecName := test.SafeResourceName("ceph-rbd-ec")
	invalidName := test.SafeResourceName("ceph-rbd-invalid")

	tests := []struct {
		name         string
		checkDestroy resource.TestCheckFunc
		step         []resource.TestStep
	}{
		{"create with its pool, update and import", checkRBDDestroyed(te, basicName, basicName), []resource.TestStep{
			{
				// A single apply builds the pool and the storage on it; the destroy removes them in reverse.
				Config: te.RenderConfig(rbdPoolConfig(basicName) + fmt.Sprintf(`
					resource "proxmox_storage_rbd" "test" {
						id   = %q
						pool = proxmox_ceph_pool.test.name
					}`, basicName)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_storage_rbd.test", map[string]string{
						"id":        basicName,
						"pool":      basicName,
						"content.#": "1",
						"krbd":      "false",
						"disable":   "false",
						"shared":    "true",
					}),
					resource.TestCheckTypeSetElemAttr("proxmox_storage_rbd.test", "content.*", "images"),
					checkStorageConfig(te, basicName, map[string]any{
						"type":      "rbd",
						"pool":      basicName,
						"data-pool": nil,
						"namespace": nil,
					}),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("proxmox_storage_rbd.test", tfjsonpath.New("data_pool"), knownvalue.Null()),
					statecheck.ExpectKnownValue("proxmox_storage_rbd.test", tfjsonpath.New("namespace"), knownvalue.Null()),
				},
			},
			{
				Config: te.RenderConfig(rbdPoolConfig(basicName) + fmt.Sprintf(`
					resource "proxmox_storage_rbd" "test" {
						id      = %q
						pool    = proxmox_ceph_pool.test.name
						nodes   = ["{{.NodeName}}"]
						content = ["images", "rootdir"]
						krbd    = true
						disable = true
					}`, basicName)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_storage_rbd.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("proxmox_ceph_pool.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_storage_rbd.test", map[string]string{
						"nodes.#":   "1",
						"content.#": "2",
						"krbd":      "true",
						"disable":   "true",
					}),
					resource.TestCheckTypeSetElemAttr("proxmox_storage_rbd.test", "content.*", "rootdir"),
					checkStorageConfig(te, basicName, map[string]any{
						"nodes":   te.NodeName,
						"krbd":    1,
						"disable": 1,
					}),
				),
			},
			{
				ResourceName:      "proxmox_storage_rbd.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     basicName,
			},
		}},
		{"namespace changes replace the storage", checkRBDDestroyed(te, namespaceName, namespaceName), []resource.TestStep{
			{
				// Proxmox doesn't need the namespace to exist until it allocates an image.
				Config: te.RenderConfig(rbdPoolConfig(namespaceName) + rbdNamespaceConfig(namespaceName, `namespace = "tenant-a"`)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_storage_rbd.test", "namespace", "tenant-a"),
					checkStorageConfig(te, namespaceName, map[string]any{"namespace": "tenant-a"}),
				),
			},
			{
				Config: te.RenderConfig(rbdPoolConfig(namespaceName) + rbdNamespaceConfig(namespaceName, `namespace = "tenant-b"`)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_storage_rbd.test", plancheck.ResourceActionReplace),
					},
				},
				Check: checkStorageConfig(te, namespaceName, map[string]any{"namespace": "tenant-b"}),
			},
			{
				Config: te.RenderConfig(rbdPoolConfig(namespaceName) + rbdNamespaceConfig(namespaceName, "")),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_storage_rbd.test", plancheck.ResourceActionReplace),
					},
				},
				Check: checkStorageConfig(te, namespaceName, map[string]any{"namespace": nil}),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("proxmox_storage_rbd.test", tfjsonpath.New("namespace"), knownvalue.Null()),
				},
			},
			{
				ResourceName:      "proxmox_storage_rbd.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     namespaceName,
			},
		}},
		{"erasure-coded pool", checkRBDDestroyed(te, ecName, ecName+"-data", ecName+"-metadata"), []resource.TestStep{
			{
				Config: te.RenderConfig(rbdECPoolConfig(ecName) + fmt.Sprintf(`
					resource "proxmox_storage_rbd" "test" {
						id        = %q
						pool      = proxmox_ceph_pool.test.metadata_pool_name
						data_pool = proxmox_ceph_pool.test.data_pool.name
					}`, ecName)),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_storage_rbd.test", map[string]string{
						"pool":      ecName + "-metadata",
						"data_pool": ecName + "-data",
					}),
					checkStorageConfig(te, ecName, map[string]any{
						"pool":      ecName + "-metadata",
						"data-pool": ecName + "-data",
					}),
				),
			},
			{
				ResourceName:      "proxmox_storage_rbd.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     ecName,
			},
		}},
		{"invalid content type", nil, []resource.TestStep{
			{
				Config: te.RenderConfig(fmt.Sprintf(`
					resource "proxmox_storage_rbd" "test" {
						id      = %q
						pool    = "rbd"
						content = ["iso"]
					}`, invalidName)),
				ExpectError: regexp.MustCompile(`(?s)content.*one of`),
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				CheckDestroy:             tt.checkDestroy,
				Steps:                    tt.step,
			})
		})
	}
}

// rbdPoolConfig renders a small replicated Ceph pool for a test storage.
func rbdPoolConfig(name string) string {
	return fmt.Sprintf(`
		resource "proxmox_ceph_pool" "test" {
			node_name         = "{{.NodeName}}"
			name              = %q
			size              = 2
			min_size          = 1
			pg_num            = 8
			pg_autoscale_mode = "warn"
			{{with .CephCrushRule}}crush_rule = "{{.}}"{{end}}
		}`, name)
}

// rbdECPoolConfig renders a small erasure-coded Ceph pool (k=2, m=2 over OSDs) for a test storage.
func rbdECPoolConfig(name string) string {
	return fmt.Sprintf(`
		resource "proxmox_ceph_pool" "test" {
			node_name = "{{.NodeName}}"
			name      = %q
			size      = 2
			min_size  = 1
			pg_num    = 8
			{{with .CephCrushRule}}crush_rule = "{{.}}"{{end}}
			erasure_coding = {
				k              = 2
				m              = 2
				failure_domain = "osd"
				{{with .CephECDeviceClass}}device_class = "{{.}}"{{end}}
			}
		}`, name)
}

// rbdNamespaceConfig renders a storage on the test pool with extra attributes appended.
func rbdNamespaceConfig(name, extra string) string {
	return fmt.Sprintf(`
		resource "proxmox_storage_rbd" "test" {
			id   = %q
			pool = proxmox_ceph_pool.test.name
			%s
		}`, name, extra)
}

// checkStorageConfig compares a storage's configuration as Proxmox returns it with the wanted values; a nil value
// means the key must be absent.
func checkStorageConfig(te *test.Environment, storageID string, want map[string]any) resource.TestCheckFunc {
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

// checkRBDDestroyed fails if the storage or any of the Ceph pools still exists.
func checkRBDDestroyed(te *test.Environment, storageID string, pools ...string) resource.TestCheckFunc {
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

		resp := &struct {
			Data []struct {
				Name string `json:"pool_name"`
			} `json:"data"`
		}{}
		if err := te.Client().DoRequest(ctx, http.MethodGet, fmt.Sprintf("nodes/%s/ceph/pool", te.NodeName), nil, resp); err != nil {
			return fmt.Errorf("listing Ceph pools: %w", err)
		}

		for _, p := range resp.Data {
			if slices.Contains(pools, p.Name) {
				return fmt.Errorf("expected Ceph pool %q to be destroyed, but it still exists", p.Name)
			}
		}

		return nil
	}
}
