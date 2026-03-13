/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package replications

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

/*
Replication used to represent a Replication in the API.

Based on docs:
  - https://pve.proxmox.com/pve-docs/api-viewer/index.html#/cluster/replication
  - https://pve.proxmox.com/pve-docs/api-viewer/index.html#/cluster/replication/{id}
*/

type Replication struct {
	ID        string            `json:"id"                   url:"id"`
	Comment   *string           `json:"comment,omitempty"    url:"comment,omitempty"`
	Disable   *types.CustomBool `json:"disable,omitempty"    url:"disable,omitempty,int"`
	Rate      *float64          `json:"rate,omitempty"       url:"rate,omitempty"`
	RemoveJob *string           `json:"remove_job,omitempty" url:"remove_job,omitempty"`
	Schedule  *string           `json:"schedule,omitempty"   url:"schedule,omitempty"`
	Source    *string           `json:"source,omitempty"     url:"source,omitempty"`
}

type ReplicationData struct {
	Replication

	Target string `json:"target" url:"target"`
	Type   string `json:"type"   url:"type"`
	Guest  int64  `json:"guest"  url:"guest"`
	JobNum int64  `json:"jobnum" url:"jobnum"`
}

type ReplicationCreate struct {
	Replication

	Target string `json:"target" url:"target"`
	Type   string `json:"type"   url:"type"`
}

type ReplicationUpdate struct {
	Replication

	Delete []string `url:"delete,omitempty,comma"`
}

type ReplicationDelete struct {
}

type replicationResponse struct {
	Data *ReplicationData `json:"data"`
}

type replicationsResponse struct {
	Data *[]ReplicationData `json:"data"`
}
