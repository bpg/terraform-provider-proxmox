/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"sort"
)

// VirtualEnvironmentDatastoreListRequestBody contains the body for a datastore list request.
type VirtualEnvironmentDatastoreListRequestBody struct {
	ContentTypes CustomCommaSeparatedList `json:"content,omitempty" url:"content,omitempty,comma"`
	Enabled      *CustomBool              `json:"enabled,omitempty" url:"enabled,omitempty,int"`
	Format       *CustomBool              `json:"format,omitempty" url:"format,omitempty,int"`
	ID           *string                  `json:"storage,omitempty" url:"storage,omitempty"`
	Target       *string                  `json:"target,omitempty" url:"target,omitempty"`
}

// VirtualEnvironmentDatastoreListResponseBody contains the body from a datastore list response.
type VirtualEnvironmentDatastoreListResponseBody struct {
	Data []*VirtualEnvironmentDatastoreListResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentDatastoreListResponseData contains the data from a datastore list response.
type VirtualEnvironmentDatastoreListResponseData struct {
	Active              *CustomBool               `json:"active,omitempty"`
	ContentTypes        *CustomCommaSeparatedList `json:"content,omitempty"`
	Enabled             *CustomBool               `json:"enabled,omitempty"`
	ID                  string                    `json:"storage,omitempty"`
	Shared              *CustomBool               `json:"shared,omitempty"`
	SpaceAvailable      *int                      `json:"avail,omitempty"`
	SpaceTotal          *int                      `json:"total,omitempty"`
	SpaceUsed           *int                      `json:"used,omitempty"`
	SpaceUsedPercentage *float64                  `json:"used_fraction,omitempty"`
	Type                string                    `json:"type,omitempty"`
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

// ListDatastores retrieves a list of nodes.
func (c *VirtualEnvironmentClient) ListDatastores(nodeName string, d *VirtualEnvironmentDatastoreListRequestBody) ([]*VirtualEnvironmentDatastoreListResponseData, error) {
	resBody := &VirtualEnvironmentDatastoreListResponseBody{}
	err := c.DoRequest(hmGET, fmt.Sprintf("nodes/%s/storage", nodeName), d, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].ID < resBody.Data[j].ID
	})

	return resBody.Data, nil
}

// UploadFileToDatastore uploads a file to a datastore.
func (c *VirtualEnvironmentClient) UploadFileToDatastore(d *VirtualEnvironmentDatastoreUploadRequestBody) (*VirtualEnvironmentDatastoreUploadResponseBody, error) {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)

	go func() {
		defer w.Close()
		defer m.Close()

		m.WriteField("content", d.ContentType)

		part, err := m.CreateFormFile("filename", d.FileName)

		if err != nil {
			return
		}

		_, err = io.Copy(part, d.FileReader)

		if err != nil {
			return
		}
	}()

	// Due to Proxmox VE not supporting chunked transfers, we sadly need to load the file into memory.
	// This is not optimal for large files but there's no alternative right now.
	workaroundReader := new(bytes.Buffer)
	workaroundReader.ReadFrom(r)

	reqBody := &VirtualEnvironmentMultiPartData{
		Boundary: m.Boundary(),
		Reader:   workaroundReader,
	}

	resBody := &VirtualEnvironmentDatastoreUploadResponseBody{}
	err := c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/storage/%s/upload", d.NodeName, d.DatastoreID), reqBody, resBody)

	if err != nil {
		return nil, err
	}

	return resBody, nil
}
