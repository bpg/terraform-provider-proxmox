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

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/firewall"
)

const (
	mkAliasesAliasNames = "alias_names"
)

func AliasesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		mkAliasesAliasNames: {
			Type:        schema.TypeList,
			Description: "Alias Names",
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func AliasesRead(ctx context.Context, fw *firewall.API, d *schema.ResourceData) diag.Diagnostics {
	list, err := fw.ListAliases(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	aliasNames := make([]interface{}, len(list))

	for i, v := range list {
		aliasNames[i] = v.Name
	}

	d.SetId(uuid.New().String())

	err = d.Set(mkAliasesAliasNames, aliasNames)

	return diag.FromErr(err)
}
