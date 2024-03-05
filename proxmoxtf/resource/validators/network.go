/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package validators

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// MACAddress is a schema validation function for MAC address.
func MACAddress() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(i interface{}, path string) ([]string, []error) {
		v, ok := i.(string)

		var ws []string

		var es []error

		if !ok {
			es = append(es, fmt.Errorf("expected type of %q to be string", path))
			return ws, es
		}

		if v != "" {
			r := regexp.MustCompile(`^[A-Fa-f0-9]{2}(:[A-Fa-f0-9]{2}){5}$`)
			ok := r.MatchString(v)

			if !ok {
				es = append(es, fmt.Errorf("expected %q to be a valid MAC address (A0:B1:C2:D3:E4:F5), got %q", path, v))
				return ws, es
			}
		}

		return ws, es
	})
}
