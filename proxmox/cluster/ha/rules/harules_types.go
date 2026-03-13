/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package rules

import (
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// HARuleListResponseBody contains the body from a HA rule list response.
type HARuleListResponseBody struct {
	Data []*HARuleListResponseData `json:"data,omitempty"`
}

// HARuleListResponseData contains the data from a HA rule list response.
type HARuleListResponseData struct {
	Rule string `json:"rule"`
	Type string `json:"type"`
}

// HARuleGetResponseBody contains the body from a HA rule get response.
type HARuleGetResponseBody struct {
	Data *HARuleGetResponseData `json:"data,omitempty"`
}

// HARuleDataBase contains fields which are both received from and sent to the HA rule API.
type HARuleDataBase struct {
	// A SHA1 digest of the rule's configuration.
	Digest *string `json:"digest,omitempty" url:"digest,omitempty"`
	// The rule's comment, if defined.
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	// Whether the HA rule is disabled.
	Disable types.CustomBool `json:"disable,omitempty" url:"disable,int"`
	// The HA rule type: node-affinity or resource-affinity.
	Type string `json:"type" url:"type"`
	// A comma-separated list of HA resource IDs (e.g. vm:100,ct:101).
	Resources string `json:"resources" url:"resources"`
}

// HARuleNodeAffinityData contains fields specific to node-affinity rules.
type HARuleNodeAffinityData struct {
	// A comma-separated list of node names with optional priorities (e.g. node1:2,node2:1).
	Nodes *string `json:"nodes,omitempty" url:"nodes,omitempty"`
	// Whether the node affinity rule is strict (resources cannot run on other nodes).
	Strict *types.CustomBool `json:"strict,omitempty" url:"strict,omitempty,int"`
}

// HARuleResourceAffinityData contains fields specific to resource-affinity rules.
type HARuleResourceAffinityData struct {
	// Whether resources should be kept on the same node (positive) or separate nodes (negative).
	Affinity *string `json:"affinity,omitempty" url:"affinity,omitempty"`
}

// HARuleGetResponseData contains the data from a HA rule get response.
type HARuleGetResponseData struct {
	HARuleDataBase
	HARuleNodeAffinityData
	HARuleResourceAffinityData

	// The rule's identifier.
	Rule string `json:"rule"`
	// The rule's order/priority (lower = higher priority).
	Order *int64 `json:"order,omitempty"`
}

// HARuleCreateRequestBody contains the data which must be sent when creating a HA rule.
type HARuleCreateRequestBody struct {
	HARuleDataBase
	HARuleNodeAffinityData
	HARuleResourceAffinityData

	// The rule's identifier.
	Rule string `url:"rule"`
}

// HARuleUpdateRequestBody contains the data which must be sent when updating a HA rule.
type HARuleUpdateRequestBody struct {
	HARuleDataBase
	HARuleNodeAffinityData
	HARuleResourceAffinityData

	// A list of settings to delete.
	Delete []string `url:"delete,omitempty,comma"`
}
