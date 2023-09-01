# Changelog

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
