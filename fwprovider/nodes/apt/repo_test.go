//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package apt_test

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/apt"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	apitypes "github.com/bpg/terraform-provider-proxmox/proxmox/types/nodes/apt/repositories"
)

// Note that some "hard-coded" values must be used because of the way how the Proxmox VE API for APT repositories works.
const (
	testAccResourceRepoSelector = "proxmox_virtual_environment_" + apt.ResourceRepoIDPrefix + ".test"

	// By default, this should be the main Debian package repository on any (new) Proxmox VE node.
	testAccResourceRepoIndex = 0

	testAccResourceStandardRepoSelector = "proxmox_virtual_environment_" + apt.ResourceStandardRepoIDPrefix + ".test"

	// Use an APT standard repository handle that is not enabled by default on any new Proxmox VE node.
	testAccResourceStandardRepoHandle = "no-subscription"
)

func TestAccDataSourceRepo(t *testing.T) {
	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{
			"read APT repository attributes",
			[]resource.TestStep{
				{
					Config: te.RenderConfig(`
					data "proxmox_virtual_environment_apt_repository" "test" {
						file_path = "/etc/apt/sources.list"
						index = 0
						node = "{{.NodeName}}"
					}`),
					// The provided attributes and computed attributes should be set.
					Check: resource.ComposeTestCheckFunc(
						resource.TestMatchResourceAttr(
							"data.proxmox_virtual_environment_apt_repository.test",
							apt.SchemaAttrNameComment,
							// Expect any value or an empty string.
							regexp.MustCompile(`(.*|^$)`),
						),
						resource.TestCheckResourceAttr(
							"data.proxmox_virtual_environment_apt_repository.test",
							apt.SchemaAttrNameTerraformID,
							"apt_repository_"+strings.ToLower(te.NodeName)+"_etc_apt_sources_list_0",
						),
						test.ResourceAttributesSet("data.proxmox_virtual_environment_apt_repository.test", []string{
							"components.#",
							"enabled",
							"file_path",
							"index",
							"node",
							"package_types.#",
							"suites.#",
							"uris.#",
						}),
					),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				resource.ParallelTest(
					t, resource.TestCase{
						ProtoV6ProviderFactories: te.AccProviders,
						Steps:                    tt.steps,
					},
				)
			},
		)
	}
}

func TestAccDataSourceStandardRepo(t *testing.T) {
	t.Helper()
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{
			"read APT standard repository attributes",
			[]resource.TestStep{
				{
					Config: te.RenderConfig(`
					data "proxmox_virtual_environment_apt_standard_repository" "test" {
						handle = "no-subscription"
						node   = "{{.NodeName}}"
					}`),
					// The provided attributes and computed attributes should be set.
					Check: resource.ComposeTestCheckFunc(
						test.ResourceAttributes("data.proxmox_virtual_environment_apt_standard_repository.test", map[string]string{
							"id": fmt.Sprintf("apt_standard_repository_%s_no_subscription", strings.ToLower(te.NodeName)),
						}),
						test.ResourceAttributesSet("data.proxmox_virtual_environment_apt_standard_repository.test", []string{
							// Note that we can not check for the following attributes because they are only available when the
							// standard repository has been added to a source list:
							//
							// - apt.SchemaAttrNameFilePath (file_path) - will be set when parsing all configured repositories in all
							//   source list files.
							// - apt.SchemaAttrNameIndex (index) - will be set when finding the repository within a source list file,
							//   based on the detected file path.
							// - apt.SchemaAttrNameStandardStatus (status) - is only available when the standard has been configured.

							apt.SchemaAttrNameStandardDescription,
							apt.SchemaAttrNameStandardHandle,
							apt.SchemaAttrNameStandardName,
							apt.SchemaAttrNameNode,
						}),
					),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				resource.ParallelTest(
					t, resource.TestCase{
						ProtoV6ProviderFactories: te.AccProviders,
						Steps:                    tt.steps,
					},
				)
			},
		)
	}
}

