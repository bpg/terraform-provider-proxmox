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

	fw "github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/datasource/firewall"
)

// FirewallAlias returns a resource that represents a single firewall alias.
func FirewallAlias() *schema.Resource {
	return &schema.Resource{
		Schema:      firewall.AliasSchema(),
		ReadContext: invokeFirewallAPI(firewall.AliasRead),
	}
}

// FirewallAliases returns a resource that represents firewall aliases.
func FirewallAliases() *schema.Resource {
	return &schema.Resource{
		Schema:      firewall.AliasesSchema(),
		ReadContext: invokeFirewallAPI(firewall.AliasesRead),
	}
}

// FirewallIPSet returns a resource that represents a single firewall IP set.
func FirewallIPSet() *schema.Resource {
	return &schema.Resource{
		Schema:      firewall.IPSetSchema(),
		ReadContext: invokeFirewallAPI(firewall.IPSetRead),
	}
}

// FirewallIPSets returns a resource that represents firewall IP sets.
func FirewallIPSets() *schema.Resource {
	return &schema.Resource{
		Schema:      firewall.IPSetsSchema(),
		ReadContext: invokeFirewallAPI(firewall.IPSetsRead),
	}
}

func invokeFirewallAPI(
	f func(context.Context, fw.API, *schema.ResourceData) diag.Diagnostics,
) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		config := m.(proxmoxtf.ProviderConfiguration)

		api, err := config.GetAPI()
		if err != nil {
			return diag.FromErr(err)
		}

		return f(ctx, api.Cluster().Firewall(), d)
	}
}
