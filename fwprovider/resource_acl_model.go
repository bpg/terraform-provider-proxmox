/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fwprovider

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gorilla/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"

	"github.com/bpg/terraform-provider-proxmox/proxmox/access"
)

type aclResourceModel struct {
	ID types.String `tfsdk:"id"`

	GroupID   types.String `tfsdk:"group_id"`
	Path      string       `tfsdk:"path"`
	Propagate bool         `tfsdk:"propagate"`
	RoleID    string       `tfsdk:"role_id"`
	TokenID   types.String `tfsdk:"token_id"`
	UserID    types.String `tfsdk:"user_id"`
}

type aclResouceIDFields struct {
	EntityID string `schema:"entity_id"`
	RoleID   string `schema:"role_id"`
}

const aclIDFormat = "{path}?entity_id={group|user@realm|user@realm!token}?role_id={role}"

func (r *aclResourceModel) generateID() error {
	encoder := schema.NewEncoder()

	fields := aclResouceIDFields{
		EntityID: r.GroupID.ValueString() + r.TokenID.ValueString() + r.UserID.ValueString(),
		RoleID:   r.RoleID,
	}
	v := url.Values{}

	err := encoder.Encode(fields, v)
	if err != nil {
		return fmt.Errorf("failed to encode ACL resource ID: %w", err)
	}

	r.ID = types.StringValue(r.Path + "?" + v.Encode())

	return nil
}

func parseACLResourceModelFromID(id string) (*aclResourceModel, error) {
	path, query, found := strings.Cut(id, "?")

	if !found {
		return nil, fmt.Errorf("invalid ACL resource ID format %#v, expected %v", id, aclIDFormat)
	}

	v, err := url.ParseQuery(query)
	if err != nil {
		return nil, fmt.Errorf("invalid ACL resource ID format %#v, expected %v: %w", id, aclIDFormat, err)
	}

	decoder := schema.NewDecoder()

	fields := aclResouceIDFields{}

	err = decoder.Decode(&fields, v)
	if err != nil {
		return nil, fmt.Errorf("invalid ACL resource ID format %#v, expected %v: %w", id, aclIDFormat, err)
	}

	model := &aclResourceModel{
		ID:        types.StringValue(id),
		GroupID:   types.StringNull(),
		Path:      path,
		Propagate: false,
		RoleID:    fields.RoleID,
		TokenID:   types.StringNull(),
		UserID:    types.StringNull(),
	}

	switch {
	case strings.Contains(fields.EntityID, "!"):
		model.TokenID = types.StringValue(fields.EntityID)
	case strings.Contains(fields.EntityID, "@"):
		model.UserID = types.StringValue(fields.EntityID)
	default:
		model.GroupID = types.StringValue(fields.EntityID)
	}

	return model, nil
}

func (r *aclResourceModel) intoUpdateBody() *access.ACLUpdateRequestBody {
	body := &access.ACLUpdateRequestBody{
		Groups:    nil,
		Path:      r.Path,
		Propagate: proxmoxtypes.CustomBoolPtr(r.Propagate),
		Roles:     []string{r.RoleID},
		Tokens:    nil,
		Users:     nil,
	}

	if !r.GroupID.IsNull() {
		body.Groups = []string{r.GroupID.ValueString()}
	}

	if !r.TokenID.IsNull() {
		body.Tokens = []string{r.TokenID.ValueString()}
	}

	if !r.UserID.IsNull() {
		body.Users = []string{r.UserID.ValueString()}
	}

	return body
}
