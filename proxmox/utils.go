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

func parseDiskSize(size *string) (int, error) {
	var diskSize int
	var err error
	if size != nil {
		if strings.HasSuffix(*size, "T") {
			diskSize, err = strconv.Atoi(strings.TrimSuffix(*size, "T"))

			if err != nil {
				return -1, err
			}

			diskSize = int(math.Ceil(float64(diskSize) * 1024))
		} else if strings.HasSuffix(*size, "G") {
			diskSize, err = strconv.Atoi(strings.TrimSuffix(*size, "G"))

			if err != nil {
				return -1, err
			}
		} else if strings.HasSuffix(*size, "M") {
			diskSize, err = strconv.Atoi(strings.TrimSuffix(*size, "M"))

			if err != nil {
				return -1, err
			}

			diskSize = int(math.Ceil(float64(diskSize) / 1024))
		} else {
			return -1, fmt.Errorf("cannot parse storage size \"%s\"", *size)
		}
	}
	return diskSize, err
}
