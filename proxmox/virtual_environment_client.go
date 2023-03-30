/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

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
)

// NewVirtualEnvironmentClient creates and initializes a VirtualEnvironmentClient instance.
func NewVirtualEnvironmentClient(
	endpoint, username, password, otp string,
	insecure bool,
) (*VirtualEnvironmentClient, error) {
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

	return &VirtualEnvironmentClient{
		Endpoint:   strings.TrimRight(u.String(), "/"),
		Insecure:   insecure,
		OTP:        pOTP,
		Password:   password,
		Username:   username,
		httpClient: httpClient,
	}, nil
}

// DoRequest performs a HTTP request against a JSON API endpoint.
func (c *VirtualEnvironmentClient) DoRequest(
	ctx context.Context,
	method, path string,
	requestBody, responseBody interface{},
) error {
	var reqBodyReader io.Reader
	var reqContentLength *int64

	modifiedPath := path
	reqBodyType := ""

	if requestBody != nil {
		multipartData, multipart := requestBody.(*VirtualEnvironmentMultiPartData)
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
				if method == HmDELETE || method == HmGET || method == HmHEAD {
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
		fmt.Sprintf("%s/%s/%s", c.Endpoint, basePathJSONAPI, modifiedPath),
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

	defer CloseOrLogError(ctx)(res.Body)

	err = c.ValidateResponseCode(res)
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

// ValidateResponseCode ensures that a response is valid.
func (c *VirtualEnvironmentClient) ValidateResponseCode(res *http.Response) error {
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		status := strings.TrimPrefix(res.Status, fmt.Sprintf("%d ", res.StatusCode))

		errRes := &VirtualEnvironmentErrorResponseBody{}
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
