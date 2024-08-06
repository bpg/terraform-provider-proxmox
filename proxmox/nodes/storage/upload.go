package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// APIUpload uploads a file to a datastore using the Proxmox API.
func (c *Client) APIUpload(
	ctx context.Context,
	d *api.FileUploadRequest,
	tempDir string,
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
	tempMultipartFile, err := os.CreateTemp(tempDir, "multipart")
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
			if errors.Is(e, os.ErrClosed) {
				// We can ignore the error in the case that the file was already closed.
				return
			}

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
		c.ExpandPath("upload"),
		reqBody,
		resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error uploading file to datastore %s: %w", c.StorageName, err)
	}

	if resBody.UploadID == nil {
		return nil, fmt.Errorf("error uploading file to datastore %s: no uploadID", c.StorageName)
	}

	err = c.Tasks().WaitForTask(ctx, *resBody.UploadID)
	if err != nil {
		return nil, fmt.Errorf("error uploading file to datastore %s: failed waiting for upload - %w", c.StorageName, err)
	}

	return resBody, nil
}
