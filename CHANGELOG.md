# Changelog

## [0.14.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.14.0...v0.14.1) (2023-03-19)


### Bug Fixes

* authentication error logging in API client ([#267](https://github.com/bpg/terraform-provider-proxmox/issues/267)) ([763527e](https://github.com/bpg/terraform-provider-proxmox/commit/763527e53584e8121b1138830ad97e8e89780322))
* **build:** Fix make example-init for TF 1.4 ([#262](https://github.com/bpg/terraform-provider-proxmox/issues/262)) ([914631f](https://github.com/bpg/terraform-provider-proxmox/commit/914631f58b40ceb25248727ac23a6400df0264a3))


### Miscellaneous

* **deps:** bump activesupport from 6.1.7.1 to 6.1.7.3 in /docs ([#263](https://github.com/bpg/terraform-provider-proxmox/issues/263)) ([ce8bd30](https://github.com/bpg/terraform-provider-proxmox/commit/ce8bd3008bc65745eb62e17dc4849d3a4b3f740a))
* **docs:** Minor documentation Improvements ([#266](https://github.com/bpg/terraform-provider-proxmox/issues/266)) ([696ecb0](https://github.com/bpg/terraform-provider-proxmox/commit/696ecb05d8796540dc21d62dce930c4a2c2d8246))

## [0.14.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.13.1...v0.14.0) (2023-03-14)


### Features

* **lxc:** Add option for nested container feature ([4d44739](https://github.com/bpg/terraform-provider-proxmox/commit/4d447390e684a90c9672528f4bdc22aa1433296b))
* **vm:** Add custom CPU models support ([82016fc](https://github.com/bpg/terraform-provider-proxmox/commit/82016fc8ff018867783839c916dce686cb38d1b6))


### Bug Fixes

* **vm:** Fix `file_format` setting for new empty disks ([#259](https://github.com/bpg/terraform-provider-proxmox/issues/259)) ([d29fd97](https://github.com/bpg/terraform-provider-proxmox/commit/d29fd97babab9a8f217b6ea0ffd89511c55624eb))


### Miscellaneous

* **deps:** bump github.com/goreleaser/goreleaser from 1.15.2 to 1.16.1 in /tools ([#258](https://github.com/bpg/terraform-provider-proxmox/issues/258)) ([9afca3b](https://github.com/bpg/terraform-provider-proxmox/commit/9afca3b88caade184e536450534666431f2c00d5))

## [0.13.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.13.0...v0.13.1) (2023-03-07)


### Miscellaneous

* **deps:** bump dependencies ([#242](https://github.com/bpg/terraform-provider-proxmox/issues/242)) ([890fb53](https://github.com/bpg/terraform-provider-proxmox/commit/890fb536846d2582cbf025f2045be3c5f903fc0a))
* **deps:** bump github.com/golangci/golangci-lint from 1.51.1 to 1.51.2 in /tools ([#244](https://github.com/bpg/terraform-provider-proxmox/issues/244)) ([e01844a](https://github.com/bpg/terraform-provider-proxmox/commit/e01844a0d73750d0ce65c76e9eaae0b3b952c206))
* **deps:** bump github.com/stretchr/testify from 1.8.1 to 1.8.2 ([#245](https://github.com/bpg/terraform-provider-proxmox/issues/245)) ([6cca133](https://github.com/bpg/terraform-provider-proxmox/commit/6cca13383527a1f33a30e5766bb520c0a575793a))
* **deps:** bump golang.org/x/crypto from 0.6.0 to 0.7.0 ([#248](https://github.com/bpg/terraform-provider-proxmox/issues/248)) ([1aa668e](https://github.com/bpg/terraform-provider-proxmox/commit/1aa668e03bcb15333772575029a07c2134d8e291))

## [0.13.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.12.1...v0.13.0) (2023-02-17)


### Features

* **vm:** update VM disc import logic ([#241](https://github.com/bpg/terraform-provider-proxmox/issues/241)) ([fcf9810](https://github.com/bpg/terraform-provider-proxmox/commit/fcf98102522821c9dfb4534731747233bd627d38))


### Bug Fixes

* **vm:** `proxmox_virtual_environment_file.changed` stored as `true` at file creation ([#240](https://github.com/bpg/terraform-provider-proxmox/issues/240)) ([197c9e5](https://github.com/bpg/terraform-provider-proxmox/commit/197c9e5152fd6524c82977001a759c36c644f8e5))


### Miscellaneous

* **deps:** bump activesupport from 6.0.6.1 to 6.1.7.1 in /docs ([#235](https://github.com/bpg/terraform-provider-proxmox/issues/235)) ([ffa39c1](https://github.com/bpg/terraform-provider-proxmox/commit/ffa39c13e0d8283da51980532c83919edcf1cbc6))
* **deps:** bump github.com/goreleaser/goreleaser from 1.15.1 to 1.15.2 in /tools ([#237](https://github.com/bpg/terraform-provider-proxmox/issues/237)) ([80dfceb](https://github.com/bpg/terraform-provider-proxmox/commit/80dfceba8433379a64a1ff86d174447e229325ab))
* **deps:** bump github.com/hashicorp/terraform-plugin-log from 0.7.0 to 0.8.0 ([#239](https://github.com/bpg/terraform-provider-proxmox/issues/239)) ([dbe26ed](https://github.com/bpg/terraform-provider-proxmox/commit/dbe26ed58f1ed668e5a059f9659bd12fd6b1a54c))
* **deps:** bump golang.org/x/crypto from 0.5.0 to 0.6.0 ([#238](https://github.com/bpg/terraform-provider-proxmox/issues/238)) ([2b99349](https://github.com/bpg/terraform-provider-proxmox/commit/2b99349f1fe89e804fb45c439470bd2474068f1c))

## [0.12.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.12.0...v0.12.1) (2023-02-07)


### Bug Fixes

* **build:** trailing space in provider's version ([#233](https://github.com/bpg/terraform-provider-proxmox/issues/233)) ([f97407d](https://github.com/bpg/terraform-provider-proxmox/commit/f97407dc00c425b8d015abf72488b5a4fd31f043))
* **vm:** ignore ssd disk flag with virtio interface ([#231](https://github.com/bpg/terraform-provider-proxmox/issues/231)) ([1de9294](https://github.com/bpg/terraform-provider-proxmox/commit/1de92947666d45fdcae881e3a6bd651bfea493a4))


### Miscellaneous

* **deps:** bump github.com/golangci/golangci-lint from 1.50.1 to 1.51.1 in /tools ([#229](https://github.com/bpg/terraform-provider-proxmox/issues/229)) ([f1022a5](https://github.com/bpg/terraform-provider-proxmox/commit/f1022a5cae0c99696292421edb28b3007d3bbb51))
* **deps:** bump github.com/goreleaser/goreleaser from 1.14.1 to 1.15.1 in /tools ([#230](https://github.com/bpg/terraform-provider-proxmox/issues/230)) ([722003e](https://github.com/bpg/terraform-provider-proxmox/commit/722003ee5ac23c4946af2257eaeb6f91028f879d))

## [0.12.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.11.0...v0.12.0) (2023-02-06)


### Features

* **core:** Add known hosts callback check for ssh connections ([#217](https://github.com/bpg/terraform-provider-proxmox/issues/217)) ([598c628](https://github.com/bpg/terraform-provider-proxmox/commit/598c62864d0a8a4e1b7dcda0cb7a3d5e380a5863))
* **lxc:** Add unprivileged option ([#225](https://github.com/bpg/terraform-provider-proxmox/issues/225)) ([1918561](https://github.com/bpg/terraform-provider-proxmox/commit/19185611b37c85a071ac4d3fd0c9a6b865b7c28d))


### Bug Fixes

* **vm:** Don't add an extra hostpci entry ([#223](https://github.com/bpg/terraform-provider-proxmox/issues/223)) ([346c92b](https://github.com/bpg/terraform-provider-proxmox/commit/346c92b2734caed90b30df423ac8019cf08c5900))
* **vm:** Fix handling of empty kvm arguments ([#228](https://github.com/bpg/terraform-provider-proxmox/issues/228)) ([e2802d0](https://github.com/bpg/terraform-provider-proxmox/commit/e2802d0654f0d6d5e99bef4987a84862e3ffbde7))


### Miscellaneous

* **deps:** bump commonmarker from 0.23.6 to 0.23.7 in /docs ([#220](https://github.com/bpg/terraform-provider-proxmox/issues/220)) ([cef0227](https://github.com/bpg/terraform-provider-proxmox/commit/cef0227ef59df55388632e775b34cc3f4644075f))
* **deps:** bump gem dependencies in /docs ([#221](https://github.com/bpg/terraform-provider-proxmox/issues/221)) ([e0864c0](https://github.com/bpg/terraform-provider-proxmox/commit/e0864c098a2e5a9d1da1c133ebaeee8650ceb4e0))
* **deps:** bump goreleaser/goreleaser-action from 4.1.0 to 4.2.0 ([#222](https://github.com/bpg/terraform-provider-proxmox/issues/222)) ([11fe9e5](https://github.com/bpg/terraform-provider-proxmox/commit/11fe9e539c56101360e5be0f5bb042f5318a4d4c))
* disable code coverage ([#227](https://github.com/bpg/terraform-provider-proxmox/issues/227)) ([a72fd27](https://github.com/bpg/terraform-provider-proxmox/commit/a72fd27a13395b9d061cdc450f68e171f1b30cbe))

## [0.11.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.10.0...v0.11.0) (2023-01-24)


### Features

* **lxc:** Add support for container tags ([#212](https://github.com/bpg/terraform-provider-proxmox/issues/212)) ([5c8ae3c](https://github.com/bpg/terraform-provider-proxmox/commit/5c8ae3c3f826969f70d5af79cfca00c0c49da418))


### Miscellaneous

* **ci:** set up code coverage ([06bd5ae](https://github.com/bpg/terraform-provider-proxmox/commit/06bd5aef0f0aac54e412e475ccdc85f8f61398d9))
* **deps:** bump dependencies ([#216](https://github.com/bpg/terraform-provider-proxmox/issues/216)) ([f157e3b](https://github.com/bpg/terraform-provider-proxmox/commit/f157e3bd532bd3b0fa728478f44986b1ef5f245b))

## [0.10.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.9.1...v0.10.0) (2023-01-18)


### Features

* **lxc:** Add option to customize RootFS size at LXC creation ([#207](https://github.com/bpg/terraform-provider-proxmox/issues/207)) ([dd9ffe1](https://github.com/bpg/terraform-provider-proxmox/commit/dd9ffe190cd9eaee7ac6a9e2c830eee45b4b69df))
* **vm:** add support for "args" flag for VM ([#205](https://github.com/bpg/terraform-provider-proxmox/issues/205)) ([8bd3fd7](https://github.com/bpg/terraform-provider-proxmox/commit/8bd3fd7b1d71e37eeee2c222e4896b857a01cabf))


### Bug Fixes

* **vm:** Add parser for CustomEFIDisk ([#208](https://github.com/bpg/terraform-provider-proxmox/issues/208)) ([b539aab](https://github.com/bpg/terraform-provider-proxmox/commit/b539aab22851817aea981727eb27a8da73edcc43))

## [0.9.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.9.0...v0.9.1) (2023-01-02)


### Bug Fixes

* **vm:** Make so that on_boot can be changed with update ([#199](https://github.com/bpg/terraform-provider-proxmox/issues/199)) ([496ab32](https://github.com/bpg/terraform-provider-proxmox/commit/496ab322be7f12257f562d53a9920377cded8aa5))

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
