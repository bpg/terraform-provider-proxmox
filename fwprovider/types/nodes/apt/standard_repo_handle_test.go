/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package apt

import (
	"testing"

	goversion "github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"

	apitypes "github.com/bpg/terraform-provider-proxmox/proxmox/types/nodes/apt/repositories"
	"github.com/bpg/terraform-provider-proxmox/proxmox/version"
)

func TestStandardRepoHandleValueIsSupportedFilePath(t *testing.T) {
	t.Parallel()

	pve8Version := &version.ProxmoxVersion{Version: *goversion.Must(goversion.NewVersion("8.4.0"))}
	pve9Version := &version.ProxmoxVersion{Version: *goversion.Must(goversion.NewVersion("9.0.0"))}
	pve10Version := &version.ProxmoxVersion{Version: *goversion.Must(goversion.NewVersion("10.0.0"))}

	tests := []struct {
		name           string
		handle         StandardRepoHandleValue
		filePath       string
		proxmoxVersion *version.ProxmoxVersion
		expected       bool
	}{
		// Ceph handle tests with PVE 8.4 (old paths only)
		{
			name:           "Ceph handle with new path and PVE 8.4",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindEnterprise, cvn: apitypes.CephVersionNameQuincy},
			filePath:       apitypes.StandardRepoFilePathCeph,
			proxmoxVersion: pve8Version,
			expected:       false,
		},
		{
			name:           "Ceph handle with old path and PVE 8.4",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindEnterprise, cvn: apitypes.CephVersionNameQuincy},
			filePath:       apitypes.OldStandardRepoFilePathCeph,
			proxmoxVersion: pve8Version,
			expected:       true,
		},
		// Ceph handle tests with PVE 9.0 (both paths accepted)
		{
			name:           "Ceph handle with new path and PVE 9.0",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindEnterprise, cvn: apitypes.CephVersionNameQuincy},
			filePath:       apitypes.StandardRepoFilePathCeph,
			proxmoxVersion: pve9Version,
			expected:       true,
		},
		{
			name:           "Ceph handle with old path and PVE 9.0",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindEnterprise, cvn: apitypes.CephVersionNameQuincy},
			filePath:       apitypes.OldStandardRepoFilePathCeph,
			proxmoxVersion: pve9Version,
			expected:       true,
		},
		// Ceph handle tests with PVE 10.0 (both paths accepted)
		{
			name:           "Ceph handle with new path and PVE 10.0",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindTest, cvn: apitypes.CephVersionNameSquid},
			filePath:       apitypes.StandardRepoFilePathCeph,
			proxmoxVersion: pve10Version,
			expected:       true,
		},
		{
			name:           "Ceph handle with old path and PVE 10.0",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindTest, cvn: apitypes.CephVersionNameSquid},
			filePath:       apitypes.OldStandardRepoFilePathCeph,
			proxmoxVersion: pve10Version,
			expected:       true,
		},
		// Ceph handle tests with nil version (both paths accepted)
		{
			name:           "Ceph handle with new path and nil version",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindNoSubscription, cvn: apitypes.CephVersionNameReef},
			filePath:       apitypes.StandardRepoFilePathCeph,
			proxmoxVersion: nil,
			expected:       true,
		},
		{
			name:           "Ceph handle with old path and nil version",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindNoSubscription, cvn: apitypes.CephVersionNameReef},
			filePath:       apitypes.OldStandardRepoFilePathCeph,
			proxmoxVersion: nil,
			expected:       true,
		},
		// Enterprise handle tests with PVE 8.4 (old paths only)
		{
			name:           "Enterprise handle with new path and PVE 8.4",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindEnterprise, cvn: apitypes.CephVersionNameUnknown},
			filePath:       apitypes.StandardRepoFilePathEnterprise,
			proxmoxVersion: pve8Version,
			expected:       false,
		},
		{
			name:           "Enterprise handle with old path and PVE 8.4",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindEnterprise, cvn: apitypes.CephVersionNameUnknown},
			filePath:       apitypes.OldStandardRepoFilePathEnterprise,
			proxmoxVersion: pve8Version,
			expected:       true,
		},
		// Enterprise handle tests with PVE 9.0 (both paths accepted)
		{
			name:           "Enterprise handle with new path and PVE 9.0",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindEnterprise, cvn: apitypes.CephVersionNameUnknown},
			filePath:       apitypes.StandardRepoFilePathEnterprise,
			proxmoxVersion: pve9Version,
			expected:       true,
		},
		{
			name:           "Enterprise handle with old path and PVE 9.0",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindEnterprise, cvn: apitypes.CephVersionNameUnknown},
			filePath:       apitypes.OldStandardRepoFilePathEnterprise,
			proxmoxVersion: pve9Version,
			expected:       true,
		},
		// Enterprise handle tests with nil version (both paths accepted)
		{
			name:           "Enterprise handle with new path and nil version",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindEnterprise, cvn: apitypes.CephVersionNameUnknown},
			filePath:       apitypes.StandardRepoFilePathEnterprise,
			proxmoxVersion: nil,
			expected:       true,
		},
		{
			name:           "Enterprise handle with old path and nil version",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindEnterprise, cvn: apitypes.CephVersionNameUnknown},
			filePath:       apitypes.OldStandardRepoFilePathEnterprise,
			proxmoxVersion: nil,
			expected:       true,
		},
		// NoSubscription handle tests with PVE 8.4 (old paths only)
		{
			name:           "NoSubscription handle with new path and PVE 8.4",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindNoSubscription, cvn: apitypes.CephVersionNameUnknown},
			filePath:       apitypes.StandardRepoFilePathMain,
			proxmoxVersion: pve8Version,
			expected:       false,
		},
		{
			name:           "NoSubscription handle with old path and PVE 8.4",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindNoSubscription, cvn: apitypes.CephVersionNameUnknown},
			filePath:       apitypes.OldStandardRepoFilePathMain,
			proxmoxVersion: pve8Version,
			expected:       true,
		},
		// NoSubscription handle tests with PVE 9.0 (both paths accepted)
		{
			name:           "NoSubscription handle with new path and PVE 9.0",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindNoSubscription, cvn: apitypes.CephVersionNameUnknown},
			filePath:       apitypes.StandardRepoFilePathMain,
			proxmoxVersion: pve9Version,
			expected:       true,
		},
		{
			name:           "NoSubscription handle with old path and PVE 9.0",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindNoSubscription, cvn: apitypes.CephVersionNameUnknown},
			filePath:       apitypes.OldStandardRepoFilePathMain,
			proxmoxVersion: pve9Version,
			expected:       true,
		},
		// NoSubscription handle tests with nil version (both paths accepted)
		{
			name:           "NoSubscription handle with new path and nil version",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindNoSubscription, cvn: apitypes.CephVersionNameUnknown},
			filePath:       apitypes.StandardRepoFilePathMain,
			proxmoxVersion: nil,
			expected:       true,
		},
		{
			name:           "NoSubscription handle with old path and nil version",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindNoSubscription, cvn: apitypes.CephVersionNameUnknown},
			filePath:       apitypes.OldStandardRepoFilePathMain,
			proxmoxVersion: nil,
			expected:       true,
		},
		// Test handle tests with PVE 8.4 (old paths only)
		{
			name:           "Test handle with new path and PVE 8.4",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindTest, cvn: apitypes.CephVersionNameUnknown},
			filePath:       apitypes.StandardRepoFilePathMain,
			proxmoxVersion: pve8Version,
			expected:       false,
		},
		{
			name:           "Test handle with old path and PVE 8.4",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindTest, cvn: apitypes.CephVersionNameUnknown},
			filePath:       apitypes.OldStandardRepoFilePathMain,
			proxmoxVersion: pve8Version,
			expected:       true,
		},
		// Test handle tests with PVE 9.0 (both paths accepted)
		{
			name:           "Test handle with new path and PVE 9.0",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindTest, cvn: apitypes.CephVersionNameUnknown},
			filePath:       apitypes.StandardRepoFilePathMain,
			proxmoxVersion: pve9Version,
			expected:       true,
		},
		{
			name:           "Test handle with old path and PVE 9.0",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindTest, cvn: apitypes.CephVersionNameUnknown},
			filePath:       apitypes.OldStandardRepoFilePathMain,
			proxmoxVersion: pve9Version,
			expected:       true,
		},
		// Invalid path tests
		{
			name:           "Ceph handle with invalid path and PVE 9.0",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindEnterprise, cvn: apitypes.CephVersionNameQuincy},
			filePath:       "/invalid/path",
			proxmoxVersion: pve9Version,
			expected:       false,
		},
		{
			name:           "Enterprise handle with Ceph path and PVE 9.0",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindEnterprise, cvn: apitypes.CephVersionNameUnknown},
			filePath:       apitypes.StandardRepoFilePathCeph,
			proxmoxVersion: pve9Version,
			expected:       false,
		},
		{
			name:           "NoSubscription handle with Enterprise path and PVE 9.0",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindNoSubscription, cvn: apitypes.CephVersionNameUnknown},
			filePath:       apitypes.StandardRepoFilePathEnterprise,
			proxmoxVersion: pve9Version,
			expected:       false,
		},
		{
			name:           "Unknown handle with any path and PVE 9.0",
			handle:         StandardRepoHandleValue{kind: apitypes.StandardRepoHandleKindUnknown, cvn: apitypes.CephVersionNameUnknown},
			filePath:       apitypes.StandardRepoFilePathMain,
			proxmoxVersion: pve9Version,
			expected:       false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := test.handle.IsSupportedFilePath(test.filePath, test.proxmoxVersion)
			require.Equal(t, test.expected, result)
		})
	}
}

