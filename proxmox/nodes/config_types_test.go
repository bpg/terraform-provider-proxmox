/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes

import (
	"net/url"
	"testing"
)

func TestACMEConfig_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  ACMEConfig
		str     string
		wantErr bool
	}{
		{
			name: "account only",
			config: ACMEConfig{
				Account: new("foo"),
				Domains: nil,
			},
			str: `"account=foo"`,
		},
		{
			name: "account and domains",
			config: ACMEConfig{
				Account: new("foo"),
				Domains: []string{"bar", "baz"},
			},
			str: `"account=foo,domains=bar;baz"`,
		},
		{
			name: "domains only",
			config: ACMEConfig{
				Account: nil,
				Domains: []string{"bar", "baz"},
			},
			str: `"domains=bar;baz"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := tt.config.UnmarshalJSON([]byte(tt.str)); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestACMEConfig_EncodeValues(t *testing.T) {
	t.Parallel()

	type args struct {
		key string
		v   *url.Values
	}

	tests := []struct {
		name    string
		config  ACMEConfig
		args    args
		wantErr bool
	}{
		{
			name: "account only",
			config: ACMEConfig{
				Account: new("foo"),
				Domains: nil,
			},
			args: args{
				"acme",
				&url.Values{
					"account": {"foo"},
				},
			},
		},
		{
			name: "account and domains",
			config: ACMEConfig{
				Account: new("foo"),
				Domains: []string{"bar", "baz"},
			},
			args: args{
				"acme",
				&url.Values{
					"account": {"foo"},
					"domains": {"bar;baz"},
				},
			},
		},
		{
			name: "domains only",
			config: ACMEConfig{
				Account: nil,
				Domains: []string{"bar", "baz"},
			},
			args: args{
				"acme",
				&url.Values{
					"domains": {"bar;baz"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := tt.config.EncodeValues(tt.args.key, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("EncodeValues() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestACMEDomainConfig_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  ACMEDomainConfig
		str     string
		wantErr bool
	}{
		{
			name: "domain only",
			config: ACMEDomainConfig{
				Domain: "foo",
			},
			str: `"foo"`,
		},
		{
			name: "domain only with key",
			config: ACMEDomainConfig{
				Domain: "foo",
			},
			str: `"domain=foo"`,
		},
		{
			name: "domain and alias",
			config: ACMEDomainConfig{
				Domain: "foo",
				Alias:  new("bar"),
			},
			str: `"domain=foo,alias=bar"`,
		},
		{
			name: "domain and plugin",
			config: ACMEDomainConfig{
				Domain: "foo",
				Plugin: new("bar"),
			},
			str: `"domain=foo,plugin=bar"`,
		},
		{
			name: "domain, alias, and plugin",
			config: ACMEDomainConfig{
				Domain: "foo",
				Alias:  new("bar"),
				Plugin: new("baz"),
			},
			str: `"domain=foo,alias=bar,plugin=baz"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := tt.config.UnmarshalJSON([]byte(tt.str)); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestACMEDomainConfig_EncodeValues(t *testing.T) {
	t.Parallel()

	type args struct {
		key string
		v   *url.Values
	}

	tests := []struct {
		name    string
		config  ACMEDomainConfig
		args    args
		wantErr bool
	}{
		{
			name: "domain only",
			config: ACMEDomainConfig{
				Domain: "foo",
			},
			args: args{
				"acme",
				&url.Values{
					"domain": {"foo"},
				},
			},
		},
		{
			name: "domain and alias",
			config: ACMEDomainConfig{
				Domain: "foo",
				Alias:  new("bar"),
			},
			args: args{
				"acme",
				&url.Values{
					"domain": {"foo"},
					"alias":  {"bar"},
				},
			},
		},
		{
			name: "domain and plugin",
			config: ACMEDomainConfig{
				Domain: "foo",
				Plugin: new("bar"),
			},
			args: args{
				"acme",
				&url.Values{
					"domain": {"foo"},
					"plugin": {"bar"},
				},
			},
		},
		{
			name: "domain, alias, and plugin",
			config: ACMEDomainConfig{
				Domain: "foo",
				Alias:  new("bar"),
				Plugin: new("baz"),
			},
			args: args{
				"acme",
				&url.Values{
					"domain": {"foo"},
					"alias":  {"bar"},
					"plugin": {"baz"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := tt.config.EncodeValues(tt.args.key, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("EncodeValues() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWakeOnLandConfig_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  WakeOnLandConfig
		str     string
		wantErr bool
	}{
		{
			name: "mac only",
			config: WakeOnLandConfig{
				MACAddress: "00:11:22:33:44:55",
			},
			str: `"00:11:22:33:44:55"`,
		},
		{
			name: "mac only with key",
			config: WakeOnLandConfig{
				MACAddress: "00:11:22:33:44:55",
			},
			str: `"mac=00:11:22:33:44:55"`,
		},
		{
			name: "mac and bind interface",
			config: WakeOnLandConfig{
				MACAddress:    "00:11:22:33:44:55",
				BindInterface: new("eth0"),
			},
			str: `"mac=00:11:22:33:44:55,bind-interface=eth0"`,
		},
		{
			name: "mac and broadcast address",
			config: WakeOnLandConfig{
				MACAddress:       "00:11:22:33:44:55",
				BroadcastAddress: new("192.168.0.155"),
			},
			str: `"mac=00:11:22:33:44:55,broadcast-address=192.168.0.255"`,
		},
		{
			name: "mac, bind interface, and broadcast address",
			config: WakeOnLandConfig{
				MACAddress:       "00:11:22:33:44:55",
				BindInterface:    new("eth0"),
				BroadcastAddress: new("192.168.0.255"),
			},
			str: `"mac=00:11:22:33:44:55,bind-interface=eth0,broadcast-address=192.168.0.255"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := tt.config.UnmarshalJSON([]byte(tt.str)); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWakeOnLandConfig_EncodeValues(t *testing.T) {
	t.Parallel()

	type args struct {
		key string
		v   *url.Values
	}

	tests := []struct {
		name    string
		config  WakeOnLandConfig
		args    args
		wantErr bool
	}{
		{
			name: "mac only",
			config: WakeOnLandConfig{
				MACAddress: "00:11:22:33:44:55",
			},
			args: args{
				"wakeonlan",
				&url.Values{
					"mac": {"00:11:22:33:44:55"},
				},
			},
		},
		{
			name: "mac and bind interface",
			config: WakeOnLandConfig{
				MACAddress:    "00:11:22:33:44:55",
				BindInterface: new("eth0"),
			},
			args: args{
				"wakeonlan",
				&url.Values{
					"mac":            {"00:11:22:33:44:55"},
					"bind-interface": {"eth0"},
				},
			},
		},
		{
			name: "mac and broadcast address",
			config: WakeOnLandConfig{
				MACAddress:       "00:11:22:33:44:55",
				BroadcastAddress: new("192.168.0.255"),
			},
			args: args{
				"wakeonlan",
				&url.Values{
					"mac":               {"00:11:22:33:44:55"},
					"broadcast-address": {"192.168.0.255"},
				},
			},
		},
		{
			name: "mac, bind interface, and broadcast address",
			config: WakeOnLandConfig{
				MACAddress:       "00:11:22:33:44:55",
				BindInterface:    new("eth0"),
				BroadcastAddress: new("10.255.255.255"),
			},
			args: args{
				"wakeonlan",
				&url.Values{
					"mac":               {"00:11:22:33:44:55"},
					"bind-interface":    {"eth0"},
					"broadcast-address": {"10.255.255.255"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := tt.config.EncodeValues(tt.args.key, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("EncodeValues() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
