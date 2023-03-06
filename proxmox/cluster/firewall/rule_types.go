/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

// RuleCreateRequestBody contains the data for a firewall rule create request.
type RuleCreateRequestBody struct {
	Action string `json:"action"  url:"action"`
	Type   string `json:"type"    url:"type"`

	Comment  *string           `json:"comment,omitempty"   url:"comment,omitempty"`
	Dest     *string           `json:"dest,omitempty"      url:"dest,omitempty"`
	Digest   *string           `json:"digest,omitempty"    url:"digest,omitempty"`
	DPort    *string           `json:"dport,omitempty"     url:"dport,omitempty"`
	Enable   *types.CustomBool `json:"enable,omitempty"    url:"enable,omitempty,int"`
	ICMPType *string           `json:"icmp-type,omitempty" url:"icmp-type,omitempty"`
	IFace    *string           `json:"iface,omitempty"     url:"iface,omitempty"`
	Log      *string           `json:"log,omitempty"       url:"log,omitempty"`
	Macro    *string           `json:"macro,omitempty"     url:"macro,omitempty"`
	Proto    *string           `json:"proto,omitempty"     url:"proto,omitempty"`
	Source   *string           `json:"source,omitempty"    url:"source,omitempty"`
	SPort    *string           `json:"sport,omitempty"     url:"sport,omitempty"`

	Pos   *int    `json:"pos,omitempty"       url:"pos,omitempty"`
	Group *string `json:"group,omitempty"   url:"group,omitempty"`
}

// RuleGetResponseBody contains the body from a firewall rule get response.
type RuleGetResponseBody struct {
	Data *RuleGetResponseData `json:"data,omitempty"`
}

// RuleGetResponseData contains the data from a firewall rule get response.
type RuleGetResponseData struct {
	BaseRule

	// NOTE: This is `int` in the PVE API docs, but it's actually a string in the response.
	Pos string `json:"pos"     url:"pos"`

	Action string `json:"action"  url:"action"`
	Type   string `json:"type"    url:"type"`
}

// RuleListResponseBody contains the data from a firewall rule get response.
type RuleListResponseBody struct {
	Data []*RuleListResponseData `json:"data,omitempty"`
}

// RuleListResponseData contains the data from a firewall rule get response.
type RuleListResponseData struct {
	Pos int `json:"pos" url:"pos"`
}

// RuleUpdateRequestBody contains the data for a firewall rule update request.
type RuleUpdateRequestBody struct {
	BaseRule

	Delete *string `json:"delete,omitempty"  url:"delete,omitempty"`

	Action *string `json:"action,omitempty"  url:"action,omitempty"`
	Type   *string `json:"type,omitempty"    url:"type,omitempty"`
	Pos    *int    `json:"pos,omitempty"       url:"pos,omitempty"`
}

type BaseRule struct {
	Comment  *string           `json:"comment,omitempty"   url:"comment,omitempty"`
	Dest     *string           `json:"dest,omitempty"      url:"dest,omitempty"`
	Digest   *string           `json:"digest,omitempty"    url:"digest,omitempty"`
	DPort    *string           `json:"dport,omitempty"     url:"dport,omitempty"`
	Enable   *types.CustomBool `json:"enable,omitempty"    url:"enable,omitempty,int"`
	ICMPType *string           `json:"icmp-type,omitempty" url:"icmp-type,omitempty"`
	IFace    *string           `json:"iface,omitempty"     url:"iface,omitempty"`
	Log      *string           `json:"log,omitempty"       url:"log,omitempty"`
	Macro    *string           `json:"macro,omitempty"     url:"macro,omitempty"`
	Proto    *string           `json:"proto,omitempty"     url:"proto,omitempty"`
	Source   *string           `json:"source,omitempty"    url:"source,omitempty"`
	SPort    *string           `json:"sport,omitempty"     url:"sport,omitempty"`
}
