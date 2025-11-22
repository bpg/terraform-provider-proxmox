/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ip

import "net"

// IsValidGlobalUnicast checks if an IP address is a valid global unicast address.
// A global unicast address is one that is:
//   - A valid unicast address (not multicast or broadcast)
//   - Not a loopback address (127.0.0.0/8 for IPv4, ::1 for IPv6)
//   - Not a link-local address (169.254.0.0/16 for IPv4, fe80::/10 for IPv6)
//   - Not an unspecified address (0.0.0.0, ::)
//
// This function accepts:
//   - Public IP addresses (e.g., 8.8.8.8, 2001:4860:4860::8888)
//   - Private IP addresses (e.g., 192.168.1.1, 10.0.0.1, 172.16.0.1)
//   - Unique local IPv6 addresses (e.g., fc00::1, fd00::1)
//
// This function rejects:
//   - Loopback addresses (127.0.0.1, ::1)
//   - Link-local addresses (169.254.x.x, fe80::)
//   - All multicast addresses (224.0.0.1, ff02::1, etc.)
//   - Broadcast addresses
//   - Unspecified addresses (0.0.0.0, ::)
//   - Invalid/malformed IP addresses
func IsValidGlobalUnicast(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	return ip.IsGlobalUnicast()
}
