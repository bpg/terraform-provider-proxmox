/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
)

func (c *Client) SetOptions(ctx context.Context, d *firewall.OptionsPutRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, "cluster/firewall/options", d, nil)
	if err != nil {
		return fmt.Errorf("error setting optionss: %w", err)
	}
	return nil
}

func (c *Client) GetOptions(ctx context.Context) (*firewall.OptionsGetResponseData, error) {
	resBody := &firewall.OptionsGetResponseBody{}
	err := c.DoRequest(ctx, http.MethodGet, "cluster/firewall/options", nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving options: %w", err)
	}

	if resBody.Data == nil {
		return nil, fmt.Errorf("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}
