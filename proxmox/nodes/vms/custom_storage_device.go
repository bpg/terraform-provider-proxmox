/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vms

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// StorageInterfaces is a list of supported storage interfaces.
//
//nolint:gochecknoglobals
var StorageInterfaces = []string{"ide", "sata", "scsi", "virtio"}

// CustomStorageDevice handles QEMU SATA device parameters.
type CustomStorageDevice struct {
	// FileVolume is the path to the storage device in format
	// "STORAGE_ID:SIZE_IN_GiB" or "STORAGE_ID:PATH_TO_FILE".
	// This is a required field.
	FileVolume string `json:"file" url:"file"`

	AIO                     *string           `json:"aio,omitempty"         url:"aio,omitempty"`
	Backup                  *types.CustomBool `json:"backup,omitempty"      url:"backup,omitempty,int"`
	BurstableReadSpeedMbps  *int              `json:"mbps_rd_max,omitempty" url:"mbps_rd_max,omitempty"`
	BurstableWriteSpeedMbps *int              `json:"mbps_wr_max,omitempty" url:"mbps_wr_max,omitempty"`
	Cache                   *string           `json:"cache,omitempty"       url:"cache,omitempty"`
	Discard                 *string           `json:"discard,omitempty"     url:"discard,omitempty"`
	Format                  *string           `json:"format,omitempty"      url:"format,omitempty"`
	IopsRead                *int              `json:"iops_rd,omitempty"     url:"iops_rd,omitempty"`
	IopsWrite               *int              `json:"iops_wr,omitempty"     url:"iops_wr,omitempty"`
	IOThread                *types.CustomBool `json:"iothread,omitempty"    url:"iothread,omitempty,int"`
	MaxIopsRead             *int              `json:"iops_rd_max,omitempty" url:"iops_rd_max,omitempty"`
	MaxIopsWrite            *int              `json:"iops_wr_max,omitempty" url:"iops_wr_max,omitempty"`
	MaxReadSpeedMbps        *int              `json:"mbps_rd,omitempty"     url:"mbps_rd,omitempty"`
	MaxWriteSpeedMbps       *int              `json:"mbps_wr,omitempty"     url:"mbps_wr,omitempty"`
	Media                   *string           `json:"media,omitempty"       url:"media,omitempty"`
	Replicate               *types.CustomBool `json:"replicate,omitempty"   url:"replicate,omitempty,int"`
	Serial                  *string           `json:"serial,omitempty"      url:"serial,omitempty"`
	Size                    *types.DiskSize   `json:"size,omitempty"        url:"size,omitempty"`
	SSD                     *types.CustomBool `json:"ssd,omitempty"         url:"ssd,omitempty,int"`
	DatastoreID             *string           `json:"-"                     url:"-"`
	FileID                  *string           `json:"-"                     url:"-"`
}

// CustomStorageDevices handles map of QEMU storage device per disk interface.
type CustomStorageDevices map[string]*CustomStorageDevice

// PathInDatastore returns path part of FileVolume or nil if it is not yet allocated.
func (d *CustomStorageDevice) PathInDatastore() *string {
	probablyDatastoreID, pathInDatastore, hasDatastoreID := strings.Cut(d.FileVolume, ":")
	if !hasDatastoreID {
		// when no ':' separator is found, 'Cut' places the whole string to 'probablyDatastoreID',
		// we want it in 'pathInDatastore' (as it is absolute filesystem path)
		pathInDatastore = probablyDatastoreID

		return &pathInDatastore
	}

	pathInDatastoreWithoutDigits := strings.Map(
		func(c rune) rune {
			if c < '0' || c > '9' {
				return -1
			}

			return c
		},
		pathInDatastore)

	if pathInDatastoreWithoutDigits == "" {
		// FileVolume is not yet allocated, it is in the "STORAGE_ID:SIZE_IN_GiB" format
		return nil
	}

	return &pathInDatastore
}

// IsOwnedBy returns true, if CustomStorageDevice is owned by given VM.
// Not yet allocated volumes are not owned by any VM.
func (d *CustomStorageDevice) IsOwnedBy(vmID int) bool {
	pathInDatastore := d.PathInDatastore()
	if pathInDatastore == nil {
		// not yet allocated volume, consider disk not owned by any VM
		// NOTE: if needed, create IsOwnedByOtherThan(vmId) instead of changing this return value.
		return false
	}

	// ZFS uses "local-zfs:vm-123-disk-0"
	if strings.HasPrefix(*pathInDatastore, fmt.Sprintf("vm-%d-", vmID)) {
		return true
	}

	// directory uses "local:123/vm-123-disk-0"
	if strings.HasPrefix(*pathInDatastore, fmt.Sprintf("%d/vm-%d-", vmID, vmID)) {
		return true
	}

	return false
}

