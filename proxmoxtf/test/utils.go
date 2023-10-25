/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertComputedAttributes checks that the given schema has the given computed attributes.
func AssertComputedAttributes(t *testing.T, s *schema.Resource, keys []string) {
	t.Helper()

	for _, v := range keys {
		require.NotNil(t, s.Schema[v], "Error in Schema: Missing definition for \"%s\"", v)
		assert.True(t, s.Schema[v].Computed, "Error in Schema: Attribute \"%s\" is not computed", v)
	}
}

// AssertNestedSchemaExistence checks that the given schema has a nested schema for the given key.
func AssertNestedSchemaExistence(t *testing.T, s *schema.Resource, key string) *schema.Resource {
	t.Helper()

	sh, ok := s.Schema[key].Elem.(*schema.Resource)

	if !ok {
		t.Fatalf("Error in Schema: Missing nested schema for \"%s\"", key)

		return nil
	}

	return sh
}

// AssertListMaxItems checks that the given schema attribute has given expected MaxItems value.
func AssertListMaxItems(t *testing.T, s *schema.Resource, key string, expectedMaxItems int) {
	t.Helper()

	require.NotNil(t, s.Schema[key], "Error in Schema: Missing definition for \"%s\"", key)
	assert.Equal(t, expectedMaxItems, s.Schema[key].MaxItems,
		"Error in Schema: Argument \"%s\" has \"MaxItems: %#v\", but value %#v is expected!",
		key, s.Schema[key].MaxItems, expectedMaxItems)
}

// AssertOptionalArguments checks that the given schema has the given optional arguments.
func AssertOptionalArguments(t *testing.T, s *schema.Resource, keys []string) {
	t.Helper()

	for _, v := range keys {
		require.NotNil(t, s.Schema[v], "Error in Schema: Missing definition for \"%s\"", v)
		assert.True(t, s.Schema[v].Optional, "Error in Schema: Argument \"%s\" is not optional", v)
	}
}

// AssertRequiredArguments checks that the given schema has the given required arguments.
func AssertRequiredArguments(t *testing.T, s *schema.Resource, keys []string) {
	t.Helper()

	for _, v := range keys {
		require.NotNil(t, s.Schema[v], "Error in Schema: Missing definition for \"%s\"", v)
		assert.True(t, s.Schema[v].Required, "Error in Schema: Argument \"%s\" is not required", v)
	}
}

// AssertValueTypes checks that the given schema has the given value types for the given fields.
func AssertValueTypes(t *testing.T, s *schema.Resource, f map[string]schema.ValueType) {
	t.Helper()

	for fn, ft := range f {
		require.NotNil(t, s.Schema[fn], "Error in Schema: Missing definition for \"%s\"", fn)
		assert.Equal(t, ft, s.Schema[fn].Type, "Error in Schema: Argument or attribute \"%s\" is not of type \"%v\"", fn, ft)
	}
}
