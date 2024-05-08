//go:build windows

package ssh

import (
	"fmt"
	"net"

	"github.com/Microsoft/go-winio"
)

// dialSocket dials a Windows named pipe.
func dialSocket(address string) (net.Conn, error) {
	conn, err := winio.DialPipe(address, nil)
	if err != nil {
		return nil, fmt.Errorf("error dialing named pipe: %w", err)
	}

	return conn, nil
}
