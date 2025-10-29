/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

var (
	_ resource.Resource                = &acmeCertificateResource{}
	_ resource.ResourceWithConfigure   = &acmeCertificateResource{}
	_ resource.ResourceWithImportState = &acmeCertificateResource{}
)

// NewACMECertificateResource creates a new resource for managing ACME certificates on nodes.
func NewACMECertificateResource() resource.Resource {
	return &acmeCertificateResource{}
}

// acmeCertificateResource contains the resource's internal data.
type acmeCertificateResource struct {
	// The Proxmox client
	client proxmox.Client
}

// acmeCertificateModel maps the schema data for the ACME certificate resource.
type acmeCertificateModel struct {
	// ID is the unique identifier for the resource (node_name)
	ID types.String `tfsdk:"id"`
	// NodeName is the name of the node for which to order the certificate
	NodeName types.String `tfsdk:"node_name"`
	// ACME account name to use
	Account types.String `tfsdk:"account"`
	// Domains to include in the certificate
	Domains types.List `tfsdk:"domains"`
	// Force certificate renewal even if not due yet
	Force types.Bool `tfsdk:"force"`
	// Certificate PEM data (computed after ordering)
	Certificate types.String `tfsdk:"certificate"`
	// Certificate fingerprint (computed)
	Fingerprint types.String `tfsdk:"fingerprint"`
	// Certificate issuer (computed)
	Issuer types.String `tfsdk:"issuer"`
	// Certificate subject (computed)
	Subject types.String `tfsdk:"subject"`
	// Certificate expiration date (computed)
	NotAfter types.String `tfsdk:"not_after"`
	// Certificate start date (computed)
	NotBefore types.String `tfsdk:"not_before"`
	// Certificate subject alternative names (computed)
	SubjectAlternativeNames types.List `tfsdk:"subject_alternative_names"`
}

// acmeDomainModel maps the schema data for an ACME domain configuration.
type acmeDomainModel struct {
	// Domain name
	Domain types.String `tfsdk:"domain"`
	// DNS plugin to use for validation (optional, if not set uses standalone http-01)
	Plugin types.String `tfsdk:"plugin"`
	// Alias domain for DNS validation (optional)
	Alias types.String `tfsdk:"alias"`
}

// Metadata defines the name of the resource.
func (r *acmeCertificateResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_acme_certificate"
}

