/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"

	proxmoxstorage "github.com/bpg/terraform-provider-proxmox/proxmox/storage"
)

func TestCIFSStorageModel_toCreateAPIRequest_IncludesPassword(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	m := &CIFSStorageModel{
		modelBase: modelBase{
			ID:           types.StringValue("cifs-test"),
			Nodes:        types.SetNull(types.StringType),
			ContentTypes: types.SetNull(types.StringType),
		},
		Server:   types.StringValue("nas.example.com"),
		Username: types.StringValue("user"),
		Password: types.StringValue("mypassword"),
		Share:    types.StringValue("backup"),
	}

	result, err := m.toCreateAPIRequest(ctx)
	require.NoError(t, err)

	req, ok := result.(proxmoxstorage.CIFSStorageCreateRequest)
	require.True(t, ok)
	require.NotNil(t, req.Password, "Password must be present in create request")
	require.Equal(t, "mypassword", *req.Password)
}
