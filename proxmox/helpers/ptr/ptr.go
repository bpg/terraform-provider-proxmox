/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ptr

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
