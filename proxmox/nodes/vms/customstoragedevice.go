package vms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// CustomStorageDevice handles QEMU SATA device parameters.
type CustomStorageDevice struct {
	AIO                     *string           `json:"aio,omitempty"         url:"aio,omitempty"`
	Backup                  *types.CustomBool `json:"backup,omitempty"      url:"backup,omitempty,int"`
	BurstableReadSpeedMbps  *int              `json:"mbps_rd_max,omitempty" url:"mbps_rd_max,omitempty"`
	BurstableWriteSpeedMbps *int              `json:"mbps_wr_max,omitempty" url:"mbps_wr_max,omitempty"`
	Cache                   *string           `json:"cache,omitempty"       url:"cache,omitempty"`
	Discard                 *string           `json:"discard,omitempty"     url:"discard,omitempty"`
	FileVolume              string            `json:"file"                  url:"file"`
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
	Size                    *types.DiskSize   `json:"size,omitempty"        url:"size,omitempty"`
	SSD                     *types.CustomBool `json:"ssd,omitempty"         url:"ssd,omitempty,int"`
	DatastoreID             *string           `json:"-"                     url:"-"`
	Enabled                 bool              `json:"-"                     url:"-"`
	FileID                  *string           `json:"-"                     url:"-"`
	Interface               *string           `json:"-"                     url:"-"`
}

// PathInDatastore returns path part of FileVolume or nil if it is not yet allocated.
func (d CustomStorageDevice) PathInDatastore() *string {
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
func (d CustomStorageDevice) IsOwnedBy(vmID int) bool {
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
func (d CustomStorageDevice) IsCloudInitDrive(vmID int) bool {
	return d.Media != nil && *d.Media == "cdrom" &&
		strings.Contains(d.FileVolume, fmt.Sprintf("vm-%d-cloudinit", vmID))
}

// StorageInterface returns the storage interface of the CustomStorageDevice,
// e.g. "virtio" or "scsi" for "virtio0" or "scsi2".
func (d CustomStorageDevice) StorageInterface() string {
	for i, r := range *d.Interface {
		if unicode.IsDigit(r) {
			return (*d.Interface)[:i]
		}
	}

	// panic(fmt.Sprintf("cannot determine storage interface for disk interface '%s'", *d.Interface))
	return ""
}

// EncodeOptions converts a CustomStorageDevice's common options a URL value.
func (d CustomStorageDevice) EncodeOptions() string {
	values := []string{}

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
func (d CustomStorageDevice) EncodeValues(key string, v *url.Values) error {
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

// CustomStorageDevices handles map of QEMU storage device per disk interface.
type CustomStorageDevices map[string]*CustomStorageDevice

// ByStorageInterface returns a map of CustomStorageDevices filtered by the given storage interface.
func (d CustomStorageDevices) ByStorageInterface(storageInterface string) CustomStorageDevices {
	result := make(CustomStorageDevices)

	for k, v := range d {
		if v.StorageInterface() == storageInterface {
			result[k] = v
		}
	}

	return result
}

// EncodeValues converts a CustomStorageDevices array to multiple URL values.
func (d CustomStorageDevices) EncodeValues(_ string, v *url.Values) error {
	for s, d := range d {
		if d.Enabled {
			if err := d.EncodeValues(s, v); err != nil {
				return fmt.Errorf("error encoding storage device %s: %w", s, err)
			}
		}
	}

	return nil
}
