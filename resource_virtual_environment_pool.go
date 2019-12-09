/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"strings"

	"github.com/danitso/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	mkResourceVirtualEnvironmentPoolComment                 = "comment"
	mkResourceVirtualEnvironmentPoolMembers                 = "members"
	mkResourceVirtualEnvironmentPoolMembersID               = "id"
	mkResourceVirtualEnvironmentPoolMembersNode             = "node"
	mkResourceVirtualEnvironmentPoolMembersStorage          = "storage"
	mkResourceVirtualEnvironmentPoolMembersType             = "type"
	mkResourceVirtualEnvironmentPoolMembersVirtualMachineID = "virtual_machine_id"
	mkResourceVirtualEnvironmentPoolPoolID                  = "pool_id"
)

func resourceVirtualEnvironmentPool() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentPoolComment: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The pool comment",
				Optional:    true,
				Default:     "",
			},
			mkResourceVirtualEnvironmentPoolMembers: &schema.Schema{
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
			mkResourceVirtualEnvironmentPoolPoolID: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The pool id",
				Required:    true,
				ForceNew:    true,
			},
		},
		Create: resourceVirtualEnvironmentPoolCreate,
		Read:   resourceVirtualEnvironmentPoolRead,
		Update: resourceVirtualEnvironmentPoolUpdate,
		Delete: resourceVirtualEnvironmentPoolDelete,
	}
}

func resourceVirtualEnvironmentPoolCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	comment := d.Get(mkResourceVirtualEnvironmentPoolComment).(string)
	poolID := d.Get(mkResourceVirtualEnvironmentPoolPoolID).(string)

	body := &proxmox.VirtualEnvironmentPoolCreateRequestBody{
		Comment: &comment,
		ID:      poolID,
	}

	err = veClient.CreatePool(body)

	if err != nil {
		return err
	}

	d.SetId(poolID)

	return resourceVirtualEnvironmentPoolRead(d, m)
}

func resourceVirtualEnvironmentPoolRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	poolID := d.Id()
	pool, err := veClient.GetPool(poolID)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")

			return nil
		}

		return err
	}

	d.SetId(poolID)

	if pool.Comment != nil {
		d.Set(mkResourceVirtualEnvironmentPoolComment, pool.Comment)
	} else {
		d.Set(mkResourceVirtualEnvironmentPoolComment, "")
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

func resourceVirtualEnvironmentPoolUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	comment := d.Get(mkResourceVirtualEnvironmentPoolComment).(string)
	poolID := d.Id()

	body := &proxmox.VirtualEnvironmentPoolUpdateRequestBody{
		Comment: &comment,
	}

	err = veClient.UpdatePool(poolID, body)

	if err != nil {
		return err
	}

	return resourceVirtualEnvironmentPoolRead(d, m)
}

func resourceVirtualEnvironmentPoolDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	poolID := d.Id()
	err = veClient.DeletePool(poolID)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")

			return nil
		}

		return err
	}

	d.SetId("")

	return nil
}
