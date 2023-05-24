/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package helper

import (
	"context"
	"io"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// CloseOrLogError closes the given io.Closer and logs any error.
func CloseOrLogError(ctx context.Context) func(io.Closer) {
	return func(c io.Closer) {
		if err := c.Close(); err != nil {
			tflog.Error(ctx, "Failed to close", map[string]interface{}{
				"error": err,
			})
		}
	}
}
