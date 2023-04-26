/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package proxmox

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/crypto/ssh"

	"github.com/pkg/sftp"
)

// GetDatastore retrieves information about a datastore.
/*
Using undocumented API endpoints is not recommended, but sometimes it's the only way to get things done.
$ pvesh get /storage/local
┌─────────┬───────────────────────────────────────────┐
│ key     │ value                                     │
╞═════════╪═══════════════════════════════════════════╡
│ content │ images,vztmpl,iso,backup,snippets,rootdir │
├─────────┼───────────────────────────────────────────┤
│ digest  │ 5b65ede80f34631d6039e6922845cfa4abc956be  │
├─────────┼───────────────────────────────────────────┤
│ path    │ /var/lib/vz                               │
├─────────┼───────────────────────────────────────────┤
│ shared  │ 0                                         │
├─────────┼───────────────────────────────────────────┤
│ storage │ local                                     │
├─────────┼───────────────────────────────────────────┤
│ type    │ dir                                       │
└─────────┴───────────────────────────────────────────┘
*/
func (c *VirtualEnvironmentClient) GetDatastore(
	ctx context.Context,
	datastoreID string,
) (*DatastoreGetResponseData, error) {
	resBody := &DatastoreGetResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("storage/%s", url.PathEscape(datastoreID)),
		nil,
		resBody,
	)
	if err != nil {
		return nil, err
	}

	return resBody.Data, nil
}

// DeleteDatastoreFile deletes a file in a datastore.
func (c *VirtualEnvironmentClient) DeleteDatastoreFile(
	ctx context.Context,
	nodeName, datastoreID, volumeID string,
) error {
	err := c.DoRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf(
			"nodes/%s/storage/%s/content/%s",
			url.PathEscape(nodeName),
			url.PathEscape(datastoreID),
			url.PathEscape(volumeID),
		),
		nil,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetDatastoreStatus gets status information for a given datastore.
func (c *VirtualEnvironmentClient) GetDatastoreStatus(
	ctx context.Context,
	nodeName, datastoreID string,
) (*DatastoreGetStatusResponseData, error) {
	resBody := &DatastoreGetStatusResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"nodes/%s/storage/%s/status",
			url.PathEscape(nodeName),
			url.PathEscape(datastoreID),
		),
		nil,
		resBody,
	)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// ListDatastoreFiles retrieves a list of the files in a datastore.
func (c *VirtualEnvironmentClient) ListDatastoreFiles(
	ctx context.Context,
	nodeName, datastoreID string,
) ([]*DatastoreFileListResponseData, error) {
	resBody := &DatastoreFileListResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"nodes/%s/storage/%s/content",
			url.PathEscape(nodeName),
			url.PathEscape(datastoreID),
		),
		nil,
		resBody,
	)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].VolumeID < resBody.Data[j].VolumeID
	})

	return resBody.Data, nil
}

// ListDatastores retrieves a list of nodes.
func (c *VirtualEnvironmentClient) ListDatastores(
	ctx context.Context,
	nodeName string,
	d *DatastoreListRequestBody,
) ([]*DatastoreListResponseData, error) {
	resBody := &DatastoreListResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("nodes/%s/storage", url.PathEscape(nodeName)),
		d,
		resBody,
	)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].ID < resBody.Data[j].ID
	})

	return resBody.Data, nil
}

// UploadFileToDatastore uploads a file to a datastore.
func (c *VirtualEnvironmentClient) UploadFileToDatastore(
	ctx context.Context,
	d *DatastoreUploadRequestBody,
) (*DatastoreUploadResponseBody, error) {
	switch d.ContentType {
	case "iso", "vztmpl":
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

			_, err = io.Copy(part, d.FileReader)

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
			err := os.Remove(name)
			if err != nil {
				tflog.Error(ctx, "failed to remove temporary file", map[string]interface{}{
					"error": err,
				})
			}
		}(tempMultipartFileName)

		// Now that the multipart data is stored in a file, we can go ahead and do a HTTP POST request.
		fileReader, err := os.Open(tempMultipartFileName)
		if err != nil {
			return nil, fmt.Errorf("failed to open temporary file: %w", err)
		}

		defer func(fileReader *os.File) {
			err := fileReader.Close()
			if err != nil {
				tflog.Error(ctx, "failed to close file reader", map[string]interface{}{
					"error": err,
				})
			}
		}(fileReader)

		fileInfo, err := fileReader.Stat()
		if err != nil {
			return nil, fmt.Errorf("failed to get file info: %w", err)
		}

		fileSize := fileInfo.Size()

		reqBody := &VirtualEnvironmentMultiPartData{
			Boundary: m.Boundary(),
			Reader:   fileReader,
			Size:     &fileSize,
		}

		resBody := &DatastoreUploadResponseBody{}
		err = c.DoRequest(
			ctx,
			http.MethodPost,
			fmt.Sprintf(
				"nodes/%s/storage/%s/upload",
				url.PathEscape(d.NodeName),
				url.PathEscape(d.DatastoreID),
			),
			reqBody,
			resBody,
		)

		if err != nil {
			return nil, err
		}

		return resBody, nil
	default:
		// We need to upload all other files using SFTP due to API limitations.
		// Hopefully, this will not be required in future releases of Proxmox VE.
		sshClient, err := c.OpenNodeShell(ctx, d.NodeName)
		if err != nil {
			return nil, err
		}

		defer func(sshClient *ssh.Client) {
			err := sshClient.Close()
			if err != nil {
				tflog.Error(ctx, "failed to close SSH client", map[string]interface{}{
					"error": err,
				})
			}
		}(sshClient)

		datastore, err := c.GetDatastore(ctx, d.DatastoreID)
		if err != nil {
			return nil, fmt.Errorf("failed to get datastore: %w", err)
		}
		if datastore.Path == nil || *datastore.Path == "" {
			return nil, errors.New("failed to determine the datastore path")
		}

		remoteFileDir := *datastore.Path
		if d.ContentType != "" {
			remoteFileDir = filepath.Join(remoteFileDir, d.ContentType)
		}

		remoteFilePath := filepath.Join(remoteFileDir, d.FileName)
		sftpClient, err := sftp.NewClient(sshClient)
		if err != nil {
			return nil, fmt.Errorf("failed to create SFTP client: %w", err)
		}

		defer func(sftpClient *sftp.Client) {
			err := sftpClient.Close()
			if err != nil {
				tflog.Error(ctx, "failed to close SFTP client", map[string]interface{}{
					"error": err,
				})
			}
		}(sftpClient)

		err = sftpClient.MkdirAll(remoteFileDir)
		if err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", remoteFileDir, err)
		}

		remoteFile, err := sftpClient.Create(remoteFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create file %s: %w", remoteFilePath, err)
		}

		defer func(remoteFile *sftp.File) {
			err := remoteFile.Close()
			if err != nil {
				tflog.Error(ctx, "failed to close remote file", map[string]interface{}{
					"error": err,
				})
			}
		}(remoteFile)

		_, err = remoteFile.ReadFrom(d.FileReader)
		if err != nil {
			return nil, fmt.Errorf("failed to upload file %s: %w", remoteFilePath, err)
		}

		return &DatastoreUploadResponseBody{}, nil
	}
}
