/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package sdn

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

// QueryParams represents optionally query parameters for all SDN API calls.
type QueryParams struct {
	Pending *types.CustomBool `url:"pending,omitempty,int"`
	Running *types.CustomBool `url:"running,omitempty,int"`
}
