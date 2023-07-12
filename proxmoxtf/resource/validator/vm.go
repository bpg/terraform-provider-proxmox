/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package validator

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// VMID returns a schema validation function for a VM ID.
func VMID() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(i interface{}, k string) ([]string, []error) {
		min := 100
		max := 2147483647

		var ws []string
		var es []error

		v, ok := i.(int)

		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be int", k))
			return ws, es
		}

		if v != -1 {
			if v < min || v > max {
				es = append(es, fmt.Errorf("expected %s to be in the range (%d - %d), got %d", k, min, max, v))
				return ws, es
			}
		}

		return ws, es
	})
}
