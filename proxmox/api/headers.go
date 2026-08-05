/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package api

import (
	"errors"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strings"

	"golang.org/x/net/http/httpguts"
)

// validateExtraHeaders rejects reserved names, invalid syntax and names that only differ in case.
// Header values are never included in the error messages, as they typically hold credentials.
func validateExtraHeaders(headers map[string]string) error {
	// headers that cannot be set by the user, grouped by the reason
	reservedHeaders := []string{
		// written by DoRequest and by the authenticators
		"accept", "authorization", "content-length", "content-type", "cookie", "csrfpreventiontoken",
		// managed by net/http: setting these breaks the client. A user-supplied Accept-Encoding
		// disables the transparent gzip handling, and Header.Set("Host") is a no-op (req.Host is
		// the real field), so accepting it would silently do nothing.
		"accept-encoding", "host",
		// hop-by-hop: ignored when serializing HTTP/1 requests, rejected outright by HTTP/2
		"connection", "keep-alive", "proxy-connection", "te", "trailer", "transfer-encoding", "upgrade",
		// ineffective here: for an HTTPS endpoint, forward proxy authentication travels in the
		// CONNECT request (Transport.ProxyConnectHeader), so these would be tunnelled to Proxmox
		"proxy-authenticate", "proxy-authorization",
	}

	seen := make(map[string]string, len(headers))

	// sorted iteration keeps the reported error stable when more than one header is invalid
	for _, name := range slices.Sorted(maps.Keys(headers)) {
		if name == "" {
			return errors.New("an API header name must not be empty")
		}

		lower := strings.ToLower(name)

		if slices.Contains(reservedHeaders, lower) {
			return fmt.Errorf("the API header %q is managed by the provider and must not be set", name)
		}

		if !httpguts.ValidHeaderFieldName(name) {
			return fmt.Errorf("the API header name %q contains invalid characters", name)
		}

		if !httpguts.ValidHeaderFieldValue(headers[name]) {
			return fmt.Errorf("the value of the API header %q contains invalid characters", name)
		}

		// http.Header canonicalizes names, so "X-Foo" and "x-foo" would collide and the winner
		// would depend on map iteration order.
		if dup, ok := seen[lower]; ok {
			return fmt.Errorf("the API headers %q and %q differ only in case", dup, name)
		}

		seen[lower] = name
	}

	return nil
}

// headerTransport adds user-defined headers to requests targeting the API endpoint only.
type headerTransport struct {
	next    http.RoundTripper
	scheme  string
	host    string
	headers map[string]string
}

func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Match scheme and host: the http.Client is shared with callers that fetch arbitrary
	// user-supplied URLs, and it follows an https -> http redirect on the same host, which
	// would put the credentials on the wire in plaintext.
	if req.URL.Scheme != t.scheme || !strings.EqualFold(req.URL.Host, t.host) {
		//nolint:wrapcheck // a RoundTripper must pass the transport error through unchanged
		return t.next.RoundTrip(req)
	}

	// Clone instead of mutating: required by the RoundTripper contract, and it keeps the headers
	// out of the snapshot http.Client takes of the original request to populate redirects.
	req = req.Clone(req.Context())

	for name, value := range t.headers {
		req.Header.Set(name, value)
	}

	//nolint:wrapcheck // a RoundTripper must pass the transport error through unchanged
	return t.next.RoundTrip(req)
}
