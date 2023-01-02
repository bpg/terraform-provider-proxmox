# Changelog

## [0.9.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.8.0...v0.9.0) (2023-01-01)


### Features

* **vm:** Add cloud-init network-config support ([#197](https://github.com/bpg/terraform-provider-proxmox/issues/197)) ([79a2101](https://github.com/bpg/terraform-provider-proxmox/commit/79a2101933d6001cb843050a83076a39cd503db8))
* **vm:** Add hostpci support ([01d2050](https://github.com/bpg/terraform-provider-proxmox/commit/01d20504a1924552611a92dd3f718bad270a7309))
* **vm:** Deletion of VM should also purge all storages and configs ([13080b4](https://github.com/bpg/terraform-provider-proxmox/commit/13080b44dcb08fbeabd0b20501631f52e022e46d))
* **vm:** OnBoot: change default to `true` ([#191](https://github.com/bpg/terraform-provider-proxmox/issues/191)) ([60a6818](https://github.com/bpg/terraform-provider-proxmox/commit/60a68184cf7c6239eb5cc540c746f11e2a78c240))

## [0.8.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.7.0...v0.8.0) (2022-12-13)


### Features

* add support for "ssd" disk flag for VM ([#181](https://github.com/bpg/terraform-provider-proxmox/issues/181)) ([2907346](https://github.com/bpg/terraform-provider-proxmox/commit/290734655ce28306ae910b76b8de5fedbd3b4bb8))
* add support for network_device MTU ([#176](https://github.com/bpg/terraform-provider-proxmox/issues/176)) ([3c02cb1](https://github.com/bpg/terraform-provider-proxmox/commit/3c02cb13895f7095ef0b0aaf58fe799e396a0715))
* add support for VM tags ([#169](https://github.com/bpg/terraform-provider-proxmox/issues/169)) ([ade1d49](https://github.com/bpg/terraform-provider-proxmox/commit/ade1d49117f5390e5ee58ddeadef0adf02143d33))
* add the ability to clone to non-shared storage on different nodes ([#178](https://github.com/bpg/terraform-provider-proxmox/issues/178)) ([0df14f9](https://github.com/bpg/terraform-provider-proxmox/commit/0df14f9d6aa139cb6478317da7ff6b632242b02d))


### Bug Fixes

* Check if any interface has global unicast address instead of all interfaces ([#182](https://github.com/bpg/terraform-provider-proxmox/issues/182)) ([722e010](https://github.com/bpg/terraform-provider-proxmox/commit/722e01053bdb51c038a7bd86d4018465417ea6fb))
* handling `datastore_id` in LXC template ([#180](https://github.com/bpg/terraform-provider-proxmox/issues/180)) ([63dc5cb](https://github.com/bpg/terraform-provider-proxmox/commit/63dc5cb8f6dbb6d273bd519c7768893df02a3b97))
* Remove cloned ide2 before creating new one ([#174](https://github.com/bpg/terraform-provider-proxmox/issues/174)) ([#175](https://github.com/bpg/terraform-provider-proxmox/issues/175)) ([2766555](https://github.com/bpg/terraform-provider-proxmox/commit/27665554de4a35ec678f5c63b529ccaa7d99bc74))

## [0.7.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.6.4...v0.7.0) (2022-11-18)


### Features

* Add support for custom cloud-init vendor data file ([#162](https://github.com/bpg/terraform-provider-proxmox/issues/162)) ([9e34dfb](https://github.com/bpg/terraform-provider-proxmox/commit/9e34dfb36213fc524957921e2d5b07cdf3585491))


### Bug Fixes

* linter issues ([#158](https://github.com/bpg/terraform-provider-proxmox/issues/158)) ([0fad160](https://github.com/bpg/terraform-provider-proxmox/commit/0fad160ed61cf763ce294a76e35b8c0f56cd33e8))

## [0.6.4](https://github.com/bpg/terraform-provider-proxmox/compare/v0.6.3...v0.6.4) (2022-10-17)


### Bug Fixes

* bump vulnerable dependencies ([#143](https://github.com/bpg/terraform-provider-proxmox/issues/143)) ([f9f357e](https://github.com/bpg/terraform-provider-proxmox/commit/f9f357e200681d56500316d204ed3c2dc836b551))

## [v0.6.3](https://github.com/bpg/terraform-provider-proxmox/compare/v0.6.2...v0.6.3) (2022-10-17)


### Bug Fixes

* Non-default VM disk format is not preserved in TF state ([#134](https://github.com/bpg/terraform-provider-proxmox/issues/134)) ([b09389f](https://github.com/bpg/terraform-provider-proxmox/commit/b09389f0a9c65f8f6ab82ae989d29951dd643ed2))

## [v0.6.2](https://github.com/bpg/terraform-provider-proxmox/tree/v0.6.2) (2022-09-28)

[Full Changelog](https://github.com/bpg/terraform-provider-proxmox/compare/v0.6.1...v0.6.2)

ENHANCEMENTS:

- Add discard option to vm disk creation [\#122](https://github.com/bpg/terraform-provider-proxmox/issues/122)

**Merged pull requests:**

- Add support for "discard" disk option for VM [\#128](https://github.com/bpg/terraform-provider-proxmox/pull/128) ([bpg](https://github.com/bpg))

## [v0.6.1](https://github.com/bpg/terraform-provider-proxmox/tree/v0.6.1) (2022-08-15)

[Full Changelog](https://github.com/bpg/terraform-provider-proxmox/compare/v0.6.0...v0.6.1)

BUG FIXES:

- Waiting for proxmox\_virtual\_environment\_vm's ipv4\_addresses does not really work [\#100](https://github.com/bpg/terraform-provider-proxmox/issues/100)

## [v0.6.0](https://github.com/bpg/terraform-provider-proxmox/tree/v0.6.0) (2022-08-09)

[Full Changelog](https://github.com/bpg/terraform-provider-proxmox/compare/v0.5.3...v0.6.0)

BREAKING CHANGES:

- Upgrade the provider codebase to use Terraform SDK v2 [\#91](https://github.com/bpg/terraform-provider-proxmox/pull/91) ([bpg](https://github.com/bpg))

ENHANCEMENTS:

- Add support for "iothread" disk option for VM [\#87](https://github.com/bpg/terraform-provider-proxmox/issues/87)

BUG FIXES:

- Powered off VM breaks plan/apply [\#105](https://github.com/bpg/terraform-provider-proxmox/issues/105)
- Disk resize causes reboot [\#102](https://github.com/bpg/terraform-provider-proxmox/issues/102)
- Typing error - dvResourceVirtualEnvironmentVMAgentEnabled instead of dvResourceVirtualEnvironmentVMAgentTrim [\#101](https://github.com/bpg/terraform-provider-proxmox/issues/101)
- Error creating VM with multiple disks on different storages [\#88](https://github.com/bpg/terraform-provider-proxmox/issues/88)

**Merged pull requests:**

- Fixed Typo  [\#107](https://github.com/bpg/terraform-provider-proxmox/pull/107) ([PrajwalBorkar](https://github.com/PrajwalBorkar))
- Avoid reboot when resizing disks. [\#104](https://github.com/bpg/terraform-provider-proxmox/pull/104) ([otopetrik](https://github.com/otopetrik))
- Add support for "iothread" disk option for VM [\#97](https://github.com/bpg/terraform-provider-proxmox/pull/97) ([bpg](https://github.com/bpg))
- Fix disk import when VM template has multiple disks [\#96](https://github.com/bpg/terraform-provider-proxmox/pull/96) ([bpg](https://github.com/bpg))

## [v0.5.3](https://github.com/bpg/terraform-provider-proxmox/tree/v0.5.3) (2022-06-02)

[Full Changelog](https://github.com/bpg/terraform-provider-proxmox/compare/v0.5.2...v0.5.3)

**Merged pull requests:**

- Bump hashicorp/go-getter for CVE-2022-30323 fix [\#82](https://github.com/bpg/terraform-provider-proxmox/pull/82) ([bpg](https://github.com/bpg))
- Update docs [\#57](https://github.com/bpg/terraform-provider-proxmox/pull/57) ([bpg](https://github.com/bpg))

## [v0.5.2](https://github.com/bpg/terraform-provider-proxmox/tree/v0.5.2) (2022-05-10)

[Full Changelog](https://github.com/bpg/terraform-provider-proxmox/compare/v0.5.1...v0.5.2)

**Closed issues:**

- proxmox\_virtual\_environment\_file json unmarshalling type issue [\#41](https://github.com/bpg/terraform-provider-proxmox/issues/41)

## [v0.5.1](https://github.com/bpg/terraform-provider-proxmox/tree/v0.5.1) (2022-03-22)

[Full Changelog](https://github.com/bpg/terraform-provider-proxmox/compare/v0.5.0...v0.5.1)

BUG FIXES:

- Version mismatch in the code [\#44](https://github.com/bpg/terraform-provider-proxmox/issues/44)
- virtual\_environment\_datastores.go: Update remote command to get datasource path [\#49](https://github.com/bpg/terraform-provider-proxmox/pull/49) ([mattburchett](https://github.com/mattburchett))

## [v0.5.0](https://github.com/bpg/terraform-provider-proxmox/tree/v0.5.0) (2021-11-06)

[Full Changelog](https://github.com/bpg/terraform-provider-proxmox/compare/v0.4.6...v0.5.0)

BREAKING CHANGES:

- Bump provider version to 0.5.0 [\#32](https://github.com/bpg/terraform-provider-proxmox/pull/32) ([bpg](https://github.com/bpg))

## [v0.4.6](https://github.com/bpg/terraform-provider-proxmox/tree/v0.4.6) (2021-09-10)

[Full Changelog](https://github.com/bpg/terraform-provider-proxmox/compare/v0.4.5...v0.4.6)

BUG FIXES:

- JSON unmarshal error when deploying LCX container [\#15](https://github.com/bpg/terraform-provider-proxmox/issues/15)
- \[BUG\] SIGSEGV if cloned VM disk is in the different storage [\#2](https://github.com/bpg/terraform-provider-proxmox/issues/2)
- fix `make test` error [\#1](https://github.com/bpg/terraform-provider-proxmox/pull/1) ([bpg](https://github.com/bpg))

## [v0.4.5](https://github.com/bpg/terraform-provider-proxmox/tree/v0.4.5) (2021-07-16)

[Full Changelog](https://github.com/bpg/terraform-provider-proxmox/compare/v0.4.4...v0.4.5)

## v0.4.4

BUG FIXES:

* resource/virtual_environment_vm: Fix watchdog deserialization issue

## v0.4.3

BUG FIXES:

* resource/virtual_environment_container: Fix IP initialization issue

## v0.4.2

BUG FIXES:

* resource/virtual_environment_vm: Fix `disk.file_id` diff issue
* resource/virtual_environment_vm: Fix disk resizing issue

OTHER:

* provider/example: Remove support for Terraform v0.11 and older
* provider/makefile: Update to use plugin caching to support local builds

## v0.4.1

OTHER:

* provider/docs: Fix issue with navigational link titles in Terraform Registry

## v0.4.0

FEATURES:

* **New Data Source:** `proxmox_virtual_environment_time`
* **New Resource:** `proxmox_virtual_environment_time`

BREAKING CHANGES:

* resource/virtual_environment_vm: `interface` is now required to create disks

    ```
      disk {
        datastore_id = "local-lvm"
        file_id      = "${proxmox_virtual_environment_file.ubuntu_cloud_image.id}"
        interface    = "scsi0"
      }
    ```

ENHANCEMENTS:

* provider/configuration: Add `virtual_environment.otp` argument for TOTP support
* resource/virtual_environment_vm: Clone supports resize and datastore_id for moving disks
* resource/virtual_environment_vm: Bulk clones can now use retries as argument to try multiple times to create a clone.
* resource/virtual_environment_vm: `on_boot` parameter can be used to start a VM after the Node has been rebooted.
* resource/virtual_environment_vm: `reboot` parameter can be used to reboot a VM after creation
* resource/virtual_environment_vm: Has now multiple new parameters to set timeouts for the vm creation/cloning `timeout_clone`, `timeout_move_disk`, `timeout_reboot`, `timeout_shutdown_vm`, `timeout_start_vm`, `timeout_stop_vm`

BUG FIXES:

* library/virtual_environment_nodes: Fix node IP address format
* library/virtual_environment_nodes: Fix WaitForNodeTask now detects errors correctly
* library/virtual_environment_vm: Fix CloneVM now waits for the task to be finished and detect errors.
* resource/virtual_environment_container: Fix VM ID collision when `vm_id` is not specified
* resource/virtual_environment_vm: Fix VM ID collision when `vm_id` is not specified
* resource/virtual_environment_vm: Fix disk import issue when importing from directory-based datastores
* resource/virtual_environment_vm: Fix handling of storage name - correct handling of `-`

WORKAROUNDS:

* resource/virtual_environment_vm: Ignore default value for `cpu.architecture` when the root account is not being used

## 0.3.0

ENHANCEMENTS:

* resource/virtual_environment_container: Add `clone` argument
* resource/virtual_environment_container: Add `disk` argument
* resource/virtual_environment_container: Add `template` argument
* resource/virtual_environment_vm: Add `agent.timeout` argument
* resource/virtual_environment_vm: Add `audio_device` argument
* resource/virtual_environment_vm: Add `clone` argument
* resource/virtual_environment_vm: Add `initialization.datastore_id` argument
* resource/virtual_environment_vm: Add `serial_device` argument
* resource/virtual_environment_vm: Add `template` argument

BUG FIXES:

* resource/virtual_environment_container: Fix `network_interface` deletion issue
* resource/virtual_environment_vm: Fix `network_device` deletion issue
* resource/virtual_environment_vm: Fix slow refresh when VM is stopped and agent is enabled
* resource/virtual_environment_vm: Fix crash caused by assuming IP addresses are always reported by the QEMU agent
* resource/virtual_environment_vm: Fix timeout issue while waiting for IP addresses to be reported by the QEMU agent

OTHER:

* provider/docs: Add HTML documentation powered by GitHub Pages

## 0.2.0

BREAKING CHANGES:

* resource/virtual_environment_vm: Rename `cloud_init` argument to `initialization`
* resource/virtual_environment_vm: Rename `os_type` argument to `operating_system.type`

FEATURES:

* **New Data Source:** `proxmox_virtual_environment_dns`
* **New Data Source:** `proxmox_virtual_environment_hosts`
* **New Resource:** `proxmox_virtual_environment_certificate`
* **New Resource:** `proxmox_virtual_environment_container`
* **New Resource:** `proxmox_virtual_environment_dns`
* **New Resource:** `proxmox_virtual_environment_hosts`

ENHANCEMENTS:

* resource/virtual_environment_vm: Add `acpi` argument
* resource/virtual_environment_vm: Add `bios` argument
* resource/virtual_environment_vm: Add `cpu.architecture`, `cpu.flags`, `cpu.type` and `cpu.units` arguments
* resource/virtual_environment_vm: Add `tablet_device` argument
* resource/virtual_environment_vm: Add `vga` argument

## 0.1.0

FEATURES:

* **New Data Source:** `proxmox_virtual_environment_datastores`
* **New Data Source:** `proxmox_virtual_environment_group`
* **New Data Source:** `proxmox_virtual_environment_groups`
* **New Data Source:** `proxmox_virtual_environment_nodes`
* **New Data Source:** `proxmox_virtual_environment_pool`
* **New Data Source:** `proxmox_virtual_environment_pools`
* **New Data Source:** `proxmox_virtual_environment_role`
* **New Data Source:** `proxmox_virtual_environment_roles`
* **New Data Source:** `proxmox_virtual_environment_user`
* **New Data Source:** `proxmox_virtual_environment_users`
* **New Data Source:** `proxmox_virtual_environment_version`
* **New Resource:** `proxmox_virtual_environment_file`
* **New Resource:** `proxmox_virtual_environment_group`
* **New Resource:** `proxmox_virtual_environment_pool`
* **New Resource:** `proxmox_virtual_environment_role`
* **New Resource:** `proxmox_virtual_environment_user`
* **New Resource:** `proxmox_virtual_environment_vm`


\* *This Changelog was automatically generated by [github_changelog_generator](https://github.com/github-changelog-generator/github-changelog-generator)*
