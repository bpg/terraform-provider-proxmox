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
)

// ApplyConfig triggers a cluster-wide SDN apply via PUT /cluster/sdn.
func (c *Client) ApplyConfig(ctx context.Context) error {
	resBody := &ApplyResponseBody{}

	if err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath(""), nil, resBody); err != nil {
		return fmt.Errorf("error applying SDN configuration: %w", err)
	}

	if resBody.Data == nil || *resBody.Data == "" {
		return fmt.Errorf("SDN apply did not return a task UPID")
	}

	err := c.Tasks().WaitForTask(ctx, *resBody.Data)
	if err != nil {
		return fmt.Errorf("error waiting for SDN apply: %w", err)
	}

	return nil
}
