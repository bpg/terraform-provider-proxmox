/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package applier

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/avast/retry-go/v5"
)

// ApplyConfig triggers a cluster-wide SDN apply via PUT /cluster/sdn.
// PVE may return a 500 error "got no worker upid - start worker failed", so we retry a few times.
func (c *Client) ApplyConfig(ctx context.Context) error {
	resBody := &ApplyResponseBody{}

	err := retry.New(
		retry.Context(ctx),
		retry.Attempts(3),
		retry.Delay(1*time.Second),
		retry.LastErrorOnly(true),
		retry.RetryIf(func(err error) bool {
			return strings.Contains(err.Error(), "got no worker upid")
		}),
	).Do(func() error {
		return c.DoRequest(ctx, http.MethodPut, c.ExpandPath(""), nil, resBody)
	})
	if err != nil {
		return fmt.Errorf("error applying SDN configuration: %w", err)
	}

	if resBody.Data == nil || *resBody.Data == "" {
		return fmt.Errorf("SDN apply did not return a task UPID")
	}

	err = c.Tasks().WaitForTask(ctx, *resBody.Data)
	if err != nil {
		return fmt.Errorf("error waiting for SDN apply: %w", err)
	}

	return nil
}
