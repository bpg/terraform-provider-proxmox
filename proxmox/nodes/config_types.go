/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// ConfigGetResponseBody contains the body from a config get response.
type ConfigGetResponseBody struct {
	Data *[]ConfigGetResponseData `json:"data,omitempty"`
}

// ConfigGetResponseData contains the data from a config get response.
type ConfigGetResponseData struct {
	// Node specific ACME settings.
	ACME *ACMEConfig `json:"acme,omitempty"`
	// ACME domain and validation plugin
	ACMEDomain0 *ACMEDomainConfig `json:"acmedomain0,omitempty"`
	// ACME domain and validation plugin
	ACMEDomain1 *ACMEDomainConfig `json:"acmedomain1,omitempty"`
	// ACME domain and validation plugin
	ACMEDomain2 *ACMEDomainConfig `json:"acmedomain2,omitempty"`
	// ACME domain and validation plugin
	ACMEDomain3 *ACMEDomainConfig `json:"acmedomain3,omitempty"`
	// ACME domain and validation plugin
	ACMEDomain4 *ACMEDomainConfig `json:"acmedomain4,omitempty"`
	// Description for the Node. Shown in the web-interface node notes panel. This is saved as comment inside the configuration file.
	Description *string `json:"description,omitempty"`
	// Prevent changes if current configuration file has different SHA1 digest. This can be used to prevent concurrent modifications.
	Digest *string `json:"digest,omitempty"`
	// Initial delay in seconds, before starting all the Virtual Guests with on-boot enabled.
	StartAllOnbootDelay *int `json:"startall-onboot-delay,omitempty"`
	// Node specific wake on LAN settings.
	WakeOnLan *WakeOnLandConfig `json:"wakeonlan,omitempty"`
}

// ConfigUpdateRequestBody contains the body for a config update request.
type ConfigUpdateRequestBody struct {
	// Node specific ACME settings.
	ACME *ACMEConfig `json:"acme,omitempty"`
	// ACME domain and validation plugin
	ACMEDomain0 *ACMEDomainConfig `json:"acmedomain0,omitempty"`
	// ACME domain and validation plugin
	ACMEDomain1 *ACMEDomainConfig `json:"acmedomain1,omitempty"`
	// ACME domain and validation plugin
	ACMEDomain2 *ACMEDomainConfig `json:"acmedomain2,omitempty"`
	// ACME domain and validation plugin
	ACMEDomain3 *ACMEDomainConfig `json:"acmedomain3,omitempty"`
	// ACME domain and validation plugin
	ACMEDomain4 *ACMEDomainConfig `json:"acmedomain4,omitempty"`
	Delete      *string           `json:"delete,omitempty"`
	// Description for the Node. Shown in the web-interface node notes panel. This is saved as comment inside the configuration file.
	Description *string `json:"description,omitempty"`
	// Prevent changes if current configuration file has different SHA1 digest. This can be used to prevent concurrent modifications.
	Digest *string `json:"digest,omitempty"`
	// Initial delay in seconds, before starting all the Virtual Guests with on-boot enabled.
	StartAllOnbootDelay *int `json:"startall-onboot-delay,omitempty"`
	// Node specific wake on LAN settings.
	WakeOnLan *WakeOnLandConfig `json:"wakeonlan,omitempty"`
}

// ACMEConfig contains the ACME account / domains configuration that use the "standalone" plugin (http challenge).
type ACMEConfig struct {
	// account name
	Account *string
	// domains
	Domains []string
}

// UnmarshalJSON unmarshals a ACMEConfig struct from JSON.
func (a *ACMEConfig) UnmarshalJSON(b []byte) error {
	config := ACMEConfig{}

	s := ""
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("error unmarshaling json: %w", err)
	}

	parts := strings.Split(s, ",")
	for _, part := range parts {
		kv := strings.Split(part, "=")
		if len(kv) != 2 {
			return fmt.Errorf("invalid key-value pair: %s", part)
		}

		switch kv[0] {
		case "account":
			config.Account = &kv[1]
		case "domains":
			config.Domains = strings.Split(kv[1], ";")
		default:
			return fmt.Errorf("unknown key: %s", kv[0])
		}
	}

	*a = config

	return nil
}

// EncodeValues encodes a ACMEConfig struct into a string.
func (a *ACMEConfig) EncodeValues(key string, v *url.Values) error {
	value := ""
	if a.Account != nil {
		value = fmt.Sprintf("account=%s", *a.Account)
	}

	value = fmt.Sprintf("%s,%s", value, strings.Join(a.Domains, ";"))
	v.Add(key, value)

	return nil
}

// ACMEDomainConfig contains the ACME domain configuration for domains using the dns challenge plugin.
type ACMEDomainConfig struct {
	// domain name
	Domain string
	// alias for the domain
	Alias *string
	// name of the plugin configuration
	Plugin *string
}

// UnmarshalJSON unmarshals a ACMEDomainConfig struct from JSON.
func (a *ACMEDomainConfig) UnmarshalJSON(b []byte) error {
	config := ACMEDomainConfig{}

	s := ""
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("error unmarshaling json: %w", err)
	}

	parts := strings.Split(s, ",")
	for _, part := range parts {
		kv := strings.Split(part, "=")
		if len(kv) == 1 {
			config.Domain = kv[0]
		}

		switch kv[0] {
		case "alias":
			config.Alias = &kv[1]
		case "plugin":
			config.Plugin = &kv[1]
		default:
			return fmt.Errorf("unknown key: %s", kv[0])
		}
	}

	*a = config

	return nil
}

// EncodeValues encodes a ACMEDomainConfig struct into a string.
func (a *ACMEDomainConfig) EncodeValues(key string, v *url.Values) error {
	value := a.Domain
	if a.Alias != nil {
		value = fmt.Sprintf("%s,alias=%s", value, *a.Alias)
	}

	if a.Plugin != nil {
		value = fmt.Sprintf("%s,plugin=%s", value, *a.Plugin)
	}

	v.Add(key, value)

	return nil
}

// WakeOnLandConfig contains the wake on LAN configuration.
type WakeOnLandConfig struct {
	// MAC address
	MACAddress string
	// bind interface
	BindInterface *string
	// IPv4 broadcast address
	BroadcastAddress *string
}

// UnmarshalJSON unmarshals a WakeOnLandConfig struct from JSON.
func (a *WakeOnLandConfig) UnmarshalJSON(b []byte) error {
	config := WakeOnLandConfig{}

	s := ""
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("error unmarshaling json: %w", err)
	}

	parts := strings.Split(s, ",")
	for _, part := range parts {
		kv := strings.Split(part, "=")
		if len(kv) == 1 {
			config.MACAddress = kv[0]
		}

		switch kv[0] {
		case "bind-interface":
			config.BindInterface = &kv[1]
		case "broadcast-address":
			config.BroadcastAddress = &kv[1]
		default:
			return fmt.Errorf("unknown key: %s", kv[0])
		}
	}

	*a = config

	return nil
}

// EncodeValues encodes a WakeOnLandConfig struct into a string.
func (a *WakeOnLandConfig) EncodeValues(key string, v *url.Values) error {
	value := a.MACAddress
	if a.BindInterface != nil {
		value = fmt.Sprintf("%s,bind-interface=%s", value, *a.BindInterface)
	}

	if a.BroadcastAddress != nil {
		value = fmt.Sprintf("%s,broadcast-address=%s", value, *a.BroadcastAddress)
	}

	v.Add(key, value)

	return nil
}
