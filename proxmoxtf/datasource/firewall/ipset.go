/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/firewall"
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

func IPSetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		mkIPSetName: {
			Type:        schema.TypeString,
			Description: "IPSet name",
			Required:    true,
		},
		mkIPSetCIDRComment: {
			Type:        schema.TypeString,
			Description: "IPSet comment",
			Computed:    true,
		},
		mkIPSetCIDR: {
			Type:        schema.TypeList,
			Description: "List of IP or Networks",
			Computed:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					mkIPSetCIDRName: {
						Type:        schema.TypeString,
						Description: "Network/IP specification in CIDR format",
						Computed:    true,
					},
					mkIPSetCIDRNoMatch: {
						Type:        schema.TypeBool,
						Description: "No match this IP/CIDR",
						Computed:    true,
					},
					mkIPSetCIDRComment: {
						Type:        schema.TypeString,
						Description: "IP/CIDR comment",
						Computed:    true,
					},
				},
			},
		},
	}
}

func IPSetRead(ctx context.Context, fw *firewall.API, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	ipSetName := d.Get(mkIPSetName).(string)
	d.SetId(ipSetName)

	ipSetList, err := fw.ListIPSets(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, ipSet := range ipSetList {
		if ipSet.Name == ipSetName {
			if ipSet.Comment != nil {
				err = d.Set(mkIPSetCIDRComment, ipSet.Comment)
			} else {
				err = d.Set(mkIPSetCIDRComment, dvIPSetCIDRComment)
			}
			diags = append(diags, diag.FromErr(err)...)
			break
		}
	}

	content, err := fw.GetIPSetContent(ctx, ipSetName)
	if err != nil {
		return diag.FromErr(err)
	}

	//nolint:prealloc
	var cidrs []interface{}

	for _, v := range content {
		cirdEntry := map[string]interface{}{}

		cirdEntry[mkIPSetCIDRName] = v.CIDR

		if v.NoMatch != nil {
			cirdEntry[mkIPSetCIDRNoMatch] = bool(*v.NoMatch)
		} else {
			cirdEntry[mkIPSetCIDRNoMatch] = dvIPSetCIDRNoMatch
		}

		if v.Comment != nil {
			cirdEntry[mkIPSetCIDRComment] = v.Comment
		} else {
			cirdEntry[mkIPSetCIDRComment] = dvIPSetCIDRComment
		}

		cidrs = append(cidrs, cirdEntry)
	}

	err = d.Set(mkIPSetCIDR, cidrs)
	diags = append(diags, diag.FromErr(err)...)

	return diags
}
