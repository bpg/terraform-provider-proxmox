/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ptr

import (
	"fmt"
	"strconv"
)

// Ptr creates a ptr from a value to use it inline.
func Ptr[T any](val T) *T {
	return &val
}

// ParseIntPtr parses a string to int and returns a pointer, with contextual error.
func ParseIntPtr(s, fieldName string) (*int, error) {
	iv, err := strconv.Atoi(s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", fieldName, err)
	}

	return &iv, nil
}

// ParseFloat64Ptr parses a string to float64 and returns a pointer, with contextual error.
func ParseFloat64Ptr(s, fieldName string) (*float64, error) {
	fv, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", fieldName, err)
	}

	return &fv, nil
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
