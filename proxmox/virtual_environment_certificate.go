/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"errors"
	"fmt"
	"net/url"
)

// DeleteCertificate deletes the custom certificate for a node.
func (c *VirtualEnvironmentClient) DeleteCertificate(nodeName string, d *VirtualEnvironmentCertificateDeleteRequestBody) error {
	return c.DoRequest(hmDELETE, fmt.Sprintf("nodes/%s/certificates/custom", url.PathEscape(nodeName)), d, nil)
}

// ListCertificates retrieves the list of certificates for a node.
func (c *VirtualEnvironmentClient) ListCertificates(nodeName string) (*[]VirtualEnvironmentCertificateListResponseData, error) {
	resBody := &VirtualEnvironmentCertificateListResponseBody{}
	err := c.DoRequest(hmGET, fmt.Sprintf("nodes/%s/certificates/info", url.PathEscape(nodeName)), nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// UpdateCertificate updates the custom certificate for a node.
func (c *VirtualEnvironmentClient) UpdateCertificate(nodeName string, d *VirtualEnvironmentCertificateUpdateRequestBody) error {
	return c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/certificates/custom", url.PathEscape(nodeName)), d, nil)
}
