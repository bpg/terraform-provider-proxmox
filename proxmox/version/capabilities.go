/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package version

import "github.com/hashicorp/go-version"

// MinimumProxmoxVersion is the minimum supported Proxmox version by the provider.
//
//nolint:gochecknoglobals
var MinimumProxmoxVersion = ProxmoxVersion{*version.Must(version.NewVersion("8.0.0"))}

// SupportImportContentType checks if the Proxmox version supports the `import` content type when uploading disk images.
// See https://bugzilla.proxmox.com/show_bug.cgi?id=2424
func (v *ProxmoxVersion) SupportImportContentType() bool {
	return v.GreaterThanOrEqual(version.Must(version.NewVersion("8.4.0")))
}
