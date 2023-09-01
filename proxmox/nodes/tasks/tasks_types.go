/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tasks

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

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

// TaskID contains the components of a PVE task ID.
type TaskID struct {
	NodeName  string
	PID       int
	PStart    int
	StartTime time.Time
	Type      string
	ID        string
	User      string
}

// ParseTaskID parses a task ID into its component parts.
// The task ID is expected to be in the format of:
//
//	UPID:<node_name>:<pid_in_hex>:<pstart_in_hex>:<starttime_in_hex>:<type>:<id (optional)>:<user>@<realm>:
func ParseTaskID(taskID string) (TaskID, error) {
	parts := strings.SplitN(taskID, ":", 9)

	if parts[0] != "UPID" || len(parts) < 8 {
		return TaskID{}, fmt.Errorf("invalid task ID format: %s", taskID)
	}

	if parts[1] == "" {
		return TaskID{}, fmt.Errorf("missing node name in task ID: %s", taskID)
	}

	pid, err := strconv.ParseInt(parts[2], 16, 32)
	if err != nil {
		return TaskID{}, fmt.Errorf("error parsing task ID: %w", err)
	}

	pstart, err := strconv.ParseInt(parts[3], 16, 32)
	if err != nil {
		return TaskID{}, fmt.Errorf("error parsing pstart in task ID: %q: %w", taskID, err)
	}

	stime, err := strconv.ParseInt(parts[4], 16, 32)
	if err != nil {
		return TaskID{}, fmt.Errorf("error parsing start time in task ID: %q: %w", taskID, err)
	}

	if parts[5] == "" {
		return TaskID{}, fmt.Errorf("missing task type in task ID: %q", taskID)
	}

	if parts[7] == "" {
		return TaskID{}, fmt.Errorf("missing user in task ID: %q", taskID)
	}

	return TaskID{
		NodeName:  parts[1],
		PID:       int(pid),
		PStart:    int(pstart),
		StartTime: time.Unix(stime, 0),
		Type:      parts[5],
		ID:        parts[6],
		User:      parts[7],
	}, nil
}
