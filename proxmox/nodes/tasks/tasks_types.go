/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tasks

// GetTaskStatusResponseBody contains the body from a node get task status response.
type GetTaskStatusResponseBody struct {
	Data *GetTaskStatusResponseData `json:"data,omitempty"`
}

// GetTaskStatusResponseData contains the data from a node get task status response.
type GetTaskStatusResponseData struct {
	PID      int    `json:"pid,omitempty"`
	Status   string `json:"status,omitempty"`
	ExitCode string `json:"exitstatus,omitempty"`
}
