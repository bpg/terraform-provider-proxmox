/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package api

import (
	"context"
	"io"
	"mime"
	"net/http"
	"net/textproto"
	"regexp"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
)

// logMaskingTransport masks credentials in the request/response logs written by the SDK logging transport it wraps, and
// keeps the body of a file upload away from it.
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

// newLoggingTransport wraps base with the SDK logging transport and the masking around it.
func newLoggingTransport(base http.RoundTripper) http.RoundTripper {
	return newLogMaskingTransport(logging.NewLoggingHTTPTransport(&withheldBodyTransport{next: base}))
}

func (t *logMaskingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := tflog.MaskFieldValuesWithFieldKeys(req.Context(), t.headers...)
	ctx = tflog.MaskAllFieldValuesRegexes(ctx, t.bodyPatterns...)

	if mediaType, _, err := mime.ParseMediaType(req.Header.Get("Content-Type")); err != nil || mediaType != "multipart/form-data" {
		//nolint:wrapcheck // a RoundTripper must pass the transport error through unchanged
		return t.next.RoundTrip(req.WithContext(ctx))
	}

	// The logging transport reads the whole body into memory and logs it, and the body of an upload is the file.
	// It gets the request without the body, withheldBodyTransport puts the body back behind it.
	ctx = context.WithValue(ctx, withheldBodyKey{}, &withheldBody{
		body:          req.Body,
		contentLength: req.ContentLength,
		getBody:       req.GetBody,
	})
	ctx = tflog.SetField(ctx, "tf_http_req_body_withheld_bytes", req.ContentLength)

	req = req.Clone(ctx)
	req.Body = http.NoBody
	req.ContentLength = 0
	req.GetBody = nil

	//nolint:wrapcheck // a RoundTripper must pass the transport error through unchanged
	return t.next.RoundTrip(req)
}

type withheldBodyKey struct{}

type withheldBody struct {
	body          io.ReadCloser
	contentLength int64
	getBody       func() (io.ReadCloser, error)
}

// withheldBodyTransport restores the request body that logMaskingTransport withheld from the logging transport.
type withheldBodyTransport struct {
	next http.RoundTripper
}

func (t *withheldBodyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if withheld, ok := req.Context().Value(withheldBodyKey{}).(*withheldBody); ok {
		req = req.Clone(req.Context())
		req.Body = withheld.body
		req.ContentLength = withheld.contentLength
		req.GetBody = withheld.getBody
	}

	//nolint:wrapcheck // a RoundTripper must pass the transport error through unchanged
	return t.next.RoundTrip(req)
}
