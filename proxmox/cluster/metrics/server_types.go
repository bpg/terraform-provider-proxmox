/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package metrics

// ServerData contains the data from a metrics server response.
type ServerData struct {
	Disable *int64  `json:"disable,omitempty" url:"disable,omitempty"`
	ID      string  `json:"id,omitempty"      url:"id,omitempty"`
	MTU     *int64  `json:"mtu"               url:"mtu,omitempty"`
	Port    int64   `json:"port"              url:"port"`
	Server  string  `json:"server"            url:"server"`
	Timeout *int64  `json:"timeout,omitempty" url:"timeout,omitempty"`
	Type    *string `json:"type,omitempty"    url:"type,omitempty"`

	// influxdb only options
	APIPathPrefix *string `json:"api-path-prefix,omitempty"    url:"api-path-prefix,omitempty"`
	Bucket        *string `json:"bucket,omitempty"             url:"bucket,omitempty"`
	InfluxDBProto *string `json:"influxdbproto,omitempty"      url:"influxdbproto,omitempty"`
	MaxBodySize   *int64  `json:"max-body-size,omitempty"      url:"max-body-size,omitempty"`
	Organization  *string `json:"organization,omitempty"       url:"organization,omitempty"`
	Token         *string `json:"token,omitempty"              url:"token,omitempty"`
	Verify        *int64  `json:"verify-certificate,omitempty" url:"verify-certificate,omitempty"`

	// graphite only options
	Path  *string `json:"path,omitempty"  url:"path,omitempty"`
	Proto *string `json:"proto,omitempty" url:"proto,omitempty"`
}

// ServerResponseBody contains the body from a metrics server response.
type ServerResponseBody struct {
	Data *ServerData `json:"data,omitempty"`
}

// ServersResponseBody contains the body from a metrics server list response.
type ServersResponseBody struct {
	Data *[]ServerData `json:"data,omitempty"`
}

// ServerRequestData contains the data for a metric server post/put request.
type ServerRequestData struct {
	ServerData
	Delete *[]string `url:"delete,omitempty"`
}
