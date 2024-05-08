//go:build !windows

package ssh

import (
	"fmt"
	"net"
)

// dialSocket dials a Unix domain socket.
func dialSocket(address string) (net.Conn, error) {
	conn, err := net.Dial("unix", address)
	if err != nil {
		return nil, fmt.Errorf("error dialing unix socket: %w", err)
	}

	return conn, nil
}
