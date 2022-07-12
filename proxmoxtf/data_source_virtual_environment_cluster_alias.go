/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	dvDataVirtualEnvironmentClusterAliasComment = ""

	mkDataSourceVirtualEnvironmentClusterAliasName    = "name"
	mkDataSourceVirtualEnvironmentClusterAliasCIDR    = "cidr"
	mkDataSourceVirtualEnvironmentClusterAliasComment = "comment"
)

func dataSourceVirtualEnvironmentClusterAlias() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentClusterAliasName: {
				Type:        schema.TypeString,
				Description: "Alias name",
				Required:    true,
			},
			mkDataSourceVirtualEnvironmentClusterAliasCIDR: {
				Type:        schema.TypeString,
				Description: "IP/CIDR block",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentClusterAliasComment: {
				Type:        schema.TypeString,
				Description: "Alias comment",
				Computed:    true,
			},
		},
		ReadContext: dataSourceVirtualEnvironmentAliasRead,
	}
}

func dataSourceVirtualEnvironmentAliasRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	AliasID := d.Get(mkDataSourceVirtualEnvironmentClusterAliasName).(string)
	Alias, err := veClient.GetAlias(AliasID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(AliasID)

	err = d.Set(mkDataSourceVirtualEnvironmentClusterAliasCIDR, Alias.CIDR)
	diags = append(diags, diag.FromErr(err)...)

	if Alias.Comment != nil {
		err = d.Set(mkDataSourceVirtualEnvironmentClusterAliasComment, Alias.Comment)
	} else {
		err = d.Set(mkDataSourceVirtualEnvironmentClusterAliasComment, dvDataVirtualEnvironmentClusterAliasComment)
	}
	diags = append(diags, diag.FromErr(err)...)

	return diags
}
