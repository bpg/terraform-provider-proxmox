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
							Description: "IPSet comment",
							Computed:    true,
						},
					},
				},
			},
		},

		ReadContext: ipSetRead,
	}
}

func ipSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	ipSetName := d.Get(mkIPSetName).(string)
	d.SetId(ipSetName)

	ipSetList, err := veClient.API().Cluster().Firewall().ListIPSets(ctx)
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

	content, err := veClient.API().Cluster().Firewall().GetIPSetContent(ctx, ipSetName)
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
