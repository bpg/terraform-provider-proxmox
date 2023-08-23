/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ssh

import (
	"context"
)

// ProxmoxNode represents node address and port for SSH connection.
type ProxmoxNode struct {
	Address string
	Port    int32
}

// NodeResolver is an interface for resolving node names to IP addresses to use for SSH connection.
type NodeResolver interface {
	Resolve(ctx context.Context, nodeName string) (ProxmoxNode, error)
}
