#!/bin/sh
#
# Clean up all test resources on the PVE host after a test run.
# Usage: ssh root@pve < .dev/pve-cleanup.sh
#

set +e  # Best-effort cleanup — don't stop on individual failures

if ! command -v pvesh >/dev/null 2>&1; then
  echo "ERROR: This script must run on a Proxmox VE host." >&2
  echo "Usage: ssh root@pve < .dev/pve-cleanup.sh" >&2
  exit 1
fi

echo "=== Cleaning VMs ==="
for id in $(qm list 2>/dev/null | awk 'NR>1 {print $1}'); do
  qm unlock "$id" 2>/dev/null || true
  qm stop "$id" 2>/dev/null || true
  qm destroy "$id" --purge 2>/dev/null && echo "  destroyed VM $id" || echo "  failed VM $id"
done

echo "=== Cleaning Containers ==="
for id in $(pct list 2>/dev/null | awk 'NR>1 {print $1}'); do
  pct unlock "$id" 2>/dev/null || true
  pct stop "$id" 2>/dev/null || true
  pct destroy "$id" --purge 2>/dev/null && echo "  destroyed CT $id" || echo "  failed CT $id"
done

echo "=== Cleaning Pools ==="
pvesh get /pools --output-format json 2>/dev/null \
  | python3 -c "import sys,json; [print(p['poolid']) for p in json.load(sys.stdin)]" 2>/dev/null \
  | while read -r pool; do
      pvesh delete "/pools/$pool" 2>/dev/null && echo "  deleted pool $pool"
    done

echo "=== Cleaning SDN VNets ==="
pvesh get /cluster/sdn/vnets --output-format json 2>/dev/null \
  | python3 -c "import sys,json; [print(v['vnet']) for v in json.load(sys.stdin)]" 2>/dev/null \
  | while read -r vnet; do
      pvesh delete "/cluster/sdn/vnets/$vnet" 2>/dev/null && echo "  deleted vnet $vnet"
    done
pvesh put /cluster/sdn 2>/dev/null
sleep 2

echo "=== Cleaning SDN Zones ==="
pvesh get /cluster/sdn/zones --output-format json 2>/dev/null \
  | python3 -c "import sys,json; [print(z['zone']) for z in json.load(sys.stdin) if z['zone'] != 'localnetwork']" 2>/dev/null \
  | while read -r zone; do
      pvesh delete "/cluster/sdn/zones/$zone" 2>/dev/null && echo "  deleted zone $zone"
    done
pvesh put /cluster/sdn 2>/dev/null

echo "=== Cleaning Firewall Security Groups ==="
pvesh get /cluster/firewall/groups --output-format json 2>/dev/null \
  | python3 -c "import sys,json; [print(g['group']) for g in json.load(sys.stdin)]" 2>/dev/null \
  | while read -r group; do
      pvesh delete "/cluster/firewall/groups/$group" 2>/dev/null && echo "  deleted group $group"
    done

echo "=== Cleaning Firewall IPSets ==="
pvesh get /cluster/firewall/ipset --output-format json 2>/dev/null \
  | python3 -c "import sys,json; [print(s['name']) for s in json.load(sys.stdin)]" 2>/dev/null \
  | while read -r ipset; do
      pvesh delete "/cluster/firewall/ipset/$ipset" 2>/dev/null && echo "  deleted ipset $ipset"
    done

echo "=== Cleaning Firewall Rules ==="
pvesh get /cluster/firewall/rules --output-format json 2>/dev/null \
  | python3 -c "import sys,json; [print(r['pos']) for r in reversed(json.load(sys.stdin))]" 2>/dev/null \
  | while read -r pos; do
      pvesh delete "/cluster/firewall/rules/$pos" 2>/dev/null && echo "  deleted rule $pos"
    done

echo "=== Cleaning VM/CT Firewall Files ==="
for f in /etc/pve/firewall/[0-9]*.fw; do
  [ -f "$f" ] && rm -f "$f" && echo "  removed $f"
done

echo "=== Resetting Cluster Firewall Config ==="
cat > /etc/pve/firewall/cluster.fw << "FWEOF"
[OPTIONS]

enable: 0

FWEOF

echo "=== Cleaning Stale NFS Mounts ==="
mount | grep 'type nfs' | grep -v '/mnt/pve/nfs ' | awk '{print $3}' | while read -r mnt; do
  umount "$mnt" 2>/dev/null && echo "  unmounted $mnt"
  rmdir "$mnt" 2>/dev/null
done

echo "=== Cleaning Test Files ==="
rm -f /var/lib/vz/template/iso/*test* /var/lib/vz/template/iso/fake_file*
rm -f /var/lib/vz/template/cache/*test* /var/lib/vz/template/cache/tpl-*
rm -f /var/lib/vz/snippets/snippet-raw-* /var/lib/vz/snippets/datasource-test-*

echo "=== Summary ==="
echo "VMs: $(qm list 2>/dev/null | tail -n+2 | wc -l)"
echo "CTs: $(pct list 2>/dev/null | tail -n+2 | wc -l)"
echo "Storage: $(pvesm status 2>/dev/null | grep -c active) active"
echo "Done"
