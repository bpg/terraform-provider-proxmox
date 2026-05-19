/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package pool

import (
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// CreateRequestBody contains the body for creating a new Ceph pool.
type CreateRequestBody struct {
	Name string `url:"name"`

	AddStorages     *types.CustomBool `url:"add_storages,omitempty,int"`
	Application     *string           `url:"application,omitempty"`
	CrushRule       *string           `url:"crush_rule,omitempty"`
	ErasureCoding   *string           `url:"erasure-coding,omitempty"`
	MinSize         *int64            `url:"min_size,omitempty"`
	PGAutoscaleMode *string           `url:"pg_autoscale_mode,omitempty"`
	PGNum           *int64            `url:"pg_num,omitempty"`
	PGNumMin        *int64            `url:"pg_num_min,omitempty"`
	Size            *int64            `url:"size,omitempty"`
	TargetSize      *string           `url:"target_size,omitempty"`
	TargetSizeRatio *float64          `url:"target_size_ratio,omitempty"`
}

// UpdateRequestBody contains the body for updating an existing Ceph pool.
type UpdateRequestBody struct {
	Application     *string  `url:"application,omitempty"`
	CrushRule       *string  `url:"crush_rule,omitempty"`
	MinSize         *int64   `url:"min_size,omitempty"`
	PGAutoscaleMode *string  `url:"pg_autoscale_mode,omitempty"`
	PGNum           *int64   `url:"pg_num,omitempty"`
	PGNumMin        *int64   `url:"pg_num_min,omitempty"`
	Size            *int64   `url:"size,omitempty"`
	TargetSize      *string  `url:"target_size,omitempty"`
	TargetSizeRatio *float64 `url:"target_size_ratio,omitempty"`
}

// DeleteRequestParams contains the query parameters for destroying a Ceph pool.
type DeleteRequestParams struct {
	Force           *types.CustomBool `url:"force,omitempty,int"`
	RemoveECProfile *types.CustomBool `url:"remove_ecprofile,omitempty,int"`
	RemoveStorages  *types.CustomBool `url:"remove_storages,omitempty,int"`
}

// ListResponseBody wraps the list-pools response.
type ListResponseBody struct {
	Data []*ListResponseData `json:"data,omitempty"`
}

// ListResponseData describes a single pool entry returned by the list endpoint.
// The list endpoint returns the full settable settings, so it doubles as the read-back source.
type ListResponseData struct {
	PoolName        string `json:"pool_name"`
	Type            string `json:"type,omitempty"`
	Size            int64  `json:"size"`
	MinSize         int64  `json:"min_size"`
	PGNum           int64  `json:"pg_num"`
	PGNumMin        *int64 `json:"pg_num_min,omitempty"`
	PGNumFinal      *int64 `json:"pg_num_final,omitempty"`
	PGAutoscaleMode string `json:"pg_autoscale_mode,omitempty"`
	// PVE returns crush_rule as a JSON string ("1") in current Squid releases despite the API
	// spec declaring it as integer; we don't need the numeric id (the human-readable name is
	// in crush_rule_name) so the field is intentionally omitted to avoid the type drift.
	CrushRuleName   string   `json:"crush_rule_name"`
	TargetSize      *int64   `json:"target_size,omitempty"`
	TargetSizeRatio *float64 `json:"target_size_ratio,omitempty"`
	// Application is derived from the application_metadata map keys. The list
	// endpoint returns `application_metadata: { "rbd": {} }`-shaped objects;
	// we capture the raw map and resolve the application name in fromAPI.
	ApplicationMetadata map[string]any `json:"application_metadata,omitempty"`
}

// CreateResponseBody wraps the create response (a UPID).
type CreateResponseBody struct {
	Data *string `json:"data,omitempty"`
}

// UpdateResponseBody wraps the update response (a UPID).
type UpdateResponseBody struct {
	Data *string `json:"data,omitempty"`
}

// DeleteResponseBody wraps the delete response (a UPID).
type DeleteResponseBody struct {
	Data *string `json:"data,omitempty"`
}
