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

func TestPBSStorageModel_FromAPI_EncryptionKeyFingerprint(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("sets fingerprint when encryption key present", func(t *testing.T) {
		t.Parallel()

		encKeyJSON := `{"fingerprint":"sha256:abc123","data":"secret","created":"2025-08-18T15:04:05Z","modified":"2025-08-18T15:04:05Z"}`
		server := "pbs.example.com"
		datastore := "backup1"
		username := "user@pbs"
		storageID := "pbs-test"

		model := &PBSStorageModel{
			modelBase: modelBase{
				ID:           types.StringValue(storageID),
				Nodes:        types.SetNull(types.StringType),
				ContentTypes: types.SetNull(types.StringType),
			},
			EncryptionKeyFingerprint: types.StringUnknown(),
			GeneratedEncryptionKey:   types.StringUnknown(),
		}

		err := model.fromAPI(ctx, &proxmoxstorage.DatastoreGetResponseData{
			ID:            &storageID,
			Server:        &server,
			Datastore:     &datastore,
			Username:      &username,
			EncryptionKey: &encKeyJSON,
		})

		require.NoError(t, err)
		require.Equal(t, "sha256:abc123", model.EncryptionKeyFingerprint.ValueString())
	})

	t.Run("sets fingerprint to null when no encryption key", func(t *testing.T) {
		t.Parallel()

		server := "pbs.example.com"
		datastore := "backup1"
		username := "user@pbs"
		storageID := "pbs-test"

		model := &PBSStorageModel{
			modelBase: modelBase{
				ID:           types.StringValue(storageID),
				Nodes:        types.SetNull(types.StringType),
				ContentTypes: types.SetNull(types.StringType),
			},
			EncryptionKeyFingerprint: types.StringUnknown(),
			GeneratedEncryptionKey:   types.StringUnknown(),
		}

		err := model.fromAPI(ctx, &proxmoxstorage.DatastoreGetResponseData{
			ID:        &storageID,
			Server:    &server,
			Datastore: &datastore,
			Username:  &username,
		})

		require.NoError(t, err)
		require.True(t, model.EncryptionKeyFingerprint.IsNull(),
			"EncryptionKeyFingerprint should be null when no encryption key, got: %s", model.EncryptionKeyFingerprint)
	})
}

func TestPBSStorageModel_toCreateAPIRequest_IncludesPassword(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	m := &PBSStorageModel{
		modelBase: modelBase{
			ID:           types.StringValue("pbs-test"),
			Nodes:        types.SetNull(types.StringType),
			ContentTypes: types.SetNull(types.StringType),
		},
		Server:    types.StringValue("pbs.example.com"),
		Datastore: types.StringValue("backup1"),
		Username:  types.StringValue("user@pbs"),
		Password:  types.StringValue("mypassword"),
	}

	result, err := m.toCreateAPIRequest(ctx)
	require.NoError(t, err)

	req, ok := result.(proxmoxstorage.PBSStorageCreateRequest)
	require.True(t, ok)
	require.NotNil(t, req.Password, "Password must be present in create request")
	require.Equal(t, "mypassword", *req.Password)
}
