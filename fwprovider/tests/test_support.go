/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testResourceAttributes(res string, attrs map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for k, v := range attrs {
			if v == "" {
				if err := resource.TestCheckResourceAttr(res, k, "")(s); err != nil {
					return fmt.Errorf("expected '%s' to be empty: %w", k, err)
				}

				continue
			}

			if err := resource.TestCheckResourceAttrWith(res, k, func(got string) error {
				match, err := regexp.Match(v, []byte(got)) //nolint:mirror
				if err != nil {
					return fmt.Errorf("error matching '%s': %w", v, err)
				}
				if !match {
					return fmt.Errorf("expected '%s' to match '%s'", got, v)
				}
				return nil
			})(s); err != nil {
				return err
			}
		}

		return nil
	}
}

func testNoResourceAttributesSet(res string, attrs []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, k := range attrs {
			if err := resource.TestCheckNoResourceAttr(res, k)(s); err != nil {
				return err
			}
		}

		return nil
	}
}

func testResourceAttributesSet(res string, attrs []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, k := range attrs {
			if err := resource.TestCheckResourceAttrSet(res, k)(s); err != nil {
				return err
			}
		}

		return nil
	}
}
