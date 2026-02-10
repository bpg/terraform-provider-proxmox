# Debugging Guide for Terraform Proxmox Provider

This guide covers debugging techniques for developing and testing the Terraform Proxmox Provider.

> **Using an LLM agent?** The `/bpg:debug-api` skill automates most of this workflow. See [README.md](README.md#working-with-llm-agents) for details.

## Using Mitmproxy to Debug API Calls

Mitmproxy is an essential tool for intercepting and analyzing HTTP/HTTPS traffic between the provider and Proxmox VE API. It helps you verify that API calls are correct, debug issues, and understand provider behavior.

### Prerequisites

```bash
# Install mitmproxy (if not already installed)
brew install mitmproxy  # macOS
# or
pip install mitmproxy  # Linux/Windows
```

### Quick Start

The provider uses proxy environment variables (`HTTP_PROXY`, `HTTPS_PROXY`) when they are set. The `testacc.env` file should already configure these for you.

```bash
# 1. Start mitmproxy proxy
mitmdump --mode regular --listen-port 8080 --flow-detail 2 > /tmp/mitmproxy.log 2>&1 &

# 2. Run acceptance tests
./testacc TestAccDatasourceFile

# 3. Analyze the captured traffic
grep "storage.*content" /tmp/mitmproxy.log

# 4. Stop the proxy
pkill -f mitmdump
```

### Flow Detail Levels

| Level | Output                                                        |
|-------|---------------------------------------------------------------|
| 0     | No output (quiet)                                             |
| 1     | Shortened URL + status code (default)                         |
| 2     | **Full URL + headers** (recommended - shows query parameters) |
| 3     | Level 2 + truncated response content                          |
| 4     | **Everything untruncated** (deep debugging)                   |

### Using the Enhanced Debug Script

For better visibility into API interactions, use the custom debug script:

```bash
# Start proxy with custom script
mitmdump --mode regular --listen-port 8080 \
  -s .dev/proxmox_debug_script.py \
  > /tmp/debug.log 2>&1 &

# Run tests
./testacc TestAccDatasourceFileContentTypeFiltering

# View formatted output
cat /tmp/debug.log

# Stop proxy
pkill -f mitmdump
```

The script provides:

- 🔍 Filtered output (only Proxmox API calls)
- 🎯 Highlighted query parameters (especially `content` type)
- 📊 Item counts for list responses
- ✅/❌ Visual success/error indicators
- JSON pretty-printing

### Real-World Example

When debugging the content type filtering feature, mitmproxy captured:

```text
127.0.0.1:58448: GET https://pve.bpghome.net:8006/api2/json/nodes/pve/storage/local/content?content=import
    Host: pve.bpghome.net:8006
    User-Agent: Go-http-client/1.1
    Authorization: PVEAPIToken=terraform@pve!provider=...
 << 200 OK 11b
```

This confirmed:

- ✅ Query parameter `?content=import` was correctly sent
- ✅ API accepted the request (200 OK)
- ✅ Empty result (11 bytes) as expected for non-existent file

## Common Debugging Scenarios

### 1. Verify New API Parameter

```bash
# Start detailed logging
mitmdump --flow-detail 2 > /tmp/test.log 2>&1 &

# Run test
./testacc TestAccYourNewFeature

# Check the parameter is sent
grep "your_param=" /tmp/test.log

# Stop proxy
pkill -f mitmdump
```

### 2. Debug API Error Responses

```bash
# Use flow-detail 4 to see full response bodies
mitmdump --flow-detail 4 > /tmp/error_debug.log 2>&1 &

# Run failing test
./testacc TestAccFailingTest

# Search for error responses
grep -A 20 "400\|401\|403\|500" /tmp/error_debug.log

pkill -f mitmdump
```

### 3. Analyze Request Body Parameters

```bash
# Custom script shows request bodies
mitmdump -s .dev/proxmox_debug_script.py > /tmp/post_debug.log 2>&1 &

# Run test that creates resources
./testacc TestAccResourceDownloadFile

# View captured POST bodies
grep -A 30 "📤 Request Body" /tmp/post_debug.log

pkill -f mitmdump
```

### 4. Compare API Calls Between Versions

```bash
# Capture from current version
mitmdump --save-stream-file /tmp/before.mitm --flow-detail 3 > /tmp/before.log 2>&1 &
./testacc TestAccSomeFeature
pkill -f mitmdump

# Apply your changes
git checkout your-feature-branch
make build

# Capture from new version
mitmdump --save-stream-file /tmp/after.mitm --flow-detail 3 > /tmp/after.log 2>&1 &
./testacc TestAccSomeFeature
pkill -f mitmdump

# Compare
diff /tmp/before.log /tmp/after.log
```

## Log Analysis Commands

```bash
# View recent requests
tail -50 /tmp/mitmproxy.log

# Find all storage content API calls
grep "storage.*content" /tmp/mitmproxy.log

# Count requests by endpoint
grep "GET\|POST\|DELETE" /tmp/mitmproxy.log | cut -d' ' -f2-3 | sort | uniq -c

# Find requests with query parameters
grep "?" /tmp/mitmproxy.log | grep -v "tasks"

# Extract URLs only
grep "GET\|POST" /tmp/mitmproxy.log | awk '{print $2}' | sort | uniq

# Find error responses
grep -E "400|401|403|404|500" /tmp/mitmproxy.log
```

## Best Practices

1. **Always start proxy before tests** - Tests will fail with "connection refused" if proxy is expected but not running
2. **Use flow-detail 2 for most cases** - Shows query parameters without overwhelming detail
3. **Save proxy PID** - Makes it easy to stop the right instance: `PROXY_PID=$!`
4. **Check proxy is running** - `pgrep -f mitmdump` before running tests
5. **Clean up after tests** - `pkill -f mitmdump` to stop the proxy
6. **Analyze logs immediately** - Review output while test context is fresh

## Troubleshooting

### Proxy Not Starting

```bash
# Check if port is already in use
lsof -i :8080

# Kill existing process
kill $(lsof -t -i:8080)
```

### Tests Timing Out

```bash
# Verify proxy is running
pgrep -f mitmdump

# Check proxy port
ps aux | grep mitmdump | grep 8080
```

### No Output in Logs

```bash
# Increase verbosity
mitmdump --flow-detail 3 -v

# Monitor log file
tail -f /tmp/mitmproxy.log
```

## Other Debugging Tools

### Go Debugging with Delve

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug a test
dlv test ./fwprovider/test -- -test.run TestAccDatasourceFile
```

### Terraform Debug Logging

```bash
# Enable Terraform debug logs
export TF_LOG=DEBUG
export TF_LOG_PATH=/tmp/terraform.log

# Run terraform
terraform plan
terraform apply

# View logs
cat /tmp/terraform.log
```

### Provider Debug Logging

The provider uses `terraform-plugin-log` for structured logging. Logs are visible when `TF_LOG` is set.

## References

- [Mitmproxy Documentation](https://docs.mitmproxy.org/stable/)
- [Mitmproxy Script API](https://docs.mitmproxy.org/stable/addons/overview/)
- [Terraform Plugin Development](https://developer.hashicorp.com/terraform/plugin)
- [Proxmox VE API Documentation](https://pve.proxmox.com/pve-docs/api-viewer/)

## Files in This Directory

- `DEBUGGING.md` - This file
- `proxmox_debug_script.py` - Enhanced mitmproxy script for API analysis
