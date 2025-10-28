//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

// TestAccResourceACMECertificate tests the ACME certificate resource.
// Note: This test requires a properly configured ACME environment:
// - Set PROXMOX_VE_ACC_ACME_ACCOUNT_NAME environment variable
// - Set PROXMOX_VE_ACC_ACME_DOMAIN environment variable
// - Set PROXMOX_VE_ACC_ACME_DNS_PLUGIN (optional, for DNS-01 challenge)
// The test will be skipped if these are not set.
func TestAccResourceACMECertificate(t *testing.T) {
	te := test.InitEnvironment(t)

	acmeAccount := utils.GetAnyStringEnv("PROXMOX_VE_ACC_ACME_ACCOUNT_NAME")
	acmeDomain := utils.GetAnyStringEnv("PROXMOX_VE_ACC_ACME_DOMAIN")
	dnsPlugin := utils.GetAnyStringEnv("PROXMOX_VE_ACC_ACME_DNS_PLUGIN")

	if acmeAccount == "" || acmeDomain == "" {
		t.Skip("Skipping ACME certificate test - set PROXMOX_VE_ACC_ACME_ACCOUNT_NAME and PROXMOX_VE_ACC_ACME_DOMAIN")
	}

	nodeName := te.NodeName
	if nodeName == "" {
		nodeName = "pve"
	}

	// Build domains config
	domainsConfig := `domains = [{
		domain = "` + acmeDomain + `"`

	if dnsPlugin != "" {
		domainsConfig += `
		plugin = "` + dnsPlugin + `"`
	}

	domainsConfig += `
	}]`

	te.AddTemplateVars(map[string]interface{}{
		"NodeName":      nodeName,
		"Account":       acmeAccount,
		"DomainsConfig": domainsConfig,
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_acme_certificate" "test_cert" {
						node_name = "{{.NodeName}}"
						account   = "{{.Account}}"
						force     = false
						{{.DomainsConfig}}
					}`, test.WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_acme_certificate.test_cert", map[string]string{
						"node_name": nodeName,
						"account":   acmeAccount,
						"force":     "false",
					}),
					test.ResourceAttributesSet("proxmox_virtual_environment_acme_certificate.test_cert", []string{
						"certificate",
						"fingerprint",
						"issuer",
						"subject",
						"not_after",
						"not_before",
					}),
				),
			},
		},
	})
}

// TestAccResourceACMECertificate_Import tests importing an ACME certificate resource.
func TestAccResourceACMECertificate_Import(t *testing.T) {
	te := test.InitEnvironment(t)

	acmeAccount := utils.GetAnyStringEnv("PROXMOX_VE_ACC_ACME_ACCOUNT_NAME")
	acmeDomain := utils.GetAnyStringEnv("PROXMOX_VE_ACC_ACME_DOMAIN")
	dnsPlugin := utils.GetAnyStringEnv("PROXMOX_VE_ACC_ACME_DNS_PLUGIN")

	if acmeAccount == "" || acmeDomain == "" {
		t.Skip("Skipping ACME certificate import test - set PROXMOX_VE_ACC_ACME_ACCOUNT_NAME and PROXMOX_VE_ACC_ACME_DOMAIN")
	}

	nodeName := te.NodeName
	if nodeName == "" {
		nodeName = "pve"
	}

	// Build domains config
	domainsConfig := `domains = [{
		domain = "` + acmeDomain + `"`

	if dnsPlugin != "" {
		domainsConfig += `
		plugin = "` + dnsPlugin + `"`
	}

	domainsConfig += `
	}]`

	te.AddTemplateVars(map[string]interface{}{
		"NodeName":      nodeName,
		"Account":       acmeAccount,
		"DomainsConfig": domainsConfig,
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_acme_certificate" "test_cert_import" {
						node_name = "{{.NodeName}}"
						account   = "{{.Account}}"
						force     = true
						{{.DomainsConfig}}
					}`, test.WithRootUser()),
			},
			{
				ResourceName:      "proxmox_virtual_environment_acme_certificate.test_cert_import",
				ImportState:       true,
				ImportStateId:     nodeName,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"force", // force is not stored in state
				},
			},
		},
	})
}
