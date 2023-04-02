/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"io"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// VirtualEnvironmentDatastoreFileListResponseBody contains the body from a datastore content list response.
type VirtualEnvironmentDatastoreFileListResponseBody struct {
	Data []*VirtualEnvironmentDatastoreFileListResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentDatastoreFileListResponseData contains the data from a datastore content list response.
type VirtualEnvironmentDatastoreFileListResponseData struct {
	ContentType    string  `json:"content"`
	FileFormat     string  `json:"format"`
	FileSize       int     `json:"size"`
	ParentVolumeID *string `json:"parent,omitempty"`
	SpaceUsed      *int    `json:"used,omitempty"`
	VMID           *int    `json:"vmid,omitempty"`
	VolumeID       string  `json:"volid"`
}

// VirtualEnvironmentDatastoreGetStatusResponseBody contains the body from a datastore status get request.
type VirtualEnvironmentDatastoreGetStatusResponseBody struct {
	Data *VirtualEnvironmentDatastoreGetStatusResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentDatastoreGetStatusResponseBody contains the data from a datastore status get request.
type VirtualEnvironmentDatastoreGetStatusResponseData struct {
	Active         *types.CustomBool               `json:"active,omitempty"`
	AvailableBytes *int64                          `json:"avail,omitempty"`
	Content        *types.CustomCommaSeparatedList `json:"content,omitempty" url:"content,omitempty,comma"`
	Enabled        *types.CustomBool               `json:"enabled,omitempty"`
	Shared         *types.CustomBool               `json:"shared,omitempty"`
	TotalBytes     *int64                          `json:"total,omitempty"`
	Type           *string                         `json:"type,omitempty"`
	UsedBytes      *int64                          `json:"used,omitempty"`
}

// VirtualEnvironmentDatastoreListRequestBody contains the body for a datastore list request.
type VirtualEnvironmentDatastoreListRequestBody struct {
	ContentTypes types.CustomCommaSeparatedList `json:"content,omitempty" url:"content,omitempty,comma"`
	Enabled      *types.CustomBool              `json:"enabled,omitempty" url:"enabled,omitempty,int"`
	Format       *types.CustomBool              `json:"format,omitempty"  url:"format,omitempty,int"`
	ID           *string                        `json:"storage,omitempty" url:"storage,omitempty"`
	Target       *string                        `json:"target,omitempty"  url:"target,omitempty"`
}

// VirtualEnvironmentDatastoreListResponseBody contains the body from a datastore list response.
type VirtualEnvironmentDatastoreListResponseBody struct {
	Data []*VirtualEnvironmentDatastoreListResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentDatastoreListResponseData contains the data from a datastore list response.
type VirtualEnvironmentDatastoreListResponseData struct {
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

// VirtualEnvironmentDatastoreUploadRequestBody contains the body for a datastore upload request.
type VirtualEnvironmentDatastoreUploadRequestBody struct {
	ContentType string    `json:"content,omitempty"`
	DatastoreID string    `json:"storage,omitempty"`
	FileName    string    `json:"filename,omitempty"`
	FileReader  io.Reader `json:"-"`
	NodeName    string    `json:"node,omitempty"`
}

// VirtualEnvironmentDatastoreUploadResponseBody contains the body from a datastore upload response.
type VirtualEnvironmentDatastoreUploadResponseBody struct {
	UploadID *string `json:"data,omitempty"`
}
