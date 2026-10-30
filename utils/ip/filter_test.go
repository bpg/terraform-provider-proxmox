/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ip

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddressFilterMatch(t *testing.T) {
	t.Parallel()

	docker := map[string][]string{
		"lo":      {"127.0.0.1", "::1"},
		"docker0": {"172.17.0.1"},
		"eth0":    {"fe80::1%eth0"},
	}

	tests := []struct {
		name       string
		filter     AddressFilter
		ifaces     map[string][]string
		ipv4, ipv6 bool
		want       bool
	}{
		{"default: docker counts (legacy behavior)", AddressFilter{}, docker, false, false, true},
		{"default: only lo and link-local", AddressFilter{}, map[string][]string{"lo": {"127.0.0.1"}, "eth0": {"fe80::1", "169.254.1.1"}}, false, false, false},
		{"cidrs exclude docker", AddressFilter{CIDRs: []string{"10.0.0.0/8"}}, docker, false, false, false},
		{"cidrs match real NIC", AddressFilter{CIDRs: []string{"10.0.0.0/8"}}, map[string][]string{"eth0": {"10.1.2.3"}}, true, false, true},
		{"cidrs disable link-local exclusion", AddressFilter{CIDRs: []string{"fe80::/10"}}, docker, false, true, true},
		{"cidrs disable loopback exclusion", AddressFilter{CIDRs: []string{"127.0.0.0/8"}}, map[string][]string{"eth0": {"127.0.0.2"}}, false, false, true},
		{"interfaces exclude docker", AddressFilter{Interfaces: []string{"eth.*"}}, docker, false, false, false},
		{"interfaces keep link-local exclusion", AddressFilter{Interfaces: []string{"eth0"}}, map[string][]string{"eth0": {"fe80::1"}}, false, false, false},
		{"interfaces full match, not substring", AddressFilter{Interfaces: []string{"eth"}}, map[string][]string{"veth123": {"10.0.0.1"}}, false, false, false},
		{"interfaces can select lo", AddressFilter{Interfaces: []string{"lo"}, CIDRs: []string{"127.0.0.0/8"}}, docker, true, false, true},
		{"AND: right iface, wrong cidr", AddressFilter{Interfaces: []string{"eth0"}, CIDRs: []string{"10.0.0.0/8"}}, map[string][]string{"eth0": {"192.168.1.5"}, "eth1": {"10.0.0.5"}}, false, false, false},
		{"split families across ifaces", AddressFilter{Interfaces: []string{"eth1", "eth2"}, CIDRs: []string{"10.0.0.0/8", "2008::/64"}}, map[string][]string{"eth1": {"10.0.0.5"}, "eth2": {"2008::5"}}, true, true, true},
		{"both families on one iface", AddressFilter{Interfaces: []string{"eth1", "eth2"}, CIDRs: []string{"10.0.0.0/8", "2008::/64"}}, map[string][]string{"eth1": {"10.0.0.5", "2008::5"}}, true, true, true},
		{"IPv4-mapped cidr matches IPv4", AddressFilter{CIDRs: []string{"::ffff:10.0.0.0/104"}}, map[string][]string{"eth0": {"10.1.2.3"}}, true, false, true},
		{"required family missing", AddressFilter{CIDRs: []string{"10.0.0.0/8", "2008::/64"}}, map[string][]string{"eth1": {"10.0.0.5"}}, true, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.filter.Match(tt.ifaces).Satisfied(tt.ipv4, tt.ipv6))
		})
	}
}

func TestAddressFilterRejectedReasons(t *testing.T) {
	t.Parallel()

	res := AddressFilter{Interfaces: []string{"eth0"}}.Match(map[string][]string{"docker0": {"172.17.0.1"}})

	assert.Equal(t, []string{"docker0 172.17.0.1: interface not matched"}, res.Rejected)

	res = AddressFilter{}.Match(map[string][]string{"lo": {"127.0.0.1"}, "eth0": {"fe80::1"}, "eth1": {"169.254.1.1"}})

	assert.Equal(t, []string{
		"eth0 fe80::1: not a global unicast address",
		"eth1 169.254.1.1: not a global unicast address",
		"lo 127.0.0.1: loopback interface",
	}, res.Rejected)
}

func TestRejectedSuffix(t *testing.T) {
	t.Parallel()

	assert.Empty(t, RejectedSuffix(nil))
	assert.Equal(t, " (rejected: a; b)", RejectedSuffix([]string{"a", "b"}))
}

func TestAddressFilterValidateFamilies(t *testing.T) {
	t.Parallel()

	v4 := AddressFilter{CIDRs: []string{"10.0.0.0/8"}}

	require.NoError(t, v4.ValidateFamilies(true, false))
	require.Error(t, v4.ValidateFamilies(false, true))
	require.NoError(t, AddressFilter{}.ValidateFamilies(true, true))
	require.NoError(t, AddressFilter{CIDRs: []string{"unknown", "10.0.0.0/8"}}.ValidateFamilies(false, true))
	require.NoError(t, AddressFilter{CIDRs: []string{"::ffff:10.0.0.0/104"}}.ValidateFamilies(true, false))
}
