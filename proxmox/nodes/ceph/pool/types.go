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

// StatusResponseBody wraps the per-pool /status response.
type StatusResponseBody struct {
	Data *StatusResponseData `json:"data,omitempty"`
}

// StatusResponseData describes the per-pool /status?verbose=1 response. The endpoint
// returns the full settable settings plus a few runtime flags we ignore. crush_rule
// arrives as the rule name (no separate _name field), and application is reported
// via application_list when verbose=1 is set.
type StatusResponseData struct {
	Name            string   `json:"name"`
	ApplicationList []string `json:"application_list,omitempty"`
	CrushRule       string   `json:"crush_rule,omitempty"`
	Size            int64    `json:"size"`
	MinSize         int64    `json:"min_size"`
	PGNum           int64    `json:"pg_num"`
	PGNumMin        *int64   `json:"pg_num_min,omitempty"`
	PGAutoscaleMode string   `json:"pg_autoscale_mode,omitempty"`
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
