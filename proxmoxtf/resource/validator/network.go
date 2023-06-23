/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package validator

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// VLANIDsValidator returns a schema validation function for VLAN IDs.
func VLANIDsValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(i interface{}, k string) ([]string, []error) {
		min := 1
		max := 4094

		var ws []string
		var es []error

		list, ok := i.([]interface{})

		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be []interface{}", k))
			return ws, es
		}

		for li, lv := range list {
			v, ok := lv.(int)

			if !ok {
				es = append(es, fmt.Errorf("expected type of %s[%d] to be int", k, li))
				return ws, es
			}

			if v != -1 {
				if v < min || v > max {
					es = append(es, fmt.Errorf("expected %s[%d] to be in the range (%d - %d), got %d", k, li, min, max, v))
					return ws, es
				}
			}
		}

		return ws, es
	})
}

// NodeNetworkInterfaceBondingModes returns a schema validation function for a node network interface bonding mode.
func NodeNetworkInterfaceBondingModes() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"balance-rr",
		"active-backup",
		"balance-xor",
		"broadcast",
		"802.3ad",
		"balance-tlb",
		"balance-alb",
		"balance-slb",
		"lacp-balance-slb",
		"lacp-balance-tcp",
	}, false))
}

// NodeNetworkInterfaceBondingTransmitHashPolicies returns a schema validation function for a node network interface
// bonding transmit hash policy.
func NodeNetworkInterfaceBondingTransmitHashPolicies() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"layer2",
		"layer2+3",
		"layer3+4",
	}, false))
}
