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
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"

	"github.com/bpg/terraform-provider-proxmox/utils"
)

// ErrNoDataObjectInResponse is returned when the server does not include a data object in the response.
var ErrNoDataObjectInResponse = errors.New("the server did not include a data object in the response")

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
}

// Connection represents a connection to the Proxmox Virtual Environment API.
type Connection struct {
	endpoint   string
	httpClient *http.Client
}

// NewConnection creates and initializes a Connection instance.
func NewConnection(endpoint string, insecure bool) (*Connection, error) {
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

	var transport http.RoundTripper = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecure, //nolint:gosec
		},
	}

	if logging.IsDebugOrHigher() {
		transport = logging.NewLoggingHTTPTransport(transport)
	}

	return &Connection{
		endpoint:   strings.TrimRight(u.String(), "/"),
		httpClient: &http.Client{Transport: transport},
	}, nil
}

// VirtualEnvironmentClient implements an API client for the Proxmox Virtual Environment API.
type client struct {
	conn *Connection
	auth Authenticator
}

// NewClient creates and initializes a VirtualEnvironmentClient instance.
func NewClient(ctx context.Context, creds *Credentials, conn *Connection) (Client, error) {
	if creds == nil {
		return nil, errors.New("credentials must not be nil")
	}

	if conn == nil {
		return nil, errors.New("connection must not be nil")
	}

	var auth Authenticator

	var err error

	if creds.APIToken != nil {
		auth, err = NewTokenAuthenticator(*creds.APIToken)
	} else {
		auth, err = NewTicketAuthenticator(ctx, conn, creds)
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

	err = c.auth.AuthenticateRequest(req)
	if err != nil {
		return fmt.Errorf("failed to authenticate HTTP %s request (path: %s) - Reason: %w",
			method,
			modifiedPath,
			err,
		)
	}

	//nolint:bodyclose
	res, err := c.conn.httpClient.Do(req)
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
		tflog.Warn(ctx, "unhandled HTTP response body", map[string]interface{}{
			"data": string(data),
		})
	}

	return nil
}

// ExpandPath expands the given path to an absolute path.
func (c *client) ExpandPath(path string) string {
	return path
}

func (c *client) IsRoot() bool {
	return c.auth.IsRoot()
}

// validateResponseCode ensures that a response is valid.
func validateResponseCode(res *http.Response) error {
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		status := strings.TrimPrefix(res.Status, fmt.Sprintf("%d ", res.StatusCode))

		errRes := &ErrorResponseBody{}
		err := json.NewDecoder(res.Body).Decode(errRes)

		if err == nil && errRes.Errors != nil {
			var errList []string

			for k, v := range *errRes.Errors {
				errList = append(errList, fmt.Sprintf("%s: %s", k, strings.TrimRight(v, "\n\r")))
			}

			status = fmt.Sprintf("%s (%s)", status, strings.Join(errList, " - "))
		}

		return fmt.Errorf("received an HTTP %d response - Reason: %s", res.StatusCode, status)
	}

	return nil
}
