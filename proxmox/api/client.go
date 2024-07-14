/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/avast/retry-go/v4"
	"github.com/google/go-querystring/query"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"golang.org/x/exp/maps"

	"github.com/bpg/terraform-provider-proxmox/utils"
)

const (
	basePathJSONAPI = "api2/json"
)

// Client is an interface for performing requests against the Proxmox API.
type Client interface {
	// DoRequest performs a request against the Proxmox API.
	DoRequest(
		ctx context.Context,
		method, path string,
		requestBody, responseBody interface{},
	) error

	// ExpandPath expands a path relative to the client's base path.
	// For example, if the client is configured for a VM and the
	// path is "firewall/options", the returned path will be
	// "/nodes/<node>/qemu/<vmid>/firewall/options".
	ExpandPath(path string) string

	// IsRoot returns true if the client is configured with the root user.
	IsRoot() bool

	// IsRootTicket returns true if the authenticator is configured to use the root directly using a login ticket.
	// (root using token is weaker, cannot change VM arch)
	IsRootTicket() bool

	// HTTP returns a lower-level HTTP client.
	HTTP() *http.Client
}

// Connection represents a connection to the Proxmox Virtual Environment API.
type Connection struct {
	endpoint   string
	httpClient *http.Client
}

// NewConnection creates and initializes a Connection instance.
func NewConnection(endpoint string, insecure bool, minTLS string) (*Connection, error) {
	u, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return nil, errors.New(
			"you must specify a valid endpoint for the Proxmox Virtual Environment API (valid: https://host:port/)",
		)
	}

	if u.Scheme != "https" {
		return nil, errors.New(
			"you must specify a secure endpoint for the Proxmox Virtual Environment API (valid: https://host:port/)",
		)
	}

	version, err := GetMinTLSVersion(minTLS)
	if err != nil {
		return nil, err
	}

	var transport http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			// deepcode ignore InsecureTLSConfig: the min TLS version is configurable
			MinVersion:         version,
			InsecureSkipVerify: insecure, //nolint:gosec
		},
	}

	if logging.IsDebugOrHigher() {
		transport = logging.NewLoggingHTTPTransport(transport)
	}

	// make sure the path does not contain "/api2/json"
	u.Path = ""

	return &Connection{
		endpoint: strings.TrimRight(u.String(), "/"),
		httpClient: &http.Client{
			Transport: transport,
		},
	}, nil
}

// VirtualEnvironmentClient implements an API client for the Proxmox Virtual Environment API.
type client struct {
	conn *Connection
	auth Authenticator
}

// NewClient creates and initializes a VirtualEnvironmentClient instance.
func NewClient(creds *Credentials, conn *Connection) (Client, error) {
	if creds == nil {
		return nil, errors.New("credentials must not be nil")
	}

	if conn == nil {
		return nil, errors.New("connection must not be nil")
	}

	var auth Authenticator

	var err error

	// todo maybe move NewTicketAuthenticator cred-input-logic to here
	// aka: creds.AuthTicket and  creds.AuthTicket != "" && creds.CSRFPreventionToken to here
	if creds.APIToken != nil {
		auth, err = NewTokenAuthenticator(*creds.APIToken)
	} else {
		auth, err = NewTicketAuthenticator(conn, creds)
	}

	if err != nil {
		return nil, err
	}

	return &client{
		conn: conn,
		auth: auth,
	}, nil
}

