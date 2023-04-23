/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// DeleteCertificate deletes the custom certificate for a node.
func (c *VirtualEnvironmentClient) DeleteCertificate(
	ctx context.Context,
	nodeName string,
	d *CertificateDeleteRequestBody,
) error {
	return c.DoRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("nodes/%s/certificates/custom", url.PathEscape(nodeName)),
		d,
		nil,
	)
}

// ListCertificates retrieves the list of certificates for a node.
func (c *VirtualEnvironmentClient) ListCertificates(
	ctx context.Context,
	nodeName string,
) (*[]CertificateListResponseData, error) {
	resBody := &CertificateListResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("nodes/%s/certificates/info", url.PathEscape(nodeName)),
		nil,
		resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving certificate list: %w", err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// UpdateCertificate updates the custom certificate for a node.
func (c *VirtualEnvironmentClient) UpdateCertificate(
	ctx context.Context,
	nodeName string,
	d *CertificateUpdateRequestBody,
) error {
	return c.DoRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("nodes/%s/certificates/custom", url.PathEscape(nodeName)),
		d,
		nil,
	)
}
