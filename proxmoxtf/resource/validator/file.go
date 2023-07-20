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

	"github.com/bpg/terraform-provider-proxmox/internal/types"
)

// FileFormat returns a schema validation function for a file format.
func FileFormat() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"qcow2",
		"raw",
		"vmdk",
	}, false))
}

// FileID returns a schema validation function for a file identifier.
func FileID() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(i interface{}, k string) ([]string, []error) {
		v, ok := i.(string)

		var ws []string
		var es []error

		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return ws, es
		}

		if v != "" {
			r := regexp.MustCompile(`^(?i)[a-z\d\-_]+:([a-z\d\-_]+/)?.+$`)
			ok := r.MatchString(v)

			if !ok {
				es = append(es, fmt.Errorf(
					"expected %s to be a valid file identifier (datastore-name:iso/some-file.img), got %s", k, v,
				))
				return ws, es
			}
		}

		return ws, es
	})
}

// FileSize is a schema validation function for file size.
func FileSize() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(i interface{}, k string) ([]string, []error) {
		v, ok := i.(string)
		var es []error

		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return nil, es
		}

		if v != "" {
			_, err := types.ParseDiskSize(v)
			if err != nil {
				es = append(es, fmt.Errorf("expected %s to be a valid file size (100, 1M, 1G), got %s", k, v))
				return nil, es
			}
		}

		return []string{}, es
	})
}
