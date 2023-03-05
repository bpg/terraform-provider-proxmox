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
	mkDataSourceVirtualEnvironmentClusterAliasesAliasNames = "alias_names"
)

func Aliases() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentClusterAliasesAliasNames: {
				Type:        schema.TypeList,
				Description: "Alias IDs",
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

	aliasIDs := make([]interface{}, len(list))

	for i, v := range list {
		aliasIDs[i] = v.Name
	}

	d.SetId("aliases")

	err = d.Set(mkDataSourceVirtualEnvironmentClusterAliasesAliasNames, aliasIDs)

	return diag.FromErr(err)
}
