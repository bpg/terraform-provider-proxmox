/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ip

import (
	"errors"
	"fmt"
	"net/netip"
	"regexp"
	"slices"
	"strings"
)

// AddressFilter selects which agent-reported addresses count when waiting for an IP.
// An empty field does not restrict. Setting CIDRs replaces the default loopback/link-local exclusion;
// setting Interfaces replaces the default skip of the "lo" interface.
type AddressFilter struct {
	CIDRs      []string // address must be inside one of these prefixes
	Interfaces []string // interface name must fully match one of these RE2 patterns
}

// IsEmpty reports whether the filter restricts nothing.
func (f AddressFilter) IsEmpty() bool {
	return len(f.CIDRs) == 0 && len(f.Interfaces) == 0
}

// FilterResult is the outcome of matching a set of interfaces against an AddressFilter.
type FilterResult struct {
	HasIPv4  bool
	HasIPv6  bool
	Rejected []string // "<iface> <addr>: <reason>" for every non-qualifying address
}

// Match checks every address of every interface (name -> addresses) against the filter.
func (f AddressFilter) Match(ifaces map[string][]string) FilterResult {
	var prefixes []netip.Prefix

	for _, c := range f.CIDRs {
		if p, err := netip.ParsePrefix(c); err == nil {
			prefixes = append(prefixes, unmapPrefix(p).Masked())
		}
	}

	var patterns []*regexp.Regexp

	for _, s := range f.Interfaces {
		if re, err := regexp.Compile("^(?:" + s + ")$"); err == nil {
			patterns = append(patterns, re)
		}
	}

	var res FilterResult

	for name, addrs := range ifaces {
		for _, a := range addrs {
			addr, err := netip.ParseAddr(a)
			if err != nil {
				res.Rejected = append(res.Rejected, fmt.Sprintf("%s %s: unparsable address", name, a))
				continue
			}

			addr = addr.WithZone("").Unmap()

			if reason := f.reject(name, addr, prefixes, patterns); reason != "" {
				res.Rejected = append(res.Rejected, fmt.Sprintf("%s %s: %s", name, a, reason))
				continue
			}

			if addr.Is4() {
				res.HasIPv4 = true
			} else {
				res.HasIPv6 = true
			}
		}
	}

	slices.Sort(res.Rejected)

	return res
}

// RejectedSuffix formats rejected addresses for a timeout error, or returns "" when there are none.
func RejectedSuffix(rejected []string) string {
	if len(rejected) == 0 {
		return ""
	}

	return " (rejected: " + strings.Join(rejected, "; ") + ")"
}

func (f AddressFilter) reject(name string, addr netip.Addr, prefixes []netip.Prefix, patterns []*regexp.Regexp) string {
	if len(f.Interfaces) == 0 {
		if name == "lo" {
			return "loopback interface"
		}
	} else if !matchesAny(name, patterns) {
		return "interface not matched"
	}

	if len(f.CIDRs) == 0 {
		if !addr.IsGlobalUnicast() {
			return "not a global unicast address"
		}
	} else if !containsAny(addr, prefixes) {
		return "outside cidrs"
	}

	return ""
}

// Satisfied reports whether the result fulfils the required families; with none required, any address does.
func (r FilterResult) Satisfied(ipv4, ipv6 bool) bool {
	if !ipv4 && !ipv6 {
		return r.HasIPv4 || r.HasIPv6
	}

	return (!ipv4 || r.HasIPv4) && (!ipv6 || r.HasIPv6)
}

// ValidateFamilies rejects required families that no configured CIDR can ever satisfy.
// Unparsable (e.g. not yet known) CIDRs make the check inconclusive, so it passes.
func (f AddressFilter) ValidateFamilies(ipv4, ipv6 bool) error {
	if len(f.CIDRs) == 0 {
		return nil
	}

	var has4, has6 bool

	for _, c := range f.CIDRs {
		p, err := netip.ParsePrefix(c)
		if err != nil {
			return nil //nolint:nilerr // inconclusive, not invalid: element validation reports bad CIDRs
		}

		if unmapPrefix(p).Addr().Is4() {
			has4 = true
		} else {
			has6 = true
		}
	}

	if ipv4 && !has4 {
		return errors.New("ipv4 = true requires at least one IPv4 entry in cidrs")
	}

	if ipv6 && !has6 {
		return errors.New("ipv6 = true requires at least one IPv6 entry in cidrs")
	}

	return nil
}

func matchesAny(s string, patterns []*regexp.Regexp) bool {
	for _, re := range patterns {
		if re.MatchString(s) {
			return true
		}
	}

	return false
}

func containsAny(addr netip.Addr, prefixes []netip.Prefix) bool {
	for _, p := range prefixes {
		if p.Contains(addr) {
			return true
		}
	}

	return false
}

// unmapPrefix turns an IPv4-mapped IPv6 prefix (::ffff:a.b.c.d/n, n >= 96) into its IPv4 form, matching
// the unmapped addresses Match compares against.
func unmapPrefix(p netip.Prefix) netip.Prefix {
	if !p.Addr().Is4In6() || p.Bits() < 96 {
		return p
	}

	return netip.PrefixFrom(p.Addr().Unmap(), p.Bits()-96)
}
