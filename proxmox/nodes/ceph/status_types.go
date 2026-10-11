/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ceph

import (
	clusterceph "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/ceph"
)

// StatusResponseBody wraps the node-scoped Ceph status response. The payload
// shape is identical to the cluster endpoint, so we reuse the typed payload
// defined alongside the cluster Ceph client.
type StatusResponseBody struct {
	Data *clusterceph.StatusResponseData `json:"data,omitempty"`
}
