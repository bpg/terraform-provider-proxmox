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

	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
)

const (
	dvAliasComment = ""

	mkAliasName    = "name"
	mkAliasCIDR    = "cidr"
	mkAliasComment = "comment"
)

func AliasSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
	}
}

func AliasRead(ctx context.Context, fw firewall.API, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	aliasName := d.Get(mkAliasName).(string)
	alias, err := fw.GetAlias(ctx, aliasName)
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
