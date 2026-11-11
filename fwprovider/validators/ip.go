/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package validators

import (
	"fmt"
	"net/netip"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// IPAddress validates that a string is a plain IPv4 or IPv6 address (no CIDR suffix, no zone), as accepted by PVE's
// "ip" format.
func IPAddress() validator.String {
	return NewParseValidator(
		func(s string) (netip.Addr, error) {
			addr, err := netip.ParseAddr(s)
			if err != nil || addr.Zone() != "" {
				return addr, fmt.Errorf("%q is not a valid IPv4 or IPv6 address", s)
			}

			return addr, nil
		},
		"must be a valid IPv4 or IPv6 address",
	)
}
