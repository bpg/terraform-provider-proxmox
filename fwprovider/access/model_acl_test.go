/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/access"
)

func TestACLResourceModelFindIn(t *testing.T) {
	t.Parallel()

	acls := []access.ACLGetResponseData{
		{Path: "/", RoleID: "PVEAuditor", Type: "group", UserOrGroupID: "admins"},
		{Path: "/vms", RoleID: "PVEVMUser", Type: "user", UserOrGroupID: "alice@pve"},
		{Path: "/vms", RoleID: "PVEVMAdmin", Type: "token", UserOrGroupID: "alice@pve!ci"},
	}

	model := func(path, role string, group, user, token types.String) *aclResourceModel {
		return &aclResourceModel{Path: path, RoleID: role, GroupID: group, UserID: user, TokenID: token}
	}

	tests := []struct {
		name     string
		model    *aclResourceModel
		expected *access.ACLGetResponseData
	}{
		{
			"group match",
			model("/", "PVEAuditor", types.StringValue("admins"), types.StringNull(), types.StringNull()),
			&acls[0],
		},
		{
			"user match",
			model("/vms", "PVEVMUser", types.StringNull(), types.StringValue("alice@pve"), types.StringNull()),
			&acls[1],
		},
		{
			"token match",
			model("/vms", "PVEVMAdmin", types.StringNull(), types.StringNull(), types.StringValue("alice@pve!ci")),
			&acls[2],
		},
		{
			"same user other role",
			model("/vms", "PVEVMAdmin", types.StringNull(), types.StringValue("alice@pve"), types.StringNull()),
			nil,
		},
		{
			"same user other path",
			model("/", "PVEVMUser", types.StringNull(), types.StringValue("alice@pve"), types.StringNull()),
			nil,
		},
		{
			"group name used as user",
			model("/", "PVEAuditor", types.StringNull(), types.StringValue("admins"), types.StringNull()),
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.expected, tt.model.findIn(acls))
		})
	}
}
