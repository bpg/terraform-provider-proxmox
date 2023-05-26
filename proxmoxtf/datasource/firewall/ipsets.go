/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
)

const (
	mkIPSetsIPSetNames = "ipset_names"
)

// IPSetsSchema defines the schema for the IP sets.
func IPSetsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		mkIPSetsIPSetNames: {
			Type:        schema.TypeList,
			Description: "IPSet Names",
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

// IPSetsRead reads the IP sets.
func IPSetsRead(ctx context.Context, fw firewall.API, d *schema.ResourceData) diag.Diagnostics {
	list, err := fw.ListIPSets(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	ipSetNames := make([]interface{}, len(list))

	for i, v := range list {
		ipSetNames[i] = v.Name
	}

	d.SetId(uuid.New().String())

	err = d.Set(mkIPSetsIPSetNames, ipSetNames)

	return diag.FromErr(err)
}
