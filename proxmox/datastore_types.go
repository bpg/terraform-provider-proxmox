/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package proxmox

import (
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// DatastoreGetResponseBody contains the body from a datastore get response.
type DatastoreGetResponseBody struct {
	Data *DatastoreGetResponseData `json:"data,omitempty"`
}

// DatastoreGetResponseData contains the data from a datastore get response.
type DatastoreGetResponseData struct {
	Content types.CustomCommaSeparatedList `json:"content,omitempty" url:"content,omitempty,comma"`
	Digest  *string                        `json:"digest,omitempty"`
	Path    *string                        `json:"path,omitempty"`
	Shared  *types.CustomBool              `json:"shared,omitempty"`
	Storage *string                        `json:"storage,omitempty"`
	Type    *string                        `json:"type,omitempty"`
}
