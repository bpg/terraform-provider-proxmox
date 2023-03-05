/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

//nolint:dupl
package firewall

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	mkAliasesAliasNames = "alias_names"
)

func Aliases() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkAliasesAliasNames: {
				Type:        schema.TypeList,
				Description: "Alias Names",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		ReadContext: aliasesRead,
	}
}

func aliasesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	list, err := veClient.API().Cluster().Firewall().ListAliases(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	aliasNames := make([]interface{}, len(list))

	for i, v := range list {
		aliasNames[i] = v.Name
	}

	d.SetId("aliases")

	err = d.Set(mkAliasesAliasNames, aliasNames)

	return diag.FromErr(err)
}
