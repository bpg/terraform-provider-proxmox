/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetDiff(t *testing.T) {
	t.Parallel()

	type args[T any] struct {
		plan  []T
		state []T
	}

	type testCase[T any] struct {
		name         string
		args         args[T]
		wantToCreate []T
		wantToUpdate []T
		wantToDelete []T
	}

	type rec struct {
		n string
		v string
	}

	tests := []testCase[rec]{
		{
			name: "empty",
			args: args[rec]{
				plan:  []rec{},
				state: []rec{},
			},
			wantToCreate: nil,
			wantToUpdate: nil,
			wantToDelete: nil,
		},
		{
			name: "create",
			args: args[rec]{
				plan:  []rec{{"a", "1"}, {"b", "2"}, {"c", "3"}},
				state: []rec{},
			},
			wantToCreate: []rec{{"a", "1"}, {"b", "2"}, {"c", "3"}},
			wantToUpdate: nil,
			wantToDelete: nil,
		},
		{
			name: "create and delete",
			args: args[rec]{
				plan:  []rec{{"a", "1"}, {"b", "2"}, {"c", "3"}},
				state: []rec{{"b", "2"}, {"c", "3"}, {"d", "4"}},
			},
			wantToCreate: []rec{{"a", "1"}},
			wantToUpdate: nil,
			wantToDelete: []rec{{"d", "4"}},
		},
		{
			name: "update",
			args: args[rec]{
				plan:  []rec{{"a", "1"}, {"b", "2"}, {"c", "3"}},
				state: []rec{{"a", "1"}, {"b", "2"}, {"c", "4"}},
			},
			wantToCreate: nil,
			wantToUpdate: []rec{{"c", "3"}},
			wantToDelete: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotToCreate, gotToUpdate, gotToDelete := SetDiff(
				tt.args.plan, tt.args.state, func(t rec) string { return t.n },
			)
			assert.Equalf(t, tt.wantToCreate, gotToCreate, "toCreate is different")
			assert.Equalf(t, tt.wantToUpdate, gotToUpdate, "toUpdate is different")
			assert.Equalf(t, tt.wantToDelete, gotToDelete, "toDelete is different")
		})
	}
}
