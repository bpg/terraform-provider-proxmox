/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/bpg/terraform-provider-proxmox/utils"
	"github.com/bpg/terraform-provider-proxmox/utils/ip"
)

const (
	mkWaitForIPCIDRs      = "cidrs"
	mkWaitForIPInterfaces = "interfaces"
)

// AddWaitForIPFilterSchema adds the `cidrs` and `interfaces` address filter attributes to a wait_for_ip block schema.
func AddWaitForIPFilterSchema(s map[string]*schema.Schema) map[string]*schema.Schema {
	s[mkWaitForIPCIDRs] = &schema.Schema{
		Type: schema.TypeSet,
		Description: "Only addresses inside one of these CIDRs count when waiting. Setting this disables the " +
			"default loopback/link-local address exclusion",
		Optional: true,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
		},
	}
	s[mkWaitForIPInterfaces] = &schema.Schema{
		Type: schema.TypeSet,
		Description: "Only addresses on interfaces whose name fully matches one of these regular expressions count " +
			"when waiting. Setting this disables the default skip of the `lo` interface",
		Optional: true,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsValidRegExp),
		},
	}

	return s
}

// WaitForIPFilterFromBlock reads the address filter from a wait_for_ip block.
func WaitForIPFilterFromBlock(block map[string]any) ip.AddressFilter {
	var filter ip.AddressFilter

	if v := waitForIPFilterValues(block[mkWaitForIPCIDRs]); len(v) > 0 {
		filter.CIDRs = utils.ConvertToStringSlice(v)
	}

	if v := waitForIPFilterValues(block[mkWaitForIPInterfaces]); len(v) > 0 {
		filter.Interfaces = utils.ConvertToStringSlice(v)
	}

	return filter
}

func waitForIPFilterValues(v any) []any {
	switch v := v.(type) {
	case *schema.Set:
		return v.List()
	case []any:
		return v
	default:
		return nil
	}
}

// WaitForIPFilterDiff rejects address filters combined with `disabled = true` and required
// IP families that no configured CIDR can satisfy. blockPath points at the wait_for_ip list.
func WaitForIPFilterDiff(blockPath string) schema.CustomizeDiffFunc {
	return func(_ context.Context, d *schema.ResourceDiff, _ any) error {
		list, _ := d.Get(blockPath).([]any)
		if len(list) == 0 || list[0] == nil {
			return nil
		}

		block := list[0].(map[string]any)
		filter := WaitForIPFilterFromBlock(block)

		if len(filter.CIDRs) == 0 && len(filter.Interfaces) == 0 {
			return nil
		}

		if disabled, _ := block["disabled"].(bool); disabled {
			return fmt.Errorf("%s: cidrs and interfaces cannot be combined with disabled = true", blockPath)
		}

		ipv4, _ := block["ipv4"].(bool)
		ipv6, _ := block["ipv6"].(bool)

		if err := filter.ValidateFamilies(ipv4, ipv6); err != nil {
			return fmt.Errorf("%s: %w", blockPath, err)
		}

		return nil
	}
}
