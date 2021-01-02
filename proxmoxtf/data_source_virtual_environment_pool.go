/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func dataSourceVirtualEnvironmentPool() *schema.Resource {
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
		Read: dataSourceVirtualEnvironmentPoolRead,
	}
}

func dataSourceVirtualEnvironmentPoolRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	poolID := d.Get(mkDataSourceVirtualEnvironmentPoolPoolID).(string)
	pool, err := veClient.GetPool(poolID)

	if err != nil {
		return err
	}

	d.SetId(poolID)

	if pool.Comment != nil {
		d.Set(mkDataSourceVirtualEnvironmentPoolComment, pool.Comment)
	} else {
		d.Set(mkDataSourceVirtualEnvironmentPoolComment, "")
	}

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

	d.Set(mkDataSourceVirtualEnvironmentPoolMembers, members)

	return nil
}