func TestStandardRepoHandleValueFromTerraform(t *testing.T) {
	t.Parallel()

	pve8Version := &version.ProxmoxVersion{Version: *goversion.Must(goversion.NewVersion("8.4.0"))}
	pve9Version := &version.ProxmoxVersion{Version: *goversion.Must(goversion.NewVersion("9.0.0"))}

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
					!val.IsSupportedFilePath(apitypes.StandardRepoFilePathCeph, nil) &&
					val.ComponentName(pve9Version) == "unknown" &&
					val.ValueString() == "ceph-foo-enterprise"
			},
		},
		`valid Ceph "enterprise" APT standard repository handle`: {
			val: tftypes.NewValue(tftypes.String, "ceph-quincy-enterprise"),
			expected: func(val StandardRepoHandleValue) bool {
				return val.kind == apitypes.StandardRepoHandleKindEnterprise &&
					val.CephVersionName() == apitypes.CephVersionNameQuincy &&
					val.IsCephHandle() &&
					val.IsSupportedFilePath(apitypes.StandardRepoFilePathCeph, nil) &&
					val.ComponentName(pve9Version) == "enterprise" &&
					val.ValueString() == "ceph-quincy-enterprise"
			},
		},
		`valid Ceph "no subscription" APT standard repository handle`: {
			val: tftypes.NewValue(tftypes.String, "ceph-reef-no-subscription"),
			expected: func(val StandardRepoHandleValue) bool {
				return val.kind == apitypes.StandardRepoHandleKindNoSubscription &&
					val.CephVersionName() == apitypes.CephVersionNameReef &&
					val.IsCephHandle() &&
					val.IsSupportedFilePath(apitypes.StandardRepoFilePathCeph, nil) &&
					val.ComponentName(pve9Version) == "no-subscription" &&
					val.ValueString() == "ceph-reef-no-subscription"
			},
		},
		`valid Ceph "test" APT standard repository handle`: {
			val: tftypes.NewValue(tftypes.String, "ceph-squid-test"),
			expected: func(val StandardRepoHandleValue) bool {
				return val.kind == apitypes.StandardRepoHandleKindTest &&
					val.CephVersionName() == apitypes.CephVersionNameSquid &&
					val.IsCephHandle() &&
					val.IsSupportedFilePath(apitypes.StandardRepoFilePathCeph, nil) &&
					val.ComponentName(pve9Version) == "test" &&
					val.ValueString() == "ceph-squid-test"
			},
		},
		"invalid APT repository handle": {
			val: tftypes.NewValue(tftypes.String, "foo-bar"),
			expected: func(val StandardRepoHandleValue) bool {
				return val.kind == apitypes.StandardRepoHandleKindUnknown &&
					!val.IsCephHandle() &&
					!val.IsSupportedFilePath(apitypes.StandardRepoFilePathCeph, nil) &&
					val.ValueString() == "foo-bar"
			},
		},
		`valid APT "no subscription" repository handle`: {
			val: tftypes.NewValue(tftypes.String, "no-subscription"),
			expected: func(val StandardRepoHandleValue) bool {
				return val.kind == apitypes.StandardRepoHandleKindNoSubscription &&
					!val.IsCephHandle() &&
					val.IsSupportedFilePath(apitypes.StandardRepoFilePathMain, nil) &&
					val.ComponentName(pve9Version) == "pve-no-subscription" &&
					val.ValueString() == "no-subscription"
			},
		},
		`valid APT "test" repository handle`: {
			val: tftypes.NewValue(tftypes.String, "test"),
			expected: func(val StandardRepoHandleValue) bool {
				return val.kind == apitypes.StandardRepoHandleKindTest &&
					!val.IsCephHandle() &&
					val.IsSupportedFilePath(apitypes.StandardRepoFilePathMain, nil) &&
					val.ComponentName(pve8Version) == "pvetest" &&
					val.ComponentName(pve9Version) == "pve-test" &&
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
