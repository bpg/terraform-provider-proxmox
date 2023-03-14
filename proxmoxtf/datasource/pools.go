/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datasource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	mkDataSourceVirtualEnvironmentPoolsPoolIDs = "pool_ids"
)

func Pools() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentPoolsPoolIDs: {
				Type:        schema.TypeList,
				Description: "The pool ids",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		ReadContext: poolsRead,
	}
}

func poolsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	list, err := veClient.ListPools(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	poolIDs := make([]interface{}, len(list))

	for i, v := range list {
		poolIDs[i] = v.ID
	}

	d.SetId("pools")

	err = d.Set(mkDataSourceVirtualEnvironmentPoolsPoolIDs, poolIDs)

	return diag.FromErr(err)
}
