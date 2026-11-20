/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cloudinitbootstrap

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"

	"github.com/diskfs/go-diskfs/backend"
	"github.com/diskfs/go-diskfs/backend/file"
	"github.com/diskfs/go-diskfs/filesystem/iso9660"
)

// noCloudVolumeLabel is the volume label cloud-init's NoCloud datasource looks for.
const noCloudVolumeLabel = "cidata"

const isoBlockSize = 2048

// seedFiles returns the NoCloud seed files keyed by their on-disk name, omitting empty optional ones.
func seedFiles(name, userData, metaData, networkConfig, vendorData string) map[string]string {
	if metaData == "" {
		metaData = fmt.Sprintf("instance-id: %s\nlocal-hostname: %s\n", name, name)
	}

	files := map[string]string{
		"meta-data": metaData,
		// cloud-init requires user-data to exist, even if empty.
		"user-data": userData,
	}

	if networkConfig != "" {
		files["network-config"] = networkConfig
	}

	if vendorData != "" {
		files["vendor-data"] = vendorData
	}

	return files
}

// buildISO writes a NoCloud seed ISO with the given files into dir, a fresh private directory, and returns its path.
// Rock Ridge and Joliet keep the lowercase, dashed file names cloud-init expects.
func buildISO(dir string, files map[string]string) (string, error) {
	workDir := filepath.Join(dir, "work")
	if err := os.Mkdir(workDir, 0o700); err != nil {
		return "", fmt.Errorf("failed to create ISO work directory: %w", err)
	}

	defer os.RemoveAll(workDir)

	var total int64
	for _, c := range files {
		total += int64(len(c))
	}

	// Headroom for system area, volume descriptors, path tables, directory records and per-file block padding.
	size := total + int64(len(files)+64)*isoBlockSize

	isoPath := filepath.Join(dir, noCloudVolumeLabel+".iso")

	b, err := file.CreateFromPath(isoPath, size)
	if err != nil {
		return "", fmt.Errorf("failed to create ISO file: %w", err)
	}

	ok := false

	defer func() {
		if !ok {
			_ = b.Close()
		}
	}()

	fs, err := iso9660.Create(b, size, 0, isoBlockSize, workDir)
	if err != nil {
		return "", fmt.Errorf("failed to create ISO filesystem: %w", err)
	}

	for name, content := range files {
		f, err := fs.OpenFile(name, os.O_CREATE|os.O_RDWR)
		if err != nil {
			return "", fmt.Errorf("failed to add %s to ISO: %w", name, err)
		}

		if _, err := f.Write([]byte(content)); err != nil {
			return "", fmt.Errorf("failed to write %s to ISO: %w", name, err)
		}

		if err := f.Close(); err != nil {
			return "", fmt.Errorf("failed to close %s in ISO: %w", name, err)
		}
	}

	err = fs.Finalize(iso9660.FinalizeOptions{
		RockRidge:        true,
		Joliet:           true,
		VolumeIdentifier: noCloudVolumeLabel,
	})
	if err != nil {
		return "", fmt.Errorf("failed to finalize ISO: %w", err)
	}

	if err := trimToVolumeSize(b, isoPath); err != nil {
		return "", err
	}

	ok = true

	return isoPath, nil
}

// trimToVolumeSize closes the backend and drops the unused headroom past the volume size recorded in the PVD.
func trimToVolumeSize(b backend.Storage, isoPath string) error {
	// Volume space size: little-endian uint32 at byte 80 of the primary volume descriptor (sector 16).
	buf := make([]byte, 4)
	if _, err := b.ReadAt(buf, 16*isoBlockSize+80); err != nil {
		return fmt.Errorf("failed to read ISO volume size: %w", err)
	}

	if err := b.Close(); err != nil {
		return fmt.Errorf("failed to close ISO file: %w", err)
	}

	if err := os.Truncate(isoPath, int64(binary.LittleEndian.Uint32(buf))*isoBlockSize); err != nil {
		return fmt.Errorf("failed to trim ISO file: %w", err)
	}

	return nil
}
