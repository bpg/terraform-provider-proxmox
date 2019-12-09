/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	mkDataSourceVirtualEnvironmentPoolComment = "comment"
	mkDataSourceVirtualEnvironmentPoolMembers = "members"
	mkDataSourceVirtualEnvironmentPoolPoolID  = "pool_id"
)

func dataSourceVirtualEnvironmentPool() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentPoolComment: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The pool comment",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentPoolMembers: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The pool members",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentPoolMembersID: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The member id",
						},
						mkResourceVirtualEnvironmentPoolMembersNode: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The node id",
						},
						mkResourceVirtualEnvironmentPoolMembersStorage: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The storage id",
							Default:     "",
						},
						mkResourceVirtualEnvironmentPoolMembersType: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The member type",
						},
						mkResourceVirtualEnvironmentPoolMembersVirtualMachineID: {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The virtual machine id",
							Default:     0,
						},
					},
				},
			},
			mkDataSourceVirtualEnvironmentPoolPoolID: &schema.Schema{
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
		values := make(map[string]interface{})

		values[mkResourceVirtualEnvironmentPoolMembersID] = v.ID
		values[mkResourceVirtualEnvironmentPoolMembersNode] = v.Node

		if v.Storage != nil {
			values[mkResourceVirtualEnvironmentPoolMembersStorage] = v.Storage
		} else {
			values[mkResourceVirtualEnvironmentPoolMembersStorage] = ""
		}

		values[mkResourceVirtualEnvironmentPoolMembersType] = v.Type

		if v.VirtualMachineID != nil {
			values[mkResourceVirtualEnvironmentPoolMembersVirtualMachineID] = v.VirtualMachineID
		} else {
			values[mkResourceVirtualEnvironmentPoolMembersVirtualMachineID] = 0
		}

		members[i] = values
	}

	d.Set(mkResourceVirtualEnvironmentPoolMembers, members)

	return nil
}
