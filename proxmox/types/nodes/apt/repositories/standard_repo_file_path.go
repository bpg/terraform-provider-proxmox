/*
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package repositories

// Note that "hard-coded" slashes are used since Proxmox VE is built on top of Linux (Debian).
const (
	// StandardRepoFilePathCeph is the default Proxmox VE pre-defined (absolute) file path for the APT source list of Ceph
	// repositories (for PVE 9.0 and above using the modern DEB822 .sources format).
	StandardRepoFilePathCeph = "/etc/apt/sources.list.d/ceph.sources"

	// OldStandardRepoFilePathCeph is the legacy Proxmox VE pre-defined (absolute) file path for the APT source list of
	// Ceph repositories (for PVE versions before 9.0 using the legacy .list format).
	OldStandardRepoFilePathCeph = "/etc/apt/sources.list.d/ceph.list"

	// StandardRepoFilePathEnterprise is the default Proxmox VE pre-defined (absolute) file path for the APT source list
	// of enterprise repositories (for PVE 9.0 and above using the modern DEB822 .sources format).
	StandardRepoFilePathEnterprise = "/etc/apt/sources.list.d/pve-enterprise.sources"

	// OldStandardRepoFilePathEnterprise is the legacy Proxmox VE pre-defined (absolute) file path for the APT source
	// list of enterprise repositories (for PVE versions before 9.0 using the legacy .list format).
	OldStandardRepoFilePathEnterprise = "/etc/apt/sources.list.d/pve-enterprise.list"

	// StandardRepoFilePathMain is the default Proxmox VE pre-defined (absolute) file path for the APT source list of main
	// OS (Debian) repositories (for PVE 9.0 and above using the modern DEB822 .sources format).
	StandardRepoFilePathMain = "/etc/apt/sources.list.d/proxmox.sources"

	// OldStandardRepoFilePathMain is the legacy Proxmox VE pre-defined (absolute) file path for the APT source list of
	// main OS (Debian) repositories (for PVE versions before 9.0 using the legacy .list format).
	OldStandardRepoFilePathMain = "/etc/apt/sources.list"
)
