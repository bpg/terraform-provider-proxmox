#!/bin/sh
export TF_ACC=1
export TF_LOG=DEBUG
export PROXMOX_VE_INSECURE="false"
export PROXMOX_VE_API_TOKEN="root@pam!terraform=f3288ecd-9abe-4f49-b6e7-97bdce1886eb"
export PROXMOX_VE_ENDPOINT="https://pve.bpglabs.net:8006/"
export PROXMOX_VE_SSH_AGENT="false"
export PROXMOX_VE_SSH_USERNAME="terraform"
export PROXMOX_VE_SSH_PRIVATE_KEY="$(cat private.key)"

go test -count=1 -v github.com/bpg/terraform-provider-proxmox/fwprovider/tests/...