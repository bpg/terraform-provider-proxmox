package disk

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

func TestModel_toAPI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		model    Model
		wantSize int64
		wantAIO  *string
	}{
		{
			name: "minimal model with size only",
			model: Model{
				Aio:         types.StringNull(),
				Backup:      types.BoolNull(),
				Cache:       types.StringNull(),
				DatastoreId: types.StringValue("local-lvm"),
				Discard:     types.StringNull(),
				FileFormat:  types.StringNull(),
				ImportFrom:  types.StringNull(),
				IOThread:    types.BoolNull(),
				Size:        types.Int64Value(8),
			},
			wantSize: 8,
			wantAIO:  nil,
		},
		{
			name: "model with all fields set",
			model: Model{
				Aio:         types.StringValue("native"),
				Backup:      types.BoolValue(false),
				Cache:       types.StringValue("writeback"),
				DatastoreId: types.StringValue("local-lvm"),
				Discard:     types.StringValue("on"),
				FileFormat:  types.StringValue("raw"),
				ImportFrom:  types.StringNull(),
				IOThread:    types.BoolValue(true),
				Size:        types.Int64Value(16),
			},
			wantSize: 16,
			wantAIO:  new("native"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tt.model.toAPI()

			assert.Equal(t, tt.wantSize, result.Size.InGigabytes())
			assert.Equal(t, tt.wantAIO, result.AIO)
			assert.NotNil(t, result.Media)
			assert.Equal(t, "disk", *result.Media)
		})
	}
}

func TestModel_fromAPI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		device   vms.CustomStorageDevice
		wantSize int64
		wantAIO  types.String
	}{
		{
			name: "device with nil optional fields",
			device: vms.CustomStorageDevice{
				Size:        proxmoxtypes.DiskSizeFromGigabytes(8),
				DatastoreID: new("local-lvm"),
			},
			wantSize: 8,
			wantAIO:  types.StringNull(),
		},
		{
			name: "device with all fields populated",
			device: vms.CustomStorageDevice{
				AIO:         new("native"),
				Backup:      proxmoxtypes.CustomBoolPtr(new(false)),
				Cache:       new("writeback"),
				DatastoreID: new("local-lvm"),
				Discard:     new("on"),
				Format:      new("raw"),
				IOThread:    proxmoxtypes.CustomBoolPtr(new(true)),
				Size:        proxmoxtypes.DiskSizeFromGigabytes(16),
			},
			wantSize: 16,
			wantAIO:  types.StringValue("native"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var m Model
			m.fromAPI(tt.device)

			assert.Equal(t, tt.wantSize, m.Size.ValueInt64())
			assert.Equal(t, tt.wantAIO, m.Aio)
		})
	}
}

func TestNewValue_NoDisks(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	diags := diag.Diagnostics{}

	config := &vms.GetResponseData{}

	result := NewValue(ctx, config, &diags)

	require.False(t, diags.HasError())
	assert.True(t, result.IsNull())
}

//go:fix inline
func ptrTo[T any](v T) *T {
	return new(v)
}
