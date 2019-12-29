/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func getContentTypeValidator() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"backup",
		"iso",
		"snippets",
		"vztmpl",
	}, false)
}

func getFileFormatValidator() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"qcow2",
		"raw",
		"vmdk",
	}, false)
}

func getFileIDValidator() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (ws []string, es []error) {
		v, ok := i.(string)

		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if v != "" {
			r := regexp.MustCompile(`^(?i)[a-z0-9\-_]+:([a-z0-9\-_]+/)?.+$`)
			ok := r.MatchString(v)

			if !ok {
				es = append(es, fmt.Errorf("expected %s to be a valid file identifier (datastore-name:iso/some-file.img), got %s", k, v))
				return
			}
		}

		return
	}
}

func getKeyboardLayoutValidator() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"da",
		"de",
		"de-ch",
		"en-gb",
		"en-us",
		"es",
		"fi",
		"fr",
		"fr-be",
		"fr-ca",
		"fr-ch",
		"hu",
		"is",
		"it",
		"ja",
		"lt",
		"mk",
		"nl",
		"no",
		"pl",
		"pt",
		"pt-br",
		"sl",
		"sv",
		"tr",
	}, false)
}

func getMACAddressValidator() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (ws []string, es []error) {
		v, ok := i.(string)

		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if v != "" {
			r := regexp.MustCompile(`^[A-Z0-9]{2}(:[A-Z0-9]{2}){5}$`)
			ok := r.MatchString(v)

			if !ok {
				es = append(es, fmt.Errorf("expected %s to be a valid MAC address (A0:B1:C2:D3:E4:F5), got %s", k, v))
				return
			}
		}

		return
	}
}

func getNetworkDeviceModelValidator() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{"e1000", "rtl8139", "virtio", "vmxnet3"}, false)
}

func getOSTypeValidator() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"l24",
		"l26",
		"other",
		"solaris",
		"w2k",
		"w2k3",
		"w2k8",
		"win7",
		"win8",
		"win10",
		"wvista",
		"wxp",
	}, false)
}

func getQEMUAgentTypeValidator() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{"isa", "virtio"}, false)
}

func getSchemaBlock(r *schema.Resource, d *schema.ResourceData, m interface{}, k []string, i int, allowDefault bool) (map[string]interface{}, error) {
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
				return resourceBlock, fmt.Errorf("Index out of bounds %d", i)
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
				return nil, err
			}

			list = listDefault.([]interface{})
		}
	}

	if len(list) > i {
		resourceBlock = list[i].(map[string]interface{})
	}

	return resourceBlock, nil
}

func getVLANIDsValidator() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (ws []string, es []error) {
		min := 1
		max := 4094

		list, ok := i.([]interface{})

		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be []interface{}", k))
			return
		}

		for li, lv := range list {
			v, ok := lv.(int)

			if !ok {
				es = append(es, fmt.Errorf("expected type of %s[%d] to be int", k, li))
				return
			}

			if v != -1 {
				if v < min || v > max {
					es = append(es, fmt.Errorf("expected %s[%d] to be in the range (%d - %d), got %d", k, li, min, max, v))
					return
				}
			}
		}

		return
	}
}

func getVMIDValidator() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (ws []string, es []error) {
		min := 1
		max := 2147483647

		v, ok := i.(int)

		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be int", k))
			return
		}

		if v != -1 {
			if v < min || v > max {
				es = append(es, fmt.Errorf("expected %s to be in the range (%d - %d), got %d", k, min, max, v))
				return
			}
		}

		return
	}
}

func testComputedAttributes(t *testing.T, s *schema.Resource, keys []string) {
	for _, v := range keys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in Schema: Missing definition for \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in Schema: Attribute \"%s\" is not computed", v)
		}
	}
}

func testNestedSchemaExistence(t *testing.T, s *schema.Resource, key string) *schema.Resource {
	schema, ok := s.Schema[key].Elem.(*schema.Resource)

	if !ok {
		t.Fatalf("Error in Schema: Missing nested schema for \"%s\"", key)

		return nil
	}

	return schema
}

func testOptionalArguments(t *testing.T, s *schema.Resource, keys []string) {
	for _, v := range keys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in Schema: Missing definition for \"%s\"", v)
		}

		if s.Schema[v].Optional != true {
			t.Fatalf("Error in Schema: Argument \"%s\" is not optional", v)
		}
	}
}

func testRequiredArguments(t *testing.T, s *schema.Resource, keys []string) {
	for _, v := range keys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in Schema: Missing definition for \"%s\"", v)
		}

		if s.Schema[v].Required != true {
			t.Fatalf("Error in Schema: Argument \"%s\" is not required", v)
		}
	}
}

func testSchemaValueTypes(t *testing.T, s *schema.Resource, keys []string, types []schema.ValueType) {
	for i, v := range keys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in Schema: Missing definition for \"%s\"", v)
		}

		if s.Schema[v].Type != types[i] {
			t.Fatalf("Error in Schema: Argument \"%s\" is not of type \"%v\"", v, types[i])
		}
	}
}
