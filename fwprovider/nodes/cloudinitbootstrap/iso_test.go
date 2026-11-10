/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cloudinitbootstrap

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/diskfs/go-diskfs/backend/file"
	"github.com/diskfs/go-diskfs/filesystem/iso9660"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSeedFilesDefaultsMetaData(t *testing.T) {
	t.Parallel()

	files := seedFiles("vm-1", "#cloud-config\n", "", "", "")

	assert.Equal(t, map[string]string{
		"meta-data": "instance-id: vm-1\nlocal-hostname: vm-1\n",
		"user-data": "#cloud-config\n",
	}, files)
}

func TestSeedFilesIncludesOptionalFiles(t *testing.T) {
	t.Parallel()

	files := seedFiles("vm-1", "u", "m", "n", "v")

	assert.Equal(t, map[string]string{
		"meta-data":      "m",
		"user-data":      "u",
		"network-config": "n",
		"vendor-data":    "v",
	}, files)
}

func TestContentHash(t *testing.T) {
	t.Parallel()

	a := contentHash(map[string]string{"user-data": "x", "meta-data": "y"})

	assert.Len(t, a, 16)
	assert.Equal(t, a, contentHash(map[string]string{"meta-data": "y", "user-data": "x"}), "must not depend on map order")
	assert.NotEqual(t, a, contentHash(map[string]string{"user-data": "x", "meta-data": "y "}))
	assert.NotEqual(t,
		contentHash(map[string]string{"user-data": "ab", "meta-data": ""}),
		contentHash(map[string]string{"user-data": "a", "meta-data": "b"}),
	)
}

func TestBuildISO(t *testing.T) {
	t.Parallel()

	files := seedFiles("vm-1", "#cloud-config\nhostname: vm-1\n", "", "version: 2\n", "")

	dir := t.TempDir()

	path, err := buildISO(dir, files)
	require.NoError(t, err)

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	require.Len(t, entries, 1, "work directory must be cleaned up")

	b, err := file.OpenFromPath(path, true)
	require.NoError(t, err)

	t.Cleanup(func() { _ = b.Close() })

	fs, err := iso9660.Read(b, 0, 0, isoBlockSize)
	require.NoError(t, err)

	// go-diskfs NUL-pads the label; blkid and cloud-init read it as "cidata".
	assert.Equal(t, noCloudVolumeLabel, strings.TrimRight(fs.Label(), "\x00 "))

	st, err := os.Stat(path)
	require.NoError(t, err)
	assert.Zero(t, st.Size()%isoBlockSize, "ISO must be trimmed to whole blocks")
	assert.Less(t, st.Size(), int64(64*isoBlockSize), "ISO must be trimmed to its volume size")

	infos, err := fs.ReadDir(".")
	require.NoError(t, err)

	got := map[string]string{}

	for _, fi := range infos {
		f, err := fs.OpenFile(fi.Name(), os.O_RDONLY)
		require.NoError(t, err)

		data, err := io.ReadAll(f)
		require.NoError(t, err)

		got[fi.Name()] = strings.TrimRight(string(data), "\x00")
	}

	assert.Equal(t, files, got)
}
