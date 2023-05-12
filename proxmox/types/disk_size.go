/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package types

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// Regex used to identify size strings. Case-insensitive. Covers megabytes, gigabytes and terabytes.
var sizeRegex = regexp.MustCompile(`(?i)^(\d+(\.\d+)?)(k|kb|kib|m|mb|mib|g|gb|gib|t|tb|tib)?$`)

// DiskSize allows a JSON integer value to also be a string. This is mapped to `<DiskSize>` data type in Proxmox API.
// Represents a disk size in bytes.
type DiskSize int64

// String returns the string representation of the disk size.
func (r DiskSize) String() string {
	return formatDiskSize(int64(r))
}

// InGigabytes returns the disk size in gigabytes.
func (r DiskSize) InGigabytes() int {
	return int(int64(r) / 1024 / 1024 / 1024)
}

// DiskSizeFromGigabytes creates a DiskSize from gigabytes.
func DiskSizeFromGigabytes(size int) DiskSize {
	return DiskSize(size * 1024 * 1024 * 1024)
}

// MarshalJSON marshals a disk size into a Proxmox API `<DiskSize>` string.
func (r DiskSize) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(formatDiskSize(int64(r)))
	if err != nil {
		return nil, fmt.Errorf("cannot marshal disk size: %w", err)
	}
	return bytes, nil
}

// UnmarshalJSON unmarshals a disk size from a Proxmox API `<DiskSize>` string.
func (r *DiskSize) UnmarshalJSON(b []byte) error {
	s := string(b)

	size, err := parseDiskSize(&s)
	if err != nil {
		return err
	}
	*r = DiskSize(size)

	return nil
}

// parseDiskSize parses a disk size string into a number of bytes.
func parseDiskSize(size *string) (int64, error) {
	if size == nil {
		return 0, nil
	}

	matches := sizeRegex.FindStringSubmatch(*size)
	if len(matches) > 0 {
		fsize, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return -1, fmt.Errorf("cannot parse disk size \"%s\": %w", *size, err)
		}
		switch strings.ToLower(matches[3]) {
		case "k", "kb", "kib":
			fsize *= 1024
		case "m", "mb", "mib":
			fsize = fsize * 1024 * 1024
		case "g", "gb", "gib":
			fsize = fsize * 1024 * 1024 * 1024
		case "t", "tb", "tib":
			fsize = fsize * 1024 * 1024 * 1024 * 1024
		}

		return int64(math.Ceil(fsize)), nil
	}

	return -1, fmt.Errorf("cannot parse disk size \"%s\"", *size)
}

func formatDiskSize(size int64) string {
	if size < 0 {
		return ""
	}

	if size < 1024 {
		return fmt.Sprintf("%d", size)
	}

	round := func(f float64) string {
		return strconv.FormatFloat(math.Ceil(f*100)/100, 'f', -1, 64)
	}

	if size < 1024*1024 {
		return round(float64(size)/1024) + "K"
	}

	if size < 1024*1024*1024 {
		return round(float64(size)/1024/1024) + "M"
	}

	if size < 1024*1024*1024*1024 {
		return round(float64(size)/1024/1024/1024) + "G"
	}

	return round(float64(size)/1024/1024/1024/1024) + "T"
}
