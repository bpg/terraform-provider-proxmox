/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cluster

import (
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type userTagAccess struct {
	UserAllowList *[]string `json:"user-allow-list,omitempty"`
	UserAllow     *string   `json:"user-allow,omitempty"`
}
type tagStyle struct {
	Shape         *string           `json:"shape,omitempty"`
	CaseSensitive *types.CustomBool `json:"case-sensitive,omitempty"`
	Ordering      *string           `json:"ordering,omitempty"`
	ColorMap      *string           `json:"color-map,omitempty"`
}
type crs struct {
	HaRebalanceOnStart *types.CustomBool `json:"ha-rebalance-on-start,omitempty"`
	HA                 *string           `json:"ha,omitempty"`
}
type notify struct {
	HAFencingMode        *string `json:"fencing,omitempty"`
	HAFencingTarget      *string `json:"target-fencing,omitempty"`
	PackageUpdates       *string `json:"package-updates,omitempty"`
	PackageUpdatesTarget *string `json:"target-package-updates,omitempty"`
	Replication          *string `json:"replication,omitempty"`
	ReplicationTarget    *string `json:"target-replication,omitempty"`
}
type migration struct {
	Network *string `json:"network,omitempty"`
	Type    *string `json:"type,omitempty"`
}

type haSettings struct {
	ShutdownPolicy *string `json:"shutdown_policy,omitempty"`
}
type nextID struct {
	Upper *types.CustomInt64 `json:"upper,omitempty"`
	Lower *types.CustomInt64 `json:"lower,omitempty"`
}
type webauthn struct {
	ID     *string `json:"id,omitempty"`
	Origin *string `json:"origin,omitempty"`
	RP     *string `json:"rp,omitempty"`
}
type optionsBaseData struct {
	BandwidthLimit *string `json:"bwlimit,omitempty"     url:"bwlimit,omitempty"`
	EmailFrom      *string `json:"email_from,omitempty"  url:"email_from,omitempty"`
	Description    *string `json:"description,omitempty" url:"description,omitempty"`
	Console        *string `json:"console,omitempty"     url:"console,omitempty"`
	HTTPProxy      *string `json:"http_proxy,omitempty"  url:"http_proxy,omitempty"`
	MacPrefix      *string `json:"mac_prefix,omitempty"  url:"mac_prefix,omitempty"`
	Keyboard       *string `json:"keyboard,omitempty"    url:"keyboard,omitempty"`
	Language       *string `json:"language,omitempty"    url:"language,omitempty"`
}

// OptionsResponseBody contains the body from a cluster options response.
type OptionsResponseBody struct {
	Data *OptionsResponseData `json:"data,omitempty"`
}

// OptionsResponseData contains the data from a cluster options response.
type OptionsResponseData struct {
	optionsBaseData

	MaxWorkers                *types.CustomInt64 `json:"max_workers,omitempty"`
	ClusterResourceScheduling *crs               `json:"crs,omitempty"`
	HASettings                *haSettings        `json:"ha,omitempty"`
	TagStyle                  *tagStyle          `json:"tag-style,omitempty"`
	Migration                 *migration         `json:"migration,omitempty"`
	Webauthn                  *webauthn          `json:"webauthn,omitempty"`
	NextID                    *nextID            `json:"next-id,omitempty"`
	Notify                    *notify            `json:"notify,omitempty"`
	UserTagAccess             *userTagAccess     `json:"user-tag-access,omitempty"`
	RegisteredTags            *[]string          `json:"registered-tags,omitempty"`
}

// OptionsRequestData contains the body for cluster options request.
type OptionsRequestData struct {
	optionsBaseData

	MaxWorkers                *int64  `json:"max_workers,omitempty"     url:"max_workers,omitempty"`
	Delete                    *string `json:"delete,omitempty"          url:"delete,omitempty"`
	ClusterResourceScheduling *string `json:"crs,omitempty"             url:"crs,omitempty"`
	HASettings                *string `json:"ha,omitempty"              url:"ha,omitempty"`
	TagStyle                  *string `json:"tag-style,omitempty"       url:"tag-style,omitempty"`
	Migration                 *string `json:"migration,omitempty"       url:"migration,omitempty"`
	Webauthn                  *string `json:"webauthn,omitempty"        url:"webauthn,omitempty"`
	NextID                    *string `json:"next-id,omitempty"         url:"next-id,omitempty"`
	Notify                    *string `json:"notify,omitempty"          url:"notify,omitempty"`
	UserTagAccess             *string `json:"user-tag-access,omitempty" url:"user-tag-access,omitempty"`
	RegisteredTags            *string `json:"registered-tags,omitempty" url:"registered-tags,omitempty"`
}
