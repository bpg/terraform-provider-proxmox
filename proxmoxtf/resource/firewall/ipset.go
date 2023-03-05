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

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	dvResourceVirtualEnvironmentFirewallIPSetCIDRComment = ""
	dvResourceVirtualEnvironmentFirewallIPSetCIDRNoMatch = false

	mkResourceVirtualEnvironmentFirewallIPSetName        = "name"
	mkResourceVirtualEnvironmentFirewallIPSetCIDR        = "cidr"
	mkResourceVirtualEnvironmentFirewallIPSetCIDRName    = "name"
	mkResourceVirtualEnvironmentFirewallIPSetCIDRComment = "comment"
	mkResourceVirtualEnvironmentFirewallIPSetCIDRNoMatch = "nomatch"
)

func IPSet() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentFirewallIPSetName: {
				Type:        schema.TypeString,
				Description: "IPSet name",
				Required:    true,
				ForceNew:    false,
			},
			mkResourceVirtualEnvironmentFirewallIPSetCIDR: {
				Type:        schema.TypeList,
				Description: "List of IP or Networks",
				Optional:    true,
				ForceNew:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentFirewallIPSetCIDRName: {
							Type:        schema.TypeString,
							Description: "Network/IP specification in CIDR format",
							Required:    true,
							ForceNew:    true,
						},
						mkResourceVirtualEnvironmentFirewallIPSetCIDRNoMatch: {
							Type:        schema.TypeBool,
							Description: "No match this IP/CIDR",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentFirewallIPSetCIDRNoMatch,
							ForceNew:    true,
						},
						mkResourceVirtualEnvironmentFirewallIPSetCIDRComment: {
							Type:        schema.TypeString,
							Description: "IP/CIDR comment",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentFirewallIPSetCIDRComment,
							ForceNew:    true,
						},
					},
				},
				MaxItems: 14,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentFirewallIPSetCIDRComment: {
				Type:        schema.TypeString,
				Description: "IPSet comment",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentFirewallIPSetCIDRComment,
			},
		},
		CreateContext: ipSetCreate,
		ReadContext:   ipSetRead,
		UpdateContext: ipSetUpdate,
		DeleteContext: ipSetDelete,
	}
}

func ipSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	comment := d.Get(mkResourceVirtualEnvironmentFirewallIPSetCIDRComment).(string)
	name := d.Get(mkResourceVirtualEnvironmentFirewallIPSetName).(string)

	IPSets := d.Get(mkResourceVirtualEnvironmentFirewallIPSetCIDR).([]interface{})
	IPSetsArray := make(firewall.IPSetContent, len(IPSets))

	for i, v := range IPSets {
		IPSetMap := v.(map[string]interface{})
		IPSetObject := firewall.IPSetGetResponseData{}

		cidr := IPSetMap[mkResourceVirtualEnvironmentFirewallIPSetCIDRName].(string)
		noMatch := IPSetMap[mkResourceVirtualEnvironmentFirewallIPSetCIDRNoMatch].(bool)
		comm := IPSetMap[mkResourceVirtualEnvironmentFirewallIPSetCIDRComment].(string)

		if comm != "" {
			IPSetObject.Comment = &comm
		}
		IPSetObject.CIDR = cidr

		if noMatch {
			noMatchBool := types.CustomBool(true)
			IPSetObject.NoMatch = &noMatchBool
		}

		IPSetsArray[i] = IPSetObject
	}

	body := &firewall.IPSetCreateRequestBody{
		Comment: comment,
		Name:    name,
	}

	err = veClient.API().Cluster().Firewall().CreateIPSet(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, v := range IPSetsArray {
		err = veClient.API().Cluster().Firewall().AddCIDRToIPSet(ctx, name, v)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(name)
	return ipSetRead(ctx, d, m)
}

func ipSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Id()

	allIPSets, err := veClient.API().Cluster().Firewall().ListIPSets(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, v := range allIPSets {
		if v.Name == name {
			err = d.Set(mkResourceVirtualEnvironmentFirewallIPSetName, v.Name)
			diags = append(diags, diag.FromErr(err)...)
			err = d.Set(mkResourceVirtualEnvironmentFirewallIPSetCIDRComment, v.Comment)
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	IPSet, err := veClient.API().Cluster().Firewall().GetIPSetContent(ctx, name)
	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")
			return nil
		}
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	//nolint:prealloc
	var entries []interface{}

	for key := range IPSet {
		entry := map[string]interface{}{}

		entry[mkResourceVirtualEnvironmentFirewallIPSetCIDRName] = IPSet[key].CIDR
		entry[mkResourceVirtualEnvironmentFirewallIPSetCIDRNoMatch] = IPSet[key].NoMatch
		entry[mkResourceVirtualEnvironmentFirewallIPSetCIDRComment] = IPSet[key].Comment

		entries = append(entries, entry)
	}

	err = d.Set(mkResourceVirtualEnvironmentFirewallIPSetCIDR, entries)
	diags = append(diags, diag.FromErr(err)...)
	return diags
}

func ipSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	comment := d.Get(mkResourceVirtualEnvironmentFirewallIPSetCIDRComment).(string)
	newName := d.Get(mkResourceVirtualEnvironmentFirewallIPSetName).(string)
	previousName := d.Id()

	body := &firewall.IPSetUpdateRequestBody{
		ReName:  previousName,
		Name:    newName,
		Comment: &comment,
	}

	err = veClient.API().Cluster().Firewall().UpdateIPSet(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newName)

	return ipSetRead(ctx, d, m)
}

func ipSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Id()

	IPSetContent, err := veClient.API().Cluster().Firewall().GetIPSetContent(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	// PVE requires content of IPSet be cleared before removal
	if len(IPSetContent) > 0 {
		for _, IPSet := range IPSetContent {
			err = veClient.API().Cluster().Firewall().DeleteIPSetContent(ctx, name, IPSet.CIDR)
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	if diags.HasError() {
		return diags
	}

	err = veClient.API().Cluster().Firewall().DeleteIPSet(ctx, name)

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