// DoRequest performs a HTTP request against a JSON API endpoint.
func (c *client) DoRequest(
	ctx context.Context,
	method, path string,
	requestBody, responseBody interface{},
) error {
	var reqBodyReader io.Reader

	var reqContentLength *int64

	modifiedPath := path
	reqBodyType := ""

	//nolint:nestif
	if requestBody != nil {
		multipartData, multipart := requestBody.(*MultiPartData)
		pipedBodyReader, pipedBody := requestBody.(*io.PipeReader)

		switch {
		case multipart:
			reqBodyReader = multipartData.Reader
			reqBodyType = fmt.Sprintf("multipart/form-data; boundary=%s", multipartData.Boundary)
			reqContentLength = multipartData.Size
		case pipedBody:
			reqBodyReader = pipedBodyReader
		default:
			v, err := query.Values(requestBody)
			if err != nil {
				return fmt.Errorf("failed to encode HTTP %s request (path: %s) - Reason: %w",
					method,
					modifiedPath,
					err,
				)
			}

			encodedValues := v.Encode()
			if encodedValues != "" {
				if method == http.MethodDelete || method == http.MethodGet || method == http.MethodHead {
					if !strings.Contains(modifiedPath, "?") {
						modifiedPath = fmt.Sprintf("%s?%s", modifiedPath, encodedValues)
					} else {
						modifiedPath = fmt.Sprintf("%s&%s", modifiedPath, encodedValues)
					}
				} else {
					reqBodyReader = bytes.NewBufferString(encodedValues)
					reqBodyType = "application/x-www-form-urlencoded"
				}
			}
		}
	} else {
		reqBodyReader = new(bytes.Buffer)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		fmt.Sprintf("%s/%s/%s", c.conn.endpoint, basePathJSONAPI, modifiedPath),
		reqBodyReader,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to create HTTP %s request (path: %s) - Reason: %w",
			method,
			modifiedPath,
			err,
		)
	}

	req.Header.Add("Accept", "application/json")

	if reqContentLength != nil {
		req.ContentLength = *reqContentLength
	}

	if reqBodyType != "" {
		req.Header.Add("Content-Type", reqBodyType)
	}

	err = c.auth.AuthenticateRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to authenticate HTTP %s request (path: %s) - Reason: %w",
			method,
			modifiedPath,
			err,
		)
	}

	//nolint:bodyclose
	res, err := retry.DoWithData(
		func() (*http.Response, error) {
			return c.conn.httpClient.Do(req)
		},
		retry.Context(ctx),
		retry.RetryIf(func(err error) bool {
			var urlErr *url.Error
			if errors.As(err, &urlErr) {
				return strings.ToUpper(urlErr.Op) == http.MethodGet
			}

			return false
		}),
		retry.LastErrorOnly(true),
		retry.Attempts(3),
	)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP %s request (path: %s) - Reason: %w",
			method,
			modifiedPath,
			err,
		)
	}

	defer utils.CloseOrLogError(ctx)(res.Body)

	err = validateResponseCode(res)
	if err != nil {
		return err
	}

	//nolint:nestif
	if responseBody != nil {
		err = json.NewDecoder(res.Body).Decode(responseBody)
		if err != nil {
			return fmt.Errorf(
				"failed to decode HTTP %s response (path: %s) - Reason: %w",
				method,
				modifiedPath,
				err,
			)
		}
	} else {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf(
				"failed to read HTTP %s response body (path: %s) - Reason: %w",
				method,
				modifiedPath,
				err,
			)
		}

		if len(data) > 0 {
			dr := dataResponse{}

			if err2 := json.NewDecoder(bytes.NewReader(data)).Decode(&dr); err2 == nil {
				if dr.Data == nil {
					return nil
				}
			}

			tflog.Warn(ctx, "unhandled HTTP response body", map[string]interface{}{
				"data": dr.Data,
			})
		}
	}

	return nil
}

type dataResponse struct {
	Data interface{} `json:"data"`
}

// ExpandPath expands the given path to an absolute path.
func (c *client) ExpandPath(path string) string {
	return path
}

func (c *client) IsRoot() bool {
	return c.auth.IsRoot()
}

func (c *client) IsRootTicket() bool {
	return c.auth.IsRootTicket()
}

func (c *client) HTTP() *http.Client {
	return c.conn.httpClient
}

// validateResponseCode ensures that a response is valid.
func validateResponseCode(res *http.Response) error {
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		if res.StatusCode == http.StatusNotFound ||
			(res.StatusCode == http.StatusInternalServerError && strings.Contains(res.Status, "does not exist")) {
			return ErrResourceDoesNotExist
		}

		msg := strings.TrimPrefix(res.Status, fmt.Sprintf("%d ", res.StatusCode))

		errRes := &ErrorResponseBody{}
		err := json.NewDecoder(res.Body).Decode(errRes)

		if err == nil && errRes.Errors != nil {
			var errList []string

			for k, v := range *errRes.Errors {
				errList = append(errList, fmt.Sprintf("%s: %s", k, strings.TrimRight(v, "\n\r")))
			}

			msg = fmt.Sprintf("%s (%s)", msg, strings.Join(errList, " - "))
		}

		return &HTTPError{
			Code:    res.StatusCode,
			Message: msg,
		}
	}

	return nil
}

// GetMinTLSVersion returns the minimum TLS version constant for the given string. If the string is empty,
// the default TLS version is returned. For unsupported TLS versions, an error is returned.
func GetMinTLSVersion(version string) (uint16, error) {
	validVersions := map[string]uint16{
		"":    tls.VersionTLS13,
		"1.3": tls.VersionTLS13,
		"1.2": tls.VersionTLS12,
		"1.1": tls.VersionTLS11,
		"1.0": tls.VersionTLS10,
	}

	if val, ok := validVersions[strings.TrimSpace(version)]; ok {
		return val, nil
	}

	valid := maps.Keys(validVersions)
	sort.Strings(valid)

	return 0, fmt.Errorf("unsupported minimal TLS version %s, must be one of: %s", version, strings.Join(valid, ", "))
}
