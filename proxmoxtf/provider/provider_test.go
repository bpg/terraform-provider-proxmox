/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// TestProviderInstantiation() tests whether the ProxmoxVirtualEnvironment instance can be instantiated.
func TestProviderInstantiation(t *testing.T) {
	t.Parallel()

	s := ProxmoxVirtualEnvironment()
	if s == nil {
		t.Fatalf("Cannot instantiate ProxmoxVirtualEnvironment")
	}
}

// TestProviderSchema() tests the ProxmoxVirtualEnvironment schema.
func TestProviderSchema(t *testing.T) {
	t.Parallel()

	s := ProxmoxVirtualEnvironment().Schema

	test.AssertOptionalArguments(t, s, []string{
		mkProviderEndpoint,
		mkProviderInsecure,
		mkProviderMinTLS,
		mkProviderAuthTicket,
		mkProviderCSRFPreventionToken,
		mkProviderOTP,
		mkProviderUsername,
		mkProviderPassword,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkProviderEndpoint:            schema.TypeString,
		mkProviderInsecure:            schema.TypeBool,
		mkProviderMinTLS:              schema.TypeString,
		mkProviderAuthTicket:          schema.TypeString,
		mkProviderCSRFPreventionToken: schema.TypeString,
		mkProviderOTP:                 schema.TypeString,
		mkProviderUsername:            schema.TypeString,
		mkProviderPassword:            schema.TypeString,
	})

	providerSSHSchema := test.AssertNestedSchemaExistence(t, s, mkProviderSSH)

	// do not limit number of nodes in the cluster
	test.AssertListMaxItems(t, providerSSHSchema, mkProviderSSHNode, 0)

	providerCloudflareAccessSchema := test.AssertNestedSchemaExistence(t, s, mkProviderCloudflareAccess)

	test.AssertOptionalArguments(t, providerCloudflareAccessSchema, []string{
		mkProviderCloudflareAccessClientID,
		mkProviderCloudflareAccessClientSecret,
	})
}

func TestCloudflareAccessConfig(t *testing.T) {
	tests := []struct {
		name       string
		envID      string
		envSecret  string
		block      []any
		wantConfig *api.CloudflareAccessConfig
		wantErr    bool
	}{
		{
			name: "block absent, no env",
		},
		{
			name:  "empty block, no env",
			block: []any{map[string]any{}},
		},
		{
			name:       "both set in block",
			block:      []any{map[string]any{"client_id": "block-id", "client_secret": "block-secret"}},
			wantConfig: &api.CloudflareAccessConfig{ClientID: "block-id", ClientSecret: "block-secret"},
		},
		{
			name:    "only client_id in block",
			block:   []any{map[string]any{"client_id": "block-id"}},
			wantErr: true,
		},
		{
			name:      "only client_secret from env",
			envSecret: "env-secret",
			wantErr:   true,
		},
		{
			name:       "client_id in block, secret from env",
			envSecret:  "env-secret",
			block:      []any{map[string]any{"client_id": "block-id"}},
			wantConfig: &api.CloudflareAccessConfig{ClientID: "block-id", ClientSecret: "env-secret"},
		},
		{
			name:       "both from env, block absent",
			envID:      "env-id",
			envSecret:  "env-secret",
			wantConfig: &api.CloudflareAccessConfig{ClientID: "env-id", ClientSecret: "env-secret"},
		},
		{
			name:       "both from env, empty block",
			envID:      "env-id",
			envSecret:  "env-secret",
			block:      []any{map[string]any{}},
			wantConfig: &api.CloudflareAccessConfig{ClientID: "env-id", ClientSecret: "env-secret"},
		},
		{
			name:       "block overrides env",
			envID:      "env-id",
			envSecret:  "env-secret",
			block:      []any{map[string]any{"client_id": "block-id", "client_secret": "block-secret"}},
			wantConfig: &api.CloudflareAccessConfig{ClientID: "block-id", ClientSecret: "block-secret"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(envProviderCloudflareAccessClientID, tt.envID)
			t.Setenv(envProviderCloudflareAccessClientSecret, tt.envSecret)

			raw := map[string]any{}
			if tt.block != nil {
				raw[mkProviderCloudflareAccess] = tt.block
			}

			d := schema.TestResourceDataRaw(t, createSchema(), raw)

			got, diags := cloudflareAccessConfig(d)

			if tt.wantErr {
				require.True(t, diags.HasError())
				require.Nil(t, got)

				return
			}

			require.False(t, diags.HasError(), diags)
			require.Equal(t, tt.wantConfig, got)
		})
	}
}
