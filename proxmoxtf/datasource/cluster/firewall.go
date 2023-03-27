/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	fw "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/datasource/firewall"
)

func FirewallAlias() *schema.Resource {
	return &schema.Resource{
		Schema:      firewall.AliasSchema(),
		ReadContext: invokeFirewallAPI(firewall.AliasRead),
	}
}

func FirewallAliases() *schema.Resource {
	return &schema.Resource{
		Schema:      firewall.AliasesSchema(),
		ReadContext: invokeFirewallAPI(firewall.AliasesRead),
	}
}

func FirewallIPSet() *schema.Resource {
	return &schema.Resource{
		Schema:      firewall.IPSetSchema(),
		ReadContext: invokeFirewallAPI(firewall.IPSetRead),
	}
}

func FirewallIPSets() *schema.Resource {
	return &schema.Resource{
		Schema:      firewall.IPSetsSchema(),
		ReadContext: invokeFirewallAPI(firewall.IPSetsRead),
	}
}

func invokeFirewallAPI(
	f func(context.Context, *fw.API, *schema.ResourceData) diag.Diagnostics,
) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		config := m.(proxmoxtf.ProviderConfiguration)
		veClient, err := config.GetVEClient()
		if err != nil {
			return diag.FromErr(err)
		}

		return f(ctx, veClient.API().Cluster().Firewall(), d)
	}
}
