/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ceph

// StatusResponseBody wraps the cluster Ceph status response.
type StatusResponseBody struct {
	Data *StatusResponseData `json:"data,omitempty"`
}

// StatusResponseData captures the typed subset of the Ceph status payload
// that callers currently care about. Unknown/extra fields in the response
// are silently dropped by json.Unmarshal, so adding a field here is the
// only step needed to expose more data.
type StatusResponseData struct {
	FSID        string       `json:"fsid"`
	Health      StatusHealth `json:"health"`
	QuorumNames []string     `json:"quorum_names"`
}

// StatusHealth is the `health` sub-object. Only `status` is exposed today;
// `checks` and other fields are version-dependent and intentionally omitted.
type StatusHealth struct {
	Status string `json:"status"`
}