// Schema defines the schema for the resource.
func (r *acmeCertificateResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages ACME SSL certificates for Proxmox VE nodes. " +
			"This resource orders and renews certificates from an ACME Certificate Authority (like Let's Encrypt) " +
			"for a specific node.",
		MarkdownDescription: "Manages ACME SSL certificates for Proxmox VE nodes.\n\n" +
			"This resource orders and renews certificates from an ACME Certificate Authority (like Let's Encrypt) " +
			"for a specific node. Before using this resource, ensure that:\n" +
			"- An ACME account is configured (using `proxmox_virtual_environment_acme_account`)\n" +
			"- DNS plugins are configured if using DNS-01 challenge (using `proxmox_virtual_environment_acme_dns_plugin`)",
		Attributes: map[string]schema.Attribute{
			"id": attribute.ResourceID(),
			"node_name": schema.StringAttribute{
				Description: "The name of the Proxmox VE node for which to order/manage the ACME certificate.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"account": schema.StringAttribute{
				Description: "The ACME account name to use for ordering the certificate.",
				Required:    true,
			},
			"domains": schema.ListNestedAttribute{
				Description: "The list of domains to include in the certificate. At least one domain is required.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"domain": schema.StringAttribute{
							Description: "The domain name to include in the certificate.",
							Required:    true,
						},
						"plugin": schema.StringAttribute{
							Description: "The DNS plugin to use for DNS-01 challenge validation. " +
								"If not specified, the standalone HTTP-01 challenge will be used.",
							Optional: true,
						},
						"alias": schema.StringAttribute{
							Description: "An optional alias domain for DNS validation. " +
								"This allows you to validate the domain using a different domain's DNS records.",
							Optional: true,
						},
					},
				},
			},
			"force": schema.BoolAttribute{
				Description: "Force certificate renewal even if the certificate is not due for renewal yet. " +
					"Setting this to true will trigger a new certificate order on every apply.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"certificate": schema.StringAttribute{
				Description: "The PEM-encoded certificate data.",
				Computed:    true,
			},
			"fingerprint": schema.StringAttribute{
				Description: "The certificate fingerprint.",
				Computed:    true,
			},
			"issuer": schema.StringAttribute{
				Description: "The certificate issuer.",
				Computed:    true,
			},
			"subject": schema.StringAttribute{
				Description: "The certificate subject.",
				Computed:    true,
			},
			"not_after": schema.StringAttribute{
				Description: "The certificate expiration timestamp.",
				Computed:    true,
			},
			"not_before": schema.StringAttribute{
				Description: "The certificate start timestamp.",
				Computed:    true,
			},
			"subject_alternative_names": schema.ListAttribute{
				Description: "The certificate subject alternative names (SANs).",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure adds the provider-configured client to the resource.
func (r *acmeCertificateResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.Resource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected config.Resource, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = cfg.Client
}

// waitForCertificateAvailable polls ListCertificates until a certificate is available.
// This replaces the fixed time.Sleep with a more robust retry mechanism.
func (r *acmeCertificateResource) waitForCertificateAvailable(
	ctx context.Context,
	nodeClient *nodes.Client,
) (*[]nodes.CertificateListResponseData, error) {
	var certificates *[]nodes.CertificateListResponseData

	err := retry.Do(
		func() error {
			certs, err := nodeClient.ListCertificates(ctx)
			if err != nil {
				return err
			}

			// Check if any certificates are found
			if certs == nil || len(*certs) == 0 {
				return fmt.Errorf("no certificates found yet")
			}

			certificates = certs
			return nil
		},
		retry.Attempts(30),                    // Maximum 30 attempts
		retry.Delay(1 * time.Second),          // Start with 1 second delay
		retry.DelayType(retry.BackOffDelay),   // Use exponential backoff
		retry.MaxJitter(500 * time.Millisecond), // Add jitter to prevent thundering herd
		retry.Context(ctx),                    // Respect context cancellation
	)

	return certificates, err
}

// Create orders a new ACME certificate for the node.
func (r *acmeCertificateResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan acmeCertificateModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	nodeName := plan.NodeName.ValueString()
	nodeClient := r.client.Node(nodeName)

	// First, configure the node with ACME settings
	if err := r.configureNodeACME(ctx, nodeClient, &plan); err != nil {
		resp.Diagnostics.AddError(
			"Unable to configure node ACME settings",
			fmt.Sprintf("An error occurred while configuring ACME settings for node %s: %s", nodeName, err.Error()),
		)

		return
	}

	// Order the certificate
	force := proxmoxtypes.CustomBool(plan.Force.ValueBool())
	orderReq := &nodes.CertificateOrderRequestBody{
		Force: &force,
	}

	taskID, err := nodeClient.OrderCertificate(ctx, orderReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to order ACME certificate",
			fmt.Sprintf("An error occurred while ordering the ACME certificate for node %s: %s", nodeName, err.Error()),
		)

		return
	}

	// Wait for the task to complete
	if taskID != nil && *taskID != "" {
		err = nodeClient.Tasks().WaitForTask(ctx, *taskID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Certificate order task failed",
				fmt.Sprintf("The certificate order task for node %s failed: %s", nodeName, err.Error()),
			)

			return
		}
	}

	// Poll for the certificate to be available using retry mechanism
	certificates, err := r.waitForCertificateAvailable(ctx, nodeClient)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read certificate information",
			fmt.Sprintf("Failed to retrieve the ordered certificate for node %s after multiple attempts: %s", nodeName, err.Error()),
		)

		return
	}

	// Update the state with certificate information
	if err := r.updateModelFromCertificates(ctx, &plan, certificates); err != nil {
		resp.Diagnostics.AddError(
			"Unable to process certificate information",
			fmt.Sprintf("An error occurred while processing certificate information: %s", err.Error()),
		)

		return
	}

	plan.ID = plan.NodeName

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read reads the current certificate information for the node.
func (r *acmeCertificateResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state acmeCertificateModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	nodeName := state.NodeName.ValueString()
	nodeClient := r.client.Node(nodeName)

	// Read certificate information
	certificates, err := nodeClient.ListCertificates(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read certificate information",
			fmt.Sprintf("An error occurred while reading certificate information for node %s: %s", nodeName, err.Error()),
		)

		return
	}

	// Update the state with certificate information
	if err := r.updateModelFromCertificates(ctx, &state, certificates); err != nil {
		resp.Diagnostics.AddError(
			"Unable to process certificate information",
			fmt.Sprintf("An error occurred while processing certificate information: %s", err.Error()),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update renews the certificate for the node.
func (r *acmeCertificateResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state acmeCertificateModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	nodeName := plan.NodeName.ValueString()
	nodeClient := r.client.Node(nodeName)

	// Update node configuration if account or domains changed
	if !plan.Account.Equal(state.Account) || !plan.Domains.Equal(state.Domains) {
		if err := r.configureNodeACME(ctx, nodeClient, &plan); err != nil {
			resp.Diagnostics.AddError(
				"Unable to update node ACME settings",
				fmt.Sprintf("An error occurred while updating ACME settings for node %s: %s", nodeName, err.Error()),
			)

			return
		}
	}

	// Order a new certificate if force is true or other changes are made
	force := proxmoxtypes.CustomBool(plan.Force.ValueBool())
	orderReq := &nodes.CertificateOrderRequestBody{
		Force: &force,
	}

	taskID, err := nodeClient.OrderCertificate(ctx, orderReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to re-order ACME certificate",
			fmt.Sprintf("An error occurred while re-ordering the ACME certificate for node %s: %s", nodeName, err.Error()),
		)

		return
	}

	// Wait for the task to complete
	if taskID != nil && *taskID != "" {
		err = nodeClient.Tasks().WaitForTask(ctx, *taskID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Certificate renewal task failed",
				fmt.Sprintf("The certificate renewal task for node %s failed: %s", nodeName, err.Error()),
			)

			return
		}
	}

	// Poll for the certificate to be available using retry mechanism
	certificates, err := r.waitForCertificateAvailable(ctx, nodeClient)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read certificate information",
			fmt.Sprintf("Failed to retrieve the renewed certificate for node %s after multiple attempts: %s", nodeName, err.Error()),
		)

		return
	}

	// Update the state with certificate information
	if err := r.updateModelFromCertificates(ctx, &plan, certificates); err != nil {
		resp.Diagnostics.AddError(
			"Unable to process certificate information",
			fmt.Sprintf("An error occurred while processing certificate information: %s", err.Error()),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete removes the certificate resource from Terraform state and cleans up ACME configuration from the node.
// The certificate files are preserved on the Proxmox node, but the ACME configuration is removed.
func (r *acmeCertificateResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state acmeCertificateModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	nodeName := state.NodeName.ValueString()
	nodeClient := r.client.Node(nodeName)

	// Clean up the ACME configuration from the node. The certificate files will remain.
	toDelete := "acme,acmedomain0,acmedomain1,acmedomain2,acmedomain3,acmedomain4"
	configUpdate := &nodes.ConfigUpdateRequestBody{
		Delete: &toDelete,
	}

	if err := nodeClient.UpdateConfig(ctx, configUpdate); err != nil {
		// Log a warning as the resource is being deleted anyway, but the user should be notified.
		resp.Diagnostics.AddWarning(
			"Failed to clean up node ACME configuration",
			fmt.Sprintf("An error occurred while cleaning up ACME settings for node %s on delete: %s. Manual cleanup of /etc/pve/nodes/%s/config may be required.", nodeName, err.Error(), nodeName),
		)
	}
}

// ImportState imports an existing certificate by node name.
func (r *acmeCertificateResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// The import ID is the node name
	nodeName := req.ID

	nodeClient := r.client.Node(nodeName)

	// Read the node configuration to get ACME settings
	config, err := nodeClient.GetConfig(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read node configuration",
			fmt.Sprintf("An error occurred while reading configuration for node %s: %s", nodeName, err.Error()),
		)

		return
	}

	if config == nil || len(*config) == 0 {
		resp.Diagnostics.AddError(
			"Unable to read node configuration",
			fmt.Sprintf("No configuration found for node %s", nodeName),
		)

		return
	}

	nodeConfig := (*config)[0]

	// Extract ACME account from config
	var accountName string
	var domains []acmeDomainModel

	// Check for standalone ACME configuration
	if nodeConfig.ACME != nil && nodeConfig.ACME.Account != nil {
		accountName = *nodeConfig.ACME.Account
		for _, domain := range nodeConfig.ACME.Domains {
			domains = append(domains, acmeDomainModel{
				Domain: types.StringValue(domain),
				Plugin: types.StringNull(),
				Alias:  types.StringNull(),
			})
		}
	}

	// Check for DNS challenge domain configurations
	acmeDomainConfigs := []*nodes.ACMEDomainConfig{
		nodeConfig.ACMEDomain0,
		nodeConfig.ACMEDomain1,
		nodeConfig.ACMEDomain2,
		nodeConfig.ACMEDomain3,
		nodeConfig.ACMEDomain4,
	}

	for _, domainConfig := range acmeDomainConfigs {
		if domainConfig != nil {
			domains = append(domains, acmeDomainModel{
				Domain: types.StringValue(domainConfig.Domain),
				Plugin: stringPtrToValue(domainConfig.Plugin),
				Alias:  stringPtrToValue(domainConfig.Alias),
			})
		}
	}

	if accountName == "" {
		resp.Diagnostics.AddWarning(
			"ACME account not found in node configuration",
			"Could not determine the ACME account name from the node configuration. "+
				"You may need to manually set the 'account' attribute after import.",
		)
	}

	// Convert domains to Terraform list
	domainsList, diag := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"domain": types.StringType,
			"plugin": types.StringType,
			"alias":  types.StringType,
		},
	}, domains)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Set all attributes
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("node_name"), nodeName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), nodeName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("force"), false)...)

	if accountName != "" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account"), accountName)...)
	}

	if len(domains) > 0 {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domains"), domainsList)...)
	}
}

