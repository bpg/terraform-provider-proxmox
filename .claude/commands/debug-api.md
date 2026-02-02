---
name: debug-api
description: Debug Proxmox API calls using mitmproxy to verify parameters are sent correctly
argument-hint: \[TestName\] \[param-to-verify\]
allowed-tools:
  - Read
  - Bash
  - Grep
  - Glob
  - AskUserQuestion
---

<objective>
Debug API interactions between the Terraform provider and Proxmox VE using mitmproxy.

Use this skill when:

- Implementing new API parameters and need to verify they're sent in correct API calls with expected values
- Debugging API errors or unexpected behavior
- Verifying a fix sends correct parameters before marking work complete
- Tests pass but you suspect API calls might be wrong

**Remember:** Tests passing ≠ API calls are correct. Always verify with mitmproxy.
</objective>

<context>
Test name: $0
Parameter to verify: $1

Reference: [DEBUGGING.md](../../.dev/DEBUGGING.md)
</context>

<process>

## Step 1: Verify Prerequisites

Check mitmproxy is available:

```bash
which mitmdump || echo "ERROR: mitmproxy not installed. Run: `brew install mitmproxy` (macOS) or `pip install mitmproxy`"
```

Check no existing proxy is running:

```bash
pgrep -f mitmdump && echo "WARNING: mitmproxy already running" || echo "OK: No existing proxy"
```

Check port 8080 is available:

```bash
lsof -i :8080 && echo "WARNING: Port 8080 in use by process above" || echo "OK: Port 8080 available"
```

If proxy already running or port in use, ask user:

- "Stop existing proxy and start fresh?"
- "Use existing proxy?"

## Step 2: Determine Test and Parameter

If `$0` (test name) not provided:

```text
AskUserQuestion(
  header: "Test Name",
  question: "Which acceptance test should I run?",
  options: [
    { label: "Let me search", description: "I'll help find the right test" }
  ]
)
```

If user wants to search, use Grep to find relevant tests:

```bash
grep -r "func TestAcc" fwprovider/ proxmoxtf/ --include="*_test.go" | grep -i "{keyword}" | head -10
```

If `$1` (parameter) not provided:

```text
AskUserQuestion(
  header: "Parameter",
  question: "What parameter or API endpoint should I verify?",
  options: [
    { label: "All parameters", description: "Capture all API traffic for analysis" }
  ]
)
```

Store test name as `$TEST_NAME` and parameter as `$PARAM`.

## Step 3: Start Mitmproxy

Choose logging approach based on need:

**Standard (recommended):**

```bash
mitmdump --mode regular --listen-port 8080 --flow-detail 2 > /tmp/api_debug.log 2>&1 &
PROXY_PID=$!
echo "Proxy started with PID: $PROXY_PID"
sleep 1
pgrep -f mitmdump > /dev/null && echo "OK: Proxy running" || echo "ERROR: Proxy failed to start"
```

**Enhanced (for detailed analysis):**

```bash
mitmdump --mode regular --listen-port 8080 -s .dev/proxmox_debug_script.py > /tmp/api_debug.log 2>&1 &
PROXY_PID=$!
echo "Proxy started with PID: $PROXY_PID (enhanced script)"
sleep 1
```

Report to user: "Mitmproxy started on port 8080. Log file: /tmp/api_debug.log"

## Step 4: Run the Test

Execute the acceptance test with proxy and TLS environment variables set explicitly:

```bash
HTTP_PROXY=http://127.0.0.1:8080 HTTPS_PROXY=http://127.0.0.1:8080 PROXMOX_VE_INSECURE=true ./testacc ${TEST_NAME}
```

Note: `PROXMOX_VE_INSECURE=true` is required because mitmproxy presents its own certificate when intercepting HTTPS traffic.

Capture the exit code:

```bash
TEST_EXIT=$?
echo "Test exit code: $TEST_EXIT"
```

Report test result to user:

- Exit 0: "Test PASSED"
- Exit non-zero: "Test FAILED (exit code: $TEST_EXIT)"

## Step 5: Verify Traffic Was Captured

```bash
[ -s /tmp/api_debug.log ] && echo "OK: Traffic captured ($(wc -l < /tmp/api_debug.log) lines)" || echo "WARNING: No traffic captured - proxy may not be working"
```

If no traffic was captured, check:

- Was the proxy running during the test?
- Are `HTTP_PROXY` and `HTTPS_PROXY` set correctly?
- Is `PROXMOX_VE_INSECURE=true` set?

## Step 6: Stop Mitmproxy

```bash
pkill -f mitmdump
sleep 1
pgrep -f mitmdump && echo "WARNING: Proxy still running" || echo "OK: Proxy stopped"
```

## Step 7: Analyze API Traffic

**If specific parameter was requested:**

```bash
echo "=== Searching for parameter: ${PARAM} ==="
grep -i "${PARAM}" /tmp/api_debug.log | head -20
```

**Show all API calls:**

```bash
echo "=== API Calls Summary ==="
grep -E "GET|POST|PUT|DELETE" /tmp/api_debug.log | grep "api2/json" | head -30
```

**Check for errors:**

```bash
echo "=== Error Responses ==="
grep -E "400|401|403|404|500|502|503" /tmp/api_debug.log | head -10
```

## Step 8: Present Findings

Summarize findings for user:

1. **Parameter verification:**
   - Found/Not found in requests
   - Show the actual request line(s) containing the parameter

2. **API calls made:**
   - List unique endpoints called
   - Note any unexpected calls

3. **Errors detected:**
   - Any 4xx/5xx responses
   - Error message content if available

4. **Recommendation:**
   - If parameter found: "API call verified. Parameter `{param}` sent correctly."
   - If parameter NOT found: "WARNING: Parameter `{param}` not found in API traffic. Check implementation."
   - If errors: "API errors detected. Review the error responses above."

## Step 9: Offer Next Steps

```text
AskUserQuestion(
  header: "Next Step",
  question: "What would you like to do next?",
  options: [
    { label: "View full log", description: "Show complete API traffic log" },
    { label: "Run another test", description: "Debug a different test" },
    { label: "Done", description: "Debugging complete" }
  ]
)
```

If "View full log":

```bash
cat /tmp/api_debug.log
```

If "Run another test": Loop back to Step 2.

</process>

<success_criteria>

- [ ] Mitmproxy started successfully
- [ ] Test executed
- [ ] Traffic capture verified
- [ ] Proxy stopped cleanly
- [ ] API traffic analyzed
- [ ] Parameter presence verified (if specified)
- [ ] Findings presented to user
</success_criteria>

<tips>
- Use `--flow-detail 4` for untruncated output when debugging complex issues
- The enhanced script (`.dev/proxmox_debug_script.py`) categorizes API calls and highlights key parameters
- Log file persists at `/tmp/api_debug.log` for later review
- If tests timeout, check proxy is running: `pgrep -f mitmdump`
</tips>
