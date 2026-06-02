/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ssh

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

func stallingListener(t *testing.T) (string, func()) {
	t.Helper()

	var lc net.ListenConfig

	ln, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to open stalling listener: %v", err)
	}

	var (
		mu    sync.Mutex
		conns []net.Conn
	)

	go func() {
		for {
			conn, aerr := ln.Accept()
			if aerr != nil {
				return
			}

			mu.Lock()

			conns = append(conns, conn)
			mu.Unlock()
		}
	}()

	return ln.Addr().String(), func() {
		_ = ln.Close()

		mu.Lock()
		for _, c := range conns {
			_ = c.Close()
		}
		mu.Unlock()
	}
}

func TestConnectHonorsContextDeadline(t *testing.T) {
	t.Parallel()

	addr, closeFn := stallingListener(t)
	defer closeFn()

	c := &client{username: "root"}

	cfg := &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{ssh.Password("irrelevant")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	type result struct {
		err     error
		elapsed time.Duration
	}

	done := make(chan result, 1)
	start := time.Now()

	go func() {
		_, err := c.connect(ctx, addr, cfg)
		done <- result{err: err, elapsed: time.Since(start)}
	}()

	select {
	case r := <-done:
		if r.err == nil {
			t.Fatalf("expected connect to fail against a stalling server, got nil error")
		}

		t.Logf("connect returned after %s with error: %v", r.elapsed, r.err)

		if r.elapsed > 5*time.Second {
			t.Fatalf("connect took %s — far past the 2s context deadline", r.elapsed)
		}
	case <-time.After(10 * time.Second):
		t.Fatalf("connect did not return within 10s: it ignored the 2s context deadline (issue #2915)")
	}
}
