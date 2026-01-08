"""
Enhanced mitmproxy script for debugging Terraform Proxmox Provider.

Usage:
    mitmdump --mode regular --listen-port 8080 -s .dev/proxmox_debug_script.py

Features:
- Filters and highlights Proxmox VE API calls
- Shows query parameters (especially 'content' type filtering)
- Pretty-prints JSON request/response bodies
- Counts items in list responses
- Visual indicators for success/error
"""
from mitmproxy import http
import json


def is_proxmox_api_call(flow: http.HTTPFlow) -> bool:
    """Check if this is a Proxmox API call"""
    path = flow.request.path
    return "/api2/json/" in path


def request(flow: http.HTTPFlow) -> None:
    """Log and highlight important API calls"""
    if not is_proxmox_api_call(flow):
        return

    path = flow.request.path
    method = flow.request.method

    # Log all Proxmox API calls
    print(f"\n{'='*80}")

    # Determine API category for better context
    api_type = "PROXMOX API"
    if "/storage/" in path:
        api_type = "STORAGE API"
    elif "/qemu/" in path:
        api_type = "VM (QEMU) API"
    elif "/lxc/" in path:
        api_type = "CONTAINER (LXC) API"
    elif "/nodes/" in path and "/network" in path:
        api_type = "NETWORK API"
    elif "/cluster/" in path:
        api_type = "CLUSTER API"
    elif "/access/" in path:
        api_type = "ACCESS CONTROL API"
    elif "/pools/" in path:
        api_type = "POOL API"

    print(f"🔍 {api_type}")
    print(f"Method: {method}")
    print(f"Host: {flow.request.pretty_host}")
    print(f"Path: {path}")

    # Parse and highlight query parameters
    if flow.request.query:
        print(f"\n📋 Query Parameters:")
        for key, value in flow.request.query.items():
            # Highlight important common parameters
            if key in ["content", "vmid", "node", "storage", "type"]:
                emoji = "🎯"
            else:
                emoji = "  "
            print(f"  {emoji} {key} = {value}")

    # Show request body for POST/PUT/PATCH
    if method in ["POST", "PUT", "PATCH"] and flow.request.content:
        try:
            body = flow.request.content.decode('utf-8')
            print(f"\n📤 Request Body:")
            try:
                json_body = json.loads(body)
                print(json.dumps(json_body, indent=2))
            except json.JSONDecodeError:
                print(body)
        except UnicodeDecodeError:
            print(f"<binary data, {len(flow.request.content)} bytes>")


def response(flow: http.HTTPFlow) -> None:
    """Log response details for all Proxmox API calls"""
    if not is_proxmox_api_call(flow):
        return

    status = flow.response.status_code
    emoji = "✅" if 200 <= status < 300 else "❌"

    print(f"\n{emoji} Response: {status} {flow.response.reason}")

    # Parse and show JSON response
    if flow.response.content:
        try:
            body = flow.response.content.decode('utf-8')
            json_body = json.loads(body)

            # Highlight data count if it's a list response
            if "data" in json_body:
                data = json_body["data"]
                if isinstance(data, list):
                    print(f"📊 Returned {len(data)} items")
                    if data:
                        print(f"\n📄 First item preview:")
                        print(json.dumps(data[0], indent=2)[:300])
                        if len(json.dumps(data[0], indent=2)) > 300:
                            print("...")
                elif isinstance(data, dict):
                    print(f"\n📄 Response Data:")
                    preview = json.dumps(json_body, indent=2)[:500]
                    print(preview)
                    if len(json.dumps(json_body, indent=2)) > 500:
                        print("...")
                else:
                    # Scalar value (string, number, bool, etc.)
                    print(f"\n📄 Response Data: {data}")
            else:
                print(f"\n📄 Response:")
                preview = json.dumps(json_body, indent=2)[:500]
                print(preview)
                if len(json.dumps(json_body, indent=2)) > 500:
                    print("...")

        except json.JSONDecodeError as e:
            print(f"⚠️  Could not parse JSON response: {e}")
            try:
                print(f"Raw response: {body[:200]}")
            except:
                print(f"<could not decode response>")
        except UnicodeDecodeError:
            print(f"<binary response, {len(flow.response.content)} bytes>")

    print(f"{'='*80}\n")


# Optional: Add more specific endpoint handlers here
# For example, to track VM operations, task status, etc.
def vm_operations(flow: http.HTTPFlow) -> None:
    """Track VM-related operations (create, update, delete)"""
    path = flow.request.path
    method = flow.request.method

    if "/qemu/" in path and method in ["POST", "PUT", "DELETE"]:
        print(f"\n🖥️  VM Operation: {method} {path}")
