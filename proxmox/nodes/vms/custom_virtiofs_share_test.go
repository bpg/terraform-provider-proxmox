/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vms

import (
	"testing"

	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

func TestCustomVirtiofsShare_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		line    string
		want    *CustomVirtiofsShare
		wantErr bool
	}{
		{
			name: "id only virtiofs share",
			line: `"test"`,
			want: &CustomVirtiofsShare{
				DirId: "test",
			},
		},
		{
			name: "virtiofs share with more details",
			line: `"folder,cache=always"`,
			want: &CustomVirtiofsShare{
				DirId: "folder",
				Cache: ptr.Ptr("always"),
			},
		},
		{
			name: "virtiofs share with flags",
			line: `"folder,cache=never,direct-io=1,expose-acl=1"`,
			want: &CustomVirtiofsShare{
				DirId:       "folder",
				Cache:       ptr.Ptr("never"),
				DirectIo:    types.CustomBool(true).Pointer(),
				ExposeAcl:   types.CustomBool(true).Pointer(),
				ExposeXattr: types.CustomBool(true).Pointer(),
			},
		},
		{
			name: "virtiofs share with xattr",
			line: `"folder,expose-xattr=1"`,
			want: &CustomVirtiofsShare{
				DirId:       "folder",
				Cache:       nil,
				DirectIo:    types.CustomBool(false).Pointer(),
				ExposeAcl:   types.CustomBool(false).Pointer(),
				ExposeXattr: types.CustomBool(true).Pointer(),
			},
		},
		{
			name:    "virtiofs share invalid combination",
			line:    `"folder,expose-acl=1,expose-xattr=0"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &CustomVirtiofsShare{}
			if err := r.UnmarshalJSON([]byte(tt.line)); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
