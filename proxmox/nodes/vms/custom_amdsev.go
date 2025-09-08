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
type CustomAMDSEV struct {
	Type         string            `json:"type"           url:"type"`
	AllowSMT     *types.CustomBool `json:"allow-smt"      url:"allow-smt,int"`
	KernelHashes *types.CustomBool `json:"kernel-hashes"  url:"kernel-hashes,int"`
	NoDebug      *types.CustomBool `json:"no-debug"       url:"no-debug,int"`
	NoKeySharing *types.CustomBool `json:"no-key-sharing" url:"no-key-sharing,int"`
}

// EncodeValues converts a CustomAMDSEV struct to a URL value.
func (r *CustomAMDSEV) EncodeValues(key string, v *url.Values) error {
	values := []string{
		fmt.Sprintf("type=%s", r.Type),
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

	pairs := strings.Split(s, ",")

	for i, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 1 && i == 0 {
			r.Type = v[0]
		}

		if len(v) == 2 {
			switch v[0] {
			case "type":
				r.Type = v[1]
			case "allow-smt":
				r.AllowSMT = types.CustomBool(v[1] == "1").Pointer()
			case "kernel-hashes":
				r.KernelHashes = types.CustomBool(v[1] == "1").Pointer()
			case "no-debug":
				r.NoDebug = types.CustomBool(v[1] == "1").Pointer()
			case "no-key-sharing":
				r.NoKeySharing = types.CustomBool(v[1] == "1").Pointer()
			}
		}
	}

	return nil
}
