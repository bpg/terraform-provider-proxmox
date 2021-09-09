/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"fmt"
	"strings"
	"time"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func resourceVirtualEnvironmentCertificate() *schema.Resource {
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
		Create: resourceVirtualEnvironmentCertificateCreate,
		Read:   resourceVirtualEnvironmentCertificateRead,
		Update: resourceVirtualEnvironmentCertificateUpdate,
		Delete: resourceVirtualEnvironmentCertificateDelete,
	}
}

func resourceVirtualEnvironmentCertificateCreate(d *schema.ResourceData, m interface{}) error {
	err := resourceVirtualEnvironmentCertificateUpdate(d, m)

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentCertificateNodeName).(string)

	d.SetId(fmt.Sprintf("%s_certificate", nodeName))

	return nil
}

func resourceVirtualEnvironmentCertificateGetUpdateBody(d *schema.ResourceData, m interface{}) (*proxmox.VirtualEnvironmentCertificateUpdateRequestBody, error) {
	certificate := d.Get(mkResourceVirtualEnvironmentCertificateCertificate).(string)
	certificateChain := d.Get(mkResourceVirtualEnvironmentCertificateCertificateChain).(string)
	overwrite := proxmox.CustomBool(d.Get(mkResourceVirtualEnvironmentCertificateOverwrite).(bool))
	privateKey := d.Get(mkResourceVirtualEnvironmentCertificatePrivateKey).(string)

	combinedCertificates := strings.TrimSpace(certificate) + "\n"

	if certificateChain != "" {
		combinedCertificates += strings.TrimSpace(certificateChain) + "\n"
	}

	force := overwrite

	if d.Id() != "" {
		force = proxmox.CustomBool(true)
	}

	restart := proxmox.CustomBool(true)

	body := &proxmox.VirtualEnvironmentCertificateUpdateRequestBody{
		Certificates: combinedCertificates,
		Force:        &force,
		PrivateKey:   &privateKey,
		Restart:      &restart,
	}

	return body, nil
}

func resourceVirtualEnvironmentCertificateRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentCertificateNodeName).(string)
	list, err := veClient.ListCertificates(nodeName)

	if err != nil {
		return err
	}

	d.Set(mkResourceVirtualEnvironmentCertificateCertificate, "")
	d.Set(mkResourceVirtualEnvironmentCertificateCertificateChain, "")

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

				d.Set(mkResourceVirtualEnvironmentCertificateCertificate, newCertificate)
				d.Set(mkResourceVirtualEnvironmentCertificateCertificateChain, newCertificateChain)
			}

			d.Set(mkResourceVirtualEnvironmentCertificateFileName, *c.FileName)

			if c.NotAfter != nil {
				t := time.Time(*c.NotAfter)

				d.Set(mkResourceVirtualEnvironmentCertificateExpirationDate, t.UTC().Format(time.RFC3339))
			} else {
				d.Set(mkResourceVirtualEnvironmentCertificateExpirationDate, "")
			}

			if c.Issuer != nil {
				d.Set(mkResourceVirtualEnvironmentCertificateIssuer, *c.Issuer)
			} else {
				d.Set(mkResourceVirtualEnvironmentCertificateIssuer, "")
			}

			if c.PublicKeyBits != nil {
				d.Set(mkResourceVirtualEnvironmentCertificatePublicKeySize, *c.PublicKeyBits)
			} else {
				d.Set(mkResourceVirtualEnvironmentCertificatePublicKeySize, 0)
			}

			if c.PublicKeyType != nil {
				pkType := *c.PublicKeyType

				for _, pkt := range []string{"ecdsa", "dsa", "rsa"} {
					if strings.Contains(pkType, pkt) {
						pkType = pkt
					}
				}

				d.Set(mkResourceVirtualEnvironmentCertificatePublicKeyType, pkType)
			} else {
				d.Set(mkResourceVirtualEnvironmentCertificatePublicKeyType, "")
			}

			if c.Fingerprint != nil {
				d.Set(mkResourceVirtualEnvironmentCertificateSSLFingerprint, *c.Fingerprint)
			} else {
				d.Set(mkResourceVirtualEnvironmentCertificateSSLFingerprint, "")
			}

			if c.NotBefore != nil {
				t := time.Time(*c.NotBefore)

				d.Set(mkResourceVirtualEnvironmentCertificateStartDate, t.UTC().Format(time.RFC3339))
			} else {
				d.Set(mkResourceVirtualEnvironmentCertificateStartDate, "")
			}

			if c.Subject != nil {
				d.Set(mkResourceVirtualEnvironmentCertificateSubject, *c.Subject)
			} else {
				d.Set(mkResourceVirtualEnvironmentCertificateSubject, "")
			}

			if c.SubjectAlternativeNames != nil {
				sanList := make([]interface{}, len(*c.SubjectAlternativeNames))

				for i, san := range *c.SubjectAlternativeNames {
					sanList[i] = san
				}

				d.Set(mkResourceVirtualEnvironmentCertificateSubjectAlternativeNames, sanList)
			} else {
				d.Set(mkResourceVirtualEnvironmentCertificateSubjectAlternativeNames, []interface{}{})
			}
		}
	}

	return nil
}

func resourceVirtualEnvironmentCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentCertificateNodeName).(string)

	body, err := resourceVirtualEnvironmentCertificateGetUpdateBody(d, m)

	if err != nil {
		return err
	}

	err = veClient.UpdateCertificate(nodeName, body)

	if err != nil {
		return err
	}

	return resourceVirtualEnvironmentCertificateRead(d, m)
}

func resourceVirtualEnvironmentCertificateDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentCertificateNodeName).(string)
	restart := proxmox.CustomBool(true)

	err = veClient.DeleteCertificate(nodeName, &proxmox.VirtualEnvironmentCertificateDeleteRequestBody{
		Restart: &restart,
	})

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
