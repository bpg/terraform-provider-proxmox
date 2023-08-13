/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datasource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	mkDataSourceVirtualEnvironmentHAGroupsGroupIDs = "group_ids"
)

// HAGroups returns a resource that lists the High Availability groups.
func HAGroups() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentHAGroupsGroupIDs: {
				Type:        schema.TypeList,
				Description: "The HA group identifiers",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		ReadContext: haGroupsRead,
	}
}

// Read the list of HA groups from the Proxmox cluster then convert it to a list of strings.
func haGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)

	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	list, err := api.Cluster().HA().Groups().List(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	groupIDs := make([]interface{}, len(list))
	for i, v := range list {
		groupIDs[i] = v.ID
	}

	d.SetId("hagroups")

	err = d.Set(mkDataSourceVirtualEnvironmentHAGroupsGroupIDs, groupIDs)
	return diag.FromErr(err)
}
