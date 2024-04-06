package utils

import (
	"reflect"
	"testing"
)

func TestMapResourceList(t *testing.T) {
	t.Parallel()

	resourceList := []interface{}{
		map[string]interface{}{"name": "resource1", "attr": "value1"},
		map[string]interface{}{"name": "resource2", "attr": "value2"},
		map[string]interface{}{"name": "resource3", "attr": "value3"},
	}

	expected := map[string]interface{}{
		"value1": map[string]interface{}{"name": "resource1", "attr": "value1"},
		"value2": map[string]interface{}{"name": "resource2", "attr": "value2"},
		"value3": map[string]interface{}{"name": "resource3", "attr": "value3"},
	}

	result := MapResourceList(resourceList, "attr")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("MapResourceList() = %v, want %v", result, expected)
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
