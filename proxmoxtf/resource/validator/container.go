/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package validator

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// MountType returns a schema validation function for a mount type on a lxc container.
func MountType() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"cifs",
		"nfs",
	}, false))
}

// ConsoleMode returns a schema validation function for a console mode on a lxc container.
func ConsoleMode() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"console",
		"shell",
		"tty",
	}, false))
}

// ContainerCPUArchitecture returns a schema validation function for a CPU architecture on a lxc container.
func ContainerCPUArchitecture() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"amd64",
		"arm64",
		"armhf",
		"i386",
	}, false))
}

// ContainerOperatingSystemType returns a schema validation function for an operating system type on a lxc container.
func ContainerOperatingSystemType() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"alpine",
		"archlinux",
		"centos",
		"debian",
		"fedora",
		"gentoo",
		"nixos",
		"opensuse",
		"ubuntu",
		"unmanaged",
	}, false))
}
