# PR #2689 — Linux Bond Interface Resource: Test Results

**Date:** 2026-03-13
**Branch:** `claude/add-network-bonding-5h2ys`

## Test Environment

| Component | Details |
|---|---|
| **Host PVE** | pve-manager/8.4.17 (3-node cluster, kernel 6.8.12-15-pve) |
| **Test PVE (nested)** | pve-manager/8.0.3 (VM 9999 on pve-1, kernel 6.2.16-3-pve) |
| **Test VM config** | 4 cores, 8GB RAM, 64GB RBD disk, `--cpu host` (nested virt) |
| **Slave interfaces** | `ens19`, `ens20` (hotplugged virtio NICs on vmbr0) |
| **Go version** | 1.25+ |

## Review Feedback Applied

| # | Issue (gemini-code-assist) | Resolution |
|---|---|---|
| 1 | Slaves list not sorted on read — causes persistent diffs | **Fixed** — added `sort.Strings()` after splitting API response in `importFromNetworkInterfaceList` |
| 2 | Comment not nullified when API returns nil | **Not applied** — existing bridge/vlan resources intentionally preserve previous state to avoid plan drift when `comment = ""` (PVE stores empty as nil) |
| 3 | Name regex doesn't allow "bond0" | **No change needed** — `^[A-Za-z][A-Za-z0-9]{0,9}$` already matches `bond0` |
| 4 | No default for `bond_mode` | **Fixed** — added `stringdefault.StaticString("balance-rr")` to schema |

### Additional Fixes Found During Testing

| # | Issue | Resolution |
|---|---|---|
| 5 | Test attribute check for `bond_xmit_hash_policy` failed — `+` in `layer3+4` treated as regex quantifier by `test.ResourceAttributes` | **Fixed** — escaped to `layer3\+4` in test |

## Test Results

### Build

```
$ make build
go build -o "./build/terraform-provider-proxmox_v0.98.1"
✅ PASS
```

### Lint

```
$ make lint
golangci-lint fmt
golangci-lint run --fix
0 issues.
✅ PASS
```

### Acceptance Test

```
$ ./testacc --no-proxy TestAccResourceLinuxBond
ok  github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/network  10.394s
✅ PASS (4/4 steps)
```

**Steps executed:**

1. **Create** — bond with `802.3ad` mode, `layer3+4` hash policy, address, autostart, comment
2. **Update** — change to `active-backup` mode with `bond_primary`, remove hash policy, change address, toggle autostart
3. **Update** — remove address, remove bond_primary, switch to `balance-rr` mode
4. **ImportState** — verify import by `node_name:iface` ID

### Docs

```
$ make docs
✅ Generated (includes network_linux_bond resource)
```

## Files Changed

- `fwprovider/nodes/network/resource_linux_bond.go` — slaves sorting on read, bond_mode default, stringdefault import
- `fwprovider/nodes/network/resource_linux_bond_test.go` — regex escaping, slave interface names
