/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

// OptionsPutRequestBody is the request body for the PUT /cluster/firewall/options API call
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

// OptionsGetResponseBody is the response body for the GET /cluster/firewall/options API call
type OptionsGetResponseBody struct {
	Data *OptionsGetResponseData `json:"data,omitempty"`
}

// OptionsGetResponseData is the data field of the response body for the GET /cluster/firewall/options API call
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
