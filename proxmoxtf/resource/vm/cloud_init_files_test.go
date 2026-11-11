/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

func TestVMCloudInitFilesExplicitlyCleared(t *testing.T) {
	t.Parallel()

	initType := cty.Object(map[string]cty.Type{
		mkInitializationUserDataFileID:    cty.String,
		mkInitializationVendorDataFileID:  cty.String,
		mkInitializationNetworkDataFileID: cty.String,
		mkInitializationMetaDataFileID:    cty.String,
	})

	null := cty.NullVal(cty.String)
	str := cty.StringVal

	config := func(user, vendor, network, meta cty.Value) cty.Value {
		block := cty.ObjectVal(map[string]cty.Value{
			mkInitializationUserDataFileID:    user,
			mkInitializationVendorDataFileID:  vendor,
			mkInitializationNetworkDataFileID: network,
			mkInitializationMetaDataFileID:    meta,
		})

		return cty.ObjectVal(map[string]cty.Value{mkInitialization: cty.ListVal([]cty.Value{block})})
	}

	tests := []struct {
		name     string
		config   cty.Value
		expected bool
	}{
		{"user id explicitly empty", config(str(""), null, null, null), true},
		{"meta id explicitly empty", config(null, null, null, str("")), true},
		{"all omitted", config(null, null, null, null), false},
		{"set to a value", config(str("local:snippets/a.yaml"), null, null, null), false},
		{"one cleared, one set", config(str(""), str("local:snippets/v.yaml"), null, null), true},
		{"unknown value", config(cty.UnknownVal(cty.String), null, null, null), false},
		{"no initialization block", cty.ObjectVal(map[string]cty.Value{mkInitialization: cty.ListValEmpty(initType)}), false},
		{"null config", cty.NullVal(cty.DynamicPseudoType), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.expected, vmCloudInitFilesExplicitlyCleared(tt.config))
		})
	}
}

func TestVMCloneCicustomDeletes(t *testing.T) {
	t.Parallel()

	user := "local:snippets/user.yaml"
	vendor := "local:snippets/vendor.yaml"
	both := &vms.CustomCloudInitFiles{UserVolume: &user, VendorVolume: &vendor}

	tests := []struct {
		name     string
		del      []string
		template *vms.CustomCloudInitFiles
		wanted   *vms.CustomCloudInitFiles
		cleared  bool
		expected []string
	}{
		{"cleared, nothing left", nil, &vms.CustomCloudInitFiles{UserVolume: &user}, nil, true, []string{"cicustom"}},
		{"cleared, empty wanted files", nil, &vms.CustomCloudInitFiles{UserVolume: &user}, &vms.CustomCloudInitFiles{}, true, []string{"cicustom"}},
		{"cleared, other file still wanted", nil, &vms.CustomCloudInitFiles{UserVolume: &user}, &vms.CustomCloudInitFiles{VendorVolume: &vendor}, true, nil},
		{"user cleared, vendor omitted drops the whole inherited cicustom", nil, both, nil, true, []string{"cicustom"}},
		{"omitted keeps template value", nil, &vms.CustomCloudInitFiles{UserVolume: &user}, nil, false, nil},
		{"template has no custom files", nil, nil, nil, true, nil},
		{"template files all empty", nil, &vms.CustomCloudInitFiles{}, nil, true, nil},
		{"keeps existing deletes", []string{"usb0"}, &vms.CustomCloudInitFiles{UserVolume: &user}, nil, true, []string{"usb0", "cicustom"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.expected, vmCloneCicustomDeletes(tt.del, tt.template, tt.wanted, tt.cleared))
		})
	}
}
