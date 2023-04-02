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

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type Options interface {
	GetOptionsID() string
	SetOptions(ctx context.Context, d *OptionsPutRequestBody) error
	GetOptions(ctx context.Context) (*OptionsGetResponseData, error)
}

type OptionsPutRequestBody struct {
	DHCP        *types.CustomBool `json:"dhcp,omitempty"          url:"dhcp,omitempty,int"`
	Enable      *types.CustomBool `json:"enable,omitempty"        url:"enable,omitempty,int"`
	IPFilter    *types.CustomBool `json:"ipfilter,omitempty"      url:"ipfilter,omitempty,int"`
	LogLevelIN  *string           `json:"log_level_in,omitempty"  url:"log_level_in,omitempty"`
	LogLevelOUT *string           `json:"log_level_out,omitempty" url:"log_level_out,omitempty"`
	MACFilter   *types.CustomBool `json:"macfilter,omitempty"     url:"macfilter,omitempty,int"`
	NDP         *types.CustomBool `json:"ndp,omitempty"           url:"ndp,omitempty,int"`
	PolicyIn    *string           `json:"policy_in,omitempty"     url:"policy_in,omitempty"`
	PolicyOut   *string           `json:"policy_out,omitempty"    url:"policy_out,omitempty"`
	RAdv        *types.CustomBool `json:"radv,omitempty"          url:"radv,omitempty,int"`
}

type OptionsGetResponseBody struct {
	Data *OptionsGetResponseData `json:"data,omitempty"`
}

type OptionsGetResponseData struct {
	DHCP        *types.CustomBool `json:"dhcp"          url:"dhcp,int"`
	Enable      *types.CustomBool `json:"enable"        url:"enable,int"`
	IPFilter    *types.CustomBool `json:"ipfilter"      url:"ipfilter,int"`
	LogLevelIN  *string           `json:"log_level_in"  url:"log_level_in"`
	LogLevelOUT *string           `json:"log_level_out" url:"log_level_out"`
	MACFilter   *types.CustomBool `json:"macfilter"     url:"macfilter,int"`
	NDP         *types.CustomBool `json:"ndp"           url:"ndp,int"`
	PolicyIn    *string           `json:"policy_in"     url:"policy_in"`
	PolicyOut   *string           `json:"policy_out"    url:"policy_out"`
	RAdv        *types.CustomBool `json:"radv"          url:"radv,int"`
}

func (c *Client) optionsPath() string {
	return c.ExpandPath("firewall/options")
}

func (c *Client) GetOptionsID() string {
	return "options-" + strconv.Itoa(schema.HashString(c.optionsPath()))
}

func (c *Client) SetOptions(ctx context.Context, d *OptionsPutRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, c.optionsPath(), d, nil)
	if err != nil {
		return fmt.Errorf("error setting optionss: %w", err)
	}
	return nil
}

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
