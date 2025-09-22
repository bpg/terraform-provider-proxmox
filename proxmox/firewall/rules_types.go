/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

// RuleCreateRequestBody contains the data for a firewall rule create request.
type RuleCreateRequestBody struct {
	BaseRule

	Action string `json:"action" url:"action"`
	Type   string `json:"type"   url:"type"`

	Group *string `json:"group,omitempty" url:"group,omitempty"`
}

// RuleGetResponseBody contains the body from a firewall rule get response.
type RuleGetResponseBody struct {
	Data *RuleGetResponseData `json:"data,omitempty"`
}

// RuleGetResponseData contains the data from a firewall rule get response.
type RuleGetResponseData struct {
	BaseRule

	// NOTE: This is `int` in the PVE API docs, but it's actually a string in the response.
	Pos    string `json:"pos"    url:"pos"`
	Action string `json:"action" url:"action"`
	Type   string `json:"type"   url:"type"`
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

	Pos    *int    `json:"pos,omitempty"    url:"pos,omitempty"`
	Action *string `json:"action,omitempty" url:"action,omitempty"`
	Type   *string `json:"type,omitempty"   url:"type,omitempty"`
	Group  *string `json:"group,omitempty"  url:"group,omitempty"`
	Delete *string `json:"delete,omitempty" url:"delete,omitempty"`
}
