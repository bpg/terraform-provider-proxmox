/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	clusterfirewall "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
)

const (
	dvSecurityGroupComment = ""

	mkSecurityGroupName    = "name"
	mkSecurityGroupComment = "comment"
)

// SecurityGroup returns a resource to manage security groups.
func SecurityGroup() *schema.Resource {
	s := map[string]*schema.Schema{
		mkSecurityGroupName: {
			Type:        schema.TypeString,
			Description: "Security group name",
			Required:    true,
			ForceNew:    false,
		},
		mkSecurityGroupComment: {
			Type:        schema.TypeString,
			Description: "Security group comment",
			Optional:    true,
			Default:     dvSecurityGroupComment,
		},
	}

	structure.MergeSchema(s, firewall.Rules().Schema)

	return &schema.Resource{
		Schema:        s,
		CreateContext: selectFirewallAPI(SecurityGroupCreate),
		ReadContext:   selectFirewallAPI(SecurityGroupRead),
		UpdateContext: selectFirewallAPI(SecurityGroupUpdate),
		DeleteContext: selectFirewallAPI(SecurityGroupDelete),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// SecurityGroupCreate creates a new security group.
func SecurityGroupCreate(ctx context.Context, api clusterfirewall.API, d *schema.ResourceData) diag.Diagnostics {
	comment := d.Get(mkSecurityGroupComment).(string)
	name := d.Get(mkSecurityGroupName).(string)

	body := &clusterfirewall.GroupCreateRequestBody{
		Comment: &comment,
		Group:   name,
	}

	err := api.CreateGroup(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	diags := firewall.RulesCreate(ctx, api.SecurityGroup(name), d)
	if diags.HasError() {
		return diags
	}

	d.SetId(name)

	return SecurityGroupRead(ctx, api, d)
}

// SecurityGroupRead reads the security group from the API and updates the state.
func SecurityGroupRead(ctx context.Context, api clusterfirewall.API, d *schema.ResourceData) diag.Diagnostics {
	diags := diag.Diagnostics{}

	name := d.Id()

	allGroups, err := api.ListGroups(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, v := range allGroups {
		if v.Group == name {
			err = d.Set(mkSecurityGroupName, v.Group)
			diags = append(diags, diag.FromErr(err)...)
			err = d.Set(mkSecurityGroupComment, v.Comment)
			diags = append(diags, diag.FromErr(err)...)

			break
		}
	}

	if diags.HasError() {
		return diags
	}

	return firewall.RulesRead(ctx, api.SecurityGroup(name), d)
}

// SecurityGroupUpdate updates a security group.
func SecurityGroupUpdate(ctx context.Context, api clusterfirewall.API, d *schema.ResourceData) diag.Diagnostics {
	comment := d.Get(mkSecurityGroupComment).(string)
	newName := d.Get(mkSecurityGroupName).(string)
	previousName := d.Id()

	body := &clusterfirewall.GroupUpdateRequestBody{
		Group:   newName,
		ReName:  &previousName,
		Comment: &comment,
	}

	err := api.UpdateGroup(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	diags := firewall.RulesUpdate(ctx, api.SecurityGroup(previousName), d)
	if diags.HasError() {
		return diags
	}

	d.SetId(newName)

	return SecurityGroupRead(ctx, api, d)
}

// SecurityGroupDelete deletes a security group.
func SecurityGroupDelete(ctx context.Context, api clusterfirewall.API, d *schema.ResourceData) diag.Diagnostics {
	group := d.Id()

	diags := firewall.RulesDelete(ctx, api.SecurityGroup(group), d)
	if diags.HasError() {
		return diags
	}

	err := api.DeleteGroup(ctx, group)
	if err != nil {
		if strings.Contains(err.Error(), "no such security group") {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
