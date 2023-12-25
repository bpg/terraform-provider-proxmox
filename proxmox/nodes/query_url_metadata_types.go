/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes

// QueryURLMetadataGetResponseBody contains the body from a QueryURLMetadata get response.
type QueryURLMetadataGetResponseBody struct {
	Data *QueryURLMetadataGetResponseData `json:"data,omitempty"`
}

// QueryURLMetadataGetResponseData contains the data from a QueryURLMetadata get response.
type QueryURLMetadataGetResponseData struct {
	Filename *string `json:"filename,omitempty" url:"filename,omitempty"`
	Mimetype *string `json:"mimetype,omitempty" url:"mimetype,omitempty"`
	Size     *int64  `json:"size,omitempty"     url:"size,omitempty"`
}

// QueryURLMetadataGetRequestBody contains the body for a QueryURLMetadata get request.
type QueryURLMetadataGetRequestBody struct {
	Verify *int    `json:"verify-certificates,omitempty" url:"verify-certificates,omitempty"`
	URL    *string `json:"url,omitempty"                 url:"url,omitempty"`
}
