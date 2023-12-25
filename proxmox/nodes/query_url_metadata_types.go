/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes

// QueryUrlMetadataGetResponseBody contains the body from a QueryUrlMetadata get response.
type QueryUrlMetadataGetResponseBody struct {
	Data *QueryUrlMetadataGetResponseData `json:"data,omitempty"`
}

// QueryUrlMetadataGetResponseData contains the data from a QueryUrlMetadata get response.
type QueryUrlMetadataGetResponseData struct {
	Filename *string `json:"filename,omitempty"   url:"filename,omitempty"`
	Mimetype *string `json:"mimetype,omitempty"   url:"mimetype,omitempty"`
	Size     *int64  `json:"size,omitempty"   url:"size,omitempty"`
}

// QueryUrlMetadataGetRequestBody contains the body for a QueryUrlMetadata get request.
type QueryUrlMetadataGetRequestBody struct {
	Verify *bool   `json:"verify-certificates,omitempty"   url:"verify-certificates,omitempty"`
	URL    *string `json:"url,omitempty"   url:"url,omitempty"`
}
