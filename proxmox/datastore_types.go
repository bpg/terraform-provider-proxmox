/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package proxmox

import (
	"os"

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

// DatastoreFileListResponseBody contains the body from a datastore content list response.
type DatastoreFileListResponseBody struct {
	Data []*DatastoreFileListResponseData `json:"data,omitempty"`
}

// DatastoreFileListResponseData contains the data from a datastore content list response.
type DatastoreFileListResponseData struct {
	ContentType    string  `json:"content"`
	FileFormat     string  `json:"format"`
	FileSize       int     `json:"size"`
	ParentVolumeID *string `json:"parent,omitempty"`
	SpaceUsed      *int    `json:"used,omitempty"`
	VMID           *int    `json:"vmid,omitempty"`
	VolumeID       string  `json:"volid"`
}

// DatastoreGetStatusResponseBody contains the body from a datastore status get request.
type DatastoreGetStatusResponseBody struct {
	Data *DatastoreGetStatusResponseData `json:"data,omitempty"`
}

// DatastoreGetStatusResponseData contains the data from a datastore status get request.
type DatastoreGetStatusResponseData struct {
	Active         *types.CustomBool               `json:"active,omitempty"`
	AvailableBytes *int64                          `json:"avail,omitempty"`
	Content        *types.CustomCommaSeparatedList `json:"content,omitempty" url:"content,omitempty,comma"`
	Enabled        *types.CustomBool               `json:"enabled,omitempty"`
	Shared         *types.CustomBool               `json:"shared,omitempty"`
	TotalBytes     *int64                          `json:"total,omitempty"`
	Type           *string                         `json:"type,omitempty"`
	UsedBytes      *int64                          `json:"used,omitempty"`
}

// DatastoreListRequestBody contains the body for a datastore list request.
type DatastoreListRequestBody struct {
	ContentTypes types.CustomCommaSeparatedList `json:"content,omitempty" url:"content,omitempty,comma"`
	Enabled      *types.CustomBool              `json:"enabled,omitempty" url:"enabled,omitempty,int"`
	Format       *types.CustomBool              `json:"format,omitempty"  url:"format,omitempty,int"`
	ID           *string                        `json:"storage,omitempty" url:"storage,omitempty"`
	Target       *string                        `json:"target,omitempty"  url:"target,omitempty"`
}

// DatastoreListResponseBody contains the body from a datastore list response.
type DatastoreListResponseBody struct {
	Data []*DatastoreListResponseData `json:"data,omitempty"`
}

// DatastoreListResponseData contains the data from a datastore list response.
type DatastoreListResponseData struct {
	Active              *types.CustomBool               `json:"active,omitempty"`
	ContentTypes        *types.CustomCommaSeparatedList `json:"content,omitempty"`
	Enabled             *types.CustomBool               `json:"enabled,omitempty"`
	ID                  string                          `json:"storage,omitempty"`
	Shared              *types.CustomBool               `json:"shared,omitempty"`
	SpaceAvailable      *int                            `json:"avail,omitempty"`
	SpaceTotal          *int                            `json:"total,omitempty"`
	SpaceUsed           *int                            `json:"used,omitempty"`
	SpaceUsedPercentage *float64                        `json:"used_fraction,omitempty"`
	Type                string                          `json:"type,omitempty"`
}

// DatastoreUploadRequestBody contains the body for a datastore upload request.
type DatastoreUploadRequestBody struct {
	ContentType string   `json:"content,omitempty"`
	DatastoreID string   `json:"storage,omitempty"`
	FileName    string   `json:"filename,omitempty"`
	NodeName    string   `json:"node,omitempty"`
	File        *os.File `json:"-"`
}

// DatastoreUploadResponseBody contains the body from a datastore upload response.
type DatastoreUploadResponseBody struct {
	UploadID *string `json:"data,omitempty"`
}
