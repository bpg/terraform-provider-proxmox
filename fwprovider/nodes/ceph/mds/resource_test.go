//go:build acceptance || all

//testacc:tier=light
//testacc:resource=ceph

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package mds_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceCephMDS(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)
	te.RequireCeph()

	basicName := test.SafeResourceName("ceph-mds-basic")
	renamedName := test.SafeResourceName("ceph-mds-renamed")
	driftName := test.SafeResourceName("ceph-mds-drift")
	firstName := test.SafeResourceName("ceph-mds-pair")
	secondName := test.SafeResourceName("ceph-mds-pair")
	dupName := test.SafeResourceName("ceph-mds-dup")

	tests := []struct {
		name string
		// every MDS the case creates; CheckDestroy requires each to be gone from the list and from ceph.conf
		created []string
		steps   []resource.TestStep
	}{
		{"create, import and rename", []string{basicName, renamedName}, []resource.TestStep{
			{
				Config: mdsConfig(te, "basic", basicName),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_ceph_mds.basic", map[string]string{
						"id":        "^" + basicName + "$",
						"node_name": "^" + te.NodeName + "$",
						"name":      "^" + basicName + "$",
					}),
					checkMDSConfigured(te, basicName),
				),
			},
			{
				ResourceName:      "proxmox_ceph_mds.basic",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     te.NodeName + "/" + basicName,
			},
			{
				// Proxmox can't rename an MDS, so a new name replaces it.
				Config: mdsConfig(te, "basic", renamedName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_ceph_mds.basic", plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_ceph_mds.basic", map[string]string{
						"id":   "^" + renamedName + "$",
						"name": "^" + renamedName + "$",
					}),
					checkMDSConfigured(te, renamedName),
					checkMDSRemoved(te, basicName),
				),
			},
		}},
		{"removed outside Terraform is created again", []string{driftName}, []resource.TestStep{
			{
				Config: mdsConfig(te, "drift", driftName),
				Check:  checkMDSConfigured(te, driftName),
			},
			{
				PreConfig: func() {
					ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
					defer cancel()

					err := te.NodeClient().Ceph().MDS().Delete(ctx, driftName).Err()
					require.NoError(t, err, "out-of-band removal of MDS %q must succeed", driftName)
				},
				Config: mdsConfig(te, "drift", driftName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("proxmox_ceph_mds.drift", plancheck.ResourceActionCreate),
					},
				},
				Check: checkMDSConfigured(te, driftName),
			},
		}},
		{"two in one apply keep both ceph.conf sections", []string{firstName, secondName}, []resource.TestStep{
			{
				// Without the client's lock, the parallel creates or destroys can drop a section.
				Config: te.RenderConfig(fmt.Sprintf(`
					resource "proxmox_ceph_mds" "first" {
						node_name = "{{.NodeName}}"
						name      = %q
					}

					resource "proxmox_ceph_mds" "second" {
						node_name = "{{.NodeName}}"
						name      = %q
					}`, firstName, secondName)),
				Check: checkMDSConfigured(te, firstName, secondName),
			},
		}},
		{"a name already in use fails and leaves the MDS alone", []string{dupName}, []resource.TestStep{
			{
				Config: mdsConfig(te, "original", dupName),
				Check:  checkMDSConfigured(te, dupName),
			},
			{
				Config: te.RenderConfig(fmt.Sprintf(`
					resource "proxmox_ceph_mds" "original" {
						node_name = "{{.NodeName}}"
						name      = %q
					}

					resource "proxmox_ceph_mds" "duplicate" {
						node_name = "{{.NodeName}}"
						name      = %q
					}`, dupName, dupName)),
				ExpectError: regexp.MustCompile(`already referenced in ceph config`),
			},
			{
				Config: mdsConfig(te, "original", dupName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: checkMDSConfigured(te, dupName),
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				CheckDestroy:             checkMDSRemoved(te, tt.created...),
				Steps:                    tt.steps,
			})
		})
	}
}

func TestUnitResourceCephMDS_Validators(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		value string
		error string
	}{
		{"leading digit", "1mds", "must start with a letter"},
		{"underscore", "ceph_mds", "must start with a letter"},
		{"trailing hyphen", "ceph-mds-", "must start with a letter"},
		{"longer than 200 characters", strings.Repeat("a", 201), "at most 200"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps: []resource.TestStep{{
					Config:      mdsConfig(te, "invalid", tt.value),
					PlanOnly:    true,
					ExpectError: regexp.MustCompile(tt.error),
				}},
			})
		})
	}
}

func mdsConfig(te *test.Environment, label, name string) string {
	return te.RenderConfig(fmt.Sprintf(`
		resource "proxmox_ceph_mds" %q {
			node_name = "{{.NodeName}}"
			name      = %q
		}`, label, name))
}

// checkMDSConfigured checks with the API that each MDS is listed on the test node and has its own ceph.conf section.
func checkMDSConfigured(te *test.Environment, names ...string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		hosts, conf, err := mdsInventory(ctx, te)
		if err != nil {
			return err
		}

		for _, name := range names {
			host, ok := hosts[name]
			if !ok {
				return fmt.Errorf("MDS %q is not listed", name)
			}

			if host != te.NodeName {
				return fmt.Errorf("MDS %q is listed on host %q, expected %q", name, host, te.NodeName)
			}

			if !strings.Contains(conf, "[mds."+name+"]") {
				return fmt.Errorf("ceph.conf has no [mds.%s] section", name)
			}
		}

		return nil
	}
}

// checkMDSRemoved checks with the API that no MDS is listed under any of the names and none has a ceph.conf section.
func checkMDSRemoved(te *test.Environment, names ...string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		hosts, conf, err := mdsInventory(ctx, te)
		if err != nil {
			return err
		}

		for _, name := range names {
			if host, ok := hosts[name]; ok {
				return fmt.Errorf("MDS %q is still listed, on host %q", name, host)
			}

			if strings.Contains(conf, "[mds."+name+"]") {
				return fmt.Errorf("ceph.conf still has an [mds.%s] section", name)
			}
		}

		return nil
	}
}

// mdsInventory returns the cluster's MDS IDs mapped to their hosts, and the raw ceph.conf.
func mdsInventory(ctx context.Context, te *test.Environment) (map[string]string, string, error) {
	list, err := te.NodeClient().Ceph().MDS().List(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("listing MDS: %w", err)
	}

	hosts := make(map[string]string, len(list))

	for _, d := range list {
		host := ""
		if d.Host != nil {
			host = *d.Host
		}

		hosts[d.Name] = host
	}

	resBody := struct {
		Data *string `json:"data"`
	}{}

	path := fmt.Sprintf("nodes/%s/ceph/cfg/raw", url.PathEscape(te.NodeName))
	if err := te.Client().DoRequest(ctx, http.MethodGet, path, nil, &resBody); err != nil {
		return nil, "", fmt.Errorf("reading ceph.conf: %w", err)
	}

	if resBody.Data == nil {
		return nil, "", fmt.Errorf("reading ceph.conf: empty response")
	}

	return hosts, *resBody.Data, nil
}
