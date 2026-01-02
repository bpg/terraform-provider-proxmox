/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeAPIClient struct{}

func (fakeAPIClient) DoRequest(_ context.Context, _, _ string, _, _ any) error { return nil }
func (fakeAPIClient) ExpandPath(path string) string                            { return path }
func (fakeAPIClient) IsRoot(_ context.Context) bool                            { return false }
func (fakeAPIClient) IsRootTicket(_ context.Context) bool                      { return false }
func (fakeAPIClient) HTTP() *http.Client                                       { return &http.Client{} }

func TestClientExpandPath(t *testing.T) {
	t.Parallel()

	c := &Client{Client: fakeAPIClient{}}

	assert.Equal(t, "storage", c.ExpandPath(""))
	assert.Equal(t, "storage/local", c.ExpandPath("local"))
	assert.Equal(t, "storage/foo%2Fbar", c.ExpandPath("foo/bar"))
}
