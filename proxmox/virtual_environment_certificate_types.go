/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

// VirtualEnvironmentCertificateDeleteRequestBody contains the data for a custom certificate delete request.
type VirtualEnvironmentCertificateDeleteRequestBody struct {
	Restart *CustomBool `json:"restart,omitempty" url:"restart,omitempty,int"`
}

// VirtualEnvironmentCertificateListResponseBody contains the body from a certificate list response.
type VirtualEnvironmentCertificateListResponseBody struct {
	Data *[]VirtualEnvironmentCertificateListResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentCertificateListResponseData contains the data from a certificate list response.
type VirtualEnvironmentCertificateListResponseData struct {
	Certificates            *string          `json:"pem,omitempty"`
	FileName                *string          `json:"filename,omitempty"`
	Fingerprint             *string          `json:"fingerprint,omitempty"`
	Issuer                  *string          `json:"issuer,omitempty"`
	NotAfter                *CustomTimestamp `json:"notafter,omitempty"`
	NotBefore               *CustomTimestamp `json:"notbefore,omitempty"`
	PublicKeyBits           *int             `json:"public-key-bits,omitempty"`
	PublicKeyType           *string          `json:"public-key-type,omitempty"`
	Subject                 *string          `json:"subject,omitempty"`
	SubjectAlternativeNames *[]string        `json:"san,omitempty"`
}

// VirtualEnvironmentCertificateUpdateRequestBody contains the body for a custom certificate update request.
type VirtualEnvironmentCertificateUpdateRequestBody struct {
	Certificates string      `json:"certificates" url:"certificates"`
	Force        *CustomBool `json:"force,omitempty" url:"force,omitempty,int"`
	PrivateKey   *string     `json:"key,omitempty" url:"key,omitempty"`
	Restart      *CustomBool `json:"restart,omitempty" url:"restart,omitempty,int"`
}