// IsCloudInitDrive returns true, if CustomStorageDevice is a cloud-init drive.
func (d *CustomStorageDevice) IsCloudInitDrive(vmID int) bool {
	return d.Media != nil && *d.Media == "cdrom" &&
		strings.Contains(d.FileVolume, fmt.Sprintf("vm-%d-cloudinit", vmID))
}

// EncodeOptions converts a CustomStorageDevice's common options a URL value.
func (d *CustomStorageDevice) EncodeOptions() string {
	var values []string

	if d.AIO != nil {
		values = append(values, fmt.Sprintf("aio=%s", *d.AIO))
	}

	if d.Backup != nil {
		if *d.Backup {
			values = append(values, "backup=1")
		} else {
			values = append(values, "backup=0")
		}
	}

	if d.IopsRead != nil {
		values = append(values, fmt.Sprintf("iops_rd=%d", *d.IopsRead))
	}

	if d.IopsWrite != nil {
		values = append(values, fmt.Sprintf("iops_wr=%d", *d.IopsWrite))
	}

	if d.MaxIopsRead != nil {
		values = append(values, fmt.Sprintf("iops_rd_max=%d", *d.MaxIopsRead))
	}

	if d.MaxIopsWrite != nil {
		values = append(values, fmt.Sprintf("iops_wr_max=%d", *d.MaxIopsWrite))
	}

	if d.IOThread != nil {
		if *d.IOThread {
			values = append(values, "iothread=1")
		} else {
			values = append(values, "iothread=0")
		}
	}

	if d.Serial != nil && *d.Serial != "" {
		values = append(values, fmt.Sprintf("serial=%s", *d.Serial))
	}

	if d.SSD != nil {
		if *d.SSD {
			values = append(values, "ssd=1")
		} else {
			values = append(values, "ssd=0")
		}
	}

	if d.Discard != nil && *d.Discard != "" {
		values = append(values, fmt.Sprintf("discard=%s", *d.Discard))
	}

	if d.Cache != nil && *d.Cache != "" {
		values = append(values, fmt.Sprintf("cache=%s", *d.Cache))
	}

	if d.BurstableReadSpeedMbps != nil {
		values = append(values, fmt.Sprintf("mbps_rd_max=%d", *d.BurstableReadSpeedMbps))
	}

	if d.BurstableWriteSpeedMbps != nil {
		values = append(values, fmt.Sprintf("mbps_wr_max=%d", *d.BurstableWriteSpeedMbps))
	}

	if d.MaxReadSpeedMbps != nil {
		values = append(values, fmt.Sprintf("mbps_rd=%d", *d.MaxReadSpeedMbps))
	}

	if d.MaxWriteSpeedMbps != nil {
		values = append(values, fmt.Sprintf("mbps_wr=%d", *d.MaxWriteSpeedMbps))
	}

	if d.Replicate != nil {
		if *d.Replicate {
			values = append(values, "replicate=1")
		} else {
			values = append(values, "replicate=0")
		}
	}

	return strings.Join(values, ",")
}

// EncodeValues converts a CustomStorageDevice struct to a URL value.
func (d *CustomStorageDevice) EncodeValues(key string, v *url.Values) error {
	values := []string{
		fmt.Sprintf("file=%s", d.FileVolume),
	}

	if d.Format != nil {
		values = append(values, fmt.Sprintf("format=%s", *d.Format))
	}

	if d.Media != nil {
		values = append(values, fmt.Sprintf("media=%s", *d.Media))
	}

	if d.Size != nil {
		values = append(values, fmt.Sprintf("size=%d", *d.Size))
	}

	values = append(values, d.EncodeOptions())

	v.Add(key, strings.Join(values, ","))

	return nil
}

