/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package apt

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"

	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types/nodes/apt"
	api "github.com/bpg/terraform-provider-proxmox/proxmox/nodes/apt/repositories"
	pveversion "github.com/bpg/terraform-provider-proxmox/proxmox/version"
)

func TestStandardRepoImportFromAPIFilePath(t *testing.T) {
	t.Parallel()

	pve9 := &pveversion.ProxmoxVersion{Version: *version.Must(version.NewVersion("9.2.0"))}
	pve8 := &pveversion.ProxmoxVersion{Version: *version.Must(version.NewVersion("8.0.0"))}

	tests := []struct {
		name         string
		version      *pveversion.ProxmoxVersion
		handle       string
		data         string
		wantFilePath types.String
		wantIndex    types.Int64
	}{
		{
			name:    "well-known path",
			version: pve9,
			handle:  "no-subscription",
			data: `{"files":[{"path":"/etc/apt/sources.list.d/proxmox.sources","file-type":"sources",
				"repositories":[{"Components":["pve-no-subscription"],"Types":["deb"],"Suites":["trixie"],"URIs":["http://x"]}]}],
				"standard-repos":[{"handle":"no-subscription","Name":"No-Subscription","status":1}]}`,
			wantFilePath: types.StringValue("/etc/apt/sources.list.d/proxmox.sources"),
			wantIndex:    types.Int64Value(0),
		},
		{
			name:    "PVE 8 with a repository in proxmox.sources uses the defining file",
			version: pve8,
			handle:  "no-subscription",
			data: `{"files":[{"path":"/etc/apt/sources.list.d/proxmox.sources","file-type":"sources",
				"repositories":[{"Components":["pve-no-subscription"],"Types":["deb"],"Suites":["bookworm"],"URIs":["http://x"]}]}],
				"standard-repos":[{"handle":"no-subscription","Name":"No-Subscription","status":1}]}`,
			wantFilePath: types.StringValue("/etc/apt/sources.list.d/proxmox.sources"),
			wantIndex:    types.Int64Value(0),
		},
		{
			name:    "repository in a custom file while proxmox.sources also exists",
			version: pve9,
			handle:  "no-subscription",
			data: `{"files":[
				{"path":"/etc/apt/sources.list.d/proxmox.sources","file-type":"sources",
				"repositories":[{"Components":["pve-enterprise"],"Types":["deb"],"Suites":["trixie"],"URIs":["http://x"]}]},
				{"path":"/etc/apt/sources.list.d/custom.sources","file-type":"sources",
				"repositories":[{"Components":["main"],"Types":["deb"],"Suites":["trixie"],"URIs":["http://y"]},
				{"Components":["pve-no-subscription"],"Types":["deb"],"Suites":["trixie"],"URIs":["http://x"]}]}],
				"standard-repos":[{"handle":"no-subscription","Name":"No-Subscription","status":1}]}`,
			wantFilePath: types.StringValue("/etc/apt/sources.list.d/custom.sources"),
			wantIndex:    types.Int64Value(1),
		},
		{
			name:    "unexpected path uses the defining file",
			version: pve9,
			handle:  "no-subscription",
			data: `{"files":[{"path":"/etc/apt/sources.list.d/pve-no-subscription.sources","file-type":"sources",
				"repositories":[{"Components":["pve-no-subscription"],"Types":["deb"],"Suites":["trixie"],"URIs":["http://x"]}]}],
				"standard-repos":[{"handle":"no-subscription","Name":"No-Subscription","status":1}]}`,
			wantFilePath: types.StringValue("/etc/apt/sources.list.d/pve-no-subscription.sources"),
			wantIndex:    types.Int64Value(0),
		},
		{
			name:         "repository not found leaves no unknown value",
			version:      pve9,
			handle:       "no-subscription",
			data:         `{"files":[],"standard-repos":[{"handle":"no-subscription","Name":"No-Subscription"}]}`,
			wantFilePath: types.StringNull(),
			wantIndex:    types.Int64Null(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var data api.GetResponseData
			require.NoError(t, json.Unmarshal([]byte(tt.data), &data))

			handle, diags := customtypes.StandardRepoHandleType{}.ValueFromString(
				context.Background(), types.StringValue(tt.handle),
			)
			require.False(t, diags.HasError())

			srp := modelStandardRepo{
				Handle:   handle.(customtypes.StandardRepoHandleValue),
				Node:     types.StringValue("pve"),
				FilePath: types.StringUnknown(),
			}

			srp.importFromAPI(context.Background(), &data, tt.version)

			require.False(t, srp.FilePath.IsUnknown())
			require.Equal(t, tt.wantFilePath, srp.FilePath)
			require.Equal(t, tt.wantIndex, srp.Index)
		})
	}
}
