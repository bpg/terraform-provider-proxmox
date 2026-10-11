/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
)

// SafeResourceName generates a Proxmox-safe resource identifier.
// It returns "{prefix}-{8 random lowercase letters}", e.g. "test-pool-abcdefgh".
// This avoids gofakeit.Word() which can produce dots, hyphens, and other
// characters that violate Proxmox identifier validation rules.
func SafeResourceName(prefix string) string {
	return fmt.Sprintf("%s-%s", prefix, strings.ToLower(gofakeit.LetterN(8)))
}

// ResourceAttributes is a helper function to test resource attributes.
func ResourceAttributes(res string, attrs map[string]string) resource.TestCheckFunc {
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

// NoResourceAttributesSet is a helper function to test that no resource attributes are set.
func NoResourceAttributesSet(res string, attrs []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, k := range attrs {
			if err := resource.TestCheckNoResourceAttr(res, k)(s); err != nil {
				return err
			}
		}

		return nil
	}
}

// ResourceAttributesSet is a helper function to test that all resource attributes are set.
func ResourceAttributesSet(res string, attrs []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, k := range attrs {
			if err := resource.TestCheckResourceAttrSet(res, k)(s); err != nil {
				return err
			}
		}

		return nil
	}
}

// onceCephStatus/errCephStatus cache the cluster-wide Ceph status probe so
// tests calling RequireCeph share one HTTP call per `go test` invocation per
// package.
//
//nolint:gochecknoglobals
var (
	onceCephStatus sync.Once
	errCephStatus  error
)

// RequireCeph skips the calling test when Ceph is not initialized on the
// testacc cluster. Probes the cluster-wide `/cluster/ceph/status` endpoint,
// which returns an error when Ceph is uninstalled or uninitialized — making
// it a cheap, cluster-scoped probe. The probe result is cached for the
// lifetime of the test binary; the first caller pays the HTTP cost,
// subsequent callers reuse the cached outcome.
//
// Default `./testacc` runs include every acceptance test in `fwprovider/...`
// regardless of `--tier` or `--resource` filters, so tests tagged
// `//testacc:resource=ceph` must self-skip on clusters without Ceph or
// they will fail the maintainer's release-cycle runs with API 500s.
func (e *Environment) RequireCeph() {
	e.t.Helper()

	onceCephStatus.Do(func() {
		// Use a fresh background context with a short timeout so the cached
		// result doesn't get bound to the first caller's test lifetime.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		_, errCephStatus = e.ClusterClient().Ceph().GetStatus(ctx)
	})

	if errCephStatus != nil {
		e.t.Skipf("Skipping: Ceph is not available on the testacc cluster (%v)", errCephStatus)
	}
}

func CreateTempFile(t *testing.T, namePattern string, content string) *os.File {
	t.Helper()

	f, err := os.CreateTemp(t.TempDir(), namePattern)
	require.NoError(t, err)

	_, err = f.WriteString(content)
	require.NoError(t, err)

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	t.Cleanup(func() {
		_ = os.Remove(f.Name())
	})

	return f
}