// configureNodeACME configures the node with ACME settings before ordering a certificate.
func (r *acmeCertificateResource) configureNodeACME(
	ctx context.Context,
	nodeClient *nodes.Client,
	model *acmeCertificateModel,
) error {
	// Parse domains from the model
	var domains []acmeDomainModel
	diag := model.Domains.ElementsAs(ctx, &domains, false)
	if diag.HasError() {
		return fmt.Errorf("error parsing domains: %v", diag.Errors())
	}

	if len(domains) == 0 {
		return fmt.Errorf("at least one domain is required")
	}

	// Build the config update request
	configUpdate := &nodes.ConfigUpdateRequestBody{}
	accountName := model.Account.ValueString()

	// Separate domains into standalone (no plugin) and DNS challenge (with plugin)
	var standaloneDomains []string
	var dnsDomains []acmeDomainModel

	for _, domain := range domains {
		if domain.Plugin.IsNull() || domain.Plugin.ValueString() == "" {
			standaloneDomains = append(standaloneDomains, domain.Domain.ValueString())
		} else {
			dnsDomains = append(dnsDomains, domain)
		}
	}

	// Always configure the ACME account
	// If we have standalone domains, include them; otherwise just set the account
	configUpdate.ACME = &nodes.ACMEConfig{
		Account: &accountName,
		Domains: standaloneDomains,
	}

	// Configure DNS challenge domains (up to 5 domains with DNS plugins)
	if len(dnsDomains) > 5 {
		return fmt.Errorf("Proxmox supports a maximum of 5 DNS challenge domains, got %d", len(dnsDomains))
	}

	for i, domain := range dnsDomains {
		domainConfig := &nodes.ACMEDomainConfig{
			Domain: domain.Domain.ValueString(),
		}

		if !domain.Plugin.IsNull() {
			plugin := domain.Plugin.ValueString()
			domainConfig.Plugin = &plugin
		}

		if !domain.Alias.IsNull() {
			alias := domain.Alias.ValueString()
			domainConfig.Alias = &alias
		}

		// Map to the appropriate acmedomain field
		switch i {
		case 0:
			configUpdate.ACMEDomain0 = domainConfig
		case 1:
			configUpdate.ACMEDomain1 = domainConfig
		case 2:
			configUpdate.ACMEDomain2 = domainConfig
		case 3:
			configUpdate.ACMEDomain3 = domainConfig
		case 4:
			configUpdate.ACMEDomain4 = domainConfig
		}
	}

	// Clean up unused acmedomain slots
	var toDelete []string
	for i := len(dnsDomains); i < 5; i++ {
		toDelete = append(toDelete, fmt.Sprintf("acmedomain%d", i))
	}
	if len(toDelete) > 0 {
		deleteValue := strings.Join(toDelete, ",")
		configUpdate.Delete = &deleteValue
	}

	// Update the node configuration
	return nodeClient.UpdateConfig(ctx, configUpdate)
}

