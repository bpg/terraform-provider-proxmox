/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package validator

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// MACAddress is a schema validation function for MAC address.
func MACAddress() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(i interface{}, k string) ([]string, []error) {
		v, ok := i.(string)

		var ws []string
		var es []error

		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return ws, es
		}

		if v != "" {
			r := regexp.MustCompile(`^[A-Z\d]{2}(:[A-Z\d]{2}){5}$`)
			ok := r.MatchString(v)

			if !ok {
				es = append(es, fmt.Errorf("expected %s to be a valid MAC address (A0:B1:C2:D3:E4:F5), got %s", k, v))
				return ws, es
			}
		}

		return ws, es
	})
}
