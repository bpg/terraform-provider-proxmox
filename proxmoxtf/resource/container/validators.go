/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// MountTypeValidator returns a schema validation function for a mount type on a lxc container.
func MountTypeValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"cifs",
		"nfs",
	}, false))
}

// ConsoleModeValidator returns a schema validation function for a console mode on a lxc container.
func ConsoleModeValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"console",
		"shell",
		"tty",
	}, false))
}

// CPUArchitectureValidator returns a schema validation function for a CPU architecture on a lxc container.
func CPUArchitectureValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"amd64",
		"arm64",
		"armhf",
		"i386",
	}, false))
}

// OperatingSystemTypeValidator returns a schema validation function for an operating system type on a lxc container.
func OperatingSystemTypeValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"alpine",
		"archlinux",
		"centos",
		"debian",
		"devuan",
		"fedora",
		"gentoo",
		"nixos",
		"opensuse",
		"ubuntu",
		"unmanaged",
	}, false))
}
