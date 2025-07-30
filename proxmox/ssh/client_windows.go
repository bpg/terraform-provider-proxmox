//go:build windows

package ssh

import (
	"context"
	"fmt"
	"net"

	"github.com/Microsoft/go-winio"
)

// dialSocket dials a Windows named pipe. If address is empty, it dials the default ssh-agent pipe.
func dialSocket(ctx context.Context, address string) (net.Conn, error) {
	if address == "" {
		address = `\\.\pipe\openssh-ssh-agent`
	}

	conn, err := winio.DialPipeContext(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("error dialing named pipe: %w", err)
	}

	return conn, nil
}
