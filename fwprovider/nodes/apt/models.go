/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package apt

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types/nodes/apt"
	api "github.com/bpg/terraform-provider-proxmox/proxmox/nodes/apt/repositories"
)

// Note that most constants are exported to allow the usage in (acceptance) tests.
const (
	// SchemaAttrNameComment is the name of the APT repository schema attribute for the associated comment.
	SchemaAttrNameComment = "comment"

	// SchemaAttrNameComponents is the name of the APT repository schema attribute for the list of components.
	SchemaAttrNameComponents = "components"

	// SchemaAttrNameEnabled is the name of the APT repository schema attribute that indicates the activation status.
	SchemaAttrNameEnabled = "enabled"

	// SchemaAttrNameFilePath is the name of the APT repository schema attribute for the path of the defining source list
	// file.
	SchemaAttrNameFilePath = "file_path"

	// SchemaAttrNameFileType is the name of the APT repository schema attribute for the format of the defining source
	// list file.
	SchemaAttrNameFileType = "file_type"

	// SchemaAttrNameIndex is the name of the APT repository schema attribute for the index within the defining source
	// list file.
	SchemaAttrNameIndex = "index"

	// SchemaAttrNameNode is the name of the APT repository schema attribute for the name of the Proxmox VE node.
	SchemaAttrNameNode = "node"

	// SchemaAttrNamePackageTypes is the name of the APT repository schema attribute for the list of package types.
	SchemaAttrNamePackageTypes = "package_types"

	// SchemaAttrNameStandardDescription is the name of the APT repository schema attribute for the description.
	SchemaAttrNameStandardDescription = "description"

	// SchemaAttrNameStandardHandle is the name of the APT repository schema attribute for the standard repository
	// handle.
	SchemaAttrNameStandardHandle = "handle"

	// SchemaAttrNameStandardName is the name of the APT repository schema attribute for the human-readable name.
	SchemaAttrNameStandardName = "name"

	// SchemaAttrNameStandardStatus is the name of the APT standard repository schema attribute that indicates the
	// configuration and activation status.
	SchemaAttrNameStandardStatus = "status"

	// SchemaAttrNameSuites is the name of the APT repository schema attribute for the list of package distributions.
	SchemaAttrNameSuites = "suites"

	// SchemaAttrNameTerraformID is the name of the APT repository schema attribute for the Terraform ID.
	SchemaAttrNameTerraformID = "id"

	// SchemaAttrNameURIs is the name of the APT repository schema attribute for the list of repository URIs.
	SchemaAttrNameURIs = "uris"
)

// RepoIDCharReplaceRegEx is a regular expression to replace characters in a Terraform resource/data source ID.
// The "^" at the beginning of the character group selects all characters not matching the group.
var RepoIDCharReplaceRegEx = regexp.MustCompile(`([^a-zA-Z1-9_])`)

// modelRepo maps the schema data for an APT repository from a parsed source list file.
type modelRepo struct {
	// Comment is the comment of the APT repository.
	Comment types.String `tfsdk:"comment"`

	// Components is the list of repository components.
	Components types.List `tfsdk:"components"`

	// Enabled indicates whether the APT repository is enabled.
	Enabled types.Bool `tfsdk:"enabled"`

	// FilePath is the path of the source list file that contains the APT repository.
	FilePath types.String `tfsdk:"file_path"`

	// FileType is the format of the packages.
	FileType types.String `tfsdk:"file_type"`

	// ID is the Terraform identifier of the APT repository.
	ID types.String `tfsdk:"id"`

	// Index is the index of the APT repository within the defining source list.
	Index types.Int64 `tfsdk:"index"`

	// Node is the name of the Proxmox VE node for the APT repository.
	Node types.String `tfsdk:"node"`

	// PackageTypes is the list of package types.
	PackageTypes types.List `tfsdk:"package_types"`

	// Suites is the list of package distributions.
	Suites types.List `tfsdk:"suites"`

	// URIs is the list of repository URIs.
	URIs types.List `tfsdk:"uris"`
}

// modelStandardRepo maps the schema data for an APT standard repository.
type modelStandardRepo struct {
	// Description is the description of the APT standard repository.
	Description types.String `tfsdk:"description"`

	// FilePath is the path of the source list file that contains the APT standard repository.
	FilePath types.String `tfsdk:"file_path"`

	// ID is the Terraform identifier of the APT standard repository.
	ID types.String `tfsdk:"id"`

	// Index is the index of the APT standard repository within the defining source list file.
	Index types.Int64 `tfsdk:"index"`

	// Handle is the handle of the APT standard repository.
	Handle customtypes.StandardRepoHandleValue `tfsdk:"handle"`

	// Name is the name of the APT standard repository.
	Name types.String `tfsdk:"name"`

	// Node is the name of the Proxmox VE node for the APT standard repository.
	Node types.String `tfsdk:"node"`

	// Status is the configuration and activation status of the APT standard repository.
	Status types.Int64 `tfsdk:"status"`
}

