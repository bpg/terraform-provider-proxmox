/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
)

const (
	dvDHCP        = false
	dvEnabled     = false
	dvLogLevelIN  = "nolog"
	dvLogLevelOUT = "nolog"
	dvMACFilter   = true
	dvNDP         = false
	dvPolicyIn    = "DROP"
	dvPolicyOut   = "ACCEPT"
	dvReadv       = true

	mkDHCP        = "dhcp"
	mkEnabled     = "enabled"
	mkIPFilter    = "ipfilter"
	mkLogLevelIN  = "log_level_in"
	mkLogLevelOUT = "log_level_out"
	mkMACFilter   = "macfilter"
	mkNDP         = "ndp"
	mkPolicyIn    = "input_policy"
	mkPolicyOut   = "output_policy"
	mkRadv        = "radv"
)

// Options returns a resource to manage firewall options.
func Options() *schema.Resource {
	s := map[string]*schema.Schema{
		mkDHCP: {
			Type:        schema.TypeBool,
			Description: "Enable DHCP",
			Optional:    true,
			Default:     dvDHCP,
		},
		mkEnabled: {
			Type:        schema.TypeBool,
			Description: "Enable or disable the firewall",
			Optional:    true,
			Default:     dvEnabled,
		},
		mkIPFilter: {
			Type: schema.TypeBool,
			Description: "Enable default IP filters. This is equivalent to adding an empty ipfilter-net<id> ipset " +
				"for every interface. Such ipsets implicitly contain sane default restrictions such as restricting " +
				"IPv6 link local addresses to the one derived from the interface's MAC address. " +
				"For containers the configured IP addresses will be implicitly added.",
			Optional: true,
		},
		mkLogLevelIN: {
			Type:             schema.TypeString,
			Description:      "Log level for incoming traffic.",
			Optional:         true,
			Default:          dvLogLevelIN,
			ValidateDiagFunc: validators.FirewallLogLevel(),
		},
		mkLogLevelOUT: {
			Type:             schema.TypeString,
			Description:      "Log level for outgoing traffic.",
			Optional:         true,
			Default:          dvLogLevelOUT,
			ValidateDiagFunc: validators.FirewallLogLevel(),
		},
		mkMACFilter: {
			Type:        schema.TypeBool,
			Description: "Enable MAC address filtering",
			Optional:    true,
			Default:     dvMACFilter,
		},
		mkNDP: {
			Type:        schema.TypeBool,
			Description: "Enable NDP (Neighbor Discovery Protocol)",
			Optional:    true,
			Default:     dvNDP,
		},
		mkPolicyIn: {
			Type:             schema.TypeString,
			Description:      "Default policy for incoming traffic",
			Optional:         true,
			Default:          dvPolicyIn,
			ValidateDiagFunc: validators.FirewallPolicy(),
		},
		mkPolicyOut: {
			Type:             schema.TypeString,
			Description:      "Default policy for outgoing traffic",
			Optional:         true,
			Default:          dvPolicyOut,
			ValidateDiagFunc: validators.FirewallPolicy(),
		},
		mkRadv: {
			Type:        schema.TypeBool,
			Description: "Allow sending Router Advertisement",
			Optional:    true,
			Default:     dvReadv,
		},
	}

	structure.MergeSchema(s, selectorSchemaMandatory())

	return &schema.Resource{
		Schema:        s,
		CreateContext: selectFirewallAPI(optionsSet),
		ReadContext:   selectFirewallAPI(optionsRead),
		UpdateContext: selectFirewallAPI(optionsUpdate),
		DeleteContext: selectFirewallAPI(optionsDelete),
		Importer: &schema.ResourceImporter{
			StateContext: optionsImport,
		},
	}
}

func optionsImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()

	switch {
	case strings.HasPrefix(id, "vm/"):
		parts := strings.SplitN(id, "/", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid import ID: %s", id)
		}

		nodeName := parts[1]

		vmID, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid import ID: %s", id)
		}

		err = d.Set(mkSelectorNodeName, nodeName)
		if err != nil {
			return nil, fmt.Errorf("failed setting state during import: %w", err)
		}

		err = d.Set(mkSelectorVMID, vmID)
		if err != nil {
			return nil, fmt.Errorf("failed setting state during import: %w", err)
		}
	case strings.HasPrefix(id, "container/"):
		parts := strings.SplitN(id, "/", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid import ID: %s", id)
		}

		nodeName := parts[1]

		containerID, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid import ID: %s", id)
		}

		err = d.Set(mkSelectorNodeName, nodeName)
		if err != nil {
			return nil, fmt.Errorf("failed setting state during import: %w", err)
		}

		err = d.Set(mkSelectorContainerID, containerID)
		if err != nil {
			return nil, fmt.Errorf("failed setting state during import: %w", err)
		}
	default:
		return nil, fmt.Errorf("invalid import ID: %s", id)
	}

	api, err := firewallApiFor(d, m)
	if err != nil {
		return nil, err
	}

	d.SetId(api.GetOptionsID())

	return []*schema.ResourceData{d}, nil
}

func optionsSet(ctx context.Context, api firewall.API, d *schema.ResourceData) diag.Diagnostics {
	dhcp := types.CustomBool(d.Get(mkDHCP).(bool))
	enabled := types.CustomBool(d.Get(mkEnabled).(bool))
	ipFilter := types.CustomBool(d.Get(mkIPFilter).(bool))
	logLevelIn := d.Get(mkLogLevelIN).(string)
	logLevelOut := d.Get(mkLogLevelOUT).(string)
	macFilter := types.CustomBool(d.Get(mkMACFilter).(bool))
	ndp := types.CustomBool(d.Get(mkNDP).(bool))
	policyIn := d.Get(mkPolicyIn).(string)
	policyOut := d.Get(mkPolicyOut).(string)
	radv := types.CustomBool(d.Get(mkRadv).(bool))

	body := &firewall.OptionsPutRequestBody{
		DHCP:        &dhcp,
		Enable:      &enabled,
		IPFilter:    &ipFilter,
		LogLevelIN:  &logLevelIn,
		LogLevelOUT: &logLevelOut,
		MACFilter:   &macFilter,
		NDP:         &ndp,
		PolicyIn:    &policyIn,
		PolicyOut:   &policyOut,
		RAdv:        &radv,
	}

	err := api.SetOptions(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(api.GetOptionsID())

	return optionsRead(ctx, api, d)
}

func optionsRead(ctx context.Context, api firewall.API, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	options, err := api.GetOptions(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if options.DHCP != nil {
		err = d.Set(mkDHCP, *options.DHCP)
		diags = append(diags, diag.FromErr(err)...)
	}

	if options.Enable != nil {
		err = d.Set(mkEnabled, *options.Enable)
		diags = append(diags, diag.FromErr(err)...)
	}

	if options.IPFilter != nil {
		err = d.Set(mkIPFilter, *options.IPFilter)
		diags = append(diags, diag.FromErr(err)...)
	}

	if options.LogLevelIN != nil {
		err = d.Set(mkLogLevelIN, *options.LogLevelIN)
		diags = append(diags, diag.FromErr(err)...)
	}

	if options.LogLevelOUT != nil {
		err = d.Set(mkLogLevelOUT, *options.LogLevelOUT)
		diags = append(diags, diag.FromErr(err)...)
	}

	if options.MACFilter != nil {
		err = d.Set(mkMACFilter, *options.MACFilter)
		diags = append(diags, diag.FromErr(err)...)
	}

	if options.NDP != nil {
		err = d.Set(mkNDP, *options.NDP)
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

	if options.RAdv != nil {
		err = d.Set(mkRadv, *options.RAdv)
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

func optionsUpdate(ctx context.Context, api firewall.API, d *schema.ResourceData) diag.Diagnostics {
	diags := optionsSet(ctx, api, d)
	if diags.HasError() {
		return diags
	}

	return optionsRead(ctx, api, d)
}

func optionsDelete(_ context.Context, _ firewall.API, d *schema.ResourceData) diag.Diagnostics {
	d.SetId("")

	return nil
}
