---
layout: page
page_title: "Setup a VM with Proxmox"
subcategory: Guides
description: |-
  This guide will help you setup a proxmox node in VM using virt-manager for a job.
---

# Setup VM with Proxmox to run examples and acceptance tests

## Who

Contributors

## Motivation

To test changes, it's best to try it on real proxmox cluster. There is dedicated `make example` command that will try to apply changes defined in `example` directory. Most resources have its examples declarations there. For example, if you add new resource, you could add new file with example resource there (ideally after adding tests). If nothing breaks, apply works fine, new resource is created and all other resources are fine, then likely change is safe.

But, proxmox node setup can be tricky task for some contributors.

## Preconditions

Be sure to install `go` and `terraform` on your system first.

## Linux (Debian/Ubuntu) with virt-manager

Goal is to have a proxmox node in VM using <https://virt-manager.org/> for a job. This text assumes some linux knowledge. Tested on Debian 12 bookworm and proxmox VE 8.1. For other distros, with any luck steps should be similar.

1. `sudo apt-get install virt-manager`.

2. Download some proxmox image from <http://download.proxmox.com/iso/>, currently latest is `proxmox-ve_8.1-1.iso`.

3. Run `virt-manager` and "create a new virtual machine", use a file you just downloaded, choose debian as a operating system, leave default network settings.

4. Give it enough RAM and disk size (required minimum is unknown for make example though I used 4GB on my 8GB laptop and 30GB disk size with success).

5. Proceed forward with installation, choose whatever you want for timezone, country, password, domain, email. Don't change other default settings.

6. After installation, log in using password from previous step and `root` username (it's proxmox default). Run: `ip a` to get assigned ip (this also appears during installation). In my case it is `192.168.122.43`.

   It may look like this:

   ```txt
   root@proxmox:~# ip a
   1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
       link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
       inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
       inet6 ::1/128 scope host noprefixroute
       valid_lft forever preferred_lft forever
   2: enp1s0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast master vmbr0 state UP group default qlen 1000
       link/ether 52:54:00:b3:22:f5 brd ff:ff:ff:ff:ff:ff
   3: vmbr0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default qlen 1000
       link/ether 52:54:00:b3:22:f5 brd ff:ff:ff:ff:ff:ff
       inet 192.168.122.43/24 scope global vmbr0
       valid_lft forever preferred_lft forever
       inet6 fe80::5054:ff:feb3:22f5/64 scope link
       valid_lft forever preferred_lft forever
   ```

7. (Optional) On **your** computer, there should be new interface created mapped to that one you see on proxmox. Again `ip a`:

   ```txt
   ...

   8: virbr0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default qlen 1000
       link/ether 52:54:00:ca:65:49 brd ff:ff:ff:ff:ff:ff
       inet 192.168.122.1/24 brd 192.168.122.255 scope global virbr0
       valid_lft forever preferred_lft forever

   ...

   ```

8. (Optional) You can SSH into proxmox node:

   ```bash
   ssh root@192.168.122.43
   ```

   You can also use browser and visit console at <https://192.168.122.43:8006>.

9. Create `terraform.tfvars` file (it will be git ignored file) in `example` folder with credentials for you new proxmox node.

   ```txt
   # example/terraform.tfvars
   virtual_environment_username = "root@pam"
   virtual_environment_endpoint = "https://192.168.122.43:8006/"
   virtual_environment_password = "your password from step 5"

   ```

10. Now you can run `make example`.

11. If you see error with proxmox_virtual_environment_file: the datastore "local" does not support content type "snippets"; supported content types are: `[backup, images, iso, rootdir, vztmpl, import]`, you need to enable them, see <https://registry.terraform.io/providers/bpg/proxmox/latest/docs/resources/virtual_environment_file#snippets>.
