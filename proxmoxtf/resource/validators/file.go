/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package validators

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// ContentType returns a schema validation function for a content type on a storage device.
func ContentType() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"backup",
		"images",
		"iso",
		"rootdir",
		"snippets",
		"vztmpl",
		"import",
	}, false))
}

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
			r := regexp.MustCompile(`^(?i)[a-z\d\-_.]+:([a-z\d\-_]+/)?.+$`)
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

// FileMode is a schema validation function for file mode.
func FileMode() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(
		func(i interface{}, k string) ([]string, []error) {
			var errs []error

			v, ok := i.(string)
			if !ok {
				errs = append(errs, fmt.Errorf(
					`expected string in octal format (e.g. "0o700" or "0700"") for %q, but got %v of type %T`, k, v, i))

				return nil, errs
			}

			mode, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to parse file mode %q: %w", v, err))
				return nil, errs
			}

			if mode < 1 || mode > int64(^uint32(0)) {
				errs = append(errs, fmt.Errorf("%q must be in the range (%d - %d), got %d", v, 1, ^uint32(0), mode))
				return nil, errs
			}

			return []string{}, errs
		},
	)
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