// isProxmoxGeneratedCertificate checks if a certificate is generated by Proxmox itself.
// This helps identify Proxmox's self-signed or auto-generated certificates that should be skipped
// when looking for ACME certificates. Proxmox-generated certificates have "Proxmox" in the issuer.
func isProxmoxGeneratedCertificate(cert *nodes.CertificateListResponseData) bool {
	if cert.Issuer == nil {
		return false
	}
	// Check if issuer contains "Proxmox" or "PVE" (Proxmox VE)
	issuer := *cert.Issuer
	return strings.Contains(issuer, "Proxmox") || strings.Contains(issuer, "PVE")
}

// findMatchingCertificate finds the certificate that matches the domains in the model.
// It prioritizes ACME certificates (issued by certificate authorities like Let's Encrypt)
// over Proxmox-generated certificates. When multiple certificates match the configured domains,
// it returns the one with the most matching domains.
func (r *acmeCertificateResource) findMatchingCertificate(
	ctx context.Context,
	model *acmeCertificateModel,
	certificates *[]nodes.CertificateListResponseData,
) (*nodes.CertificateListResponseData, error) {
	if certificates == nil || len(*certificates) == 0 {
		return nil, fmt.Errorf("no certificates found")
	}

	// Extract domains from the model
	var domainModels []acmeDomainModel
	diag := model.Domains.ElementsAs(ctx, &domainModels, false)
	if diag.HasError() {
		// If we can't parse domains, try to find an ACME certificate (not Proxmox-generated)
		for i := range *certificates {
			if !isProxmoxGeneratedCertificate(&(*certificates)[i]) {
				return &(*certificates)[i], nil
			}
		}
		// Fall back to first certificate if all are Proxmox-generated
		return &(*certificates)[0], nil
	}

	// Extract domain strings for matching
	configDomains := make([]string, len(domainModels))
	for i, dm := range domainModels {
		configDomains[i] = dm.Domain.ValueString()
	}

	// Convert to a map for faster lookup
	domainMap := make(map[string]bool)
	for _, domain := range configDomains {
		domainMap[domain] = true
	}

	// Find the certificate that matches the most domains, preferring ACME certificates
	var bestMatch *nodes.CertificateListResponseData
	bestMatchCount := 0
	bestMatchIsProxmoxGen := true

	for i := range *certificates {
		cert := &(*certificates)[i]
		isProxmoxGen := isProxmoxGeneratedCertificate(cert)
		matchCount := 0

		// Check Subject Alternative Names (primary matching criteria for ACME certs)
		if cert.SubjectAlternativeNames != nil {
			for _, san := range *cert.SubjectAlternativeNames {
				if domainMap[san] {
					matchCount++
				}
			}
		}

		// Check Subject field (CN) if SANs don't have matches
		if cert.Subject != nil && matchCount == 0 {
			// Extract CN from Subject string (format: CN=domain.com,...)
			// Simple extraction: look for CN= and take until the next comma
			subject := *cert.Subject
			if cnIdx := strings.Index(subject, "CN="); cnIdx != -1 {
				cnStart := cnIdx + 3
				cnEnd := strings.Index(subject[cnStart:], ",")
				if cnEnd == -1 {
					cnEnd = len(subject[cnStart:])
				} else {
					cnEnd += cnStart
				}
				cn := subject[cnStart:cnEnd]
				if domainMap[cn] {
					matchCount++
				}
			}
		}

		// Update best match if:
		// 1. This certificate matches more domains, OR
		// 2. It matches the same domains but is ACME (not Proxmox-generated)
		if matchCount > bestMatchCount || (matchCount > 0 && matchCount == bestMatchCount && !isProxmoxGen && bestMatchIsProxmoxGen) {
			bestMatch = cert
			bestMatchCount = matchCount
			bestMatchIsProxmoxGen = isProxmoxGen

			// If we found an ACME certificate with all domains, we can stop searching
			if bestMatchCount == len(domainMap) && !isProxmoxGen {
				break
			}
		}
	}

	// If we found a certificate with matching domains, return it
	if bestMatch != nil && bestMatchCount > 0 {
		return bestMatch, nil
	}

	// If no domain matches found, prefer ACME certificates (not Proxmox-generated)
	for i := range *certificates {
		if !isProxmoxGeneratedCertificate(&(*certificates)[i]) {
			return &(*certificates)[i], nil
		}
	}

	// Last resort: return first certificate (shouldn't reach here in normal cases)
	return &(*certificates)[0], nil
}

