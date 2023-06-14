/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/bpg/terraform-provider-proxmox/internal/provider"
)

const (
	// ProviderConfig is a shared configuration to combine with the actual
	// test configuration so the Proxmox VE client is properly configured.
	// It is also possible to use the PROXMOX_VE_ environment variables instead,.
	ProviderConfig = `
provider "proxmox" {
  endpoint  = "https://localhost:8006"  
  username = "root@pam"
  password = "password"
  insecure = true
}
`
	// such as updating the Makefile and running the testing through that tool.
)

// AccTestProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var AccTestProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"proxmox": providerserver.NewProtocol6WithError(provider.New("test")()),
}
