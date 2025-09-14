/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"
	"fmt"
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

// findActualGroupName finds the actual group name as stored by Proxmox (case-insensitive lookup).
func findActualGroupName(ctx context.Context, api clusterfirewall.API, targetName string) (string, error) {
	allGroups, err := api.ListGroups(ctx)
	if err != nil {
		return "", fmt.Errorf("error retrieving security groups: %w", err)
	}

	for _, group := range allGroups {
		if strings.EqualFold(group.Group, targetName) {
			return group.Group, nil
		}
	}

	return "", fmt.Errorf("security group '%s' not found after creation", targetName)
}

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

	actualName, err := findActualGroupName(ctx, api, name)
	if err != nil {
		return diag.FromErr(err)
	}

	diags := firewall.RulesCreate(ctx, api.SecurityGroup(actualName), d)
	if diags.HasError() {
		return diags
	}

	d.SetId(actualName)

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

	var foundGroup *clusterfirewall.GroupListResponseData

	for _, v := range allGroups {
		if strings.EqualFold(v.Group, name) {
			foundGroup = v
			break
		}
	}

	if foundGroup == nil {
		d.SetId("")
		return nil
	}

	err = d.Set(mkSecurityGroupName, foundGroup.Group)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkSecurityGroupComment, foundGroup.Comment)
	diags = append(diags, diag.FromErr(err)...)

	if diags.HasError() {
		return diags
	}

	return firewall.RulesRead(ctx, api.SecurityGroup(foundGroup.Group), d)
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

	actualName, err := findActualGroupName(ctx, api, newName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(actualName)

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
