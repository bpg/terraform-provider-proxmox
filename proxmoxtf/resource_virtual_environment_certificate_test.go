/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

// TestResourceVirtualEnvironmentCertificateInstantiation tests whether the ResourceVirtualEnvironmentCertificate instance can be instantiated.
func TestResourceVirtualEnvironmentCertificateInstantiation(t *testing.T) {
	s := resourceVirtualEnvironmentCertificate()

	if s == nil {
		t.Fatalf("Cannot instantiate resourceVirtualEnvironmentCertificate")
	}
}

// TestResourceVirtualEnvironmentCertificateSchema tests the resourceVirtualEnvironmentCertificate schema.
func TestResourceVirtualEnvironmentCertificateSchema(t *testing.T) {
	s := resourceVirtualEnvironmentCertificate()

	testRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentCertificateCertificate,
		mkResourceVirtualEnvironmentCertificateNodeName,
		mkResourceVirtualEnvironmentCertificatePrivateKey,
	})

	testOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentCertificateCertificateChain,
	})

	testSchemaValueTypes(t, s, []string{
		mkResourceVirtualEnvironmentCertificateCertificate,
		mkResourceVirtualEnvironmentCertificateCertificateChain,
		mkResourceVirtualEnvironmentCertificateNodeName,
		mkResourceVirtualEnvironmentCertificatePrivateKey,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeString,
		schema.TypeString,
		schema.TypeString,
	})
}
