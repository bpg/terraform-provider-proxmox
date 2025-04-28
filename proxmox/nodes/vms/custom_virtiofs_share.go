/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vms

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// CustomVirtiofsShare handles Virtiofs directory shares.
type CustomVirtiofsShare struct {
	DirId       string            `json:"dirid"                  url:"dirid"`
	Cache       *string           `json:"cache,omitempty"        url:"cache,omitempty"`
	DirectIo    *types.CustomBool `json:"direct-io,omitempty"    url:"direct-io,omitempty,int"`
	ExposeAcl   *types.CustomBool `json:"expose-acl,omitempty"   url:"expose-acl,omitempty,int"`
	ExposeXattr *types.CustomBool `json:"expose-xattr,omitempty" url:"expose-xattr,omitempty,int"`
}

// CustomVirtiofsShares handles Virtiofs directory shares.
type CustomVirtiofsShares map[string]*CustomVirtiofsShare

// EncodeValues converts a CustomVirtiofsShare struct to a URL value.
func (r *CustomVirtiofsShare) EncodeValues(key string, v *url.Values) error {
	if r.ExposeAcl != nil && *r.ExposeAcl && r.ExposeXattr != nil && !*r.ExposeXattr {
		// expose-xattr implies expose-acl
		return errors.New("expose_xattr must be omitted or true when expose_acl is enabled")
	}

	var values []string
	values = append(values, fmt.Sprintf("dirid=%s", r.DirId))

	if r.Cache != nil {
		values = append(values, fmt.Sprintf("cache=%s", *r.Cache))
	}

	if r.DirectIo != nil {
		if *r.DirectIo {
			values = append(values, "direct-io=1")
		} else {
			values = append(values, "direct-io=0")
		}
	}

	if r.ExposeAcl != nil {
		if *r.ExposeAcl {
			values = append(values, "expose-acl=1")
		} else {
			values = append(values, "expose-acl=0")
		}
	}

	if r.ExposeXattr != nil && (r.ExposeAcl == nil || !*r.ExposeAcl) {
		// expose-acl implies expose-xattr, omit it when unnecessary for consistency
		if *r.ExposeXattr {
			values = append(values, "expose-xattr=1")
		} else {
			values = append(values, "expose-xattr=0")
		}
	}

	v.Add(key, strings.Join(values, ","))

	return nil
}

// EncodeValues converts a CustomVirtiofsShares dict to multiple URL values.
func (r CustomVirtiofsShares) EncodeValues(_ string, v *url.Values) error {
	for s, d := range r {
		if err := d.EncodeValues(s, v); err != nil {
			return fmt.Errorf("failed to encode virtiofs share %s: %w", s, err)
		}
	}

	return nil
}

// UnmarshalJSON converts a CustomVirtiofsShare string to an object.
func (r *CustomVirtiofsShare) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal CustomVirtiofsShare: %w", err)
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 1 {
			r.DirId = v[0]
		} else if len(v) == 2 {
			switch v[0] {
			case "dirid":
				r.DirId = v[1]
			case "cache":
				r.Cache = &v[1]
			case "direct-io":
				bv := types.CustomBool(v[1] == "1")
				r.DirectIo = &bv
			case "expose-acl":
				bv := types.CustomBool(v[1] == "1")
				r.ExposeAcl = &bv
			case "expose-xattr":
				bv := types.CustomBool(v[1] == "1")
				r.ExposeXattr = &bv
			}
		}
	}

	// expose-acl implies expose-xattr
	if r.ExposeAcl != nil && *r.ExposeAcl {
		if r.ExposeXattr == nil {
			bv := types.CustomBool(true)
			r.ExposeAcl = &bv
		} else if !*r.ExposeXattr {
			return fmt.Errorf("failed to unmarshal CustomVirtiofsShare: expose-xattr contradicts the value of expose-acl")
		}
	}

	return nil
}
