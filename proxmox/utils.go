/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"context"
	"fmt"
	"io"
	"math"
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

var sizeregex = regexp.MustCompile(`(?i)^([\d\.]+)\s?(m|mb|g|gb|t|tb|p|pb)?$`)

func ParseDiskSize(size *string) (int, error) {
	if size == nil {
		return 0, nil
	}

	matches := sizeregex.FindStringSubmatch(*size)
	if len(matches) > 0 {
		fsize, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return -1, fmt.Errorf("cannot parse disk size \"%s\"", *size)
		}
		switch strings.ToLower(matches[2]) {
		case "m", "mb", "mib":
			return int(math.Ceil(fsize / 1024)), nil
		case "g", "gb", "gib":
			return int(math.Ceil(fsize)), nil
		case "t", "tb", "tib":
			return int(math.Ceil(fsize * 1024)), nil
		case "p", "pb", "pib":
			return int(math.Ceil(fsize * 1024 * 1024)), nil
		default:
			return int(math.Ceil(fsize)), nil
		}
	}

	return -1, fmt.Errorf("cannot parse disk size \"%s\"", *size)
}
