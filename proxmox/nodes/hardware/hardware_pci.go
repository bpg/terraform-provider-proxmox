/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package hardware

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// ListPCIDevices retrieves the list of PCI devices on the node.
func (c *Client) ListPCIDevices(
	ctx context.Context,
	d *ListPCIDevicesRequestBody,
) ([]*PCIDeviceData, error) {
	resBody := &ListPCIDevicesResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath("pci"), d, resBody)
	if err != nil {
		return nil, fmt.Errorf("error listing PCI devices: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].ID < resBody.Data[j].ID
	})

	return resBody.Data, nil
}
