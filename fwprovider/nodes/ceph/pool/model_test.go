/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package pool

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func erasureCodingAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"k":              types.Int64Type,
		"m":              types.Int64Type,
		"failure_domain": types.StringType,
		"device_class":   types.StringType,
		"profile":        types.StringType,
	}
}

func erasureCodingValue(t *testing.T, ec erasureCodingModel) types.Object {
	t.Helper()

	v, diags := types.ObjectValueFrom(context.Background(), erasureCodingAttrTypes(), ec)
	require.False(t, diags.HasError(), "%v", diags)

	return v
}

func TestErasureCodingModel_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ec   erasureCodingModel
		want string
	}{
		{
			name: "k and m only",
			ec: erasureCodingModel{
				K: types.Int64Value(2), M: types.Int64Value(1),
				FailureDomain: types.StringNull(), DeviceClass: types.StringNull(), Profile: types.StringNull(),
			},
			want: "k=2,m=1",
		},
		{
			name: "failure domain and device class",
			ec: erasureCodingModel{
				K: types.Int64Value(4), M: types.Int64Value(2),
				FailureDomain: types.StringValue("osd"), DeviceClass: types.StringValue("hdd"), Profile: types.StringNull(),
			},
			want: "k=4,m=2,failure-domain=osd,device-class=hdd",
		},
		{
			name: "existing profile",
			ec: erasureCodingModel{
				K: types.Int64Value(2), M: types.Int64Value(2),
				FailureDomain: types.StringNull(), DeviceClass: types.StringNull(), Profile: types.StringValue("shared"),
			},
			want: "k=2,m=2,profile=shared",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, tt.ec.String())
		})
	}
}

func TestCephPoolModel_ToDeleteParams_RemovesOnlyGeneratedProfiles(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name          string
		erasureCoding types.Object
		wantRemove    bool
	}{
		{
			name:          "replicated pool",
			erasureCoding: types.ObjectNull(erasureCodingAttrTypes()),
			wantRemove:    false,
		},
		{
			name: "profile created by Proxmox",
			erasureCoding: erasureCodingValue(t, erasureCodingModel{
				K: types.Int64Value(2), M: types.Int64Value(2),
				FailureDomain: types.StringValue("osd"), DeviceClass: types.StringNull(), Profile: types.StringNull(),
			}),
			wantRemove: true,
		},
		{
			name: "existing profile named by the user",
			erasureCoding: erasureCodingValue(t, erasureCodingModel{
				K: types.Int64Value(2), M: types.Int64Value(2),
				FailureDomain: types.StringNull(), DeviceClass: types.StringNull(), Profile: types.StringValue("shared"),
			}),
			wantRemove: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := cephPoolModel{ErasureCoding: tt.erasureCoding, ForceDestroy: types.BoolValue(false)}

			params, diags := m.toDeleteParams(ctx)
			require.False(t, diags.HasError(), "%v", diags)
			require.NotNil(t, params.RemoveECProfile)
			require.Equal(t, tt.wantRemove, bool(*params.RemoveECProfile))
			require.Nil(t, params.RemoveStorages, "remove_storages must never be sent")
		})
	}
}
