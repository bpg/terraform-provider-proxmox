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

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/validator"
)

const (
	dvLogRatelimiEnabled = true
	dvLogRatelimitBurst  = 5
	dvLogRatelimitRate   = "1/second"
	dvPolicyIn           = "DROP"
	dvPolicyOut          = "ACCEPT"

	mkEBTables            = "ebtables"
	mkEnabled             = "enabled"
	mkLogRatelimit        = "log_ratelimit"
	mkLogRatelimitEnabled = "enabled"
	mkLogRatelimitBurst   = "burst"
	mkLogRatelimitRate    = "rate"
	mkPolicyIn            = "input_policy"
	mkPolicyOut           = "output_policy"
)

// Firewall returns a resource to manage firewall options.
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
				Optional:    true,
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkLogRatelimitEnabled: {
							Type:        schema.TypeBool,
							Description: "Enable or disable log ratelimiting",
							Optional:    true,
							Default:     dvLogRatelimiEnabled,
						},
						mkLogRatelimitBurst: {
							Type:        schema.TypeInt,
							Description: "Initial burst of packages which will always get logged before the rate is applied",
							Optional:    true,
							Default:     dvLogRatelimitBurst,
						},
						mkLogRatelimitRate: {
							Type:             schema.TypeString,
							Description:      "Frequency with which the burst bucket gets refilled",
							Optional:         true,
							Default:          dvLogRatelimitRate,
							ValidateDiagFunc: validator.FirewallRate(),
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkPolicyIn: {
				Type:             schema.TypeString,
				Description:      "Default policy for incoming traffic",
				Optional:         true,
				Default:          dvPolicyIn,
				ValidateDiagFunc: validator.FirewallPolicy(),
			},
			mkPolicyOut: {
				Type:             schema.TypeString,
				Description:      "Default policy for outgoing traffic",
				Optional:         true,
				Default:          dvPolicyOut,
				ValidateDiagFunc: validator.FirewallPolicy(),
			},
		},
		CreateContext: selectFirewallAPI(firewallCreate),
		ReadContext:   selectFirewallAPI(firewallRead),
		UpdateContext: selectFirewallAPI(firewallUpdate),
		DeleteContext: selectFirewallAPI(firewallDelete),
	}
}

func firewallCreate(ctx context.Context, api firewall.API, d *schema.ResourceData) diag.Diagnostics {
	diags := setOptions(ctx, api, d)
	if diags.HasError() {
		return diags
	}

	return firewallRead(ctx, api, d)
}

func setOptions(ctx context.Context, api firewall.API, d *schema.ResourceData) diag.Diagnostics {
	policyIn := d.Get(mkPolicyIn).(string)
	policyOut := d.Get(mkPolicyOut).(string)
	body := &firewall.OptionsPutRequestBody{
		PolicyIn:  &policyIn,
		PolicyOut: &policyOut,
	}

	logRatelimit := d.Get(mkLogRatelimit).([]interface{})
	if len(logRatelimit) > 0 {
		m := logRatelimit[0].(map[string]interface{})
		burst := m[mkLogRatelimitBurst].(int)
		rate := m[mkLogRatelimitRate].(string)
		rl := firewall.CustomLogRateLimit{
			Enable: types.CustomBool(m[mkLogRatelimitEnabled].(bool)),
			Burst:  &burst,
			Rate:   &rate,
		}
		body.LogRateLimit = &rl
	}

	ebtablesBool := types.CustomBool(d.Get(mkEBTables).(bool))
	body.EBTables = &ebtablesBool

	enabledBool := types.CustomBool(d.Get(mkEnabled).(bool))
	body.Enable = &enabledBool

	err := api.SetGlobalOptions(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("cluster-firewall")

	return nil
}

func firewallRead(ctx context.Context, api firewall.API, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	options, err := api.GetGlobalOptions(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if options.EBTables != nil {
		err = d.Set(mkEBTables, *options.EBTables)
		diags = append(diags, diag.FromErr(err)...)
	}

	if options.Enable != nil {
		err = d.Set(mkEnabled, *options.Enable)
		diags = append(diags, diag.FromErr(err)...)
	}

	if options.LogRateLimit != nil {
		err = d.Set(mkLogRatelimit, []interface{}{
			map[string]interface{}{
				mkLogRatelimitEnabled: options.LogRateLimit.Enable,
				mkLogRatelimitBurst:   *options.LogRateLimit.Burst,
				mkLogRatelimitRate:    *options.LogRateLimit.Rate,
			},
		})
		diags = append(diags, diag.FromErr(err)...)
	}

	if options.PolicyIn != nil {
		err = d.Set(mkPolicyIn, *options.PolicyIn)
		diags = append(diags, diag.FromErr(err)...)
	}

	if options.PolicyOut != nil {
		err = d.Set(mkPolicyOut, *options.PolicyOut)
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

func firewallUpdate(ctx context.Context, api firewall.API, d *schema.ResourceData) diag.Diagnostics {
	diags := setOptions(ctx, api, d)
	if diags.HasError() {
		return diags
	}

	return firewallRead(ctx, api, d)
}

func firewallDelete(_ context.Context, _ firewall.API, d *schema.ResourceData) diag.Diagnostics {
	d.SetId("")

	return nil
}

func selectFirewallAPI(
	f func(context.Context, firewall.API, *schema.ResourceData) diag.Diagnostics,
) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		config := m.(proxmoxtf.ProviderConfiguration)
		veClient, err := config.GetVEClient()
		if err != nil {
			return diag.FromErr(err)
		}

		api := veClient.API().Cluster().Firewall()

		return f(ctx, api, d)
	}
}
