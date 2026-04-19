<!-- markdownlint-disable MD060 -->

# Audit: `proxmox_vm` against ADRs 001–007 (#1231)

**Issue:** [bpg/terraform-provider-proxmox#1231](https://github.com/bpg/terraform-provider-proxmox/issues/1231)
**Design:** [1231_DESIGN.md](1231_DESIGN.md)
**Tracker:** [1231_TRACKER.md](1231_TRACKER.md)
**Scope:** PR #1 (Phase 1)
**Status:** In progress — frozen after PR #1 merges

## Methodology

| Source | What it produces |
|---|---|
| File:line scan of `fwprovider/nodes/vm/` | ADR compliance findings (Section 1) |
| Walk of `proxmoxtf/resource/vm/vm.go` | Capabilities inventory (Section 2) |
| Walk of `fwprovider/test/resource_vm_*.go` + `proxmoxtf/resource/vm/**/*_test.go` | Legacy test inventory (Section 3) |
| Per-attribute mitmproxy trace + reasoning | Per-attribute classification (Section 4) |
| Validator-by-validator review against ADR-004 amendment | Validator inventory (Section 5) |
| Open questions resolution from design | Q5 power_state resolution (Section 6) |
| Grep of `proxmox/nodes/vms` consumers | Shared-types catalog (Section 7) |

Severity tags used in Section 1:

| Tag | Meaning |
|---|---|
| `blocker` | Must fix before Phase 2 (or ship it as PR #3 alongside the contract port) |
| `should-fix` | Fix during Phase 1 (PRs #3–#5) |
| `nit` | Defer; capture in gap matrix as deferred-cleanup |

---

## Section 1 — ADR compliance findings

### Scope of scan

| Path | Files | Lines (approx) |
|---|---|---|
| `fwprovider/nodes/vm/` (top-level) | `resource.go`, `resource_schema.go`, `resource_short.go`, `model.go`, `datasource.go`, `datasource_schema.go`, `datasource_short.go`, `concurrency_test.go`, `datasource_test.go`, `resource_test.go` | ~1500 |
| `fwprovider/nodes/vm/cpu/` | `resource.go`, `resource_schema.go`, `model.go`, `datasource_schema.go`, `resource_test.go` | ~830 |
| `fwprovider/nodes/vm/cdrom/` | `resource.go`, `resource_schema.go`, `model.go`, `datasource_schema.go`, `resource_test.go` | ~340 |
| `fwprovider/nodes/vm/vga/` | `resource.go`, `resource_schema.go`, `model.go`, `datasource_schema.go`, `resource_test.go` | ~400 |
| `fwprovider/nodes/vm/rng/` | `resource.go`, `resource_schema.go`, `model.go`, `datasource_schema.go`, `resource_test.go` | ~400 |
| `fwprovider/nodes/vm/memory/` | `resource.go`, `resource_schema.go`, `model.go` | ~270 |
| `fwprovider/nodes/vm/network/` | (empty placeholder) | 0 |

### Findings

> Each finding: `area`, `severity`, `ADR`, `file:line`, `description`, `target PR`. Pre-resolved findings from the design grilling (P1–P7) are listed at the bottom of this document; the table below holds *new* findings.

#### Top-level package (`fwprovider/nodes/vm/`)

| # | Area | Severity | ADR | File:line | Description | Target PR |
|---|---|---|---|---|---|---|
| F1 | error msg | should-fix | 005 | `resource.go:116` | `"Failed to generate VM ID"` — should be `"Unable to Generate VM ID"` | #5 |
| F2 | error msg | should-fix | 005 | `resource.go:136` | `"VM does not exist after creation", ""` — empty detail; should be `"Unable to Create VM N"` with detail | #5 |
| F3 | error msg | should-fix | 005 | `resource.go:168, 177, 363, 365, 378` | Generic context strings (`"VM create"`, `"VM template conversion"`, `"VM stop/shutdown"`, `"VM delete"`) for `AddDiags`/`AddDiagsAsWarnings`. Format inconsistent with ADR-005. | #5 |
| F4 | error msg | should-fix | 005 | `resource.go:237` | `"VM does not exist after update", ""` — empty detail | #5 |
| F5 | error msg | should-fix | 005 | `resource.go:295` | `"Failed to update VM"` — should be `"Unable to Update VM N"` | #5 |
| F6 | error msg | should-fix | 005 | `resource.go:310, 353` | `"Failed to get VM status"` — should be `"Unable to Get VM %d Status"` (two callsites) | #5 |
| F7 | error msg | should-fix | 005 | `resource.go:326` | `"Cannot convert template back to VM"` — should be `"Unable to Convert Template Back to VM"` | #5 |
| F8 | error msg | should-fix | 005 | `resource.go:374` | `"Unable to Delete VM"` — correct prefix but missing VM ID per CLAUDE.md ADR-005 note | #5 |
| F9 | error msg | should-fix | 005 | `resource.go:425` | `fmt.Sprintf("VM %d does not exist on node %s", ...)` — should be `"Unable to Import VM N"` with detail | #5 |
| F10 | sub-block contract | should-fix | 008 | `resource.go:256–282` | Hand-rolled `del()` closure for top-level scalars (Description, Name, Tags). Should use `attribute.CheckDelete(plan, state, body, "FieldName")` per ADR-008. | #3 |
| F11 | code clarity | nit | — | `resource.go:432–435` | Comment `"not clear why this is needed, but ImportStateVerify fails without it"` for setting StopOnDestroy/PurgeOnDestroy/DeleteUnreferencedDisksOnDestroy on import. Either replace with proper explanation or fix root cause. | #4 (during rename pass) |
| F12 | resource type name | blocker | 007 | `resource_short.go:28–34, 49–57` | `resourceShort` wrapper struct + constructor + `MoveState` method exist. Per design D2, this collapses into single `NewResource()` returning `proxmox_vm` in PR #4. (Line 17 is just the `shortResourceTypeName` const — fine in itself; the issue is the wrapper structure.) | #4 |
| F13 | datasource type name | blocker | 007 | `datasource_short.go:15, 29–46` | `datasourceShort` wrapper. Collapses in PR #4 alongside resourceShort. | #4 |
| F14 | datasource type name | should-fix | 007 | `datasource.go:41` | `req.ProviderTypeName + "_vm2"` — long datasource name still in use. Parallels F12/F13. | #4 |
| F15 | error msg | should-fix | 005 | `datasource.go:89–92` | `"VM Not Found"` summary + custom detail; should be `"Unable to Read VM N"` per ADR-005 | #5 |
| F16 | error msg | should-fix | 005 | `model.go:74, 82` | `"Unable to Read VM"`, `"Unable to Read VM Status"` — correct prefix, missing VM ID | #5 |
| F17 | error msg | should-fix | 005 | `model.go:87, 144` | `"VM ID is missing in status API response"` — wrong format; should be `"Unable to Read VM N"` with detail (two callsites) | #5 |
| F18 | error msg | should-fix | 005 | `model.go:131` | `"Failed to get VM"` — should be `"Unable to Read VM N"` | #5 |
| F19 | error msg | should-fix | 005 | `model.go:139` | `"Failed to get VM status"` — duplicate of F6 in resource read path | #5 |

#### `cpu/` sub-package

| # | Area | Severity | ADR | File:line | Description | Target PR |
|---|---|---|---|---|---|---|
| F20 | sentinel | should-fix | 004 | `cpu/resource.go:38–42` | (Confirms P2) Nil-substitution sentinel `Cores == nil → 1`. Drop per ADR-004 PVE-defaults rule. | #3 |
| F21 | sentinel | should-fix | 004 | `cpu/resource.go:44–48` | **Additional sentinel** beyond P2: `Sockets == nil → 1`. Drop per ADR-004. | #3 |
| F22 | sentinel | should-fix | 004 | `cpu/resource.go:50–60` | **Additional sentinel** beyond P2: `CPUEmulation == nil → Type "kvm64", Flags null`. Drop per ADR-004. | #3 |
| F23 | sub-block contract | should-fix | 008 | `cpu/resource.go:159–225` | (Confirms P4) Hand-rolled `del()` + `ShouldBeRemoved` + `IsDefined` cascades for 8 fields. Replace with `attribute.CheckDelete`. | #3 |
| F24 | bug | should-fix | 008 | `cpu/resource.go:190` | (Confirms P1) `IsDefined(plan.Sockets)` copy-paste bug in Limit branch — should be `IsDefined(plan.Limit)` | #3 |
| F25 | error msg | nit | 005 | `cpu/resource.go:250` | `"Cannot have CPU flags without explicit definition of CPU type", ""` — empty detail, summary phrasing inconsistent with ADR-005 format | #5 |
| F26 | sub-block contract | should-fix | 008 | `cpu/resource.go:227–255` | Special-case for `CPUEmulation` compound update (delType/delFlags switch). Doesn't fit standard CheckDelete shape — ADR-008 should call out compound types as a recognized pattern (or refactor) | #3 |
| F27 | validator | should-fix | 004 | `cpu/resource_schema.go:124–204` | (Confirms P5) Long enum validator (~75 CPU types) for `cpu.type`. Drop per ADR-004 enum rule. | #3 |
| F28 | classification | should-fix | 004 | `cpu/resource_schema.go:31–122` | All 10 CPU attributes are `Optional+Computed`. Per-attribute classification needed against PVE Read behavior (Section 4) | #3 |
| F29 | rehome | should-fix | — | `cpu/resource_schema.go:83` | (Confirms P3) `cpu.hotplugged` exists. Drop in #3, rehome as top-level `vcpus` in #14. | #3 (drop) / #14 (rehome) |
| F30 | rehome | should-fix | — | `cpu/resource_schema.go:101` | (Confirms P3) `cpu.numa` exists. Drop in #3, rehome as `numa.enabled` in #13. | #3 (drop) / #13 (rehome) |

#### `vga/` sub-package

| # | Area | Severity | ADR | File:line | Description | Target PR |
|---|---|---|---|---|---|---|
| F31 | sub-block contract | should-fix | 008 | `vga/resource.go:107–129` | Hand-rolled cascades for Clipboard/Type/Memory in FillUpdateBody. Replace with CheckDelete pattern. | #3 |
| F32 | sub-block contract | should-fix | 008 | `vga/resource.go:101–105` | `vgaDevice` initialized from state (not zero) and mutated. Different mutation pattern from cpu — ADR-008 should normalize. | #3 |
| F32a | sub-block anti-pattern | should-fix | 008 | `vga/resource.go:72, 131` | `if !reflect.DeepEqual(vgaDevice, &vms.CustomVGADevice{}) { ... }` zero-struct comparison. In FillUpdateBody (line 131) this is always true because `vgaDevice` is initialized from state — every Update sends the entire vga block to PVE. ADR-008 should explicitly reject this anti-pattern. | #3 |
| F33 | validator | should-fix | 004 | `vga/resource_schema.go:55–72` | Long enum validator (14 VGA types, version-evolving). Drop per ADR-004 enum rule. | #3 |
| F34 | classification | should-fix | 004 | `vga/resource_schema.go:33–82` | All 3 VGA attributes are `Optional+Computed`. Section 4 classification: all drop Computed → Optional only (no PVE auto-populate per Finding 3). | #3 |

#### `rng/` sub-package

| # | Area | Severity | ADR | File:line | Description | Target PR |
|---|---|---|---|---|---|---|
| F35 | sub-block contract | should-fix | 008 | `rng/resource.go:117–141` | Hand-rolled cascades for Source/MaxBytes/Period in FillUpdateBody. Replace with CheckDelete pattern. | #3 |
| F36 | sub-block contract | should-fix | 008 | `rng/resource.go:115` | `rngDevice = createRNGDevice(state, true)` — same state-initialized mutation pattern as vga (F32) | #3 |
| F36a | sub-block anti-pattern | should-fix | 008 | `rng/resource.go:86, 143` | Same `reflect.DeepEqual(&vms.CustomRNGDevice{})` anti-pattern as F32a in both `FillCreateBody` (line 86) and `FillUpdateBody` (line 143). ADR-008 should explicitly reject. | #3 |
| F37 | int-zero trap | nit | 004 | `rng/resource.go:54, 59` | `MaxBytes.ValueInt64() != 0` and `Period.ValueInt64() != 0` use 0 as "not set" sentinel. Documented as `"Use 0 to disable limiting"` in schema, but the FillCreateBody never sends 0 to PVE — meaning user-set 0 is silently dropped. ADR-004 integer-0 trap. | #3 |
| F38 | classification | should-fix | 004 | `rng/resource_schema.go:31–67` | All 3 RNG attributes `Optional+Computed`. Section 4 classification: all drop Computed → Optional only (no PVE auto-populate per Finding 3). | #3 |

#### `memory/` sub-package

| # | Area | Severity | ADR | File:line | Description | Target PR |
|---|---|---|---|---|---|---|
| F39 | provider default | should-fix (PR-#3-blocker) | 004 | `memory/resource_schema.go:52, 66, 79` | `Default(...)` for `size=512`, `balloon=0`, `shares=1000` — these duplicate PVE defaults. Will violate ADR-004 amendment once PR #2 lands. **Currently functional**: `clonedvm` (the only consumer today) works fine; this is a forward-looking violation, not a production regression. Must be fixed in PR #3 alongside the contract port. | #3 |
| F40 | NewValue sentinel | should-fix | 004 | `memory/resource.go:35–40` | Nil-substitution sentinel **in `NewValue`** (API → state): `DedicatedMemory == nil → types.Int64Value(512)`. Should be `types.Int64PointerValue(nil)` so PVE's absent → state null. Drop per ADR-004. | #3 |
| F41 | NewValue sentinel | should-fix | 004 | `memory/resource.go:42–48` | Same pattern in `NewValue` for `FloatingMemory == nil → 0`. Drop per ADR-004. | #3 |
| F42 | NewValue sentinel | should-fix | 004 | `memory/resource.go:50–56` | Same pattern in `NewValue` for `FloatingMemoryShares == nil → 1000`. Drop per ADR-004. | #3 |
| F43 | sub-block contract | should-fix (PR-#6-blocker) | 008 | `memory/resource.go` (no `FillCreateBody`) | `memory/` package has **no** `FillCreateBody`. **Currently functional**: `clonedvm` calls only `FillUpdateBody` (clone semantics never call Create with config), so the absence isn't blocking. Becomes blocking when PR #6 wires memory into `proxmox_vm` (which uses Create). PR #3 must add `FillCreateBody`. PR #3 fix is two-part with F40-F42: rewrite `NewValue` to return null on nil, **and** add `FillCreateBody` that handles null/unknown plan values. | #3 |
| F44 | sub-block contract | should-fix | 008 | `memory/resource.go:75–122` | `FillUpdateBody` signature is `(ctx, planValue, body, diags)` — no `stateValue` parameter. Fields are set if present in plan but never deleted. Cannot remove `hugepages` or `keep_hugepages` once set. Diverges from ADR-008 update-body shape. PR #3 must add stateValue + CheckDelete logic. | #3 |
| F45 | classification | should-fix | 004 | `memory/resource_schema.go:43–105` | After F39 `Default` removal, all 5 attributes need ADR-004 classification (Section 4). Note: Section 4's memory classifications are predicted from vga/rng pattern — verify with mitmproxy in PR #6. | #3 |

#### `cdrom/` sub-package (reference for map-keyed pattern)

| # | Area | Severity | ADR | File:line | Description | Target PR |
|---|---|---|---|---|---|---|
| F46 | regex bound | nit | 008 | `cdrom/resource_schema.go:33` | Slot regex `^(ide[0-3]\|sata[0-5]\|scsi([0-9]\|1[0-3]))$` — `scsi` only goes to 13, but PVE bound is `MAX_SCSI_DISKS=31`. Restrict-then-relax behavior is OK (additive), but design's slot-regex table specifies `scsi([0-9]\|[12][0-9]\|30)`. Tighten or relax to match. | #3 |
| F47 | provider default | nit | 004 | `cdrom/resource_schema.go:46` | `Default: stringdefault.StaticString("cdrom")` for `file_id`. Verify against PVE Read behavior — does PVE auto-default file_id to `"cdrom"` for an attached storage device, or is this a provider sentinel? | #3 |

### Summary by severity

| Severity | Count |
|---|---|
| blocker | 2 (F12, F13 — actual production-blocking issues) |
| should-fix (PR-#3-blocker) | 2 (F39, F43 — must fix during PR #3 / #6) |
| should-fix | 35 |
| nit | 6 |
| Total new findings | 47 (added F32a, F36a after scrutiny pass) |

Plus 7 pre-resolved findings (P1–P7) from grilling. Combined: 54 findings.

### Summary by target PR

| PR | New findings |
|---|---|
| #3 (port sub-packages) | 27 (added F32a, F36a) |
| #4 (rename) | 4 (F11, F12, F13, F14) |
| #5 (error sweep) | 16 |

---

## Section 2 — Capabilities inventory

Every attribute and block of the legacy SDK `proxmox_virtual_environment_vm`
classified as one of:

| Status | Meaning |
|---|---|
| `done` | Already implemented in `proxmox_vm` |
| `planned` | Target Phase 2 PR identified |
| `deliberately dropped` | Out of scope; document why |
| `open question` | Needs maintainer decision before PR can land |

### Top-level scalars

| SDK key | SDK source | Status | Target PR | Notes |
|---|---|---|---|---|
| `description` | `vm.go:206` | done | — | Already in `proxmox_vm` |
| `name` | `vm.go:268` | done | — | Already in `proxmox_vm` (with DNS validator) |
| `node_name` | `vm.go:270` | done | — | Already in `proxmox_vm` (Required) |
| `tags` | `vm.go:295` | done | — | Already in `proxmox_vm` (stringset) |
| `template` | `vm.go:296` | done | — | Already in `proxmox_vm` (RequiresReplace planmodifier) |
| `vm_id` | `vm.go:313` | done (as `id`) | — | Already in `proxmox_vm`; renamed |
| `pool_id` | `vm.go:273` | planned | #18 | — |
| `protection` | `vm.go:274` | planned | #18 | — |
| `migrate` | `vm.go:267` | planned | #19 | — |
| `acpi` | `vm.go:165` | planned | #14 | — |
| `bios` | `vm.go:184` | planned | #8 | — |
| `boot_order` | `vm.go:164` | planned | #8 | — |
| `hook_script_file_id` | `vm.go:315` | planned | #18 | — |
| `hotplug` | `vm.go:232` | planned | #14 | — |
| `keyboard_layout` | `vm.go:258` | planned | #14 | — |
| `kvm_arguments` | `vm.go:259` | planned | #14 | — |
| `machine` | `vm.go:260` | planned | #8 | — |
| `scsi_hardware` | `vm.go:314` | planned | #9 | — |
| `tablet_device` | `vm.go:294` | planned | #14 | — |
| `vm_id` (clone source) | n/a | dropped | — | Belongs to `proxmox_cloned_vm`, out of scope (D4) |

### Top-level blocks

| SDK key | SDK source | Status | Target PR | Notes |
|---|---|---|---|---|
| `cpu` | `vm.go:195` | done | — | In production. PR #3 ports to ADR-008. `numa`/`hotplugged` rehoming P3. |
| `vga` | `vm.go:309` | done | — | In production; PR #3 ports to ADR-008. Drop long-enum `type` validator. |
| `rng` | `vm.go:275` | done | — | In production; PR #3 ports to ADR-008. |
| `cdrom` | `vm.go:185` | done | — | In production (map-keyed, reference impl). PR #3 tightens slot regex (F46). |
| `memory` | `vm.go:261` | wired in #6 | #6 | Package exists with critical violations (F39, F43); PR #3 fixes contract; PR #6 wires into `proxmox_vm` |
| `disk` | `disk/schema.go:30` (MkDisk) | planned | #7 | First new map-keyed application |
| `network_device` | `network/schema.go:32` (MkNetworkDevice) | planned | #10 | Map-keyed |
| `agent` | `vm.go:166` | planned | #13 | Includes `enabled`, `timeout`, `trim`, `type`, `wait_for_ip` |
| `numa` (NUMA topology block) | `vm.go:208` | planned | #13 | Distinct from `cpu.numa` boolean; map-keyed `numa[N]` per `MAX_NUMA=8` |
| `efi_disk` | `vm.go:215` | planned | #9 | Single-nested per ADR-008 architectural-single rule |
| `tpm_state` | `vm.go:220` | planned | #9 | Single-nested per ADR-008 architectural-single rule |
| `hostpci` | `vm.go:223` | planned | #16 | Map-keyed per `MAX_HOSTPCI_DEVICES=16` |
| `usb` | `vm.go:305` | planned | #15 | Map-keyed per `MAX_USB_DEVICES=14` |
| `serial_device` | `vm.go:279` | planned | #17 | Map-keyed per `MAX_SERIAL_PORTS=4` |
| `audio_device` | `vm.go:180` | planned | #17 | Map-keyed (one-key today) per ADR-008 forward-compat rule |
| `virtiofs` | `vm.go:319` | planned | #17 | Map-keyed |
| `watchdog` | `vm.go:325` | planned | #13 | Single-nested (one watchdog per VM) |
| `initialization` (cloud-init) | `vm.go:233` | planned | #11 | Single-nested with nested DNS, IP config, user account, file ID blocks |
| `operating_system` | `vm.go:271` | planned | #12 | Single-nested with `type` field |
| `smbios` | `vm.go:281` | planned | #12 | Single-nested with family/manufacturer/product/SKU/serial/uuid/version |
| `amd_sev` | `vm.go:174` | planned | #18 | Single-nested with `type`, `allow_smt`, `kernel_hashes`, `no_debug`, `no_key_sharing` |
| `startup` | `vm.go:290` | planned | #18 | Single-nested with `order`, `up_delay`, `down_delay` |
| `clone` | `vm.go:189` | dropped | — | Out of scope — belongs to `proxmox_cloned_vm` (D4) |

### Disk-family sub-attributes (under map-keyed `disk[slot]` block, PR #7)

| SDK key | SDK source | Status | Target PR | Notes |
|---|---|---|---|---|
| `aio` | `disk/schema.go:31` | planned | #7 | — |
| `backup` | `disk/schema.go:32` | planned | #7 | — |
| `cache` | `disk/schema.go:33` | planned | #7 | — |
| `datastore_id` | `disk/schema.go:34` | planned | #7 | — |
| `discard` | `disk/schema.go:35` | planned | #7 | — |
| `file_format` | `disk/schema.go:36` | planned | #7 | — |
| `file_id` | `disk/schema.go:37` | planned | #7 | — |
| `import_from` | `disk/schema.go:38` | planned | #7 | — |
| `interface` (legacy slot field) | `disk/schema.go:39` | dropped | — | Replaced by map key per ADR-008 map-keyed pattern |
| `iops_read` / `iops_read_burstable` / `iops_write` / `iops_write_burstable` | `disk/schema.go:40–43` | planned | #7 | — |
| `iothread` | `disk/schema.go:44` | planned | #7 | — |
| `path_in_datastore` | `disk/schema.go:45` | planned | #7 | — |
| `replicate` | `disk/schema.go:46` | planned | #7 | — |
| `serial` | `disk/schema.go:47` | planned | #7 | — |
| `size` | `disk/schema.go:48` | planned | #7 | — |
| `speed` (read/read_burstable/write/write_burstable nested) | `disk/schema.go:49–53` | planned | #7 | — |
| `ssd` | `disk/schema.go:54` | planned | #7 | — |

### Network-family sub-attributes (under map-keyed `network_device[slot]` block, PR #10)

| SDK key | SDK source | Status | Target PR | Notes |
|---|---|---|---|---|
| `bridge` | `network/schema.go:33` | planned | #10 | — |
| `disconnected` | `network/schema.go:34` | planned | #10 | — |
| `enabled` | `network/schema.go:35` | planned | #10 | — |
| `firewall` | `network/schema.go:36` | planned | #10 | — |
| `mac_address` | `network/schema.go:37` | planned | #10 | — |
| `mtu` | `network/schema.go:38` | planned | #10 | — |
| `model` | `network/schema.go:39` | planned | #10 | — |
| `queues` | `network/schema.go:40` | planned | #10 | — |
| `rate_limit` | `network/schema.go:41` | planned | #10 | — |
| `trunks` | `network/schema.go:42` | planned | #10 | — |
| `vlan_id` | `network/schema.go:43` | planned | #10 | — |
| `ipv4_addresses` (read-only) | `network/schema.go:27` | planned | #10 | Computed; read from agent |
| `ipv6_addresses` (read-only) | `network/schema.go:28` | planned | #10 | Computed; read from agent |
| `mac_addresses` (read-only) | `network/schema.go:29` | planned | #10 | Computed |
| `network_interface_names` (read-only) | `network/schema.go:44` | planned | #10 | Computed |

### Watchdog sub-attributes (under single-nested `watchdog` block, PR #13)

| SDK key | SDK source | Status | Target PR | Notes |
|---|---|---|---|---|
| `enabled` | `vm.go:327` | planned | #13 | Standard enable/disable |
| `model` | `vm.go:328` | planned | #13 | Watchdog hardware model (e.g., `i6300esb`, `ib700`) |
| `action` | `vm.go:329` | planned | #13 | Action on watchdog timeout (e.g., `reset`, `shutdown`, `poweroff`) |

### Agent sub-attributes (under single-nested `agent` block, PR #13)

| SDK key | SDK source | Status | Target PR | Notes |
|---|---|---|---|---|
| `enabled` | `vm.go:167` | planned | #13 | — |
| `timeout` | `vm.go:168` | planned | #13 | Per OQ4: keep as PVE pass-through |
| `trim` | `vm.go:169` | planned | #13 | — |
| `type` | `vm.go:170` | planned | #13 | — |
| `wait_for_ip` | `vm.go:171` | planned | #13 | Nested |
| `wait_for_ip.ipv4` | `vm.go:172` | planned | #13 | Nested |
| `wait_for_ip.ipv6` | `vm.go:173` | planned | #13 | Nested |

### AMD SEV sub-attributes (under single-nested `amd_sev` block, PR #18)

| SDK key | SDK source | Status | Target PR | Notes |
|---|---|---|---|---|
| `type` | `vm.go:175` | planned | #18 | — |
| `allow_smt` | `vm.go:176` | planned | #18 | — |
| `kernel_hashes` | `vm.go:177` | planned | #18 | — |
| `no_debug` | `vm.go:178` | planned | #18 | — |
| `no_key_sharing` | `vm.go:179` | planned | #18 | — |

### Audio device sub-attributes (under map-keyed `audio_device` block, PR #17)

| SDK key | SDK source | Status | Target PR | Notes |
|---|---|---|---|---|
| `device` | `vm.go:181` | planned | #17 | — |
| `driver` | `vm.go:182` | planned | #17 | Long enum — drop validator per ADR-004 (Q4) |
| `enabled` | `vm.go:183` | planned | #17 | — |

### NUMA sub-attributes (under map-keyed `numa[N]` block, PR #13)

| SDK key | SDK source | Status | Target PR | Notes |
|---|---|---|---|---|
| `device` | `vm.go:209` | planned | #13 | — |
| `cpus` | `vm.go:210` | planned | #13 | — |
| `hostnodes` | `vm.go:211` | planned | #13 | — |
| `memory` | `vm.go:212` | planned | #13 | — |
| `policy` | `vm.go:213` | planned | #13 | — |
| `enabled` (rehomed `cpu.numa`) | (new) | planned | #13 | Per design D7/P3 |

### EFI Disk + TPM State + HostPCI + USB + Serial + Virtiofs + SMBIOS + OS + Startup sub-attributes

> Sub-attribute tables omitted for brevity — each Phase 2 PR's first commit must enumerate sub-attributes from the corresponding `mk*` constants (`mkEFIDisk*` lines 215-219, `mkTPMState*` lines 220-222, `mkHostPCI*` lines 223-231, `mkHostUSB*` lines 305-308, `mkSerialDevice*` lines 279-280, `mkVirtiofs*` lines 319-324, `mkSMBIOS*` lines 281-288, `mkOperatingSystem*` lines 271-272, `mkStartup*` lines 290-293).

### Cloud-init sub-attributes (under `initialization` block, PR #11)

| SDK key | SDK source | Status | Target PR | Notes |
|---|---|---|---|---|
| `datastore_id` | `vm.go:234` | planned | #11 | — |
| `interface` | `vm.go:235` | planned | #11 | — |
| `file_format` | `vm.go:236` | planned | #11 | — |
| `dns.domain` / `dns.servers` | `vm.go:237–239` | planned | #11 | Nested |
| `ip_config.ipv4.address` / `.ipv4.gateway` | `vm.go:241–243` | planned | #11 | Nested, map-keyed by interface |
| `ip_config.ipv6.address` / `.ipv6.gateway` | `vm.go:244–246` | planned | #11 | Nested, map-keyed by interface |
| `type` | `vm.go:247` | planned | #11 | Cloud-init type (`nocloud`, `configdrive2`, `opennebula`) |
| `upgrade` | `vm.go:248` | planned | #11 | — |
| `user_account.username` / `.password` / `.keys` | `vm.go:249–252` | planned | #11 | Nested |
| `user_data_file_id` / `vendor_data_file_id` / `network_data_file_id` / `meta_data_file_id` | `vm.go:253–256` | planned | #11 | — |

### Runtime / lifecycle attributes

| SDK key | SDK source | Status | Target PR | Notes |
|---|---|---|---|---|
| `started` | `vm.go:289` | dropped | — | Replaced by `power_state` (Q5/PR #6) |
| `reboot` (after creation) | `vm.go:161` | dropped | — | Provider decides from pending changes (Q5/PR #6) |
| `reboot_after_update` | `vm.go:162` | dropped | — | Same as above |
| `on_boot` | `vm.go:163` | planned | #6 | PVE "Start at boot" — `Optional` only per ADR-004 amendment |
| `stop_on_destroy` | `vm.go:316` | done | — | Already in `proxmox_vm` (provider-only `Optional+Default`) |
| `purge_on_destroy` | `vm.go:317` | done | — | Already in `proxmox_vm` |
| `delete_unreferenced_disks_on_destroy` | `vm.go:318` | done | — | Already in `proxmox_vm` |
| `power_state` | (new) | planned | #6 | New attribute per Q5 |
| `status` (Computed) | (new) | planned | #6 | Mirror of PVE runtime status |
| `vcpus` (top-level, rehomed `cpu.hotplugged`) | (new) | planned | #14 | Per design D7/P3 |
| `numa.enabled` (rehomed `cpu.numa`) | (new) | planned | #13 | Per design D7/P3 |

### Timeouts (folded into `timeouts` block under ADR-006)

| SDK key | SDK source | Status | Target PR | Notes |
|---|---|---|---|---|
| `timeout_clone` | `vm.go:297` | dropped | — | Belongs to clonedvm |
| `timeout_create` | `vm.go:298` | done | — | Folded into `timeouts.create` |
| `timeout_migrate` | `vm.go:299` | planned | #19 | Folded into `timeouts.update` (per SDK comment that this is essentially a "timeout_update") |
| `timeout_reboot` | `vm.go:300` | dropped | — | Reboot is provider-internal (Q5); reuse `timeouts.update` |
| `timeout_shutdown_vm` | `vm.go:301` | planned | #6 | Internal to `power_state` transitions; not user-facing |
| `timeout_start_vm` | `vm.go:302` | planned | #6 | Internal to `power_state` transitions; not user-facing |
| `timeout_stop_vm` | `vm.go:303` | planned | #6 | Internal to `power_state` transitions; not user-facing |
| `timeout_move_disk` | `vm.go:304` | planned | #19 | Folded into migrate timeout |

### Open questions (pending maintainer decision)

| # | Question | Notes |
|---|---|---|
| OQ1 | Should `timeout_*` granular controls survive (as `timeouts.{create,update,delete,read}`) or do we adopt PVE-like granularity (separate per phase)? | Recommend: stick with framework `timeouts` block (ADR-006); fold internal phase timeouts into the parent op. |
| OQ2 | Read-only network attributes (`ipv4_addresses`, `ipv6_addresses`, `mac_addresses`, `network_interface_names`) — surface in resource as Computed, only in datasource, or both? | Recommend: only in datasource (Q3 deferral spirit — datasource exposes runtime state). Resource focuses on configuration. |
| OQ3 | Cloud-init `ip_config` — keep as ordered list (SDK) or convert to map keyed by interface name? | Recommend: map-keyed by interface name (`net0`, `net1`, …) — natural correspondence with `network_device` slots. |
| OQ4 | `agent.timeout` and `agent.wait_for_ip.{ipv4,ipv6}` — keep as PVE pass-through or fold into provider `timeouts.create` semantics? | Recommend: keep as PVE pass-through (they map to PVE config keys). |

### Summary by status

| Status | Count |
|---|---|
| done (already in `proxmox_vm`) | 14 |
| planned (Phase 2) | ~70 attributes + ~20 nested fields |
| dropped (out of scope or replaced) | 8 |
| open (maintainer decision) | 4 |

---

## Section 3 — Legacy test inventory

Every test in `proxmoxtf/resource/vm/**/*_test.go` and
`fwprovider/test/resource_vm_*.go` mapped to the user-visible behavior it
exercises and the Phase 2 PR that will own its port.

### Acceptance tests on the SDK resource (in `fwprovider/test/`)

Despite the directory name, these tests exercise `proxmox_virtual_environment_vm` (SDK), driven through the framework's test harness. They are the primary source of behavioral coverage to port.

#### `resource_vm_test.go` (2087 LOC, ~40 sub-cases)

| Test function | Behavior cluster | Target PR(s) for port |
|---|---|---|
| `TestAccResourceVM` | Description, name, node_name, protection, cpu update, memory update, vga update, watchdog, rng, virtiofs, purge/delete-disks/stop-on-destroy defaults + updates, hotplug variants, timeout persistence | Spread across #6 (memory, lifecycle), #13 (watchdog), #14 (hotplug), #17 (virtiofs); plus #3 (cpu/vga/rng update tests survive contract port) |
| `TestAccResourceVMImport` (`vm_test.go:715`) | Import round-trip with various attribute coverage | Carries forward into every PR — every block must include "import round-trip plan empty diff" per design Mandatory Test Scenarios |
| `TestAccResourceVMInitialization` (`vm_test.go:790`) | Cloud-init: custom + native; SCSI interface; username updates; upgrade flag variants; running-VM update | #11 (initialization) |
| `TestAccResourceVMNetwork` (`vm_test.go:1166`) | network_device interfaces, wait-for-IPv4, disconnected, removal (single + multi) | #10 (network_device) |
| `TestAccResourceVMClone` (`vm_test.go:1500`) | Clone scenarios | dropped — clone is `proxmox_cloned_vm`'s domain (D4) |
| `TestAccResourceVMVirtioSCSISingleWithAgent` (`vm_test.go:1804`) | SCSI single + agent enabled together | #7 (disk) + #13 (agent) |
| `TestAccResourceVMUpdateWhileStopped` (`vm_test.go:1880`) | Update operations while VM is stopped | #6 (power_state) — interaction with stopped state |

#### `resource_vm_disks_test.go` (2494 LOC)

| Test function | Behavior | Target PR |
|---|---|---|
| `TestAccResourceVMDisks` | Core disk CRUD (create/update/delete/import) across multiple interfaces | #7 |
| `TestAccResourceVMDiskCloneNFSResize` | NFS storage + clone resize | #7 (resource scope) / dropped clone parts |
| `TestAccResourceVMDiskRemovalReuseIssue2218` | Disk slot rename/reuse — regression for #2218 | #7 (mandatory map-keyed scenario) |
| `TestAccResourceVMDiskSpeedPerDisk` | Per-disk speed limits | #7 |
| `TestAccResourceVMDiskSpeedUpdate` | Update speed settings | #7 |
| `TestAccResourceVMDiskResizeWithOptionChange` | Resize + option change in one apply | #7 |
| `TestAccResourceVMDiskRemoval` | Plain removal scenario | #7 (covered by mandatory "remove middle slot") |
| `TestAccResourceVMDiskCDROMNotInDiskBlock` | CD-ROM excluded from disk block | #7 + #3 (cdrom keeps separate block) |
| `TestAccResourceVMDiskResizeNonHotpluggable` | Resize when hotplug disabled | #7 |
| `TestAccResourceVMDiskResizeDefaultHotplug` | Resize with default hotplug | #7 |
| `TestAccResourceVMEFIDiskStorageMigration` | EFI disk storage migration | #9 (efi_disk) + #19 (migrate) |

#### `resource_vm_hotplug_test.go` (1176 LOC)

| Test function | Behavior | Target PR |
|---|---|---|
| `TestAccResourceVMHotplug` | Hotplug attribute variants and behavior under update | #14 |

#### `resource_vm_pool_test.go` (562 LOC)

| Test function | Behavior | Target PR |
|---|---|---|
| `TestAccResourceVMPoolDetection` | Pool auto-detection | #18 |
| `TestAccResourceVMPoolDetectionLegacy` | Legacy pool detection | #18 |
| `TestAccResourceVMPoolDetectionManual` | Manual pool assignment | #18 |
| `TestAccResourceVMPoolMembership` | Pool membership lifecycle | #18 |
| `TestAccResourceVMPoolMembershipLegacy` | Legacy pool membership | #18 |
| `TestAccResourceVMPoolMembershipWithExplicitPoolID` | Explicit `pool_id` | #18 |

#### `resource_vm_reboot_after_creation_test.go` + `resource_vm_reboot_after_update_test.go` (606 LOC combined)

| Test function | Behavior | Target PR | Port status |
|---|---|---|---|
| `TestAccResourceVMRebootAfterCreationWithAgent` | Reboot policy after create | #6 (power_state) | **rewritten** — `reboot` user-facing attribute dropped (Q5); test becomes "after create, pending changes trigger reboot" |
| `TestAccResourceVMRebootAfterUpdateTPMStatePolicy` | TPM update triggers reboot | #6 + #9 | rewritten same way |
| `TestAccResourceVMRebootAfterUpdateCloudInitMovePolicy` | Cloud-init move triggers reboot | #6 + #11 | rewritten |
| `TestAccResourceVMRebootAfterUpdateTemplatePolicy` | Template flag change | #6 | rewritten |
| `TestAccResourceVMRebootAfterUpdateDiskMovePolicy` | Disk move triggers reboot | #6 + #7 | rewritten |
| `TestAccResourceVMRebootAfterUpdateDiskResizePolicy` | Disk resize triggers reboot | #6 + #7 | rewritten |

These six are explicitly **rewritten** (not ported byte-level) — the user-facing `reboot` attribute is gone; tests assert that pending changes drive provider-internal reboots.

#### `resource_vm_template_test.go`

| Test function | Behavior | Target PR |
|---|---|---|
| `TestAccResourceVMTemplateConversion` | Convert VM to template | done — already in `proxmox_vm` Create + Update paths |

#### `resource_vm_tpm_state_test.go`

| Test function | Behavior | Target PR |
|---|---|---|
| `TestAccResourceVMTpmState` | TPM state lifecycle | #9 (tpm_state) |

#### `resource_vm_cdrom_test.go`

| Test function | Behavior | Target PR |
|---|---|---|
| `TestAccResourceVMCDROM` | CD-ROM lifecycle | done — `cdrom/` package already has acceptance tests |

#### Datasource tests

| Test function | Behavior | Target PR |
|---|---|---|
| `TestAccDatasourceSDKVMNotFound` (`fwprovider/test/datasource_vm_test.go`) | SDK-side datasource not-found | (SDK-only — unaffected by #1231) |
| FW datasource coverage | Equivalent functional coverage required per ADR-006 | every Phase 2 PR adds coverage — datasource gets the same map-keyed blocks (per Q1 + design Datasource Parity table) |

### Unit tests on the SDK package (in `proxmoxtf/resource/vm/`)

| Test function | Source | Behavior | Target PR | Port status |
|---|---|---|---|---|
| `TestVMInstantiation` | `vm_test.go:21` | Schema instantiation smoke test | — | dropped — Framework provider has its own schema-validation patterns |
| `TestVMSchema` | `vm_test.go:31` | Schema field-by-field assertions | — | dropped — Framework schema is the source of truth |
| `TestHotplugContains` | `vm_test.go:421` | `hotplug` flag parsing | #14 | port to FW unit test |
| `Test_parseImportIDWIthNodeName` | `vm_test.go:457` | Import ID parser | #4 | port — already done equivalent at `fwprovider/nodes/vm/resource.go:397` |
| `TestCPUType` | `validators_test.go:16` | Validates CPU type enum | — | **dropped** — long enum validator dropped per ADR-004 (Q4/F27) |
| `TestMachineType` | `validators_test.go:47` | Validates machine type enum | — | **dropped** — long enum validator dropped per ADR-004 |
| `TestVmHostname` | `validators_test.go:82` | DNS name validator | #4/#5 | port — `name` validator survives (kept per Section 5 inventory) |
| `TestDiskOrderingDeterministic` | `disk/disk_test.go:27` | Map ordering for disk slots | — | **dropped** — map-keyed pattern eliminates ordering concern |
| `TestDiskOrderingVariousInterfaces` | `disk/disk_test.go:112` | Cross-interface ordering | — | **dropped** — same |
| `TestDiskDevicesEqual` | `disk/disk_test.go:195` | `CustomStorageDevice.Equals` | #7 | port if FW disk implementation needs equality helper; otherwise drop |
| `TestDiskUpdateSkipsUnchangedDisks` | `disk/disk_test.go:268` | No-op update for unchanged | #7 | port — covered by mandatory map-keyed scenario "Apply, re-plan with same config — assert empty diff" |
| `TestImportFromDiskNotReimportedOnSizeChange` | `disk/disk_test.go:369` | Import-from semantics on resize | #7 | port — specific behavior |
| `TestDiskDeletionDetectionInGetDiskDeviceObjects` | `disk/disk_test.go:448` | Detection during read | #7 | port |
| `TestDiskDeletionWithBootDiskProtection` | `disk/disk_test.go:638` | Boot disk protection during deletion | #7 | port |
| `TestOriginalBugScenario` | `disk/disk_test.go:716` | Regression for original bug | #7 | port |
| `TestDiskSpeedSettingsPerDisk` | `disk/disk_test.go:817` | Speed setting per disk | #7 | port |
| `TestVMSchema` (disk subpkg) | `disk/schema_test.go:11` | Disk schema assertions | — | dropped — FW schema is source of truth |
| `TestNetworkSchema` | `network/schema_test.go:11` | Network schema assertions | — | dropped — FW schema is source of truth |

### Summary by port action

| Class | Count | Port action |
|---|---|---|
| Acceptance tests to port (Phase 2) | ~30 | port; carry forward as `TestAccResource...` in `fwprovider/nodes/vm/` |
| Acceptance tests to **rewrite** (reboot semantics) | 6 | rewritten per Q5 — provider-driven reboots, not user-facing |
| Acceptance tests already done | 3 | `TestAccResourceVMTemplateConversion`, `TestAccResourceVMCDROM`, `TestAccResourceVMShort` |
| Acceptance tests dropped (out of scope) | 1 | `TestAccResourceVMClone` (clonedvm domain) |
| SDK unit tests to port | ~7 | hotplug parsing, hostname validator, disk-specific behavior |
| SDK unit tests to **drop** | ~7 | Long-enum validators, schema-instantiation tests, ordering tests (map-keyed eliminates) |

---

## Section 4 — Per-attribute classification (ADR-004 amendment)

Every existing `Optional+Computed` attribute in the current `proxmox_vm` (and
its sub-packages) reclassified per the new ADR-004 PVE-defaults rule.

> **Rule recap (ADR-004 amendment, drafted in PR #2):**
>
> | PVE Read behavior | Schema target |
> |---|---|
> | Auto-populates default value | `Optional + Computed` (no Default) |
> | Returns null/absent when unset | `Optional` only |
> | Provider-only attribute (no PVE counterpart) | `Optional + Default` |

### Methodology (Section 4)

Three data sources reconciled (a single test wasn't enough — see "Methodology limitations" below):

1. **Empirical mitmproxy traces** (this audit, 2026-04-19):
   - **Pass 1**: `TestAccResourceVMShort` (`fwprovider/nodes/vm/resource_test.go`) — minimal VM with only top-level scalars. 25 GET responses, log `/tmp/api_debug.log`.
   - **Pass 2**: `TestAccResourceVM2CPU` (`fwprovider/nodes/vm/cpu/resource_test.go`) — VMs with explicit `cpu.*` fields set across 8 sub-cases. 33 GET responses, log `/tmp/api_debug_cpu.log`. **This pass surfaced the cores/sockets auto-populate carve-out.**
   - **Pass 3**: `TestAccResourceVM2VGA`, `TestAccResourceVM2RNG`, `TestAccResourceVM2CDROM` together — VMs with explicit vga, rng, cdrom blocks. 33 GET responses, log `/tmp/api_debug_blocks.log`.
2. **PVE source** — `qemu-server.git src/PVE/QemuServer.pm` `$confdesc` hashref documents internal defaults per field.
3. **Existing provider sentinels** — `cpu/resource.go:38-60` substitutes `1` for nil `Cores`/`Sockets` and `kvm64` for nil `CPUEmulation.Type`. The existence of these sentinels and their purpose ("PVE does not return actual value for cores VM, etc is using default") is itself evidence about PVE behavior — but only partial evidence (see cpu carve-out below).

### Key empirical findings

**Finding 1 — minimal VM (no sub-blocks set).** PVE returns ONLY:

```json
{
    "boot": " ",
    "smbios1": "uuid=<random-uuid>",
    "vmgenid": "<random-guid>",
    "meta": "creation-qemu=...,ctime=...",
    "digest": "<config-digest>"
}
```

Every other field is absent.

**Finding 2 — cpu carve-out.** When the user sets ANY `cpu.*` field, PVE auto-populates `cores=1` and `sockets=1` in the GET response, even if the user didn't set those. Direct evidence from VM 102 in `/tmp/api_debug_cpu.log`:

```text
config (user set cpu.type="x86-64-v4" only):
{ "cpu": "x86-64-v4", "sockets": 1, "cores": 1, ... }
```

VM 103 with `cpu.limit=64` only also returned `cores: 1` (no sockets in this case — auto-populate behavior depends on which field triggered it). The existing provider sentinels (cpu/resource.go:38-60) handle BOTH the "block absent" case AND the "block has fields" case — they were defensive against this carve-out.

**Finding 3 — vga/rng/cdrom: no auto-populate.** Direct evidence from `/tmp/api_debug_blocks.log`:

| Block | Test config | PVE response |
|---|---|---|
| vga | `vga = { type = "std" }` | `"vga": "type=std"` (only what user set) |
| vga | `vga = { type = "qxl", clipboard = "vnc" }` | `"vga": "clipboard=vnc,type=qxl"` |
| rng | `rng = { source = "/dev/urandom" }` | `"rng0": "source=/dev/urandom"` (only what user set) |
| rng | `rng = { source, period }` | `"rng0": "source=...,period=1000"` |
| cdrom | `cdrom = { ide2 = { file_id = "cdrom" } }` | `"ide2": "cdrom,media=cdrom"` (PVE adds the implicit `media=cdrom` qualifier; `file_id` is always present when slot exists) |

When the user sets nothing in vga/rng/cdrom, PVE returns nothing for those blocks.

**Finding 4 — memory: not directly tested.** No `TestAccResourceVM2Memory` exists. memory is shared with `clonedvm` only. Pattern likely follows vga/rng (no auto-populate) since memory's qemu-server `$confdesc` shows `memory: default=none` and `balloon: default=none`. Verify in PR #6 when memory is wired into `proxmox_vm`.

**Finding 5 — PVE Perl source.** Cross-reference with `qemu-server.git src/PVE/QemuServer.pm`:

- `cores`: default=1 (matches Finding 2 auto-populate)
- `sockets`: default=1 (matches Finding 2)
- `cpu` (property string): default=`type=kvm64` — but **NOT auto-populated** in our traces (PVE returned the user-set value or nothing; never substituted kvm64). The provider's sentinel for `Type → "kvm64"` is therefore **not corroborated by PVE behavior**; it appears to be a provider invention.
- `vga`/`rng`: defaults documented in source but NOT surfaced in GET (matches Finding 3).

**Implication:** The wholesale "drop Computed from all sub-block attributes" rule needs a carve-out: `cpu.cores` and `cpu.sockets` keep `Optional+Computed` because PVE actively auto-populates them. Everything else (vga, rng, cdrom-block-level, memory predicted, all other cpu fields) drops Computed.

### `cpu` attributes (`fwprovider/nodes/vm/cpu/resource_schema.go`)

| Attribute | Current | PVE Read | qemu-server default | Target schema | Target PR |
|---|---|---|---|---|---|
| `cpu` (block) | Optional+Computed | absent when block has no fields; populated when block has any field (cores/sockets auto-added) | n/a | **Optional** (drop Computed at block level — when user provides nothing, PVE returns nothing) | #3 |
| `cpu.affinity` | Optional+Computed | absent (only present when set) | none | **Optional** | #3 |
| `cpu.architecture` | Optional+Computed | absent | none (root@pam only) | **Optional** | #3 |
| `cpu.cores` | Optional+Computed | **AUTO-POPULATED to 1** when block has any field | 1 | **Optional+Computed (KEEP)** — preserve current sentinel behavior, but ensure `NewValue` reads PVE's actual returned value | #3 |
| `cpu.flags` | Optional+Computed | absent | none | **Optional** | #3 |
| `cpu.limit` | Optional+Computed | absent | 0 | **Optional** | #3 |
| `cpu.numa` (bool) | Optional+Computed | absent | 0 | dropped (rehomed `numa.enabled` per P3) | #3 / #13 |
| `cpu.sockets` | Optional+Computed | **AUTO-POPULATED to 1** when block has any field | 1 | **Optional+Computed (KEEP)** — same as cores | #3 |
| `cpu.type` | Optional+Computed | absent when not set (provider sentinel was wrong — see Finding 5) | kvm64 (qemu-server default; not surfaced by PVE) | **Optional** — drop the `Type→"kvm64"` sentinel; provider was over-reaching | #3 |
| `cpu.units` | Optional+Computed | absent | 1024 (cgroup v1) / 100 (cgroup v2) | **Optional** | #3 |
| `cpu.hotplugged` (vcpus) | Optional+Computed | absent | 0 | dropped (rehomed `vcpus` per P3) | #3 / #14 |

**cpu carve-out summary**: 2 attributes keep `Optional+Computed` (`cores`, `sockets`); 7 drop to `Optional`; 2 dropped via rehoming.

### `vga` attributes (`fwprovider/nodes/vm/vga/resource_schema.go`)

| Attribute | Current | PVE Read | qemu-server default | Target schema | Target PR |
|---|---|---|---|---|---|
| `vga` (block) | Optional+Computed (with `UseStateForUnknown`) | absent | std (type only) | **Optional** (drop UseStateForUnknown planmodifier too) | #3 |
| `vga.clipboard` | Optional+Computed | absent | none | **Optional** | #3 |
| `vga.type` | Optional+Computed | absent | std | **Optional** | #3 |
| `vga.memory` | Optional+Computed | absent | none | **Optional** | #3 |

### `rng` attributes (`fwprovider/nodes/vm/rng/resource_schema.go`)

| Attribute | Current | PVE Read | qemu-server default | Target schema | Target PR |
|---|---|---|---|---|---|
| `rng` (block) | Optional+Computed (with `UseStateForUnknown`) | absent | none (root@pam only) | **Optional** (drop UseStateForUnknown planmodifier too) | #3 |
| `rng.source` | Optional+Computed | absent | none | **Optional** | #3 |
| `rng.max_bytes` | Optional+Computed | absent | none | **Optional** | #3 |
| `rng.period` | Optional+Computed | absent | none | **Optional** | #3 |

### `memory` attributes (`fwprovider/nodes/vm/memory/resource_schema.go`)

| Attribute | Current | PVE Read | qemu-server default | Target schema | Target PR |
|---|---|---|---|---|---|
| `memory` (block) | Optional+Computed | absent | n/a | **Optional** | #3 |
| `memory.size` | Optional+Computed+`Default(512)` | absent | none (PVE applies 512 at launch only) | **Optional** (drop Default per F39) | #3 |
| `memory.balloon` | Optional+Computed+`Default(0)` | absent | none | **Optional** (drop Default) | #3 |
| `memory.shares` | Optional+Computed+`Default(1000)` | absent | 1000 (PVE source — but absent in GET) | **Optional** (drop Default) | #3 |
| `memory.hugepages` | Optional+Computed | absent | none | **Optional** | #3 |
| `memory.keep_hugepages` | Optional+Computed | absent | 0 | **Optional** | #3 |

### `cdrom` attributes (`fwprovider/nodes/vm/cdrom/resource_schema.go`)

| Attribute | Current | PVE Read | qemu-server default | Target schema | Target PR |
|---|---|---|---|---|---|
| `cdrom` (map-level) | Optional+Computed | absent (no IDE devices auto-attached) | n/a | **Optional** (drop Computed) | #3 |
| `cdrom[slot].file_id` | Optional+Computed+`Default("cdrom")` | varies (when device set, returns FileVolume; when unset, slot is absent entirely) | n/a (per-slot value is the storage:path string) | Optional+Computed (kept — per-slot block always has a `file_id` value when the slot exists; the block-level Optional handles "no slot") | #3 |

### Top-level scalars (already in `proxmox_vm`)

| Attribute | Current | PVE Read | qemu-server default | Target schema | Notes |
|---|---|---|---|---|---|
| `description` | Optional | absent when unset | none | confirmed `Optional` | No change |
| `name` | Optional (with DNS validator) | absent when unset | none | confirmed `Optional` | No change |
| `tags` | stringset (Optional) | absent when unset | none | confirmed `Optional` | No change |
| `template` | Optional (RequiresReplace) | absent when 0 (false), present when 1 | 0 | confirmed `Optional` | No change. Note: PVE returns `template=1` only when set; this is a "presence as truthiness" pattern. |
| `id` (SDK `vmid`) | Computed+Optional+RequiresReplace | always present (it's the path key) | n/a | confirmed | No change |
| `node_name` | Required | n/a | n/a | confirmed `Required` | No change |
| `stop_on_destroy` | Optional+Computed+Default=false | n/a (provider-only) | n/a | confirmed `Optional+Default` | No change (provider-only attribute per ADR-004) |
| `purge_on_destroy` | Optional+Computed+Default=true | n/a (provider-only) | n/a | confirmed `Optional+Default` | No change |
| `delete_unreferenced_disks_on_destroy` | Optional+Computed+Default=true | n/a (provider-only) | n/a | confirmed `Optional+Default` | No change |

### Future fields with auto-population behavior (Phase 2 PRs to verify when implemented)

These will likely **keep** `Optional+Computed` (without provider Default) per the ADR-004 rule. Mitmproxy verification at PR-time:

| Attribute | qemu-server default | PVE Read prediction | Target schema (predicted) | Target PR |
|---|---|---|---|---|
| `boot_order` | cdn (legacy), nested order=none | always present (`boot=" "` or `boot="order=..."`) | Optional+Computed | #8 |
| `smbios.uuid` | autogenerated UUID | always present in `smbios1` | Optional+Computed | #12 |
| `vmgenid` (if surfaced) | autogenerated | always present | Optional+Computed (or hidden) | maintainer decision |
| `acpi` | 1 | likely absent unless changed (verify) | TBD — predict Optional only | #14 |
| `tablet` | 1 | likely absent unless changed (verify) | TBD — predict Optional only | #14 |
| `kvm` | 1 | likely absent unless changed (verify) | TBD — predict Optional only | (out of scope today) |
| `bios` | seabios | likely absent unless changed (verify) | TBD — predict Optional only | #8 |
| `scsihw` | lsi | likely absent unless changed (verify) | TBD — predict Optional only | #9 |
| `hotplug` | network,disk,usb | likely returns the list always (verify) | TBD — predict Optional+Computed | #14 |
| `protection` | 0 | likely absent unless changed | TBD — predict Optional only | #18 |
| `onboot` | 0 | likely absent unless changed | predict Optional only | #6 |

### Cross-validation with PVE source

The empirical and source-code data are consistent: PVE Perl's `$confdesc` documents internal defaults that are applied at QEMU launch time, but the Perl Web API does not write those defaults back to the on-disk config. So `parse_vm_config()` returns only what's literally in the config file. The mitmproxy GET responses confirm this: only user-set fields and PVE-auto-generated fields (`smbios1`, `vmgenid`, `boot`, `meta`, `digest`) appear.

### Summary by classification action

| Action | Attribute count | PRs |
|---|---|---|
| Drop `Computed` (Optional+Computed → Optional) | 21 in existing sub-packages | #3 |
| Keep `Optional+Computed` (PVE auto-populates) | 2 (cpu.cores, cpu.sockets) + cdrom[slot].file_id | #3 |
| Drop provider `Default` (per F39) | 3 (memory.size, memory.balloon, memory.shares) | #3 |
| Drop `UseStateForUnknown` planmodifier (consequence of dropping Computed) | 2 (vga, rng blocks) | #3 |
| Drop attribute (rehome) | 2 (cpu.numa, cpu.hotplugged) | #3 / #13 / #14 |
| Confirmed no change (already correct) | 9 top-level scalars | — |
| Predicted but verify in Phase 2 | 11 future fields | #6, #8, #9, #12, #14, #18 |

### Methodology limitations

Three caveats on Section 4's confidence:

1. **`memory/` not directly tested** — no `TestAccResourceVM2Memory` exists (memory is only consumed by `clonedvm` today). The Section 4 classification for memory attributes is *predicted* from the vga/rng pattern (no auto-populate) and the qemu-server source. PR #6 (when memory is wired into `proxmox_vm`) must re-verify with mitmproxy and adjust if PVE auto-populates anything.
2. **"Set then unset" path not explicitly tested** — the mitmproxy traces cover "never set" and "set" but not the "set then unset" transition. PVE's config file model (key=value pairs in `/etc/pve/qemu-server/{vmid}.conf`) makes "absent" and "unset" equivalent at the storage layer, so the conclusion is sound, but a dedicated test in PR #3 should cover the round-trip.
3. **Future-field predictions are predictions** — the table for fields not yet in the codebase (acpi, tablet, bios, etc.) is reasoning from `$confdesc` defaults, not direct observation. Each Phase 2 PR that adds these fields must re-verify with mitmproxy at PR-time.

### Mitmproxy session details

| Pass | Date | Test(s) | Captures | Log |
|---|---|---|---|---|
| 1 | 2026-04-19 | `TestAccResourceVMShort` (9 sub-cases, all PASS) | 25 GET /config | `/tmp/api_debug.log` (3473 lines) |
| 2 | 2026-04-19 | `TestAccResourceVM2CPU` (8 sub-cases, all PASS) | 33 GET /config | `/tmp/api_debug_cpu.log` (4062 lines) |
| 3 | 2026-04-19 | `TestAccResourceVM2VGA` + `TestAccResourceVM2RNG` + `TestAccResourceVM2CDROM` (combined) | 33 GET /config | `/tmp/api_debug_blocks.log` (5979 lines) |

| Detail | Value |
|---|---|
| Proxy mode | `mitmdump --mode regular@8082 --flow-detail 4` (port 8082 — Docker holds 8080) |
| Cluster | PVE 10.1.2 at `pve.bpghome.net:8006` (per `meta: "creation-qemu=10.1.2"`) |
| Auth | `terraform@pve!provider` API token |
| Reproducibility | Re-run the same tests with `HTTP_PROXY=http://127.0.0.1:8082 HTTPS_PROXY=... PROXMOX_VE_INSECURE=true ./testacc <TestName>` |

---

## Section 5 — Validator inventory (ADR-004 enum rule)

Every validator currently in use, classified per the ADR-004 amendment enum
rule:

> **Rule recap:** Use `OneOf` for short, stable PVE enums (≤5 values, unlikely
> to extend). For long or version-evolving enums (CPU types, VGA types,
> machine types, BIOS modes, scsi_hardware), defer to PVE; use a regex
> validator only if format-only validation is meaningful.

### Top-level (`resource_schema.go`)

| Attribute | Validator | Type | Decision | Reason | Target PR |
|---|---|---|---|---|---|
| `name` | `stringvalidator.RegexMatches` DNS regex (`resource_schema.go:64`) | format | keep | DNS format check is format-only validation, meaningful and stable | — |

### `cpu/`

| Attribute | Validator | Type | Decision | Reason | Target PR |
|---|---|---|---|---|---|
| `affinity` | `stringvalidator.RegexMatches ^\d+[\d-,]*$` (`cpu/resource_schema.go:38`) | format | keep | Format-only validation is meaningful (catches typos before PVE call) | — |
| `architecture` | `stringvalidator.OneOf("aarch64", "x86_64")` (`cpu/resource_schema.go:50`) | short enum (2) | keep | Short, stable, no growth pressure | — |
| `cores` | `int64validator.Between(1, 1024)` (`cpu/resource_schema.go:59`) | range | keep | PVE-source bound | — |
| `flags` | `setvalidator.AlsoRequires(parent.type)` (`cpu/resource_schema.go:73`) | cross-attribute | keep | Frequently hit; compound `CPUEmulation` requires `Type` when `Flags` set (see F26) | — |
| `flags` (elements) | `stringvalidator.RegexMatches (.\|\s)*\S(.\|\s)*` + `LengthAtLeast(1)` (`cpu/resource_schema.go:75–79`) | format | keep | Non-empty/non-whitespace check | — |
| `hotplugged` | `int64validator.Between(1, 1024)` (`cpu/resource_schema.go:89`) | range | **drop** | Attribute itself being rehomed to top-level `vcpus` (P3); will get its own validator there | #3 (drop with attr) / #14 (replace) |
| `limit` | `float64validator.Between(0, 128)` (`cpu/resource_schema.go:98`) | range | keep | PVE-source bound | — |
| `sockets` | `int64validator.Between(1, 16)` (`cpu/resource_schema.go:113`) | range | keep | PVE-source bound | — |
| `type` | `stringvalidator.OneOf(...75 CPU types...)` (`cpu/resource_schema.go:125–204`) | **long enum** | **drop** | (Confirms F27/P5) Long, version-evolving — defer to PVE per ADR-004 enum rule | #3 |
| `units` | `int64validator.Between(1, 262144)` (`cpu/resource_schema.go:214`) | range | open question | PVE source: cgroup v2 default is `100`, with `0` valid (disables limiting). Current `Between(1, 262144)` rejects `0`. Acceptable on cgroup v1 hosts; restrictive on cgroup v2. PR #3: relax to `AtLeast(0)` or document the cgroup v1 assumption. | #3 |

### `vga/`

| Attribute | Validator | Type | Decision | Reason | Target PR |
|---|---|---|---|---|---|
| `clipboard` | `stringvalidator.OneOf("vnc")` (`vga/resource_schema.go:47`) | enum (1) | keep | Single accepted value today; revisit if PVE adds | — |
| `type` | `stringvalidator.OneOf(...14 VGA types...)` (`vga/resource_schema.go:56–72`) | **long enum** | **drop** | (Confirms F33) Version-evolving (PVE adds qxl variants over time) | #3 |
| `memory` | `int64validator.Between(4, 512)` (`vga/resource_schema.go:80`) | range | keep | PVE-source bound | — |

### `rng/`

| Attribute | Validator | Type | Decision | Reason | Target PR |
|---|---|---|---|---|---|
| `source` | `stringvalidator.LengthAtLeast(1)` (`rng/resource_schema.go:45`) | format | keep | Trivial non-empty check | — |
| `max_bytes` | `int64validator.AtLeast(0)` (`rng/resource_schema.go:55`) | range | keep | Trivial bound; intersects F37 (int-zero trap is in `FillCreateBody`, not validator) | — |
| `period` | `int64validator.AtLeast(0)` (`rng/resource_schema.go:65`) | range | keep | Trivial bound; same int-zero trap context | — |

### `memory/`

| Attribute | Validator | Type | Decision | Reason | Target PR |
|---|---|---|---|---|---|
| `size` | `int64validator.Between(64, 268435456)` (`memory/resource_schema.go:54`) | range | keep | PVE-source bound | — |
| `balloon` | `int64validator.Between(0, 268435456)` (`memory/resource_schema.go:68`) | range | keep | PVE-source bound | — |
| `shares` | `int64validator.Between(0, 50000)` (`memory/resource_schema.go:81`) | range | keep | PVE-source bound | — |
| `hugepages` | `stringvalidator.OneOf("2", "1024", "any")` (`memory/resource_schema.go:95`) | short enum (3) | keep | Short, stable | — |

### `cdrom/`

| Attribute | Validator | Type | Decision | Reason | Target PR |
|---|---|---|---|---|---|
| (map keys) | `mapvalidator.KeysAre + stringvalidator.RegexMatches ^(ide[0-3]\|sata[0-5]\|scsi([0-9]\|1[0-3]))$` (`cdrom/resource_schema.go:31–37`) | slot regex | **tighten** | (Confirms F46) `scsi` only goes to 13; PVE bound is `MAX_SCSI_DISKS=31`. Update to design slot-regex table value `scsi([0-9]\|[12][0-9]\|30)`. | #3 |
| `file_id` | `stringvalidator.Any(OneOf("cdrom", "none"), validators.FileID())` (`cdrom/resource_schema.go:48–49`) | short enum + format | keep | Sentinel values + file ID format check | — |

### Summary by decision

| Decision | Count | Rationale |
|---|---|---|
| keep | 18 | Short stable enums, cross-attribute, range bounds, format checks |
| drop | 2 | Long version-evolving enums (`cpu.type`, `vga.type`) |
| tighten | 1 | `cdrom` slot regex (`scsi` upper bound) |
| open question | 1 | `cpu.units` lower-bound (cgroup v1 vs v2) |
| drop-with-attr | 1 | `cpu.hotplugged` validator dies with the attribute (rehome) |

All target PR #3 (sub-package port).

---

## Section 6 — Q5 resolution: `power_state` redesign

| Aspect | Resolution |
|---|---|
| Drop `started` (boolean) | Yes — replaced by `power_state` |
| Add `power_state` (string) | Values: `"running"`, `"stopped"`. Default `"running"`. |
| Add Computed `status` | For runtime drift visibility |
| Drop user-facing `reboot` | Provider decides reboot-vs-restart from pending changes |
| Keep `on_boot` (boolean) | Yes — corresponds to PVE "Start at boot" |
| Implementation PR | PR #6 |
| Audit deliverable | This section + an entry in Section 2 (capabilities) |

Implementation notes for PR #6:

### API call sequence (Create / Update reaching desired `power_state`)

1. **After CRUD config write** (Create or Update path), call `GetVMStatus` (`vmAPI.GetVMStatus(ctx)`).
2. Compare `current.Status` vs `plan.power_state.ValueString()`.
3. **`plan="running"`, `current="stopped"`**: dispatch `StartVM(ctx, startTimeoutSec)` then `WaitForVMStatus(ctx, "running")`. Mirrors SDK `vmStart` (`proxmoxtf/resource/vm/vm.go:1980`).
4. **`plan="stopped"`, `current="running"`**: prefer graceful shutdown if QEMU guest agent is enabled — `ShutdownVM(ctx, &vms.ShutdownRequestBody{ForceStop, Timeout})` then `WaitForVMStatus(ctx, "stopped")`. Fall back to forceful `StopVM(ctx)` + `WaitForVMStatus(ctx, "stopped")` if the agent is not available. Mirrors SDK `EnsureStopped` (`proxmoxtf/resource/vm/vm.go:2043`).
5. **Equal**: no-op.

Use the existing `proxmox/retry` operation types per ADR-005:

- `retry.NewTaskOperation` for `StartVM` (returns UPID async task)
- `retry.NewAPICallOperation` for `ShutdownVM` (synchronous)
- `retry.NewPollOperation` for `WaitForVMStatus`

### Reboot-detection heuristic

After a successful `UpdateVM`, **the provider** decides whether a reboot is required:

1. Re-fetch config via `GetVM(ctx)` (already required by ADR-005 read-back rule).
2. Inspect `vmConfig.PendingChanges` (the `pending` field PVE returns when applied changes require a reboot to take effect).
3. If `pending != nil && len(*pending) > 0 && plan.power_state == "running"` → reboot.
4. Reboot = stop + start (reuse the EnsureStopped/EnsureRunning pattern). Skip if `power_state == "stopped"` (the change will apply on next start).

**No user-facing `reboot` attribute.** Confirmed by Q5 resolution. Replaces SDK's `mkRebootAfterCreation` / `mkRebootAfterUpdate` controls.

### Interaction with `stop_on_destroy`

- `stop_on_destroy` continues to mean "during Delete: skip graceful shutdown, force-stop instead".
- It is **independent** of `power_state`. A VM may have `power_state="stopped"` (already stopped) and `stop_on_destroy=true` (the latter has no effect on already-stopped VMs).
- A VM with `power_state="running"` and `stop_on_destroy=false` is shut down gracefully on destroy, matching today's behavior.

### Computed `status` attribute

- Read-only mirror of PVE runtime `Status` field (`running`, `stopped`, `paused`, `prelaunch`, `migrating`, `unknown`).
- Distinct from `power_state` (desired state, user-controlled).
- `status` drift is informational only; provider does not auto-correct on Read (consistent with how today's drift detection works).

### `on_boot` semantics

- Maps directly to PVE `onboot` (start at PVE host boot).
- No interaction with `power_state` (one is "should the VM be running now", the other is "should it autostart later").
- `Optional` only (per ADR-004 amendment — PVE returns null when not set).

---

## Section 7 — Shared-types catalog (R3 mitigation)

Every type in `proxmox/nodes/vms` consumed by the SDK resource. Reviewers
of any client-touching PR check this table to flag cascade risk.

> **Allowed-changes rules (from R3):** add freely; internal cleanup freely;
> rename public fields with same-PR SDK callsite updates; remove public
> fields forbidden until post-Phase-2.

### Public types (consumed by both SDK and Framework — high cascade risk)

> Renaming a field on these types requires updating both `proxmoxtf/resource/vm/` and `fwprovider/nodes/{vm,clonedvm}/` callsites in the same PR.

| Type | SDK callsites | FW VM callsites | FW ClonedVM callsites | Notes |
|---|---|---|---|---|
| `vms.GetResponseData` | `vm.go:1936, 1937, 1950`, `disk/disk.go:34`, `network/network.go:109, 123, 134, 190` | `model.go:120` (read), every sub-package `NewValue` | `clonedvm/resource.go:508, 677, 791, 846, 902, 971` | API Read shape — touched everywhere |
| `vms.UpdateRequestBody` | `disk/disk.go:71, 94, 599`, `vm.go:1962` | `resource.go:252`, every sub-package `FillUpdateBody` | `clonedvm/resource.go:509, 587, 678` | API Update shape — touched everywhere |
| `vms.Client` | `disk/disk.go:63`, `network/network.go:188`, `vm.go:1980, 1998, 2019, 2043, 2076, 2122, 2190, 2246` | `resource.go:166, 250, 348, 443, 471` | `clonedvm/resource.go:508, 791, 1013, 1041` | RPC surface (`vms.Client`) — method renames cascade |
| `vms.CustomStorageDevice` / `vms.CustomStorageDevices` | `disk/disk.go:34, 61, 62, 159, 168, 171, 281, 302, 325, 374, 597–607`, `disk/disk_test.go` (>20 callsites), `vm.go:1937, 1950` | `cdrom/model.go:32–39`, `cdrom/resource.go:25` | `clonedvm/resource.go:688, 695, 902` | Storage shape — disk PR (#7) will touch heavily |
| `vms.CustomNetworkDevice` / `vms.CustomNetworkDevices` | `network/network.go:24, 26, 43, 109, 110, 123, 134` | (none yet — added in #10) | `clonedvm/resource.go:606, 611, 846, 971` | Network shape — network PR (#10) and ClonedVM share |
| `vms.ShutdownRequestBody` | `vm.go:2007` | `resource.go:452` | `clonedvm/resource.go:1022` | Shutdown shape — small surface |
| `vms.CloneRequestBody` | (clone is in clonedvm only on FW side) | `concurrency_test.go:88` (test) | `clonedvm/resource.go:463` | Clone shape |
| `vms.CreateRequestBody` | (none — SDK uses different patterns) | `resource.go:148`, sub-package `FillCreateBody` | (none — clones don't use Create) | Framework-only — safe to rename |

### Framework-only types (low cascade risk)

| Type | FW callsites | SDK callsites | Notes |
|---|---|---|---|
| `vms.CustomCPUEmulation` | `cpu/resource.go:123, 229, 253` | (none) | FW-only — safe to rename in #3 |
| `vms.CustomVGADevice` | `vga/resource.go:57, 72, 101, 131` | (none) | FW-only — safe to rename in #3 |
| `vms.CustomRNGDevice` | `rng/resource.go:48, 86, 143` | (none) | FW-only — safe to rename in #3 |

### Public type aliases / constants

| Identifier | File:line | Consumers | Notes |
|---|---|---|---|
| `vms.StorageInterfaces` | `proxmox/nodes/vms/...` (slice) | `proxmoxtf/resource/vm/disk/disk.go:281, 284` | SDK-only consumer; FW uses regex-based slot validation per ADR-008 |
| `vms.MoveDiskRequestBody` | `proxmox/nodes/vms/...` | `proxmoxtf/resource/vm/disk/disk.go:113` | SDK-only today; will become FW-shared when #19 (migrate) lands |
| `vms.ResizeDiskRequestBody` | `proxmox/nodes/vms/...` | `proxmoxtf/resource/vm/disk/disk.go:128, 362`, `clonedvm/resource.go:680, 685, 737` | Shared between SDK and ClonedVM today; #7 (disk MVP) will add FW VM consumer |
| `vms.WaitForIPConfig` | `proxmox/nodes/vms/...` | `proxmoxtf/resource/vm/network/network.go:192` | SDK-only today; #10 (network) likely adds FW consumer |

### Already-confirmed sensitive shapes (carried forward from prior FWK audits)

| Type | Field | Constraint |
|---|---|---|
| `*vms.GetResponseData.EFIDisk` | `*CustomEFIDisk`, json `efidisk0,omitempty` | Single-instance pointer — preserves D5 architectural-single decision |
| `*vms.GetResponseData.TPMState` | `*CustomTPMState`, json `tpmstate0,omitempty` | Single-instance pointer — same |
| `*vms.UpdateRequestBody.ToDelete(string) error` | method | Used by `attribute.CheckDelete`; do not rename without ADR-008 update |

---

## Gaps surfaced for the gap matrix

> Anything found here that needs Phase-2 tracking is mirrored into
> `1231_GAP_MATRIX.md`. The audit (this file) is frozen after PR #1; the gap
> matrix lives through Phase 2 and becomes PR #20's parity report.

---

## Pre-resolved findings (carried from grilling pass)

These are the findings that the design pass already locked in — the audit
confirms but does not re-litigate:

| # | Finding | Severity | Target PR | Source |
|---|---|---|---|---|
| P1 | `IsDefined(plan.Sockets)` copy-paste bug in cpu Limit branch (`fwprovider/nodes/vm/cpu/resource.go:190`) | should-fix | #3 | Direct read of file during grilling |
| P2 | Nil-substitution sentinel `cpu.Cores == nil → 1` violates ADR-004 PVE-defaults rule | should-fix | #3 | Read of `cpu/resource.go:38–42` |
| P3 | `cpu.numa` and `cpu.hotplugged` are SDK-inherited misnomers; should rehome | should-fix | #3 (drop) + #13 (rehome `numa.enabled`) + #14 (rehome `vcpus`) | Design D7 + ADR-008 single-vs-map rule |
| P4 | Hand-rolled `ShouldBeRemoved` + `IsDefined` cascades in cpu/, vga/, rng/, memory/ instead of `attribute.CheckDelete` | should-fix | #3 | Design ADR-008 contract |
| P5 | Long-enum validators (`cpu.type`, `vga.type`) duplicate PVE validation | should-fix | #3 | Design ADR-004 enum rule (Q4) |
| P6 | `proxmox_virtual_environment_vm2` long resource type name + `MoveState` exist | blocker | #4 | Design D2 |
| P7 | Error message format inconsistency vs ADR-005 `"Unable to [Action] VM %d"` | should-fix | #5 | Design Phase 1 PR #5 |
