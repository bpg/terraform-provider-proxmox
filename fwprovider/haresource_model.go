/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fwprovider

import (
	"fmt"

	haresources "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/ha/resources"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// haResourceModel maps the schema data for the High Availability resource data source.
type haResourceModel struct {
	// The Terraform resource identifier
	ID types.String `tfsdk:"id"`
	// The Proxmox HA resource identifier
	ResourceID types.String `tfsdk:"resource_id"`
	// The type of HA resources to fetch. If unset, all resources will be fetched.
	Type types.String `tfsdk:"type"`
	// The desired state of the resource.
	State types.String `tfsdk:"state"`
	// The comment associated with this resource.
	Comment types.String `tfsdk:"comment"`
	// The identifier of the High Availability group this resource is a member of.
	Group types.String `tfsdk:"group"`
	// The maximal number of relocation attempts.
	MaxRelocate types.Int64 `tfsdk:"max_relocate"`
	// The maximal number of restart attempts.
	MaxRestart types.Int64 `tfsdk:"max_restart"`
}

// importFromAPI imports the contents of a HA resource model from the API's response data.
func (d *haResourceModel) importFromAPI(data *haresources.HAResourceGetResponseData) {
	d.ID = data.ID.ToValue()
	d.ResourceID = data.ID.ToValue()
	d.Type = data.Type.ToValue()
	d.State = data.State.ToValue()
	d.Comment = types.StringPointerValue(data.Comment)
	d.Group = types.StringPointerValue(data.Group)
	d.MaxRelocate = types.Int64PointerValue(data.MaxRelocate)
	d.MaxRestart = types.Int64PointerValue(data.MaxRestart)
}

// toRequestBase builds the common request data structure for HA resource creation or update API calls.
func (d *haResourceModel) toRequestBase() haresources.HAResourceDataBase {
	var state proxmoxtypes.HAResourceState

	if d.State.IsNull() {
		state = proxmoxtypes.HAResourceStateStarted
	} else {
		var err error

		state, err = proxmoxtypes.ParseHAResourceState(d.State.ValueString())
		if err != nil {
			panic(fmt.Errorf(
				"state string '%s' wrongly assumed to be valid; error: %w",
				d.State.ValueString(), err,
			))
		}
	}

	return haresources.HAResourceDataBase{
		State:       state,
		Comment:     d.Comment.ValueStringPointer(),
		Group:       d.Group.ValueStringPointer(),
		MaxRelocate: d.MaxRelocate.ValueInt64Pointer(),
		MaxRestart:  d.MaxRestart.ValueInt64Pointer(),
	}
}

// toCreateRequest builds the request data structure for creating a new HA resource.
func (d *haResourceModel) toCreateRequest(resID proxmoxtypes.HAResourceID) *haresources.HAResourceCreateRequestBody {
	return &haresources.HAResourceCreateRequestBody{
		ID:                 resID,
		Type:               &resID.Type,
		HAResourceDataBase: d.toRequestBase(),
	}
}

// toUpdateRequest builds the request data structure for updating an existing HA resource.
func (d *haResourceModel) toUpdateRequest(state *haResourceModel) *haresources.HAResourceUpdateRequestBody {
	var del []string

	if d.Comment.IsNull() && !state.Comment.IsNull() {
		del = append(del, "comment")
	}

	if d.Group.IsNull() && !state.Group.IsNull() {
		del = append(del, "group")
	}

	if d.MaxRelocate.IsNull() && !state.MaxRelocate.IsNull() {
		del = append(del, "max_relocate")
	}

	if d.MaxRestart.IsNull() && !state.MaxRestart.IsNull() {
		del = append(del, "max_restart")
	}

	return &haresources.HAResourceUpdateRequestBody{
		HAResourceDataBase: d.toRequestBase(),
		Delete:             del,
	}
}
