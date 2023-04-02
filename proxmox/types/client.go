package types

import "context"

type Client interface {
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
