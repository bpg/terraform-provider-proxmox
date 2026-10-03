/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/vm/disk"
)

// efiDiskInterface is the fixed interface the provider uses for the efi_disk block.
const efiDiskInterface = "efidisk0"

// initializationInterfaceCollisionDiff rejects a configuration where initialization.interface
// is also used by a disk or efi_disk block.
func initializationInterfaceCollisionDiff(_ context.Context, d *schema.ResourceDiff, _ any) error {
	return validateInitializationInterface(d.GetRawConfig())
}

func validateInitializationInterface(config cty.Value) error {
	if config.IsNull() || !config.IsKnown() || !config.Type().IsObjectType() {
		return nil
	}

	initInterface := initializationInterface(config)
	if initInterface == "" {
		return nil
	}

	if config.Type().HasAttribute(mkEFIDisk) {
		efi := config.GetAttr(mkEFIDisk)
		if efi.IsKnown() && !efi.IsNull() && efi.LengthInt() > 0 && initInterface == efiDiskInterface {
			return fmt.Errorf(
				"%s.0.%s %q collides with the %s device, which always uses that interface",
				mkInitialization, mkInitializationInterface, initInterface, mkEFIDisk,
			)
		}
	}

	if !config.Type().HasAttribute(disk.MkDisk) {
		return nil
	}

	disks := config.GetAttr(disk.MkDisk)
	if !disks.IsKnown() || disks.IsNull() {
		return nil
	}

	for i, block := range disks.AsValueSlice() {
		if !block.IsKnown() || block.IsNull() || !block.Type().HasAttribute("interface") {
			continue
		}

		diskInterface, _ := block.GetAttr("interface").Unmark()
		if !diskInterface.IsKnown() || diskInterface.IsNull() {
			continue
		}

		if diskInterface.AsString() == initInterface {
			return fmt.Errorf(
				"%s.0.%s %q collides with %s.%d.interface, each device needs its own interface",
				mkInitialization, mkInitializationInterface, initInterface, disk.MkDisk, i,
			)
		}
	}

	return nil
}

// initializationInterface returns the configured initialization interface, or an empty string
// when it is not set or not known yet.
func initializationInterface(config cty.Value) string {
	if !config.Type().HasAttribute(mkInitialization) {
		return ""
	}

	blocks := config.GetAttr(mkInitialization)
	if !blocks.IsKnown() || blocks.IsNull() || blocks.LengthInt() == 0 {
		return ""
	}

	block := blocks.AsValueSlice()[0]
	if !block.IsKnown() || block.IsNull() || !block.Type().HasAttribute(mkInitializationInterface) {
		return ""
	}

	v, _ := block.GetAttr(mkInitializationInterface).Unmark()
	if !v.IsKnown() || v.IsNull() {
		return ""
	}

	return v.AsString()
}
