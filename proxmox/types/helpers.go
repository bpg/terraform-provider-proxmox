/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package types

// StrPtr returns a pointer to a string.
func StrPtr(s string) *string {
	return &s
}

// IntPtr returns a pointer to an int.
func IntPtr(i int) *int {
	return &i
}

// BoolPtr returns a pointer to a bool.
func BoolPtr(s bool) *CustomBool {
	customBool := CustomBool(s)
	return &customBool
}

// CopyString copies content of a string pointer.
func CopyString(s *string) *string {
	if s == nil {
		return nil
	}

	return StrPtr(*s)
}

// CopyInt copies content of an int pointer.
func CopyInt(i *int) *int {
	if i == nil {
		return nil
	}

	return IntPtr(*i)
}

// Int64PtrToIntPtr converts an int64 pointer to an int pointer.
func Int64PtrToIntPtr(i *int64) *int {
	if i == nil {
		return nil
	}

	return IntPtr(int(*i))
}
