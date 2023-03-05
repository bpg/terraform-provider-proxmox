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
	mkIPSetsIPSetNames = "ipset_names"
)

func IPSets() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkIPSetsIPSetNames: {
				Type:        schema.TypeList,
				Description: "IPSet Names",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		ReadContext: ipSetsRead,
	}
}

func ipSetsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	list, err := veClient.API().Cluster().Firewall().ListIPSets(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	ipSetNames := make([]interface{}, len(list))

	for i, v := range list {
		ipSetNames[i] = v.Name
	}

	d.SetId("ipsets")

	err = d.Set(mkIPSetsIPSetNames, ipSetNames)

	return diag.FromErr(err)
}
