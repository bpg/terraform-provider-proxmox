/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package structure

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// MergeSchema merges the map[string]*schema.Schema from src into dst.
// Panicking enforces safety against conflicts.
func MergeSchema(dst, src map[string]*schema.Schema) {
	for k, v := range src {
		if _, ok := dst[k]; ok {
			panic(fmt.Errorf("conflicting schema key: %s", k))
		}

		dst[k] = v
	}
}

// GetSchemaBlock returns a map[string]interface{} of a nested resource by key(s) from a schema.ResourceData.
func GetSchemaBlock(
	r *schema.Resource,
	d *schema.ResourceData,
	k []string,
	i int,
	allowDefault bool,
) (map[string]interface{}, error) {
	var resourceBlock map[string]interface{}

	var resourceData interface{}

	var resourceSchema *schema.Schema

	for ki, kv := range k {
		if ki == 0 {
			resourceData = d.Get(kv)
			resourceSchema = r.Schema[kv]
		} else {
			mapValues := resourceData.([]interface{})

			if len(mapValues) <= i {
				return resourceBlock, fmt.Errorf("index out of bounds %d", i)
			}

			mapValue := mapValues[i].(map[string]interface{})

			resourceData = mapValue[kv]
			resourceSchema = resourceSchema.Elem.(*schema.Resource).Schema[kv]
		}
	}

	list := resourceData.([]interface{})

	if len(list) == 0 {
		if allowDefault {
			listDefault, err := resourceSchema.DefaultValue()
			if err != nil {
				return nil, fmt.Errorf("failed to get default value for %s: %w", strings.Join(k, "."), err)
			}

			list = listDefault.([]interface{})
		}
	}

	if len(list) > i {
		resourceBlock = list[i].(map[string]interface{})
	}

	return resourceBlock, nil
}

// SuppressIfListsAreEqualIgnoringOrder is a customdiff.SuppressionFunc that suppresses
// changes to a list if the old and new lists are equal, ignoring the order of the
// elements.
// It will be called for each list item, so it is not super efficient. It is
// recommended to use it only for small lists.
// Ref: https://github.com/hashicorp/terraform-plugin-sdk/issues/477
func SuppressIfListsAreEqualIgnoringOrder(key, _, _ string, d *schema.ResourceData) bool {
	// the key is a path to the list item, not the list itself, e.g. "tags.0"
	lastDotIndex := strings.LastIndex(key, ".")
	if lastDotIndex != -1 {
		key = key[:lastDotIndex]
	}

	oldData, newData := d.GetChange(key)
	if oldData == nil || newData == nil {
		return false
	}

	oldArray := oldData.([]interface{})
	newArray := newData.([]interface{})

	if len(oldArray) != len(newArray) {
		return false
	}

	oldEvents := make([]string, len(oldArray))
	newEvents := make([]string, len(newArray))

	for i, oldEvt := range oldArray {
		oldEvents[i] = fmt.Sprint(oldEvt)
	}

	for j, newEvt := range newArray {
		newEvents[j] = fmt.Sprint(newEvt)
	}

	sort.Strings(oldEvents)
	sort.Strings(newEvents)

	return reflect.DeepEqual(oldEvents, newEvents)
}