// updateModelFromCertificates updates the model with certificate information.
func (r *acmeCertificateResource) updateModelFromCertificates(
	ctx context.Context,
	model *acmeCertificateModel,
	certificates *[]nodes.CertificateListResponseData,
) error {
	if certificates == nil || len(*certificates) == 0 {
		return fmt.Errorf("no certificates found")
	}

	// Find the certificate that matches the configured domains
	cert, err := r.findMatchingCertificate(ctx, model, certificates)
	if err != nil {
		return err
	}

	// Update basic certificate fields
	model.Certificate = stringPtrToValue(cert.Certificates)
	model.Fingerprint = stringPtrToValue(cert.Fingerprint)
	model.Issuer = stringPtrToValue(cert.Issuer)
	model.Subject = stringPtrToValue(cert.Subject)

	// Update timestamps
	if cert.NotAfter != nil {
		model.NotAfter = types.StringValue(time.Time(*cert.NotAfter).Format(time.RFC3339))
	} else {
		model.NotAfter = types.StringNull()
	}

	if cert.NotBefore != nil {
		model.NotBefore = types.StringValue(time.Time(*cert.NotBefore).Format(time.RFC3339))
	} else {
		model.NotBefore = types.StringNull()
	}

	// Handle subject alternative names
	if cert.SubjectAlternativeNames != nil {
		sanList := make([]types.String, 0, len(*cert.SubjectAlternativeNames))
		for _, san := range *cert.SubjectAlternativeNames {
			sanList = append(sanList, types.StringValue(san))
		}

		list, diag := types.ListValueFrom(ctx, types.StringType, sanList)
		if diag.HasError() {
			return fmt.Errorf("error creating subject_alternative_names list: %v", diag.Errors())
		}
		model.SubjectAlternativeNames = list
	} else {
		model.SubjectAlternativeNames = types.ListNull(types.StringType)
	}

	return nil
}

// stringPtrToValue converts a string pointer to a types.String value.
func stringPtrToValue(ptr *string) types.String {
	if ptr != nil {
		return types.StringValue(*ptr)
	}

	return types.StringNull()
}
