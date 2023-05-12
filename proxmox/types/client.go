/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package types

import (
	"context"
	"errors"
)

// ErrNoDataObjectInResponse is returned when the server does not include a data object in the response.
var ErrNoDataObjectInResponse = errors.New("the server did not include a data object in the response")

// Client is an interface for performing requests against the Proxmox API.
type Client interface {
	// DoRequest performs a request against the Proxmox API.
	DoRequest(
		ctx context.Context,
		method, path string,
		requestBody, responseBody interface{},
	) error

	// ExpandPath expands a path relative to the client's base path.
	// For example, if the client is configured for a VM and the
	// path is "firewall/options", the returned path will be
	// "/nodes/<node>/qemu/<vmid>/firewall/options".
	ExpandPath(path string) string
}
