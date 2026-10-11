/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cluster

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindContainerResourceIgnoresQemuEntries(t *testing.T) {
	t.Parallel()

	// A qemu guest and a container may not share a VMID in PVE, but filtering on type keeps the
	// lookup honest if that ever changes.
	resources := []*ResourcesListResponseData{
		{Type: "qemu", VMID: 100, NodeName: "pve1"},
		{Type: "lxc", VMID: 101, NodeName: "pve2"},
	}

	res, err := FindContainerResource(resources, 101)
	require.NoError(t, err)
	require.Equal(t, "pve2", res.NodeName)

	_, err = FindContainerResource(resources, 100)
	require.ErrorIs(t, err, ErrVMDoesNotExist)
}