// Run tests for APT repository resource definitions with valid input where all required attributes are specified.
// Only the [Create], [Read] and [Update] method implementations of the
// [github.com/hashicorp/terraform-plugin-framework/resource.Resource] interface are tested in sequential steps because
// [Delete] is no-op due to the non-existing capability of the Proxmox VE API of deleting configured APT repository.
//
// [Create]: https://developer.hashicorp.com/terraform/plugin/framework/resources/create
// [Delete]: https://developer.hashicorp.com/terraform/plugin/framework/resources/delete
// [Read]: https://developer.hashicorp.com/terraform/plugin/framework/resources/read
// [Update]: https://developer.hashicorp.com/terraform/plugin/framework/resources/update
func TestAccResourceRepoValidInput(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	resource.Test(
		t, resource.TestCase{
			ProtoV6ProviderFactories: te.AccProviders,
			Steps: []resource.TestStep{
				// Test the "Create" and "Read" implementations.
				{
					Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_apt_repository" "test" {
						enabled   = true
						file_path = "/etc/apt/sources.list"
						index     = 0
						node      = "{{.NodeName}}"
					}`),
					// The computed attributes should be set.
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							testAccResourceRepoSelector,
							tfjsonpath.New(apt.SchemaAttrNameComponents),
							knownvalue.ListPartial(
								map[int]knownvalue.Check{
									// Use the same check for both entries because the sort order cannot be guaranteed.
									0: knownvalue.StringRegexp(regexp.MustCompile(`(contrib|main)`)),
									1: knownvalue.StringRegexp(regexp.MustCompile(`(contrib|main)`)),
								},
							),
						),
						statecheck.ExpectKnownValue(
							testAccResourceRepoSelector,
							tfjsonpath.New(apt.SchemaAttrNamePackageTypes),
							knownvalue.ListPartial(
								map[int]knownvalue.Check{
									0: knownvalue.StringRegexp(regexp.MustCompile(`(deb)`)),
								},
							),
						),
						statecheck.ExpectKnownValue(
							testAccResourceRepoSelector,
							tfjsonpath.New(apt.SchemaAttrNameSuites),
							knownvalue.ListPartial(
								map[int]knownvalue.Check{
									// The possible Debian version is based on the official table of the Proxmox VE FAQ page:
									// - https://pve.proxmox.com/wiki/FAQ#faq-support-table
									// - https://www.thomas-krenn.com/en/wiki/Proxmox_VE#Proxmox_VE_8.x
									//
									// The required Proxmox VE version for this provider is of course also taken into account:
									// - https://github.com/bpg/terraform-provider-proxmox?tab=readme-ov-file#requirements
									0: knownvalue.StringRegexp(regexp.MustCompile(`(bookworm)`)),
								},
							),
						),
						statecheck.ExpectKnownValue(
							testAccResourceRepoSelector,
							tfjsonpath.New(apt.SchemaAttrNameURIs),
							knownvalue.ListPartial(
								map[int]knownvalue.Check{
									0: knownvalue.StringRegexp(regexp.MustCompile(`https?://ftp\.([a-z]+\.)?debian\.org/debian`)),
								},
							),
						),
					},
					// The provided attributes and computed attributes should be set.
					Check: resource.ComposeTestCheckFunc(
						test.ResourceAttributes("proxmox_virtual_environment_apt_repository.test", map[string]string{
							"enabled":   strconv.FormatBool(true),
							"file_path": "/etc/apt/sources.list",
							"index":     strconv.FormatInt(0, 10),
							"node":      te.NodeName,
							"id": fmt.Sprintf(
								"apt_repository_%s_%s_%d",
								strings.ToLower(te.NodeName),
								apt.RepoIDCharReplaceRegEx.ReplaceAllString(
									strings.TrimPrefix("/etc/apt/sources.list", "/"),
									"_",
								),
								0,
							),
						}),
						resource.TestMatchResourceAttr("proxmox_virtual_environment_apt_repository.test", "comment", regexp.MustCompile(`(.*|^$)`)),
						resource.TestCheckResourceAttrSet("proxmox_virtual_environment_apt_repository.test", "file_type"),
					),
				},

				// Test the "ImportState" implementation.
				{
					ImportState: true,
					ImportStateId: fmt.Sprintf(
						"%s,%s,%d",
						strings.ToLower(te.NodeName),
						apitypes.StandardRepoFilePathMain,
						testAccResourceRepoIndex,
					),
					ImportStateVerify: true,
					ResourceName:      testAccResourceRepoSelector,
				},

				// Test the "Update" implementation by toggling the activation status.
				{
					Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_apt_repository" "test" {
						enabled    = false
						file_path  = "/etc/apt/sources.list"
						index     = 0
						node      = "{{.NodeName}}"
					}`),
					// The provides attributes and some computed attributes should be set.
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("proxmox_virtual_environment_apt_repository.test", "enabled", "false"),
					),
				},
			},
		},
	)
}

// Run tests for APT standard repository resource definitions with valid input where all required attributes are
// specified.
// Only the [Create] and [Read] method implementations of the
// [github.com/hashicorp/terraform-plugin-framework/resource.Resource] interface are tested in sequential steps because
// [Delete] and [Update] are no-op due to the non-existing capability of the Proxmox VE API of deleting or updating a
// configured APT standard repository.
//
// [Create]: https://developer.hashicorp.com/terraform/plugin/framework/resources/create
// [Delete]: https://developer.hashicorp.com/terraform/plugin/framework/resources/delete
// [Read]: https://developer.hashicorp.com/terraform/plugin/framework/resources/read
// [Update]: https://developer.hashicorp.com/terraform/plugin/framework/resources/update
func TestAccResourceStandardRepoValidInput(t *testing.T) {
	t.Helper()
	t.Parallel()

	te := test.InitEnvironment(t)

	resource.Test(
		t, resource.TestCase{
			ProtoV6ProviderFactories: te.AccProviders,
			Steps: []resource.TestStep{
				// Test the "Create" and "Read" implementations.
				{
					// 	PUT /api2/json/nodes/{node}/apt/repositories with handle = "no-subscription" will create a new
					// entry in /etc/apt/sources.list on each call :/
					SkipFunc: func() (bool, error) {
						return true, nil
					},
					Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_apt_standard_repository" "test" {
						handle = "no-subscription"
						node   = "{{.NodeName}}"
					}`),
					// The provided attributes and computed attributes should be set.
					Check: resource.ComposeTestCheckFunc(
						test.ResourceAttributes("proxmox_virtual_environment_apt_standard_repository.test", map[string]string{
							"file_path": "/etc/apt/sources.list",
							"handle":    "no-subscription",
							"node":      te.NodeName,
							"status":    "1",
							"id":        fmt.Sprintf("apt_standard_repository_%s_no_subscription", strings.ToLower(te.NodeName)),
						}),
						test.ResourceAttributesSet("proxmox_virtual_environment_apt_standard_repository.test", []string{
							"description",
							"index",
							"name",
						}),
					),
				},

				// Test the "ImportState" implementation.
				{
					SkipFunc: func() (bool, error) {
						return true, nil
					},
					ImportState:       true,
					ImportStateId:     fmt.Sprintf("%s,no-subscription", strings.ToLower(te.NodeName)),
					ImportStateVerify: true,
					ResourceName:      "proxmox_virtual_environment_apt_standard_repository.test",
				},
			},
		},
	)
}
