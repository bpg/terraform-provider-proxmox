package utils

import "sort"

// OrderedListFromMap generates a list from a map's values. The values are sorted based on the map's keys.
func OrderedListFromMap(inputMap map[string]interface{}) []interface{} {
	itemCount := len(inputMap)
	keyList := make([]string, itemCount)
	i := 0

	for key := range inputMap {
		keyList[i] = key
		i++
	}

	sort.Strings(keyList)

	orderedList := make([]interface{}, itemCount)
	for i, k := range keyList {
		orderedList[i] = inputMap[k]
	}

	return orderedList
}

// MapResourceList generates a list of strings from a Terraform resource list (list of maps).
// The list is generated from the value of the specified attribute.
//
// "Map" in this context is a functional programming term, not a Go map.
// "Resource" in this context is a Terraform resource, i.e. a map of attributes.
func MapResourceList(resourceList []interface{}, attrName string) map[string]interface{} {
	m := make(map[string]interface{}, len(resourceList))

	for _, resource := range resourceList {
		r := resource.(map[string]interface{})
		key := r[attrName].(string)
		m[key] = r
	}

	return m
}

// OrderedListFromMapByKeyValues generates a list from a map's values.
// The values are sorted based on the provided key list. If a key is not found in the map, it is skipped.
func OrderedListFromMapByKeyValues(inputMap map[string]interface{}, keyList []string) []interface{} {
	orderedList := make([]interface{}, len(keyList))

	for i, k := range keyList {
		val, ok := inputMap[k]
		if ok {
			orderedList[i] = val
		}
	}

	return orderedList
}
