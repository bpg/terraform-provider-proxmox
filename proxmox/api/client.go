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

const (
	basePathJSONAPI = "api2/json"
)

// VirtualEnvironmentClient implements an API client for the Proxmox Virtual Environment API.
type client struct {
	endpoint string
	insecure bool
	otp      *string
	password string
	username string

	authenticationData *AuthenticationResponseData
	httpClient         *http.Client
}

// NewClient creates and initializes a VirtualEnvironmentClient instance.
func NewClient(
	endpoint, username, password, otp string, insecure bool,
) (Client, error) {
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

	if password == "" {
		return nil, errors.New(
			"you must specify a password for the Proxmox Virtual Environment API",
		)
	}

	if username == "" {
		return nil, errors.New(
			"you must specify a username for the Proxmox Virtual Environment API",
		)
	}

	if !strings.Contains(username, "@") {
		return nil, errors.New(
			"make sure the username for the Proxmox Virtual Environment API ends in '@pve or @pam'",
		)
	}

	var pOTP *string

	if otp != "" {
		pOTP = &otp
	}

	var transport http.RoundTripper = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecure, //nolint:gosec
		},
	}

	if logging.IsDebugOrHigher() {
		transport = logging.NewLoggingHTTPTransport(transport)
	}

	httpClient := &http.Client{Transport: transport}

	return &client{
		endpoint:   strings.TrimRight(u.String(), "/"),
		insecure:   insecure,
		otp:        pOTP,
		password:   password,
		username:   username,
		httpClient: httpClient,
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

	if requestBody != nil {
		multipartData, multipart := requestBody.(*MultiPartData)
		pipedBodyReader, pipedBody := requestBody.(*io.PipeReader)

		if multipart {
			reqBodyReader = multipartData.Reader
			reqBodyType = fmt.Sprintf("multipart/form-data; boundary=%s", multipartData.Boundary)
			reqContentLength = multipartData.Size
		} else if pipedBody {
			reqBodyReader = pipedBodyReader
		} else {
			v, err := query.Values(requestBody)
			if err != nil {
				fErr := fmt.Errorf("failed to encode HTTP %s request (path: %s) - Reason: %s", method, modifiedPath, err.Error())
				tflog.Warn(ctx, fErr.Error())
				return fErr
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
		fmt.Sprintf("%s/%s/%s", c.endpoint, basePathJSONAPI, modifiedPath),
		reqBodyReader,
	)
	if err != nil {
		fErr := fmt.Errorf(
			"failed to create HTTP %s request (path: %s) - Reason: %w",
			method,
			modifiedPath,
			err,
		)
		tflog.Warn(ctx, fErr.Error())
		return fErr
	}

	req.Header.Add("Accept", "application/json")

	if reqContentLength != nil {
		req.ContentLength = *reqContentLength
	}

	if reqBodyType != "" {
		req.Header.Add("Content-Type", reqBodyType)
	}

	err = c.AuthenticateRequest(ctx, req)

	if err != nil {
		tflog.Warn(ctx, err.Error())
		return err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		fErr := fmt.Errorf(
			"failed to perform HTTP %s request (path: %s) - Reason: %w",
			method,
			modifiedPath,
			err,
		)
		tflog.Warn(ctx, fErr.Error())
		return fErr
	}

	defer utils.CloseOrLogError(ctx)(res.Body)

	err = c.validateResponseCode(res)
	if err != nil {
		tflog.Warn(ctx, err.Error())
		return err
	}

	if responseBody != nil {
		err = json.NewDecoder(res.Body).Decode(responseBody)

		if err != nil {
			fErr := fmt.Errorf(
				"failed to decode HTTP %s response (path: %s) - Reason: %w",
				method,
				modifiedPath,
				err,
			)
			tflog.Warn(ctx, fErr.Error())
			return fErr
		}
	} else {
		data, _ := io.ReadAll(res.Body)
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
	return c.username == "root@pam"
}

// validateResponseCode ensures that a response is valid.
func (c *client) validateResponseCode(res *http.Response) error {
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
