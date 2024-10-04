# Changelog

## [0.66.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.65.0...v0.66.0) (2024-10-04)


### Features

* **provider:** add support for pre(external) auth'd session tokens ([#1441](https://github.com/bpg/terraform-provider-proxmox/issues/1441)) ([eb2f36b](https://github.com/bpg/terraform-provider-proxmox/commit/eb2f36be21bc8a74bbd7e74bc7a7f3f6ed2d8daf))
* **provider:** reliable sequential and random `vm_id` generation ([#1557](https://github.com/bpg/terraform-provider-proxmox/issues/1557)) ([72f7cb8](https://github.com/bpg/terraform-provider-proxmox/commit/72f7cb81a8bbe0fcd9aa739073db6e552820a7de))
* **vm:** add support for `watchdog` ([#1556](https://github.com/bpg/terraform-provider-proxmox/issues/1556)) ([d226b59](https://github.com/bpg/terraform-provider-proxmox/commit/d226b59e2e37ffebc737e27cc9bb0182d4bda993))


### Bug Fixes

* **docs:** correct url for template download in container resource example ([#1563](https://github.com/bpg/terraform-provider-proxmox/issues/1563)) ([7266583](https://github.com/bpg/terraform-provider-proxmox/commit/72665830be975f0671c6aecca47fc52906cb21e2))
* **lxc:** add new storage-backed mount point to existing container ([#1553](https://github.com/bpg/terraform-provider-proxmox/issues/1553)) ([89d72cd](https://github.com/bpg/terraform-provider-proxmox/commit/89d72cd37f63fd6238dbb09709ce90020e42e18a))
* **provider:** handle IPv6 in ssh client ([#1558](https://github.com/bpg/terraform-provider-proxmox/issues/1558)) ([18a7f8e](https://github.com/bpg/terraform-provider-proxmox/commit/18a7f8ec358a158ccc4f927e0c07b5eab77b01dd))


### Miscellaneous

* **ci:** update actions/checkout action (v4.1.7 → v4.2.0) ([#1555](https://github.com/bpg/terraform-provider-proxmox/issues/1555)) ([471c5f2](https://github.com/bpg/terraform-provider-proxmox/commit/471c5f23bb67323abab77e57a359e07050ac536c))
* **ci:** update peter-evans/create-issue-from-file action (v5.0.0 → v5.0.1) ([#1554](https://github.com/bpg/terraform-provider-proxmox/issues/1554)) ([049362a](https://github.com/bpg/terraform-provider-proxmox/commit/049362a94e585c17359a61ad410ad0c20b207d62))
* **docs:** document specifics of "ballooning device" in the `vm.memory` block ([#1550](https://github.com/bpg/terraform-provider-proxmox/issues/1550)) ([0b29ae2](https://github.com/bpg/terraform-provider-proxmox/commit/0b29ae281ed5075aa1349fbe262176d8f8791dfc))

## [0.65.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.64.0...v0.65.0) (2024-09-21)


### Features

* **vm:** remove deprecated `vga.enabled` ([#1533](https://github.com/bpg/terraform-provider-proxmox/issues/1533)) ([f4ee629](https://github.com/bpg/terraform-provider-proxmox/commit/f4ee6291ccd42ef6c6c7d0e8d6780f68a46ac218))


### Bug Fixes

* **file:** respect `download_file.override` attr when checking for the file size change ([#1537](https://github.com/bpg/terraform-provider-proxmox/issues/1537)) ([44e2814](https://github.com/bpg/terraform-provider-proxmox/commit/44e2814adf69588e219c522fa415fa38cef9dc3a))
* **vm:** handle more than 4 `hostpci` devices ([#1543](https://github.com/bpg/terraform-provider-proxmox/issues/1543)) ([569db6d](https://github.com/bpg/terraform-provider-proxmox/commit/569db6df3bcacd86d0c1dab77512ca73d7b04560))


### Miscellaneous

* **ci:** update actions/create-github-app-token action (v1.10.4 → v1.11.0) ([#1532](https://github.com/bpg/terraform-provider-proxmox/issues/1532)) ([831207b](https://github.com/bpg/terraform-provider-proxmox/commit/831207b0e30ed0b91ec612681c130fb412210e61))
* **ci:** update jetbrains/qodana-action action (v2024.1.9 → v2024.2.3) ([#1544](https://github.com/bpg/terraform-provider-proxmox/issues/1544)) ([13ceab8](https://github.com/bpg/terraform-provider-proxmox/commit/13ceab820501933c5b3b5d24aae1e9db69b2816b))
* **deps:** update github.com/hashicorp/terraform-plugin-* ([#1545](https://github.com/bpg/terraform-provider-proxmox/issues/1545)) ([b5a4e32](https://github.com/bpg/terraform-provider-proxmox/commit/b5a4e32d7a9b7dc0328e114b60eaf42875d4fab6))
* **deps:** update golang.org/x/exp digest (e7e105d → 701f63a) ([#1526](https://github.com/bpg/terraform-provider-proxmox/issues/1526)) ([b96532f](https://github.com/bpg/terraform-provider-proxmox/commit/b96532f884af73b562199d3b049b4ee4251fbf78))
* **deps:** update module github.com/golangci/golangci-lint (v1.60.3 → v1.61.0) in /tools ([#1527](https://github.com/bpg/terraform-provider-proxmox/issues/1527)) ([356c978](https://github.com/bpg/terraform-provider-proxmox/commit/356c978ece6c60e3b130d753f9bcd644d8507fce))
* **docs:** add example for host disk pass-through to VM ([#1547](https://github.com/bpg/terraform-provider-proxmox/issues/1547)) ([e23d2f3](https://github.com/bpg/terraform-provider-proxmox/commit/e23d2f375cb72a1e8a7f67490cbcbe1a3ea5d86f))
* **docs:** improve `network_linux_bridge` documentation ([#1538](https://github.com/bpg/terraform-provider-proxmox/issues/1538)) ([8f9b036](https://github.com/bpg/terraform-provider-proxmox/commit/8f9b03685791088353b7ac0872130ebc7c832626))
* **docs:** update terraform local (2.5.1 → 2.5.2) ([#1535](https://github.com/bpg/terraform-provider-proxmox/issues/1535)) ([f1ce62e](https://github.com/bpg/terraform-provider-proxmox/commit/f1ce62eeeb26d37ae1be2580869cd5baa95741be))
* **docs:** update terraform tls (4.0.5 → 4.0.6) ([#1536](https://github.com/bpg/terraform-provider-proxmox/issues/1536)) ([a89cfbf](https://github.com/bpg/terraform-provider-proxmox/commit/a89cfbf045f2bddde3ac0d769686434467bb61d6))
* **lxc:** update list of supported OS-es in `operating_system.type` ([#1548](https://github.com/bpg/terraform-provider-proxmox/issues/1548)) ([14a62a7](https://github.com/bpg/terraform-provider-proxmox/commit/14a62a712ae4109bc5345b6e09248bd5e68002b5))

## [0.64.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.63.0...v0.64.0) (2024-09-08)


### Features

* **acme:** implement resources and data sources for ACME plugins ([#1479](https://github.com/bpg/terraform-provider-proxmox/issues/1479)) ([a6eb81a](https://github.com/bpg/terraform-provider-proxmox/commit/a6eb81af08f175ab078ae0369e0246425e77de6a))
* **lxc:** Implement support for container's "protection flag" ([#1512](https://github.com/bpg/terraform-provider-proxmox/issues/1512)) ([9b34c48](https://github.com/bpg/terraform-provider-proxmox/commit/9b34c485ab6ee8e29e8df2269d7789820a980558))


### Bug Fixes

* **docs:** improve "SSH User" section describing `sudo` privileges ([#1522](https://github.com/bpg/terraform-provider-proxmox/issues/1522)) ([84d55cf](https://github.com/bpg/terraform-provider-proxmox/commit/84d55cf5690df0d7a5231e2d664f5fec38064df8))
* **vm:** `cpu.architecture` showed as new attribute at re-apply after creation ([#1524](https://github.com/bpg/terraform-provider-proxmox/issues/1524)) ([c20d79d](https://github.com/bpg/terraform-provider-proxmox/commit/c20d79dfbe34b4f161d49a6c9b595354444bacd9))


### Miscellaneous

* **ci:** update actions/create-github-app-token action (v1.10.3 → v1.10.4) ([#1518](https://github.com/bpg/terraform-provider-proxmox/issues/1518)) ([112f058](https://github.com/bpg/terraform-provider-proxmox/commit/112f058e66a5bb410696d1d40b11e7216ab3986e))
* **deps:** update golang.org/x/exp digest (0cdaa3a → 9b4947d) ([#1506](https://github.com/bpg/terraform-provider-proxmox/issues/1506)) ([0623bfd](https://github.com/bpg/terraform-provider-proxmox/commit/0623bfd1561a7a30771685d8fa3f46c9bfb4d1d8))
* **deps:** update golang.org/x/exp digest (9b4947d → e7e105d) ([#1517](https://github.com/bpg/terraform-provider-proxmox/issues/1517)) ([273b3c9](https://github.com/bpg/terraform-provider-proxmox/commit/273b3c93b3ffafa448076ca07f6d53066844f164))
* **deps:** update module github.com/golangci/golangci-lint (v1.60.1 → v1.60.2) in /tools ([#1499](https://github.com/bpg/terraform-provider-proxmox/issues/1499)) ([35cbe98](https://github.com/bpg/terraform-provider-proxmox/commit/35cbe98d6b26b47c0b8797dd6d4f084585f846c7))
* **deps:** update module github.com/golangci/golangci-lint (v1.60.2 → v1.60.3) in /tools ([#1507](https://github.com/bpg/terraform-provider-proxmox/issues/1507)) ([5707fea](https://github.com/bpg/terraform-provider-proxmox/commit/5707feab06750a4782e66c0be0e5cc2edd612aa0))
* **deps:** update module go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp (v0.53.0 → v0.54.0) in /tools ([#1509](https://github.com/bpg/terraform-provider-proxmox/issues/1509)) ([9172ec6](https://github.com/bpg/terraform-provider-proxmox/commit/9172ec6a9fd01e876e18fa3e8664aee3b639c891))
* **deps:** update module golang.org/x/crypto (v0.26.0 → v0.27.0) ([#1519](https://github.com/bpg/terraform-provider-proxmox/issues/1519)) ([3951bda](https://github.com/bpg/terraform-provider-proxmox/commit/3951bdab4a958761e6978ab0dc03857b797bd417))
* **deps:** update module golang.org/x/net (v0.28.0 → v0.29.0) ([#1520](https://github.com/bpg/terraform-provider-proxmox/issues/1520)) ([a0bce8d](https://github.com/bpg/terraform-provider-proxmox/commit/a0bce8d699b6d633acb309c55e43e8743afb7122))
* **vm:** add acc test for no cdrom by default ([#1523](https://github.com/bpg/terraform-provider-proxmox/issues/1523)) ([35b6571](https://github.com/bpg/terraform-provider-proxmox/commit/35b6571bf468d09902c899f3e1782db9a13b78a0))

## [0.63.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.62.0...v0.63.0) (2024-08-21)


### Features

* **file:** add support to set the file mode ([#1478](https://github.com/bpg/terraform-provider-proxmox/issues/1478)) ([cc9d0e7](https://github.com/bpg/terraform-provider-proxmox/commit/cc9d0e713190d736861883f97058c1ec7f0a179d))


### Bug Fixes

* **doc:** update virtual_environment_vm.md ([#1495](https://github.com/bpg/terraform-provider-proxmox/issues/1495)) ([4f8e18d](https://github.com/bpg/terraform-provider-proxmox/commit/4f8e18da04419c0077d53c3eaf90e5d18a8d947f))
* **firewall:** remove limit for number of CIRDs in `ipset` ([#1487](https://github.com/bpg/terraform-provider-proxmox/issues/1487)) ([55bfe14](https://github.com/bpg/terraform-provider-proxmox/commit/55bfe14ce18a432849e544360561d7b9091ac4d8))


### Miscellaneous

* **ci:** update hashicorp/setup-terraform action (v3.1.1 → v3.1.2) ([#1498](https://github.com/bpg/terraform-provider-proxmox/issues/1498)) ([3f6bc5e](https://github.com/bpg/terraform-provider-proxmox/commit/3f6bc5eba223c789e9b87bf74af6ec8c49bb6b60))
* **ci:** update jetbrains/qodana-action action (v2024.1.8 → v2024.1.9) ([#1484](https://github.com/bpg/terraform-provider-proxmox/issues/1484)) ([dd6d5ec](https://github.com/bpg/terraform-provider-proxmox/commit/dd6d5ec88ffefc41cf68aeac6260915a05450463))
* **code:** align resource/datasource names in the fwprovider code ([#1488](https://github.com/bpg/terraform-provider-proxmox/issues/1488)) ([8f82d1a](https://github.com/bpg/terraform-provider-proxmox/commit/8f82d1a3848571f09aea400960b252d4b6c39fb1))
* **deps:** update go (1.22.6 → 1.23.0) ([#1489](https://github.com/bpg/terraform-provider-proxmox/issues/1489)) ([3343a93](https://github.com/bpg/terraform-provider-proxmox/commit/3343a938409adc2d61658e7d2c2a52b452c50d48))
* **deps:** update golang.org/x/exp digest (8a7402a → 0cdaa3a) ([#1483](https://github.com/bpg/terraform-provider-proxmox/issues/1483)) ([e566c4e](https://github.com/bpg/terraform-provider-proxmox/commit/e566c4e35e6680c2d24cbc42fb78aa1e753f1dd1))
* **deps:** update module github.com/hashicorp/terraform-plugin-testing (v1.9.0 → v1.10.0) ([#1485](https://github.com/bpg/terraform-provider-proxmox/issues/1485)) ([767f680](https://github.com/bpg/terraform-provider-proxmox/commit/767f680cc3428484a6c5fe15c22c6ab7acf1dce8))
* **deps:** update tools ([#1490](https://github.com/bpg/terraform-provider-proxmox/issues/1490)) ([6635c89](https://github.com/bpg/terraform-provider-proxmox/commit/6635c89f632f61041382ea08344380f777ee54ba))
* **devcontainer:** configure provider override in post attach command ([#1481](https://github.com/bpg/terraform-provider-proxmox/issues/1481)) ([daaedf2](https://github.com/bpg/terraform-provider-proxmox/commit/daaedf2d99ee56a1873aeb6da4af94eb06fc7d5f))
* **example:** update links to container images, use faster mirror ([#1497](https://github.com/bpg/terraform-provider-proxmox/issues/1497)) ([abb0be4](https://github.com/bpg/terraform-provider-proxmox/commit/abb0be4684fa21dd5e7b600f238348c509b6c9af))

## [0.62.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.61.1...v0.62.0) (2024-08-08)


### Features

* **acme:** implement resources and data sources for ACME accounts ([#1455](https://github.com/bpg/terraform-provider-proxmox/issues/1455)) ([9de4037](https://github.com/bpg/terraform-provider-proxmox/commit/9de4037a8222fbf9d524df989d0b670c0ae13771))


### Bug Fixes

* **docs:** document role privileges required by `proxmox_virtual_environment_download_file` ([#1442](https://github.com/bpg/terraform-provider-proxmox/issues/1442)) ([2475d99](https://github.com/bpg/terraform-provider-proxmox/commit/2475d99317bbef3f4a2372ed1377c64590ad5df8))
* **storage:** ignore `os.ErrClosed` from deferred `fileReader.Close()` in storage.APIUpload ([#1468](https://github.com/bpg/terraform-provider-proxmox/issues/1468)) ([e8bd1fc](https://github.com/bpg/terraform-provider-proxmox/commit/e8bd1fcda7c60de485714af6ab788f6e6a99d28f))
* **tools:** ensure docs/guides files are copied correctly ([#1464](https://github.com/bpg/terraform-provider-proxmox/issues/1464)) ([f7c3560](https://github.com/bpg/terraform-provider-proxmox/commit/f7c35602608a7cfcdebceadd725ebae02c9b6cf6))


### Miscellaneous

* **deps:** update go (1.22.5 → 1.22.6) ([#1469](https://github.com/bpg/terraform-provider-proxmox/issues/1469)) ([966ea7a](https://github.com/bpg/terraform-provider-proxmox/commit/966ea7aa40e89990177ead67598a642bce355b27))
* **deps:** update go (1.22.5 → 1.22.6) in /tools ([#1470](https://github.com/bpg/terraform-provider-proxmox/issues/1470)) ([10b2340](https://github.com/bpg/terraform-provider-proxmox/commit/10b2340c5b820260ffd3efdb930f830783aea330))
* **deps:** update golang.org/x/exp digest (46b0784 → e3f2596) ([#1448](https://github.com/bpg/terraform-provider-proxmox/issues/1448)) ([1ba422a](https://github.com/bpg/terraform-provider-proxmox/commit/1ba422a68ef86ef8903aafd0030a4dcdb76023e1))
* **deps:** update golang.org/x/exp digest (e3f2596 → 8a7402a) ([#1452](https://github.com/bpg/terraform-provider-proxmox/issues/1452)) ([76ce03f](https://github.com/bpg/terraform-provider-proxmox/commit/76ce03f6142f5811400ba3f3724a9a35decbefc3))
* **deps:** update module github.com/hashicorp/terraform-plugin-framework (v1.10.0 → v1.11.0) ([#1471](https://github.com/bpg/terraform-provider-proxmox/issues/1471)) ([b7a110e](https://github.com/bpg/terraform-provider-proxmox/commit/b7a110e8fbe341f4347f517153ef36da286e1ae8))
* **deps:** update module github.com/skeema/knownhosts (v1.2.2 → v1.3.0) ([#1449](https://github.com/bpg/terraform-provider-proxmox/issues/1449)) ([f116127](https://github.com/bpg/terraform-provider-proxmox/commit/f1161276ad4106033a88de0ed99ceb1258817ea7))
* **deps:** update module golang.org/x/crypto (v0.25.0 → v0.26.0) ([#1472](https://github.com/bpg/terraform-provider-proxmox/issues/1472)) ([e641b55](https://github.com/bpg/terraform-provider-proxmox/commit/e641b55734944ebe8b359c158e5ea2ee798f433d))
* **deps:** update module golang.org/x/net (v0.27.0 → v0.28.0) ([#1474](https://github.com/bpg/terraform-provider-proxmox/issues/1474)) ([56b98bd](https://github.com/bpg/terraform-provider-proxmox/commit/56b98bd4e698ab26aef773dcd7ec2aa94be9ae35))
* **docs:** minor documentation updates ([#1463](https://github.com/bpg/terraform-provider-proxmox/issues/1463)) ([2bf9a94](https://github.com/bpg/terraform-provider-proxmox/commit/2bf9a9456210b7b0f60f7a5059b3dca108eef7b3))

## [0.61.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.61.0...v0.61.1) (2024-07-11)


### Bug Fixes

* **time:** embed timezone data to fix set/get TZ on windows  ([#1436](https://github.com/bpg/terraform-provider-proxmox/issues/1436)) ([b4d4e4f](https://github.com/bpg/terraform-provider-proxmox/commit/b4d4e4f9d4d1c38a5954e76beb0067f20652da19))
* **vm:** improve reliability of VM create / get operations ([#1431](https://github.com/bpg/terraform-provider-proxmox/issues/1431)) ([d193abd](https://github.com/bpg/terraform-provider-proxmox/commit/d193abd33e8f09aae0dc00d2405297908b21b7d5))
* **vm:** panic when reading VM with SPICE USB passthrough ([#1438](https://github.com/bpg/terraform-provider-proxmox/issues/1438)) ([fc42ac3](https://github.com/bpg/terraform-provider-proxmox/commit/fc42ac3a5e49daf391af21ec7cce7d97ecf8af20))


### Miscellaneous

* **deps:** update github.com/hashicorp/terraform-plugin-* ([#1437](https://github.com/bpg/terraform-provider-proxmox/issues/1437)) ([b63c98b](https://github.com/bpg/terraform-provider-proxmox/commit/b63c98b13f7f4dde06f3fbc873939be9898571f8))
* **deps:** update golang.org/x/exp digest (7f521ea → 46b0784) ([#1432](https://github.com/bpg/terraform-provider-proxmox/issues/1432)) ([0a978e5](https://github.com/bpg/terraform-provider-proxmox/commit/0a978e566417156b55a5973b2535a5314db992fd))

## [0.61.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.60.1...v0.61.0) (2024-07-06)


### Features

* **node:** implement initial support to manage APT repositories ([#1325](https://github.com/bpg/terraform-provider-proxmox/issues/1325)) ([357f7c7](https://github.com/bpg/terraform-provider-proxmox/commit/357f7c70a701c42dda62e0adef1337d319aaed04))
* **vm:** implement filtering in vms data source. ([#1423](https://github.com/bpg/terraform-provider-proxmox/issues/1423)) ([65f8ba5](https://github.com/bpg/terraform-provider-proxmox/commit/65f8ba5bfe3a8ca3872fd87830b1a858f24a7237))


### Bug Fixes

* **repo:** fix source.list path validation on Windows ([#1429](https://github.com/bpg/terraform-provider-proxmox/issues/1429)) ([6103e67](https://github.com/bpg/terraform-provider-proxmox/commit/6103e67ae97434991e1d85d368dbea738f6c8041))


### Miscellaneous

* **ci:** update actions/create-github-app-token action (v1.10.1 → v1.10.2) ([#1416](https://github.com/bpg/terraform-provider-proxmox/issues/1416)) ([692ea94](https://github.com/bpg/terraform-provider-proxmox/commit/692ea944c34d85ec54d40c8b1d5f15ed887f42d9))
* **ci:** update actions/create-github-app-token action (v1.10.2 → v1.10.3) ([#1419](https://github.com/bpg/terraform-provider-proxmox/issues/1419)) ([877251c](https://github.com/bpg/terraform-provider-proxmox/commit/877251c86eaa01aa1cce3f69afeb0e0214e87486))
* **ci:** update jetbrains/qodana-action action (v2024.1.5 → v2024.1.8) ([#1425](https://github.com/bpg/terraform-provider-proxmox/issues/1425)) ([91a16af](https://github.com/bpg/terraform-provider-proxmox/commit/91a16af747332215b9bfccc58cccb13a478222b2))
* **deps:** update go (1.22.4 → 1.22.5) ([#1420](https://github.com/bpg/terraform-provider-proxmox/issues/1420)) ([a521c60](https://github.com/bpg/terraform-provider-proxmox/commit/a521c60f0c7e9b36eac94ac78c9cfac83856e294))
* **deps:** update go (1.22.4 → 1.22.5) in /tools ([#1421](https://github.com/bpg/terraform-provider-proxmox/issues/1421)) ([d55f670](https://github.com/bpg/terraform-provider-proxmox/commit/d55f6706a49232869efd7319a6d3fda14e9e9372))
* **deps:** update module go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp (v0.52.0 → v0.53.0) in /tools ([#1426](https://github.com/bpg/terraform-provider-proxmox/issues/1426)) ([cf1142e](https://github.com/bpg/terraform-provider-proxmox/commit/cf1142e24b4c308637bfa481107f82d149a3b585))
* **deps:** update module golang.org/x/crypto (v0.24.0 → v0.25.0) ([#1427](https://github.com/bpg/terraform-provider-proxmox/issues/1427)) ([3ad5b0a](https://github.com/bpg/terraform-provider-proxmox/commit/3ad5b0a9e10ff8d6a5e707dd2dfb2a9f0d70bd4b))
* **deps:** update module golang.org/x/net (v0.26.0 → v0.27.0) ([#1430](https://github.com/bpg/terraform-provider-proxmox/issues/1430)) ([6436967](https://github.com/bpg/terraform-provider-proxmox/commit/6436967903970e7cdece9f85bbddc00afa9ce8f9))

## [0.60.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.60.0...v0.60.1) (2024-06-25)


### Bug Fixes

* **docs:** Mention alternate valid use of Volid's for containers ([#1407](https://github.com/bpg/terraform-provider-proxmox/issues/1407)) ([763654a](https://github.com/bpg/terraform-provider-proxmox/commit/763654ac8f88787de7c4d64e18aaf33e82990719))
* **lxc:** use default rootfs size (4Gb) prevents creation of mount points ([#1398](https://github.com/bpg/terraform-provider-proxmox/issues/1398)) ([4a8bf8d](https://github.com/bpg/terraform-provider-proxmox/commit/4a8bf8da27295d6fcca4be8cae05aedfe48d7cf3))


### Miscellaneous

* **code:** fix code quality warnings ([#1409](https://github.com/bpg/terraform-provider-proxmox/issues/1409)) ([4f70459](https://github.com/bpg/terraform-provider-proxmox/commit/4f70459c264e1d9d3e6541fa7276d5de57fa5d46))
* **deps:** update module github.com/brianvoe/gofakeit/v7 (v7.0.3 → v7.0.4) ([#1405](https://github.com/bpg/terraform-provider-proxmox/issues/1405)) ([c926484](https://github.com/bpg/terraform-provider-proxmox/commit/c926484249af93b75c64d2acba3e1460f16f35da))
* refactor container acc test ([#1408](https://github.com/bpg/terraform-provider-proxmox/issues/1408)) ([a0d9300](https://github.com/bpg/terraform-provider-proxmox/commit/a0d9300f0ff1f49b4df08691f4c2286b276f84b5))
* **vm:** fix regression after storage devices refactoring ([#1399](https://github.com/bpg/terraform-provider-proxmox/issues/1399)) ([5d1c8c6](https://github.com/bpg/terraform-provider-proxmox/commit/5d1c8c606d61d91360325228655ccd372848602a))
* **vm:** refactor storage devices handling from/to API ([#1394](https://github.com/bpg/terraform-provider-proxmox/issues/1394)) ([cc7fc63](https://github.com/bpg/terraform-provider-proxmox/commit/cc7fc63ec1c167344191cd2d3496139747e5c112))

## [0.60.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.59.1...v0.60.0) (2024-06-14)


### Features

* **vm:** add support for `disk.serial` attribute  ([#1385](https://github.com/bpg/terraform-provider-proxmox/issues/1385)) ([de34952](https://github.com/bpg/terraform-provider-proxmox/commit/de349523fedac8185468620a7088c55deefdbd79))


### Bug Fixes

* **docs:** add more details about the `usb` block in VM ([#1382](https://github.com/bpg/terraform-provider-proxmox/issues/1382)) ([1eb64a2](https://github.com/bpg/terraform-provider-proxmox/commit/1eb64a25226d6c2ae48ac5cc4c24499e0bbb84cf))
* **docs:** update virtual_environment_vm.md ([#1377](https://github.com/bpg/terraform-provider-proxmox/issues/1377)) ([63f0cbf](https://github.com/bpg/terraform-provider-proxmox/commit/63f0cbfd27387a9d4fd0b5aa853be9a5a52475f2))
* **lxc:** make container vm_id computed ([#1386](https://github.com/bpg/terraform-provider-proxmox/issues/1386)) ([bc079c0](https://github.com/bpg/terraform-provider-proxmox/commit/bc079c0c192d1b3cb99c2b646f639a57cfbff7f6))


### Miscellaneous

* **ci:** update actions/checkout action (v4.1.6 → v4.1.7) ([#1390](https://github.com/bpg/terraform-provider-proxmox/issues/1390)) ([fbf41dd](https://github.com/bpg/terraform-provider-proxmox/commit/fbf41dde93ff344950e4b879102b0070406bb227))
* **ci:** update googleapis/release-please-action action (v4.1.1 → v4.1.3) ([#1380](https://github.com/bpg/terraform-provider-proxmox/issues/1380)) ([77f62e1](https://github.com/bpg/terraform-provider-proxmox/commit/77f62e1500e7d8468a4dcf60fba3bda907d5c7d9))
* **code:** remove goreleaser from `/tools` ([#1384](https://github.com/bpg/terraform-provider-proxmox/issues/1384)) ([cfe1b1c](https://github.com/bpg/terraform-provider-proxmox/commit/cfe1b1cc4b151a7ab5fee7a051429ecf57b497d6))
* **deps:** update golang.org/x/exp digest (fc45aab → 7f521ea) ([#1391](https://github.com/bpg/terraform-provider-proxmox/issues/1391)) ([f22858c](https://github.com/bpg/terraform-provider-proxmox/commit/f22858ce4e311395f3127da4dd22d45c533caa2e))
* **deps:** Update module github.com/goreleaser/goreleaser (v1.26.2 → v2.0.0) in /tools ([#1381](https://github.com/bpg/terraform-provider-proxmox/issues/1381)) ([095be69](https://github.com/bpg/terraform-provider-proxmox/commit/095be69d8d6ad281db0670d4a22d88c860eff0ff))

## [0.59.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.59.0...v0.59.1) (2024-06-10)


### Bug Fixes

* **vm:** ignore 599 task status response when waiting for VM start ([#1365](https://github.com/bpg/terraform-provider-proxmox/issues/1365)) ([5d4193b](https://github.com/bpg/terraform-provider-proxmox/commit/5d4193b6beb2b4b71b9f5be7507d4ddc6ea3b6b7))


### Miscellaneous

* **api:** refactor `nodes/vms/vms_types.go`: split into multiple files ([#1368](https://github.com/bpg/terraform-provider-proxmox/issues/1368)) ([580381f](https://github.com/bpg/terraform-provider-proxmox/commit/580381f89248f40ad645d8f9208e2223c66d6e4f))
* **code:** fix code inspection warnings ([#1375](https://github.com/bpg/terraform-provider-proxmox/issues/1375)) ([e647b21](https://github.com/bpg/terraform-provider-proxmox/commit/e647b21e3a90d49f0f1300fcaaa9e23312f67b2b))
* **deps:** update module github.com/golangci/golangci-lint (v1.59.0 → v1.59.1) in /tools ([#1373](https://github.com/bpg/terraform-provider-proxmox/issues/1373)) ([24f9990](https://github.com/bpg/terraform-provider-proxmox/commit/24f9990249037893312636abdae8f30391cff207))
* **docs:** clarify `disk.file_id` format ([#1374](https://github.com/bpg/terraform-provider-proxmox/issues/1374)) ([f3ac15a](https://github.com/bpg/terraform-provider-proxmox/commit/f3ac15ad12b0638ed421e25b561f1d16f7d3a41d))
* **vm2:** add initial support for `cdrom` ([#1370](https://github.com/bpg/terraform-provider-proxmox/issues/1370)) ([80cafa6](https://github.com/bpg/terraform-provider-proxmox/commit/80cafa689f3795a808dc0c2ecbcf81608181e370))

## [0.59.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.58.1...v0.59.0) (2024-06-05)


### Features

* **vm:** allow `none` `cdrom` ([#1347](https://github.com/bpg/terraform-provider-proxmox/issues/1347)) ([994f367](https://github.com/bpg/terraform-provider-proxmox/commit/994f367cc2bb444a61ecc5b04e1280394aa6b345))


### Bug Fixes

* **access:** fix `user` and `role` import ([#1346](https://github.com/bpg/terraform-provider-proxmox/issues/1346)) ([bdd9225](https://github.com/bpg/terraform-provider-proxmox/commit/bdd9225041d1dbbc834f7df6adc7463bb3c09b9f))
* **lxc,vm:** fix panics when a resource block is empty ([#1345](https://github.com/bpg/terraform-provider-proxmox/issues/1345)) ([9bbbf8d](https://github.com/bpg/terraform-provider-proxmox/commit/9bbbf8db390ed43291126eca2d3fc391b7308acc))


### Miscellaneous

* **ci:** update actions/create-github-app-token action (v1.10.0 → v1.10.1) ([#1356](https://github.com/bpg/terraform-provider-proxmox/issues/1356)) ([92476d4](https://github.com/bpg/terraform-provider-proxmox/commit/92476d4f4024c9746c3da987a0d7e210525b14ab))
* **ci:** Update goreleaser/goreleaser-action action (v5.1.0 → v6.0.0) ([#1362](https://github.com/bpg/terraform-provider-proxmox/issues/1362)) ([9dbee1a](https://github.com/bpg/terraform-provider-proxmox/commit/9dbee1a12c2435db00fbf2011d45fbeffdb2b15c))
* **deps:** update go (1.22.3 → 1.22.4) ([#1357](https://github.com/bpg/terraform-provider-proxmox/issues/1357)) ([ea509b2](https://github.com/bpg/terraform-provider-proxmox/commit/ea509b289ec002a085d500521519c93873eb0f65))
* **deps:** update golang.org/x/exp digest (23cca88 → fd00a4e) ([#1343](https://github.com/bpg/terraform-provider-proxmox/issues/1343)) ([eaeff08](https://github.com/bpg/terraform-provider-proxmox/commit/eaeff08571e390cd29f037f00b0c08086e704f70))
* **deps:** update golang.org/x/exp digest (fd00a4e → fc45aab) ([#1355](https://github.com/bpg/terraform-provider-proxmox/issues/1355)) ([8fb6410](https://github.com/bpg/terraform-provider-proxmox/commit/8fb64105075e370945c007411b0e74f703e51f4d))
* **deps:** Update module github.com/goreleaser/goreleaser (v1.26.2 → v2.0.0) in /tools ([#1363](https://github.com/bpg/terraform-provider-proxmox/issues/1363)) ([7064f12](https://github.com/bpg/terraform-provider-proxmox/commit/7064f1259a7830e493e7a3e07a2dce79eef030fc))
* **deps:** update module github.com/hashicorp/terraform-plugin-framework (v1.8.0 → v1.9.0) ([#1359](https://github.com/bpg/terraform-provider-proxmox/issues/1359)) ([27ad93e](https://github.com/bpg/terraform-provider-proxmox/commit/27ad93e793864310d0944591802876efcbfafbc2))
* **deps:** update tools ([#1358](https://github.com/bpg/terraform-provider-proxmox/issues/1358)) ([9c4b9db](https://github.com/bpg/terraform-provider-proxmox/commit/9c4b9dbcbad247d2926d9ee8b6407efc0ab91cb8))

## [0.58.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.58.0...v0.58.1) (2024-05-30)


### Bug Fixes

* **vm:** adding disks causes VM to be re-created ([#1336](https://github.com/bpg/terraform-provider-proxmox/issues/1336)) ([d02dc1e](https://github.com/bpg/terraform-provider-proxmox/commit/d02dc1eb0a524d89a06b4b96b1d36c504ab737fe))
* **vm:** host parameter is optional if using mapping for usb ([#1338](https://github.com/bpg/terraform-provider-proxmox/issues/1338)) ([76d9806](https://github.com/bpg/terraform-provider-proxmox/commit/76d980683c5349ac2376d4f8523df241c5c05e0b))
* **vm:** regression: `vga.clipboard` attribute check ([#1340](https://github.com/bpg/terraform-provider-proxmox/issues/1340)) ([2785c40](https://github.com/bpg/terraform-provider-proxmox/commit/2785c40d44ca950742cbd78e6b720b25bb6e7826))


### Miscellaneous

* **deps:** update golang.org/x/exp digest (4c93da0 → 23cca88) ([#1333](https://github.com/bpg/terraform-provider-proxmox/issues/1333)) ([804b498](https://github.com/bpg/terraform-provider-proxmox/commit/804b498453621d308a3d3a4caae168076d494a08))
* **deps:** update module github.com/hashicorp/terraform-plugin-docs (v0.19.2 → v0.19.3) in /tools ([#1334](https://github.com/bpg/terraform-provider-proxmox/issues/1334)) ([2f87941](https://github.com/bpg/terraform-provider-proxmox/commit/2f87941f2b370bdda681d4f51dde78ca0e40154a))

## [0.58.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.57.1...v0.58.0) (2024-05-28)


### Features

* **vm:** add support for `vga.clipboard`, `virtio-gl` type ([#1326](https://github.com/bpg/terraform-provider-proxmox/issues/1326)) ([32de050](https://github.com/bpg/terraform-provider-proxmox/commit/32de05060527b68198a4746ec1402f531fe9c77f))


### Bug Fixes

* **lxc:** `mount_point.backup` should be `false` by default ([#1327](https://github.com/bpg/terraform-provider-proxmox/issues/1327)) ([206fc4b](https://github.com/bpg/terraform-provider-proxmox/commit/206fc4b03ddfcc7ecf43b01ee22df2e66a4505b5))
* **vm:** check task warning status against exit code ([#1330](https://github.com/bpg/terraform-provider-proxmox/issues/1330)) ([b7afd02](https://github.com/bpg/terraform-provider-proxmox/commit/b7afd02cc3bccfc43ab91f2484bf0bca4f532cbc))


### Miscellaneous

* **deps:** update golang.org/x/exp digest (9bf2ced → 4c93da0) ([#1321](https://github.com/bpg/terraform-provider-proxmox/issues/1321)) ([1c6b9e4](https://github.com/bpg/terraform-provider-proxmox/commit/1c6b9e4935d448a00552323c0e8679656bf6978d))
* **deps:** update module github.com/golangci/golangci-lint (v1.58.2 → v1.59.0) in /tools ([#1332](https://github.com/bpg/terraform-provider-proxmox/issues/1332)) ([e0097d9](https://github.com/bpg/terraform-provider-proxmox/commit/e0097d9880698203b9b293c4b9c76e46a697f662))
* **deps:** update module github.com/goreleaser/goreleaser (v1.26.1 → v1.26.2) in /tools ([#1322](https://github.com/bpg/terraform-provider-proxmox/issues/1322)) ([93f86b1](https://github.com/bpg/terraform-provider-proxmox/commit/93f86b1e13ca092bb77ad0972927e52463358b5d))
* **docs:** document expectations for the default `root` shell on PVE nodes ([#1329](https://github.com/bpg/terraform-provider-proxmox/issues/1329)) ([ffffa56](https://github.com/bpg/terraform-provider-proxmox/commit/ffffa56518aed38ae93706ebdb22f5386baffcb1))
* **vm2:** add support for `vga` ([#1328](https://github.com/bpg/terraform-provider-proxmox/issues/1328)) ([d843e46](https://github.com/bpg/terraform-provider-proxmox/commit/d843e46b375e27a9cb0364537adad272beb3bcdb))

## [0.57.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.57.0...v0.57.1) (2024-05-22)


### Bug Fixes

* **file:** upload timeout with 596 response code ([#1315](https://github.com/bpg/terraform-provider-proxmox/issues/1315)) ([a6c6b98](https://github.com/bpg/terraform-provider-proxmox/commit/a6c6b98d449e2039edb76a9bdbb0a71ed7fdeb16))
* **vm:** allow `clone` to pass with warnings ([#1317](https://github.com/bpg/terraform-provider-proxmox/issues/1317)) ([fb1105d](https://github.com/bpg/terraform-provider-proxmox/commit/fb1105d93bd2fec4d9fe0cfa9155539d17d8b418))


### Miscellaneous

* **ci:** update actions/checkout action (v4.1.5 → v4.1.6) ([#1305](https://github.com/bpg/terraform-provider-proxmox/issues/1305)) ([0494fcb](https://github.com/bpg/terraform-provider-proxmox/commit/0494fcb16772cd69d510d42594badfc83ae2e0ac))
* **ci:** update jetbrains/qodana-action action (v2024.1.4 → v2024.1.5) ([#1306](https://github.com/bpg/terraform-provider-proxmox/issues/1306)) ([48fdb03](https://github.com/bpg/terraform-provider-proxmox/commit/48fdb03c8aac90f15040e7a7cad317acb535c797))
* **deps:** update github.com/hashicorp/terraform-plugin-* ([#1309](https://github.com/bpg/terraform-provider-proxmox/issues/1309)) ([7d1b400](https://github.com/bpg/terraform-provider-proxmox/commit/7d1b4006412962aeef40bd64d38021c807543e0e))
* **deps:** update tools ([#1316](https://github.com/bpg/terraform-provider-proxmox/issues/1316)) ([3f97a7b](https://github.com/bpg/terraform-provider-proxmox/commit/3f97a7b5c3eb8926375e02f87f84d23a9b644d2c))
* **vm2:** add datasource implementation  ([#1318](https://github.com/bpg/terraform-provider-proxmox/issues/1318)) ([e3dd31f](https://github.com/bpg/terraform-provider-proxmox/commit/e3dd31f55e41047e004ebff34a713a4e9ed8616f))
* **vm2:** add support for `cpu` ([#1310](https://github.com/bpg/terraform-provider-proxmox/issues/1310)) ([aa309fd](https://github.com/bpg/terraform-provider-proxmox/commit/aa309fd9eabad00e7694b8bd7d63cad6180d7118))
* **vm2:** cleanup `cpu` implementation, refactor rearrange acc tests ([#1311](https://github.com/bpg/terraform-provider-proxmox/issues/1311)) ([2e34c57](https://github.com/bpg/terraform-provider-proxmox/commit/2e34c57f6c95bff2f3e35ee617e783a78616c900))

## [0.57.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.56.1...v0.57.0) (2024-05-16)


### Features

* **vm:** remove support for `initialization.upgrade` attribute ([#1295](https://github.com/bpg/terraform-provider-proxmox/issues/1295)) ([2b7dd08](https://github.com/bpg/terraform-provider-proxmox/commit/2b7dd08020c542b5e08d241aa577be1627f19301))


### Bug Fixes

* **access:** user token `expiration_date` handling ([#1293](https://github.com/bpg/terraform-provider-proxmox/issues/1293)) ([a436c24](https://github.com/bpg/terraform-provider-proxmox/commit/a436c24923d314575bebcebb4e2173dffc753651))
* **vm:** yet another fix for disk reordering ([#1297](https://github.com/bpg/terraform-provider-proxmox/issues/1297)) ([3180f81](https://github.com/bpg/terraform-provider-proxmox/commit/3180f81b4a14868bb0076ec8ab6299b3f6945510))


### Miscellaneous

* **ci:** update google-github-actions/release-please-action action (v4.1.0 → v4.1.1) ([#1296](https://github.com/bpg/terraform-provider-proxmox/issues/1296)) ([71480bf](https://github.com/bpg/terraform-provider-proxmox/commit/71480bf9340d11d3a06c72a4eb5425f920ea03e3))
* **ci:** update release-please action ref ([bfb20e0](https://github.com/bpg/terraform-provider-proxmox/commit/bfb20e0bfdd7a326844351f6ece8d47c1a343085))
* **deps:** update module github.com/goreleaser/goreleaser (v1.26.0 → v1.26.1) in /tools [security] ([#1301](https://github.com/bpg/terraform-provider-proxmox/issues/1301)) ([ea6bb93](https://github.com/bpg/terraform-provider-proxmox/commit/ea6bb93be5e8ae0ba645995165548a780424f008))

## [0.56.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.56.0...v0.56.1) (2024-05-13)


### Bug Fixes

* **vm:** `cdrom` is not attached when creating a VM from scratch ([#1287](https://github.com/bpg/terraform-provider-proxmox/issues/1287)) ([be4a3b7](https://github.com/bpg/terraform-provider-proxmox/commit/be4a3b7c6f73a441f5e562f2b97c1fafef06bad0))


### Miscellaneous

* **deps:** update module github.com/goreleaser/goreleaser (v1.25.1 → v1.26.0) in /tools ([#1289](https://github.com/bpg/terraform-provider-proxmox/issues/1289)) ([7853a53](https://github.com/bpg/terraform-provider-proxmox/commit/7853a535270aa2a0782de8fad20a9d0c57f0b322))

## [0.56.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.55.1...v0.56.0) (2024-05-12)


### Features

* **access:** add ACL resource ([#1166](https://github.com/bpg/terraform-provider-proxmox/issues/1166)) ([afcbb41](https://github.com/bpg/terraform-provider-proxmox/commit/afcbb415a30fdf67912fc0c8ecace338c88fd003))
* **access:** add proxmox user token ([#1159](https://github.com/bpg/terraform-provider-proxmox/issues/1159)) ([8220271](https://github.com/bpg/terraform-provider-proxmox/commit/8220271eee5755e1d0d87f5064bfb00304a307bb))
* **provider:** add support for `ssh-agent` on Windows ([#1270](https://github.com/bpg/terraform-provider-proxmox/issues/1270)) ([ccf4834](https://github.com/bpg/terraform-provider-proxmox/commit/ccf4834c16797e099782a13254ab1ef011833133))


### Miscellaneous

* **ci:** update actions/checkout action (v4.1.4 → v4.1.5) ([#1278](https://github.com/bpg/terraform-provider-proxmox/issues/1278)) ([10774f4](https://github.com/bpg/terraform-provider-proxmox/commit/10774f49e17ec4585508989ed225b9d10c2710fb))
* **ci:** update goreleaser/goreleaser-action action (v5.0.0 → v5.1.0) ([#1284](https://github.com/bpg/terraform-provider-proxmox/issues/1284)) ([78f9d6e](https://github.com/bpg/terraform-provider-proxmox/commit/78f9d6e2bfecebd47018c16ab1ec22e6f00eb878))
* **ci:** update hashicorp/setup-terraform action (v3.1.0 → v3.1.1) ([#1279](https://github.com/bpg/terraform-provider-proxmox/issues/1279)) ([e79cf12](https://github.com/bpg/terraform-provider-proxmox/commit/e79cf12f44602d9ca9c6e90b1554ac602d79f2b7))
* **ci:** update jetbrains/qodana-action action (v2024.1.3 → v2024.1.4) ([#1280](https://github.com/bpg/terraform-provider-proxmox/issues/1280)) ([6ae9b58](https://github.com/bpg/terraform-provider-proxmox/commit/6ae9b581df6e82cfe9bae2bb9174d8cc3b46ea6c))
* **deps:** update go (1.22.2 → 1.22.3) ([#1273](https://github.com/bpg/terraform-provider-proxmox/issues/1273)) ([097f3f9](https://github.com/bpg/terraform-provider-proxmox/commit/097f3f947ca97708a6acf67e036205d459b9c0bd))
* **deps:** update module github.com/brianvoe/gofakeit/v7 (v7.0.2 → v7.0.3) ([#1274](https://github.com/bpg/terraform-provider-proxmox/issues/1274)) ([2be83eb](https://github.com/bpg/terraform-provider-proxmox/commit/2be83eba1ea4512579ea0188c77709fd5a409c29))
* **deps:** update module github.com/gorilla/schema (v1.2.1 → v1.3.0) ([#1276](https://github.com/bpg/terraform-provider-proxmox/issues/1276)) ([1f80e9d](https://github.com/bpg/terraform-provider-proxmox/commit/1f80e9dc3132f16eca7533c918f2928e0eaf33e5))
* **deps:** update module github.com/hashicorp/terraform-plugin-mux (v0.15.0 → v0.16.0) ([#1277](https://github.com/bpg/terraform-provider-proxmox/issues/1277)) ([ff98439](https://github.com/bpg/terraform-provider-proxmox/commit/ff9843945943ab0287535f84c3fe1bd60ec2b4d5))
* **deps:** update tools ([#1275](https://github.com/bpg/terraform-provider-proxmox/issues/1275)) ([d3ea1ac](https://github.com/bpg/terraform-provider-proxmox/commit/d3ea1ac560620552ac4d4f67a9a0a0eb4de43012))

## [0.55.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.55.0...v0.55.1) (2024-05-06)


### Bug Fixes

* **code:** update acc tests to PVE 8.2 ([#1269](https://github.com/bpg/terraform-provider-proxmox/issues/1269)) ([0b35f47](https://github.com/bpg/terraform-provider-proxmox/commit/0b35f47780a15d5c89b3d229e90c3af4ac35f7ce))
* **vm:** fix timeout when resizing custom disk at create ([#1260](https://github.com/bpg/terraform-provider-proxmox/issues/1260)) ([10790f6](https://github.com/bpg/terraform-provider-proxmox/commit/10790f668d4b57b8bfe0f62cb691bdb75e12fca3))
* **vm:** only set `initialization.upgrade` attribute when not using custom cloud-init file ([#1253](https://github.com/bpg/terraform-provider-proxmox/issues/1253)) ([9762405](https://github.com/bpg/terraform-provider-proxmox/commit/97624052f152d43cf325908ed717ef444306cbd7))
* **vm:** use `raw` file format for `efi_disk` by default ([#1249](https://github.com/bpg/terraform-provider-proxmox/issues/1249)) ([b3caf6a](https://github.com/bpg/terraform-provider-proxmox/commit/b3caf6af8252e1fc737489f71534b96188f50631))


### Miscellaneous

* **ci:** update actions/create-github-app-token action (v1.9.3 → v1.10.0) ([#1262](https://github.com/bpg/terraform-provider-proxmox/issues/1262)) ([ce47411](https://github.com/bpg/terraform-provider-proxmox/commit/ce47411f99115a85d900a2671b2e57e21d680189))
* **ci:** update jetbrains/qodana-action action (v2024.1.2 → v2024.1.3) ([#1261](https://github.com/bpg/terraform-provider-proxmox/issues/1261)) ([f08d139](https://github.com/bpg/terraform-provider-proxmox/commit/f08d139ba8b8e5bf3fececeda852a2d02ffd80cd))
* **deps:** update module github.com/golangci/golangci-lint (v1.57.2 → v1.58.0) in /tools ([#1258](https://github.com/bpg/terraform-provider-proxmox/issues/1258)) ([452b24b](https://github.com/bpg/terraform-provider-proxmox/commit/452b24b36b15eb98cddba313227a37a770553f64))
* **deps:** update module github.com/hashicorp/terraform-plugin-go (v0.22.2 → v0.23.0) ([#1266](https://github.com/bpg/terraform-provider-proxmox/issues/1266)) ([02302cb](https://github.com/bpg/terraform-provider-proxmox/commit/02302cb498999c09fe71479f2f35676a472c98fa))
* **deps:** update module golang.org/x/crypto (v0.22.0 → v0.23.0) ([#1267](https://github.com/bpg/terraform-provider-proxmox/issues/1267)) ([c21da7e](https://github.com/bpg/terraform-provider-proxmox/commit/c21da7e1831a929703bac0ac98b3b4f3db55ea4d))
* **deps:** update module golang.org/x/exp (v0.0.0-20240416160154-fe59bbe5cc7f → v0.0.0-20240506185415-9bf2ced13842) ([#1265](https://github.com/bpg/terraform-provider-proxmox/issues/1265)) ([6bacbe8](https://github.com/bpg/terraform-provider-proxmox/commit/6bacbe8b56562c5c814c65bbb44f03602fae41be))
* **deps:** update module golang.org/x/net (v0.24.0 → v0.25.0) ([#1268](https://github.com/bpg/terraform-provider-proxmox/issues/1268)) ([2672bf3](https://github.com/bpg/terraform-provider-proxmox/commit/2672bf3ca12a134bc6daa5943a01e66a87ed77d6))

## [0.55.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.54.0...v0.55.0) (2024-04-30)


### Features

* **cluster:** Extend the available attributes for the `proxmox_virtual_environment_cluster_options` resource ([#1241](https://github.com/bpg/terraform-provider-proxmox/issues/1241)) ([2eb36f4](https://github.com/bpg/terraform-provider-proxmox/commit/2eb36f4134a7f89e2e3fb1f69f39d359730c2acc))
* **vm:** add initial support for IDE ([#1237](https://github.com/bpg/terraform-provider-proxmox/issues/1237)) ([bd195d6](https://github.com/bpg/terraform-provider-proxmox/commit/bd195d66069485fc290e13fc66d2b80b220ea731))


### Miscellaneous

* **ci:** update actions/checkout action (v4.1.1 → v4.1.3) ([bb44a8d](https://github.com/bpg/terraform-provider-proxmox/commit/bb44a8d84491cdca770bbbf58df4d8b0c5605f1e))
* **ci:** update actions/checkout action (v4.1.3 → v4.1.4) ([229883b](https://github.com/bpg/terraform-provider-proxmox/commit/229883b110f6e6ab22b9e79a52167d044971c4c8))
* **ci:** update hashicorp/setup-terraform action (v3.0.0 → v3.1.0) ([84778a3](https://github.com/bpg/terraform-provider-proxmox/commit/84778a39656cc8c2bb36f767d1ce55916353073c))
* **ci:** update jetbrains/qodana-action action (v2024.1.1 → v2024.1.2) ([4da06ec](https://github.com/bpg/terraform-provider-proxmox/commit/4da06ecb94b1d3097e2a0fb12c7aeafc907d60a1))
* **ci:** update lycheeverse/lychee-action action (v1.9.3 → v1.10.0) ([0839c92](https://github.com/bpg/terraform-provider-proxmox/commit/0839c92ba69eb3ba53aa098a0b3a970ba405a533))
* **deps:** update go (1.22.0 → 1.22.2) ([#1234](https://github.com/bpg/terraform-provider-proxmox/issues/1234)) ([ee939a3](https://github.com/bpg/terraform-provider-proxmox/commit/ee939a38a35bf4a490cbb549c14a8a5538e47e23))
* **deps:** update module github.com/hashicorp/terraform-plugin-docs (v0.19.1 → v0.19.2) in /tools ([#1247](https://github.com/bpg/terraform-provider-proxmox/issues/1247)) ([b9429cb](https://github.com/bpg/terraform-provider-proxmox/commit/b9429cbd682ab1fff65f05fd67bc9f34bd6f5906))
* **deps:** update tools ([#1238](https://github.com/bpg/terraform-provider-proxmox/issues/1238)) ([e13bb03](https://github.com/bpg/terraform-provider-proxmox/commit/e13bb03ddafac41e7e9ca815ab3e615b839b56e3))
* remove unused code ([#1236](https://github.com/bpg/terraform-provider-proxmox/issues/1236)) ([efd0d7b](https://github.com/bpg/terraform-provider-proxmox/commit/efd0d7b4662b0bb0bebf7827f9b3b508a2f298b6))
* **vm2:** experiment with sets, add `tags` support ([#1232](https://github.com/bpg/terraform-provider-proxmox/issues/1232)) ([961404d](https://github.com/bpg/terraform-provider-proxmox/commit/961404dd0b2be71c1b612201cc454226ac12c04e))
* **vm2:** experimental support for `clone` and inherited attributes ([#1235](https://github.com/bpg/terraform-provider-proxmox/issues/1235)) ([7209fe0](https://github.com/bpg/terraform-provider-proxmox/commit/7209fe03215c7fda32dd74ab9647b9824e0b8d61))

## [0.54.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.53.1...v0.54.0) (2024-04-19)


### Features

* **cluster:** Implement initial support for "hardware mappings" ([#1213](https://github.com/bpg/terraform-provider-proxmox/issues/1213)) ([76c1c22](https://github.com/bpg/terraform-provider-proxmox/commit/76c1c2282bb24ccc01cd5dd138ad19e830633d10))


### Bug Fixes

* **vm,lxc,file:** improve timeouts handling ([#1222](https://github.com/bpg/terraform-provider-proxmox/issues/1222)) ([0c9c206](https://github.com/bpg/terraform-provider-proxmox/commit/0c9c2066fd88abf2ec422423113de29651be0417))
* **vm:** state drift due to disk re-ordering ([#1215](https://github.com/bpg/terraform-provider-proxmox/issues/1215)) ([ad036a6](https://github.com/bpg/terraform-provider-proxmox/commit/ad036a67e8129ce6a50096792c33cb1ce364967e))


### Miscellaneous

* **code:** reorganize HA & node network code ([#1218](https://github.com/bpg/terraform-provider-proxmox/issues/1218)) ([163a773](https://github.com/bpg/terraform-provider-proxmox/commit/163a7730888dac2d436988499af87d4ff8f7d118))
* **deps:** bump golang.org/x/net from 0.22.0 to 0.23.0 in /tools ([#1225](https://github.com/bpg/terraform-provider-proxmox/issues/1225)) ([7e4710d](https://github.com/bpg/terraform-provider-proxmox/commit/7e4710d3f01ab076d2ed702885f660fd1d14b90f))
* **deps:** update github.com/hashicorp/terraform-plugin-* ([#1228](https://github.com/bpg/terraform-provider-proxmox/issues/1228)) ([5a606ec](https://github.com/bpg/terraform-provider-proxmox/commit/5a606ec18e27a89289880cff06d767758f716ce6))
* **deps:** update go (1.22.0 → 1.22.2) in /tools ([#1227](https://github.com/bpg/terraform-provider-proxmox/issues/1227)) ([c7e36ee](https://github.com/bpg/terraform-provider-proxmox/commit/c7e36ee1a18375e15fbc365e7014f9fa9d7be24b))
* **deps:** update golang.org/x/exp digest (93d18d7 → fe59bbe) ([#1226](https://github.com/bpg/terraform-provider-proxmox/issues/1226)) ([83d835f](https://github.com/bpg/terraform-provider-proxmox/commit/83d835ff3ebaf05e960ab500192feb12e4ef6cc6))
* **deps:** update module github.com/avast/retry-go/v4 (v4.5.1 → v4.6.0) ([#1229](https://github.com/bpg/terraform-provider-proxmox/issues/1229)) ([9387845](https://github.com/bpg/terraform-provider-proxmox/commit/938784531f331afddc1e0639a11b75e96e96af17))
* **deps:** update tools ([#1216](https://github.com/bpg/terraform-provider-proxmox/issues/1216)) ([1bf3bf0](https://github.com/bpg/terraform-provider-proxmox/commit/1bf3bf0f68bf73efd82e153d3efdc1ed4db75dd1))
* **vm2:** initial experimental VM resource implementation using Plugin Framework ([#1230](https://github.com/bpg/terraform-provider-proxmox/issues/1230)) ([d8202dd](https://github.com/bpg/terraform-provider-proxmox/commit/d8202dd7a1c4f41ba9503f83b97633bd1cba30d4))
* **vm:** fix acceptance test for hardware mapping ([#1223](https://github.com/bpg/terraform-provider-proxmox/issues/1223)) ([653e6f2](https://github.com/bpg/terraform-provider-proxmox/commit/653e6f25f848ccd13549b6e5d6ca7ba2e713ecc4))

## [0.53.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.53.0...v0.53.1) (2024-04-12)


### Bug Fixes

* **vm:** error 'resizing of efidisks is not supported' when clonning a VM with re-defined `efi_disk` ([#1206](https://github.com/bpg/terraform-provider-proxmox/issues/1206)) ([46d323d](https://github.com/bpg/terraform-provider-proxmox/commit/46d323d31b1568af363bbee84f7ea7e6a46a62cc))
* **vm:** regression: force VM recreation due to `initialization.upgrade` when it is not in the plan ([#1209](https://github.com/bpg/terraform-provider-proxmox/issues/1209)) ([628d0f1](https://github.com/bpg/terraform-provider-proxmox/commit/628d0f15dd7499a15b7ae1a00441a72191b130be))


### Miscellaneous

* **deps:** bump github.com/sigstore/cosign/v2 from 2.2.3 to 2.2.4 in /tools ([#1207](https://github.com/bpg/terraform-provider-proxmox/issues/1207)) ([19a13be](https://github.com/bpg/terraform-provider-proxmox/commit/19a13bee0a023c850854be02310a424762b7ed8c))

## [0.53.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.52.0...v0.53.0) (2024-04-10)


### Features

* **vm:** add Intel E1000E support ([#1191](https://github.com/bpg/terraform-provider-proxmox/issues/1191)) ([f0a7afa](https://github.com/bpg/terraform-provider-proxmox/commit/f0a7afa57d7dc538320efe436890e23cc89993e3))
* **vm:** add support for `initialization.upgrade` attribute ([#1203](https://github.com/bpg/terraform-provider-proxmox/issues/1203)) ([59972dc](https://github.com/bpg/terraform-provider-proxmox/commit/59972dc1b5c4171224c4e3ad3d9790a6efc66759))


### Bug Fixes

* **firewall:** remove strict validation for `rule.iface` attribute ([#1199](https://github.com/bpg/terraform-provider-proxmox/issues/1199)) ([ae88cf6](https://github.com/bpg/terraform-provider-proxmox/commit/ae88cf603d1d7e3ca86ba1506944be0f23f9301c))
* **vm:** panic if `initialization.user_account_keys` contains `null` ([#1197](https://github.com/bpg/terraform-provider-proxmox/issues/1197)) ([b4174a7](https://github.com/bpg/terraform-provider-proxmox/commit/b4174a7c82398970f9f318955fdd07d7ab378360))
* **vm:** panic if `numa` block is empty ([#1196](https://github.com/bpg/terraform-provider-proxmox/issues/1196)) ([c273111](https://github.com/bpg/terraform-provider-proxmox/commit/c27311183ab74bec3e12228917e1b5468a4b941c))


### Miscellaneous

* **deps:** update golang.org/x/exp digest (c0f41cb → 93d18d7) ([#1201](https://github.com/bpg/terraform-provider-proxmox/issues/1201)) ([ead36dc](https://github.com/bpg/terraform-provider-proxmox/commit/ead36dcafbe1a475748dd846b0a05c94b831c600))
* **docs:** minor documentation updates ([#1200](https://github.com/bpg/terraform-provider-proxmox/issues/1200)) ([c14b7eb](https://github.com/bpg/terraform-provider-proxmox/commit/c14b7ebe1d524783badf299cb4e146168011729e))
* refactor acceptance tests ([#1195](https://github.com/bpg/terraform-provider-proxmox/issues/1195)) ([c772fb3](https://github.com/bpg/terraform-provider-proxmox/commit/c772fb3cf658fcda9051f8305125fdfc78ce3b3b))

## [0.52.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.51.1...v0.52.0) (2024-04-07)


### Features

* **vm:** add support for numa architecture attribute ([#1156](https://github.com/bpg/terraform-provider-proxmox/issues/1156)) ([#1175](https://github.com/bpg/terraform-provider-proxmox/issues/1175)) ([dbbd966](https://github.com/bpg/terraform-provider-proxmox/commit/dbbd966736716ab730da04ef1099cc82d2e85019))


### Bug Fixes

* **vm:** regression: provider always tries to update `memory.*hugepages` even if it is not specified ([#1188](https://github.com/bpg/terraform-provider-proxmox/issues/1188)) ([31b6812](https://github.com/bpg/terraform-provider-proxmox/commit/31b6812ce281da214261786f633ed7e93c197fb0))


### Miscellaneous

* **ci:** update actions/create-github-app-token action (v1.9.2 → v1.9.3) ([#1190](https://github.com/bpg/terraform-provider-proxmox/issues/1190)) ([24daca5](https://github.com/bpg/terraform-provider-proxmox/commit/24daca59410bb3974045e7dd5825fdbe0ef1356e))
* **deps:** update tools ([#1184](https://github.com/bpg/terraform-provider-proxmox/issues/1184)) ([f4e5ef8](https://github.com/bpg/terraform-provider-proxmox/commit/f4e5ef8aaacdee4063eb65cd5fb53d197229335a))
* improve acceptance tests on CI ([#1173](https://github.com/bpg/terraform-provider-proxmox/issues/1173)) ([268c861](https://github.com/bpg/terraform-provider-proxmox/commit/268c8616494219463fc2f20b2279a166f17c791a))

## [0.51.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.51.0...v0.51.1) (2024-04-04)


### Bug Fixes

* **vm:** regressions: provider always tries to update `cpu.affinity` even if it is not specified ([#1182](https://github.com/bpg/terraform-provider-proxmox/issues/1182)) ([82d435f](https://github.com/bpg/terraform-provider-proxmox/commit/82d435f575576e495e0558d4f32b1620bbc4f842))


### Miscellaneous

* **ci:** update actions/create-github-app-token action (v1.9.0 → v1.9.2) ([#1177](https://github.com/bpg/terraform-provider-proxmox/issues/1177)) ([7358257](https://github.com/bpg/terraform-provider-proxmox/commit/7358257dbaa0f5d4c9de2aaaa24c7142d390baf6))
* **deps:** update golang.org/x/exp digest (a685a6e → c0f41cb) ([#1183](https://github.com/bpg/terraform-provider-proxmox/issues/1183)) ([e680aa9](https://github.com/bpg/terraform-provider-proxmox/commit/e680aa9318b638d0ea7d235893eccb6ef102c925))

## [0.51.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.50.0...v0.51.0) (2024-03-30)


### Features

* **vm:** add proxmox_virtual_environment_node datasource ([#1151](https://github.com/bpg/terraform-provider-proxmox/issues/1151)) ([9fee8c1](https://github.com/bpg/terraform-provider-proxmox/commit/9fee8c1e2e24a61cd8e75c3744227f5d44662c0f))
* **vm:** add support for cpu `affinity` attribute ([#1148](https://github.com/bpg/terraform-provider-proxmox/issues/1148)) ([#1153](https://github.com/bpg/terraform-provider-proxmox/issues/1153)) ([83a5971](https://github.com/bpg/terraform-provider-proxmox/commit/83a5971cf64451f41d9a24b6245f59cdfae14569))
* **vm:** add support for memory `hugepages` and `keep_hugepages` attributes ([#1158](https://github.com/bpg/terraform-provider-proxmox/issues/1158)) ([a2f055a](https://github.com/bpg/terraform-provider-proxmox/commit/a2f055a916104b1ba8008b7b359080acacb605f6))


### Bug Fixes

* **lxc:** improve configurable timeouts for containers operations ([#1161](https://github.com/bpg/terraform-provider-proxmox/issues/1161)) ([c45e336](https://github.com/bpg/terraform-provider-proxmox/commit/c45e3367e7128f82e6931fdf88bd31f9736147fd))
* prevent authenticating multiple times in parallel ([#1168](https://github.com/bpg/terraform-provider-proxmox/issues/1168)) ([e87bc4b](https://github.com/bpg/terraform-provider-proxmox/commit/e87bc4b941564aace95d60ff987b7ec0d508b437))
* **vm:** re-use ide, sata, and scsi cloud-init storage ([#1083](https://github.com/bpg/terraform-provider-proxmox/issues/1083)) ([#1141](https://github.com/bpg/terraform-provider-proxmox/issues/1141)) ([d41463d](https://github.com/bpg/terraform-provider-proxmox/commit/d41463de79dbae512a5560926cc95c2e46aba70f))


### Miscellaneous

* **deps:** update golang.org/x/exp digest (a85f2c6 → a685a6e) ([#1171](https://github.com/bpg/terraform-provider-proxmox/issues/1171)) ([e60ee99](https://github.com/bpg/terraform-provider-proxmox/commit/e60ee990fa691228bdfe64c32d52ee8e381689a0))
* **deps:** update module github.com/golangci/golangci-lint (v1.56.2 → v1.57.1) in /tools ([#1149](https://github.com/bpg/terraform-provider-proxmox/issues/1149)) ([48756b6](https://github.com/bpg/terraform-provider-proxmox/commit/48756b69ba8f1fcafd88b1186a925236cb93cc30))
* **deps:** update module github.com/golangci/golangci-lint (v1.57.1 → v1.57.2) in /tools ([#1172](https://github.com/bpg/terraform-provider-proxmox/issues/1172)) ([02efdc4](https://github.com/bpg/terraform-provider-proxmox/commit/02efdc40bf28d85f08d7b1b185386063db81cd7a))
* **deps:** update module github.com/hashicorp/terraform-plugin-framework (v1.6.1 → v1.7.0) ([#1150](https://github.com/bpg/terraform-provider-proxmox/issues/1150)) ([e5ef91e](https://github.com/bpg/terraform-provider-proxmox/commit/e5ef91e72d896b45fbe3e3d0c7179028db632b1d))

## [0.50.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.49.0...v0.50.0) (2024-03-21)


### Features

* **lxc:** add configurable timeout for container creation ([#1146](https://github.com/bpg/terraform-provider-proxmox/issues/1146)) ([84440d8](https://github.com/bpg/terraform-provider-proxmox/commit/84440d8329c607d7a1fd7ced8a078251cf23908c))
* **lxc:** support hook script for LXC ([#1140](https://github.com/bpg/terraform-provider-proxmox/issues/1140)) ([0deaf18](https://github.com/bpg/terraform-provider-proxmox/commit/0deaf1801ade7896500c43d4903a69ee5da240a3))
* **vm:** add support for `disconnected` attribute in network interface ([#1129](https://github.com/bpg/terraform-provider-proxmox/issues/1129)) ([11ca880](https://github.com/bpg/terraform-provider-proxmox/commit/11ca880f28b68c7beec1feb7e9dc7a6fbf00bf40))


### Bug Fixes

* **docs:** fix `sudo` configuration steps ([#1132](https://github.com/bpg/terraform-provider-proxmox/issues/1132)) ([0c6ecb8](https://github.com/bpg/terraform-provider-proxmox/commit/0c6ecb85ed0cc236db13874f87bf5f1b17c6c933))
* **docs:** typo in `container.console` attribute `mode` -&gt; `type` ([#1133](https://github.com/bpg/terraform-provider-proxmox/issues/1133)) ([5cd6b4a](https://github.com/bpg/terraform-provider-proxmox/commit/5cd6b4acdfd4777fb7eaa48b74328284a88e53ec))
* **examples:** `datastore_id` is missing in `initialization` block for vm ([#1136](https://github.com/bpg/terraform-provider-proxmox/issues/1136)) ([6a5becc](https://github.com/bpg/terraform-provider-proxmox/commit/6a5beccf9133c3c8306be37d6812f90f97508799))


### Miscellaneous

* **deps:** bump github.com/docker/docker from 24.0.7+incompatible to 24.0.9+incompatible in /tools ([#1142](https://github.com/bpg/terraform-provider-proxmox/issues/1142)) ([f4c80ad](https://github.com/bpg/terraform-provider-proxmox/commit/f4c80ad297c11e139d7a0b2bb851b00dd02b7ff7))
* **deps:** update golang.org/x/exp digest (814bf88 → c7f7c64) ([#1130](https://github.com/bpg/terraform-provider-proxmox/issues/1130)) ([ab51262](https://github.com/bpg/terraform-provider-proxmox/commit/ab512622f7f21bbfcbae48793cf9d0b24af5c2d5))
* **deps:** update golang.org/x/exp digest (c7f7c64 → a85f2c6) ([#1139](https://github.com/bpg/terraform-provider-proxmox/issues/1139)) ([6f8a472](https://github.com/bpg/terraform-provider-proxmox/commit/6f8a472981a3d1f478f9a0fd44abb83f2e19d641))
* **vm:** refactoring: extract network device code from vm.go ([#1127](https://github.com/bpg/terraform-provider-proxmox/issues/1127)) ([29b5438](https://github.com/bpg/terraform-provider-proxmox/commit/29b5438faf7c1d09291c549516ac359952e2a87b))

## [0.49.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.48.4...v0.49.0) (2024-03-15)


### Features

* **vm:** add support for `protection` attribute ([#1126](https://github.com/bpg/terraform-provider-proxmox/issues/1126)) ([cea65a8](https://github.com/bpg/terraform-provider-proxmox/commit/cea65a8a129168581451ecda5b84be8ddd2e3a3e))
* **vm:** add support for disk `aio`, `backup`, `iops*` attributes ([#1124](https://github.com/bpg/terraform-provider-proxmox/issues/1124)) ([014b59e](https://github.com/bpg/terraform-provider-proxmox/commit/014b59e04f30fc08bc512f68cf471fe2cfdc481c))


### Miscellaneous

* **deps:** bump google.golang.org/protobuf from 1.31.0 to 1.33.0 in /tools ([#1121](https://github.com/bpg/terraform-provider-proxmox/issues/1121)) ([1ffc40b](https://github.com/bpg/terraform-provider-proxmox/commit/1ffc40b77ae6d76c8eccf64373ef982e2b14f4d1))
* **deps:** update module github.com/skeema/knownhosts (v1.2.1 → v1.2.2) ([#1123](https://github.com/bpg/terraform-provider-proxmox/issues/1123)) ([732f018](https://github.com/bpg/terraform-provider-proxmox/commit/732f0188d798fc082fcc01312b41a9f5ded2ae0c))

## [0.48.4](https://github.com/bpg/terraform-provider-proxmox/compare/v0.48.3...v0.48.4) (2024-03-13)


### Bug Fixes

* **ci:** missing releases from HashiCorp Registry ([#1118](https://github.com/bpg/terraform-provider-proxmox/issues/1118)) ([ffc64d2](https://github.com/bpg/terraform-provider-proxmox/commit/ffc64d209a392afb3198acee3ee5449b7392e579))

## [0.48.3](https://github.com/bpg/terraform-provider-proxmox/compare/v0.48.2...v0.48.3) (2024-03-12)


### Bug Fixes

* **provider:** EOF error when closing SSH session ([#1113](https://github.com/bpg/terraform-provider-proxmox/issues/1113)) ([b63f1b7](https://github.com/bpg/terraform-provider-proxmox/commit/b63f1b7889287558510526f8392cfdaa9d22524b))
* **vm:** timeout when resizing a disk during clone ([#1103](https://github.com/bpg/terraform-provider-proxmox/issues/1103)) ([449f9fc](https://github.com/bpg/terraform-provider-proxmox/commit/449f9fc31c0d737d2094b4c0db7a207b3e764122))


### Miscellaneous

* **ci:** update google-github-actions/release-please-action action (v4.0.2 → v4.1.0) ([#1115](https://github.com/bpg/terraform-provider-proxmox/issues/1115)) ([04e7421](https://github.com/bpg/terraform-provider-proxmox/commit/04e74219e3cac4805c3ae9cedced42f7f64ed461))
* **deps:** update module github.com/hashicorp/terraform-plugin-go (v0.22.0 → v0.22.1) ([#1114](https://github.com/bpg/terraform-provider-proxmox/issues/1114)) ([a059728](https://github.com/bpg/terraform-provider-proxmox/commit/a0597289b56219a07ece3296e587a8317b1251e9))
* **docs:** update terraform local (2.4.1 → 2.5.1) ([#1116](https://github.com/bpg/terraform-provider-proxmox/issues/1116)) ([29d60f5](https://github.com/bpg/terraform-provider-proxmox/commit/29d60f593232f08440f7e2c9426d12c24eacd572))
* minor cleanups and doc updates ([#1108](https://github.com/bpg/terraform-provider-proxmox/issues/1108)) ([27dbcad](https://github.com/bpg/terraform-provider-proxmox/commit/27dbcad5cdd732a4777e886806c5eeb1a06129a4))

## [0.48.2](https://github.com/bpg/terraform-provider-proxmox/compare/v0.48.1...v0.48.2) (2024-03-09)


### Bug Fixes

* **network:** multiple fixes to `network_linux_bridge` resource ([#1095](https://github.com/bpg/terraform-provider-proxmox/issues/1095)) ([4aed7cb](https://github.com/bpg/terraform-provider-proxmox/commit/4aed7cb085c87aed68f3dc426644a9d76c075db1))
* **provider:** allow LDAP realm API tokens ([#1101](https://github.com/bpg/terraform-provider-proxmox/issues/1101)) ([461321c](https://github.com/bpg/terraform-provider-proxmox/commit/461321cf5e15dfd1b89a506ecf6a410e24bc8c5d))


### Miscellaneous

* **deps:** bump gopkg.in/go-jose/go-jose.v2 from 2.6.1 to 2.6.3 in /tools ([#1102](https://github.com/bpg/terraform-provider-proxmox/issues/1102)) ([4bc7b29](https://github.com/bpg/terraform-provider-proxmox/commit/4bc7b291fe9cedd41eb0a1c26107ff77710cf517))
* **deps:** update github.com/hashicorp/terraform-plugin-* ([#1096](https://github.com/bpg/terraform-provider-proxmox/issues/1096)) ([4cac320](https://github.com/bpg/terraform-provider-proxmox/commit/4cac320ff9008b1274a944af4a6b3b302af276e0))
* **deps:** update module golang.org/x/crypto (v0.20.0 → v0.21.0) ([#1097](https://github.com/bpg/terraform-provider-proxmox/issues/1097)) ([3b9739a](https://github.com/bpg/terraform-provider-proxmox/commit/3b9739ab5986acc4e0b25772704ce57d20818384))
* **deps:** update module golang.org/x/net (v0.21.0 → v0.22.0) ([#1105](https://github.com/bpg/terraform-provider-proxmox/issues/1105)) ([7dea6f6](https://github.com/bpg/terraform-provider-proxmox/commit/7dea6f6ee8d3d8b5ff8027eea388c62a559ef4d6))

## [0.48.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.48.0...v0.48.1) (2024-03-05)


### Bug Fixes

* **ci:** TestAccResourceVMNetwork acc test fails when run on CI ([#1092](https://github.com/bpg/terraform-provider-proxmox/issues/1092)) ([61a0fcd](https://github.com/bpg/terraform-provider-proxmox/commit/61a0fcd936c3c88e6eb0b7b5d5517795b4c3c092))
* **docs:** fix wrong startup delay attributes ([#1088](https://github.com/bpg/terraform-provider-proxmox/issues/1088)) ([85705fd](https://github.com/bpg/terraform-provider-proxmox/commit/85705fdd51b5e64662bea169d86922ff85f062cb))


### Miscellaneous

* **deps:** update tools ([#1017](https://github.com/bpg/terraform-provider-proxmox/issues/1017)) ([fbd04ed](https://github.com/bpg/terraform-provider-proxmox/commit/fbd04ed95061f23747e4bb7224901f6a409f7547))
* **docs:** minor improvements around SSH private key usage ([#1091](https://github.com/bpg/terraform-provider-proxmox/issues/1091)) ([171dd2f](https://github.com/bpg/terraform-provider-proxmox/commit/171dd2f234b7e1effe00bbe66bc42c30f78f9e2d))

## [0.48.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.47.0...v0.48.0) (2024-03-03)


### ⚠ BREAKING CHANGES

* **file:** snippets upload using SSH input stream ([#1085](https://github.com/bpg/terraform-provider-proxmox/issues/1085))

### Features

* **file:** snippets upload using SSH input stream ([#1085](https://github.com/bpg/terraform-provider-proxmox/issues/1085)) ([3195b3c](https://github.com/bpg/terraform-provider-proxmox/commit/3195b3cdf4c7c9d0c9e23177b4bd097de3b1fa65))
* **provider:** add support for private key authentication for SSH ([#1076](https://github.com/bpg/terraform-provider-proxmox/issues/1076)) ([2c6d3ad](https://github.com/bpg/terraform-provider-proxmox/commit/2c6d3ad01d7b6882597415d032380cd32cbaa68f))
* **vm:** add `VLAN` trunk support ([#1086](https://github.com/bpg/terraform-provider-proxmox/issues/1086)) ([cb5fc27](https://github.com/bpg/terraform-provider-proxmox/commit/cb5fc279cd44de9b9782aff5749a771975f72f51))


### Miscellaneous

* **ci:** setup acceptance tests ([#1071](https://github.com/bpg/terraform-provider-proxmox/issues/1071)) ([0bf42d5](https://github.com/bpg/terraform-provider-proxmox/commit/0bf42d52e5af26c423730bd5c339bd295abf2533))
* **ci:** split acceptance tests into a separate workflow ([#1084](https://github.com/bpg/terraform-provider-proxmox/issues/1084)) ([e38b45f](https://github.com/bpg/terraform-provider-proxmox/commit/e38b45f033a147f216228df0bf9a527665bbd808))
* **ci:** update actions/create-github-app-token action (v1.8.1 → v1.9.0) ([66ec9f4](https://github.com/bpg/terraform-provider-proxmox/commit/66ec9f4b9b027eb963be6b9d1e8a56c6a4610fc4))
* **ci:** update dorny/paths-filter action (v3.0.1 → v3.0.2) ([3d6cc75](https://github.com/bpg/terraform-provider-proxmox/commit/3d6cc75107c52d8eb42a46e83cd21673770968be))
* **deps:** update github.com/hashicorp/terraform-plugin-* ([#1078](https://github.com/bpg/terraform-provider-proxmox/issues/1078)) ([2398f6c](https://github.com/bpg/terraform-provider-proxmox/commit/2398f6c339c891d78eae501648c673af470793a8))
* **deps:** update module github.com/brianvoe/gofakeit/v7 (v7.0.1 → v7.0.2) ([#1080](https://github.com/bpg/terraform-provider-proxmox/issues/1080)) ([0d4740f](https://github.com/bpg/terraform-provider-proxmox/commit/0d4740fb90dad40c16994269e03de8b73ffee5dd))
* **deps:** update module github.com/stretchr/testify (v1.8.4 → v1.9.0) ([#1081](https://github.com/bpg/terraform-provider-proxmox/issues/1081)) ([dbd1655](https://github.com/bpg/terraform-provider-proxmox/commit/dbd1655974b31f1fae1f4c02766ef35cca77fa1e))
* **deps:** update module golang.org/x/crypto (v0.19.0 → v0.20.0) ([#1082](https://github.com/bpg/terraform-provider-proxmox/issues/1082)) ([e3ddd6f](https://github.com/bpg/terraform-provider-proxmox/commit/e3ddd6f5fa70728607849077fdc426d71bcf2338))
* switch to `terraform-plugin-testing` for acceptance tests ([#1067](https://github.com/bpg/terraform-provider-proxmox/issues/1067)) ([14fce33](https://github.com/bpg/terraform-provider-proxmox/commit/14fce3366da5cf3bca04511535a2898026c3210c))

## [0.47.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.46.6...v0.47.0) (2024-02-27)


### Features

* **file:** add `overwrite_unmanaged` attribute to `virtual_environment_download_file` resource ([#1064](https://github.com/bpg/terraform-provider-proxmox/issues/1064)) ([c64fcd2](https://github.com/bpg/terraform-provider-proxmox/commit/c64fcd2948bf6ffbcf6c907fb3f15931a5595596))


### Bug Fixes

* **provider:** race condition in`~/.ssh` path existence check ([#1052](https://github.com/bpg/terraform-provider-proxmox/issues/1052)) ([f7f67db](https://github.com/bpg/terraform-provider-proxmox/commit/f7f67dbd3d3edb2b6e092b77c898962d7641256f))
* **user:** `expiration_date` attribute handling ([#1066](https://github.com/bpg/terraform-provider-proxmox/issues/1066)) ([3c52760](https://github.com/bpg/terraform-provider-proxmox/commit/3c5276093a6edc2282512aa8a489b7d5ad4eee51))


### Miscellaneous

* **ci:** update actions/create-github-app-token action (v1.8.0 → v1.8.1) ([#1063](https://github.com/bpg/terraform-provider-proxmox/issues/1063)) ([9b52c12](https://github.com/bpg/terraform-provider-proxmox/commit/9b52c127ba11a4e01f7d63e2b1d06d7090cbadcb))
* **deps:** update golang.org/x/exp digest (v0.0.0-20240213143201-ec583247a57a → ) ([#1057](https://github.com/bpg/terraform-provider-proxmox/issues/1057)) ([4959480](https://github.com/bpg/terraform-provider-proxmox/commit/4959480f02f08354bfc009128ddc33c25aa22cae))
* **deps:** update module github.com/brianvoe/gofakeit/v7 (v7.0.0 → v7.0.1) ([#1058](https://github.com/bpg/terraform-provider-proxmox/issues/1058)) ([190ec39](https://github.com/bpg/terraform-provider-proxmox/commit/190ec39234bbe9a2d51f8cafa343dfa66df88e66))
* **vm:** refactor: move disks code out of vm.go ([#1062](https://github.com/bpg/terraform-provider-proxmox/issues/1062)) ([493ad1c](https://github.com/bpg/terraform-provider-proxmox/commit/493ad1c1219e666e61e05a6ad50a5fe746b4a69c))

## [0.46.6](https://github.com/bpg/terraform-provider-proxmox/compare/v0.46.5...v0.46.6) (2024-02-21)


### Bug Fixes

* **vm:** regression: `mac_addresses` list is missing some interfaces ([#1049](https://github.com/bpg/terraform-provider-proxmox/issues/1049)) ([518e25e](https://github.com/bpg/terraform-provider-proxmox/commit/518e25efaf6db6863d34ea3d83432eb0cd54d18a))


### Miscellaneous

* **lxc,vm:** refactor: move vm and container code to subpackages ([#1046](https://github.com/bpg/terraform-provider-proxmox/issues/1046)) ([0791194](https://github.com/bpg/terraform-provider-proxmox/commit/079119444d9f5a4c1266a4859c1aabe416c70b5d))

## [0.46.5](https://github.com/bpg/terraform-provider-proxmox/compare/v0.46.4...v0.46.5) (2024-02-20)


### Bug Fixes

* **lxc:** panic on empty `initialization.ip_config.ipv4|6` block ([#1043](https://github.com/bpg/terraform-provider-proxmox/issues/1043)) ([69c4a66](https://github.com/bpg/terraform-provider-proxmox/commit/69c4a66345547b79f4e1add7cb34d04125c6d451))
* **lxc:** panic when handling `network_interface.firewall` attribute ([#1042](https://github.com/bpg/terraform-provider-proxmox/issues/1042)) ([eb3e374](https://github.com/bpg/terraform-provider-proxmox/commit/eb3e3744321c2f5abc796b5e21e263703cff8916))


### Miscellaneous

* **deps:** Update module github.com/brianvoe/gofakeit/v6 (v6.28.0 → v7.0.0) ([#1044](https://github.com/bpg/terraform-provider-proxmox/issues/1044)) ([7fda43f](https://github.com/bpg/terraform-provider-proxmox/commit/7fda43f4ea78695d4c962b99df196fa0a1535dc5))
* **docs:** update README.md ([#1045](https://github.com/bpg/terraform-provider-proxmox/issues/1045)) ([8e620dc](https://github.com/bpg/terraform-provider-proxmox/commit/8e620dc59b3562de84d94e9088c82158663a3b8c))
* **vm:** refactoring, add acceptance tests ([#1040](https://github.com/bpg/terraform-provider-proxmox/issues/1040)) ([b648e5b](https://github.com/bpg/terraform-provider-proxmox/commit/b648e5bcb0ca21874aa7d7a081995ff0d7bc1040))

## [0.46.4](https://github.com/bpg/terraform-provider-proxmox/compare/v0.46.3...v0.46.4) (2024-02-16)


### Bug Fixes

* **vm:** fix panic when a config block is empty ([#1033](https://github.com/bpg/terraform-provider-proxmox/issues/1033)) ([027cf1e](https://github.com/bpg/terraform-provider-proxmox/commit/027cf1e81a2ab25f9d934921c6510d091870e3ee))
* **vm:** multi-line description field is always marked as changed ([#1030](https://github.com/bpg/terraform-provider-proxmox/issues/1030)) ([797873b](https://github.com/bpg/terraform-provider-proxmox/commit/797873b257614246fbadf167e7649cc5ed8e17e8))


### Miscellaneous

* **ci:** update actions/create-github-app-token action (v1.7.0 → v1.8.0) ([#1022](https://github.com/bpg/terraform-provider-proxmox/issues/1022)) ([0469192](https://github.com/bpg/terraform-provider-proxmox/commit/046919275607986c4ff380a846171f0c56e5e5f2))
* **ci:** update dorny/paths-filter action (v3.0.0 → v3.0.1) ([#1032](https://github.com/bpg/terraform-provider-proxmox/issues/1032)) ([d444202](https://github.com/bpg/terraform-provider-proxmox/commit/d444202ab8b2f80f6d144d46ea2af55f25aa8af7))
* **ci:** update mergify config to auto-approve renovate PRs ([#1023](https://github.com/bpg/terraform-provider-proxmox/issues/1023)) ([dfb95a8](https://github.com/bpg/terraform-provider-proxmox/commit/dfb95a85f437c3e414f2e8c7020d0077ebe01bc7))
* **deps:** update golang.org/x/exp digest (v0.0.0-20240205201215-2c58cdc269a3 → ) ([#1031](https://github.com/bpg/terraform-provider-proxmox/issues/1031)) ([4fab30e](https://github.com/bpg/terraform-provider-proxmox/commit/4fab30e5dfd62d63e29986b86dca57943f13d8af))
* **deps:** update module golang.org/x/crypto (v0.18.0 → v0.19.0) ([#1018](https://github.com/bpg/terraform-provider-proxmox/issues/1018)) ([34d31e2](https://github.com/bpg/terraform-provider-proxmox/commit/34d31e2ed080dc944900f5219338dbe9846a3aad))
* **deps:** update module golang.org/x/net (v0.20.0 → v0.21.0) ([#1020](https://github.com/bpg/terraform-provider-proxmox/issues/1020)) ([ed3bdb5](https://github.com/bpg/terraform-provider-proxmox/commit/ed3bdb5187dbf5588eedfc8d9ed193ab108edd64))
* **docs:** update links disk image link in examples ([#1028](https://github.com/bpg/terraform-provider-proxmox/issues/1028)) ([62a2130](https://github.com/bpg/terraform-provider-proxmox/commit/62a2130554c9ad09a7406d40e19678c4471f9364))

## [0.46.3](https://github.com/bpg/terraform-provider-proxmox/compare/v0.46.2...v0.46.3) (2024-02-07)


### Bug Fixes

* **file:** `error moving file` when uploading snippets ([#1013](https://github.com/bpg/terraform-provider-proxmox/issues/1013)) ([b6fbdcf](https://github.com/bpg/terraform-provider-proxmox/commit/b6fbdcf5ab3c191136c60814404153785aec806b))


### Miscellaneous

* **deps:** update module go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp (v0.47.0 → v0.48.0) in /tools ([#1014](https://github.com/bpg/terraform-provider-proxmox/issues/1014)) ([303b7da](https://github.com/bpg/terraform-provider-proxmox/commit/303b7da684dcf2f986fd6a70a74e40b75d71911a))

## [0.46.2](https://github.com/bpg/terraform-provider-proxmox/compare/v0.46.1...v0.46.2) (2024-02-06)


### Bug Fixes

* **docs:** update sudo configuration to a more restrictive variant ([#1001](https://github.com/bpg/terraform-provider-proxmox/issues/1001)) ([6bd8ba5](https://github.com/bpg/terraform-provider-proxmox/commit/6bd8ba566a60c18121d9a66f1cdd056878fe6114))
* **file:** use `sudo` for snippets upload ([#1004](https://github.com/bpg/terraform-provider-proxmox/issues/1004)) ([60fb679](https://github.com/bpg/terraform-provider-proxmox/commit/60fb679e9f31b3be3e05bb9b25a0deb0ab37c48c))
* **vm:** error when creating custom disks on PVE with non-default shell ([#983](https://github.com/bpg/terraform-provider-proxmox/issues/983)) ([1f333ea](https://github.com/bpg/terraform-provider-proxmox/commit/1f333ea097f43097e3847d08153145ac2a44faad))
* **vm:** panic at import / state refresh if disk size is not set ([#994](https://github.com/bpg/terraform-provider-proxmox/issues/994)) ([363e502](https://github.com/bpg/terraform-provider-proxmox/commit/363e502a567f8c75c45b682795ce5974e993d082))


### Miscellaneous

* **ci:** update lycheeverse/lychee-action action (v1.9.2 → v1.9.3) ([#999](https://github.com/bpg/terraform-provider-proxmox/issues/999)) ([f8004b0](https://github.com/bpg/terraform-provider-proxmox/commit/f8004b0e2a35616b94804fdb272df598cd2b88a2))
* **deps:** update golang.org/x/exp digest (v0.0.0-20240119083558-1b970713d09a → ) ([9cfd383](https://github.com/bpg/terraform-provider-proxmox/commit/9cfd3833da3a8c38ef5800fae1e65f3cd6d3b696))
* **deps:** update golang.org/x/exp digest (v0.0.0-20240119083558-1b970713d09a 1b97071 → 2c58cdc) ([#1007](https://github.com/bpg/terraform-provider-proxmox/issues/1007)) ([9cfd383](https://github.com/bpg/terraform-provider-proxmox/commit/9cfd3833da3a8c38ef5800fae1e65f3cd6d3b696))
* **deps:** update module github.com/goreleaser/goreleaser (v1.23.0 → v1.24.0) in /tools [security] ([#1006](https://github.com/bpg/terraform-provider-proxmox/issues/1006)) ([e132f5a](https://github.com/bpg/terraform-provider-proxmox/commit/e132f5af4bbb892efd130777f7046269ddb0cfa6))
* **deps:** update module github.com/hashicorp/terraform-plugin-mux (v0.13.0 → v0.14.0) ([#989](https://github.com/bpg/terraform-provider-proxmox/issues/989)) ([eb6377e](https://github.com/bpg/terraform-provider-proxmox/commit/eb6377e6fdbd84d3cbd59b254106f94325dbc479))
* **deps:** update module github.com/hashicorp/terraform-plugin-sdk/v2 (v2.31.0 → v2.32.0) ([#990](https://github.com/bpg/terraform-provider-proxmox/issues/990)) ([c1eeefb](https://github.com/bpg/terraform-provider-proxmox/commit/c1eeefbb1214ef9c14341eac94b9469e7161e96f))
* **deps:** update module golang.org/x/net (v0.18.0 → v0.20.0) ([#994](https://github.com/bpg/terraform-provider-proxmox/issues/994)) ([b196cdb](https://github.com/bpg/terraform-provider-proxmox/commit/b196cdb65bed27c34a755c3bab1654f71ef4a5e6))

## [0.46.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.46.0...v0.46.1) (2024-01-28)


### Bug Fixes

* **docs:** fix documentation tree structure in the TF registry ([#980](https://github.com/bpg/terraform-provider-proxmox/issues/980)) ([49a76bb](https://github.com/bpg/terraform-provider-proxmox/commit/49a76bb1a10c56ab2537e83b4b9fb20d2c7c9b9e))

## [0.46.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.45.1...v0.46.0) (2024-01-28)


### Features

* **docs:** rename howtos -&gt; guides and publish to the Terraform Registry ([#971](https://github.com/bpg/terraform-provider-proxmox/issues/971)) ([c39494b](https://github.com/bpg/terraform-provider-proxmox/commit/c39494b939afb0e1316776eea8730f2545135b4b))
* **provider:** add SOCKS5 proxy support for SSH connections ([#970](https://github.com/bpg/terraform-provider-proxmox/issues/970)) ([da1d780](https://github.com/bpg/terraform-provider-proxmox/commit/da1d7804af6b2ad6d6a1d698e52d19de3c1d5cb6))


### Bug Fixes

* **docs:** fix broken links ([#976](https://github.com/bpg/terraform-provider-proxmox/issues/976)) ([0e2eb80](https://github.com/bpg/terraform-provider-proxmox/commit/0e2eb80e9f4f2e13678ad9ca6afb3cf5de4d5f19))
* **lxc:** panic on empty `initialization.ip_config` block ([#977](https://github.com/bpg/terraform-provider-proxmox/issues/977)) ([0253eb9](https://github.com/bpg/terraform-provider-proxmox/commit/0253eb97576c6f6b06e4cf652b5c1e74ad20639d))
* **pool:** missing `pool_id` after import ([#974](https://github.com/bpg/terraform-provider-proxmox/issues/974)) ([ed33a18](https://github.com/bpg/terraform-provider-proxmox/commit/ed33a18c9b6499ff33bacb79cebfd510b24a29c8))
* **vm:** `timeout_start_vm` is ignored ([#978](https://github.com/bpg/terraform-provider-proxmox/issues/978)) ([625bdb6](https://github.com/bpg/terraform-provider-proxmox/commit/625bdb696f5c41f76c12f5572c89bb4594f81853))


### Miscellaneous

* **ci:** update actions/create-github-app-token action (v1.6.4 → v1.7.0) ([2fad644](https://github.com/bpg/terraform-provider-proxmox/commit/2fad644ffd732283875ab38f70c852bf9723c409))
* **ci:** Update dorny/paths-filter action (v2.12.0 → v3.0.0) ([#959](https://github.com/bpg/terraform-provider-proxmox/issues/959)) ([3790b52](https://github.com/bpg/terraform-provider-proxmox/commit/3790b522e71ccbab69ea0549b0a12b390bfb8848))
* **ci:** update lycheeverse/lychee-action action (v1.9.1 → v1.9.2) ([105a694](https://github.com/bpg/terraform-provider-proxmox/commit/105a694ddf8ff3c8e60440d995ffa2f42cc70788))
* **ci:** Update peter-evans/create-issue-from-file action (v4.0.1 → v5.0.0) ([#960](https://github.com/bpg/terraform-provider-proxmox/issues/960)) ([2ec8c1d](https://github.com/bpg/terraform-provider-proxmox/commit/2ec8c1d2cdd518d34245ce120937e852a6eedea0))

## [0.45.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.45.0...v0.45.1) (2024-01-27)


### Bug Fixes

* **docs:** inconsistent indentations in docs ([#961](https://github.com/bpg/terraform-provider-proxmox/issues/961)) ([0d548a7](https://github.com/bpg/terraform-provider-proxmox/commit/0d548a78078ee9ae3e0653ea8f5e75b228dc17ac))
* **docs:** update HOW-TOs for cloud-init ([#955](https://github.com/bpg/terraform-provider-proxmox/issues/955)) ([d91ec25](https://github.com/bpg/terraform-provider-proxmox/commit/d91ec25bfae08e6f24bb9923c0ba962792e765db))
* **vm:** regression: `sudo: command not found` when creating a VM ([#966](https://github.com/bpg/terraform-provider-proxmox/issues/966)) ([01a8f97](https://github.com/bpg/terraform-provider-proxmox/commit/01a8f9779c87a844f7d74ccaa8f9a3d4bc28bb55))


### Miscellaneous

* **ci:** update dorny/paths-filter action (v2.11.1 → v2.12.0) ([#958](https://github.com/bpg/terraform-provider-proxmox/issues/958)) ([3a5e69d](https://github.com/bpg/terraform-provider-proxmox/commit/3a5e69d9c8e647a72b9e6141fe3e2d0f2363c991))
* **deps:** update module github.com/google/uuid (v1.5.0 → v1.6.0) ([#954](https://github.com/bpg/terraform-provider-proxmox/issues/954)) ([b6474f8](https://github.com/bpg/terraform-provider-proxmox/commit/b6474f8ddbd8c1d3564c7d2f2bbe5a996862d443))
* **deps:** update module github.com/hashicorp/terraform-plugin-docs (v0.17.0 → v0.18.0) in /tools ([#957](https://github.com/bpg/terraform-provider-proxmox/issues/957)) ([4a03a78](https://github.com/bpg/terraform-provider-proxmox/commit/4a03a78dcd6d350b3a17fccabffb85b23c7f9fc3))
* **deps:** update module github.com/hashicorp/terraform-plugin-go (v0.20.0 → v0.21.0) ([#964](https://github.com/bpg/terraform-provider-proxmox/issues/964)) ([63e7bfc](https://github.com/bpg/terraform-provider-proxmox/commit/63e7bfc042bd5f4b60f6fbf70c7fdfd344b91b05))

## [0.45.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.44.0...v0.45.0) (2024-01-22)


### Features

* **provider:** use `sudo` when running commands over SSH ([#950](https://github.com/bpg/terraform-provider-proxmox/issues/950)) ([9d764e5](https://github.com/bpg/terraform-provider-proxmox/commit/9d764e588976a5b1d35662501de2a6bc804fb693))


### Miscellaneous

* **ci:** add link checker, reformat actions code ([#944](https://github.com/bpg/terraform-provider-proxmox/issues/944)) ([a030542](https://github.com/bpg/terraform-provider-proxmox/commit/a030542da0524caa7f7bdb996892b18c78d45804))
* **ci:** update actions/create-github-app-token action (v1.6.3 → v1.6.4) ([#939](https://github.com/bpg/terraform-provider-proxmox/issues/939)) ([25db34b](https://github.com/bpg/terraform-provider-proxmox/commit/25db34b149f29b25935ae14245fc97837bffc0d6))
* **ci:** Update google-github-actions/release-please-action action (v3.7.13 → v4.0.2) ([#905](https://github.com/bpg/terraform-provider-proxmox/issues/905)) ([d4832b3](https://github.com/bpg/terraform-provider-proxmox/commit/d4832b3d5991c6b4610dacae7c43a31dea3f94ee))
* **ci:** update issue templates, renovate config ([#951](https://github.com/bpg/terraform-provider-proxmox/issues/951)) ([9644590](https://github.com/bpg/terraform-provider-proxmox/commit/96445909989fbb65a8a28aad4f98ce072db93e79))
* **docs:** move list of contributors to CONTRIBUTORS.md ([#945](https://github.com/bpg/terraform-provider-proxmox/issues/945)) ([aabfeb8](https://github.com/bpg/terraform-provider-proxmox/commit/aabfeb86a204bdd109b885e7e1cda84eff42d8a5))
* **docs:** update README.md, add note about OpenTofu support ([#943](https://github.com/bpg/terraform-provider-proxmox/issues/943)) ([b926c57](https://github.com/bpg/terraform-provider-proxmox/commit/b926c57a53002f955651dde8e95ac3734d453e8f))

## [0.44.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.43.3...v0.44.0) (2024-01-20)


### Features

* **lxc:** add container startup options ([#923](https://github.com/bpg/terraform-provider-proxmox/issues/923)) ([c9c3067](https://github.com/bpg/terraform-provider-proxmox/commit/c9c3067b61bf6fe7930c6d8281040aa382eac09d))
* **provider:** add min_tls option to provider config ([#931](https://github.com/bpg/terraform-provider-proxmox/issues/931)) ([01ff2cb](https://github.com/bpg/terraform-provider-proxmox/commit/01ff2cb7dba6e74e5aae51114dd13883740d028f))


### Bug Fixes

* **vm:** panic on empty `initialization.dns` block ([#928](https://github.com/bpg/terraform-provider-proxmox/issues/928)) ([e5bccbc](https://github.com/bpg/terraform-provider-proxmox/commit/e5bccbc53de66f73b95e92f00a80ba98af6becf1))


### Miscellaneous

* **deps:** update golang.org/x/exp digest (v0.0.0-20240112132812-db7319d0e0e3 → ) ([#934](https://github.com/bpg/terraform-provider-proxmox/issues/934)) ([3ffd230](https://github.com/bpg/terraform-provider-proxmox/commit/3ffd2306828af30ffd25aaa753ed086700bd71a2))
* **deps:** update module github.com/brianvoe/gofakeit/v6 (v6.27.0 → v6.28.0) ([#937](https://github.com/bpg/terraform-provider-proxmox/issues/937)) ([c1e9c08](https://github.com/bpg/terraform-provider-proxmox/commit/c1e9c089ba921bc522a363531cb8835dd14fc30a))
* **deps:** update module github.com/hashicorp/terraform-plugin-docs (v0.16.0 → v0.17.0) in /tools ([#922](https://github.com/bpg/terraform-provider-proxmox/issues/922)) ([c8e298c](https://github.com/bpg/terraform-provider-proxmox/commit/c8e298cc4c6071f47fdaf90328548b5e690b674b))
* **deps:** update module go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp (v0.46.1 → v0.47.0) in /tools ([#933](https://github.com/bpg/terraform-provider-proxmox/issues/933)) ([9326131](https://github.com/bpg/terraform-provider-proxmox/commit/932613110dc0ee0f4a6c438a910d14b64761f8a2))
* **docs:** remove static website generator with ruby dependencies ([#929](https://github.com/bpg/terraform-provider-proxmox/issues/929)) ([7d94bf7](https://github.com/bpg/terraform-provider-proxmox/commit/7d94bf73ec37bed1802cc2a37399832498ee35e7))
* **docs:** update activesupport (7.1.2 → 7.1.3) ([#925](https://github.com/bpg/terraform-provider-proxmox/issues/925)) ([85109cb](https://github.com/bpg/terraform-provider-proxmox/commit/85109cbe3d3c2cfa8068e56978eae7a1472f9cc5))
* **docs:** update jekyll (3.9.3 → 3.9.4) ([#921](https://github.com/bpg/terraform-provider-proxmox/issues/921)) ([93283ef](https://github.com/bpg/terraform-provider-proxmox/commit/93283ef3ab684d6155202d1fd62190a73fee1792))
* **docs:** update terraform proxmox (0.43.2 → 0.43.3) ([#919](https://github.com/bpg/terraform-provider-proxmox/issues/919)) ([5cffafc](https://github.com/bpg/terraform-provider-proxmox/commit/5cffafc26e3d9e7b668cd53232c78006f757faea))

## [0.43.3](https://github.com/bpg/terraform-provider-proxmox/compare/v0.43.2...v0.43.3) (2024-01-16)


### Bug Fixes

* **docs:** fix indentation in `virtual_environment_container.md` ([#882](https://github.com/bpg/terraform-provider-proxmox/issues/882)) ([10dbfdd](https://github.com/bpg/terraform-provider-proxmox/commit/10dbfddc57c5dc3245b9a1827ec3f5d43f783e21))


### Miscellaneous

* **ci:** switch to renovate ([#891](https://github.com/bpg/terraform-provider-proxmox/issues/891)) ([01e6698](https://github.com/bpg/terraform-provider-proxmox/commit/01e669854bb4044afcf22144a1b6e3c4cbfe92b5))
* **ci:** update ([#890](https://github.com/bpg/terraform-provider-proxmox/issues/890)) ([c635044](https://github.com/bpg/terraform-provider-proxmox/commit/c635044db341422b458202a62538cffdaadb5fcc))
* **ci:** update dorny/paths-filter action ( v2.2.1 → v2.11.1 ) ([#911](https://github.com/bpg/terraform-provider-proxmox/issues/911)) ([daa94d4](https://github.com/bpg/terraform-provider-proxmox/commit/daa94d4f8791d68747cf9be0dc7451fb466833bb))
* **ci:** update dorny/paths-filter digest ( 4512585 → 3b817c9 ) ([#910](https://github.com/bpg/terraform-provider-proxmox/issues/910)) ([5574e60](https://github.com/bpg/terraform-provider-proxmox/commit/5574e60542861e5b3010f427aa295f780ce90437))
* **ci:** update renovate config ([8226c42](https://github.com/bpg/terraform-provider-proxmox/commit/8226c421f5e99fc4e8ca6254a9b0738d6395df1b))
* **deps:** bump github.com/hashicorp/terraform-plugin-framework from 1.4.2 to 1.5.0 ([#889](https://github.com/bpg/terraform-provider-proxmox/issues/889)) ([c7bbb47](https://github.com/bpg/terraform-provider-proxmox/commit/c7bbb47223bce86a653f236841fa1285f2a5dfe5))
* **deps:** update activesupport (7.0.7.1 → 7.1.2) ([#897](https://github.com/bpg/terraform-provider-proxmox/issues/897)) ([9a0b897](https://github.com/bpg/terraform-provider-proxmox/commit/9a0b8979befe5875b67315f231fbebc4ed7f0d63))
* **deps:** update github.com/hashicorp/go-cty digest ( d3edf31 → 8598007 ) ([#892](https://github.com/bpg/terraform-provider-proxmox/issues/892)) ([34cb5a7](https://github.com/bpg/terraform-provider-proxmox/commit/34cb5a7c4eae16695b0e2ae84d078c783d6bf78f))
* **deps:** update golang.org/x/exp digest ( 9212866 → db7319d ) ([#893](https://github.com/bpg/terraform-provider-proxmox/issues/893)) ([21264c0](https://github.com/bpg/terraform-provider-proxmox/commit/21264c039af6ee714006e3f6f84798ade22cf46a))
* **deps:** update jekyll (3.9.3 → 3.9.4) ([#894](https://github.com/bpg/terraform-provider-proxmox/issues/894)) ([65f429e](https://github.com/bpg/terraform-provider-proxmox/commit/65f429e81cae3f83cf989eac4bf5d0a846459d90))
* **deps:** update just-the-docs (0.5.4 → 0.7.0) ([#898](https://github.com/bpg/terraform-provider-proxmox/issues/898)) ([2f074d6](https://github.com/bpg/terraform-provider-proxmox/commit/2f074d6b2d5bc046036b8eaac66396ef924dc0a7))
* **deps:** update module github.com/brianvoe/gofakeit/v6 (v6.26.4 → v6.27.0) ([#900](https://github.com/bpg/terraform-provider-proxmox/issues/900)) ([c500cc5](https://github.com/bpg/terraform-provider-proxmox/commit/c500cc5b9b54257cb8fd551034efa5a206ce84bb))
* **deps:** update module github.com/nats-io/nkeys (v0.4.6 → v0.4.7) ([#895](https://github.com/bpg/terraform-provider-proxmox/issues/895)) ([268722c](https://github.com/bpg/terraform-provider-proxmox/commit/268722c1214627298aa1e34da5cd8d7ba73b20a2))
* **deps:** update module go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp (v0.44.0 → v0.46.1) ([#901](https://github.com/bpg/terraform-provider-proxmox/issues/901)) ([8d5ef9a](https://github.com/bpg/terraform-provider-proxmox/commit/8d5ef9a73fbf99b112729ea34235d8c0ee259bf9))
* **deps:** update terraform local (2.2.2 → 2.4.1) ([#902](https://github.com/bpg/terraform-provider-proxmox/issues/902)) ([c116db5](https://github.com/bpg/terraform-provider-proxmox/commit/c116db592b7b0d7dc345f84589075e1c9e679811))
* **deps:** update tzinfo-data (1.2023.3 → 1.2023.4) ([#896](https://github.com/bpg/terraform-provider-proxmox/issues/896)) ([2edf2cb](https://github.com/bpg/terraform-provider-proxmox/commit/2edf2cbb1d5d4996c4f33f8a5be18983ffe3e9d9))
* **docs:** update terraform proxmox (0.38.1 → 0.43.2) in docs ([#903](https://github.com/bpg/terraform-provider-proxmox/issues/903)) ([5d9f41c](https://github.com/bpg/terraform-provider-proxmox/commit/5d9f41c877e82a500da212c03cb611d5530f290a))
* **docs:** update terraform tls (3.1.0 → 3.4.0) in docs ([#904](https://github.com/bpg/terraform-provider-proxmox/issues/904)) ([699f19d](https://github.com/bpg/terraform-provider-proxmox/commit/699f19d135529eb9772bd0dcacd8169e01abd1d7))
* **docs:** Update Terraform tls (3.4.0 → 4.0.5) in docs ([#908](https://github.com/bpg/terraform-provider-proxmox/issues/908)) ([9e7d7d1](https://github.com/bpg/terraform-provider-proxmox/commit/9e7d7d17cb62bffcacacbfab64e2473f33d2f086))

## [0.43.2](https://github.com/bpg/terraform-provider-proxmox/compare/v0.43.1...v0.43.2) (2024-01-11)


### Bug Fixes

* **provider:** node DNS lookup fallback does not produce an IP ([#874](https://github.com/bpg/terraform-provider-proxmox/issues/874)) ([e436427](https://github.com/bpg/terraform-provider-proxmox/commit/e436427e00bd39ffe0df4ae7d6c3f445f0a0cb31))
* **vm:** missing disks when importing VM to a TF state ([#877](https://github.com/bpg/terraform-provider-proxmox/issues/877)) ([a8bf497](https://github.com/bpg/terraform-provider-proxmox/commit/a8bf497c7f3331e0c92501d479ccf04a8481e926))


### Miscellaneous

* **deps:** bump github.com/brianvoe/gofakeit/v6 from 6.26.3 to 6.26.4 ([#879](https://github.com/bpg/terraform-provider-proxmox/issues/879)) ([6aa56b3](https://github.com/bpg/terraform-provider-proxmox/commit/6aa56b3f6e7639ddd1190d84e4d53bd6c516a977))
* **deps:** bump golang.org/x/crypto from 0.17.0 to 0.18.0 ([#878](https://github.com/bpg/terraform-provider-proxmox/issues/878)) ([0f198eb](https://github.com/bpg/terraform-provider-proxmox/commit/0f198eb66b44d6aea23f6489587c404ef3d7ffdf))

## [0.43.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.43.0...v0.43.1) (2024-01-10)


### Bug Fixes

* **docs:** typos in `proxmox_virtual_environment_file` resource ([#872](https://github.com/bpg/terraform-provider-proxmox/issues/872)) ([74e0ef3](https://github.com/bpg/terraform-provider-proxmox/commit/74e0ef3b1e37c02b8671fb650b4593c378bf96d1))
* **vm:** optimize retrieval of VM volume attributes from a datastore ([#862](https://github.com/bpg/terraform-provider-proxmox/issues/862)) ([613be84](https://github.com/bpg/terraform-provider-proxmox/commit/613be842bee37eec4d0f74ddfa91a3a0bf8db43a))


### Miscellaneous

* **deps:** bump github.com/cloudflare/circl from 1.3.3 to 1.3.7 ([#869](https://github.com/bpg/terraform-provider-proxmox/issues/869)) ([ea653e1](https://github.com/bpg/terraform-provider-proxmox/commit/ea653e1f253655c0a97677376cbab8544b2a9c3c))
* **deps:** bump github.com/cloudflare/circl from 1.3.5 to 1.3.7 in /tools ([#870](https://github.com/bpg/terraform-provider-proxmox/issues/870)) ([ffafa06](https://github.com/bpg/terraform-provider-proxmox/commit/ffafa063af28e4c7b7d5180c82a937b9abd17ccb))

## [0.43.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.42.1...v0.43.0) (2024-01-04)


### Features

* **provider:** add DNS lookup fallback for node IP resolution ([#848](https://github.com/bpg/terraform-provider-proxmox/issues/848)) ([d398c9c](https://github.com/bpg/terraform-provider-proxmox/commit/d398c9c102fc2f6741b3e6d574fcfb8a4f7f49aa))
* **storage:** add new resource `proxmox_virtual_environment_download_file`  ([#837](https://github.com/bpg/terraform-provider-proxmox/issues/837)) ([58347c0](https://github.com/bpg/terraform-provider-proxmox/commit/58347c09fe012e35025613923d95c5aa8340318a))


### Miscellaneous

* **deps:** bump crazy-max/ghaction-import-gpg from 6.0.0 to 6.1.0 ([#855](https://github.com/bpg/terraform-provider-proxmox/issues/855)) ([620bb84](https://github.com/bpg/terraform-provider-proxmox/commit/620bb84635d38181e841a8d41cd8a01dd4afd83b))
* **deps:** bump github.com/goreleaser/goreleaser from 1.22.1 to 1.23.0 in /tools ([#854](https://github.com/bpg/terraform-provider-proxmox/issues/854)) ([3914bc2](https://github.com/bpg/terraform-provider-proxmox/commit/3914bc28b64decbcc853c0e6aa3188b3343ebf81))
* **docs:** update provider documentation with more details about token use ([#846](https://github.com/bpg/terraform-provider-proxmox/issues/846)) ([2677445](https://github.com/bpg/terraform-provider-proxmox/commit/2677445802bf792fbd2b92ec8120f1ddacdb299a))

## [0.42.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.42.0...v0.42.1) (2023-12-29)


### Bug Fixes

* **lxc:** add missing `onboot` param on container clone create ([#838](https://github.com/bpg/terraform-provider-proxmox/issues/838)) ([40102a6](https://github.com/bpg/terraform-provider-proxmox/commit/40102a6a501a5ead3219492f34255a25c4f21371))
* **vm,lxc:** accept IPv6 in `initialization.dns.servers` attribute ([#842](https://github.com/bpg/terraform-provider-proxmox/issues/842)) ([bf5cbd9](https://github.com/bpg/terraform-provider-proxmox/commit/bf5cbd9dad116a4515bd2eb193c296097b1e4b84))
* **vm,lxc:** unexpected state drift when using `initialization.dns.servers` ([#844](https://github.com/bpg/terraform-provider-proxmox/issues/844)) ([ac923cd](https://github.com/bpg/terraform-provider-proxmox/commit/ac923cd1b42c0c64d9829beb1ab552680b21d98b))
* **vm:** Fixed missing default for disk discard ([#840](https://github.com/bpg/terraform-provider-proxmox/issues/840)) ([5281ac2](https://github.com/bpg/terraform-provider-proxmox/commit/5281ac24921795ed933047e5d9ca953add15bdd0))


### Miscellaneous

* **deps:** bump github.com/go-git/go-git/v5 from 5.7.0 to 5.11.0 in /tools ([#839](https://github.com/bpg/terraform-provider-proxmox/issues/839)) ([f860c4b](https://github.com/bpg/terraform-provider-proxmox/commit/f860c4bab54344beb4fd54366adcf940ea1463fe))
* **tests:** Update acceptance tests to PVE 8.1, add docs ([#834](https://github.com/bpg/terraform-provider-proxmox/issues/834)) ([d8f82d4](https://github.com/bpg/terraform-provider-proxmox/commit/d8f82d47b3a74e4b64a26757522c067a635e4fa3))

## [0.42.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.41.0...v0.42.0) (2023-12-23)


### Features

* **vm, lxc:** add new `initialization.dns.servers` param to vm and container ([#832](https://github.com/bpg/terraform-provider-proxmox/issues/832)) ([16e571d](https://github.com/bpg/terraform-provider-proxmox/commit/16e571dc199c8977b3954e3c56c9b96cc351503e))
* **vm:** add new dns servers param to vm and container, deprecated server param ([16e571d](https://github.com/bpg/terraform-provider-proxmox/commit/16e571dc199c8977b3954e3c56c9b96cc351503e))
* **vm:** add support for up to 32 network interfaces ([#822](https://github.com/bpg/terraform-provider-proxmox/issues/822)) ([4113bec](https://github.com/bpg/terraform-provider-proxmox/commit/4113bec1b5184cd30c0435ae50470a8f7ab3ba39))


### Bug Fixes

* **provider:** allow FQDN for `ssh.node.address` in provider's config ([#824](https://github.com/bpg/terraform-provider-proxmox/issues/824)) ([34df977](https://github.com/bpg/terraform-provider-proxmox/commit/34df9773c34b43ba39b5d8505b5916b52f87ff3e))
* **vm:** update `smbios` during clone ([#827](https://github.com/bpg/terraform-provider-proxmox/issues/827)) ([0ffe75a](https://github.com/bpg/terraform-provider-proxmox/commit/0ffe75afa44995d4b648687281974e990029977e))


### Miscellaneous

* **deps:** bump golang.org/x/crypto from 0.14.0 to 0.17.0 in /tools ([#819](https://github.com/bpg/terraform-provider-proxmox/issues/819)) ([21a4b01](https://github.com/bpg/terraform-provider-proxmox/commit/21a4b01cd16fec6d3f04ae0e2c7eab9a021ee1e6))
* **deps:** bump golang.org/x/crypto from 0.16.0 to 0.17.0 ([#820](https://github.com/bpg/terraform-provider-proxmox/issues/820)) ([ec31d75](https://github.com/bpg/terraform-provider-proxmox/commit/ec31d75fe1a93e110f4e21108c4f69d12c9a38d7))
* **docs:** improve make example docs and add proxmox setup how-to ([#829](https://github.com/bpg/terraform-provider-proxmox/issues/829)) ([4f54f89](https://github.com/bpg/terraform-provider-proxmox/commit/4f54f89b5db4cf37321b9c021a411b747093325f))

## [0.41.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.40.0...v0.41.0) (2023-12-18)


### Features

* **vm:** add `cpu.limit` attribute ([#814](https://github.com/bpg/terraform-provider-proxmox/issues/814)) ([9712952](https://github.com/bpg/terraform-provider-proxmox/commit/9712952e2614a9af6a5a35a4cf318af44684f063))
* **vm:** support stopping (rather than shutting down) VMs on resource destroy ([#783](https://github.com/bpg/terraform-provider-proxmox/issues/783)) ([6ebe8dc](https://github.com/bpg/terraform-provider-proxmox/commit/6ebe8dcc60be12276d9f2847fb9242e93be98441))


### Bug Fixes

* **docs:** add clone/full parameter for vms ([#797](https://github.com/bpg/terraform-provider-proxmox/issues/797)) ([86d0f07](https://github.com/bpg/terraform-provider-proxmox/commit/86d0f07e9b6023d0c3627f45f9944f26d26a4e1d))
* **provider:** typo in provider example ([#785](https://github.com/bpg/terraform-provider-proxmox/issues/785)) ([32bdc21](https://github.com/bpg/terraform-provider-proxmox/commit/32bdc2175076a4a3cc89bc0ff18035fb9b8aa4d6))
* **vm:** hostpci devices not showing up in refresh plan ([#578](https://github.com/bpg/terraform-provider-proxmox/issues/578)) ([aa939c7](https://github.com/bpg/terraform-provider-proxmox/commit/aa939c731f7bc36213b6d0abc51cc284a1295338))
* **vm:** panic at read when cloud-init drive is on directory storage ([#811](https://github.com/bpg/terraform-provider-proxmox/issues/811)) ([3e0ef1d](https://github.com/bpg/terraform-provider-proxmox/commit/3e0ef1d08b036297a5d8326aedce6c43c1200bb2))


### Miscellaneous

* **deps:** bump actions/setup-go from 4 to 5 ([#791](https://github.com/bpg/terraform-provider-proxmox/issues/791)) ([164a72d](https://github.com/bpg/terraform-provider-proxmox/commit/164a72d19d9c3a952364bfbedb5a4295e2fd48ea))
* **deps:** bump actions/stale from 8 to 9 ([#790](https://github.com/bpg/terraform-provider-proxmox/issues/790)) ([02b5da7](https://github.com/bpg/terraform-provider-proxmox/commit/02b5da705da682b9325a1ac882f30993c8f96bb0))
* **deps:** bump github.com/brianvoe/gofakeit/v6 from 6.26.0 to 6.26.3 ([#807](https://github.com/bpg/terraform-provider-proxmox/issues/807)) ([1d69c69](https://github.com/bpg/terraform-provider-proxmox/commit/1d69c691acdacf062406f27e5daf993c70ed04d8))
* **deps:** bump github.com/google/uuid from 1.4.0 to 1.5.0 ([#805](https://github.com/bpg/terraform-provider-proxmox/issues/805)) ([3b4a69e](https://github.com/bpg/terraform-provider-proxmox/commit/3b4a69edfae45eaa7edb0b5bff7310f79fe542be))
* **deps:** bump github.com/hashicorp/terraform-plugin-mux from 0.12.0 to 0.13.0 ([#806](https://github.com/bpg/terraform-provider-proxmox/issues/806)) ([53270e2](https://github.com/bpg/terraform-provider-proxmox/commit/53270e23108657e0e878859950205a6dcc7e9b1c))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.30.0 to 2.31.0 ([#808](https://github.com/bpg/terraform-provider-proxmox/issues/808)) ([5c91b91](https://github.com/bpg/terraform-provider-proxmox/commit/5c91b91938e0f7e3101a3472ae3be866b3ec0f26))

## [0.40.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.39.0...v0.40.0) (2023-12-06)


### ⚠ BREAKING CHANGES

* **lxc:** allow to update `features`, add mount type support ([#765](https://github.com/bpg/terraform-provider-proxmox/issues/765))

### Features

* **lxc:** allow to update `features`, add mount type support ([#765](https://github.com/bpg/terraform-provider-proxmox/issues/765)) ([8bf2609](https://github.com/bpg/terraform-provider-proxmox/commit/8bf26099e0da5db85f1997789cb867aa11db9906))
* **vm:** Add support for setting the VM TPM State device ([#743](https://github.com/bpg/terraform-provider-proxmox/issues/743)) ([66bba2a](https://github.com/bpg/terraform-provider-proxmox/commit/66bba2a0275e2e9e3a2c5c2de7414d89be89a53c))


### Bug Fixes

* **docs:** add more details about local testing of the provider ([#698](https://github.com/bpg/terraform-provider-proxmox/issues/698)) ([f1450cb](https://github.com/bpg/terraform-provider-proxmox/commit/f1450cb6dd13e291ce885130f0550cd26e97e99f))
* **lxc:** description is always showed as changed ([#762](https://github.com/bpg/terraform-provider-proxmox/issues/762)) ([d1f2093](https://github.com/bpg/terraform-provider-proxmox/commit/d1f2093d3977ff9d30b1af95f97e1fe601d22991))
* **lxc:** fixes for datastore-backed volume mounts ([#772](https://github.com/bpg/terraform-provider-proxmox/issues/772)) ([25deebb](https://github.com/bpg/terraform-provider-proxmox/commit/25deebba265ccea0031ea2261ee2e03f1c09f5d7))


### Miscellaneous

* configure vscode's linter to use proper .golangci.yml file ([#774](https://github.com/bpg/terraform-provider-proxmox/issues/774)) ([d0f43e1](https://github.com/bpg/terraform-provider-proxmox/commit/d0f43e1497325a5aafd915771c5da5d99f2c7ead))
* **deps:** bump github.com/brianvoe/gofakeit/v6 from 6.25.0 to 6.26.0 ([#775](https://github.com/bpg/terraform-provider-proxmox/issues/775)) ([006b5e9](https://github.com/bpg/terraform-provider-proxmox/commit/006b5e9caa51a24ae3a573abdb3bd7e21506e974))
* **docs:** update CONTRIBUTING.md and other project docs ([#771](https://github.com/bpg/terraform-provider-proxmox/issues/771)) ([7505b37](https://github.com/bpg/terraform-provider-proxmox/commit/7505b377087f08773c88819d287364c0f5be8d20))
* **docs:** update PR and issue templates ([#777](https://github.com/bpg/terraform-provider-proxmox/issues/777)) ([54288dd](https://github.com/bpg/terraform-provider-proxmox/commit/54288ddd76c4e96542921123c2f081aff51075be))

## [0.39.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.38.1...v0.39.0) (2023-11-30)


### Features

* **docs:** add initial mini-howtos for VM ([#730](https://github.com/bpg/terraform-provider-proxmox/issues/730)) ([e2717a9](https://github.com/bpg/terraform-provider-proxmox/commit/e2717a9a9ee542e7e17c0b518ccd1da78d5abdea))
* **provider:** modify the proxmox api client to support connecting through an https proxy ([#748](https://github.com/bpg/terraform-provider-proxmox/issues/748)) ([13e911c](https://github.com/bpg/terraform-provider-proxmox/commit/13e911cf592770740d88880e9e12d16f5d9bb8b6))
* **vm:** Support hook script ([#733](https://github.com/bpg/terraform-provider-proxmox/issues/733)) ([0eb04b2](https://github.com/bpg/terraform-provider-proxmox/commit/0eb04b2a250999893996ce62ea1b9109081494a7))


### Bug Fixes

* **cluster:** can't read back cluster options on PVE 8.1 ([#755](https://github.com/bpg/terraform-provider-proxmox/issues/755)) ([cd24cf2](https://github.com/bpg/terraform-provider-proxmox/commit/cd24cf238cb11c3ffb4b9e5378c52b04c6961068))
* **docs:** improve documentation for container feature flags ([#747](https://github.com/bpg/terraform-provider-proxmox/issues/747)) ([d5193b3](https://github.com/bpg/terraform-provider-proxmox/commit/d5193b3e9b0ddbaf0fde45381b2bc9d9e28bca18))
* **vm,lxc:** file ID validator to allow . in a storage name ([#750](https://github.com/bpg/terraform-provider-proxmox/issues/750)) ([a6fa40e](https://github.com/bpg/terraform-provider-proxmox/commit/a6fa40e1772dd29d919dc62be071a95b871facbd))
* **vm:** resize image once imported ([#753](https://github.com/bpg/terraform-provider-proxmox/issues/753)) ([d16b8e1](https://github.com/bpg/terraform-provider-proxmox/commit/d16b8e1696a50d5aca6bcce02a220069bfed0e87))
* **vm:** unable to clone as non-root due to `hook_script` ([#756](https://github.com/bpg/terraform-provider-proxmox/issues/756)) ([728eceb](https://github.com/bpg/terraform-provider-proxmox/commit/728eceb5e9fc342984218a0baf0155afae53fd71))


### Miscellaneous

* **deps:** bump github.com/brianvoe/gofakeit/v6 from 6.24.0 to 6.25.0 ([#741](https://github.com/bpg/terraform-provider-proxmox/issues/741)) ([9016641](https://github.com/bpg/terraform-provider-proxmox/commit/9016641c34d839ee1b5eb34892171897fac75880))
* **deps:** bump golang.org/x/crypto from 0.15.0 to 0.16.0 ([#752](https://github.com/bpg/terraform-provider-proxmox/issues/752)) ([a4ac84a](https://github.com/bpg/terraform-provider-proxmox/commit/a4ac84a78cca4b3f972d38e856adb294bb60fb2b))

## [0.38.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.38.0...v0.38.1) (2023-11-17)


### Bug Fixes

* **vm:** type error when unmarshalling `GetResponseData.data.memory` ([#728](https://github.com/bpg/terraform-provider-proxmox/issues/728)) ([b429f95](https://github.com/bpg/terraform-provider-proxmox/commit/b429f95ca578c530d08caae95228f20e57de0c03))

## [0.38.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.37.1...v0.38.0) (2023-11-17)


### Features

* **file:** rename content type `backup` -&gt; `dump` for backups ([#724](https://github.com/bpg/terraform-provider-proxmox/issues/724)) ([3280370](https://github.com/bpg/terraform-provider-proxmox/commit/3280370155ef339e5cf05ac4c94eb0e412d81d5c))
* **vm:** Add Win 11 as os type ([#720](https://github.com/bpg/terraform-provider-proxmox/issues/720)) ([0eeb7a7](https://github.com/bpg/terraform-provider-proxmox/commit/0eeb7a7fd924f4cd09e424219c532e55cc3ea721))


### Bug Fixes

* **vm:** memory size datatype conversion causing `null` on read ([#715](https://github.com/bpg/terraform-provider-proxmox/issues/715)) ([2bbf228](https://github.com/bpg/terraform-provider-proxmox/commit/2bbf228eecd4f34120f38b32102688d4b78eb220))
* **vm:** use int64 for resource memory and disk size ([#694](https://github.com/bpg/terraform-provider-proxmox/issues/694)) ([5fe6892](https://github.com/bpg/terraform-provider-proxmox/commit/5fe6892724a74906400d30a67c3047e1a0e86781))


### Miscellaneous

* **deps:** bump github.com/avast/retry-go/v4 from 4.5.0 to 4.5.1 ([#722](https://github.com/bpg/terraform-provider-proxmox/issues/722)) ([b0fea6d](https://github.com/bpg/terraform-provider-proxmox/commit/b0fea6d681301501826b95ae60e3f701e7bf79c2))
* **deps:** bump github.com/hashicorp/terraform-plugin-go from 0.19.0 to 0.19.1 ([#723](https://github.com/bpg/terraform-provider-proxmox/issues/723)) ([6c83e07](https://github.com/bpg/terraform-provider-proxmox/commit/6c83e07bdf65367f23326f8bd9f6f35cf254509f))
* **deps:** bump google-github-actions/release-please-action from 3.7.12 to 3.7.13 ([#716](https://github.com/bpg/terraform-provider-proxmox/issues/716)) ([4898d4d](https://github.com/bpg/terraform-provider-proxmox/commit/4898d4d80c352b47455431f1408366e4504ac524))

## [0.37.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.37.0...v0.37.1) (2023-11-12)


### Bug Fixes

* **docs:** add SSH info box to ressources needing it ([#690](https://github.com/bpg/terraform-provider-proxmox/issues/690)) ([e45c1c8](https://github.com/bpg/terraform-provider-proxmox/commit/e45c1c81263c723d1665c26c36d57f8c570b6ca3))
* **file:** display warning if directory is not found ([#703](https://github.com/bpg/terraform-provider-proxmox/issues/703)) ([e10b4b5](https://github.com/bpg/terraform-provider-proxmox/commit/e10b4b561793fd462f18ff1aa616b62ccfe586f2))
* **provider:** do not blindly use first IP for SSH ([#704](https://github.com/bpg/terraform-provider-proxmox/issues/704)) ([a586d03](https://github.com/bpg/terraform-provider-proxmox/commit/a586d0381e9c892b4b9aa2a0699f6c039c151ad2))
* **provider:** sanitize PVE endpoint value ([#686](https://github.com/bpg/terraform-provider-proxmox/issues/686)) ([3f582d8](https://github.com/bpg/terraform-provider-proxmox/commit/3f582d816334d4db370e8a5124f27ae4842c93f1))
* **storage:** unmarshal error when list storage containing large files ([#688](https://github.com/bpg/terraform-provider-proxmox/issues/688)) ([64c67d9](https://github.com/bpg/terraform-provider-proxmox/commit/64c67d947362e2653feaab7fa9ffb3b6016d0650))
* **vm:** update validation and docs for `machine` attribute ([#681](https://github.com/bpg/terraform-provider-proxmox/issues/681)) ([3fd6b6b](https://github.com/bpg/terraform-provider-proxmox/commit/3fd6b6b2ce36fa4bead31fa737f1137cd43cc16e))


### Miscellaneous

* **build:** add devcontainer ([#699](https://github.com/bpg/terraform-provider-proxmox/issues/699)) ([5bf9d1b](https://github.com/bpg/terraform-provider-proxmox/commit/5bf9d1b9da7359d3cc38ac123cfeb0629f215eca))
* **deps:** bump github.com/golangci/golangci-lint from 1.55.1 to 1.55.2 in /tools ([#680](https://github.com/bpg/terraform-provider-proxmox/issues/680)) ([2b8fd1a](https://github.com/bpg/terraform-provider-proxmox/commit/2b8fd1ad48540ad4552ab54b28e1d12379703c77))
* **deps:** bump github.com/goreleaser/goreleaser from 1.21.2 to 1.22.1 in /tools ([#709](https://github.com/bpg/terraform-provider-proxmox/issues/709)) ([48c89ff](https://github.com/bpg/terraform-provider-proxmox/commit/48c89ffe1327fe29ff8ccb22d298153d44d7669a))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.29.0 to 2.30.0 ([#708](https://github.com/bpg/terraform-provider-proxmox/issues/708)) ([817e43c](https://github.com/bpg/terraform-provider-proxmox/commit/817e43c912134be89c389cce1e718922ab993dde))
* **deps:** bump github.com/sigstore/cosign/v2 from 2.0.3-0.20230523133326-0544abd8fc8a to 2.2.1 in /tools ([#705](https://github.com/bpg/terraform-provider-proxmox/issues/705)) ([3e6fe4d](https://github.com/bpg/terraform-provider-proxmox/commit/3e6fe4db5598bfe475a0a844c3bd5937bc83aec3))
* **docs:** update hostpci id to mentions requirement around root user ([#710](https://github.com/bpg/terraform-provider-proxmox/issues/710)) ([0bf3a2a](https://github.com/bpg/terraform-provider-proxmox/commit/0bf3a2aea3d8d5e2821a24b3613bc44cd60b7b2d))
* **docs:** update VM ip address to mention the CIDR suffic requirement ([#697](https://github.com/bpg/terraform-provider-proxmox/issues/697)) ([d61cdc2](https://github.com/bpg/terraform-provider-proxmox/commit/d61cdc2b5c7efa50cc8228d5cebf789cc3f1cb5e))

## [0.37.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.36.0...v0.37.0) (2023-10-31)


### Features

* **vm:** add support for USB devices passthrough ([#666](https://github.com/bpg/terraform-provider-proxmox/issues/666)) ([cec4e65](https://github.com/bpg/terraform-provider-proxmox/commit/cec4e6586834feb876321520b93caf7ce4cb68d7))


### Bug Fixes

* **docs:** document qemu-guest-agent behavior ([#670](https://github.com/bpg/terraform-provider-proxmox/issues/670)) ([e2e5b4e](https://github.com/bpg/terraform-provider-proxmox/commit/e2e5b4e3441f46fbaef36751c6e5a6d1bc5ad671))
* **docs:** update `README.md` and file resource documentation ([#659](https://github.com/bpg/terraform-provider-proxmox/issues/659)) ([f6f05a5](https://github.com/bpg/terraform-provider-proxmox/commit/f6f05a56e4c9296491044ec8d5d0215a44da6f56))
* **vm:** MAC address validator should allow lowercase hex ([#660](https://github.com/bpg/terraform-provider-proxmox/issues/660)) ([7867e66](https://github.com/bpg/terraform-provider-proxmox/commit/7867e66d531484815b529fcdf0b8607fa837dc89))


### Miscellaneous

* **deps:** bump github.com/docker/docker from 24.0.2+incompatible to 24.0.7+incompatible in /tools ([#667](https://github.com/bpg/terraform-provider-proxmox/issues/667)) ([aea4a6f](https://github.com/bpg/terraform-provider-proxmox/commit/aea4a6f1cb0848a3274799dc605446c96dc192df))
* **deps:** bump github.com/golangci/golangci-lint from 1.55.0 to 1.55.1 in /tools ([#664](https://github.com/bpg/terraform-provider-proxmox/issues/664)) ([6ab1d5f](https://github.com/bpg/terraform-provider-proxmox/commit/6ab1d5fffbe265291e2991db162fb68fd1b50b02))
* **deps:** bump github.com/google/uuid from 1.3.1 to 1.4.0 ([#662](https://github.com/bpg/terraform-provider-proxmox/issues/662)) ([0ec8c24](https://github.com/bpg/terraform-provider-proxmox/commit/0ec8c2498b56a7f3d206409590f943e6764f8586))
* **deps:** bump github.com/hashicorp/terraform-plugin-framework from 1.4.1 to 1.4.2 ([#663](https://github.com/bpg/terraform-provider-proxmox/issues/663)) ([b389080](https://github.com/bpg/terraform-provider-proxmox/commit/b38908063f523fc5738a9e1c987b848b80ecb5d3))

## [0.36.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.35.1...v0.36.0) (2023-10-26)


### Features

* **vm:** add configurable timeout for VM create operation ([#648](https://github.com/bpg/terraform-provider-proxmox/issues/648)) ([a30f96c](https://github.com/bpg/terraform-provider-proxmox/commit/a30f96c348888522ca9278d8fef4bd9b12b1b634))


### Bug Fixes

* **file:** handle missing file on state refresh ([#649](https://github.com/bpg/terraform-provider-proxmox/issues/649)) ([2a56c23](https://github.com/bpg/terraform-provider-proxmox/commit/2a56c23f52abda293f328196a0d80b9becd749a7))
* **vm:** better handle of ctrl+c when qemu is not responding  ([#627](https://github.com/bpg/terraform-provider-proxmox/issues/627)) ([aec09e4](https://github.com/bpg/terraform-provider-proxmox/commit/aec09e4ecd8f9df937a04845162a679098f0c480))


### Miscellaneous

* **deps:** bump github.com/brianvoe/gofakeit/v6 from 6.23.2 to 6.24.0 ([#642](https://github.com/bpg/terraform-provider-proxmox/issues/642)) ([72951dc](https://github.com/bpg/terraform-provider-proxmox/commit/72951dc65691bdb44bae5f31b218a31811ffdfb7))
* **deps:** bump google.golang.org/grpc from 1.57.0 to 1.57.1 ([#652](https://github.com/bpg/terraform-provider-proxmox/issues/652)) ([4740da0](https://github.com/bpg/terraform-provider-proxmox/commit/4740da0d1f413743c252be963f7b6252ed3f0d96))
* **deps:** bump google.golang.org/grpc from 1.57.0 to 1.57.1 in /tools ([#653](https://github.com/bpg/terraform-provider-proxmox/issues/653)) ([db9140d](https://github.com/bpg/terraform-provider-proxmox/commit/db9140d05ef71753c8a8e1310c259717fef6e417))
* fix linter error ([#645](https://github.com/bpg/terraform-provider-proxmox/issues/645)) ([1056180](https://github.com/bpg/terraform-provider-proxmox/commit/1056180ca571ef171870be5e864461fb49732bdf))

## [0.35.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.35.0...v0.35.1) (2023-10-22)


### Bug Fixes

* **vm:** better check for disk ownership ([#633](https://github.com/bpg/terraform-provider-proxmox/issues/633)) ([6753582](https://github.com/bpg/terraform-provider-proxmox/commit/6753582e4b1999fdf2fd9ea0f499c0cd0f7cd64c))
* **vm:** set FileVolume for disks with file_id ([#635](https://github.com/bpg/terraform-provider-proxmox/issues/635)) ([d1d7bd3](https://github.com/bpg/terraform-provider-proxmox/commit/d1d7bd39c741d99b2395ef858bf739cb067f0542))


### Miscellaneous

* **deps:** bump github.com/golangci/golangci-lint from 1.54.2 to 1.55.0 in /tools ([#636](https://github.com/bpg/terraform-provider-proxmox/issues/636)) ([bcd33bb](https://github.com/bpg/terraform-provider-proxmox/commit/bcd33bb139d20ea4986d7dadf145e6ebbe497e79))

## [0.35.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.34.0...v0.35.0) (2023-10-17)


### Features

* **vm:** add 'path_in_datastore' disk argument ([#606](https://github.com/bpg/terraform-provider-proxmox/issues/606)) ([aeb5e88](https://github.com/bpg/terraform-provider-proxmox/commit/aeb5e88bc9112686675c7058501fa9378b69af93))


### Bug Fixes

* **lxc:** unmarshal string/int vmid as int when read container status ([#622](https://github.com/bpg/terraform-provider-proxmox/issues/622)) ([b90445a](https://github.com/bpg/terraform-provider-proxmox/commit/b90445a12c31c970c1cd1d2f37508ffcee586bf8))
* **provider:** add informative error around ssh-agent  ([#620](https://github.com/bpg/terraform-provider-proxmox/issues/620)) ([388ce7c](https://github.com/bpg/terraform-provider-proxmox/commit/388ce7ce8d37964da427d2430c9e03b14f790856))


### Miscellaneous

* **deps:** bump github.com/google/go-cmp from 0.5.9 to 0.6.0 ([#624](https://github.com/bpg/terraform-provider-proxmox/issues/624)) ([21e48c7](https://github.com/bpg/terraform-provider-proxmox/commit/21e48c7fb8aef8b5f5a48fea76ca9a030ccd59cc))
* **deps:** bump golang.org/x/net from 0.13.0 to 0.17.0 ([#616](https://github.com/bpg/terraform-provider-proxmox/issues/616)) ([29894bd](https://github.com/bpg/terraform-provider-proxmox/commit/29894bda234baca2645fc5e0d5d6f05101406b18))
* **deps:** bump golang.org/x/net from 0.15.0 to 0.17.0 in /tools ([#617](https://github.com/bpg/terraform-provider-proxmox/issues/617)) ([7287f5d](https://github.com/bpg/terraform-provider-proxmox/commit/7287f5de4801d0f56faa8e3b99c80d41ac2f1f01))

## [0.34.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.33.0...v0.34.0) (2023-10-10)


### Features

* **lxc:** add support for the `start_on_boot` option ([#605](https://github.com/bpg/terraform-provider-proxmox/issues/605)) ([d36cf4e](https://github.com/bpg/terraform-provider-proxmox/commit/d36cf4eab81955184c926c86ce692bcf6c01b840))
* **provider:** configure temp directory ([#607](https://github.com/bpg/terraform-provider-proxmox/issues/607)) ([06ad004](https://github.com/bpg/terraform-provider-proxmox/commit/06ad00463c8ec0426f72a559924e6a0adfe4e2a8))
* **vm:** add option to enable multiqueue in network devices ([#614](https://github.com/bpg/terraform-provider-proxmox/issues/614)) ([be5251d](https://github.com/bpg/terraform-provider-proxmox/commit/be5251dd5ad535be6bdf8f9ef73c43f54a9dc2c7))


### Bug Fixes

* **lxc:** cloned container does not start by default ([#615](https://github.com/bpg/terraform-provider-proxmox/issues/615)) ([d5994a2](https://github.com/bpg/terraform-provider-proxmox/commit/d5994a2bd5323cef34b71f3fea895539a0cfccd8))
* **lxc:** create container when authenticated with API token ([#610](https://github.com/bpg/terraform-provider-proxmox/issues/610)) ([32bdc94](https://github.com/bpg/terraform-provider-proxmox/commit/32bdc94167253b7b3ec6eaecbccc2d2cc0104b61))
* **lxc:** multi-line description always shows as changed ([#611](https://github.com/bpg/terraform-provider-proxmox/issues/611)) ([088ad09](https://github.com/bpg/terraform-provider-proxmox/commit/088ad09e356e1baf17b7cb84656155d192d2909d))


### Miscellaneous

* **deps:** bump github.com/hashicorp/terraform-plugin-framework from 1.4.0 to 1.4.1 ([#612](https://github.com/bpg/terraform-provider-proxmox/issues/612)) ([a266496](https://github.com/bpg/terraform-provider-proxmox/commit/a266496fcbf9c044712896ea1af5827f47869be1))
* **deps:** bump golang.org/x/crypto from 0.13.0 to 0.14.0 ([#613](https://github.com/bpg/terraform-provider-proxmox/issues/613)) ([0150a97](https://github.com/bpg/terraform-provider-proxmox/commit/0150a97cd4a2489311db943459f7d41b8ef8e61e))

## [0.33.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.32.2...v0.33.0) (2023-10-02)


### Features

* **file:** add optional `overwrite` flag to the file resource ([#593](https://github.com/bpg/terraform-provider-proxmox/issues/593)) ([5e24a75](https://github.com/bpg/terraform-provider-proxmox/commit/5e24a75d09b930aef07a067b37be0507c1948de1))
* **vm:** allow `scsi` and `sata` interfaces for CloudInit Drive ([#598](https://github.com/bpg/terraform-provider-proxmox/issues/598)) ([0b8f2e2](https://github.com/bpg/terraform-provider-proxmox/commit/0b8f2e2c6f80b0370290e6b32ba1e7add977018c))


### Bug Fixes

* **api:** set min TLS version 1.3, secure HTTP-only cookie ([#596](https://github.com/bpg/terraform-provider-proxmox/issues/596)) ([16ebf30](https://github.com/bpg/terraform-provider-proxmox/commit/16ebf30a79e8e3cc2df48787b210fd78950f8260))


### Miscellaneous

* **ci:** cleanup CI flows ([#595](https://github.com/bpg/terraform-provider-proxmox/issues/595)) ([bd09fd3](https://github.com/bpg/terraform-provider-proxmox/commit/bd09fd3d6ec954e6d2c8d01e51050faf5677d422))

## [0.32.2](https://github.com/bpg/terraform-provider-proxmox/compare/v0.32.1...v0.32.2) (2023-09-28)


### Bug Fixes

* **tasks:** fix UPID (task id) parsing error ([#591](https://github.com/bpg/terraform-provider-proxmox/issues/591)) ([294a9da](https://github.com/bpg/terraform-provider-proxmox/commit/294a9daa8711f7a2dbb054f1de750bf9f1bb4f3a))


### Miscellaneous

* **deps:** bump github.com/goreleaser/goreleaser from 1.20.0 to 1.21.0 in /tools ([#587](https://github.com/bpg/terraform-provider-proxmox/issues/587)) ([2573323](https://github.com/bpg/terraform-provider-proxmox/commit/257332393f48dc2c5367c8934923bea28964ffdc))
* **deps:** bump github.com/goreleaser/goreleaser from 1.21.0 to 1.21.2 in /tools ([#592](https://github.com/bpg/terraform-provider-proxmox/issues/592)) ([2621aad](https://github.com/bpg/terraform-provider-proxmox/commit/2621aadb5f089a88b6ddf027ce906c20031ee2a0))

## [0.32.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.32.0...v0.32.1) (2023-09-23)


### Bug Fixes

* **cluster:** inconsistencies in applying cluster options ([#573](https://github.com/bpg/terraform-provider-proxmox/issues/573)) ([03f3ed7](https://github.com/bpg/terraform-provider-proxmox/commit/03f3ed7871e2a2fe653d6cfe9dcb28196738e1e2))
* **network:** remove computed flag from mtu attribute ([#572](https://github.com/bpg/terraform-provider-proxmox/issues/572)) ([5720fe4](https://github.com/bpg/terraform-provider-proxmox/commit/5720fe4673873166e7dbf7bc687b57837b99b117))


### Miscellaneous

* **code:** bump go to v1.21 ([#585](https://github.com/bpg/terraform-provider-proxmox/issues/585)) ([11c0940](https://github.com/bpg/terraform-provider-proxmox/commit/11c09405ea2f6d9dfc28191ce50739f811b5f0c4))
* **code:** re-organize and cleanup "fwk provider"'s code ([#568](https://github.com/bpg/terraform-provider-proxmox/issues/568)) ([7d064a8](https://github.com/bpg/terraform-provider-proxmox/commit/7d064a8b27d78a1564c9da914f17340966d955d1))
* **deps:** bump github.com/skeema/knownhosts from 1.2.0 to 1.2.1 ([#584](https://github.com/bpg/terraform-provider-proxmox/issues/584)) ([7890212](https://github.com/bpg/terraform-provider-proxmox/commit/7890212a566036bf448f4db149a7f8816de187ab))
* **docs:** add "Proof of Work" section to the PR template ([#583](https://github.com/bpg/terraform-provider-proxmox/issues/583)) ([de1eb2b](https://github.com/bpg/terraform-provider-proxmox/commit/de1eb2b950ae6a001ad07a93f27a90858500749b))
* **docs:** add a note about DCO to CONTRIBUTING.md ([#574](https://github.com/bpg/terraform-provider-proxmox/issues/574)) ([d0c9b45](https://github.com/bpg/terraform-provider-proxmox/commit/d0c9b4594d46c327b5a0f09288ac5b88a48af61e))
* **docs:** update `proxmox_virtual_environment_file` documentation ([#580](https://github.com/bpg/terraform-provider-proxmox/issues/580)) ([7dde53c](https://github.com/bpg/terraform-provider-proxmox/commit/7dde53cf1cee1127cecb86cab3b6e75331410c56))

## [0.32.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.31.0...v0.32.0) (2023-09-13)


### Features

* **cluster:** add cluster options resource ([#548](https://github.com/bpg/terraform-provider-proxmox/issues/548)) ([de8b4ec](https://github.com/bpg/terraform-provider-proxmox/commit/de8b4ec41ada527b5a14883b5dcacdab2684fc37))


### Bug Fixes

* **lxc,vm:** error unmarshalling string `cpulimit` ([#563](https://github.com/bpg/terraform-provider-proxmox/issues/563)) ([11a8ec0](https://github.com/bpg/terraform-provider-proxmox/commit/11a8ec0c9594c1b9ff305edcd47f090309bc1466))


### Miscellaneous

* **ci:** cleanup and update project configs ([#549](https://github.com/bpg/terraform-provider-proxmox/issues/549)) ([edec5bf](https://github.com/bpg/terraform-provider-proxmox/commit/edec5bfd1cc25886fa36e1344a6de4a6d2427786))
* **code:** remove redundant `types2` import aliases ([#564](https://github.com/bpg/terraform-provider-proxmox/issues/564)) ([2dee65b](https://github.com/bpg/terraform-provider-proxmox/commit/2dee65bd0b872b795f559530cbd5b12c856e5771))
* **deps:** bump crazy-max/ghaction-import-gpg from 5 to 6 ([#558](https://github.com/bpg/terraform-provider-proxmox/issues/558)) ([1f8330a](https://github.com/bpg/terraform-provider-proxmox/commit/1f8330afc7f189964ab09fa652b39e2123e6187e))
* **deps:** bump github.com/hashicorp/terraform-plugin-* dependencies ([#561](https://github.com/bpg/terraform-provider-proxmox/issues/561)) ([3d7fbaa](https://github.com/bpg/terraform-provider-proxmox/commit/3d7fbaa7c7f8ce7a1cedf1dae3d31fceecad5ea1))
* **deps:** bump golang.org/x/crypto from 0.12.0 to 0.13.0 ([#554](https://github.com/bpg/terraform-provider-proxmox/issues/554)) ([1040aab](https://github.com/bpg/terraform-provider-proxmox/commit/1040aabb23d6eb7ff3841315aa5f608b24437e26))
* **deps:** bump goreleaser/goreleaser-action from 4.4.0 to 5.0.0 ([#560](https://github.com/bpg/terraform-provider-proxmox/issues/560)) ([ac556b5](https://github.com/bpg/terraform-provider-proxmox/commit/ac556b55150d271c098916b2134d3991f765891a))

## [0.31.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.30.3...v0.31.0) (2023-09-04)


### Features

* **file:** FORMAT CHANGE: update import id, so it matches the resource's format: `&lt;node_name&gt;/<datastore_id>:<content_type>/<file>` ([#543](https://github.com/bpg/terraform-provider-proxmox/issues/543)) ([7ace07d](https://github.com/bpg/terraform-provider-proxmox/commit/7ace07dfa47c4a6750973d04cb8d853fc9640047))
* **lxc:** add support for `keyctl` and `fuse` features ([#537](https://github.com/bpg/terraform-provider-proxmox/issues/537)) ([8ce9006](https://github.com/bpg/terraform-provider-proxmox/commit/8ce9006eed15dadc6f051464b8b98e3a1abd7d6d))
* **provider:** add optional SSH port param to node in provider ssh block ([#520](https://github.com/bpg/terraform-provider-proxmox/issues/520)) ([124cac2](https://github.com/bpg/terraform-provider-proxmox/commit/124cac247ce34e2603b0d1c1c94106d958185708))


### Bug Fixes

* **provider:** panic crash in provider, interface conversion error ([#545](https://github.com/bpg/terraform-provider-proxmox/issues/545)) ([13326bb](https://github.com/bpg/terraform-provider-proxmox/commit/13326bbd33648391f0f87d339db272145e3066ac))
* **vm:** explicitly allow `""` as a value for CloudInit interfaces ([#546](https://github.com/bpg/terraform-provider-proxmox/issues/546)) ([0233053](https://github.com/bpg/terraform-provider-proxmox/commit/0233053dd8f8aa0fbfae8f7c11bb8ce359576bce))


### Miscellaneous

* **code:** fix `proxmox` package dependencies ([#536](https://github.com/bpg/terraform-provider-proxmox/issues/536)) ([5ecf135](https://github.com/bpg/terraform-provider-proxmox/commit/5ecf13539862bb9602696a7575568f228fc85e29))
* **deps:** bump actions/checkout from 3 to 4 ([#541](https://github.com/bpg/terraform-provider-proxmox/issues/541)) ([44d6d6b](https://github.com/bpg/terraform-provider-proxmox/commit/44d6d6b080c534ad16b3d9911ae445d4e16acfa3))

## [0.30.3](https://github.com/bpg/terraform-provider-proxmox/compare/v0.30.2...v0.30.3) (2023-09-01)


### Bug Fixes

* **file:** file upload in multi-node PVE cluster ([#533](https://github.com/bpg/terraform-provider-proxmox/issues/533)) ([ef2f2c1](https://github.com/bpg/terraform-provider-proxmox/commit/ef2f2c115976dfd97de2ce557be899927672f4b8))

## [0.30.2](https://github.com/bpg/terraform-provider-proxmox/compare/v0.30.1...v0.30.2) (2023-08-31)


### Bug Fixes

* **core:** improve error handling while waiting for PVE tasks to complete ([#526](https://github.com/bpg/terraform-provider-proxmox/issues/526)) ([6f02df4](https://github.com/bpg/terraform-provider-proxmox/commit/6f02df4440566ed1d97e0c6d016311b91bd53125))
* **file:** forced replacement of file resources that missing `timeout_upload` attribute ([#528](https://github.com/bpg/terraform-provider-proxmox/issues/528)) ([11d8261](https://github.com/bpg/terraform-provider-proxmox/commit/11d82614e628d24d9ee8db5cccc33427bf5a811c))
* **node:** creating linux_bridge with 'vlan_aware=false' or 'autostart=false' ([#529](https://github.com/bpg/terraform-provider-proxmox/issues/529)) ([f00e48a](https://github.com/bpg/terraform-provider-proxmox/commit/f00e48a51e1618bccf1d1800590b81696db15071))
* **provider:** User-settable VLAN ID and name ([#518](https://github.com/bpg/terraform-provider-proxmox/issues/518)) ([5599c7a](https://github.com/bpg/terraform-provider-proxmox/commit/5599c7afe45dbea217457b1452186c02b07db90f))


### Miscellaneous

* **deps:** bump activesupport from 7.0.6 to 7.0.7.1 in /docs ([#522](https://github.com/bpg/terraform-provider-proxmox/issues/522)) ([cd7927b](https://github.com/bpg/terraform-provider-proxmox/commit/cd7927bec347f22ecce500147866fbe01e742b51))
* **deps:** bump github.com/hashicorp/terraform-plugin-framework-validators from 0.11.0 to 0.12.0 ([#530](https://github.com/bpg/terraform-provider-proxmox/issues/530)) ([e35443a](https://github.com/bpg/terraform-provider-proxmox/commit/e35443a23b9528290952c24db573971d115e9877))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.27.0 to 2.28.0 ([#524](https://github.com/bpg/terraform-provider-proxmox/issues/524)) ([5556b17](https://github.com/bpg/terraform-provider-proxmox/commit/5556b17a1ed1e4e92343d17d534461348d3da30c))

## [0.30.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.30.0...v0.30.1) (2023-08-22)


### Bug Fixes

* **vm:** fix PCI device resource mapping changed ([#517](https://github.com/bpg/terraform-provider-proxmox/issues/517)) ([b1ac87d](https://github.com/bpg/terraform-provider-proxmox/commit/b1ac87df1df96a9172fee7cb4aa5934c6afb4ef1))


### Miscellaneous

* **deps:** bump github.com/golangci/golangci-lint from 1.54.1 to 1.54.2 in /tools ([#514](https://github.com/bpg/terraform-provider-proxmox/issues/514)) ([731dad8](https://github.com/bpg/terraform-provider-proxmox/commit/731dad87945335ebd3f897ff747edfc3e30607c4))
* **deps:** bump github.com/google/uuid from 1.3.0 to 1.3.1 ([#515](https://github.com/bpg/terraform-provider-proxmox/issues/515)) ([79c7f10](https://github.com/bpg/terraform-provider-proxmox/commit/79c7f100f6cfd2ea52d50aa69b92f5c99a0deded))

## [0.30.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.29.0...v0.30.0) (2023-08-21)


### Features

* **ha:** add support for Proxmox High Availability objects ([#498](https://github.com/bpg/terraform-provider-proxmox/issues/498)) ([03c9b36](https://github.com/bpg/terraform-provider-proxmox/commit/03c9b36b86914583c1709e99db305682b7b7dc99))
* **vm:** add support for migration when the node name is modified ([#501](https://github.com/bpg/terraform-provider-proxmox/issues/501)) ([a285360](https://github.com/bpg/terraform-provider-proxmox/commit/a2853606ad294476e9b5f17a994cb230643e9277))
* **vm:** add support for non-default CloudInit interface and CloudInit storage change ([#486](https://github.com/bpg/terraform-provider-proxmox/issues/486)) ([5475936](https://github.com/bpg/terraform-provider-proxmox/commit/547593661f5bcab1141edc9a7203dca65c6b539d))
* **vm:** add support for pool update ([#505](https://github.com/bpg/terraform-provider-proxmox/issues/505)) ([e6c15ec](https://github.com/bpg/terraform-provider-proxmox/commit/e6c15eccc6fd2076afb2f521e28f27976abba892))
* **vm:** fix adding/removing hostpci devices forcing vm recreation ([#504](https://github.com/bpg/terraform-provider-proxmox/issues/504)) ([a038fd3](https://github.com/bpg/terraform-provider-proxmox/commit/a038fd31420fe23963c7d68198ed5f40b6583058))
* **vm:** support PCI device resource mapping ([#500](https://github.com/bpg/terraform-provider-proxmox/issues/500)) ([2697054](https://github.com/bpg/terraform-provider-proxmox/commit/26970541c48495b7b9fd220960c83f54956e8132))


### Bug Fixes

* **vm:** fix CloudInit datastore change support ([#509](https://github.com/bpg/terraform-provider-proxmox/issues/509)) ([73c1294](https://github.com/bpg/terraform-provider-proxmox/commit/73c1294979b956939b755ac05796fb1a68f92f75))
* **vm:** fix index out of range when unmarshalling custompcidevice ([#496](https://github.com/bpg/terraform-provider-proxmox/issues/496)) ([78d6683](https://github.com/bpg/terraform-provider-proxmox/commit/78d668377f383badd8a53a18dbd4cb65e67176c2))
* **vm:** fixed startup / shutdown behaviour on HA clusters ([#508](https://github.com/bpg/terraform-provider-proxmox/issues/508)) ([148a9e0](https://github.com/bpg/terraform-provider-proxmox/commit/148a9e0c9c3f8d78645846b39646ad7d8c78c4a5))
* **vm:** no IP address detection when VM contains bridges ([#493](https://github.com/bpg/terraform-provider-proxmox/issues/493)) ([9fd9d21](https://github.com/bpg/terraform-provider-proxmox/commit/9fd9d211d75e760ef1c7e44d13de9ce8d38bf834))


### Miscellaneous

* **deps:** bump github.com/golangci/golangci-lint from 1.54.0 to 1.54.1 in /tools ([#489](https://github.com/bpg/terraform-provider-proxmox/issues/489)) ([e4f9888](https://github.com/bpg/terraform-provider-proxmox/commit/e4f9888f6f6db835d425e52878631517ef4d5e14))
* **deps:** bump github.com/hashicorp/terraform-plugin-framework from 1.3.4 to 1.3.5 ([#512](https://github.com/bpg/terraform-provider-proxmox/issues/512)) ([98ae6a8](https://github.com/bpg/terraform-provider-proxmox/commit/98ae6a8d8f489b98c05d88598594d43c004b6316))
* **deps:** bump github.com/pkg/sftp from 1.13.5 to 1.13.6 ([#488](https://github.com/bpg/terraform-provider-proxmox/issues/488)) ([9045183](https://github.com/bpg/terraform-provider-proxmox/commit/9045183c1dd18ba2ebd3c1afcd7c16e73213bf27))
* **vm:** fix linter errors ([#506](https://github.com/bpg/terraform-provider-proxmox/issues/506)) ([1896ea0](https://github.com/bpg/terraform-provider-proxmox/commit/1896ea08f09ec4e684a886d11a5915c6e573eac1))

## [0.29.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.28.0...v0.29.0) (2023-08-10)


### Features

* **file:** ensure upload of ISO/VSTMPL is completed upon resource creation ([#471](https://github.com/bpg/terraform-provider-proxmox/issues/471)) ([f901e71](https://github.com/bpg/terraform-provider-proxmox/commit/f901e711dd4e8cd59b3a1e34c58a1a03564bd13a))


### Bug Fixes

* **user:** make `password` attribute optional ([#474](https://github.com/bpg/terraform-provider-proxmox/issues/474)) ([244e061](https://github.com/bpg/terraform-provider-proxmox/commit/244e061779f05752bd0760ea6b5a15c869e26505))
* **vm:** default disk cache is not set to `none` if not specified for an existing disk ([#478](https://github.com/bpg/terraform-provider-proxmox/issues/478)) ([8d0b3ed](https://github.com/bpg/terraform-provider-proxmox/commit/8d0b3ed25fa1c2dcc0d319d725aea34f3e18aef8))
* **vm:** ensure startup / shutdown delay is applied when order is not configured ([#479](https://github.com/bpg/terraform-provider-proxmox/issues/479)) ([2cf64b8](https://github.com/bpg/terraform-provider-proxmox/commit/2cf64b88c35991db19f83a1fa69ed41cbceebd32))


### Miscellaneous

* **deps-dev:** bump commonmarker from 0.23.9 to 0.23.10 in /docs ([#472](https://github.com/bpg/terraform-provider-proxmox/issues/472)) ([2e16fbb](https://github.com/bpg/terraform-provider-proxmox/commit/2e16fbb44bf45e58c9296d1e1e28d3fbea9d732c))
* **deps:** bump github.com/golangci/golangci-lint from 1.53.3 to 1.54.0 in /tools ([#482](https://github.com/bpg/terraform-provider-proxmox/issues/482)) ([390f03c](https://github.com/bpg/terraform-provider-proxmox/commit/390f03c1590725d7f89a1f38c3848269bbe4c402))
* **deps:** bump github.com/goreleaser/goreleaser from 1.19.2 to 1.20.0 in /tools ([#481](https://github.com/bpg/terraform-provider-proxmox/issues/481)) ([eb3d847](https://github.com/bpg/terraform-provider-proxmox/commit/eb3d8473acd593ef0e876b711c3ffc3441fc4b54))
* **deps:** bump github.com/hashicorp/terraform-plugin-framework from 1.3.3 to 1.3.4 ([#466](https://github.com/bpg/terraform-provider-proxmox/issues/466)) ([8a5a533](https://github.com/bpg/terraform-provider-proxmox/commit/8a5a53301b3e2e7ecad9322c80f2700726ea0504))
* **deps:** bump github.com/hashicorp/terraform-plugin-framework-validators from 0.10.0 to 0.11.0 ([#467](https://github.com/bpg/terraform-provider-proxmox/issues/467)) ([7c9e3ed](https://github.com/bpg/terraform-provider-proxmox/commit/7c9e3ed1afaa3bf50b78be1aefd77cc76fc3d06d))
* **deps:** bump golang.org/x/crypto from 0.11.0 to 0.12.0 ([#465](https://github.com/bpg/terraform-provider-proxmox/issues/465)) ([185e98f](https://github.com/bpg/terraform-provider-proxmox/commit/185e98fe802119ab0de53bb2eeb34d7510517475))
* **deps:** bump goreleaser/goreleaser-action from 4.3.0 to 4.4.0 ([#480](https://github.com/bpg/terraform-provider-proxmox/issues/480)) ([a7047da](https://github.com/bpg/terraform-provider-proxmox/commit/a7047dac7269143ce833da1310ca6f03646ccbf1))

## [0.28.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.27.0...v0.28.0) (2023-08-06)


### Features

* **vm:** add support for SMBIOS settings ([#454](https://github.com/bpg/terraform-provider-proxmox/issues/454)) ([85ff60d](https://github.com/bpg/terraform-provider-proxmox/commit/85ff60d4bd928880eebeb6bbd9440a65a3e2cc9d))


### Bug Fixes

* **api:** remove HTTP client timeout ([#464](https://github.com/bpg/terraform-provider-proxmox/issues/464)) ([824e51c](https://github.com/bpg/terraform-provider-proxmox/commit/824e51c6508fe0e5905b143ef6d8dd161b1acbfe))
* **user:** make `password` attribute optional ([#463](https://github.com/bpg/terraform-provider-proxmox/issues/463)) ([5a3b1cc](https://github.com/bpg/terraform-provider-proxmox/commit/5a3b1ccaf703db260ba25e564c04506ea0de6247))
* **vm:** give `cache` the correct default value ([#450](https://github.com/bpg/terraform-provider-proxmox/issues/450)) ([0d3227a](https://github.com/bpg/terraform-provider-proxmox/commit/0d3227a890b4df12ecb71fbd3215e5f1d4babff8))


### Miscellaneous

* **doc:** add all-contributors to README.md ([#455](https://github.com/bpg/terraform-provider-proxmox/issues/455)) ([d885e64](https://github.com/bpg/terraform-provider-proxmox/commit/d885e643728c1da30deca3f26150a57ba75593db))
* **doc:** add existing contributors ([#459](https://github.com/bpg/terraform-provider-proxmox/issues/459)) ([cb71d73](https://github.com/bpg/terraform-provider-proxmox/commit/cb71d731f1903ec9fbfa2eb5d4b78c53c961f86f))
* **doc:** cleanup readme ([#461](https://github.com/bpg/terraform-provider-proxmox/issues/461)) ([368b133](https://github.com/bpg/terraform-provider-proxmox/commit/368b133427e14753b287469b814591141126913d))

## [0.27.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.26.0...v0.27.0) (2023-07-30)


### Features

* **vm:** add support for disk `cache` option ([#443](https://github.com/bpg/terraform-provider-proxmox/issues/443)) ([cfe3d96](https://github.com/bpg/terraform-provider-proxmox/commit/cfe3d96576b521cb294f217fb3f7caf45347e58e))
* **vm:** add support for start/shutdown order configuration ([#445](https://github.com/bpg/terraform-provider-proxmox/issues/445)) ([b045746](https://github.com/bpg/terraform-provider-proxmox/commit/b045746a94d2717b69fc48234b9ece101b53bdcd))


### Bug Fixes

* **vm:** cloned VM with `efi_disk` got re-created at re-apply ([#447](https://github.com/bpg/terraform-provider-proxmox/issues/447)) ([c1e7cea](https://github.com/bpg/terraform-provider-proxmox/commit/c1e7cea21ed7d49375de8850f9cd3737d485c3d2))


### Miscellaneous

* update dependencies, cleanup docs ([#446](https://github.com/bpg/terraform-provider-proxmox/issues/446)) ([a3b95c8](https://github.com/bpg/terraform-provider-proxmox/commit/a3b95c80536c7b69ab5b4c10e434d410ec5e05e5))

## [0.26.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.25.0...v0.26.0) (2023-07-29)


### Features

* **core:** migrate `version` datasource to TF plugin framework ([#440](https://github.com/bpg/terraform-provider-proxmox/issues/440)) ([a9a7329](https://github.com/bpg/terraform-provider-proxmox/commit/a9a7329d9fef42466f6fe2a7eeff9645100459c6))


### Miscellaneous

* **deps:** bump github.com/hashicorp/terraform-plugin-framework from 1.3.2 to 1.3.3 ([#439](https://github.com/bpg/terraform-provider-proxmox/issues/439)) ([d82a08d](https://github.com/bpg/terraform-provider-proxmox/commit/d82a08dcb434e3b2aa0241332aeb3b43eac372d1))
* **docs:** Update README.md  ([#442](https://github.com/bpg/terraform-provider-proxmox/issues/442)) ([8e2d180](https://github.com/bpg/terraform-provider-proxmox/commit/8e2d18053f0fca807ecd81cbf2c4a3b5169f0d49))

## [0.25.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.24.2...v0.25.0) (2023-07-20)


### Features

* **lxc:** add support for lxc mount points ([#394](https://github.com/bpg/terraform-provider-proxmox/issues/394)) ([beef9b1](https://github.com/bpg/terraform-provider-proxmox/commit/beef9b1219dc078cc7a3adeae9e6162235c603f8))


### Bug Fixes

* **vm:** Don't add an extra efi_disk entry ([#435](https://github.com/bpg/terraform-provider-proxmox/issues/435)) ([6781c03](https://github.com/bpg/terraform-provider-proxmox/commit/6781c03ca1eb794ed9e5ab322e1b73d57969b721))
* **vm:** fix for the api call upon empty disks ([#436](https://github.com/bpg/terraform-provider-proxmox/issues/436)) ([aea9846](https://github.com/bpg/terraform-provider-proxmox/commit/aea9846c6f5399d721458f21e94f253922103432))


### Miscellaneous

* cleanup resource validators & utility code ([#438](https://github.com/bpg/terraform-provider-proxmox/issues/438)) ([b2a27f3](https://github.com/bpg/terraform-provider-proxmox/commit/b2a27f3ccfa6e318d2243bae2c855f47e5523240))
* **deps:** bump github.com/hashicorp/terraform-plugin-mux from 0.11.1 to 0.11.2 ([#432](https://github.com/bpg/terraform-provider-proxmox/issues/432)) ([4324b29](https://github.com/bpg/terraform-provider-proxmox/commit/4324b294239bca04de550027402deabe1e6f1615))
* **deps:** bump github.com/skeema/knownhosts from 1.1.1 to 1.2.0 ([#433](https://github.com/bpg/terraform-provider-proxmox/issues/433)) ([b9ee3ae](https://github.com/bpg/terraform-provider-proxmox/commit/b9ee3ae10d942b3700fa057553471c9ed47ce4d4))

## [0.24.2](https://github.com/bpg/terraform-provider-proxmox/compare/v0.24.1...v0.24.2) (2023-07-16)


### Bug Fixes

* **vm:** do not reboot VM on config change if it is not running ([#430](https://github.com/bpg/terraform-provider-proxmox/issues/430)) ([0281bc8](https://github.com/bpg/terraform-provider-proxmox/commit/0281bc83e2d64fdfe2782feb6f21395706dbcc32))

## [0.24.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.24.0...v0.24.1) (2023-07-16)


### Bug Fixes

* **firewall:** add VM / container ID validation to firewall rules ([#424](https://github.com/bpg/terraform-provider-proxmox/issues/424)) ([6a3bc03](https://github.com/bpg/terraform-provider-proxmox/commit/6a3bc034706cef4190651118bfc2e8f62de8aecd))
* **vm:** add `interface` argument to `cdrom` block ([#429](https://github.com/bpg/terraform-provider-proxmox/issues/429)) ([b86fa23](https://github.com/bpg/terraform-provider-proxmox/commit/b86fa239ddd29f0cfc60d66ac4cede39b0167985))
* **vm:** add missing unmarshal for vm custom startup order ([#428](https://github.com/bpg/terraform-provider-proxmox/issues/428)) ([e59b06e](https://github.com/bpg/terraform-provider-proxmox/commit/e59b06e5195da90da837f5b660e6b76cca9fd632))

## [0.24.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.23.0...v0.24.0) (2023-07-09)


### Features

* add import support for a lot of resources ([#390](https://github.com/bpg/terraform-provider-proxmox/issues/390)) ([4147ff6](https://github.com/bpg/terraform-provider-proxmox/commit/4147ff6a29500dd47cd905a0239abdc28cffc596))
* **vm:** add more valid cpu types ([#411](https://github.com/bpg/terraform-provider-proxmox/issues/411)) ([e9a9fd7](https://github.com/bpg/terraform-provider-proxmox/commit/e9a9fd76dae22be24767cdf44cb9668f96c9ea90))


### Bug Fixes

* **firewall:** ignore non-existent rules at read/delete ([#415](https://github.com/bpg/terraform-provider-proxmox/issues/415)) ([fc3bbc3](https://github.com/bpg/terraform-provider-proxmox/commit/fc3bbc3d92466fc069db69619b5f1a7f338fc391))
* **node:** fix error when listing network interfaces of a node ([#412](https://github.com/bpg/terraform-provider-proxmox/issues/412)) ([16ee6a9](https://github.com/bpg/terraform-provider-proxmox/commit/16ee6a9f955f0452b80ba4ee88667edd4bd34fde))
* **node:** ignore field `bridge_fd` when listing network interfaces of a node ([#414](https://github.com/bpg/terraform-provider-proxmox/issues/414)) ([01a8456](https://github.com/bpg/terraform-provider-proxmox/commit/01a845636ae7242ea78b52365468f496fc52372b))


### Miscellaneous

* **deps:** bump github.com/goreleaser/goreleaser from 1.18.2 to 1.19.1 in /tools ([#403](https://github.com/bpg/terraform-provider-proxmox/issues/403)) ([0597217](https://github.com/bpg/terraform-provider-proxmox/commit/059721741ac5508bb98a1ca50b83a67e6a86c206))
* **deps:** bump github.com/goreleaser/goreleaser from 1.19.1 to 1.19.2 in /tools ([#417](https://github.com/bpg/terraform-provider-proxmox/issues/417)) ([7240715](https://github.com/bpg/terraform-provider-proxmox/commit/72407157614179f4368698235f957ead68dd51b1))
* **deps:** bump github.com/hashicorp/terraform-plugin-docs from 0.15.0 to 0.16.0 in /tools ([#418](https://github.com/bpg/terraform-provider-proxmox/issues/418)) ([6a309ac](https://github.com/bpg/terraform-provider-proxmox/commit/6a309ac4abec72e54529a84e17a360763110dfaa))
* **deps:** bump github.com/hashicorp/terraform-plugin-framework from 1.3.1 to 1.3.2 ([#401](https://github.com/bpg/terraform-provider-proxmox/issues/401)) ([908713a](https://github.com/bpg/terraform-provider-proxmox/commit/908713a08493e46796438c7bc5585efab25fc4e0))
* **deps:** bump github.com/hashicorp/terraform-plugin-go from 0.16.0 to 0.17.0 ([#399](https://github.com/bpg/terraform-provider-proxmox/issues/399)) ([24ee318](https://github.com/bpg/terraform-provider-proxmox/commit/24ee318cc33a0faad76045644ac03394a13c7605))
* **deps:** bump github.com/hashicorp/terraform-plugin-go from 0.17.0 to 0.18.0 ([#408](https://github.com/bpg/terraform-provider-proxmox/issues/408)) ([f494525](https://github.com/bpg/terraform-provider-proxmox/commit/f49452543c6c88f90fb46d245b39ea9942eca5ea))
* **deps:** bump github.com/hashicorp/terraform-plugin-mux from 0.10.0 to 0.11.1 ([#400](https://github.com/bpg/terraform-provider-proxmox/issues/400)) ([1a6cfb2](https://github.com/bpg/terraform-provider-proxmox/commit/1a6cfb2cf1039694e594b45cb79ac5bba7810383))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.26.1 to 2.27.0 ([#402](https://github.com/bpg/terraform-provider-proxmox/issues/402)) ([af56c4b](https://github.com/bpg/terraform-provider-proxmox/commit/af56c4b2a75b6611e6dbcddc755a00ccddfd5248))
* **deps:** bump golang.org/x/crypto from 0.10.0 to 0.11.0 ([#416](https://github.com/bpg/terraform-provider-proxmox/issues/416)) ([5e173e0](https://github.com/bpg/terraform-provider-proxmox/commit/5e173e0bc9d2d7219e385e8b64ae82b3fcfdb42f))
* **refactoring:** remove accidentally added `types2` import alias ([#409](https://github.com/bpg/terraform-provider-proxmox/issues/409)) ([feac6b0](https://github.com/bpg/terraform-provider-proxmox/commit/feac6b0128520a16c4ecd4850d4f73e311ec1f7b))

## [0.23.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.22.0...v0.23.0) (2023-07-03)


### Features

* **vm:** efi disk, cpu numa ([#384](https://github.com/bpg/terraform-provider-proxmox/issues/384)) ([e9a74e9](https://github.com/bpg/terraform-provider-proxmox/commit/e9a74e90374570dee4af93fed4454209157bcbb7))


### Bug Fixes

* **docs:** minor firewall options page improvement ([#396](https://github.com/bpg/terraform-provider-proxmox/issues/396)) ([b0b5fa1](https://github.com/bpg/terraform-provider-proxmox/commit/b0b5fa153253102ecf4bcae896426296188f83be))
* **file:** spurious unsupported content type warning ([#395](https://github.com/bpg/terraform-provider-proxmox/issues/395)) ([4da2b68](https://github.com/bpg/terraform-provider-proxmox/commit/4da2b682de1f2c7f456c6f7c7bc06048881cb8b9))
* **lxc:** add support for 'nixos' ([#387](https://github.com/bpg/terraform-provider-proxmox/issues/387)) ([23a5194](https://github.com/bpg/terraform-provider-proxmox/commit/23a519475d2eddf2f2145166ff2593c60c807f53))
* **provider:** better handling of root@pam token ([#386](https://github.com/bpg/terraform-provider-proxmox/issues/386)) ([03eaf72](https://github.com/bpg/terraform-provider-proxmox/commit/03eaf72767082ca4b5642538f64730dc9c4e34aa))
* **provider:** config environment variables handling caused "rpc error" ([#397](https://github.com/bpg/terraform-provider-proxmox/issues/397)) ([d748a7d](https://github.com/bpg/terraform-provider-proxmox/commit/d748a7de7b16fd792e6e3d8d6b60a951f6031ac3))
* **vm:** do not error on `read` at state refresh if VM is missing ([#398](https://github.com/bpg/terraform-provider-proxmox/issues/398)) ([253a59e](https://github.com/bpg/terraform-provider-proxmox/commit/253a59ece6c8f505362d7cd40f62a076b7caa590))
* **vm:** search for vm in cluster resources before calling node api ([#393](https://github.com/bpg/terraform-provider-proxmox/issues/393)) ([99fda9c](https://github.com/bpg/terraform-provider-proxmox/commit/99fda9cbcdbd2f254cd4c8e48559a0f7ce7a3b01))


### Miscellaneous

* **deps:** bump github.com/hashicorp/terraform-plugin-framework from 1.3.0 to 1.3.1 ([#381](https://github.com/bpg/terraform-provider-proxmox/issues/381)) ([c1219ec](https://github.com/bpg/terraform-provider-proxmox/commit/c1219ecd3c5bb3ef3c728b479056f5309a02b6a8))
* **deps:** bump github.com/hashicorp/terraform-plugin-go from 0.15.0 to 0.16.0 ([#380](https://github.com/bpg/terraform-provider-proxmox/issues/380)) ([9146703](https://github.com/bpg/terraform-provider-proxmox/commit/91467037d56bcf612c32648c0fcb5ceb2df547df))

## [0.22.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.21.1...v0.22.0) (2023-06-24)


### Features

* **vm:** add network device resources ([#376](https://github.com/bpg/terraform-provider-proxmox/issues/376)) ([343e804](https://github.com/bpg/terraform-provider-proxmox/commit/343e8045c125a4e216443855be3fc794e56399cd))
* **vm:** add support for meta-data in cloud-init ([#378](https://github.com/bpg/terraform-provider-proxmox/issues/378)) ([7aa25b8](https://github.com/bpg/terraform-provider-proxmox/commit/7aa25b8d058ae6f1807252fb731fab0aec3a4814))


### Bug Fixes

* **file:** add check for supported content types when uploading file to a storage ([#379](https://github.com/bpg/terraform-provider-proxmox/issues/379)) ([4e1ce30](https://github.com/bpg/terraform-provider-proxmox/commit/4e1ce30619ccf7db141874f4daa5873ca9f012f1))


### Miscellaneous

* **deps:** bump github.com/golangci/golangci-lint from 1.53.2 to 1.53.3 in /tools ([#375](https://github.com/bpg/terraform-provider-proxmox/issues/375)) ([2863aa6](https://github.com/bpg/terraform-provider-proxmox/commit/2863aa6e2d1a472259c8f60bd63a934c0161f598))
* **deps:** bump golang.org/x/crypto from 0.9.0 to 0.10.0 ([#374](https://github.com/bpg/terraform-provider-proxmox/issues/374)) ([f6e20bd](https://github.com/bpg/terraform-provider-proxmox/commit/f6e20bd787977b99b9d934bb6ba4d7d06244ef42))
* **deps:** bump goreleaser/goreleaser-action from 4.2.0 to 4.3.0 ([#371](https://github.com/bpg/terraform-provider-proxmox/issues/371)) ([e3a62d7](https://github.com/bpg/terraform-provider-proxmox/commit/e3a62d79ad0fc319d4f57c9ae12cfae14f8e25f6))

## [0.21.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.21.0...v0.21.1) (2023-06-07)


### Bug Fixes

* **core:** Do not limit cluster size to 1 in provider's `ssh` config ([#369](https://github.com/bpg/terraform-provider-proxmox/issues/369)) ([926382c](https://github.com/bpg/terraform-provider-proxmox/commit/926382c155169f1be07cba26b3fda0572fdc1002))
* **doc:** Update documentation for resource `proxmox_virtual_environment_firewall_ipset` ([#366](https://github.com/bpg/terraform-provider-proxmox/issues/366)) ([0aa33f0](https://github.com/bpg/terraform-provider-proxmox/commit/0aa33f0929c4b9588cce8bcde67d297137c4ddc0))
* **firewall:** Improve firewall resources argument validation ([#359](https://github.com/bpg/terraform-provider-proxmox/issues/359)) ([8c1f246](https://github.com/bpg/terraform-provider-proxmox/commit/8c1f246b5a6288195dce25ab2417ae5218b7888d))
* **vm:** fix incorrect disk interface ref when reading VM info from PVE ([#365](https://github.com/bpg/terraform-provider-proxmox/issues/365)) ([de3935d](https://github.com/bpg/terraform-provider-proxmox/commit/de3935d462cd074ae8f1bfa2ead655efec8256b7))
* **vm:** Make `vm_id` computed ([#367](https://github.com/bpg/terraform-provider-proxmox/issues/367)) ([2a5abb1](https://github.com/bpg/terraform-provider-proxmox/commit/2a5abb10fc43603d2c786ad806cba056886c7f29))


### Miscellaneous

* **deps:** bump github.com/golangci/golangci-lint from 1.52.2 to 1.53.2 in /tools ([#363](https://github.com/bpg/terraform-provider-proxmox/issues/363)) ([a546a82](https://github.com/bpg/terraform-provider-proxmox/commit/a546a8292803d9645bdba48e2aa2d6c845c70a0a))
* **deps:** bump github.com/hashicorp/terraform-plugin-log from 0.8.0 to 0.9.0 ([#362](https://github.com/bpg/terraform-provider-proxmox/issues/362)) ([170ec8a](https://github.com/bpg/terraform-provider-proxmox/commit/170ec8ad924cc4bb9683ec87cbd39d8f1e8a1ee3))
* **deps:** bump github.com/stretchr/testify from 1.8.3 to 1.8.4 ([#361](https://github.com/bpg/terraform-provider-proxmox/issues/361)) ([46e0f8f](https://github.com/bpg/terraform-provider-proxmox/commit/46e0f8f6e79ba0c859507aee97aaf3be931640cd))
* **doc:** project documentation update ([#356](https://github.com/bpg/terraform-provider-proxmox/issues/356)) ([9587c63](https://github.com/bpg/terraform-provider-proxmox/commit/9587c6383c37be894f6ca5a8d8f3edbb1826c219))

## [0.21.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.20.1...v0.21.0) (2023-06-01)


### Features

* API client cleanup and refactoring ([#323](https://github.com/bpg/terraform-provider-proxmox/issues/323)) ([1f006aa](https://github.com/bpg/terraform-provider-proxmox/commit/1f006aa82bc79125a63543dbbf765629692b7b38))
* **core:** Add ability to override node IP used for SSH connection ([80c94a5](https://github.com/bpg/terraform-provider-proxmox/commit/80c94a51262df7c5cd49a938f58c7fd09a1a3540))
* **core:** Add API Token authentication ([#350](https://github.com/bpg/terraform-provider-proxmox/issues/350)) ([ab54aa1](https://github.com/bpg/terraform-provider-proxmox/commit/ab54aa1092534c323b85c46571de33bee80ae950))


### Bug Fixes

* **vm:** Make mac_address computed, fix [#339](https://github.com/bpg/terraform-provider-proxmox/issues/339) ([#354](https://github.com/bpg/terraform-provider-proxmox/issues/354)) ([e15c4a6](https://github.com/bpg/terraform-provider-proxmox/commit/e15c4a678409a378ce10ed47fb73051e9dcdae61))


### Miscellaneous

* **deps:** bump github.com/goreleaser/nfpm/v2 from 2.28.0 to 2.29.0 in /tools ([#347](https://github.com/bpg/terraform-provider-proxmox/issues/347)) ([2358557](https://github.com/bpg/terraform-provider-proxmox/commit/23585570ab379948edab12fff85542a8aadf1792))
* **deps:** bump github.com/sigstore/rekor from 1.1.1 to 1.2.0 in /tools ([#349](https://github.com/bpg/terraform-provider-proxmox/issues/349)) ([6e59360](https://github.com/bpg/terraform-provider-proxmox/commit/6e593607bb40419d27cfeebfc3d7c11ea1061828))

## [0.20.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.20.0...v0.20.1) (2023-05-23)


### Bug Fixes

* **vm:** Regression: wait for 'net.IsGlobalUnicast' IP address  ([#100](https://github.com/bpg/terraform-provider-proxmox/issues/100)) ([#345](https://github.com/bpg/terraform-provider-proxmox/issues/345)) ([20131b0](https://github.com/bpg/terraform-provider-proxmox/commit/20131b0ffcad256835256fb28bf177c20d344482))

## [0.20.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.19.1...v0.20.0) (2023-05-22)


### Features

* bump Go to 1.20 to resolve MacOS DNS resolution issues ([#342](https://github.com/bpg/terraform-provider-proxmox/issues/342)) ([1c920de](https://github.com/bpg/terraform-provider-proxmox/commit/1c920de71d4e07bd7a29700cdffd4e6b319f95c3))
* SSH-Agent Support ([#306](https://github.com/bpg/terraform-provider-proxmox/issues/306)) ([9fa9242](https://github.com/bpg/terraform-provider-proxmox/commit/9fa92423b5b3960ee7f46fb66fc18d12fcc8af29))


### Miscellaneous

* **deps:** bump github.com/skeema/knownhosts from 1.1.0 to 1.1.1 ([#336](https://github.com/bpg/terraform-provider-proxmox/issues/336)) ([0d8e6d3](https://github.com/bpg/terraform-provider-proxmox/commit/0d8e6d31584b418677c5e8579d9e41649e9790a7))
* **deps:** bump github.com/stretchr/testify from 1.8.2 to 1.8.3 ([#343](https://github.com/bpg/terraform-provider-proxmox/issues/343)) ([fc1e03f](https://github.com/bpg/terraform-provider-proxmox/commit/fc1e03f0949ee730a8c14cf2346a053d8f1a28e2))
* **deps:** bump golang.org/x/crypto from 0.8.0 to 0.9.0 ([#337](https://github.com/bpg/terraform-provider-proxmox/issues/337)) ([b1cb49c](https://github.com/bpg/terraform-provider-proxmox/commit/b1cb49cf7a829a6112ce6e7d363e7bf537e5c52c))

## [0.19.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.19.0...v0.19.1) (2023-05-14)


### Bug Fixes

* **vm:** Regression: cannot create disks larger than 99G ([#335](https://github.com/bpg/terraform-provider-proxmox/issues/335)) ([79e5a8e](https://github.com/bpg/terraform-provider-proxmox/commit/79e5a8ebb07d9c7858a32dbef280dfab5e78c19e))


### Miscellaneous

* **deps:** bump github.com/docker/distribution from 2.8.1+incompatible to 2.8.2+incompatible in /tools ([#331](https://github.com/bpg/terraform-provider-proxmox/issues/331)) ([37a1234](https://github.com/bpg/terraform-provider-proxmox/commit/37a1234bb05cb57229d276ec568096043c2073e0))

## [0.19.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.18.2...v0.19.0) (2023-05-11)


### Features

* **vm,lxc:** Improved support for different disk size units ([#326](https://github.com/bpg/terraform-provider-proxmox/issues/326)) ([4be9914](https://github.com/bpg/terraform-provider-proxmox/commit/4be9914757cb9fee38f1c2c08772daca364b1ac9))


### Bug Fixes

* **vm,lxc:** Add validation for non-empty tags ([#330](https://github.com/bpg/terraform-provider-proxmox/issues/330)) ([8359c03](https://github.com/bpg/terraform-provider-proxmox/commit/8359c03aa8069d8816e0802e41fb36a220040673))


### Miscellaneous

* **deps:** bump github.com/goreleaser/goreleaser from 1.17.2 to 1.18.1 in /tools ([#324](https://github.com/bpg/terraform-provider-proxmox/issues/324)) ([aea079e](https://github.com/bpg/terraform-provider-proxmox/commit/aea079e0b11e2f1a07b734ed1edde3c468518429))
* **deps:** bump github.com/goreleaser/goreleaser from 1.18.1 to 1.18.2 in /tools ([#327](https://github.com/bpg/terraform-provider-proxmox/issues/327)) ([d94a4ce](https://github.com/bpg/terraform-provider-proxmox/commit/d94a4ce7cf86974b93987904d430130108d9984c))

## [0.18.2](https://github.com/bpg/terraform-provider-proxmox/compare/v0.18.1...v0.18.2) (2023-05-05)


### Bug Fixes

* **vm,lxc:** Fix tags reordering on plan re-apply ([#322](https://github.com/bpg/terraform-provider-proxmox/issues/322)) ([f0b88e3](https://github.com/bpg/terraform-provider-proxmox/commit/f0b88e336c48d76f5119fba78af9ce8b087d240e))
* **vm:** Fix IPv6 handling ([#319](https://github.com/bpg/terraform-provider-proxmox/issues/319)) ([97ca22a](https://github.com/bpg/terraform-provider-proxmox/commit/97ca22abbba4bf50895b56324ce3c3e693b46e2f))


### Miscellaneous

* **deps:** bump github.com/goreleaser/goreleaser from 1.17.1 to 1.17.2 in /tools ([#313](https://github.com/bpg/terraform-provider-proxmox/issues/313)) ([2a03818](https://github.com/bpg/terraform-provider-proxmox/commit/2a03818d4034a6ae9a0ca9153fdb2d3012cd4b97))
* **deps:** bump github.com/sigstore/rekor from 1.0.1 to 1.1.1 in /tools ([#320](https://github.com/bpg/terraform-provider-proxmox/issues/320)) ([b8184e4](https://github.com/bpg/terraform-provider-proxmox/commit/b8184e47c1af12423202385aeb7eb456e98bb42d))
* **make:** Add `lint`, `release-build` targets ([#317](https://github.com/bpg/terraform-provider-proxmox/issues/317)) ([aa99290](https://github.com/bpg/terraform-provider-proxmox/commit/aa9929066491765cfe421a7f3ede163b74473149))

## [0.18.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.18.0...v0.18.1) (2023-04-23)


### Bug Fixes

* **file:** fix SSH file upload on Windows ([#308](https://github.com/bpg/terraform-provider-proxmox/issues/308)) ([7c9505d](https://github.com/bpg/terraform-provider-proxmox/commit/7c9505d11f7cc99f6052e814f549004fe97e8b49))
* **firewall:** use correct default value for firewall ([#312](https://github.com/bpg/terraform-provider-proxmox/issues/312)) ([496bda4](https://github.com/bpg/terraform-provider-proxmox/commit/496bda4edcab1e52a3877581828080f62ef525d7))

## [0.18.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.17.1...v0.18.0) (2023-04-18)


### Features

* **vm:** Wait for the VM creation task to complete ([#305](https://github.com/bpg/terraform-provider-proxmox/issues/305)) ([8addb1d](https://github.com/bpg/terraform-provider-proxmox/commit/8addb1d1d547197ab7502b33105a17737c06788a))


### Miscellaneous

* **deps:** bump commonmarker from 0.23.8 to 0.23.9 in /docs ([#298](https://github.com/bpg/terraform-provider-proxmox/issues/298)) ([fc4a6e8](https://github.com/bpg/terraform-provider-proxmox/commit/fc4a6e8ace24db9a44102b878f09bd5d30329bd8))
* **deps:** bump github.com/goreleaser/goreleaser from 1.16.2 to 1.17.1 in /tools ([#303](https://github.com/bpg/terraform-provider-proxmox/issues/303)) ([d24f60a](https://github.com/bpg/terraform-provider-proxmox/commit/d24f60aaa22867d536a35712a0f7d8209a7d1ac2))
* **deps:** bump golang.org/x/crypto from 0.7.0 to 0.8.0 ([#296](https://github.com/bpg/terraform-provider-proxmox/issues/296)) ([a896b50](https://github.com/bpg/terraform-provider-proxmox/commit/a896b5051ec6103cd8449234102917d9f17a1011))
* **deps:** bump nokogiri from 1.14.2 to 1.14.3 in /docs ([#299](https://github.com/bpg/terraform-provider-proxmox/issues/299)) ([6722492](https://github.com/bpg/terraform-provider-proxmox/commit/672249246f921eaba7ffa5020c0427bf95c3ed29))

## [0.17.1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.17.0...v0.17.1) (2023-04-10)


### Bug Fixes

* **core:** Error when open SSH session on Windows ([#293](https://github.com/bpg/terraform-provider-proxmox/issues/293)) ([be3995e](https://github.com/bpg/terraform-provider-proxmox/commit/be3995e969e16eac08c3e1d0fbaadb60244a5576))
* **file:** "Permission denied" error when creating a file by a non-root user ([#291](https://github.com/bpg/terraform-provider-proxmox/issues/291)) ([401b397](https://github.com/bpg/terraform-provider-proxmox/commit/401b39782f857382b30ab71b3e49a8ab44fbac48))
* **firewall:** Add support for `firewall` flag for LXC/VM net adapters ([#295](https://github.com/bpg/terraform-provider-proxmox/issues/295)) ([f4783f8](https://github.com/bpg/terraform-provider-proxmox/commit/f4783f8cda701b6800403d50840240da6469fd38))

## [0.17.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.17.0-rc2...v0.17.0) (2023-04-07)


### Features

* **vm:** add support for `boot_order` argument for VM ([#219](https://github.com/bpg/terraform-provider-proxmox/issues/219)) ([f9e263a](https://github.com/bpg/terraform-provider-proxmox/commit/f9e263ad5edf47fb12f5321af0090d928da50d42))


### Bug Fixes

* **provider:** Deprecate `virtual_environment` block ([#288](https://github.com/bpg/terraform-provider-proxmox/issues/288)) ([ed3dfea](https://github.com/bpg/terraform-provider-proxmox/commit/ed3dfeae9907757f42c0cce63fe1f00a4e2ec0a2))


### Miscellaneous

* cleanup and fix linter errors ([#290](https://github.com/bpg/terraform-provider-proxmox/issues/290)) ([2fa9229](https://github.com/bpg/terraform-provider-proxmox/commit/2fa922930f6c3b6c1e0c32789b44ef6ab9189e6d))

## [0.17.0-rc2](https://github.com/bpg/terraform-provider-proxmox/compare/v0.17.0-rc1...v0.17.0-rc2) (2023-04-04)


### Bug Fixes

* **firewall:** fw controls bugfixes ([#287](https://github.com/bpg/terraform-provider-proxmox/issues/287)) ([1bfc29e](https://github.com/bpg/terraform-provider-proxmox/commit/1bfc29e2cc3342699f491d0225da474078220ecd))


### Miscellaneous

* **deps:** bump activesupport from 7.0.4.2 to 7.0.4.3 in /docs ([#285](https://github.com/bpg/terraform-provider-proxmox/issues/285)) ([fc08e19](https://github.com/bpg/terraform-provider-proxmox/commit/fc08e19f867ef652ae7597e89fd49fb3ecc3a9a8))

## [0.17.0-rc1](https://github.com/bpg/terraform-provider-proxmox/compare/v0.16.0...v0.17.0-rc1) (2023-04-02)


### Features

* Add firewall resources ([#246](https://github.com/bpg/terraform-provider-proxmox/issues/246)) ([98e1cff](https://github.com/bpg/terraform-provider-proxmox/commit/98e1cff7fef0f24d932933bcba56ebc5b6ca7548))

## [0.16.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.15.0...v0.16.0) (2023-04-02)


### Features

* Update to Go 1.19 ([#280](https://github.com/bpg/terraform-provider-proxmox/issues/280)) ([8edfe9c](https://github.com/bpg/terraform-provider-proxmox/commit/8edfe9c7c54c9554adca52ffcf31c091f1fce11f))
* **vm:** Add scsi_hardware field ([#282](https://github.com/bpg/terraform-provider-proxmox/issues/282)) ([f0f31ee](https://github.com/bpg/terraform-provider-proxmox/commit/f0f31eee470dc954fdd5d1c952ea3067a2a68f1b))


### Miscellaneous

* add missing docs ([#283](https://github.com/bpg/terraform-provider-proxmox/issues/283)) ([db7afe2](https://github.com/bpg/terraform-provider-proxmox/commit/db7afe2e4a93ae3f97c0533dba0adaf82123c49f))
* **deps:** bump actions/stale from 7 to 8 ([#276](https://github.com/bpg/terraform-provider-proxmox/issues/276)) ([edd9685](https://github.com/bpg/terraform-provider-proxmox/commit/edd96857e64f73b041ed76a9c1818a864b4a0cca))
* **deps:** bump github.com/golangci/golangci-lint from 1.52.1 to 1.52.2 in /tools ([#278](https://github.com/bpg/terraform-provider-proxmox/issues/278)) ([d8c1fb3](https://github.com/bpg/terraform-provider-proxmox/commit/d8c1fb3573de553bc2eb26d8e37cdcfe8a78f384))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.25.0 to 2.26.1 ([#277](https://github.com/bpg/terraform-provider-proxmox/issues/277)) ([b403a49](https://github.com/bpg/terraform-provider-proxmox/commit/b403a4940fd31eacef90deaa11f2696fa7c03910))

## [0.15.0](https://github.com/bpg/terraform-provider-proxmox/compare/v0.14.1...v0.15.0) (2023-03-25)


### Features

* **vm:** Add bare minimum VM datasource ([#268](https://github.com/bpg/terraform-provider-proxmox/issues/268)) ([c2d3f46](https://github.com/bpg/terraform-provider-proxmox/commit/c2d3f46474fc0d0603c34596eb81b82c06713b17))


### Bug Fixes

* **vm:** Prevent `file_format` override with default `qcow2` in TF state ([#275](https://github.com/bpg/terraform-provider-proxmox/issues/275)) ([17dca98](https://github.com/bpg/terraform-provider-proxmox/commit/17dca987eb240454dbd980ed8f0c4a939e327ff0))


### Miscellaneous

* **deps:** bump actions/setup-go from 3 to 4 ([#269](https://github.com/bpg/terraform-provider-proxmox/issues/269)) ([fdb9dc7](https://github.com/bpg/terraform-provider-proxmox/commit/fdb9dc7714f12a1682f47e57ab319753fc18e4f4))
* **deps:** bump github.com/golangci/golangci-lint from 1.51.2 to 1.52.1 in /tools ([#274](https://github.com/bpg/terraform-provider-proxmox/issues/274)) ([1150163](https://github.com/bpg/terraform-provider-proxmox/commit/1150163b4b66249b79e446cf20f9df54fc204f7c))
* **deps:** bump github.com/goreleaser/goreleaser from 1.16.1 to 1.16.2 in /tools ([#271](https://github.com/bpg/terraform-provider-proxmox/issues/271)) ([7a0e1db](https://github.com/bpg/terraform-provider-proxmox/commit/7a0e1db6c4a117e0599b4899df781eb3c3fe600f))

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
