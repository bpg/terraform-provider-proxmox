/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package utils

import "reflect"

// SetDiff compares two slices of elements and returns the elements that are in the plan but not
// in the state (toCreate), the elements that are in the plan and in the state but are different (toUpdate),
// and the elements that are in the state but not in the plan (toDelete).
// The keyFunc is used to extract a unique key from each element to compare them.
func SetDiff[T any](plan []T, state []T, keyFunc func(t T) string) ([]T, []T, []T) {
	var toCreate, toUpdate, toDelete []T

	stateMap := map[string]T{}
	for _, s := range state {
		stateMap[keyFunc(s)] = s
	}

	planMap := map[string]T{}
	for _, p := range plan {
		planMap[keyFunc(p)] = p
	}

	for _, p := range plan {
		s, ok := stateMap[keyFunc(p)]
		if !ok {
			toCreate = append(toCreate, p)
		} else if !reflect.DeepEqual(p, s) {
			toUpdate = append(toUpdate, p)
		}
	}

	for _, s := range state {
		_, ok := planMap[keyFunc(s)]
		if !ok {
			toDelete = append(toDelete, s)
		}
	}

	return toCreate, toUpdate, toDelete
}
