/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// DeleteCertificate deletes the custom certificate for a node.
func (c *Client) DeleteCertificate(ctx context.Context, d *CertificateDeleteRequestBody) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath("certificates/custom"), d, nil)
	if err != nil {
		return fmt.Errorf("error deleting certificate: %w", err)
	}

	return nil
}

// ListCertificates retrieves the list of certificates for a node.
func (c *Client) ListCertificates(ctx context.Context) (*[]CertificateListResponseData, error) {
	resBody := &CertificateListResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath("certificates/info"), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving certificate list: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// UpdateCertificate updates the custom certificate for a node.
func (c *Client) UpdateCertificate(ctx context.Context, d *CertificateUpdateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("certificates/custom"), d, nil)
	if err != nil {
		return fmt.Errorf("error updating certificate: %w", err)
	}

	return nil
}

// OrderCertificate orders a new certificate from ACME CA for a node.
func (c *Client) OrderCertificate(ctx context.Context, d *CertificateOrderRequestBody) (*string, error) {
	resBody := &CertificateOrderResponseBody{}

	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("certificates/acme/certificate"), d, resBody)
	if err != nil {
		return nil, fmt.Errorf("error ordering ACME certificate: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// RenewCertificate renews an existing ACME certificate for a node.
func (c *Client) RenewCertificate(ctx context.Context, d *CertificateRenewRequestBody) (*string, error) {
	resBody := &CertificateOrderResponseBody{}

	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath("certificates/acme/certificate"), d, resBody)
	if err != nil {
		return nil, fmt.Errorf("error renewing ACME certificate: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}
