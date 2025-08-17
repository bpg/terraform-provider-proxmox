//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDatasourceVersion(t *testing.T) {
	te := InitEnvironment(t)

	datasourceName := "data.proxmox_virtual_environment_version.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: `data "proxmox_virtual_environment_version" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
					resource.TestCheckResourceAttrSet(datasourceName, "repository_id"),
					resource.TestCheckResourceAttrWith(datasourceName, "release", validateReleaseVersion),
					resource.TestCheckResourceAttrWith(datasourceName, "version", validateFullVersion),
					validateVersionReleaseConsistency(datasourceName),
				),
			},
		},
	})
}

// validateReleaseVersion validates that the release field matches the expected pattern (e.g., "8.4", "9.0").
func validateReleaseVersion(value string) error {
	releasePattern := regexp.MustCompile(`^[0-9]+\.[0-9]+$`)
	if !releasePattern.MatchString(value) {
		return fmt.Errorf("release %q does not match expected pattern (major.minor)", value)
	}

	// Ensure it's at least the minimum supported version
	releaseVer, err := version.NewVersion(value + ".0") // Add patch version for comparison
	if err != nil {
		return fmt.Errorf("failed to parse release version %q: %w", value, err)
	}

	minVersion := version.Must(version.NewVersion("8.0.0"))
	if releaseVer.LessThan(minVersion) {
		return fmt.Errorf("release version %q is below minimum supported version 8.0", value)
	}

	return nil
}

// validateFullVersion validates that the version field is a valid semantic version.
func validateFullVersion(value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("version cannot be empty")
	}

	_, err := version.NewVersion(value)
	if err != nil {
		return fmt.Errorf("version %q is not a valid semantic version: %w", value, err)
	}

	return nil
}

// validateVersionReleaseConsistency returns a TestCheckFunc that validates version and release consistency.
func validateVersionReleaseConsistency(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}

		release := rs.Primary.Attributes["release"]
		version := rs.Primary.Attributes["version"]

		if release == "" {
			return fmt.Errorf("release attribute is empty")
		}
		if version == "" {
			return fmt.Errorf("version attribute is empty")
		}

		if !strings.HasPrefix(version, release) {
			return fmt.Errorf("version %q does not start with release %q", version, release)
		}

		return nil
	}
}
