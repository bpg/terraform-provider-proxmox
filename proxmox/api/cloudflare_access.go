/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package api

import (
	"fmt"
	"net/http"
	"net/url"
)

// CloudflareAccessConfig holds the Cloudflare Access service-token credentials.
type CloudflareAccessConfig struct {
	ClientID     string
	ClientSecret string
}

// cloudflareAccessTransport wraps an http.RoundTripper to inject Cloudflare Access
// service-token headers into requests scoped to the configured endpoint host.
type cloudflareAccessTransport struct {
	base         http.RoundTripper
	config       CloudflareAccessConfig
	endpointHost string
}

// NewCloudflareAccessTransport returns a RoundTripper that adds
// CF-Access-Client-Id and CF-Access-Client-Secret headers to requests
// whose Host matches the configured endpoint host.
// Port differences are ignored.
func NewCloudflareAccessTransport(
	base http.RoundTripper,
	config CloudflareAccessConfig,
	endpointURL string,
) http.RoundTripper {
	host := extractHost(endpointURL)

	return &cloudflareAccessTransport{
		base:         base,
		config:       config,
		endpointHost: host,
	}
}

// RoundTrip implements http.RoundTripper.
func (t *cloudflareAccessTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}

	if !matchesEndpointHost(req, t.endpointHost) {
		resp, err := base.RoundTrip(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send Cloudflare Access passthrough request: %w", err)
		}

		return resp, nil
	}

	cloned := req.Clone(req.Context())
	cloned.Header = req.Header.Clone()

	cloned.Header.Set("CF-Access-Client-Id", t.config.ClientID)
	cloned.Header.Set("CF-Access-Client-Secret", t.config.ClientSecret)

	resp, err := base.RoundTrip(cloned)
	if err != nil {
		return nil, fmt.Errorf("failed to send Cloudflare Access authenticated request: %w", err)
	}

	return resp, nil
}

// extractHost parses the endpoint URL and returns the hostname, stripping any port.
func extractHost(endpointURL string) string {
	u, err := url.Parse(endpointURL)
	if err != nil {
		return ""
	}

	return u.Hostname()
}

// matchesEndpointHost returns true when the request is addressed to the same
// origin as the Proxmox endpoint. Port differences are ignored so that
// "pve.example.com" matches "pve.example.com:8006".
func matchesEndpointHost(req *http.Request, endpointHost string) bool {
	if endpointHost == "" {
		return false
	}

	return req.URL.Hostname() == endpointHost
}
