# Terraform Provider for Proxmox

[![Go Report Card](https://goreportcard.com/badge/github.com/bpg/terraform-provider-proxmox)](https://goreportcard.com/report/github.com/bpg/terraform-provider-proxmox)
[![GoDoc](https://godoc.org/github.com/bpg/terraform-provider-proxmox?status.svg)](http://godoc.org/github.com/bpg/terraform-provider-proxmox)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/bpg/terraform-provider-proxmox)](https://github.com/bpg/terraform-provider-proxmox/releases/latest)
[![GitHub Release Date](https://img.shields.io/github/release-date/bpg/terraform-provider-proxmox)](https://github.com/bpg/terraform-provider-proxmox/releases/latest)
[![GitHub stars](https://img.shields.io/github/stars/bpg/terraform-provider-proxmox)](https://github.com/bpg/terraform-provider-proxmox/stargazers)
[![All Contributors](https://img.shields.io/github/all-contributors/bpg/terraform-provider-proxmox)](#contributors)
[![Conventional Commits](https://img.shields.io/badge/conventional%20commits-v1.0.0-ff69b4)](https://www.conventionalcommits.org/en/v1.0.0/)
[![Buy Me A Coffee](https://img.shields.io/badge/-buy%20me%20a%20coffee-5F7FFF?logo=buymeacoffee&labelColor=gray&logoColor=FFDD00)](https://www.buymeacoffee.com/bpgca)

A Terraform Provider which adds support for Proxmox solutions.

This repository is a fork
of <https://github.com/danitso/terraform-provider-proxmox>
which is no longer maintained.

## Compatibility promise

This provider is compatible with the latest version of Proxmox VE (currently
8.0). While it may work with older 7.x versions, it is not guaranteed to do so.

While provider is on version 0.x, it is not guaranteed to be backwards
compatible with all previous minor versions. However, we will try to keep the
backwards compatibility between provider versions as much as possible.

## Requirements

- [Proxmox Virtual Environment](https://www.proxmox.com/en/proxmox-virtual-environment/) 8.x
- TLS 1.3 for the Proxmox API endpoint
- [Terraform](https://www.terraform.io/downloads.html) 1.4+
- [Go](https://golang.org/doc/install) 1.21 (to build the provider plugin)

## Building the provider

- Clone the repository
  to `$GOPATH/src/github.com/bpg/terraform-provider-proxmox`:

  ```sh
  mkdir -p "${GOPATH}/src/github.com/bpg"
  cd "${GOPATH}/src/github.com/bpg"
  git clone git@github.com:bpg/terraform-provider-proxmox
  ```

- Enter the provider directory and build it:

  ```sh
  cd "${GOPATH}/src/github.com/bpg/terraform-provider-proxmox"
  make build
  ```

## Using the provider

You can find the latest release and its documentation in
the [Terraform Registry](https://registry.terraform.io/providers/bpg/proxmox/latest).

## Testing the provider

In order to test the provider, you can simply run `make test`.

```sh
make test
```

Tests are limited to regression tests, ensuring backwards compatibility.

A limited number of acceptance tests are available in the `proxmoxtf/test` directory, mostly
for "new" functionality implemented using the Terraform Provider Framework. These tests
are not run by default, as they require a Proxmox VE environment to be available.
They can be run using `make testacc`, the Proxmox connection can be configured using
environment variables, see provider documentation for details.

## Deploying the example resources

There are number of TF examples in the `example` directory, which can be used
to deploy a Container, VM, or other Proxmox resources on your test Proxmox
environment. The following assumptions are made about the test environment:

- It has one node named `pve`
- The node has local storages named `local` and `local-lvm`
- The "Snippets" content type is enabled in `local` storage

Create `examples/terraform.tfvars` with the following variables:

```sh
virtual_environment_username = "root@pam"
virtual_environment_password = "put-your-password-here"
virtual_environment_endpoint = "https://<your-cluster-endpoint>:8006/"
```

Then run `make example` to deploy the example resources.

## Future work

The provider is using
the [Terraform SDKv2](https://developer.hashicorp.com/terraform/plugin/sdkv2),
which is considered legacy and is in maintenance mode.
The work has started to migrate the provider to the
new [Terraform Plugin Framework](https://www.terraform.io/docs/extend/plugin-sdk.html),
with aim to release it as a new major version **1.0**.

## Known issues

### Disk images cannot be imported by non-PAM accounts

Due to limitations in the Proxmox VE API, certain actions need to be performed
using SSH. This requires the use of a PAM account (standard Linux account).

### Disk images from VMware cannot be uploaded or imported

Proxmox VE is not currently supporting VMware disk images directly. However, you
can still use them as disk images by using this workaround:

```hcl
resource "proxmox_virtual_environment_file" "vmdk_disk_image" {
  content_type = "iso"
  datastore_id = "datastore-id"
  node_name    = "node-name"

  source_file {
    # We must override the file extension to bypass the validation code
    # in the Proxmox VE API.
    file_name = "vmdk-file-name.img"
    path      = "path-to-vmdk-file"
  }
}

resource "proxmox_virtual_environment_vm" "example" {
  //...

  disk {
    datastore_id = "datastore-id"
    # We must tell the provider that the file format is vmdk instead of qcow2.
    file_format  = "vmdk"
    file_id      = "${proxmox_virtual_environment_file.vmdk_disk_image.id}"
  }

  //...
}
```

### Snippets cannot be uploaded by non-PAM accounts

Due to limitations in the Proxmox VE API, certain files need to be uploaded
using SFTP. This requires the use of a PAM account (standard Linux account).

## Contributors

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://danitso.com/"><img src="https://avatars.githubusercontent.com/u/7096448?v=4?s=100" width="100px;" alt="Dan R. Petersen"/><br /><sub><b>Dan R. Petersen</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=danitso-dp" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/bpg"><img src="https://avatars.githubusercontent.com/u/627562?v=4?s=100" width="100px;" alt="Pavel Boldyrev"/><br /><sub><b>Pavel Boldyrev</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=bpg" title="Code">💻</a> <a href="#maintenance-bpg" title="Maintenance">🚧</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/luhahn"><img src="https://avatars.githubusercontent.com/u/61747797?v=4?s=100" width="100px;" alt="Lucas Hahn"/><br /><sub><b>Lucas Hahn</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=luhahn" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/blz-ea"><img src="https://avatars.githubusercontent.com/u/19339605?v=4?s=100" width="100px;" alt="Alex Kulikovskikh"/><br /><sub><b>Alex Kulikovskikh</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=blz-ea" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/otopetrik"><img src="https://avatars.githubusercontent.com/u/972298?v=4?s=100" width="100px;" alt="Oto Petřík"/><br /><sub><b>Oto Petřík</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=otopetrik" title="Code">💻</a> <a href="#question-otopetrik" title="Answering Questions">💬</a> <a href="https://github.com/bpg/terraform-provider-proxmox/issues?q=author%3Aotopetrik" title="Bug reports">🐛</a> <a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=otopetrik" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://www.patreon.com/boik"><img src="https://avatars.githubusercontent.com/u/6451933?v=4?s=100" width="100px;" alt="Boik"/><br /><sub><b>Boik</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=qazbnm456" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/abdo-farag"><img src="https://avatars.githubusercontent.com/u/10170837?v=4?s=100" width="100px;" alt="Abdelfadeel Farag"/><br /><sub><b>Abdelfadeel Farag</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=abdo-farag" title="Code">💻</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/kugo12"><img src="https://avatars.githubusercontent.com/u/15050771?v=4?s=100" width="100px;" alt="Szczepan Wiśniowski"/><br /><sub><b>Szczepan Wiśniowski</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=kugo12" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/xonvanetta"><img src="https://avatars.githubusercontent.com/u/11271952?v=4?s=100" width="100px;" alt="Fabian Heib"/><br /><sub><b>Fabian Heib</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=xonvanetta" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/kaje783"><img src="https://avatars.githubusercontent.com/u/120482249?v=4?s=100" width="100px;" alt="kaje783"/><br /><sub><b>kaje783</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=kaje783" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/michalg91"><img src="https://avatars.githubusercontent.com/u/63045346?v=4?s=100" width="100px;" alt="michalg91"/><br /><sub><b>michalg91</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=michalg91" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/1-cameron"><img src="https://avatars.githubusercontent.com/u/68611194?v=4?s=100" width="100px;" alt="Cameron"/><br /><sub><b>Cameron</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=1-cameron" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://christopherjones.us/"><img src="https://avatars.githubusercontent.com/u/115515?v=4?s=100" width="100px;" alt="Chris Jones"/><br /><sub><b>Chris Jones</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=magikid" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://dominik.wombacher.cc/"><img src="https://avatars.githubusercontent.com/u/16312366?v=4?s=100" width="100px;" alt="Dominik Wombacher"/><br /><sub><b>Dominik Wombacher</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=wombelix" title="Code">💻</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="http://frank.villaro-dixon.eu/"><img src="https://avatars.githubusercontent.com/u/17879459?v=4?s=100" width="100px;" alt="Frank Villaro-Dixon"/><br /><sub><b>Frank Villaro-Dixon</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=Frankkkkk" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/groggemans"><img src="https://avatars.githubusercontent.com/u/11381284?v=4?s=100" width="100px;" alt="Gertjan Roggemans"/><br /><sub><b>Gertjan Roggemans</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=groggemans" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/HenriAW"><img src="https://avatars.githubusercontent.com/u/24527359?v=4?s=100" width="100px;" alt="Henri Williams"/><br /><sub><b>Henri Williams</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=HenriAW" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/zeddD1abl0"><img src="https://avatars.githubusercontent.com/u/8335605?v=4?s=100" width="100px;" alt="Jordan Keith"/><br /><sub><b>Jordan Keith</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=zeddD1abl0" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/shortmann"><img src="https://avatars.githubusercontent.com/u/20142334?v=4?s=100" width="100px;" alt="Kai Kahllund"/><br /><sub><b>Kai Kahllund</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=shortmann" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/kevinglasson"><img src="https://avatars.githubusercontent.com/u/22187575?v=4?s=100" width="100px;" alt="Kevin"/><br /><sub><b>Kevin</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=kevinglasson" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/krzysztof-magosa"><img src="https://avatars.githubusercontent.com/u/6112411?v=4?s=100" width="100px;" alt="Krzysztof Magosa"/><br /><sub><b>Krzysztof Magosa</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=krzysztof-magosa" title="Code">💻</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://burchett.dev/"><img src="https://avatars.githubusercontent.com/u/783042?v=4?s=100" width="100px;" alt="Matt Burchett"/><br /><sub><b>Matt Burchett</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=mattburchett" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/moyiz"><img src="https://avatars.githubusercontent.com/u/8603313?v=4?s=100" width="100px;" alt="Moyiz"/><br /><sub><b>Moyiz</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=moyiz" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/pescobar"><img src="https://avatars.githubusercontent.com/u/103797?v=4?s=100" width="100px;" alt="Pablo Escobar Lopez"/><br /><sub><b>Pablo Escobar Lopez</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=pescobar" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://hrmny.sh/"><img src="https://avatars.githubusercontent.com/u/8845940?v=4?s=100" width="100px;" alt="Leah"/><br /><sub><b>Leah</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=ForsakenHarmony" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/wbpascal"><img src="https://avatars.githubusercontent.com/u/9532590?v=4?s=100" width="100px;" alt="Pascal Wiedenbeck"/><br /><sub><b>Pascal Wiedenbeck</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=wbpascal" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/Patricol"><img src="https://avatars.githubusercontent.com/u/13428020?v=4?s=100" width="100px;" alt="Patrick Collins"/><br /><sub><b>Patrick Collins</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=Patricol" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://prajwal-portfolio.netlify.app/"><img src="https://avatars.githubusercontent.com/u/48290911?v=4?s=100" width="100px;" alt="Prajwal"/><br /><sub><b>Prajwal</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=PrajwalBorkar" title="Code">💻</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/rafsaf"><img src="https://avatars.githubusercontent.com/u/51059348?v=4?s=100" width="100px;" alt="Rafał Safin"/><br /><sub><b>Rafał Safin</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=rafsaf" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/RemkoMolier"><img src="https://avatars.githubusercontent.com/u/16520301?v=4?s=100" width="100px;" alt="Remko Molier"/><br /><sub><b>Remko Molier</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=RemkoMolier" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="http://www.tuomassalmi.com/"><img src="https://avatars.githubusercontent.com/u/3398165?v=4?s=100" width="100px;" alt="Tuomas Salmi"/><br /><sub><b>Tuomas Salmi</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=Tumetsu" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/ikiris"><img src="https://avatars.githubusercontent.com/u/4852950?v=4?s=100" width="100px;" alt="ikiris"/><br /><sub><b>ikiris</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=ikiris" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/mleone87"><img src="https://avatars.githubusercontent.com/u/807457?v=4?s=100" width="100px;" alt="mleone87"/><br /><sub><b>mleone87</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=mleone87" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://thiscute.world/en/"><img src="https://avatars.githubusercontent.com/u/22363274?v=4?s=100" width="100px;" alt="Ryan Yin"/><br /><sub><b>Ryan Yin</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=ryan4yin" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/zoop-btc"><img src="https://avatars.githubusercontent.com/u/101409458?v=4?s=100" width="100px;" alt="zoop"/><br /><sub><b>zoop</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=zoop-btc" title="Code">💻</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://www.si458.co.uk"><img src="https://avatars.githubusercontent.com/u/765314?v=4?s=100" width="100px;" alt="Simon Smith"/><br /><sub><b>Simon Smith</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/issues?q=author%3Asi458" title="Bug reports">🐛</a> <a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=si458" title="Tests">⚠️</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/grzeg1"><img src="https://avatars.githubusercontent.com/u/8179857?v=4?s=100" width="100px;" alt="grzeg1"/><br /><sub><b>grzeg1</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/issues?q=author%3Agrzeg1" title="Bug reports">🐛</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/moustafab"><img src="https://avatars.githubusercontent.com/u/27738648?v=4?s=100" width="100px;" alt="Moustafa Baiou"/><br /><sub><b>Moustafa Baiou</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/issues?q=author%3Amoustafab" title="Bug reports">🐛</a> <a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=moustafab" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/dandaolrian"><img src="https://avatars.githubusercontent.com/u/86479955?v=4?s=100" width="100px;" alt="dandaolrian"/><br /><sub><b>dandaolrian</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=dandaolrian" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/yoshikakbudto"><img src="https://avatars.githubusercontent.com/u/10331946?v=4?s=100" width="100px;" alt="Dmitry"/><br /><sub><b>Dmitry</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/issues?q=author%3Ayoshikakbudto" title="Bug reports">🐛</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://michael.franzl.name"><img src="https://avatars.githubusercontent.com/u/72123?v=4?s=100" width="100px;" alt="Michael Franzl"/><br /><sub><b>Michael Franzl</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/issues?q=author%3Amichaelfranzl" title="Bug reports">🐛</a></td>
      <td align="center" valign="top" width="14.28%"><a href="http://www.ebenoit.info"><img src="https://avatars.githubusercontent.com/u/1409844?v=4?s=100" width="100px;" alt="Emmanuel Benoît"/><br /><sub><b>Emmanuel Benoît</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=tseeker" title="Code">💻</a> <a href="https://github.com/bpg/terraform-provider-proxmox/issues?q=author%3Atseeker" title="Bug reports">🐛</a> <a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=tseeker" title="Tests">⚠️</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/mandrav"><img src="https://avatars.githubusercontent.com/u/1273530?v=4?s=100" width="100px;" alt="mandrav"/><br /><sub><b>mandrav</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=mandrav" title="Code">💻</a> <a href="https://github.com/bpg/terraform-provider-proxmox/issues?q=author%3Amandrav" title="Bug reports">🐛</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/michaelze"><img src="https://avatars.githubusercontent.com/u/673902?v=4?s=100" width="100px;" alt="Michael Iseli"/><br /><sub><b>Michael Iseli</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=michaelze" title="Code">💻</a> <a href="https://github.com/bpg/terraform-provider-proxmox/issues?q=author%3Amichaelze" title="Bug reports">🐛</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/zharalim"><img src="https://avatars.githubusercontent.com/u/1004061?v=4?s=100" width="100px;" alt="Risto Oikarinen"/><br /><sub><b>Risto Oikarinen</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=zharalim" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/dawidole"><img src="https://avatars.githubusercontent.com/u/37155335?v=4?s=100" width="100px;" alt="dawidole"/><br /><sub><b>dawidole</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/issues?q=author%3Adawidole" title="Bug reports">🐛</a></td>
      <td align="center" valign="top" width="14.28%"><a href="http://www.krupa.me.uk/"><img src="https://avatars.githubusercontent.com/u/5756726?v=4?s=100" width="100px;" alt="Gerard Krupa"/><br /><sub><b>Gerard Krupa</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=GJKrupa" title="Tests">⚠️</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://simoncaron.com"><img src="https://avatars.githubusercontent.com/u/8635747?v=4?s=100" width="100px;" alt="Simon Caron"/><br /><sub><b>Simon Caron</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=simoncaron" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/ishioni"><img src="https://avatars.githubusercontent.com/u/50323052?v=4?s=100" width="100px;" alt="Piotr Maksymiuk"/><br /><sub><b>Piotr Maksymiuk</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/issues?q=author%3Aishioni" title="Bug reports">🐛</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/0xinterface"><img src="https://avatars.githubusercontent.com/u/890207?v=4?s=100" width="100px;" alt="Kristopher"/><br /><sub><b>Kristopher</b></sub></a><br /><a href="#ideas-0xinterface" title="Ideas, Planning, & Feedback">🤔</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/mritalian"><img src="https://avatars.githubusercontent.com/u/15789014?v=4?s=100" width="100px;" alt="Eric B"/><br /><sub><b>Eric B</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=mritalian" title="Tests">⚠️</a> <a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=mritalian" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/2b"><img src="https://avatars.githubusercontent.com/u/829041?v=4?s=100" width="100px;" alt="2b"/><br /><sub><b>2b</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/issues?q=author%3A2b" title="Bug reports">🐛</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/TheNotary"><img src="https://avatars.githubusercontent.com/u/799247?v=4?s=100" width="100px;" alt="TheNotary"/><br /><sub><b>TheNotary</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=TheNotary" title="Code">💻</a> <a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=TheNotary" title="Tests">⚠️</a> <a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=TheNotary" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/zamrih"><img src="https://avatars.githubusercontent.com/u/1061718?v=4?s=100" width="100px;" alt="zamrih"/><br /><sub><b>zamrih</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/issues?q=author%3Azamrih" title="Bug reports">🐛</a> <a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=zamrih" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/ratiborusx"><img src="https://avatars.githubusercontent.com/u/123507924?v=4?s=100" width="100px;" alt="Ratiborus"/><br /><sub><b>Ratiborus</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/issues?q=author%3Aratiborusx" title="Bug reports">🐛</a> <a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=ratiborusx" title="Tests">⚠️</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/skleinjung"><img src="https://avatars.githubusercontent.com/u/17599474?v=4?s=100" width="100px;" alt="Sean Kleinjung"/><br /><sub><b>Sean Kleinjung</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/issues?q=author%3Askleinjung" title="Bug reports">🐛</a> <a href="#financial-skleinjung" title="Financial">💵</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/muhlba91"><img src="https://avatars.githubusercontent.com/u/653739?v=4?s=100" width="100px;" alt="Daniel Mühlbachler-Pietrzykowski"/><br /><sub><b>Daniel Mühlbachler-Pietrzykowski</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=muhlba91" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/windowsrefund"><img src="https://avatars.githubusercontent.com/u/512222?v=4?s=100" width="100px;" alt="windowsrefund"/><br /><sub><b>windowsrefund</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=windowsrefund" title="Tests">⚠️</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/Fabiosilvero"><img src="https://avatars.githubusercontent.com/u/22865938?v=4?s=100" width="100px;" alt="Fabiosilvero"/><br /><sub><b>Fabiosilvero</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=Fabiosilvero" title="Tests">⚠️</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://danielhabenicht.github.io/"><img src="https://avatars.githubusercontent.com/u/13590797?v=4?s=100" width="100px;" alt="DanielHabenicht"/><br /><sub><b>DanielHabenicht</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/issues?q=author%3ADanielHabenicht" title="Bug reports">🐛</a> <a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=DanielHabenicht" title="Documentation">📖</a></td>
    </tr>
  </tbody>
  <tfoot>
    <tr>
      <td align="center" size="13px" colspan="7">
        <img src="https://raw.githubusercontent.com/all-contributors/all-contributors-cli/1b8533af435da9854653492b1327a23a4dbd0a10/assets/logo-small.svg">
          <a href="https://all-contributors.js.org/docs/en/bot/usage">Add your contributions</a>
        </img>
      </td>
    </tr>
  </tfoot>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

## Stargazers

[![Star History Chart](https://api.star-history.com/svg?repos=bpg/terraform-provider-proxmox&type=Date)](https://star-history.com/#bpg/terraform-provider-proxmox&Date)

## Sponsorship

❤️ This project is sponsored by:

- [TJ Zimmerman](https://github.com/zimmertr)

Thanks again for your support, it is much appreciated! 🙏
