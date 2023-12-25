/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodestorage

// DatastoreUploadResponseBody contains the body from a datastore upload response.
type DatastoreUploadResponseBody struct {
	UploadID *string `json:"data,omitempty"`
}
