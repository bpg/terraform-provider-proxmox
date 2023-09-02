/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"sort"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// DeleteDatastoreFile deletes a file in a datastore.
func (c *Client) DeleteDatastoreFile(
	ctx context.Context,
	datastoreID, volumeID string,
) error {
	err := c.DoRequest(
		ctx,
		http.MethodDelete,
		c.ExpandPath(
			fmt.Sprintf(
				"storage/%s/content/%s",
				url.PathEscape(datastoreID),
				url.PathEscape(volumeID),
			),
		),
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error deleting file %s from datastore %s: %w", volumeID, datastoreID, err)
	}

	return nil
}

// GetDatastoreStatus gets status information for a given datastore.
func (c *Client) GetDatastoreStatus(
	ctx context.Context,
	datastoreID string,
) (*DatastoreGetStatusResponseData, error) {
	resBody := &DatastoreGetStatusResponseBody{}

	err := c.DoRequest(
		ctx,
		http.MethodGet,
		c.ExpandPath(
			fmt.Sprintf(
				"storage/%s/status",
				url.PathEscape(datastoreID),
			),
		),
		nil,
		resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving status for datastore %s: %w", datastoreID, err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// ListDatastoreFiles retrieves a list of the files in a datastore.
func (c *Client) ListDatastoreFiles(
	ctx context.Context,
	datastoreID string,
) ([]*DatastoreFileListResponseData, error) {
	resBody := &DatastoreFileListResponseBody{}

	err := c.DoRequest(
		ctx,
		http.MethodGet,
		c.ExpandPath(
			fmt.Sprintf(
				"storage/%s/content",
				url.PathEscape(datastoreID),
			),
		),
		nil,
		resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving files from datastore %s: %w", datastoreID, err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].VolumeID < resBody.Data[j].VolumeID
	})

	return resBody.Data, nil
}

// ListDatastores retrieves a list of nodes.
func (c *Client) ListDatastores(
	ctx context.Context,
	d *DatastoreListRequestBody,
) ([]*DatastoreListResponseData, error) {
	resBody := &DatastoreListResponseBody{}

	err := c.DoRequest(
		ctx,
		http.MethodGet,
		c.ExpandPath("storage"),
		d,
		resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving datastores: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].ID < resBody.Data[j].ID
	})

	return resBody.Data, nil
}

// APIUpload uploads a file to a datastore using the Proxmox API.
func (c *Client) APIUpload(
	ctx context.Context,
	datastoreID string,
	d *api.FileUploadRequest,
	uploadTimeout int,
) (*DatastoreUploadResponseBody, error) {
	tflog.Debug(ctx, "uploading file to datastore using PVE API", map[string]interface{}{
		"file_name":    d.FileName,
		"content_type": d.ContentType,
	})

	r, w := io.Pipe()

	defer func(r *io.PipeReader) {
		err := r.Close()
		if err != nil {
			tflog.Error(ctx, "failed to close pipe reader", map[string]interface{}{
				"error": err,
			})
		}
	}(r)

	m := multipart.NewWriter(w)

	go func() {
		defer func(w *io.PipeWriter) {
			err := w.Close()
			if err != nil {
				tflog.Error(ctx, "failed to close pipe writer", map[string]interface{}{
					"error": err,
				})
			}
		}(w)
		defer func(m *multipart.Writer) {
			err := m.Close()
			if err != nil {
				tflog.Error(ctx, "failed to close multipart writer", map[string]interface{}{
					"error": err,
				})
			}
		}(m)

		err := m.WriteField("content", d.ContentType)
		if err != nil {
			tflog.Error(ctx, "failed to write 'content' field", map[string]interface{}{
				"error": err,
			})

			return
		}

		part, err := m.CreateFormFile("filename", d.FileName)
		if err != nil {
			return
		}

		_, err = io.Copy(part, d.File)

		if err != nil {
			return
		}
	}()

	// We need to store the multipart content in a temporary file to avoid using high amounts of memory.
	// This is necessary due to Proxmox VE not supporting chunked transfers in v6.1 and earlier versions.
	tempMultipartFile, err := os.CreateTemp("", "multipart")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}

	tempMultipartFileName := tempMultipartFile.Name()

	_, err = io.Copy(tempMultipartFile, r)
	if err != nil {
		return nil, fmt.Errorf("failed to copy multipart data to temporary file: %w", err)
	}

	err = tempMultipartFile.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close temporary file: %w", err)
	}

	defer func(name string) {
		e := os.Remove(name)
		if e != nil {
			tflog.Error(ctx, "failed to remove temporary file", map[string]interface{}{
				"error": e,
			})
		}
	}(tempMultipartFileName)

	// Now that the multipart data is stored in a file, we can go ahead and do an HTTP POST request.
	fileReader, err := os.Open(tempMultipartFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open temporary file: %w", err)
	}

	defer func(fileReader *os.File) {
		e := fileReader.Close()
		if e != nil {
			tflog.Error(ctx, "failed to close file reader", map[string]interface{}{
				"error": e,
			})
		}
	}(fileReader)

	fileInfo, err := fileReader.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	fileSize := fileInfo.Size()

	reqBody := &api.MultiPartData{
		Boundary: m.Boundary(),
		Reader:   fileReader,
		Size:     &fileSize,
	}

	resBody := &DatastoreUploadResponseBody{}
	err = c.DoRequest(
		ctx,
		http.MethodPost,
		c.ExpandPath(
			fmt.Sprintf(
				"storage/%s/upload",
				url.PathEscape(datastoreID),
			),
		),
		reqBody,
		resBody,
	)

	if err != nil {
		return nil, fmt.Errorf("error uploading file to datastore %s: %w", datastoreID, err)
	}

	if resBody.UploadID == nil {
		return nil, fmt.Errorf("error uploading file to datastore %s: no uploadID", datastoreID)
	}

	err = c.Tasks().WaitForTask(ctx, *resBody.UploadID, uploadTimeout, 5)
	if err != nil {
		return nil, fmt.Errorf("error uploading file to datastore %s: failed waiting for upload - %w", datastoreID, err)
	}

	return resBody, nil
}
