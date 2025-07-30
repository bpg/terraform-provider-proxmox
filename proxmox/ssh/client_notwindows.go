//go:build !windows

package ssh

import (
	"context"
	"errors"
	"fmt"
	"net"
)

// dialSocket dials a Unix domain socket.
func dialSocket(ctx context.Context, address string) (net.Conn, error) {
	if address == "" {
		return nil, errors.New("failed connecting to SSH agent socket: the socket file is not defined, " +
			"authentication will fall back to password")
	}

	dialer := &net.Dialer{}

	conn, err := dialer.DialContext(ctx, "unix", address)
	if err != nil {
		return nil, fmt.Errorf("error dialing unix socket: %w", err)
	}

	return conn, nil
}
