/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	dvResourceVirtualEnvironmentCertificateCertificateChain = ""
	dvResourceVirtualEnvironmentCertificateOverwrite        = false

	mkResourceVirtualEnvironmentCertificateCertificate             = "certificate"
	mkResourceVirtualEnvironmentCertificateCertificateChain        = "certificate_chain"
	mkResourceVirtualEnvironmentCertificateFileName                = "file_name"
	mkResourceVirtualEnvironmentCertificateIssuer                  = "issuer"
	mkResourceVirtualEnvironmentCertificateNodeName                = "node_name"
	mkResourceVirtualEnvironmentCertificateExpirationDate          = "expiration_date"
	mkResourceVirtualEnvironmentCertificateOverwrite               = "overwrite"
	mkResourceVirtualEnvironmentCertificatePrivateKey              = "private_key"
	mkResourceVirtualEnvironmentCertificatePublicKeySize           = "public_key_size"
	mkResourceVirtualEnvironmentCertificatePublicKeyType           = "public_key_type"
	mkResourceVirtualEnvironmentCertificateSSLFingerprint          = "ssl_fingerprint"
	mkResourceVirtualEnvironmentCertificateStartDate               = "start_date"
	mkResourceVirtualEnvironmentCertificateSubject                 = "subject"
	mkResourceVirtualEnvironmentCertificateSubjectAlternativeNames = "subject_alternative_names"
)

// Certificate returns a resource that manages a certificate.
func Certificate() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentCertificateCertificate: {
				Type:        schema.TypeString,
				Description: "The PEM encoded certificate",
				Required:    true,
			},
			mkResourceVirtualEnvironmentCertificateCertificateChain: {
				Type:        schema.TypeString,
				Description: "The PEM encoded certificate chain",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentCertificateCertificateChain,
			},
			mkResourceVirtualEnvironmentCertificateExpirationDate: {
				Type:        schema.TypeString,
				Description: "The expiration date",
				Computed:    true,
			},
			mkResourceVirtualEnvironmentCertificateFileName: {
				Type:        schema.TypeString,
				Description: "The file name",
				Computed:    true,
			},
			mkResourceVirtualEnvironmentCertificateIssuer: {
				Type:        schema.TypeString,
				Description: "The issuer",
				Computed:    true,
			},
			mkResourceVirtualEnvironmentCertificateNodeName: {
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentCertificateOverwrite: {
				Type:        schema.TypeBool,
				Description: "Whether to overwrite an existing certificate",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentCertificateOverwrite,
			},
			mkResourceVirtualEnvironmentCertificatePrivateKey: {
				Type:        schema.TypeString,
				Description: "The PEM encoded private key",
				Required:    true,
				Sensitive:   true,
			},
			mkResourceVirtualEnvironmentCertificatePublicKeySize: {
				Type:        schema.TypeInt,
				Description: "The public key size",
				Computed:    true,
			},
			mkResourceVirtualEnvironmentCertificatePublicKeyType: {
				Type:        schema.TypeString,
				Description: "The public key type",
				Computed:    true,
			},
			mkResourceVirtualEnvironmentCertificateSSLFingerprint: {
				Type:        schema.TypeString,
				Description: "The SSL fingerprint",
				Computed:    true,
			},
			mkResourceVirtualEnvironmentCertificateStartDate: {
				Type:        schema.TypeString,
				Description: "The start date",
				Computed:    true,
			},
			mkResourceVirtualEnvironmentCertificateSubject: {
				Type:        schema.TypeString,
				Description: "The subject",
				Computed:    true,
			},
			mkResourceVirtualEnvironmentCertificateSubjectAlternativeNames: {
				Type:        schema.TypeList,
				Description: "The subject alternative names",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		CreateContext: certificateCreate,
		ReadContext:   certificateRead,
		UpdateContext: certificateUpdate,
		DeleteContext: certificateDelete,
	}
}

func certificateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	diags := certificateUpdate(ctx, d, m)
	if diags.HasError() {
		return diags
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentCertificateNodeName).(string)

	d.SetId(fmt.Sprintf("%s_certificate", nodeName))

	return nil
}

func certificateGetUpdateBody(d *schema.ResourceData) *nodes.CertificateUpdateRequestBody {
	certificate := d.Get(mkResourceVirtualEnvironmentCertificateCertificate).(string)
	certificateChain := d.Get(mkResourceVirtualEnvironmentCertificateCertificateChain).(string)
	overwrite := types.CustomBool(d.Get(mkResourceVirtualEnvironmentCertificateOverwrite).(bool))
	privateKey := d.Get(mkResourceVirtualEnvironmentCertificatePrivateKey).(string)

	combinedCertificates := strings.TrimSpace(certificate) + "\n"

	if certificateChain != "" {
		combinedCertificates += strings.TrimSpace(certificateChain) + "\n"
	}

	force := overwrite

	if d.Id() != "" {
		force = true
	}

	restart := types.CustomBool(true)

	body := &nodes.CertificateUpdateRequestBody{
		Certificates: combinedCertificates,
		Force:        &force,
		PrivateKey:   &privateKey,
		Restart:      &restart,
	}

	return body
}

func certificateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentCertificateNodeName).(string)
	api := veClient.API().Node(nodeName)

	list, err := api.ListCertificates(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(mkResourceVirtualEnvironmentCertificateCertificate, "")
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkResourceVirtualEnvironmentCertificateCertificateChain, "")
	diags = append(diags, diag.FromErr(err)...)

	certificateChain := d.Get(mkResourceVirtualEnvironmentCertificateCertificateChain).(string)

	for _, c := range *list {
		if c.FileName != nil && *c.FileName == "pveproxy-ssl.pem" {
			if c.Certificates != nil {
				newCertificate := ""
				newCertificateChain := ""

				if certificateChain != "" {
					certificates := strings.Split(strings.TrimSpace(*c.Certificates), "\n")
					newCertificate = certificates[0] + "\n"

					if len(certificates) > 1 {
						newCertificateChain = strings.Join(certificates[1:], "\n") + "\n"
					}
				} else {
					newCertificate = *c.Certificates
				}

				err = d.Set(mkResourceVirtualEnvironmentCertificateCertificate, newCertificate)
				diags = append(diags, diag.FromErr(err)...)
				err = d.Set(
					mkResourceVirtualEnvironmentCertificateCertificateChain,
					newCertificateChain,
				)
				diags = append(diags, diag.FromErr(err)...)
			}

			err = d.Set(mkResourceVirtualEnvironmentCertificateFileName, *c.FileName)
			diags = append(diags, diag.FromErr(err)...)

			if c.NotAfter != nil {
				t := time.Time(*c.NotAfter)
				err = d.Set(
					mkResourceVirtualEnvironmentCertificateExpirationDate,
					t.UTC().Format(time.RFC3339),
				)
			} else {
				err = d.Set(mkResourceVirtualEnvironmentCertificateExpirationDate, "")
			}
			diags = append(diags, diag.FromErr(err)...)

			if c.Issuer != nil {
				err = d.Set(mkResourceVirtualEnvironmentCertificateIssuer, *c.Issuer)
			} else {
				err = d.Set(mkResourceVirtualEnvironmentCertificateIssuer, "")
			}
			diags = append(diags, diag.FromErr(err)...)

			if c.PublicKeyBits != nil {
				err = d.Set(mkResourceVirtualEnvironmentCertificatePublicKeySize, *c.PublicKeyBits)
			} else {
				err = d.Set(mkResourceVirtualEnvironmentCertificatePublicKeySize, 0)
			}
			diags = append(diags, diag.FromErr(err)...)

			if c.PublicKeyType != nil {
				pkType := *c.PublicKeyType
				for _, pkt := range []string{"ecdsa", "dsa", "rsa"} {
					if strings.Contains(pkType, pkt) {
						pkType = pkt
					}
				}
				err = d.Set(mkResourceVirtualEnvironmentCertificatePublicKeyType, pkType)
			} else {
				err = d.Set(mkResourceVirtualEnvironmentCertificatePublicKeyType, "")
			}
			diags = append(diags, diag.FromErr(err)...)

			if c.Fingerprint != nil {
				err = d.Set(mkResourceVirtualEnvironmentCertificateSSLFingerprint, *c.Fingerprint)
			} else {
				err = d.Set(mkResourceVirtualEnvironmentCertificateSSLFingerprint, "")
			}
			diags = append(diags, diag.FromErr(err)...)

			if c.NotBefore != nil {
				t := time.Time(*c.NotBefore)
				err = d.Set(
					mkResourceVirtualEnvironmentCertificateStartDate,
					t.UTC().Format(time.RFC3339),
				)
			} else {
				err = d.Set(mkResourceVirtualEnvironmentCertificateStartDate, "")
			}
			diags = append(diags, diag.FromErr(err)...)

			if c.Subject != nil {
				err = d.Set(mkResourceVirtualEnvironmentCertificateSubject, *c.Subject)
			} else {
				err = d.Set(mkResourceVirtualEnvironmentCertificateSubject, "")
			}
			diags = append(diags, diag.FromErr(err)...)

			if c.SubjectAlternativeNames != nil {
				sanList := make([]interface{}, len(*c.SubjectAlternativeNames))
				for i, san := range *c.SubjectAlternativeNames {
					sanList[i] = san
				}
				err = d.Set(mkResourceVirtualEnvironmentCertificateSubjectAlternativeNames, sanList)
			} else {
				err = d.Set(mkResourceVirtualEnvironmentCertificateSubjectAlternativeNames, []interface{}{})
			}
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}

func certificateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentCertificateNodeName).(string)
	api := veClient.API().Node(nodeName)

	body := certificateGetUpdateBody(d)

	err = api.UpdateCertificate(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	return certificateRead(ctx, d, m)
}

func certificateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentCertificateNodeName).(string)
	api := veClient.API().Node(nodeName)

	restart := types.CustomBool(true)

	err = api.DeleteCertificate(
		ctx,
		&nodes.CertificateDeleteRequestBody{
			Restart: &restart,
		},
	)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
