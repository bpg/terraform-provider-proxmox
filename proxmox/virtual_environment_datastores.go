/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/url"
	"os"
	"sort"
)

// DeleteDatastoreFile deletes a file in a datastore.
func (c *VirtualEnvironmentClient) DeleteDatastoreFile(nodeName, datastoreID, volumeID string) error {
	err := c.DoRequest(hmDELETE, fmt.Sprintf("nodes/%s/storage/%s/content/%s", url.PathEscape(nodeName), url.PathEscape(datastoreID), url.PathEscape(volumeID)), nil, nil)

	if err != nil {
		return err
	}

	return nil
}

// ListDatastoreFiles retrieves a list of the files in a datastore.
func (c *VirtualEnvironmentClient) ListDatastoreFiles(nodeName, datastoreID string) ([]*VirtualEnvironmentDatastoreFileListResponseData, error) {
	resBody := &VirtualEnvironmentDatastoreFileListResponseBody{}
	err := c.DoRequest(hmGET, fmt.Sprintf("nodes/%s/storage/%s/content", url.PathEscape(nodeName), url.PathEscape(datastoreID)), nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].VolumeID < resBody.Data[j].VolumeID
	})

	return resBody.Data, nil
}

// ListDatastores retrieves a list of nodes.
func (c *VirtualEnvironmentClient) ListDatastores(nodeName string, d *VirtualEnvironmentDatastoreListRequestBody) ([]*VirtualEnvironmentDatastoreListResponseData, error) {
	resBody := &VirtualEnvironmentDatastoreListResponseBody{}
	err := c.DoRequest(hmGET, fmt.Sprintf("nodes/%s/storage", url.PathEscape(nodeName)), d, resBody)

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
	err = c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/storage/%s/upload", url.PathEscape(d.NodeName), url.PathEscape(d.DatastoreID)), reqBody, resBody)

	if err != nil {
		return nil, err
	}

	return resBody, nil
}
