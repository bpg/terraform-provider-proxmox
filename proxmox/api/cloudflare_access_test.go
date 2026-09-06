/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package api

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"
	"sync"
	"testing"
)

// recordingTransport captures the last request that was sent through it.
type recordingTransport struct {
	mu       sync.Mutex
	lastReq  *http.Request
	respCode int
	respBody string
}

func (t *recordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.lastReq = req.Clone(req.Context())

	return &http.Response{
		StatusCode: t.respCode,
		Body:       io.NopCloser(strings.NewReader(t.respBody)),
		Header:     make(http.Header),
	}, nil
}

func (t *recordingTransport) lastRequest() *http.Request {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.lastReq
}

func TestCloudflareAccessTransport_HeadersAdded(t *testing.T) {
	t.Parallel()

	rec := &recordingTransport{respCode: http.StatusOK, respBody: `{}`}

	config := CloudflareAccessConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}

	transport := NewCloudflareAccessTransport(rec, config, "https://pve.example.com:8006/")

	client := &http.Client{Transport: transport}

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"https://pve.example.com:8006/api2/json/version",
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	lastReq := rec.lastRequest()
	if lastReq == nil {
		t.Fatal("no request was recorded")
	}

	if got := lastReq.Header.Get("CF-Access-Client-Id"); got != "test-client-id" {
		t.Errorf("CF-Access-Client-Id = %q, want %q", got, "test-client-id")
	}

	if got := lastReq.Header.Get("CF-Access-Client-Secret"); got != "test-client-secret" {
		t.Errorf("CF-Access-Client-Secret = %q, want %q", got, "test-client-secret")
	}
}

func TestCloudflareAccessTransport_ExistingHeadersPreserved(t *testing.T) {
	t.Parallel()

	rec := &recordingTransport{respCode: http.StatusOK, respBody: `{}`}

	config := CloudflareAccessConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}

	transport := NewCloudflareAccessTransport(rec, config, "https://pve.example.com/")

	client := &http.Client{Transport: transport}

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"https://pve.example.com/api2/json/version",
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set("X-Custom-Header", "custom-value")
	req.Header.Set("Authorization", "PVE token=value")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	lastReq := rec.lastRequest()
	if lastReq == nil {
		t.Fatal("no request was recorded")
	}

	if got := lastReq.Header.Get("X-Custom-Header"); got != "custom-value" {
		t.Errorf("X-Custom-Header = %q, want %q", got, "custom-value")
	}

	if got := lastReq.Header.Get("Authorization"); got != "PVE token=value" {
		t.Errorf("Authorization = %q, want %q", got, "PVE token=value")
	}
}

func TestCloudflareAccessTransport_OriginalRequestUnchanged(t *testing.T) {
	t.Parallel()

	rec := &recordingTransport{respCode: http.StatusOK, respBody: `{}`}

	config := CloudflareAccessConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}

	transport := NewCloudflareAccessTransport(rec, config, "https://pve.example.com/")

	client := &http.Client{Transport: transport}

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"https://pve.example.com/api2/json/version",
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	originalHeaderCount := len(req.Header)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if got := req.Header.Get("CF-Access-Client-Id"); got != "" {
		t.Errorf("original request was mutated: CF-Access-Client-Id = %q", got)
	}

	if len(req.Header) != originalHeaderCount {
		t.Errorf("original request header count changed from %d to %d", originalHeaderCount, len(req.Header))
	}
}

func TestCloudflareAccessTransport_DuplicateHeadersOverwritten(t *testing.T) {
	t.Parallel()

	rec := &recordingTransport{respCode: http.StatusOK, respBody: `{}`}

	config := CloudflareAccessConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}

	transport := NewCloudflareAccessTransport(rec, config, "https://pve.example.com/")

	client := &http.Client{Transport: transport}

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"https://pve.example.com/api2/json/version",
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set("CF-Access-Client-Id", "existing-value")
	req.Header.Set("CF-Access-Client-Secret", "existing-secret")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	lastReq := rec.lastRequest()
	if lastReq == nil {
		t.Fatal("no request was recorded")
	}

	if got := lastReq.Header.Get("CF-Access-Client-Id"); got != "test-client-id" {
		t.Errorf("CF-Access-Client-Id = %q, want %q (should overwrite)", got, "test-client-id")
	}

	if got := lastReq.Header[textproto.CanonicalMIMEHeaderKey("CF-Access-Client-Id")]; len(got) != 1 {
		t.Errorf("CF-Access-Client-Id has %d values, want 1", len(got))
	}
}

func TestCloudflareAccessTransport_DifferentHost_NoHeaders(t *testing.T) {
	t.Parallel()

	rec := &recordingTransport{respCode: http.StatusOK, respBody: `{}`}

	config := CloudflareAccessConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}

	transport := NewCloudflareAccessTransport(rec, config, "https://pve.example.com/")

	client := &http.Client{Transport: transport}

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"https://external-host.example.com/api2/json/version",
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	lastReq := rec.lastRequest()
	if lastReq == nil {
		t.Fatal("no request was recorded")
	}

	if got := lastReq.Header.Get("CF-Access-Client-Id"); got != "" {
		t.Errorf("CF-Access-Client-Id should not be set for different host, got %q", got)
	}

	if got := lastReq.Header.Get("CF-Access-Client-Secret"); got != "" {
		t.Errorf("CF-Access-Client-Secret should not be set for different host, got %q", got)
	}
}

func TestCloudflareAccessTransport_NilBase(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cid := r.Header.Get("CF-Access-Client-Id")
		csec := r.Header.Get("CF-Access-Client-Secret")

		if cid != "test-client-id" {
			t.Errorf("CF-Access-Client-Id = %q, want %q", cid, "test-client-id")
		}

		if csec != "test-client-secret" {
			t.Errorf("CF-Access-Client-Secret = %q, want %q", csec, "test-client-secret")
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := CloudflareAccessConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}

	transport := NewCloudflareAccessTransport(nil, config, server.URL)

	client := &http.Client{Transport: transport}

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		server.URL+"/api2/json/version",
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
}

