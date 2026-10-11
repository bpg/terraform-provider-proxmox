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

// defaultInitializationInterface is the interface the provider picks for the cloud-init drive
// on create when initialization.interface is not set.
const defaultInitializationInterface = "ide2"

// initializationInterfaceCollisionDiff rejects a configuration where initialization.interface
// is also used by a disk block.
func initializationInterfaceCollisionDiff(_ context.Context, d *schema.ResourceDiff, _ any) error {
	// A clone reuses the template's cloud-init drive when the interface is unset, so ide2 is only certain for a
	// plain create.
	assumeDefault := d.Id() == "" && len(d.Get(mkClone).([]any)) == 0

	return validateInitializationInterface(d.GetRawConfig(), assumeDefault)
}

// validateInitializationInterface checks the raw config for a disk that uses the same interface as the
// initialization block. An unset initialization.interface is checked as ide2 only when assumeDefault is set; on
// update and on clone the provider detects an existing drive instead.
func validateInitializationInterface(config cty.Value, assumeDefault bool) error {
	if config.IsNull() || !config.IsKnown() || !config.Type().IsObjectType() {
		return nil
	}

	initInterface, known := initializationInterface(config)
	if !known {
		return nil
	}

	if initInterface == "" {
		if !assumeDefault {
			return nil
		}

		initInterface = defaultInitializationInterface
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

// initializationInterface returns the configured initialization interface. The second value is false
// when there is no initialization block or the interface is not known yet, so nothing can be checked.
// An unset interface is returned as an empty string.
func initializationInterface(config cty.Value) (string, bool) {
	if !config.Type().HasAttribute(mkInitialization) {
		return "", false
	}

	blocks := config.GetAttr(mkInitialization)
	if !blocks.IsKnown() || blocks.IsNull() || blocks.LengthInt() == 0 {
		return "", false
	}

	block := blocks.AsValueSlice()[0]
	if !block.IsKnown() || block.IsNull() {
		return "", false
	}

	if !block.Type().HasAttribute(mkInitializationInterface) {
		return "", true
	}

	v, _ := block.GetAttr(mkInitializationInterface).Unmark()
	if !v.IsKnown() {
		return "", false
	}

	if v.IsNull() {
		return "", true
	}

	return v.AsString(), true
}
