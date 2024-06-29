/*
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package repositories

// Note that "hard-coded" slashes are used since Proxmox VE is built on top of Linux (Debian).
const (
	// StandardRepoFilePathCeph is the default Proxmox VE pre-defined (absolute) file path for the APT source list of Ceph
	// repositories.
	StandardRepoFilePathCeph = "/etc/apt/sources.list.d/ceph.list"

	// StandardRepoFilePathEnterprise is the default Proxmox VE pre-defined (absolute) file path for the APT source list
	// of enterprise repositories.
	StandardRepoFilePathEnterprise = "/etc/apt/sources.list.d/pve-enterprise.list"

	// StandardRepoFilePathMain is the default Proxmox VE pre-defined (absolute) file path for the APT source list of main
	// OS (Debian) repositories.
	StandardRepoFilePathMain = "/etc/apt/sources.list"
)
