/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ptr

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Ptr creates a ptr from a value to use it inline.
func Ptr[T any](val T) *T {
	return &val
}

// Or will dereference a pointer and return the given value if it's nil.
func Or[T any](p *T, or T) T {
	if p != nil {
		return *p
	}

	return or
}

// Eq compares two pointers and returns true if they are equal.
func Eq[T comparable](a, b *T) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	return *a == *b
}

// UpdateIfChanged updates dst with src if src is not nil and different from dst.
// Returns true if an update was made.
func UpdateIfChanged[T comparable](dst **T, src *T) bool {
	if src != nil && !Eq(*dst, src) {
		*dst = src
		return true
	}

	return false
}

// PtrOrNil safely gets a value of any type from schema.ResourceData.
// If the key is missing, returns nil. For strings, also returns nil if empty or whitespace.
func PtrOrNil[T any](d *schema.ResourceData, key string) *T {
	if v, ok := d.GetOk(key); ok {
		val := v.(T)

		// Special case: skip empty/whitespace-only strings
		if s, ok := any(val).(string); ok && strings.TrimSpace(s) == "" {
			return nil
		}

		return &val
	}

	return nil
}
