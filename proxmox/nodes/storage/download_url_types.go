/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

// DownloadURLResponseBody contains the body from a DownloadURL post response.
type DownloadURLResponseBody struct {
	TaskID *string `json:"data,omitempty"`
}

// DownloadURLPostRequestBody contains the body for a DownloadURL post request.
type DownloadURLPostRequestBody struct {
	Content           *string           `json:"content,omitempty"             url:"content,omitempty"`
	FileName          *string           `json:"filename,omitempty"            url:"filename,omitempty"`
	Node              *string           `json:"node,omitempty"                url:"node,omitempty"`
	Storage           *string           `json:"storage,omitempty"             url:"storage,omitempty"`
	URL               *string           `json:"url,omitempty"                 url:"url,omitempty"`
	Checksum          *string           `json:"checksum,omitempty"            url:"checksum,omitempty"`
	ChecksumAlgorithm *string           `json:"checksum-algorithm,omitempty"  url:"checksum-algorithm,omitempty"`
	Compression       *string           `json:"compression,omitempty"         url:"compression,omitempty"`
	Verify            *types.CustomBool `json:"verify-certificates,omitempty" url:"verify-certificates,omitempty,int"`
}
