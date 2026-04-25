/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/bpg/terraform-provider-proxmox/utils"
)

// testFile represents a file served by TestFileServer.
type testFile struct {
	content      []byte
	reportedSize int64 // if > 0, overrides Content-Length
	filename     string
}

// TestFileServer is a controllable HTTP server for testing file downloads.
// It can serve multiple files at different paths, each with configurable content
// and Content-Length headers. This allows tests to simulate upstream file changes
// (e.g., new cloud image releases) without depending on external services.
//
// The server binds to 0.0.0.0 to be accessible from Proxmox. Set the environment
// variable PROXMOX_VE_ACC_TEST_FILE_SERVER_IP to the IP address that Proxmox
// should use to reach this server.
type TestFileServer struct {
	t        *testing.T
	server   *http.Server
	listener net.Listener
	mu       sync.RWMutex

	// externalIP is the IP that Proxmox will use to reach this server
	externalIP string
	// port is the port the server is listening on
	port int

	// files maps URL path -> file content
	files map[string]*testFile
}

// NewTestFileServer creates a new test file server.
// The server starts serving immediately on a random available port, bound to 0.0.0.0.
//
// The server needs to know the IP that Proxmox can use to reach this machine.
// It tries these sources in order:
//  1. PROXMOX_VE_ACC_TEST_FILE_SERVER_IP environment variable (explicit override)
//  2. Auto-detect: the local IP on the route to the Proxmox node (from PROXMOX_VE_ACC_NODE_SSH_ADDRESS or PROXMOX_VE_ENDPOINT)
//
// Returns nil if the IP cannot be determined.
func NewTestFileServer(t *testing.T) *TestFileServer {
	t.Helper()

	externalIP := utils.GetAnyStringEnv("PROXMOX_VE_ACC_TEST_FILE_SERVER_IP")
	if externalIP == "" {
		externalIP = detectLocalIPForPVE(t)
	}

	if externalIP == "" {
		return nil
	}

	s := &TestFileServer{
		t:          t,
		externalIP: externalIP,
		files:      make(map[string]*testFile),
	}

	// add a default file for backwards compatibility
	s.files["/file"] = &testFile{
		content:  []byte("initial content"),
		filename: "test_file.iso",
	}

	lc := net.ListenConfig{}

	// bind to all interfaces so Proxmox can reach us
	listener, err := lc.Listen(context.Background(), "tcp", "0.0.0.0:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}

	s.listener = listener
	s.port = listener.Addr().(*net.TCPAddr).Port
	s.server = &http.Server{
		Handler:           s,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		if err := s.server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Logf("HTTP server error: %v", err)
		}
	}()

	t.Cleanup(func() {
		s.Close()
	})

	t.Logf("Test file server started at %s (internal: %s)", s.URL(), s.listener.Addr().String())

	return s
}

// ServeHTTP implements http.Handler to serve files from the files map.
func (s *TestFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	file, ok := s.files[r.URL.Path]
	s.mu.RUnlock()

	if !ok {
		http.NotFound(w, r)
		return
	}

	content := file.content
	reportedSize := file.reportedSize
	filename := file.filename

	// determine the size to report in Content-Length
	size := int64(len(content))
	if reportedSize > 0 {
		size = reportedSize
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.FormatInt(size, 10))

	if filename != "" {
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	}

	// for HEAD requests (used by Proxmox query-url-metadata), don't send body
	if r.Method == http.MethodHead {
		return
	}

	// if reportedSize is larger than actual content, pad with zeros
	if reportedSize > int64(len(content)) {
		padded := make([]byte, reportedSize)
		copy(padded, content)
		content = padded
	}

	if _, err := w.Write(content); err != nil {
		s.t.Logf("Failed to write content: %v", err)
	}
}

// URL returns the base URL that Proxmox should use to download files.
func (s *TestFileServer) URL() string {
	return "http://" + net.JoinHostPort(s.externalIP, strconv.Itoa(s.port))
}

// FileURL returns the URL to download the default test file at /file.
func (s *TestFileServer) FileURL() string {
	return s.URL() + "/file"
}

// AddFile adds a file to be served at the given path.
// The path should start with "/" (e.g., "/fake_file.iso").
// Returns the full URL to access the file.
func (s *TestFileServer) AddFile(path, filename string, content []byte) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.files[path] = &testFile{
		content:  content,
		filename: filename,
	}

	return s.URL() + path
}

// SetContent sets the content for the default file at /file.
func (s *TestFileServer) SetContent(content []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if f, ok := s.files["/file"]; ok {
		f.content = content
		f.reportedSize = 0
	}
}

// SetReportedSize sets the Content-Length header for the default file at /file.
func (s *TestFileServer) SetReportedSize(size int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if f, ok := s.files["/file"]; ok {
		f.reportedSize = size
	}
}

// SetFilename sets the filename for the default file at /file.
func (s *TestFileServer) SetFilename(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if f, ok := s.files["/file"]; ok {
		f.filename = name
	}
}

// SetFileReportedSize sets the Content-Length header for a specific file path.
func (s *TestFileServer) SetFileReportedSize(path string, size int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if f, ok := s.files[path]; ok {
		f.reportedSize = size
	}
}

// GetActualSize returns the actual content length of the default file.
func (s *TestFileServer) GetActualSize() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if f, ok := s.files["/file"]; ok {
		return int64(len(f.content))
	}

	return 0
}

// GetFileSHA256 returns the SHA256 checksum of a file's content.
func (s *TestFileServer) GetFileSHA256(path string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if f, ok := s.files[path]; ok {
		hash := sha256.Sum256(f.content)
		return hex.EncodeToString(hash[:])
	}

	return ""
}

// Close shuts down the test server.
func (s *TestFileServer) Close() {
	if s.server != nil {
		_ = s.server.Close()
	}
}

// detectLocalIPForPVE finds the local IP address that routes to the Proxmox node.
// It uses a UDP dial (no actual traffic) to determine which local interface would be
// used to reach the PVE host.
func detectLocalIPForPVE(t *testing.T) string {
	t.Helper()

	// Try SSH address first (most direct), then parse endpoint URL
	pveHost := utils.GetAnyStringEnv("PROXMOX_VE_ACC_NODE_SSH_ADDRESS")
	if pveHost == "" {
		endpoint := utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT")
		if endpoint != "" {
			if u, err := url.Parse(endpoint); err == nil {
				pveHost = u.Hostname()
			}
		}
	}

	if pveHost == "" {
		t.Log("Cannot auto-detect test file server IP: no PVE host configured")
		return ""
	}

	dialer := net.Dialer{}

	conn, err := dialer.DialContext(context.Background(), "udp", net.JoinHostPort(pveHost, "80"))
	if err != nil {
		t.Logf("Cannot auto-detect test file server IP: %v", err)
		return ""
	}

	defer func() { _ = conn.Close() }()

	localIP := conn.LocalAddr().(*net.UDPAddr).IP.String()
	t.Logf("Auto-detected test file server IP: %s (route to %s)", localIP, pveHost)

	return localIP
}
