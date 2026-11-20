/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/utils"
)

// NewEmptyClusterListProxy fronts the real PVE endpoint and answers every /cluster/resources request with
// an empty list, which is what a multi-node cluster returns until the guest list has synced to the node
// serving the request.
func NewEmptyClusterListProxy(t *testing.T) string {
	t.Helper()

	target, err := url.Parse(utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT"))
	require.NoError(t, err)

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
			// keep the upstream response uncompressed so the body below can be swapped in
			r.Out.Header.Del("Accept-Encoding")
		},
		Transport: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // in-test proxy to a lab PVE with a self-signed certificate
		},
		ModifyResponse: func(resp *http.Response) error {
			if !strings.HasSuffix(resp.Request.URL.Path, "/cluster/resources") {
				return nil
			}

			_ = resp.Body.Close()

			body := `{"data":[]}`
			resp.Body = io.NopCloser(strings.NewReader(body))
			resp.ContentLength = int64(len(body))
			resp.Header.Set("Content-Length", strconv.Itoa(len(body)))

			return nil
		},
	}

	server := httptest.NewTLSServer(proxy)
	t.Cleanup(server.Close)

	return server.URL
}