// importFromAPI imports the contents of an APT repository model from the Proxmox VE API's response data.
func (rp *modelRepo) importFromAPI(ctx context.Context, data *api.GetResponseData) diag.Diagnostics {
	diags := diag.Diagnostics{}

	// We can only ensure a unique ID by using the name of the Proxmox VE node and the absolute file path because custom
	// source list files can be loaded by Proxmox VE from every path on a node.
	rp.ID = types.StringValue(
		fmt.Sprintf(
			"%s_%s_%s_%d",
			ResourceRepoIDPrefix,
			strings.ToLower(rp.Node.ValueString()),
			strings.ToLower(RepoIDCharReplaceRegEx.ReplaceAllString(strings.TrimPrefix(rp.FilePath.ValueString(), "/"), "_")),
			rp.Index.ValueInt64(),
		),
	)

	// We must ensure that the type definitions for lists and other attributes are set since Terraform must know these
	// during the planning phase. This is important when the resource was imported where only the ID is known.
	rp.Comment = types.StringNull()
	rp.Enabled = types.BoolNull()
	rp.FileType = types.StringNull()
	rp.Components = types.ListNull(types.StringType)
	rp.PackageTypes = types.ListNull(types.StringType)
	rp.Suites = types.ListNull(types.StringType)
	rp.URIs = types.ListNull(types.StringType)

	// Iterate through all repository files…
	for _, repoFile := range data.Files {
		// …and the defined repositories when the file path matches.
		if repoFile.Path == rp.FilePath.ValueString() {
			// Handle situations where an APT repository might have been removed manually which is currently the only way to
			// solve this with the capabilities of the Proxmox VE API.
			if int64(len(repoFile.Repositories)) > rp.Index.ValueInt64() {
				repo := repoFile.Repositories[rp.Index.ValueInt64()]

				// Strip the unnecessary new line control character (\n) from the end of the comment that is, for whatever
				// reason, returned this way by the Proxmox VE API.
				if repo.Comment != nil {
					rp.Comment = types.StringValue(strings.TrimSuffix(*repo.Comment, "\n"))
				}

				rp.Enabled = repo.Enabled.ToValue()
				rp.FileType = types.StringValue(repo.FileType)

				components, convDiags := types.ListValueFrom(ctx, types.StringType, repo.Components)
				if convDiags.HasError() {
					diags.AddError("Terraform list value conversion", "Convert list of APT repository components")
				} else {
					rp.Components = components
				}

				pkgTypes, convDiags := types.ListValueFrom(ctx, types.StringType, repo.PackageTypes)
				if convDiags.HasError() {
					diags.AddError("Terraform list value conversion", "Convert list of APT repository package types")
				} else {
					rp.PackageTypes = pkgTypes
				}

				suites, convDiags := types.ListValueFrom(ctx, types.StringType, repo.Suites)
				if convDiags.HasError() {
					diags.AddError("Terraform list value conversion", "Convert list of APT repository suites")
				} else {
					rp.Suites = suites
				}

				uris, convDiags := types.ListValueFrom(ctx, types.StringType, repo.URIs)
				if convDiags.HasError() {
					diags.AddError("Terraform list value conversion", "Convert list of APT repository URIs")
				} else {
					rp.URIs = uris
				}
			}
		}
	}

	return diags
}

// importFromAPI imports the contents of an APT standard repository from the Proxmox VE API's response data.
func (srp *modelStandardRepo) importFromAPI(_ context.Context, data *api.GetResponseData) {
	for _, repo := range data.StandardRepos {
		if repo.Handle == srp.Handle.ValueString() {
			srp.Description = types.StringPointerValue(repo.Description)
			// We can only ensure a unique ID by using the name of the Proxmox VE node in combination with the unique standard
			// handle.
			srp.ID = types.StringValue(
				fmt.Sprintf(
					"%s_%s_%s",
					ResourceStandardRepoIDPrefix,
					strings.ToLower(srp.Node.ValueString()),
					RepoIDCharReplaceRegEx.ReplaceAllString(srp.Handle.ValueString(), "_"),
				),
			)

			srp.Name = types.StringValue(repo.Name)
			srp.Status = types.Int64PointerValue(repo.Status)
		}
	}

	// Set the index…
	srp.setIndex(data)
	// … and then the file path when the index is valid…
	if !srp.Index.IsNull() {
		// …by iterating through all repository files…
		for _, repoFile := range data.Files {
			// …and get the repository when the file path matches.
			if srp.Handle.IsSupportedFilePath(repoFile.Path) {
				srp.FilePath = types.StringValue(repoFile.Path)
			}
		}
	}
}

// setIndex sets the index of the APT standard repository derived from the defining source list file.
func (srp *modelStandardRepo) setIndex(data *api.GetResponseData) {
	for _, file := range data.Files {
		for idx, repo := range file.Repositories {
			if slices.Contains(repo.Components, srp.Handle.ComponentName()) {
				// Return early for non-Ceph repositories…
				if !srp.Handle.IsCephHandle() {
					srp.Index = types.Int64Value(int64(idx))

					return
				}

				// …and find the index for Ceph repositories based on the version name within the list of URIs.
				for _, uri := range repo.URIs {
					if strings.Contains(uri, srp.Handle.CephVersionName().String()) {
						srp.Index = types.Int64Value(int64(idx))

						return
					}
				}
			}
		}
	}

	srp.Index = types.Int64Null()
}
