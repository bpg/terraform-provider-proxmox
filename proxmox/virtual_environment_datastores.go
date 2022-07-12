/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/pkg/sftp"
)

// DeleteDatastoreFile deletes a file in a datastore.
func (c *VirtualEnvironmentClient) DeleteDatastoreFile(ctx context.Context, nodeName, datastoreID, volumeID string) error {
	err := c.DoRequest(ctx, hmDELETE, fmt.Sprintf("nodes/%s/storage/%s/content/%s", url.PathEscape(nodeName), url.PathEscape(datastoreID), url.PathEscape(volumeID)), nil, nil)

	if err != nil {
		return err
	}

	return nil
}

// ListDatastoreFiles retrieves a list of the files in a datastore.
func (c *VirtualEnvironmentClient) ListDatastoreFiles(ctx context.Context, nodeName, datastoreID string) ([]*VirtualEnvironmentDatastoreFileListResponseData, error) {
	resBody := &VirtualEnvironmentDatastoreFileListResponseBody{}
	err := c.DoRequest(ctx, hmGET, fmt.Sprintf("nodes/%s/storage/%s/content", url.PathEscape(nodeName), url.PathEscape(datastoreID)), nil, resBody)

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
func (c *VirtualEnvironmentClient) ListDatastores(ctx context.Context, nodeName string, d *VirtualEnvironmentDatastoreListRequestBody) ([]*VirtualEnvironmentDatastoreListResponseData, error) {
	resBody := &VirtualEnvironmentDatastoreListResponseBody{}
	err := c.DoRequest(ctx, hmGET, fmt.Sprintf("nodes/%s/storage", url.PathEscape(nodeName)), d, resBody)

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
func (c *VirtualEnvironmentClient) UploadFileToDatastore(ctx context.Context, d *VirtualEnvironmentDatastoreUploadRequestBody) (*VirtualEnvironmentDatastoreUploadResponseBody, error) {
	switch d.ContentType {
	case "iso", "vztmpl":
		r, w := io.Pipe()

		defer r.Close()

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

		// We need to store the multipart content in a temporary file to avoid using high amounts of memory.
		// This is necessary due to Proxmox VE not supporting chunked transfers in v6.1 and earlier versions.
		tempMultipartFile, err := ioutil.TempFile("", "multipart")

		if err != nil {
			return nil, err
		}

		tempMultipartFileName := tempMultipartFile.Name()

		io.Copy(tempMultipartFile, r)

		err = tempMultipartFile.Close()

		if err != nil {
			return nil, err
		}

		defer os.Remove(tempMultipartFileName)

		// Now that the multipart data is stored in a file, we can go ahead and do a HTTP POST request.
		fileReader, err := os.Open(tempMultipartFileName)

		if err != nil {
			return nil, err
		}

		defer fileReader.Close()

		fileInfo, err := fileReader.Stat()

		if err != nil {
			return nil, err
		}

		fileSize := fileInfo.Size()

		reqBody := &VirtualEnvironmentMultiPartData{
			Boundary: m.Boundary(),
			Reader:   fileReader,
			Size:     &fileSize,
		}

		resBody := &VirtualEnvironmentDatastoreUploadResponseBody{}
		err = c.DoRequest(ctx, hmPOST, fmt.Sprintf("nodes/%s/storage/%s/upload", url.PathEscape(d.NodeName), url.PathEscape(d.DatastoreID)), reqBody, resBody)

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

		defer sshClient.Close()

		sshSession, err := sshClient.NewSession()

		if err != nil {
			return nil, err
		}

		buf, err := sshSession.CombinedOutput(
			fmt.Sprintf(`awk "/.+: %s$/,/^$/" /etc/pve/storage.cfg | grep -oP '(?<=path[ ])[^\s]+' | head -c -1`, d.DatastoreID),
		)

		if err != nil {
			sshSession.Close()

			return nil, err
		}

		sshSession.Close()

		datastorePath := strings.Trim(string(buf), "\000")

		if datastorePath == "" {
			return nil, errors.New("failed to determine the datastore path")
		}

		remoteFileDir := datastorePath

		switch d.ContentType {
		default:
			remoteFileDir += fmt.Sprintf("/%s", d.ContentType)
		}

		remoteFilePath := fmt.Sprintf("%s/%s", remoteFileDir, d.FileName)
		sftpClient, err := sftp.NewClient(sshClient)

		if err != nil {
			return nil, err
		}

		defer sftpClient.Close()

		err = sftpClient.MkdirAll(remoteFileDir)

		if err != nil {
			return nil, err
		}

		remoteFile, err := sftpClient.Create(remoteFilePath)

		if err != nil {
			return nil, err
		}

		defer remoteFile.Close()

		_, err = remoteFile.ReadFrom(d.FileReader)

		if err != nil {
			return nil, err
		}

		return &VirtualEnvironmentDatastoreUploadResponseBody{}, nil
	}
}
