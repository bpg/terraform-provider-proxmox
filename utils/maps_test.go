/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedListFromMap(t *testing.T) {
	t.Parallel()

	inputMap := map[string]interface{}{
		"value1": map[string]interface{}{"name": "resource1", "attr": "value1"},
		"value3": map[string]interface{}{"name": "resource3", "attr": "value3"},
		"value2": map[string]interface{}{"name": "resource2", "attr": "value2"},
	}

	expected := []interface{}{
		map[string]interface{}{"name": "resource1", "attr": "value1"},
		map[string]interface{}{"name": "resource2", "attr": "value2"},
		map[string]interface{}{"name": "resource3", "attr": "value3"},
	}

	result := OrderedListFromMap(inputMap)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ListResourcesAttributeValue() = %v, want %v", result, expected)
	}
}

func TestMapResourceList(t *testing.T) {
	t.Parallel()

	resourceList := []interface{}{
		map[string]interface{}{"name": "resource1", "attr": "value1"},
		map[string]interface{}{"name": "resource2", "attr": "value2"},
		nil,
		map[string]interface{}{"name": "resource3", "attr": "value3"},
		map[string]interface{}{"name": "resource4", "attr": "value4"},
	}

	expected := []string{
		"value1",
		"value2",
		"value3",
		"value4",
	}

	result := ListResourcesAttributeValue(resourceList, "attr")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ListResourcesAttributeValue() = %v, want %v", result, expected)
	}
}

func TestOrderedListFromMapByKeyValues(t *testing.T) {
	t.Parallel()

	inputMap := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
	}

	keyList := []string{"key2", "key1", "key4"}

	expected := []interface{}{"value2", "value1", "value4"}

	result := OrderedListFromMapByKeyValues(inputMap, keyList)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("OrderedListFromMapByKeyValues() = %v, want %v", result, expected)
	}
}

func TestCompareWithPrefix(t *testing.T) {
	t.Parallel()

	type args struct {
		a string
		b string
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{"equal", args{"a", "a"}, 0},
		{"a < b", args{"a", "b"}, -1},
		{"b > a", args{"b", "a"}, 1},
		{"a < b with prefix", args{"a1", "a2"}, -1},
		{"b > a with prefix", args{"a2", "a1"}, 1},
		{"a < b with different prefix", args{"a1", "b1"}, -1},
		{"b > a with different prefix", args{"b1", "a1"}, 1},
		{"a < b with different prefix and numbers", args{"a1", "a10"}, -1},
		{"b > a with different prefix and numbers", args{"a10", "a1"}, 1},
		{"a < b with different prefix and numbers", args{"a10", "b1"}, -1},
		{"b > a with different prefix and numbers", args{"b1", "a10"}, 1},
		{"a < b with different prefix and numbers", args{"a10", "b10"}, -1},
		{"b > a with different prefix and numbers", args{"b10", "a10"}, 1},
		{"b > a with different numbers leading zeros", args{"dev01", "dev001"}, 1},
		{"a > b with different numbers leading zeros", args{"dev001", "dev01"}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equalf(t, tt.want, compareWithPrefix(tt.args.a, tt.args.b), "compareWithPrefix(%v, %v)", tt.args.a, tt.args.b)
		})
	}
}
