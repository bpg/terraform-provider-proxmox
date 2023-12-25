/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodestorage

// DatastoreFileListResponseBody contains the body from a datastore content list response.
type DatastoreFileListResponseBody struct {
	Data []*DatastoreFileListResponseData `json:"data,omitempty"`
}

// DatastoreFileListResponseData contains the data from a datastore content list response.
type DatastoreFileListResponseData struct {
	ContentType    string  `json:"content"`
	FileFormat     string  `json:"format"`
	FileSize       int64   `json:"size"`
	ParentVolumeID *string `json:"parent,omitempty"`
	SpaceUsed      *int    `json:"used,omitempty"`
	VMID           *int    `json:"vmid,omitempty"`
	VolumeID       string  `json:"volid"`
}

// DatastoreFileGetRequestData contains the body from a datastore content get request.
type DatastoreFileGetRequestData struct {
	Node     string `json:"node"`
	VolumeID string `json:"volume"`
}

// DatastoreFileGetResponseBody contains the body from a datastore content get response.
type DatastoreFileGetResponseBody struct {
	Data *DatastoreFileGetResponseData `json:"data,omitempty"`
}

// DatastoreFileGetResponseData contains the data from a datastore content get response.
type DatastoreFileGetResponseData struct {
	Path       *string `json:"path"`
	FileFormat *string `json:"format"`
	FileSize   *int64  `json:"size"`
	SpaceUsed  *int64  `json:"used,omitempty"`
}
