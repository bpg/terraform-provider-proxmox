/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package pools

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePoolMembershipID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		testName           string
		id                 string
		expectedPoolId     string
		expectedType       string
		expectedResourceId any
		expectError        bool
	}{
		{"correct vm id", "test-pool/vm/102", "test-pool", "vm", 102, false},
		{"correct storage id", "test-pool/storage/local-lvm", "test-pool", "storage", "local-lvm", false},
		{"wrong vm id format", "test-pool/vm/asdlasd", "", "", 0, true},
		{"missing pool id", "vm/102", "", "", 0, true},
		{"wrong id format", "test-pool/lxc/102/hello", "", "", 0, true},
		{"unknown type", "test-pool/hello/102", "", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			model, err := createMembershipModelFromID(tt.id)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, model)
				assert.Equal(t, tt.id, model.ID.ValueString())
				assert.Equal(t, tt.expectedPoolId, model.PoolID.ValueString())
				assert.Equal(t, tt.expectedType, model.Type.ValueString())

				var value any
				if model.VmID.IsNull() {
					value = model.StorageID.ValueString()
				} else {
					value = model.VmID.ValueInt64()
				}

				assert.EqualValues(t, tt.expectedResourceId, value)
			}
		})
	}
}

func TestGeneratePoolMembershipID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		model       poolMembershipModel
		expectedId  string
		expectError bool
	}{
		{
			"vm pool membership",
			poolMembershipModel{VmID: types.Int64Value(102), PoolID: types.StringValue("test-pool")},
			"test-pool/vm/102",
			false,
		},
		{
			"storage pool membership",
			poolMembershipModel{StorageID: types.StringValue("local-lvm"), PoolID: types.StringValue("test-pool")},
			"test-pool/storage/local-lvm",
			false,
		},
		{
			"missing any resource id",
			poolMembershipModel{PoolID: types.StringValue("test-pool")},
			"",
			true,
		},
		{"empty model", poolMembershipModel{}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			membershipType, err := tt.model.deduceMembershipType()
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				tt.model.Type = types.StringValue(membershipType)
			}

			id, err := tt.model.generateID()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expectedId, id.ValueString())
			}
		})
	}
}