func TestCloudflareAccessTransport_TransportError_NoCredentialsInError(t *testing.T) {
	t.Parallel()

	errTransport := &errorTransport{err: &net.OpError{
		Op:  "dial",
		Net: "tcp",
	}}

	config := CloudflareAccessConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}

	transport := NewCloudflareAccessTransport(errTransport, config, "https://pve.example.com/")

	client := &http.Client{Transport: transport}

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"https://pve.example.com/api2/json/version",
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	errStr := err.Error()
	if strings.Contains(errStr, "test-client-id") || strings.Contains(errStr, "test-client-secret") {
		t.Errorf("error contains credentials: %s", errStr)
	}
}

func TestCloudflareAccessTransport_Redirect_CrossHost(t *testing.T) {
	t.Parallel()

	var (
		mu            sync.Mutex
		sourceHeaders http.Header
		targetHeaders http.Header
	)

	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		targetHeaders = r.Header.Clone()
		mu.Unlock()

		w.WriteHeader(http.StatusOK)
	}))
	defer target.Close()

	// httptest binds to 127.0.0.1; redirecting to "localhost" gives a different hostname on the same machine.
	targetURL := strings.Replace(target.URL, "127.0.0.1", "localhost", 1) + "/target"

	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		sourceHeaders = r.Header.Clone()
		mu.Unlock()

		http.Redirect(w, r, targetURL, http.StatusFound)
	}))
	defer source.Close()

	config := CloudflareAccessConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}

	transport := NewCloudflareAccessTransport(nil, config, source.URL)

	client := &http.Client{Transport: transport}

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		source.URL+"/source",
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected redirect to be followed, got status %d", resp.StatusCode)
	}

	mu.Lock()
	defer mu.Unlock()

	if sourceHeaders == nil {
		t.Fatal("endpoint host was never requested")
	}

	if got := sourceHeaders.Get("CF-Access-Client-Id"); got != "test-client-id" {
		t.Errorf("endpoint host request CF-Access-Client-Id = %q, want %q", got, "test-client-id")
	}

	if targetHeaders == nil {
		t.Fatal("redirect target was never requested")
	}

	if got := targetHeaders.Get("CF-Access-Client-Id"); got != "" {
		t.Errorf("redirect target should not receive CF-Access-Client-Id, got %q", got)
	}

	if got := targetHeaders.Get("CF-Access-Client-Secret"); got != "" {
		t.Errorf("redirect target should not receive CF-Access-Client-Secret, got %q", got)
	}
}

func TestExtractHost(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  string
	}{
		{"https://pve.example.com:8006/", "pve.example.com"},
		{"https://pve.example.com/", "pve.example.com"},
		{"https://10.0.111.15:8006/api2/json", "10.0.111.15"},
		{"https://10.0.111.15/", "10.0.111.15"},
		{"invalid", ""},
	}

	for _, tt := range tests {
		got := extractHost(tt.input)
		if got != tt.want {
			t.Errorf("extractHost(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestMatchesEndpointHost(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		requestURL   string
		endpointHost string
		want         bool
	}{
		{"exact match", "https://pve.example.com:8006/api", "pve.example.com", true},
		{"host only", "https://pve.example.com/api", "pve.example.com", true},
		{"different host", "https://other.example.com/api", "pve.example.com", false},
		{"empty host", "https://pve.example.com/api", "", false},
		{"ip match", "https://10.0.111.15:8006/api", "10.0.111.15", true},
		{"ip different port", "https://10.0.111.15:9999/api", "10.0.111.15", true},
		{"same host different port", "https://pve.example.com:9999/api", "pve.example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequestWithContext(
				context.Background(),
				http.MethodGet,
				tt.requestURL,
				nil,
			)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			got := matchesEndpointHost(req, tt.endpointHost)
			if got != tt.want {
				t.Errorf("matchesEndpointHost(%v, %q) = %v, want %v",
					req.URL, tt.endpointHost, got, tt.want)
			}
		})
	}
}

type errorTransport struct {
	err error
}

func (t *errorTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, t.err
}

func TestNewConnection_CloudflareAccessHeadersSent(t *testing.T) {
	t.Parallel()

	var (
		mu  sync.Mutex
		got http.Header
	)

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		got = r.Header.Clone()
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")

		if _, err := w.Write([]byte(`{"data":{"version":"8.4.1","release":"8.4","repoid":"abc"}}`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	config := &CloudflareAccessConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}

	conn, err := NewConnection(server.URL, true, "", config)
	if err != nil {
		t.Fatalf("NewConnection: %v", err)
	}

	creds, err := NewCredentials("", "", "", "user@pve!token=test", "", "")
	if err != nil {
		t.Fatalf("NewCredentials: %v", err)
	}

	client, err := NewClient(creds, conn)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	var resp map[string]any

	if err := client.DoRequest(context.Background(), http.MethodGet, "version", nil, &resp); err != nil {
		t.Fatalf("DoRequest: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if got == nil {
		t.Fatal("server was never called")
	}

	if v := got.Get("CF-Access-Client-Id"); v != "test-client-id" {
		t.Errorf("CF-Access-Client-Id = %q, want %q", v, "test-client-id")
	}

	if v := got.Get("CF-Access-Client-Secret"); v != "test-client-secret" {
		t.Errorf("CF-Access-Client-Secret = %q, want %q", v, "test-client-secret")
	}

	if v := got.Get("Authorization"); v == "" {
		t.Error("Authorization header missing")
	}
}
