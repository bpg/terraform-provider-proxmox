/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func ParseDiskSize(size *string) (int, error) {
	if size == nil {
		return 0, nil
	}

	if strings.HasSuffix(*size, "T") {
		diskSize, err := strconv.Atoi(strings.TrimSuffix(*size, "T"))
		if err != nil {
			return -1, fmt.Errorf("failed to parse disk size: %w", err)
		}
		return int(math.Ceil(float64(diskSize) * 1024)), nil
	}

	if strings.HasSuffix(*size, "G") {
		diskSize, err := strconv.Atoi(strings.TrimSuffix(*size, "G"))
		if err != nil {
			return -1, fmt.Errorf("failed to parse disk size: %w", err)
		}
		return diskSize, nil
	}

	if strings.HasSuffix(*size, "M") {
		diskSize, err := strconv.Atoi(strings.TrimSuffix(*size, "M"))
		if err != nil {
			return -1, fmt.Errorf("failed to parse disk size: %w", err)
		}
		return int(math.Ceil(float64(diskSize) / 1024)), nil
	}

	return -1, fmt.Errorf("cannot parse disk size \"%s\"", *size)
}
