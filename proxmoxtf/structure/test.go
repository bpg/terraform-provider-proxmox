/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package structure

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertComputedAttributes asserts that the given keys are present in the schema and are computed.
func AssertComputedAttributes(t *testing.T, s map[string]*schema.Schema, keys []string) {
	t.Helper()

	for _, v := range keys {
		require.NotNil(t, s[v], "Error in Schema: Missing definition for \"%s\"", v)
		assert.True(t, s[v].Computed, "Error in Schema: Attribute \"%s\" is not computed", v)
	}
}

// AssertNestedSchemaExistence asserts that the given key is present in the schema and is a nested schema.
func AssertNestedSchemaExistence(t *testing.T, s map[string]*schema.Schema, key string) *schema.Resource {
	t.Helper()

	sh, ok := s[key].Elem.(*schema.Resource)

	if !ok {
		t.Fatalf("Error in Schema: Missing nested schema for \"%s\"", key)

		return nil
	}

	return sh
}

// AssertOptionalArguments asserts that the given keys are present in the schema and are optional.
func AssertOptionalArguments(t *testing.T, s map[string]*schema.Schema, keys []string) {
	t.Helper()

	for _, v := range keys {
		require.NotNil(t, s[v], "Error in Schema: Missing definition for \"%s\"", v)
		assert.True(t, s[v].Optional, "Error in Schema: Argument \"%s\" is not optional", v)
	}
}

// AssertRequiredArguments asserts that the given keys are present in the schema and are required.
func AssertRequiredArguments(t *testing.T, s map[string]*schema.Schema, keys []string) {
	t.Helper()

	for _, v := range keys {
		require.NotNil(t, s[v], "Error in Schema: Missing definition for \"%s\"", v)
		assert.True(t, s[v].Required, "Error in Schema: Argument \"%s\" is not required", v)
	}
}

// AssertValueTypes asserts that the given keys are present in the schema and are of the given type.
func AssertValueTypes(t *testing.T, s map[string]*schema.Schema, f map[string]schema.ValueType) {
	t.Helper()

	for fn, ft := range f {
		require.NotNil(t, s[fn], "Error in Schema: Missing definition for \"%s\"", fn)
		assert.Equal(t, ft, s[fn].Type, "Error in Schema: Argument or attribute \"%s\" is not of type \"%v\"", fn, ft)
	}
}
