/*
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package apt

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	apitypes "github.com/bpg/terraform-provider-proxmox/proxmox/types/nodes/apt/repositories"
)

// Ensure the implementations satisfy the required interfaces.
var (
	_ basetypes.StringTypable  = StandardRepoHandleType{}
	_ basetypes.StringValuable = StandardRepoHandleValue{}
)

// StandardRepoHandleType is a type that represents an APT standard repository handle.
type StandardRepoHandleType struct {
	basetypes.StringType
}

// StandardRepoHandleValue is a type that represents the value of an APT standard repository handle.
type StandardRepoHandleValue struct {
	basetypes.StringValue
	cvn  apitypes.CephVersionName
	kind apitypes.StandardRepoHandleKind
}

// Equal returns true if the two types are equal.
func (t StandardRepoHandleType) Equal(o attr.Type) bool {
	other, ok := o.(StandardRepoHandleType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

// String returns a string representation of the type.
func (t StandardRepoHandleType) String() string {
	return "StandardRepoHandleType"
}

// ValueFromString converts a string value to a basetypes.StringValuable.
func (t StandardRepoHandleType) ValueFromString(_ context.Context, in basetypes.StringValue) (
	basetypes.StringValuable,
	diag.Diagnostics,
) {
	value := StandardRepoHandleValue{
		StringValue: in,
		cvn:         apitypes.CephVersionNameUnknown,
		kind:        apitypes.StandardRepoHandleKindUnknown,
	}

	// Parse the Ceph version name when the handle has the prefix.
	if strings.HasPrefix(value.ValueString(), apitypes.CephStandardRepoHandlePrefix) {
		parts := strings.Split(value.ValueString(), "-")
		// Only continue when there is at least the Ceph prefix and the major version name in the handle.
		if len(parts) > 2 {
			cvn, err := apitypes.ParseCephVersionName(parts[1])
			if err == nil {
				value.cvn = cvn
			}
		}
	}

	// Parse the handle kind…
	handleString := value.ValueString()

	if value.IsCephHandle() {
		// …but ensure to strip Ceph specific parts from the handle string.
		name, ok := strings.CutPrefix(handleString, fmt.Sprintf("%s-%s-", apitypes.CephStandardRepoHandlePrefix, value.cvn))
		if ok {
			handleString = name
		}
	}

	switch handleString {
	case apitypes.StandardRepoHandleKindEnterprise.String():
		value.kind = apitypes.StandardRepoHandleKindEnterprise
	case apitypes.StandardRepoHandleKindNoSubscription.String():
		value.kind = apitypes.StandardRepoHandleKindNoSubscription
	case apitypes.StandardRepoHandleKindTest.String():
		value.kind = apitypes.StandardRepoHandleKindTest
	}

	return value, nil
}

// ValueFromTerraform converts a Terraform value to a basetypes.StringValuable.
func (t StandardRepoHandleType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, errors.Join(
			ErrValueConversion("unexpected error converting Terraform value to StringValue"),
			err,
		)
	}

	stringValue, ok := attrValue.(basetypes.StringValue)
	if !ok {
		return nil, ErrValueConversion("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)
	if diags.HasError() {
		return nil, ErrValueConversion(
			"unexpected error converting StringValue to StringValuable: %v",
			diags,
		)
	}

	return stringValuable, nil
}

// ValueType returns the underlying value type.
func (t StandardRepoHandleType) ValueType(_ context.Context) attr.Value {
	return StandardRepoHandleValue{}
}

// CephVersionName returns the corresponding Ceph major version name.
// Note that the version will be [apitypes.CephVersionNameUnknown] when not a Ceph specific handle!
func (v StandardRepoHandleValue) CephVersionName() apitypes.CephVersionName {
	return v.cvn
}

// ComponentName returns the corresponding component name.
func (v StandardRepoHandleValue) ComponentName() string {
	if v.cvn == apitypes.CephVersionNameUnknown && v.kind != apitypes.StandardRepoHandleKindUnknown {
		// For whatever reason the non-Ceph handle "test" kind does not use a dash in between the "pve" prefix.
		if v.kind == apitypes.StandardRepoHandleKindTest {
			return fmt.Sprintf("pve%s", v.kind)
		}

		return fmt.Sprintf("pve-%s", v.kind)
	}

	return v.kind.String()
}

// Equal returns true if the two values are equal.
func (v StandardRepoHandleValue) Equal(o attr.Value) bool {
	other, ok := o.(StandardRepoHandleValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

// IsCephHandle indicates if this is a Ceph APT standard repository.
func (v StandardRepoHandleValue) IsCephHandle() bool {
	return v.cvn != apitypes.CephVersionNameUnknown
}

// IsSupportedFilePath returns whether the handle is supported for the given source list file path.
func (v StandardRepoHandleValue) IsSupportedFilePath(filePath string) bool {
	switch filePath {
	case apitypes.StandardRepoFilePathCeph:
		return v.IsCephHandle()
	case apitypes.StandardRepoFilePathEnterprise:
		return !v.IsCephHandle() && v.kind == apitypes.StandardRepoHandleKindEnterprise
	case apitypes.StandardRepoFilePathMain:
		return !v.IsCephHandle() && v.kind != apitypes.StandardRepoHandleKindEnterprise
	default:
		return false
	}
}

// Type returns the type of the value.
func (v StandardRepoHandleValue) Type(_ context.Context) attr.Type {
	return StandardRepoHandleType{}
}
