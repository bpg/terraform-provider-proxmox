/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/firewall"
)

const (
	mkSecurityGroupsSecurityGroupNames = "security_group_names"
)

func SecurityGroupsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		mkSecurityGroupsSecurityGroupNames: {
			Type:        schema.TypeList,
			Description: "Security Group Names",
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func SecurityGroupsRead(ctx context.Context, fw *firewall.API, d *schema.ResourceData) diag.Diagnostics {
	groups, err := fw.ListGroups(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	groupNames := make([]interface{}, len(groups))

	for i, v := range groups {
		groupNames[i] = v.Group
	}

	d.SetId(uuid.New().String())

	err = d.Set(mkSecurityGroupsSecurityGroupNames, groupNames)

	return diag.FromErr(err)
}
