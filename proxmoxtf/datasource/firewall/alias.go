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
	dvDataVirtualEnvironmentFirewallAliasComment = ""

	mkDataSourceVirtualEnvironmentFirewallAliasName    = "name"
	mkDataSourceVirtualEnvironmentFirewallAliasCIDR    = "cidr"
	mkDataSourceVirtualEnvironmentFirewallAliasComment = "comment"
)

func Alias() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentFirewallAliasName: {
				Type:        schema.TypeString,
				Description: "Alias name",
				Required:    true,
			},
			mkDataSourceVirtualEnvironmentFirewallAliasCIDR: {
				Type:        schema.TypeString,
				Description: "IP/CIDR block",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentFirewallAliasComment: {
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

	AliasID := d.Get(mkDataSourceVirtualEnvironmentFirewallAliasName).(string)
	Alias, err := veClient.API().Cluster().Firewall().GetAlias(ctx, AliasID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(AliasID)

	err = d.Set(mkDataSourceVirtualEnvironmentFirewallAliasCIDR, Alias.CIDR)
	diags = append(diags, diag.FromErr(err)...)

	if Alias.Comment != nil {
		err = d.Set(mkDataSourceVirtualEnvironmentFirewallAliasComment, Alias.Comment)
	} else {
		err = d.Set(mkDataSourceVirtualEnvironmentFirewallAliasComment, dvDataVirtualEnvironmentFirewallAliasComment)
	}
	diags = append(diags, diag.FromErr(err)...)

	return diags
}
