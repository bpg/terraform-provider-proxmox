/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// OptionsGetResponseBody is the response body for the GET /nodes/{node}/firewall/options API call.
type OptionsGetResponseBody struct {
	Data *OptionsGetResponseData `json:"data,omitempty"`
}

// OptionsGetResponseData is the data field of the response body for the GET /nodes/{node}/firewall/options API call.
type OptionsGetResponseData struct {
	Enable                           *types.CustomBool `json:"enable,omitempty"                               url:"enable,omitempty,int"`
	LogLevelIn                       *string           `json:"log_level_in,omitempty"                         url:"log_level_in,omitempty"`
	LogLevelOut                      *string           `json:"log_level_out,omitempty"                        url:"log_level_out,omitempty"`
	LogLevelForward                  *string           `json:"log_level_forward,omitempty"                    url:"log_level_forward,omitempty"`
	LogNFConntrack                   *types.CustomBool `json:"log_nf_conntrack,omitempty"                     url:"log_nf_conntrack,omitempty,int"`
	NDP                              *types.CustomBool `json:"ndp,omitempty"                                  url:"ndp,omitempty,int"`
	NFConntrackAllowInvalid          *types.CustomBool `json:"nf_conntrack_allow_invalid,omitempty"           url:"nf_conntrack_allow_invalid,omitempty,int"`
	NFConntrackHelpers               *string           `json:"nf_conntrack_helpers,omitempty"                 url:"nf_conntrack_helpers,omitempty"`
	NFConntrackMax                   *int64            `json:"nf_conntrack_max,omitempty"                     url:"nf_conntrack_max,omitempty"`
	NFConntrackTCPTimeoutEstablished *int64            `json:"nf_conntrack_tcp_timeout_established,omitempty" url:"nf_conntrack_tcp_timeout_established,omitempty"`
	NFConntrackTCPTimeoutSynRecv     *int64            `json:"nf_conntrack_tcp_timeout_syn_recv,omitempty"    url:"nf_conntrack_tcp_timeout_syn_recv,omitempty"`
	NFTables                         *types.CustomBool `json:"nftables,omitempty"                             url:"nftables,omitempty,int"`
	NoSMURFs                         *types.CustomBool `json:"nosmurfs,omitempty"                             url:"nosmurfs,omitempty,int"`
	ProtectionSynflood               *types.CustomBool `json:"protection_synflood,omitempty"                  url:"protection_synflood,omitempty,int"`
	ProtectionSynfloodBurst          *int64            `json:"protection_synflood_burst,omitempty"            url:"protection_synflood_burst,omitempty"`
	ProtectionSynfloodRate           *int64            `json:"protection_synflood_rate,omitempty"             url:"protection_synflood_rate,omitempty"`
	SMURFLogLevel                    *string           `json:"smurf_log_level,omitempty"                      url:"smurf_log_level,omitempty"`
	TCPFlagsLogLevel                 *string           `json:"tcp_flags_log_level,omitempty"                  url:"tcp_flags_log_level,omitempty"`
}

// OptionsPutRequestBody is the request body for the PUT /nodes/{node}/firewall/options API call.
type OptionsPutRequestBody struct {
	OptionsGetResponseData

	Delete *[]string `url:"delete,omitempty"`
}
