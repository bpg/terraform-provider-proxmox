/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package api

import (
	"net/http"
	"net/textproto"
	"regexp"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// logMaskingTransport masks credentials in the request/response logs written by the SDK logging transport it wraps.
// The masks must be on the request context, because that is the context the logging transport logs with.
type logMaskingTransport struct {
	next         http.RoundTripper
	headers      []string
	bodyPatterns []*regexp.Regexp
}

func newLogMaskingTransport(next http.RoundTripper) *logMaskingTransport {
	return &logMaskingTransport{
		next: next,
		// the SDK logging transport logs each header as a field keyed by its canonical MIME name
		headers: []string{
			textproto.CanonicalMIMEHeaderKey("Authorization"),
			textproto.CanonicalMIMEHeaderKey("Cookie"),
			textproto.CanonicalMIMEHeaderKey("CSRFPreventionToken"),
		},
		bodyPatterns: []*regexp.Regexp{
			// password login form body (POST /access/ticket)
			regexp.MustCompile(`password=[^&\s]*`),
			// ticket and CSRF token in the login response
			regexp.MustCompile(`"(ticket|CSRFPreventionToken)"\s*:\s*"[^"]*"`),
		},
	}
}

func (t *logMaskingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := tflog.MaskFieldValuesWithFieldKeys(req.Context(), t.headers...)
	ctx = tflog.MaskAllFieldValuesRegexes(ctx, t.bodyPatterns...)

	//nolint:wrapcheck // a RoundTripper must pass the transport error through unchanged
	return t.next.RoundTrip(req.WithContext(ctx))
}
