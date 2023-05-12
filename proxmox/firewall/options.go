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
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Options is an interface for the Proxmox firewall options API
type Options interface {
	GetOptionsID() string
	SetOptions(ctx context.Context, d *OptionsPutRequestBody) error
	GetOptions(ctx context.Context) (*OptionsGetResponseData, error)
}

func (c *Client) optionsPath() string {
	return c.ExpandPath("firewall/options")
}

// GetOptionsID returns the ID of the options object
func (c *Client) GetOptionsID() string {
	return "options-" + strconv.Itoa(schema.HashString(c.optionsPath()))
}

// SetOptions sets the options object
func (c *Client) SetOptions(ctx context.Context, d *OptionsPutRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, c.optionsPath(), d, nil)
	if err != nil {
		return fmt.Errorf("error setting optionss: %w", err)
	}
	return nil
}

// GetOptions retrieves the options object
func (c *Client) GetOptions(ctx context.Context) (*OptionsGetResponseData, error) {
	resBody := &OptionsGetResponseBody{}
	err := c.DoRequest(ctx, http.MethodGet, c.optionsPath(), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving options: %w", err)
	}

	if resBody.Data == nil {
		return nil, fmt.Errorf("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}
