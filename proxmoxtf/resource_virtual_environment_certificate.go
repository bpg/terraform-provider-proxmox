/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"fmt"
	"strings"

	"github.com/danitso/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	dvResourceVirtualEnvironmentCertificateCertificateChain = ""
	dvResourceVirtualEnvironmentCertificateOverwrite        = false

	mkResourceVirtualEnvironmentCertificateCertificate      = "certificate"
	mkResourceVirtualEnvironmentCertificateCertificateChain = "certificate_chain"
	mkResourceVirtualEnvironmentCertificateNodeName         = "node_name"
	mkResourceVirtualEnvironmentCertificateOverwrite        = "overwrite"
	mkResourceVirtualEnvironmentCertificatePrivateKey       = "private_key"
)

func resourceVirtualEnvironmentCertificate() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentCertificateCertificate: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The PEM encoded certificate",
				Required:    true,
			},
			mkResourceVirtualEnvironmentCertificateCertificateChain: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The PEM encoded certificate chain",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentCertificateCertificateChain,
			},
			mkResourceVirtualEnvironmentCertificateNodeName: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentCertificateOverwrite: &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Whether to overwrite an existing certificate",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentCertificateOverwrite,
			},
			mkResourceVirtualEnvironmentCertificatePrivateKey: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The PEM encoded private key",
				Required:    true,
				Sensitive:   true,
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
