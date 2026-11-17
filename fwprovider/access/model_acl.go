/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/proxmox/access"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
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

const aclIDFormat = "{path}?{group|user@realm|user@realm!token}?{role}"

func (r *aclResourceModel) entityID() string {
	return r.GroupID.ValueString() + r.TokenID.ValueString() + r.UserID.ValueString()
}

func (r *aclResourceModel) generateID() types.String {
	return types.StringValue(r.Path + "?" + r.entityID() + "?" + r.RoleID)
}

// findIn returns the entry of acls that matches the path, role and group, user or token of the model.
func (r *aclResourceModel) findIn(acls []access.ACLGetResponseData) *access.ACLGetResponseData {
	for i := range acls {
		acl := &acls[i]

		var entityID types.String

		switch acl.Type {
		case "group":
			entityID = r.GroupID
		case "token":
			entityID = r.TokenID
		case "user":
			entityID = r.UserID
		default:
			continue
		}

		if acl.UserOrGroupID == entityID.ValueString() && acl.Path == r.Path && acl.RoleID == r.RoleID {
			return acl
		}
	}

	return nil
}

func parseACLResourceModelFromID(id string) (*aclResourceModel, error) {
	parts := strings.Split(id, "?")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid ACL resource ID format %#v, expected %v", id, aclIDFormat)
	}

	path := parts[0]
	entityID := parts[1]
	roleID := parts[2]

	model := &aclResourceModel{
		ID:        types.StringValue(id),
		GroupID:   types.StringNull(),
		Path:      path,
		Propagate: false,
		RoleID:    roleID,
		TokenID:   types.StringNull(),
		UserID:    types.StringNull(),
	}

	switch {
	case strings.Contains(entityID, "!"):
		model.TokenID = types.StringValue(entityID)
	case strings.Contains(entityID, "@"):
		model.UserID = types.StringValue(entityID)
	default:
		model.GroupID = types.StringValue(entityID)
	}

	return model, nil
}

func (r *aclResourceModel) intoUpdateBody() *access.ACLUpdateRequestBody {
	body := &access.ACLUpdateRequestBody{
		Groups:    nil,
		Path:      r.Path,
		Propagate: proxmoxtypes.CustomBool(r.Propagate).Pointer(),
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
