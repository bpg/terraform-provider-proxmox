/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	proxmoxapi "github.com/bpg/terraform-provider-proxmox/proxmox/api"
	clusterfirewall "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
)

const (
	dvSecurityGroupComment = ""

	mkSecurityGroupName    = "name"
	mkSecurityGroupComment = "comment"
)

// validateSecurityGroupName rejects names with uppercase letters. PVE lowercases security group names on
// creation, so a mixed-case name would make the follow-up rule calls fail with a "no such security group" error.
func validateSecurityGroupName(i any, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if v != strings.ToLower(v) {
		return nil, []error{fmt.Errorf(
			"%q must be lowercase, got %q: PVE lowercases security group names, use %q instead",
			k, v, strings.ToLower(v),
		)}
	}

	return nil, nil
}

// SecurityGroup returns a resource to manage security groups.
func SecurityGroup() *schema.Resource {
	s := map[string]*schema.Schema{
		mkSecurityGroupName: {
			Type:             schema.TypeString,
			Description:      "Security group name. Must be lowercase, because PVE lowercases security group names.",
			Required:         true,
			ForceNew:         false,
			ValidateDiagFunc: validation.ToDiagFunc(validateSecurityGroupName),
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
		if !errors.Is(err, proxmoxapi.ErrResourceAlreadyExists) {
			return diag.FromErr(err)
		}

		// adopt existing group into state; next plan will reconcile.
		d.SetId(name)

		return SecurityGroupRead(ctx, api, d)
	}

	d.SetId(name)

	diags := firewall.RulesCreate(ctx, api.SecurityGroup(name), d)
	if diags.HasError() {
		return diags
	}

	// RulesCreate overwrites d.SetId with a rules hash; restore the group name.
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

	found := false

	for _, v := range allGroups {
		if v.Group == name {
			found = true

			err = d.Set(mkSecurityGroupName, v.Group)
			diags = append(diags, diag.FromErr(err)...)
			err = d.Set(mkSecurityGroupComment, v.Comment)
			diags = append(diags, diag.FromErr(err)...)

			break
		}
	}

	if !found {
		d.SetId("")
		return nil
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
