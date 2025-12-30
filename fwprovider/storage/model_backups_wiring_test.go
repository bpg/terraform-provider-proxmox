package storage

import (
	"context"
	"testing"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestStorageModels_ToCreateAPIRequest_EncodesBackups(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("nfs", func(t *testing.T) {
		t.Parallel()

		m := &NFSStorageModel{
			StorageModelBase: StorageModelBase{
				ID:           types.StringValue("nfs-test"),
				Nodes:        types.SetNull(types.StringType),
				ContentTypes: types.SetNull(types.StringType),
				Disable:      types.BoolValue(false),
				Shared:       types.BoolValue(true),
			},
			Server:                 types.StringValue("127.0.0.1"),
			Export:                 types.StringValue("/export"),
			Options:                types.StringNull(),
			Preallocation:          types.StringNull(),
			SnapshotsAsVolumeChain: types.BoolValue(false),
			Backups: &BackupModel{
				MaxProtectedBackups: types.Int64Value(5),
				KeepAll:             types.BoolValue(false),
				KeepDaily:           types.Int64Value(7),
				KeepHourly:          types.Int64Null(),
				KeepLast:            types.Int64Null(),
				KeepWeekly:          types.Int64Null(),
				KeepMonthly:         types.Int64Null(),
				KeepYearly:          types.Int64Null(),
			},
		}

		reqAny, err := m.toCreateAPIRequest(ctx)
		require.NoError(t, err)

		values, err := query.Values(reqAny)
		require.NoError(t, err)

		require.Equal(t, "5", values.Get("max-protected-backups"))
		require.Equal(t, "keep-daily=7", values.Get("prune-backups"))
	})

	t.Run("cifs", func(t *testing.T) {
		t.Parallel()

		m := &CIFSStorageModel{
			StorageModelBase: StorageModelBase{
				ID:           types.StringValue("cifs-test"),
				Nodes:        types.SetNull(types.StringType),
				ContentTypes: types.SetNull(types.StringType),
				Disable:      types.BoolValue(false),
				Shared:       types.BoolValue(true),
			},
			Server:                 types.StringValue("127.0.0.1"),
			Username:               types.StringValue("user"),
			Password:               types.StringValue("pass"),
			Share:                  types.StringValue("share"),
			Domain:                 types.StringNull(),
			SubDirectory:           types.StringNull(),
			Preallocation:          types.StringNull(),
			SnapshotsAsVolumeChain: types.BoolValue(false),
			Backups: &BackupModel{
				MaxProtectedBackups: types.Int64Value(5),
				KeepAll:             types.BoolValue(false),
				KeepDaily:           types.Int64Value(7),
				KeepHourly:          types.Int64Null(),
				KeepLast:            types.Int64Null(),
				KeepWeekly:          types.Int64Null(),
				KeepMonthly:         types.Int64Null(),
				KeepYearly:          types.Int64Null(),
			},
		}

		reqAny, err := m.toCreateAPIRequest(ctx)
		require.NoError(t, err)

		values, err := query.Values(reqAny)
		require.NoError(t, err)

		require.Equal(t, "5", values.Get("max-protected-backups"))
		require.Equal(t, "keep-daily=7", values.Get("prune-backups"))
	})

	t.Run("pbs", func(t *testing.T) {
		t.Parallel()

		m := &PBSStorageModel{
			StorageModelBase: StorageModelBase{
				ID:           types.StringValue("pbs-test"),
				Nodes:        types.SetNull(types.StringType),
				ContentTypes: types.SetNull(types.StringType),
				Disable:      types.BoolValue(false),
				Shared:       types.BoolNull(),
			},
			Server:                   types.StringValue("127.0.0.1"),
			Datastore:                types.StringValue("ds"),
			Username:                 types.StringValue("user"),
			Password:                 types.StringValue("pass"),
			Namespace:                types.StringNull(),
			Fingerprint:              types.StringNull(),
			EncryptionKey:            types.StringNull(),
			EncryptionKeyFingerprint: types.StringNull(),
			GenerateEncryptionKey:    types.BoolNull(),
			GeneratedEncryptionKey:   types.StringNull(),
			Backups: &BackupModel{
				MaxProtectedBackups: types.Int64Value(5),
				KeepAll:             types.BoolValue(false),
				KeepDaily:           types.Int64Value(7),
				KeepHourly:          types.Int64Null(),
				KeepLast:            types.Int64Null(),
				KeepWeekly:          types.Int64Null(),
				KeepMonthly:         types.Int64Null(),
				KeepYearly:          types.Int64Null(),
			},
		}

		reqAny, err := m.toCreateAPIRequest(ctx)
		require.NoError(t, err)

		values, err := query.Values(reqAny)
		require.NoError(t, err)

		require.Equal(t, "5", values.Get("max-protected-backups"))
		require.Equal(t, "keep-daily=7", values.Get("prune-backups"))
	})
}
