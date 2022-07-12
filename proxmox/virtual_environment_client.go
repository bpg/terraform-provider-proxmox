/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

// NewVirtualEnvironmentClient creates and initializes a VirtualEnvironmentClient instance.
func NewVirtualEnvironmentClient(endpoint, username, password, otp string, insecure bool) (*VirtualEnvironmentClient, error) {
	u, err := url.ParseRequestURI(endpoint)

	if err != nil {
		return nil, errors.New("you must specify a valid endpoint for the Proxmox Virtual Environment API (valid: https://host:port/)")
	}

	if u.Scheme != "https" {
		return nil, errors.New("you must specify a secure endpoint for the Proxmox Virtual Environment API (valid: https://host:port/)")
	}

	if password == "" {
		return nil, errors.New("you must specify a password for the Proxmox Virtual Environment API")
	}

	if username == "" {
		return nil, errors.New("you must specify a username for the Proxmox Virtual Environment API")
	}

	var pOTP *string

	if otp != "" {
		pOTP = &otp
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
			},
		},
	}

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
func (c *VirtualEnvironmentClient) DoRequest(method, path string, requestBody interface{}, responseBody interface{}) error {
	var reqBodyReader io.Reader
	var reqContentLength *int64

	log.Printf("[DEBUG] Performing HTTP %s request (path: %s)", method, path)

	modifiedPath := path
	reqBodyType := ""

	if requestBody != nil {
		multipartData, multipart := requestBody.(*VirtualEnvironmentMultiPartData)
		pipedBodyReader, pipedBody := requestBody.(*io.PipeReader)

		if multipart {
			reqBodyReader = multipartData.Reader
			reqBodyType = fmt.Sprintf("multipart/form-data; boundary=%s", multipartData.Boundary)
			reqContentLength = multipartData.Size

			log.Printf("[DEBUG] Added multipart request body to HTTP %s request (path: %s)", method, modifiedPath)
		} else if pipedBody {
			reqBodyReader = pipedBodyReader

			log.Printf("[DEBUG] Added piped request body to HTTP %s request (path: %s)", method, modifiedPath)
		} else {
			v, err := query.Values(requestBody)

			if err != nil {
				fErr := fmt.Errorf("failed to encode HTTP %s request (path: %s) - Reason: %s", method, modifiedPath, err.Error())
				log.Printf("[DEBUG] WARNING: %s", fErr.Error())
				return fErr
			}

			encodedValues := v.Encode()

			if encodedValues != "" {
				if method == hmDELETE || method == hmGET || method == hmHEAD {
					if !strings.Contains(modifiedPath, "?") {
						modifiedPath = fmt.Sprintf("%s?%s", modifiedPath, encodedValues)
					} else {
						modifiedPath = fmt.Sprintf("%s&%s", modifiedPath, encodedValues)
					}
				} else {
					reqBodyReader = bytes.NewBufferString(encodedValues)
					reqBodyType = "application/x-www-form-urlencoded"
				}

				log.Printf("[DEBUG] Added request body to HTTP %s request (path: %s) - Body: %s", method, modifiedPath, encodedValues)
			}
		}
	} else {
		reqBodyReader = new(bytes.Buffer)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s/%s", c.Endpoint, basePathJSONAPI, modifiedPath), reqBodyReader)

	if err != nil {
		fErr := fmt.Errorf("failed to create HTTP %s request (path: %s) - Reason: %s", method, modifiedPath, err.Error())
		log.Printf("[DEBUG] WARNING: %s", fErr.Error())
		return fErr
	}

	req.Header.Add("Accept", "application/json")

	if reqContentLength != nil {
		req.ContentLength = *reqContentLength
	}

	if reqBodyType != "" {
		req.Header.Add("Content-Type", reqBodyType)
	}

	err = c.AuthenticateRequest(req)

	if err != nil {
		log.Printf("[DEBUG] WARNING: %s", err.Error())
		return err
	}

	res, err := c.httpClient.Do(req)

	if err != nil {
		fErr := fmt.Errorf("failed to perform HTTP %s request (path: %s) - Reason: %s", method, modifiedPath, err.Error())
		log.Printf("[DEBUG] WARNING: %s", fErr.Error())
		return fErr
	}

	defer res.Body.Close()

	err = c.ValidateResponseCode(res)

	if err != nil {
		log.Printf("[DEBUG] WARNING: %s", err.Error())
		return err
	}

	if responseBody != nil {
		err = json.NewDecoder(res.Body).Decode(responseBody)

		if err != nil {
			fErr := fmt.Errorf("failed to decode HTTP %s response (path: %s) - Reason: %s", method, modifiedPath, err.Error())
			log.Printf("[DEBUG] WARNING: %s", fErr.Error())
			return fErr
		}
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		log.Printf("[DEBUG] WARNING: Unhandled HTTP response body: %s", string(data))
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
