---
layout: page
page_title: "Development Troubleshooting"
subcategory: Guides
description: |-
  Common issues and solutions when developing the Proxmox provider.
---

# Development troubleshooting

This guide covers common issues encountered during provider development and their solutions.

## Acceptance tests

### Sandbox permission errors

When running acceptance tests in a sandboxed environment (e.g., Cursor IDE), you may see errors like:

```text
operation not permitted
```

or

```text
xargs: sysconf(_SC_ARG_MAX) failed
```

**Cause:** The sandbox restricts access to system paths like `~/Library/Caches/go-build`.

**Solutions:**

1. Run tests without sandboxing (request "all" permissions).
2. Set Go cache paths inside the workspace:

   ```sh
   export GOCACHE="$PWD/.cache/go-build"
   export GOMODCACHE="$PWD/.cache/go-mod"
   ./testacc TestName
   ```

### Proxy configuration issues

If you use an HTTP proxy and see errors like:

```text
Request cancelled
```

or provider reattach failures, the issue is that Terraform tries to route localhost traffic through the proxy.

**Solution:** Ensure `NO_PROXY` includes localhost addresses:

```sh
export NO_PROXY="127.0.0.1,localhost,::1"
```

The `./testacc` script automatically adds these when proxy environment variables are set (unless you pass `--no-proxy`).

### Stuck test VMs

Test VMs can get stuck if:

- They lack a boot disk (stuck in boot loop).
- They have `onboot = 1` and auto-restart after being stopped.
- A lock file prevents destruction.

**Cleanup procedure:**

SSH to the Proxmox node and run:

```sh
# List test VMs
qm list | grep test

# Disable auto-start
qm set <vmid> --onboot 0 --skiplock

# Kill the QEMU process
kill -9 $(cat /var/run/qemu-server/<vmid>.pid)

# Remove lock file
rm -f /var/lock/qemu-server/lock-<vmid>.conf

# Destroy the VM
qm destroy <vmid> --purge --skiplock
```

### Test timeout issues

If tests hang or timeout, you can pass additional flags to the test runner:

```sh
./testacc TestName -- -timeout 10m -count 1
```

## Build issues

### Linter errors

`make lint` automatically fixes formatting errors detected by `gofmt`, `gofumpt`, and `goimports`.
If it reports errors, most likely they require a non-trivial code change / manual fix.
Inspect the errors and fix them accordingly.

### Documentation generation

If `make docs` fails or produces unexpected output, ensure you have the correct version of `tfplugindocs`:

```sh
go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
make docs
```

## Provider development

### Changes not reflected in Terraform

If your code changes aren't showing up when running `terraform plan`:

1. Rebuild and reinstall the provider:

   ```sh
   go install .
   ```

2. Verify your `~/.terraformrc` (or `%APPDATA%/terraform.rc` on Windows) points to the correct `$GOPATH/bin`.

3. Check that no cached provider binary exists in `.terraform/providers/`.

### API debugging with mitmproxy

To inspect Proxmox API calls:

1. Start mitmproxy:

   ```sh
   mitmproxy --mode regular --listen-port 8080
   ```

2. Configure the provider to use the proxy:

   ```sh
   export HTTPS_PROXY="http://localhost:8080"
   export PROXMOX_VE_INSECURE="true"
   ```

3. Run your Terraform commands and inspect traffic in mitmproxy.

> [!WARNING]
> Never commit proxy configurations, captured traffic, or credentials to the repository.
