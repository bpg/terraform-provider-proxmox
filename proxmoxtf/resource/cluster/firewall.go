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
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/firewall"
)

const (
	dvLogRatelimiEnabled = true
	dvLogRatelimitBurst  = 5
	dvLogRatelimitRate   = 1
	dvPolicyIn           = "DROP"
	dvPolicyOut          = "ACCEPT"

	mkEBTables            = "ebtables"
	mkEnabled             = "enabled"
	mkLogRatelimit        = "log_ratelimit"
	mkLogRatelimitEnabled = "enabled"
	mkLogRatelimitBurst   = "burst"
	mkLogRatelimitRate    = "rate"
	mkPolicyIn            = "policy_in"
	mkPolicyOut           = "policy_out"

	mkRule = "rule"
)

func Firewall() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkEBTables: {
				Type:        schema.TypeBool,
				Description: "Enable ebtables cluster-wide",
				Optional:    true,
			},
			mkEnabled: {
				Type:        schema.TypeBool,
				Description: "Enable or disable the firewall cluster-wide",
				Required:    true,
			},
			mkLogRatelimit: {
				Type:        schema.TypeList,
				Description: "Log ratelimiting settings",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{
						map[string]interface{}{
							mkLogRatelimitEnabled: dvLogRatelimiEnabled,
							mkLogRatelimitBurst:   dvLogRatelimitBurst,
							mkLogRatelimitRate:    dvLogRatelimitRate,
						},
					}, nil
				},
			},
			mkPolicyIn: {
				Type:        schema.TypeString,
				Description: "Default policy for incoming traffic",
				Optional:    true,
				Default:     dvPolicyIn,
			},
			mkPolicyOut: {
				Type:        schema.TypeString,
				Description: "Default policy for outgoing traffic",
				Optional:    true,
				Default:     dvPolicyOut,
			},
			mkRule: {
				Type:        schema.TypeList,
				Description: "List of rules",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				ForceNew: true,
				Elem:     &schema.Resource{Schema: firewall.RuleSchema()},
			},
		},
		CreateContext: invokeFirewallAPI(firewallCreate),
		ReadContext:   invokeFirewallAPI(firewallRead),
		UpdateContext: invokeFirewallAPI(firewallUpdate),
		DeleteContext: invokeFirewallAPI(firewallDelete),
	}
}

func firewallCreate(_ context.Context, _ fw.API, _ *schema.ResourceData) diag.Diagnostics {
	return nil
}

func firewallRead(_ context.Context, _ fw.API, _ *schema.ResourceData) diag.Diagnostics {
	return nil
}

func firewallUpdate(_ context.Context, _ fw.API, _ *schema.ResourceData) diag.Diagnostics {
	return nil
}

func firewallDelete(_ context.Context, _ fw.API, _ *schema.ResourceData) diag.Diagnostics {
	return nil
}

func FirewallAlias() *schema.Resource {
	return &schema.Resource{
		Schema:        firewall.AliasSchema(),
		CreateContext: invokeFirewallAPI(firewall.AliasCreate),
		ReadContext:   invokeFirewallAPI(firewall.AliasRead),
		UpdateContext: invokeFirewallAPI(firewall.AliasUpdate),
		DeleteContext: invokeFirewallAPI(firewall.AliasDelete),
	}
}

func FirewallIPSet() *schema.Resource {
	return &schema.Resource{
		Schema:        firewall.IPSetSchema(),
		CreateContext: invokeFirewallAPI(firewall.IPSetCreate),
		ReadContext:   invokeFirewallAPI(firewall.IPSetRead),
		UpdateContext: invokeFirewallAPI(firewall.IPSetUpdate),
		DeleteContext: invokeFirewallAPI(firewall.IPSetDelete),
	}
}

func FirewallSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Schema:        firewall.SecurityGroupSchema(),
		CreateContext: invokeFirewallAPI(firewall.SecurityGroupCreate),
		ReadContext:   invokeFirewallAPI(firewall.SecurityGroupRead),
		UpdateContext: invokeFirewallAPI(firewall.SecurityGroupUpdate),
		DeleteContext: invokeFirewallAPI(firewall.SecurityGroupDelete),
	}
}

func invokeFirewallAPI(
	f func(context.Context, fw.API, *schema.ResourceData) diag.Diagnostics,
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
