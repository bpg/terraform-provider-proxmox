/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package api

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type countingReader struct {
	reader io.Reader
	read   int
}

func (r *countingReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	r.read += n

	//nolint:wrapcheck // an io.Reader must pass io.EOF through unchanged
	return n, err
}

func TestLoggingTransportKeepsUploadBodyFromLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		contentType     string
		loggerReadsBody bool
	}{
		{"multipart upload", "multipart/form-data; boundary=x", false},
		{"form body", "application/x-www-form-urlencoded", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			payload := strings.Repeat("x", 4096)
			body := &countingReader{reader: strings.NewReader(payload)}

			var (
				readBeforeTransport int
				sentLength          int64
				sent                []byte
			)

			base := roundTripFunc(func(req *http.Request) (*http.Response, error) {
				readBeforeTransport = body.read
				sentLength = req.ContentLength

				var err error

				sent, err = io.ReadAll(req.Body)
				require.NoError(t, err)

				return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Header: http.Header{}, Request: req}, nil
			})

			req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "https://pve.example:8006/api2/json/upload", body)
			require.NoError(t, err)

			req.ContentLength = int64(len(payload))
			req.Header.Set("Content-Type", tt.contentType)

			resp, err := newLoggingTransport(base).RoundTrip(req)
			require.NoError(t, err)
			require.NoError(t, resp.Body.Close())

			require.Equal(t, tt.loggerReadsBody, readBeforeTransport > 0, "bytes read by the logger: %d", readBeforeTransport)
			require.Equal(t, int64(len(payload)), sentLength)
			require.Equal(t, payload, string(sent))
		})
	}
}
