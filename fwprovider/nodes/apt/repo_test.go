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
					Config: fmt.Sprintf(
						`
					data %q %q {
						%s = %q # file_path
						%s = %d # index
						%s = %q # node
					}
					`,
						strings.Split(testAccResourceRepoSelector, ".")[0],
						strings.Split(testAccResourceRepoSelector, ".")[1],
						// To ensure stable acceptance tests we must use one of the Proxmox VE default source lists that always
						// exists on any (new) Proxmox VE node.
						apt.SchemaAttrNameFilePath, apitypes.StandardRepoFilePathMain,
						apt.SchemaAttrNameIndex, testAccResourceRepoIndex,
						apt.SchemaAttrNameNode, te.NodeName,
					),
					// The provided attributes and computed attributes should be set.
					Check: resource.ComposeTestCheckFunc(
						resource.TestMatchResourceAttr(
							fmt.Sprintf("data.%s", testAccResourceRepoSelector),
							apt.SchemaAttrNameComment,
							// Expect any value or an empty string.
							regexp.MustCompile(`(.*|^$)`),
						),
						resource.TestCheckResourceAttr(
							fmt.Sprintf("data.%s", testAccResourceRepoSelector),
							apt.SchemaAttrNameTerraformID,
							fmt.Sprintf(
								"%s_%s_%s_%d",
								apt.ResourceRepoIDPrefix,
								strings.ToLower(te.NodeName),
								apt.RepoIDCharReplaceRegEx.ReplaceAllString(
									strings.TrimPrefix(apitypes.StandardRepoFilePathMain, "/"),
									"_",
								),
								testAccResourceRepoIndex,
							),
						),
						test.ResourceAttributesSet(
							fmt.Sprintf("data.%s", testAccResourceRepoSelector),
							[]string{
								fmt.Sprintf("%s.#", apt.SchemaAttrNameComponents),
								apt.SchemaAttrNameEnabled,
								apt.SchemaAttrNameFilePath,
								apt.SchemaAttrNameIndex,
								apt.SchemaAttrNameNode,
								fmt.Sprintf("%s.#", apt.SchemaAttrNamePackageTypes),
								fmt.Sprintf("%s.#", apt.SchemaAttrNameSuites),
								fmt.Sprintf("%s.#", apt.SchemaAttrNameURIs),
							},
						),
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
					Config: fmt.Sprintf(
						`
					data %q %q {
						%s = %q # handle
						%s = %q # node
					}
					`,
						strings.Split(testAccResourceStandardRepoSelector, ".")[0],
						strings.Split(testAccResourceStandardRepoSelector, ".")[1],
						apt.SchemaAttrNameStandardHandle, apitypes.StandardRepoHandleKindNoSubscription,
						apt.SchemaAttrNameNode, te.NodeName,
					),
					// The provided attributes and computed attributes should be set.
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(
							fmt.Sprintf("data.%s", testAccResourceStandardRepoSelector),
							apt.SchemaAttrNameTerraformID,
							fmt.Sprintf(
								"%s_%s_%s",
								apt.ResourceStandardRepoIDPrefix,
								strings.ToLower(te.NodeName),
								apt.RepoIDCharReplaceRegEx.ReplaceAllString(
									strings.TrimPrefix(apitypes.StandardRepoHandleKindNoSubscription.String(), "/"),
									"_",
								),
							),
						),
						test.ResourceAttributesSet(
							fmt.Sprintf("data.%s", testAccResourceStandardRepoSelector),
							// Note that we can not check for the following attributes because they are only available when the
							// standard repository has been added to a source list:
							//
							// - apt.SchemaAttrNameFilePath (file_path) - will be set when parsing all configured repositories in all
							//   source list files.
							// - apt.SchemaAttrNameIndex (index) - will be set when finding the repository within a source list file,
							//   based on the detected file path.
							// - apt.SchemaAttrNameStandardStatus (status) - is only available when the standard has been configured.
							[]string{
								apt.SchemaAttrNameStandardDescription,
								apt.SchemaAttrNameStandardHandle,
								apt.SchemaAttrNameStandardName,
								apt.SchemaAttrNameNode,
							},
						),
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
					Config: fmt.Sprintf(
						`
					resource %q %q {
						%s = %t # enabled
						%s = %q # file_path
						%s = %d # index
						%s = %q # node
					}
					`,
						strings.Split(testAccResourceRepoSelector, ".")[0],
						strings.Split(testAccResourceRepoSelector, ".")[1],
						apt.SchemaAttrNameEnabled, apt.ResourceRepoActivationStatus,
						// To ensure stable acceptance tests we must use one of the Proxmox VE default source lists that always
						// exists on any (new) Proxmox VE node.
						apt.SchemaAttrNameFilePath, apitypes.StandardRepoFilePathMain,
						apt.SchemaAttrNameIndex, testAccResourceRepoIndex,
						apt.SchemaAttrNameNode, te.NodeName,
					),
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
						resource.TestMatchResourceAttr(
							testAccResourceRepoSelector,
							apt.SchemaAttrNameComment,
							// Expect any value or an empty string.
							regexp.MustCompile(`(.*|^$)`),
						),
						resource.TestCheckResourceAttr(
							testAccResourceRepoSelector,
							apt.SchemaAttrNameEnabled,
							strconv.FormatBool(apt.ResourceRepoActivationStatus),
						),
						resource.TestCheckResourceAttr(
							testAccResourceRepoSelector,
							apt.SchemaAttrNameFilePath,
							apitypes.StandardRepoFilePathMain,
						),
						resource.TestCheckResourceAttrSet(testAccResourceRepoSelector, apt.SchemaAttrNameFileType),
						resource.TestCheckResourceAttr(
							testAccResourceRepoSelector,
							apt.SchemaAttrNameIndex,
							strconv.FormatInt(testAccResourceRepoIndex, 10),
						),
						resource.TestCheckResourceAttr(testAccResourceRepoSelector, apt.SchemaAttrNameNode, te.NodeName),
						resource.TestCheckResourceAttr(
							testAccResourceRepoSelector,
							apt.SchemaAttrNameTerraformID,
							fmt.Sprintf(
								"%s_%s_%s_%d",
								apt.ResourceRepoIDPrefix,
								strings.ToLower(te.NodeName),
								apt.RepoIDCharReplaceRegEx.ReplaceAllString(
									strings.TrimPrefix(apitypes.StandardRepoFilePathMain, "/"),
									"_",
								),
								testAccResourceRepoIndex,
							),
						),
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
					Config: fmt.Sprintf(
						`
					resource %q %q {
						%s = %t # enabled
						%s = %q # file_path
						%s = %d # index
						%s = %q # node
					}
					`,
						strings.Split(testAccResourceRepoSelector, ".")[0],
						strings.Split(testAccResourceRepoSelector, ".")[1],
						// Disable the repository which is enabled by default for created or imported resources.
						apt.SchemaAttrNameEnabled, !apt.ResourceRepoActivationStatus,
						// To ensure stable acceptance tests we must use one of the Proxmox VE default source lists that always
						// exists on any (new) Proxmox VE node.s
						apt.SchemaAttrNameFilePath, apitypes.StandardRepoFilePathMain,
						apt.SchemaAttrNameIndex, testAccResourceRepoIndex,
						apt.SchemaAttrNameNode, te.NodeName,
					),
					// The provides attributes and some computed attributes should be set.
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(
							testAccResourceRepoSelector,
							apt.SchemaAttrNameEnabled,
							strconv.FormatBool(!apt.ResourceRepoActivationStatus),
						),
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
					Config: fmt.Sprintf(
						`
					resource %q %q {
						%s = %q # handle
						%s = %q # node
					}
					`,
						strings.Split(testAccResourceStandardRepoSelector, ".")[0],
						strings.Split(testAccResourceStandardRepoSelector, ".")[1],
						apt.SchemaAttrNameStandardHandle, testAccResourceStandardRepoHandle,
						apt.SchemaAttrNameNode, te.NodeName,
					),
					// The provided attributes and computed attributes should be set.
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet(
							testAccResourceStandardRepoSelector,
							apt.SchemaAttrNameStandardDescription,
						),
						resource.TestCheckResourceAttr(
							testAccResourceStandardRepoSelector,
							apt.SchemaAttrNameFilePath,
							apitypes.StandardRepoFilePathMain,
						),
						resource.TestCheckResourceAttr(
							testAccResourceStandardRepoSelector,
							apt.SchemaAttrNameStandardHandle,
							testAccResourceStandardRepoHandle,
						),
						resource.TestCheckResourceAttrSet(testAccResourceStandardRepoSelector, apt.SchemaAttrNameIndex),
						resource.TestCheckResourceAttrSet(testAccResourceStandardRepoSelector, apt.SchemaAttrNameStandardName),
						resource.TestCheckResourceAttr(testAccResourceStandardRepoSelector, apt.SchemaAttrNameNode, te.NodeName),
						resource.TestCheckResourceAttr(
							testAccResourceStandardRepoSelector,
							apt.SchemaAttrNameStandardStatus,
							// By default, newly added APT standard repositories are enabled.
							strconv.Itoa(1),
						),
						resource.TestCheckResourceAttr(
							testAccResourceStandardRepoSelector,
							apt.SchemaAttrNameTerraformID,
							fmt.Sprintf(
								"%s_%s_%s",
								apt.ResourceStandardRepoIDPrefix,
								strings.ToLower(te.NodeName),
								apt.RepoIDCharReplaceRegEx.ReplaceAllString(testAccResourceStandardRepoHandle, "_"),
							),
						),
					),
				},

				// Test the "ImportState" implementation.
				{
					ImportState:       true,
					ImportStateId:     fmt.Sprintf("%s,%s", strings.ToLower(te.NodeName), testAccResourceStandardRepoHandle),
					ImportStateVerify: true,
					ResourceName:      testAccResourceStandardRepoSelector,
				},
			},
		},
	)
}
