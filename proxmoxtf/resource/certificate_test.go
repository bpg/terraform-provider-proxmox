/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// TestCertificateInstantiation tests whether the Certificate instance can be instantiated.
func TestCertificateInstantiation(t *testing.T) {
	t.Parallel()

	s := Certificate()
	if s == nil {
		t.Fatalf("Cannot instantiate Certificate")
	}
}

// TestCertificateSchema tests the Certificate schema.
func TestCertificateSchema(t *testing.T) {
	t.Parallel()

	s := Certificate()

	test.AssertRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentCertificateCertificate,
		mkResourceVirtualEnvironmentCertificateNodeName,
		mkResourceVirtualEnvironmentCertificatePrivateKey,
	})

	test.AssertOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentCertificateCertificateChain,
	})

	test.AssertComputedAttributes(t, s, []string{
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

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentCertificateCertificate:             schema.TypeString,
		mkResourceVirtualEnvironmentCertificateCertificateChain:        schema.TypeString,
		mkResourceVirtualEnvironmentCertificateExpirationDate:          schema.TypeString,
		mkResourceVirtualEnvironmentCertificateFileName:                schema.TypeString,
		mkResourceVirtualEnvironmentCertificateIssuer:                  schema.TypeString,
		mkResourceVirtualEnvironmentCertificateNodeName:                schema.TypeString,
		mkResourceVirtualEnvironmentCertificatePrivateKey:              schema.TypeString,
		mkResourceVirtualEnvironmentCertificatePublicKeySize:           schema.TypeInt,
		mkResourceVirtualEnvironmentCertificatePublicKeyType:           schema.TypeString,
		mkResourceVirtualEnvironmentCertificateSSLFingerprint:          schema.TypeString,
		mkResourceVirtualEnvironmentCertificateStartDate:               schema.TypeString,
		mkResourceVirtualEnvironmentCertificateSubject:                 schema.TypeString,
		mkResourceVirtualEnvironmentCertificateSubjectAlternativeNames: schema.TypeList,
	})
}
