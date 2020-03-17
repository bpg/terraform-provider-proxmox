/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"strings"

	"github.com/danitso/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	dvResourceVirtualEnvironmentGroupComment = ""

	mkResourceVirtualEnvironmentGroupACL          = "acl"
	mkResourceVirtualEnvironmentGroupACLPath      = "path"
	mkResourceVirtualEnvironmentGroupACLPropagate = "propagate"
	mkResourceVirtualEnvironmentGroupACLRoleID    = "role_id"
	mkResourceVirtualEnvironmentGroupComment      = "comment"
	mkResourceVirtualEnvironmentGroupID           = "group_id"
	mkResourceVirtualEnvironmentGroupMembers      = "members"
)

func resourceVirtualEnvironmentGroup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentGroupACL: {
				Type:        schema.TypeSet,
				Description: "The access control list",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentGroupACLPath: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The path",
						},
						mkResourceVirtualEnvironmentGroupACLPropagate: {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to propagate to child paths",
							Default:     false,
						},
						mkResourceVirtualEnvironmentGroupACLRoleID: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The role id",
						},
					},
				},
			},
			mkResourceVirtualEnvironmentGroupComment: {
				Type:        schema.TypeString,
				Description: "The group comment",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentGroupComment,
			},
			mkResourceVirtualEnvironmentGroupID: {
				Type:        schema.TypeString,
				Description: "The group id",
				Required:    true,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentGroupMembers: {
				Type:        schema.TypeSet,
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

	comment := d.Get(mkResourceVirtualEnvironmentGroupComment).(string)
	groupID := d.Get(mkResourceVirtualEnvironmentGroupID).(string)

	body := &proxmox.VirtualEnvironmentGroupCreateRequestBody{
		Comment: &comment,
		ID:      groupID,
	}

	err = veClient.CreateGroup(body)

	if err != nil {
		return err
	}

	d.SetId(groupID)

	aclParsed := d.Get(mkResourceVirtualEnvironmentGroupACL).(*schema.Set).List()

	for _, v := range aclParsed {
		aclDelete := proxmox.CustomBool(false)
		aclEntry := v.(map[string]interface{})
		aclPropagate := proxmox.CustomBool(aclEntry[mkResourceVirtualEnvironmentGroupACLPropagate].(bool))

		aclBody := &proxmox.VirtualEnvironmentACLUpdateRequestBody{
			Delete:    &aclDelete,
			Groups:    []string{groupID},
			Path:      aclEntry[mkResourceVirtualEnvironmentGroupACLPath].(string),
			Propagate: &aclPropagate,
			Roles:     []string{aclEntry[mkResourceVirtualEnvironmentGroupACLRoleID].(string)},
		}

		err := veClient.UpdateACL(aclBody)

		if err != nil {
			return err
		}
	}

	return resourceVirtualEnvironmentGroupRead(d, m)
}

func resourceVirtualEnvironmentGroupRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	groupID := d.Id()
	group, err := veClient.GetGroup(groupID)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")

			return nil
		}

		return err
	}

	acl, err := veClient.GetACL()

	if err != nil {
		return err
	}

	aclParsed := []interface{}{}

	for _, v := range acl {
		if v.Type == "group" && v.UserOrGroupID == groupID {
			aclEntry := map[string]interface{}{}

			aclEntry[mkResourceVirtualEnvironmentGroupACLPath] = v.Path

			if v.Propagate != nil {
				aclEntry[mkResourceVirtualEnvironmentGroupACLPropagate] = bool(*v.Propagate)
			} else {
				aclEntry[mkResourceVirtualEnvironmentGroupACLPropagate] = false
			}

			aclEntry[mkResourceVirtualEnvironmentGroupACLRoleID] = v.RoleID

			aclParsed = append(aclParsed, aclEntry)
		}
	}

	d.Set(mkResourceVirtualEnvironmentGroupACL, aclParsed)

	if group.Comment != nil {
		d.Set(mkResourceVirtualEnvironmentGroupComment, group.Comment)
	} else {
		d.Set(mkResourceVirtualEnvironmentGroupComment, "")
	}

	d.Set(mkResourceVirtualEnvironmentGroupMembers, group.Members)

	return nil
}

func resourceVirtualEnvironmentGroupUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	comment := d.Get(mkResourceVirtualEnvironmentGroupComment).(string)
	groupID := d.Id()

	body := &proxmox.VirtualEnvironmentGroupUpdateRequestBody{
		Comment: &comment,
	}

	err = veClient.UpdateGroup(groupID, body)

	if err != nil {
		return err
	}

	aclArgOld, aclArg := d.GetChange(mkResourceVirtualEnvironmentGroupACL)
	aclParsedOld := aclArgOld.(*schema.Set).List()

	for _, v := range aclParsedOld {
		aclDelete := proxmox.CustomBool(true)
		aclEntry := v.(map[string]interface{})
		aclPropagate := proxmox.CustomBool(aclEntry[mkResourceVirtualEnvironmentGroupACLPropagate].(bool))

		aclBody := &proxmox.VirtualEnvironmentACLUpdateRequestBody{
			Delete:    &aclDelete,
			Groups:    []string{groupID},
			Path:      aclEntry[mkResourceVirtualEnvironmentGroupACLPath].(string),
			Propagate: &aclPropagate,
			Roles:     []string{aclEntry[mkResourceVirtualEnvironmentGroupACLRoleID].(string)},
		}

		err := veClient.UpdateACL(aclBody)

		if err != nil {
			return err
		}
	}

	aclParsed := aclArg.(*schema.Set).List()

	for _, v := range aclParsed {
		aclDelete := proxmox.CustomBool(false)
		aclEntry := v.(map[string]interface{})
		aclPropagate := proxmox.CustomBool(aclEntry[mkResourceVirtualEnvironmentGroupACLPropagate].(bool))

		aclBody := &proxmox.VirtualEnvironmentACLUpdateRequestBody{
			Delete:    &aclDelete,
			Groups:    []string{groupID},
			Path:      aclEntry[mkResourceVirtualEnvironmentGroupACLPath].(string),
			Propagate: &aclPropagate,
			Roles:     []string{aclEntry[mkResourceVirtualEnvironmentGroupACLRoleID].(string)},
		}

		err := veClient.UpdateACL(aclBody)

		if err != nil {
			return err
		}
	}

	return resourceVirtualEnvironmentGroupRead(d, m)
}

func resourceVirtualEnvironmentGroupDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	aclParsed := d.Get(mkResourceVirtualEnvironmentGroupACL).(*schema.Set).List()
	groupID := d.Id()

	for _, v := range aclParsed {
		aclDelete := proxmox.CustomBool(true)
		aclEntry := v.(map[string]interface{})
		aclPropagate := proxmox.CustomBool(aclEntry[mkResourceVirtualEnvironmentGroupACLPropagate].(bool))

		aclBody := &proxmox.VirtualEnvironmentACLUpdateRequestBody{
			Delete:    &aclDelete,
			Groups:    []string{groupID},
			Path:      aclEntry[mkResourceVirtualEnvironmentGroupACLPath].(string),
			Propagate: &aclPropagate,
			Roles:     []string{aclEntry[mkResourceVirtualEnvironmentGroupACLRoleID].(string)},
		}

		err := veClient.UpdateACL(aclBody)

		if err != nil {
			return err
		}
	}

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
