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
	dvIPSetCIDRComment = ""
	dvIPSetCIDRNoMatch = false

	mkIPSetName        = "name"
	mkIPSetCIDR        = "cidr"
	mkIPSetCIDRName    = "name"
	mkIPSetCIDRComment = "comment"
	mkIPSetCIDRNoMatch = "nomatch"
)

func IPSet() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkIPSetName: {
				Type:        schema.TypeString,
				Description: "IPSet name",
				Required:    true,
				ForceNew:    false,
			},
			mkIPSetCIDR: {
				Type:        schema.TypeList,
				Description: "List of IP or Networks",
				Optional:    true,
				ForceNew:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkIPSetCIDRName: {
							Type:        schema.TypeString,
							Description: "Network/IP specification in CIDR format",
							Required:    true,
							ForceNew:    true,
						},
						mkIPSetCIDRNoMatch: {
							Type:        schema.TypeBool,
							Description: "No match this IP/CIDR",
							Optional:    true,
							Default:     dvIPSetCIDRNoMatch,
							ForceNew:    true,
						},
						mkIPSetCIDRComment: {
							Type:        schema.TypeString,
							Description: "IP/CIDR comment",
							Optional:    true,
							Default:     dvIPSetCIDRComment,
							ForceNew:    true,
						},
					},
				},
				MaxItems: 14,
				MinItems: 0,
			},
			mkIPSetCIDRComment: {
				Type:        schema.TypeString,
				Description: "IPSet comment",
				Optional:    true,
				Default:     dvIPSetCIDRComment,
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

	comment := d.Get(mkIPSetCIDRComment).(string)
	name := d.Get(mkIPSetName).(string)

	ipSets := d.Get(mkIPSetCIDR).([]interface{})
	ipSetsArray := make(firewall.IPSetContent, len(ipSets))

	for i, v := range ipSets {
		ipSetMap := v.(map[string]interface{})
		ipSetObject := firewall.IPSetGetResponseData{}

		cidr := ipSetMap[mkIPSetCIDRName].(string)
		noMatch := ipSetMap[mkIPSetCIDRNoMatch].(bool)
		comm := ipSetMap[mkIPSetCIDRComment].(string)

		if comm != "" {
			ipSetObject.Comment = &comm
		}
		ipSetObject.CIDR = cidr

		if noMatch {
			noMatchBool := types.CustomBool(true)
			ipSetObject.NoMatch = &noMatchBool
		}

		ipSetsArray[i] = ipSetObject
	}

	body := &firewall.IPSetCreateRequestBody{
		Comment: comment,
		Name:    name,
	}

	err = veClient.API().Cluster().Firewall().CreateIPSet(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, v := range ipSetsArray {
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
			err = d.Set(mkIPSetName, v.Name)
			diags = append(diags, diag.FromErr(err)...)
			err = d.Set(mkIPSetCIDRComment, v.Comment)
			diags = append(diags, diag.FromErr(err)...)
			break
		}
	}

	ipSet, err := veClient.API().Cluster().Firewall().GetIPSetContent(ctx, name)
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

	for key := range ipSet {
		entry := map[string]interface{}{}

		entry[mkIPSetCIDRName] = ipSet[key].CIDR
		entry[mkIPSetCIDRNoMatch] = ipSet[key].NoMatch
		entry[mkIPSetCIDRComment] = ipSet[key].Comment

		entries = append(entries, entry)
	}

	err = d.Set(mkIPSetCIDR, entries)
	diags = append(diags, diag.FromErr(err)...)
	return diags
}

func ipSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	comment := d.Get(mkIPSetCIDRComment).(string)
	newName := d.Get(mkIPSetName).(string)
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
