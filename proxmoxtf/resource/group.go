/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/internal/types"
	"github.com/bpg/terraform-provider-proxmox/proxmox/access"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
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

// Group returns a resource that manages a group in the Proxmox VE access control list.
func Group() *schema.Resource {
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
		CreateContext: groupCreate,
		ReadContext:   groupRead,
		UpdateContext: groupUpdate,
		DeleteContext: groupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func groupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	comment := d.Get(mkResourceVirtualEnvironmentGroupComment).(string)
	groupID := d.Get(mkResourceVirtualEnvironmentGroupID).(string)

	body := &access.GroupCreateRequestBody{
		Comment: &comment,
		ID:      groupID,
	}

	err = api.Access().CreateGroup(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(groupID)

	aclParsed := d.Get(mkResourceVirtualEnvironmentGroupACL).(*schema.Set).List()

	for _, v := range aclParsed {
		aclDelete := types.CustomBool(false)
		aclEntry := v.(map[string]interface{})
		aclPropagate := types.CustomBool(
			aclEntry[mkResourceVirtualEnvironmentGroupACLPropagate].(bool),
		)

		aclBody := &access.ACLUpdateRequestBody{
			Delete:    &aclDelete,
			Groups:    []string{groupID},
			Path:      aclEntry[mkResourceVirtualEnvironmentGroupACLPath].(string),
			Propagate: &aclPropagate,
			Roles:     []string{aclEntry[mkResourceVirtualEnvironmentGroupACLRoleID].(string)},
		}

		err := api.Access().UpdateACL(ctx, aclBody)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return groupRead(ctx, d, m)
}

func groupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	groupID := d.Id()
	group, err := api.Access().GetGroup(ctx, groupID)
	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")

			return nil
		}
		return diag.FromErr(err)
	}

	acl, err := api.Access().GetACL(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	var aclParsed []interface{}

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

	err = d.Set(mkResourceVirtualEnvironmentGroupACL, aclParsed)
	diags = append(diags, diag.FromErr(err)...)

	if group.Comment != nil {
		err = d.Set(mkResourceVirtualEnvironmentGroupComment, group.Comment)
	} else {
		err = d.Set(mkResourceVirtualEnvironmentGroupComment, "")
	}
	diags = append(diags, diag.FromErr(err)...)

	err = d.Set(mkResourceVirtualEnvironmentGroupMembers, group.Members)
	diags = append(diags, diag.FromErr(err)...)

	return diags
}

func groupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	comment := d.Get(mkResourceVirtualEnvironmentGroupComment).(string)
	groupID := d.Id()

	body := &access.GroupUpdateRequestBody{
		Comment: &comment,
	}

	err = api.Access().UpdateGroup(ctx, groupID, body)
	if err != nil {
		return diag.FromErr(err)
	}

	aclArgOld, aclArg := d.GetChange(mkResourceVirtualEnvironmentGroupACL)
	aclParsedOld := aclArgOld.(*schema.Set).List()

	for _, v := range aclParsedOld {
		aclDelete := types.CustomBool(true)
		aclEntry := v.(map[string]interface{})
		aclPropagate := types.CustomBool(
			aclEntry[mkResourceVirtualEnvironmentGroupACLPropagate].(bool),
		)

		aclBody := &access.ACLUpdateRequestBody{
			Delete:    &aclDelete,
			Groups:    []string{groupID},
			Path:      aclEntry[mkResourceVirtualEnvironmentGroupACLPath].(string),
			Propagate: &aclPropagate,
			Roles:     []string{aclEntry[mkResourceVirtualEnvironmentGroupACLRoleID].(string)},
		}

		err := api.Access().UpdateACL(ctx, aclBody)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	aclParsed := aclArg.(*schema.Set).List()

	for _, v := range aclParsed {
		aclDelete := types.CustomBool(false)
		aclEntry := v.(map[string]interface{})
		aclPropagate := types.CustomBool(
			aclEntry[mkResourceVirtualEnvironmentGroupACLPropagate].(bool),
		)

		aclBody := &access.ACLUpdateRequestBody{
			Delete:    &aclDelete,
			Groups:    []string{groupID},
			Path:      aclEntry[mkResourceVirtualEnvironmentGroupACLPath].(string),
			Propagate: &aclPropagate,
			Roles:     []string{aclEntry[mkResourceVirtualEnvironmentGroupACLRoleID].(string)},
		}

		err := api.Access().UpdateACL(ctx, aclBody)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return groupRead(ctx, d, m)
}

func groupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	aclParsed := d.Get(mkResourceVirtualEnvironmentGroupACL).(*schema.Set).List()
	groupID := d.Id()

	for _, v := range aclParsed {
		aclDelete := types.CustomBool(true)
		aclEntry := v.(map[string]interface{})
		aclPropagate := types.CustomBool(
			aclEntry[mkResourceVirtualEnvironmentGroupACLPropagate].(bool),
		)

		aclBody := &access.ACLUpdateRequestBody{
			Delete:    &aclDelete,
			Groups:    []string{groupID},
			Path:      aclEntry[mkResourceVirtualEnvironmentGroupACLPath].(string),
			Propagate: &aclPropagate,
			Roles:     []string{aclEntry[mkResourceVirtualEnvironmentGroupACLRoleID].(string)},
		}

		err = api.Access().UpdateACL(ctx, aclBody)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = api.Access().DeleteGroup(ctx, groupID)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")

			return nil
		}
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
