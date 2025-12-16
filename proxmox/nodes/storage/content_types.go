/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import "slices"

// Content type constants for Proxmox VE datastores.
const (
	ContentTypeBackup   = "backup"   // VM backups
	ContentTypeImages   = "images"   // VM disk images
	ContentTypeImport   = "import"   // VM disk images for import
	ContentTypeISO      = "iso"      // ISO images
	ContentTypeRootDir  = "rootdir"  // Container root directories
	ContentTypeSnippets = "snippets" // Snippets (cloud-init, hook scripts, etc.)
	ContentTypeVZTmpl   = "vztmpl"   // Container templates
)

// validContentTypes returns all valid content type values.
func validContentTypes() []string {
	return []string{
		ContentTypeBackup,
		ContentTypeImages,
		ContentTypeImport,
		ContentTypeISO,
		ContentTypeRootDir,
		ContentTypeSnippets,
		ContentTypeVZTmpl,
	}
}

// IsValidContentType checks if the given content type is valid.
func IsValidContentType(contentType string) bool {
	return slices.Contains(validContentTypes(), contentType)
}

// DatastoreFileListRequestBody contains the request parameters for listing datastore files.
type DatastoreFileListRequestBody struct {
	ContentType *string `json:"content,omitempty" url:"content,omitempty"`
}

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
	Node     string `json:"node,omitempty"   url:"node,omitempty"`
	VolumeID string `json:"volume,omitempty" url:"volume,omitempty"`
}

// DatastoreFileGetResponseBody contains the body from a datastore content get response.
type DatastoreFileGetResponseBody struct {
	Data *DatastoreFileGetResponseData `json:"data,omitempty" url:"data,omitempty"`
}

// DatastoreFileGetResponseData contains the data from a datastore content get response.
type DatastoreFileGetResponseData struct {
	Path       *string `json:"path"           url:"path,omitempty"`
	FileFormat *string `json:"format"         url:"format,omitempty"`
	FileSize   *int64  `json:"size"           url:"size,omitempty"`
	SpaceUsed  *int64  `json:"used,omitempty" url:"used,omitempty"`
}
