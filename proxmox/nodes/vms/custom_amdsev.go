/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vms

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// CustomAMDSEV handles AMDSEV parameters.
type CustomAMDSEV struct { // TODO verify the json parameters, "omitempty"?
	Type         *string           `json:"type"           url:"type"`
	AllowSMT     *types.CustomBool `json:"allow-smt" 			url:"allow-smt,int"`
	KernelHashes *types.CustomBool `json:"kernel-hashes"  url:"kernel-hashes,int"`
	NoDebug      *types.CustomBool `json:"no-debug"       url:"no-debug,int"`
	NoKeySharing *types.CustomBool `json:"no-key-sharing" url:"no-key-sharing,int"`
}

// EncodeValues converts a CustomAMDSEV struct to a URL value.
func (r *CustomAMDSEV) EncodeValues(key string, v *url.Values) error {
	var values []string

	if r.Type != nil {
		values = append(values, fmt.Sprintf("type=%s", *r.Type))
	}

	if r.AllowSMT != nil {
		if *r.AllowSMT {
			values = append(values, "allow-smt=1")
		} else {
			values = append(values, "allow-smt=0")
		}
	}

	if r.KernelHashes != nil {
		if *r.KernelHashes {
			values = append(values, "kernel-hashes=1")
		} else {
			values = append(values, "kernel-hashes=0")
		}
	}

	if r.NoDebug != nil {
		if *r.NoDebug {
			values = append(values, "no-debug=1")
		} else {
			values = append(values, "no-debug=0")
		}
	}

	if r.NoKeySharing != nil {
		if *r.NoKeySharing {
			values = append(values, "no-key-sharing=1")
		} else {
			values = append(values, "no-key-sharing=0")
		}
	}

	if len(values) > 0 {
		v.Add(key, strings.Join(values, ","))
	}

	return nil
}

// UnmarshalJSON converts a CustomAMDSEV string to an object.
func (r *CustomAMDSEV) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("error unmarshalling CustomAMDSEV: %w", err)
	}

	// TODO the get schema from proxmox only has `pve-qemu-sev-fmt` described, instead of each field
	// is it still the same fields or is the format different?
	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 1 {
			// TODO: can args be passed without values? most likely bools in that case
			continue
		} else if len(v) == 2 {
			switch v[0] {
			case "type":
				r.Type = &v[1]
			case "allow-smt":
				allow_smt := types.CustomBool(v[1] == "1")
				r.AllowSMT = &allow_smt
			case "kernel-hashes":
				kernel_hashes := types.CustomBool(v[1] == "1")
				r.KernelHashes = &kernel_hashes
			case "no-debug":
				no_debug := types.CustomBool(v[1] == "1")
				r.NoDebug = &no_debug
			case "no-key-sharing":
				no_key_sharing := types.CustomBool(v[1] == "1")
				r.NoKeySharing = &no_key_sharing
			}
		}
	}

	return nil
}
