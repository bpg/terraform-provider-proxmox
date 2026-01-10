/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
)

func ptr[T any](v T) *T {
	return &v
}

func createDomainsListValue(t *testing.T, domains []acmeDomainModel) types.List {
	t.Helper()

	if len(domains) == 0 {
		return types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"domain": types.StringType,
				"plugin": types.StringType,
				"alias":  types.StringType,
			},
		})
	}

	list, diag := types.ListValueFrom(context.Background(), types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"domain": types.StringType,
			"plugin": types.StringType,
			"alias":  types.StringType,
		},
	}, domains)
	require.False(t, diag.HasError(), "failed to create domains list: %v", diag.Errors())

	return list
}

func TestFindMatchingCertificate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		domains      []acmeDomainModel
		certificates []nodes.CertificateListResponseData
		wantIssuer   string
		wantErr      bool
	}{
		{
			name: "match certificate by SAN",
			domains: []acmeDomainModel{
				{Domain: types.StringValue("example.com"), Plugin: types.StringNull(), Alias: types.StringNull()},
			},
			certificates: []nodes.CertificateListResponseData{
				{
					Issuer:                  ptr("CN=Let's Encrypt"),
					Subject:                 ptr("CN=example.com"),
					SubjectAlternativeNames: ptr([]string{"example.com"}),
				},
			},
			wantIssuer: "CN=Let's Encrypt",
			wantErr:    false,
		},
		{
			name: "match certificate by CN when SANs empty",
			domains: []acmeDomainModel{
				{Domain: types.StringValue("example.com"), Plugin: types.StringNull(), Alias: types.StringNull()},
			},
			certificates: []nodes.CertificateListResponseData{
				{
					Issuer:  ptr("CN=Let's Encrypt"),
					Subject: ptr("CN=example.com"),
				},
			},
			wantIssuer: "CN=Let's Encrypt",
			wantErr:    false,
		},
		{
			name: "prefer ACME cert over Proxmox-generated",
			domains: []acmeDomainModel{
				{Domain: types.StringValue("pve.example.com"), Plugin: types.StringNull(), Alias: types.StringNull()},
			},
			certificates: []nodes.CertificateListResponseData{
				{
					Issuer:                  ptr("CN=Proxmox Virtual Environment"),
					Subject:                 ptr("CN=pve.example.com"),
					SubjectAlternativeNames: ptr([]string{"pve.example.com"}),
				},
				{
					Issuer:                  ptr("CN=Let's Encrypt"),
					Subject:                 ptr("CN=pve.example.com"),
					SubjectAlternativeNames: ptr([]string{"pve.example.com"}),
				},
			},
			wantIssuer: "CN=Let's Encrypt",
			wantErr:    false,
		},
		{
			name: "prefer cert with more matching domains",
			domains: []acmeDomainModel{
				{Domain: types.StringValue("example.com"), Plugin: types.StringNull(), Alias: types.StringNull()},
				{Domain: types.StringValue("www.example.com"), Plugin: types.StringNull(), Alias: types.StringNull()},
			},
			certificates: []nodes.CertificateListResponseData{
				{
					Issuer:                  ptr("CN=Partial Match"),
					Subject:                 ptr("CN=example.com"),
					SubjectAlternativeNames: ptr([]string{"example.com"}),
				},
				{
					Issuer:                  ptr("CN=Full Match"),
					Subject:                 ptr("CN=example.com"),
					SubjectAlternativeNames: ptr([]string{"example.com", "www.example.com"}),
				},
			},
			wantIssuer: "CN=Full Match",
			wantErr:    false,
		},
		{
			name: "no matching domains - prefer ACME over Proxmox",
			domains: []acmeDomainModel{
				{Domain: types.StringValue("other.com"), Plugin: types.StringNull(), Alias: types.StringNull()},
			},
			certificates: []nodes.CertificateListResponseData{
				{
					Issuer:                  ptr("CN=Proxmox Virtual Environment"),
					Subject:                 ptr("CN=pve.local"),
					SubjectAlternativeNames: ptr([]string{"pve.local"}),
				},
				{
					Issuer:                  ptr("CN=Let's Encrypt"),
					Subject:                 ptr("CN=example.com"),
					SubjectAlternativeNames: ptr([]string{"example.com"}),
				},
			},
			wantIssuer: "CN=Let's Encrypt",
			wantErr:    false,
		},
		{
			name:         "empty certificates list",
			domains:      []acmeDomainModel{},
			certificates: []nodes.CertificateListResponseData{},
			wantErr:      true,
		},
		{
			name:    "fallback to first cert when all are Proxmox-generated",
			domains: []acmeDomainModel{},
			certificates: []nodes.CertificateListResponseData{
				{
					Issuer:  ptr("CN=Proxmox Virtual Environment"),
					Subject: ptr("CN=pve.local"),
				},
				{
					Issuer:  ptr("CN=PVE Cluster Node"),
					Subject: ptr("CN=pve2.local"),
				},
			},
			wantIssuer: "CN=Proxmox Virtual Environment",
			wantErr:    false,
		},
		{
			name: "match with multiple SANs",
			domains: []acmeDomainModel{
				{Domain: types.StringValue("api.example.com"), Plugin: types.StringNull(), Alias: types.StringNull()},
			},
			certificates: []nodes.CertificateListResponseData{
				{
					Issuer:  ptr("CN=Let's Encrypt"),
					Subject: ptr("CN=example.com"),
					SubjectAlternativeNames: ptr([]string{
						"example.com",
						"www.example.com",
						"api.example.com",
						"mail.example.com",
					}),
				},
			},
			wantIssuer: "CN=Let's Encrypt",
			wantErr:    false,
		},
		{
			name: "match CN with X.509 subject format (slash delimiter)",
			domains: []acmeDomainModel{
				{Domain: types.StringValue("pve.bpghome.net"), Plugin: types.StringNull(), Alias: types.StringNull()},
			},
			certificates: []nodes.CertificateListResponseData{
				{
					Issuer:  ptr("/C=US/O=Let's Encrypt/CN=R12"),
					Subject: ptr("/CN=pve.bpghome.net"),
				},
			},
			wantIssuer: "/C=US/O=Let's Encrypt/CN=R12",
			wantErr:    false,
		},
		{
			name: "match CN with complex X.509 subject format",
			domains: []acmeDomainModel{
				{Domain: types.StringValue("pve.bpglabs.net"), Plugin: types.StringNull(), Alias: types.StringNull()},
			},
			certificates: []nodes.CertificateListResponseData{
				{
					Issuer:  ptr("/CN=Proxmox Virtual Environment"),
					Subject: ptr("/OU=PVE Cluster Node/O=Proxmox Virtual Environment/CN=pve.bpglabs.net"),
				},
			},
			wantIssuer: "/CN=Proxmox Virtual Environment",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &acmeCertificateResource{}
			model := &acmeCertificateModel{
				Domains: createDomainsListValue(t, tt.domains),
			}

			certs := tt.certificates
			result, err := r.findMatchingCertificate(context.Background(), model, &certs)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)

			if tt.wantIssuer != "" {
				assert.Equal(t, tt.wantIssuer, *result.Issuer)
			}
		})
	}
}

func TestIsProxmoxGeneratedCertificate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		issuer *string
		want   bool
	}{
		{
			name:   "Proxmox in issuer",
			issuer: ptr("CN=Proxmox Virtual Environment"),
			want:   true,
		},
		{
			name:   "PVE in issuer",
			issuer: ptr("CN=PVE Cluster Node"),
			want:   true,
		},
		{
			name:   "Let's Encrypt issuer",
			issuer: ptr("CN=Let's Encrypt"),
			want:   false,
		},
		{
			name:   "nil issuer",
			issuer: nil,
			want:   false,
		},
		{
			name:   "empty issuer",
			issuer: ptr(""),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cert := &nodes.CertificateListResponseData{
				Issuer: tt.issuer,
			}

			got := isProxmoxGeneratedCertificate(cert)
			assert.Equal(t, tt.want, got)
		})
	}
}
