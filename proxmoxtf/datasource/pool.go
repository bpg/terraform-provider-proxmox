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
	mkDataSourceVirtualEnvironmentPoolComment            = "comment"
	mkDataSourceVirtualEnvironmentPoolMembers            = "members"
	mkDataSourceVirtualEnvironmentPoolMembersDatastoreID = "datastore_id"
	mkDataSourceVirtualEnvironmentPoolMembersID          = "id"
	mkDataSourceVirtualEnvironmentPoolMembersNodeName    = "node_name"
	mkDataSourceVirtualEnvironmentPoolMembersType        = "type"
	mkDataSourceVirtualEnvironmentPoolMembersVMID        = "vm_id"
	mkDataSourceVirtualEnvironmentPoolPoolID             = "pool_id"
)

// Pool returns a resource for a single Proxmox pool.
func Pool() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentPoolComment: {
				Type:        schema.TypeString,
				Description: "The pool comment",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentPoolMembers: {
				Type:        schema.TypeList,
				Description: "The pool members",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkDataSourceVirtualEnvironmentPoolMembersDatastoreID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The datastore id",
						},
						mkDataSourceVirtualEnvironmentPoolMembersID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The member id",
						},
						mkDataSourceVirtualEnvironmentPoolMembersNodeName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The node name",
						},
						mkDataSourceVirtualEnvironmentPoolMembersType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The member type",
						},
						mkDataSourceVirtualEnvironmentPoolMembersVMID: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The virtual machine id",
						},
					},
				},
			},
			mkDataSourceVirtualEnvironmentPoolPoolID: {
				Type:        schema.TypeString,
				Description: "The pool id",
				Required:    true,
			},
		},
		ReadContext: poolRead,
	}
}

func poolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	poolID := d.Get(mkDataSourceVirtualEnvironmentPoolPoolID).(string)
	pool, err := veClient.GetPool(ctx, poolID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(poolID)

	if pool.Comment != nil {
		err = d.Set(mkDataSourceVirtualEnvironmentPoolComment, pool.Comment)
	} else {
		err = d.Set(mkDataSourceVirtualEnvironmentPoolComment, "")
	}
	diags = append(diags, diag.FromErr(err)...)

	members := make([]interface{}, len(pool.Members))

	for i, v := range pool.Members {
		values := map[string]interface{}{}

		values[mkDataSourceVirtualEnvironmentPoolMembersID] = v.ID
		values[mkDataSourceVirtualEnvironmentPoolMembersNodeName] = v.Node

		if v.DatastoreID != nil {
			values[mkDataSourceVirtualEnvironmentPoolMembersDatastoreID] = v.DatastoreID
		} else {
			values[mkDataSourceVirtualEnvironmentPoolMembersDatastoreID] = ""
		}

		values[mkDataSourceVirtualEnvironmentPoolMembersType] = v.Type

		if v.VMID != nil {
			values[mkDataSourceVirtualEnvironmentPoolMembersVMID] = v.VMID
		} else {
			values[mkDataSourceVirtualEnvironmentPoolMembersVMID] = 0
		}

		members[i] = values
	}

	err = d.Set(mkDataSourceVirtualEnvironmentPoolMembers, members)
	diags = append(diags, diag.FromErr(err)...)

	return diags
}
