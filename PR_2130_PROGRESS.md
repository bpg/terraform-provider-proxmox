# PR 2130 progress (storage provisioning)

This file tracks **current state**, **proof of work**, and **next steps** for PR 2130 so work can resume in a future session.

## Current branch

- Branch: `feat/support-cluster-storage`

## Completed phases (commits)

- `76a67beb` `fix(storage): correct cluster storage client encoding`
  - fixed `/storage` path building and encoding (`shared`, backups retention via query encoder) and added unit tests.
- `c643978a` `fix(storage): improve framework storage CRUD`
  - framework storage resources: import support, safer read error handling, re-read after create/update, avoid empty `nodes`/`content`, wire directory `shared` + backups.
- `2fee2010` `test(storage): add directory storage acceptance test`
  - baseline acceptance test for `proxmox_virtual_environment_storage_directory`.
- `16a9dc1f` `test(storage): add backups and import coverage`
  - extend directory storage acceptance test with backups update + import verify (retention ignored on import).
  - set default `disable/shared` values on read when the API omits them.

## Proof of work

### Lint/unit tests

Run before each phase commit:

```bash
make lint
go test ./... -count=1
```

### Acceptance test (directory storage)

```bash
./testacc TestAccResourceStorageDirectory -- -- -count 1
```

Optional env override for directory path:

- `PROXMOX_VE_ACC_STORAGE_DIR_PATH` (defaults to `/var/lib/vz`)

### Mitmproxy verification (sanitized)

We validated the `/api2/json/storage` calls using `mitmdump` and only logged **request path + parameter keys** (no values).

Commands used:

```bash
# start mitmdump (regular proxy on 127.0.0.1:8080)
/opt/homebrew/bin/mitmdump --mode regular --listen-host 127.0.0.1 --listen-port 8080 \
  --set ssl_insecure=true --set flow_detail=0 --set console_eventlog_verbosity=error \
  -s /tmp/mitm_storage_sanitize.py > /tmp/mitm_storage.out 2>&1 &

# run the test through the proxy
HTTP_PROXY=http://127.0.0.1:8080 HTTPS_PROXY=http://127.0.0.1:8080 NO_PROXY=127.0.0.1,localhost,::1 \
  ./testacc TestAccResourceStorageDirectory -- -- -count 1
```

Observed (sanitized) calls:

- `POST /api2/json/storage` body keys: `content, disable, nodes, path, shared, storage, type`
- `PUT /api2/json/storage/<id>` body keys: `content, disable, nodes, shared`
- `PUT /api2/json/storage/<id>` body keys: `content, disable, max-protected-backups, nodes, prune-backups, shared`
- `GET /api2/json/storage/<id>` (multiple refreshes)
- `DELETE /api2/json/storage/<id>`

Notes:
- Terraform/plugin tooling may attempt calls to `checkpoint-api.hashicorp.com` which will fail TLS interception unless the mitm CA is trusted. This does not affect Proxmox API verification.

## Remaining gaps / next steps

### Acceptance tests

- Add acceptance coverage for **other storage types** (likely gated by env):
  - NFS (`proxmox_virtual_environment_storage_nfs`)
  - SMB/CIFS (`proxmox_virtual_environment_storage_smb`)
  - PBS (`proxmox_virtual_environment_storage_pbs`)
  - LVM (`proxmox_virtual_environment_storage_lvm`)
  - LVMThin (`proxmox_virtual_environment_storage_lvmthin`)
  - ZFS (`proxmox_virtual_environment_storage_zfspool`)

### API round-tripping / drift

- Ensure fields like `snapshot_as_volume_chain` are readable back from `/storage/<id>` (otherwise they can’t be verified after import/refresh).
- Consider whether backup retention should be importable (depends on whether the API returns it on GET).

### Docs/examples

- Add at least one example under `examples/resources/` for directory storage.
- Run `make docs` and verify generated docs align with the new resources.


