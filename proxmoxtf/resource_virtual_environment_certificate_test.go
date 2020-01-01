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

	testComputedAttributes(t, s, []string{
		mkResourceVirtualEnvironmentCertificateExpirationDate,
		mkResourceVirtualEnvironmentCertificateFileName,
		mkResourceVirtualEnvironmentCertificateIssuer,
		mkResourceVirtualEnvironmentCertificatePublicKeySize,
		mkResourceVirtualEnvironmentCertificatePublicKeyType,
		mkResourceVirtualEnvironmentCertificateSSLFingerprint,
		mkResourceVirtualEnvironmentCertificateStartDate,
		mkResourceVirtualEnvironmentCertificateSubject,
		mkResourceVirtualEnvironmentCertificateSubjectAlternativeNames,
	})

	testSchemaValueTypes(t, s, []string{
		mkResourceVirtualEnvironmentCertificateCertificate,
		mkResourceVirtualEnvironmentCertificateCertificateChain,
		mkResourceVirtualEnvironmentCertificateExpirationDate,
		mkResourceVirtualEnvironmentCertificateFileName,
		mkResourceVirtualEnvironmentCertificateIssuer,
		mkResourceVirtualEnvironmentCertificateNodeName,
		mkResourceVirtualEnvironmentCertificatePrivateKey,
		mkResourceVirtualEnvironmentCertificatePublicKeySize,
		mkResourceVirtualEnvironmentCertificatePublicKeyType,
		mkResourceVirtualEnvironmentCertificateSSLFingerprint,
		mkResourceVirtualEnvironmentCertificateStartDate,
		mkResourceVirtualEnvironmentCertificateSubject,
		mkResourceVirtualEnvironmentCertificateSubjectAlternativeNames,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeString,
		schema.TypeString,
		schema.TypeString,
		schema.TypeString,
		schema.TypeString,
		schema.TypeString,
		schema.TypeInt,
		schema.TypeString,
		schema.TypeString,
		schema.TypeString,
		schema.TypeString,
		schema.TypeList,
	})
}
