/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
	resource "github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/vm"
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
			ForceNew:         true,
			RequiredWith:     []string{mkSelectorNodeName},
			ValidateDiagFunc: resource.VMIDValidator(),
		},
		mkSelectorContainerID: {
			Type:             schema.TypeInt,
			Optional:         true,
			Description:      "The ID of the container to manage the firewall for.",
			ForceNew:         true,
			RequiredWith:     []string{mkSelectorNodeName},
			ValidateDiagFunc: resource.VMIDValidator(),
		},
	}
}

// selectorCustomizeDiff replaces node and cluster level resources when node_name changes, as an in-place update would
// only reach the new target and leave the objects on the old one. vm_id and container_id are ForceNew in the schema
// instead, which also covers values unknown at plan time. A guest's node_name change is not a new target: guest firewall
// config is stored cluster-wide per VMID, so it survives a migration.
func selectorCustomizeDiff() schema.CustomizeDiffFunc {
	return customdiff.ForceNewIf(mkSelectorNodeName, func(_ context.Context, d *schema.ResourceDiff, _ any) bool {
		if !d.HasChange(mkSelectorNodeName) {
			return false
		}

		for _, key := range []string{mkSelectorVMID, mkSelectorContainerID} {
			o, n := d.GetChange(key)
			if o.(int) != 0 || n.(int) != 0 || !d.NewValueKnown(key) {
				return false
			}
		}

		return true
	})
}

func optionsSelectorSchema() map[string]*schema.Schema {
	s := selectorSchema()
	s[mkSelectorNodeName].Optional = false
	s[mkSelectorNodeName].Required = true
	// options are a set of values on the target, an in-place update on a new target is fine
	s[mkSelectorVMID].ForceNew = false
	s[mkSelectorContainerID].ForceNew = false
	// required attributes can't be included in RequiredWith
	s[mkSelectorVMID].RequiredWith = nil
	s[mkSelectorContainerID].RequiredWith = nil
	// make sure exactly one of vm_id/container_id is present. neither is invalid. both is invalid.
	s[mkSelectorVMID].ExactlyOneOf = []string{mkSelectorVMID, mkSelectorContainerID}
	s[mkSelectorContainerID].ExactlyOneOf = []string{mkSelectorVMID, mkSelectorContainerID}

	return s
}

func firewallAPIFor(d *schema.ResourceData, m any) (firewall.API, error) {
	config := m.(proxmoxtf.ProviderConfiguration)

	api, err := config.GetClient()
	if err != nil {
		return nil, err
	}

	var fwAPI firewall.API = api.Cluster().Firewall()

	if nn, ok := d.GetOk(mkSelectorNodeName); ok {
		nodeName := nn.(string)
		nodeAPI := api.Node(nodeName)
		fwAPI = nodeAPI.Firewall()

		if v, ok := d.GetOk(mkSelectorVMID); ok {
			fwAPI = nodeAPI.VM(v.(int)).Firewall()
		} else if v, ok := d.GetOk(mkSelectorContainerID); ok {
			fwAPI = nodeAPI.Container(v.(int)).Firewall()
		}
	}

	return fwAPI, nil
}

func selectFirewallAPI(
	f func(context.Context, firewall.API, *schema.ResourceData) diag.Diagnostics,
) func(context.Context, *schema.ResourceData, any) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
		fwAPI, err := firewallAPIFor(d, m)
		if err != nil {
			return diag.FromErr(err)
		}

		return f(ctx, fwAPI, d)
	}
}
