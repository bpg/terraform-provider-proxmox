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
	mkResourceVirtualEnvironmentGroupComment = "comment"
	mkResourceVirtualEnvironmentGroupID      = "group_id"
	mkResourceVirtualEnvironmentGroupMembers = "members"
)

func resourceVirtualEnvironmentGroup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentGroupComment: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The group comment",
				Optional:    true,
				Default:     "",
			},
			mkResourceVirtualEnvironmentGroupID: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The group id",
				Required:    true,
			},
			mkResourceVirtualEnvironmentGroupMembers: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The group members",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Create: resourceVirtualEnvironmentGroupCreate,
		Read:   resourceVirtualEnvironmentGroupRead,
		Update: resourceVirtualEnvironmentGroupUpdate,
		Delete: resourceVirtualEnvironmentGroupDelete,
	}
}

func resourceVirtualEnvironmentGroupCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	groupID := d.Get(mkResourceVirtualEnvironmentGroupID).(string)
	body := &proxmox.VirtualEnvironmentGroupCreateRequestBody{
		Comment: d.Get(mkResourceVirtualEnvironmentGroupComment).(string),
		ID:      groupID,
	}

	err = veClient.CreateGroup(body)

	if err != nil {
		return err
	}

	d.SetId(groupID)

	return resourceVirtualEnvironmentGroupRead(d, m)
}

func resourceVirtualEnvironmentGroupRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	groupID := d.Id()
	accessGroup, err := veClient.GetGroup(groupID)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")

			return nil
		}

		return err
	}

	d.SetId(groupID)

	d.Set(mkResourceVirtualEnvironmentGroupComment, accessGroup.Comment)
	d.Set(mkResourceVirtualEnvironmentGroupMembers, accessGroup.Members)

	return nil
}

func resourceVirtualEnvironmentGroupUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	body := &proxmox.VirtualEnvironmentGroupUpdateRequestBody{
		Comment: d.Get(mkResourceVirtualEnvironmentGroupComment).(string),
	}

	groupID := d.Id()
	err = veClient.UpdateGroup(groupID, body)

	if err != nil {
		return err
	}

	return resourceVirtualEnvironmentGroupRead(d, m)
}

func resourceVirtualEnvironmentGroupDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	groupID := d.Id()
	err = veClient.DeleteGroup(groupID)

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
