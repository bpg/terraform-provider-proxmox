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

	"github.com/bpg/proxmox-api/firewall"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/validator"
)

const (
	mkSelectorNodeName    = "node_name"
	mkSelectorVMID        = "vm_id"
	mkSelectorContainerID = "container_id"
)

func selectorSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		mkSelectorNodeName: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the node.",
		},
		mkSelectorVMID: {
			Type:             schema.TypeInt,
			Optional:         true,
			Description:      "The ID of the VM to manage the firewall for.",
			RequiredWith:     []string{mkSelectorNodeName},
			ValidateDiagFunc: validator.VMID(),
		},
		mkSelectorContainerID: {
			Type:             schema.TypeInt,
			Optional:         true,
			Description:      "The ID of the container to manage the firewall for.",
			RequiredWith:     []string{mkSelectorNodeName},
			ValidateDiagFunc: validator.VMID(),
		},
	}
}

func selectorSchemaMandatory() map[string]*schema.Schema {
	s := selectorSchema()
	s[mkSelectorNodeName].Optional = false
	s[mkSelectorNodeName].Required = true
	// required attributes can't be included in RequiredWith
	s[mkSelectorVMID].RequiredWith = nil
	s[mkSelectorContainerID].RequiredWith = nil

	return s
}

func selectFirewallAPI(
	f func(context.Context, firewall.API, *schema.ResourceData) diag.Diagnostics,
) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		config := m.(proxmoxtf.ProviderConfiguration)

		api, err := config.GetClient()
		if err != nil {
			return diag.FromErr(err)
		}

		var fwAPI firewall.API = api.Cluster().Firewall()

		if nn, ok := d.GetOk(mkSelectorNodeName); ok {
			nodeName := nn.(string)
			nodeAPI := api.Node(nodeName)

			if v, ok := d.GetOk(mkSelectorVMID); ok {
				fwAPI = nodeAPI.VM(v.(int)).Firewall()
			} else if v, ok := d.GetOk(mkSelectorContainerID); ok {
				fwAPI = nodeAPI.Container(v.(int)).Firewall()
			}
		}

		return f(ctx, fwAPI, d)
	}
}
