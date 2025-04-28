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

	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
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

// IPSet returns a resource to manage IP sets.
func IPSet() *schema.Resource {
	s := map[string]*schema.Schema{
		mkIPSetName: {
			Type:        schema.TypeString,
			Description: "IPSet name",
			Required:    true,
			ForceNew:    false,
		},
		mkIPSetCIDRComment: {
			Type:        schema.TypeString,
			Description: "IPSet comment",
			Optional:    true,
			Default:     dvIPSetCIDRComment,
		},
		mkIPSetCIDR: {
			Type:        schema.TypeList,
			Description: "List of IP or Networks",
			Optional:    true,
			ForceNew:    true,
			DefaultFunc: func() (interface{}, error) {
				return []interface{}{}, nil
			},
			DiffSuppressFunc: structure.SuppressIfListsOfMapsAreEqualIgnoringOrderByKey(mkIPSetCIDRName),
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
		},
	}

	structure.MergeSchema(s, selectorSchema())

	return &schema.Resource{
		Schema:        s,
		CreateContext: selectFirewallAPI(ipSetCreate),
		ReadContext:   selectFirewallAPI(ipSetRead),
		UpdateContext: selectFirewallAPI(ipSetUpdate),
		DeleteContext: selectFirewallAPI(ipSetDelete),
	}
}

func ipSetCreate(ctx context.Context, api firewall.API, d *schema.ResourceData) diag.Diagnostics {
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

	err := api.CreateIPSet(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, v := range ipSetsArray {
		err = api.AddCIDRToIPSet(ctx, name, v)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(name)

	return ipSetRead(ctx, api, d)
}

func ipSetRead(ctx context.Context, api firewall.API, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	name := d.Id()

	allIPSets, err := api.ListIPSets(ctx)
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

	ipSet, err := api.GetIPSetContent(ctx, name)
	if err != nil {
		if strings.Contains(err.Error(), "no such IPSet") {
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

func ipSetUpdate(ctx context.Context, api firewall.API, d *schema.ResourceData) diag.Diagnostics {
	comment := d.Get(mkIPSetCIDRComment).(string)
	newName := d.Get(mkIPSetName).(string)
	previousName := d.Id()

	body := &firewall.IPSetUpdateRequestBody{
		ReName:  previousName,
		Name:    newName,
		Comment: &comment,
	}

	err := api.UpdateIPSet(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newName)

	return ipSetRead(ctx, api, d)
}

func ipSetDelete(ctx context.Context, api firewall.API, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	name := d.Id()

	IPSetContent, err := api.GetIPSetContent(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	// PVE requires content of IPSet be cleared before removal
	if len(IPSetContent) > 0 {
		for _, ipSet := range IPSetContent {
			err = api.DeleteIPSetContent(ctx, name, ipSet.CIDR)
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	if diags.HasError() {
		return diags
	}

	err = api.DeleteIPSet(ctx, name)
	if err != nil {
		if strings.Contains(err.Error(), "no such IPSet") {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
