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
	dvAliasComment = ""

	mkAliasName    = "name"
	mkAliasCIDR    = "cidr"
	mkAliasComment = "comment"
)

func Alias() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkAliasName: {
				Type:        schema.TypeString,
				Description: "Alias name",
				Required:    true,
			},
			mkAliasCIDR: {
				Type:        schema.TypeString,
				Description: "IP/CIDR block",
				Computed:    true,
			},
			mkAliasComment: {
				Type:        schema.TypeString,
				Description: "Alias comment",
				Computed:    true,
			},
		},
		ReadContext: aliasRead,
	}
}

func aliasRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	aliasName := d.Get(mkAliasName).(string)
	alias, err := veClient.API().Cluster().Firewall().GetAlias(ctx, aliasName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aliasName)

	err = d.Set(mkAliasCIDR, alias.CIDR)
	diags = append(diags, diag.FromErr(err)...)

	if alias.Comment != nil {
		err = d.Set(mkAliasComment, alias.Comment)
	} else {
		err = d.Set(mkAliasComment, dvAliasComment)
	}
	diags = append(diags, diag.FromErr(err)...)

	return diags
}
