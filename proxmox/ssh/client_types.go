/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ssh

import (
	"context"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// Client is an interface for performing SSH requests against the Proxmox Nodes.
type Client interface {
	// ExecuteNodeCommands executes a command on a node.
	ExecuteNodeCommands(
		ctx context.Context, nodeAddress string,
		commands []string,
	) error

	// NodeUpload uploads a file to a node.
	NodeUpload(
		ctx context.Context, nodeAddress string,
		remoteFileDir string, fileUploadRequest *api.FileUploadRequest,
	) error
}
