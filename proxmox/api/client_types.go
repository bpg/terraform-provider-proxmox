/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package api

import (
	"io"
	"os"
)

// MultiPartData enables multipart uploads in DoRequest.
type MultiPartData struct {
	Boundary string
	Reader   io.Reader
	Size     *int64
}

// ErrorResponseBody contains the body of an error response.
type ErrorResponseBody struct {
	Data   *string            `json:"data"`
	Errors *map[string]string `json:"errors"`
}

// FileUploadRequest is a request for uploading a file.
type FileUploadRequest struct {
	ContentType string
	FileName    string
	File        *os.File
}
