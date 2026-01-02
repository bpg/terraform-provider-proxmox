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

// ACME Certificate Test Environment Setup
//
// These tests require a properly configured ACME environment. Here's how to set it up:
//
// 1. ACME Account: Create an ACME account in Proxmox VE (Datacenter -> ACME -> Add Account)
//    - Use Let's Encrypt staging for testing: https://acme-staging-v02.api.letsencrypt.org/directory
//    - For production: https://acme-v02.api.letsencrypt.org/directory
//    - See: https://pve.proxmox.com/wiki/Certificate_Management
//
// 2. DNS Plugin (for DNS-01 challenge): Configure a DNS plugin in Proxmox VE
//    - Datacenter -> ACME -> Challenge Plugins -> Add
//    - Supported plugins: https://pve.proxmox.com/pve-docs/pve-admin-guide.html#sysadmin_certs_acme_plugins
//    - Example providers: Cloudflare, Desec, DigitalOcean, etc.
//
// 3. Environment Variables:
//    - PROXMOX_VE_ACC_ACME_ACCOUNT_NAME: Name of the ACME account configured in step 1
//    - PROXMOX_VE_ACC_ACME_DOMAIN: Domain name to use for the certificate (must be resolvable)
//    - PROXMOX_VE_ACC_ACME_DNS_PLUGIN: (Optional) Name of the DNS plugin for DNS-01 challenge
//
// Note: HTTP-01 challenge requires the Proxmox node to be reachable on port 80 from the internet.
// DNS-01 challenge is recommended for testing as it doesn't require inbound connectivity.

// TestAccResourceACMECertificate tests the ACME certificate resource.
func TestAccResourceACMECertificate(t *testing.T) {
	te := test.InitEnvironment(t)

	acmeAccount := utils.GetAnyStringEnv("PROXMOX_VE_ACC_ACME_ACCOUNT_NAME")
	acmeDomain := utils.GetAnyStringEnv("PROXMOX_VE_ACC_ACME_DOMAIN")
	dnsPlugin := utils.GetAnyStringEnv("PROXMOX_VE_ACC_ACME_DNS_PLUGIN")

	if acmeAccount == "" || acmeDomain == "" {
		t.Skip("Skipping ACME certificate test - set PROXMOX_VE_ACC_ACME_ACCOUNT_NAME and PROXMOX_VE_ACC_ACME_DOMAIN")
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
						force     = true
						{{.DomainsConfig}}
					}`, test.WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_acme_certificate.test_cert", map[string]string{
						"node_name": te.NodeName,
						"account":   acmeAccount,
						"force":     "true",
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
				ImportStateId:     te.NodeName,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"force", // force is not stored in state
				},
			},
		},
	})
}
