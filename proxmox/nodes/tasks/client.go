/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tasks

import (
	"fmt"
	"net/url"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// Client is an interface for performing requests against the Proxmox 'tasks' API.
type Client struct {
	api.Client
}

// ExpandPath expands a path relative to the client's base path.
func (c *Client) ExpandPath(_ string) string {
	panic("ExpandPath of tasks.Client must not be used. Use BuildPath instead.")
}

// BuildPath builds a path using information from Task ID.
func (c *Client) BuildPath(taskID string, path string) (string, error) {
	tid, err := ParseTaskID(taskID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("nodes/%s/tasks/%s/%s",
		url.PathEscape(tid.NodeName), url.PathEscape(taskID), url.PathEscape(path)), nil
}
