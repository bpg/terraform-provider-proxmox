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
