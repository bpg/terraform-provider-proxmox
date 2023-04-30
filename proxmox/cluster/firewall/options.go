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
)

type Options interface {
	SetGlobalOptions(ctx context.Context, d *OptionsPutRequestBody) error
	GetGlobalOptions(ctx context.Context) (*OptionsGetResponseData, error)
}

func (c *Client) SetGlobalOptions(ctx context.Context, d *OptionsPutRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, "cluster/firewall/options", d, nil)
	if err != nil {
		return fmt.Errorf("error setting optionss: %w", err)
	}
	return nil
}

func (c *Client) GetGlobalOptions(ctx context.Context) (*OptionsGetResponseData, error) {
	resBody := &OptionsGetResponseBody{}
	err := c.DoRequest(ctx, http.MethodGet, "cluster/firewall/options", nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving options: %w", err)
	}

	if resBody.Data == nil {
		return nil, fmt.Errorf("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}
