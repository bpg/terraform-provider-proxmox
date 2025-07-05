/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package apt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	apitypes "github.com/bpg/terraform-provider-proxmox/proxmox/types/nodes/apt/repositories"
)

func TestStandardRepoHandleValueFromTerraform(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		val         tftypes.Value
		expected    func(val StandardRepoHandleValue) bool
		expectError bool
	}{
		"null value": {
			val: tftypes.NewValue(tftypes.String, nil),
			expected: func(val StandardRepoHandleValue) bool {
				return val.IsNull()
			},
		},
		"unknown value": {
			val: tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			expected: func(val StandardRepoHandleValue) bool {
				return val.IsUnknown()
			},
		},
		"invalid Ceph APT standard repository handle": {
			val: tftypes.NewValue(tftypes.String, "ceph-foo-enterprise"),
			expected: func(val StandardRepoHandleValue) bool {
				return val.kind == apitypes.StandardRepoHandleKindUnknown &&
					!val.IsCephHandle() &&
					!val.IsSupportedFilePath(apitypes.StandardRepoFilePathCeph) &&
					val.ComponentName() == "unknown" &&
					val.ValueString() == "ceph-foo-enterprise"
			},
		},
		"valid Ceph APT standard repository handle": {
			val: tftypes.NewValue(tftypes.String, "ceph-quincy-enterprise"),
			expected: func(val StandardRepoHandleValue) bool {
				return val.kind == apitypes.StandardRepoHandleKindEnterprise &&
					val.CephVersionName() == apitypes.CephVersionNameQuincy &&
					val.IsCephHandle() &&
					val.IsSupportedFilePath(apitypes.StandardRepoFilePathCeph) &&
					val.ComponentName() == "enterprise" &&
					val.ValueString() == "ceph-quincy-enterprise"
			},
		},
		`valid Ceph "test" APT standard repository handle`: {
			val: tftypes.NewValue(tftypes.String, "ceph-reef-test"),
			expected: func(val StandardRepoHandleValue) bool {
				return val.kind == apitypes.StandardRepoHandleKindTest &&
					val.CephVersionName() == apitypes.CephVersionNameReef &&
					val.IsCephHandle() &&
					val.IsSupportedFilePath(apitypes.StandardRepoFilePathCeph) &&
					val.ComponentName() == "test" &&
					val.ValueString() == "ceph-reef-test"
			},
		},
		"invalid APT repository handle": {
			val: tftypes.NewValue(tftypes.String, "foo-bar"),
			expected: func(val StandardRepoHandleValue) bool {
				return val.kind == apitypes.StandardRepoHandleKindUnknown &&
					!val.IsCephHandle() &&
					!val.IsSupportedFilePath(apitypes.StandardRepoFilePathCeph) &&
					val.ValueString() == "foo-bar"
			},
		},
		`valid APT "no subscription" repository handle`: {
			val: tftypes.NewValue(tftypes.String, "no-subscription"),
			expected: func(val StandardRepoHandleValue) bool {
				return val.kind == apitypes.StandardRepoHandleKindNoSubscription &&
					!val.IsCephHandle() &&
					val.IsSupportedFilePath(apitypes.StandardRepoFilePathMain) &&
					val.ComponentName() == "pve-no-subscription" &&
					val.ValueString() == "no-subscription"
			},
		},
		`valid APT "test" repository handle`: {
			val: tftypes.NewValue(tftypes.String, "test"),
			expected: func(val StandardRepoHandleValue) bool {
				return val.kind == apitypes.StandardRepoHandleKindTest &&
					!val.IsCephHandle() &&
					val.IsSupportedFilePath(apitypes.StandardRepoFilePathMain) &&
					val.ComponentName() == "pvetest" &&
					val.ValueString() == "test"
			},
		},
	}

	for name, test := range tests {
		t.Run(
			name, func(t *testing.T) {
				t.Parallel()

				ctx := t.Context()

				val, err := StandardRepoHandleType{}.ValueFromTerraform(ctx, test.val)
				if err == nil && test.expectError {
					t.Fatal("expected error, got no error")
				}

				if err != nil && !test.expectError {
					t.Fatalf("got unexpected error: %s", err)
				}

				if !test.expected(val.(StandardRepoHandleValue)) {
					t.Errorf("unexpected result")
				}
			},
		)
	}
}
