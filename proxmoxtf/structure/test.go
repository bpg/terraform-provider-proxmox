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

func AssertComputedAttributes(t *testing.T, s map[string]*schema.Schema, keys []string) {
	for _, v := range keys {
		require.NotNil(t, s[v], "Error in Schema: Missing definition for \"%s\"", v)
		assert.True(t, s[v].Computed, "Error in Schema: Attribute \"%s\" is not computed", v)
	}
}

func AssertNestedSchemaExistence(t *testing.T, s map[string]*schema.Schema, key string) *schema.Resource {
	sh, ok := s[key].Elem.(*schema.Resource)

	if !ok {
		t.Fatalf("Error in Schema: Missing nested schema for \"%s\"", key)

		return nil
	}

	return sh
}

func AssertOptionalArguments(t *testing.T, s map[string]*schema.Schema, keys []string) {
	for _, v := range keys {
		require.NotNil(t, s[v], "Error in Schema: Missing definition for \"%s\"", v)
		assert.True(t, s[v].Optional, "Error in Schema: Argument \"%s\" is not optional", v)
	}
}

func AssertRequiredArguments(t *testing.T, s map[string]*schema.Schema, keys []string) {
	for _, v := range keys {
		require.NotNil(t, s[v], "Error in Schema: Missing definition for \"%s\"", v)
		assert.True(t, s[v].Required, "Error in Schema: Argument \"%s\" is not required", v)
	}
}

func AssertValueTypes(t *testing.T, s map[string]*schema.Schema, f map[string]schema.ValueType) {
	for fn, ft := range f {
		require.NotNil(t, s[fn], "Error in Schema: Missing definition for \"%s\"", fn)
		assert.Equal(t, ft, s[fn].Type, "Error in Schema: Argument or attribute \"%s\" is not of type \"%v\"", fn, ft)
	}
}
