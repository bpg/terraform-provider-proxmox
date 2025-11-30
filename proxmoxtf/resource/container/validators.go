/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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

// EnvironmentVariablesValidator validates environment variable map based on Proxmox API requirements.
func EnvironmentVariablesValidator() schema.SchemaValidateDiagFunc {
	return func(v any, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		envVars, ok := v.(map[string]any)
		if !ok {
			return diags
		}

		// Proxmox regex for keys: \w+ (one or more word characters)
		keyRegex := regexp.MustCompile(`^\w+$`)

		// Proxmox rejects these control characters in values: \x00-\x08, \x10-\x1F, \x7F
		invalidValueCharsRegex := regexp.MustCompile(`[\x00-\x08\x10-\x1F\x7F]`)

		for key, val := range envVars {
			if !keyRegex.MatchString(key) {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Invalid environment variable key",
					Detail:   fmt.Sprintf("Environment variable key '%s' is invalid. Keys must contain only letters, digits, and underscores (matching \\w+)", key),
				})
			}

			if valStr, ok := val.(string); ok && invalidValueCharsRegex.MatchString(valStr) {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Invalid environment variable value",
					Detail:   fmt.Sprintf("Environment variable '%s' has an invalid value. Values cannot contain control characters", key),
				})
			}
		}

		return diags
	}
}
