/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package proxmox

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func CloseOrLogError(ctx context.Context) func(io.Closer) {
	return func(c io.Closer) {
		if err := c.Close(); err != nil {
			tflog.Error(ctx, "Failed to close", map[string]interface{}{
				"error": err,
			})
		}
	}
}

// Regex used to identify size strings. Case-insensitive. Covers megabytes, gigabytes and terabytes
var sizeRegex = regexp.MustCompile(`(?i)^(\d+)(m|mb|mib|g|gb|gib|t|tb|tib)$`)

// ParseDiskSize parses a disk size string into a number of gigabytes
func ParseDiskSize(size *string) (int, error) {
	if size == nil {
		return 0, nil
	}

	matches := sizeRegex.FindStringSubmatch(*size)
	if len(matches) > 0 {
		fsize, err := strconv.Atoi(matches[1])
		if err != nil {
			return -1, fmt.Errorf("cannot parse disk size \"%s\": %w", *size, err)
		}
		switch strings.ToLower(matches[2]) {
		case "m", "mb", "mib":
			return fsize / 1024, nil
		case "g", "gb", "gib":
			return fsize, nil
		case "t", "tb", "tib":
			return fsize * 1024, nil
		}
	}

	return -1, fmt.Errorf("cannot parse disk size \"%s\"", *size)
}
