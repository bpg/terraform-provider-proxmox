/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package helper

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloseOrLogError(t *testing.T) {
	t.Parallel()

	f := CloseOrLogError(context.Background())

	c := &testCloser{}
	b := &badCloser{}

	func() {
		defer f(c)
		defer f(b)
		assert.Equal(t, false, c.isClosed)
	}()

	assert.Equal(t, true, c.isClosed)
}

type testCloser struct {
	isClosed bool
}

func (t *testCloser) Close() error {
	t.isClosed = true
	return nil
}

type badCloser struct{}

func (t *badCloser) Close() error {
	return fmt.Errorf("bad")
}
