/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package plugins

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// BaseACMEPluginData contains common fields for ACME plugin data.
type BaseACMEPluginData struct {
	// ACME challenge type (dns, standalone).
	Type string `json:"type,omitempty" url:"type,omitempty"`
	// Prevent changes if current configuration file has a different digest. This can be used to prevent concurrent modifications.
	Digest string `json:"digest,omitempty" url:"digest,omitempty"`
	// API plugin name
	API string `json:"api,omitempty" url:"api,omitempty"`
	// Extra delay in seconds to wait before requesting validation. Allows to cope with a long TTL of DNS records (0 - 172800).
	ValidationDelay int64 `json:"validation-delay,omitempty" url:"validation-delay,omitempty"`
}

// ACMEPluginsListResponseBody contains the body from an ACME plugins list response.
type ACMEPluginsListResponseBody struct {
	// Unique identifier for ACME plugin instance.
	Data []*ACMEPluginsListResponseData `json:"data,omitempty"`
}

// ACMEPluginsListResponseData contains the data from an ACME plugins list response.
type ACMEPluginsListResponseData struct {
	BaseACMEPluginData
	// ACME Plugin ID name
	Plugin string `json:"plugin" url:"plugin"`
	// DNS plugin data.
	Data *DNSPluginData `json:"data,omitempty"`
}

// ACMEPluginsGetResponseBody contains the body from an ACME plugins get response.
type ACMEPluginsGetResponseBody struct {
	Data *ACMEPluginsGetResponseData `json:"data,omitempty"`
}

// ACMEPluginsGetResponseData contains the data from an ACME plugins get response.
type ACMEPluginsGetResponseData struct {
	BaseACMEPluginData
	// ACME Plugin ID name
	Plugin string `json:"plugin" url:"plugin"`
	// DNS plugin data.
	Data *DNSPluginData `json:"data"`
}

// ACMEPluginsCreateRequestBody contains the body for creating a new ACME plugin.
type ACMEPluginsCreateRequestBody struct {
	BaseACMEPluginData
	// ACME Plugin ID name
	Plugin string `json:"id" url:"id"`
	// DNS plugin data. (base64 encoded)
	Data *DNSPluginData `url:"data,omitempty"`
	// Flag to disable the config.
	Disable bool `url:"disable,omitempty"`
	// List of cluster node names.
	Nodes string `url:"nodes,omitempty"`
}

// ACMEPluginsUpdateRequestBody contains the body for updating an existing ACME plugin.
type ACMEPluginsUpdateRequestBody struct {
	BaseACMEPluginData
	// DNS plugin data. (base64 encoded)
	Data *DNSPluginData `url:"data,omitempty"`
	// A list of settings you want to delete.
	Delete string `url:"delete,omitempty"`
	// Flag to disable the config.
	Disable bool `url:"disable,omitempty"`
	// List of cluster node names.
	Nodes string `url:"nodes,omitempty"`
}

// DNSPluginData is a map of DNS plugin data.
type DNSPluginData map[string]string

// EncodeValues encodes the DNSPluginData into the URL values.
func (d DNSPluginData) EncodeValues(key string, v *url.Values) error {
	values := make([]string, 0, len(d))

	for key, value := range d {
		values = append(values, fmt.Sprintf("%s=%s", key, value))
	}

	v.Add(key, base64.StdEncoding.EncodeToString([]byte(strings.Join(values, `\n`))))

	return nil
}

// UnmarshalJSON unmarshals a DNSPluginData struct from JSON.
func (d *DNSPluginData) UnmarshalJSON(b []byte) error {
	mapData := make(map[string]string)

	s := ""
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("error unmarshaling json: %w", err)
	}

	for _, line := range strings.Split(s, `\n`) {
		before, after, found := strings.Cut(line, "=")
		if !found {
			return fmt.Errorf("invalid DNS plugin data: %s", line)
		}

		mapData[before] = after
	}

	*d = mapData

	return nil
}