// UnmarshalJSON converts a CustomStorageDevice string to an object.
func (d *CustomStorageDevice) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomStorageDevice: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		//nolint:nestif
		if len(v) == 1 {
			d.FileVolume = v[0]

			// split file volume into datastore ID and path
			_, pathInDatastore, hasDatastoreID := strings.Cut(v[0], ":")
			if hasDatastoreID {
				// we don't set them here,... but probably should
				// d.DatastoreID = &probablyDatastoreID
				// d.FileID = &pathInDatastore
				ext := filepath.Ext(pathInDatastore)
				if ext != "" {
					format := string([]byte(ext)[1:])
					d.Format = &format
				}
			}
		} else if len(v) == 2 {
			switch v[0] {
			case "aio":
				d.AIO = &v[1]

			case "backup":
				bv := types.CustomBool(v[1] == "1")
				d.Backup = &bv

			case "cache":
				d.Cache = &v[1]

			case "discard":
				d.Discard = &v[1]

			case "file":
				d.FileVolume = v[1]

			case "format":
				d.Format = &v[1]

			case "iops_rd":
				iv, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to convert iops_rd to int: %w", err)
				}

				d.IopsRead = &iv

			case "iops_rd_max":
				iv, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to convert iops_rd_max to int: %w", err)
				}

				d.MaxIopsRead = &iv

			case "iops_wr":
				iv, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to convert iops_wr to int: %w", err)
				}

				d.IopsWrite = &iv

			case "iops_wr_max":
				iv, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to convert iops_wr_max to int: %w", err)
				}

				d.MaxIopsWrite = &iv

			case "iothread":
				bv := types.CustomBool(v[1] == "1")
				d.IOThread = &bv

			case "mbps_rd":
				iv, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to convert mbps_rd to int: %w", err)
				}

				d.MaxReadSpeedMbps = &iv

			case "mbps_rd_max":
				iv, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to convert mbps_rd_max to int: %w", err)
				}

				d.BurstableReadSpeedMbps = &iv

			case "mbps_wr":
				iv, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to convert mbps_wr to int: %w", err)
				}

				d.MaxWriteSpeedMbps = &iv

			case "mbps_wr_max":
				iv, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("failed to convert mbps_wr_max to int: %w", err)
				}

				d.BurstableWriteSpeedMbps = &iv

			case "media":
				d.Media = &v[1]

			case "replicate":
				bv := types.CustomBool(v[1] == "1")
				d.Replicate = &bv

			case "serial":
				d.Serial = &v[1]

			case "size":
				d.Size = new(types.DiskSize)

				err := d.Size.UnmarshalJSON([]byte(v[1]))
				if err != nil {
					return fmt.Errorf("failed to unmarshal disk size: %w", err)
				}

			case "ssd":
				bv := types.CustomBool(v[1] == "1")
				d.SSD = &bv
			}
		}
	}

	return nil
}

// MergeWith merges attributes of the given CustomStorageDevice with the current one.
// It will overwrite the current attributes with the given ones if they are not nil.
// The attributes that are not merged are:
//   - DatastoreID
//   - FileID
//   - FileVolume
//   - Format
//
// It will return true if any attribute of the current CustomStorageDevice was changed.
func (d *CustomStorageDevice) MergeWith(m CustomStorageDevice) bool {
	updated := ptr.UpdateIfChanged(&d.AIO, m.AIO)
	updated = ptr.UpdateIfChanged(&d.Backup, m.Backup) || updated
	updated = ptr.UpdateIfChanged(&d.BurstableReadSpeedMbps, m.BurstableReadSpeedMbps) || updated
	updated = ptr.UpdateIfChanged(&d.BurstableWriteSpeedMbps, m.BurstableWriteSpeedMbps) || updated
	updated = ptr.UpdateIfChanged(&d.Cache, m.Cache) || updated
	updated = ptr.UpdateIfChanged(&d.Discard, m.Discard) || updated
	updated = ptr.UpdateIfChanged(&d.IOThread, m.IOThread) || updated
	updated = ptr.UpdateIfChanged(&d.IopsRead, m.IopsRead) || updated
	updated = ptr.UpdateIfChanged(&d.IopsWrite, m.IopsWrite) || updated
	updated = ptr.UpdateIfChanged(&d.Media, m.Media) || updated
	updated = ptr.UpdateIfChanged(&d.MaxIopsRead, m.MaxIopsRead) || updated
	updated = ptr.UpdateIfChanged(&d.MaxIopsWrite, m.MaxIopsWrite) || updated
	updated = ptr.UpdateIfChanged(&d.MaxReadSpeedMbps, m.MaxReadSpeedMbps) || updated
	updated = ptr.UpdateIfChanged(&d.MaxWriteSpeedMbps, m.MaxWriteSpeedMbps) || updated
	updated = ptr.UpdateIfChanged(&d.Replicate, m.Replicate) || updated
	updated = ptr.UpdateIfChanged(&d.SSD, m.SSD) || updated
	updated = ptr.UpdateIfChanged(&d.Serial, m.Serial) || updated

	return updated
}

// Filter returns a map of CustomStorageDevices filtered by the given function.
func (d CustomStorageDevices) Filter(fn func(*CustomStorageDevice) bool) CustomStorageDevices {
	result := make(CustomStorageDevices)

	for k, v := range d {
		if fn(v) {
			result[k] = v
		}
	}

	return result
}

// EncodeValues converts a CustomStorageDevices array to multiple URL values.
func (d CustomStorageDevices) EncodeValues(_ string, v *url.Values) error {
	for s, d := range d {
		// Explicitly skip disks which have FileID set, so it won't be encoded in "Create" or "Update" operations.
		if d.FileID == nil || *d.FileID == "" {
			if err := d.EncodeValues(s, v); err != nil {
				return fmt.Errorf("error encoding storage device %s: %w", s, err)
			}
		}
	}

	return nil
}
