# Unattended Proxmox VE VM script

This script automates the whole [Option 2](./README.md#option-2---using-virsh-and-an-unattended-installation) workflow. It creates a Proxmox VE VM on a local libvirt/KVM host, without any interaction:

1. Creates the storage pool if it does not exist yet, and starts the libvirt network if needed
2. Reserves a fixed IP address for the VM (DHCP reservation, in the subnet of the libvirt network)
3. Builds the Proxmox ISO with your `answer.toml` and uploads it to the storage pool
4. Creates the VM, waits for the unattended installation to finish and makes sure the VM is running
5. Waits for the Proxmox web console to answer, then prints its address

## Requirements

- KVM and libvirt are installed and working, and your user is in the `libvirt` group (see the first steps of Option 2).
- `virsh`, `virt-install`, `xorriso` and `proxmox-auto-install-assistant` are installed.
- A Proxmox VE ISO and an `answer.toml` file that uses `source = "from-dhcp"` for the network.

## Usage

Save the script below as `run.sh` and make it executable:

```shell
chmod +x run.sh
```

Create the VM:

```shell
./run.sh /path/to/proxmox-ve_9.2-1.iso answer.toml
```

When the script ends, it prints the address of the Proxmox web console. You can also retrieve the VM IP address with the following command. The IP address is the one reserved for the VM, so it is always the same for a given network.

```shell
virsh domifaddr pve
```

Remove the VM and everything that was created for it (disk, uploaded ISO, DHCP reservation):

```shell
./run.sh --destroy
```

## Configuration

The defaults work for most setups. Every value can be overridden with an environment variable:

| Variable     | Default             | Description                                                 |
| ------------ | ------------------- | ----------------------------------------------------------- |
| `VM_NAME`    | `pve`               | Name of the VM                                              |
| `VM_POOL`    | `default`           | libvirt storage pool (created if it does not exist)         |
| `VM_NETWORK` | `default`           | libvirt network the VM is attached to                       |
| `VM_MAC`     | `52:54:00:00:00:43` | MAC address of the VM, bound to the reserved IP             |
| `VM_HOST_ID` | `43`                | Last octet of the reserved IP, in the subnet of the network |
| `VM_MEMORY`  | `4096`              | Memory in MiB                                               |
| `VM_VCPUS`   | `4`                 | Number of virtual CPUs                                      |
| `VM_DISK_GB` | `30`                | Disk size in GiB                                            |
| `VM_WAIT`    | `30`                | Maximum duration of the installation, in minutes            |

For example, to create a bigger VM:

```shell
VM_NAME=pve2 VM_MEMORY=16384 VM_DISK_GB=100 ./run.sh proxmox-ve_9.2-1.iso answer.toml
```

> **Note:** to run several VMs at the same time, give each one its own `VM_NAME`, `VM_MAC` and `VM_HOST_ID`.
>
> `--destroy` removes the DHCP reservation using the current values of `VM_MAC` and `VM_HOST_ID`, so pass the same values you used to create the VM.

## How it works

The script reserves a fixed IP address for the VM with a DHCP reservation in the libvirt network. The address is built from the subnet of the network (`virsh net-dumpxml default` shows it), so `VM_HOST_ID=43` gives `192.168.122.43` on a `192.168.122.0/24` network and `192.168.123.43` on a `192.168.123.0/24` one. The subnet can differ from one host to another: libvirt picks another range when `192.168.122.0/24` is already taken, for example when you run inside a VM that also uses libvirt.

To do the same reservation by hand, `net-update` applies it immediately without restarting the network:

```shell
virsh net-update default add ip-dhcp-host \
  "<host mac='52:54:00:00:00:43' name='pve' ip='192.168.122.43'/>" \
  --live --config
```

You can also run `virsh net-edit default` and add the `<host .../>` line inside the `<dhcp>` block, but the network must then be restarted, which briefly interrupts the other VMs attached to it.

A reserved address also shows up in `virsh domifaddr pve`. An address configured statically inside the guest does not, because libvirt only sees DHCP leases (use `virsh domifaddr pve --source arp` in that case).

## Troubleshooting

The script uses `--os-variant debian13`. If your system does not list it (`osinfo-query os | grep -i debian`), edit the value in the script. `debiantesting` is a safe fallback.

## Script

```shell
#!/usr/bin/env bash
#
# Create (or destroy) an unattended Proxmox VE VM on a local libvirt/KVM host.
#
# Usage:
#   ./run.sh /path/to/proxmox-ve.iso [answer.toml]   create and start the VM
#   ./run.sh --destroy                               remove the VM and what was created for it
#
# Optional environment variables (defaults in brackets):
#   VM_NAME [pve]            VM name
#   VM_POOL [default]        libvirt storage pool (created if missing)
#   VM_NETWORK [default]     libvirt network
#   VM_MAC [52:54:00:00:00:43]  MAC address, bound to the reserved IP
#   VM_HOST_ID [43]          last octet of the reserved IP in the network subnet
#   VM_MEMORY [4096]         memory in MiB
#   VM_VCPUS [4]             number of vCPUs
#   VM_DISK_GB [30]          disk size in GiB
#   VM_WAIT [30]             maximum duration of the installation, in minutes

set -euo pipefail

export LIBVIRT_DEFAULT_URI=qemu:///system

VM=${VM_NAME:-pve}
POOL=${VM_POOL:-default}
NET=${VM_NETWORK:-default}
MAC=${VM_MAC:-52:54:00:00:00:43}
HOST_ID=${VM_HOST_ID:-43}
MEMORY=${VM_MEMORY:-4096}
VCPUS=${VM_VCPUS:-4}
DISK_GB=${VM_DISK_GB:-30}
WAIT=${VM_WAIT:-30}
ISO_VOL="$VM-auto.iso"

die() { echo "error: $*" >&2; exit 1; }
log() { echo "==> $*"; }

need() {
  for cmd in "$@"; do
    command -v "$cmd" >/dev/null 2>&1 || die "missing command: $cmd"
  done
}

# Prints the reserved IP, built from the subnet of the libvirt network
reserved_ip() {
  local prefix
  prefix=$(virsh net-dumpxml "$NET" \
    | sed -n "s/.*<ip address='\([0-9.]*\)\.[0-9]*'.*/\1/p" | head -n1)
  [ -n "$prefix" ] || die "cannot find the IPv4 subnet of network '$NET'"
  echo "$prefix.$HOST_ID"
}

host_xml() { echo "<host mac='$MAC' name='$VM' ip='$1'/>"; }

destroy() {
  need virsh
  local ip
  ip=$(reserved_ip 2>/dev/null || true)
  log "Removing VM '$VM'"
  virsh destroy "$VM" >/dev/null 2>&1 || true
  virsh undefine "$VM" --remove-all-storage --nvram >/dev/null 2>&1 \
    || virsh undefine "$VM" --remove-all-storage >/dev/null 2>&1 || true
  virsh vol-delete "$ISO_VOL" --pool "$POOL" >/dev/null 2>&1 || true
  if [ -n "$ip" ]; then
    virsh net-update "$NET" delete ip-dhcp-host "$(host_xml "$ip")" \
      --live --config >/dev/null 2>&1 || true
  fi
  log "Done"
}

if [ "${1:-}" = "--destroy" ]; then
  destroy
  exit 0
fi

ISO_SRC=${1:-}
ANSWER=${2:-answer.toml}
[ -n "$ISO_SRC" ] || die "usage: $0 /path/to/proxmox-ve.iso [answer.toml] | --destroy"
[ -f "$ISO_SRC" ] || die "ISO not found: $ISO_SRC"
[ -f "$ANSWER" ] || die "answer file not found: $ANSWER"

need virsh virt-install proxmox-auto-install-assistant xorriso
[ -e /dev/kvm ] || die "/dev/kvm not found: KVM is not available on this host"
virsh uri >/dev/null 2>&1 || die "cannot connect to libvirt: is your user in the 'libvirt' group?"

if virsh dominfo "$VM" >/dev/null 2>&1; then
  die "VM '$VM' already exists (use '$0 --destroy' to remove it, or set VM_NAME)"
fi

WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT

# Storage pool: create it if missing, start it if inactive
if ! virsh pool-info "$POOL" >/dev/null 2>&1; then
  log "Creating storage pool '$POOL'"
  virsh pool-define-as "$POOL" dir --target "/var/lib/libvirt/images/$POOL" >/dev/null
  virsh pool-build "$POOL" >/dev/null
  virsh pool-autostart "$POOL" >/dev/null
fi
if [ "$(virsh pool-info "$POOL" | awk '/^State:/{print $2}')" != "running" ]; then
  virsh pool-start "$POOL" >/dev/null
fi

# Network: make sure it is running
virsh net-info "$NET" >/dev/null 2>&1 || die "libvirt network '$NET' does not exist"
if [ "$(virsh net-info "$NET" | awk '/^Active:/{print $2}')" != "yes" ]; then
  log "Starting network '$NET'"
  virsh net-start "$NET" >/dev/null
fi

# Fixed IP: DHCP reservation in the subnet of the network, whatever it is
IP=$(reserved_ip)
if ! virsh net-dumpxml "$NET" | grep -q "mac='$MAC'"; then
  log "Reserving $IP for $MAC on network '$NET'"
  virsh net-update "$NET" add ip-dhcp-host "$(host_xml "$IP")" --live --config >/dev/null
fi

# Build the ISO with the answer file
log "Preparing the ISO"
proxmox-auto-install-assistant prepare-iso "$ISO_SRC" \
  --fetch-from iso --answer-file "$ANSWER" --output "$WORK/$ISO_VOL"

# Upload it to the storage pool (no sudo needed, and QEMU can read it)
log "Uploading the ISO to pool '$POOL'"
virsh vol-delete "$ISO_VOL" --pool "$POOL" >/dev/null 2>&1 || true
virsh vol-create-as "$POOL" "$ISO_VOL" "$(stat -c%s "$WORK/$ISO_VOL")" --format raw >/dev/null
virsh vol-upload --pool "$POOL" "$ISO_VOL" "$WORK/$ISO_VOL"

# Create the VM; libvirt creates the disk in the pool.
# virt-install waits until the unattended installation is finished.
log "Installing Proxmox VE (up to $WAIT minutes)"
virt-install --name "$VM" --memory "$MEMORY" --vcpus "$VCPUS" \
  --cpu host-passthrough --os-variant debian13 \
  --disk pool="$POOL",size="$DISK_GB",bus=virtio,format=qcow2 \
  --cdrom "$(virsh vol-path --pool "$POOL" "$ISO_VOL")" \
  --network network="$NET",model=virtio,mac="$MAC" \
  --graphics vnc,listen=127.0.0.1 \
  --noautoconsole --wait "$WAIT"

# virt-install may restart the VM itself once the installation is complete.
# Start it only if it is not already running.
log "Making sure the VM is running"
for _ in $(seq 1 30); do
  state=$(virsh domstate "$VM")
  [ "$state" = "running" ] && break
  if [ "$state" = "shut off" ]; then
    virsh start "$VM" >/dev/null
    break
  fi
  sleep 2
done

# Wait for the Proxmox web console to answer
log "Waiting for the web console on $IP:8006"
ready=0
for _ in $(seq 1 60); do
  if (exec 3<>"/dev/tcp/$IP/8006") 2>/dev/null; then ready=1; break; fi
  sleep 5
done

if [ "$ready" = 1 ]; then
  log "Ready: https://$IP:8006 (ssh root@$IP)"
else
  echo "warning: no answer on $IP:8006 after 5 minutes; check the VM with 'virsh domifaddr $VM'" >&2
fi
```
