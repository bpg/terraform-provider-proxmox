#!/usr/bin/env bash

# Display welcome banner
echo -e "\033[1;36m"
echo "════════════════════════════════════════════════════════════════════════════════════════════"
echo 
echo "  🚀 Terraform Provider For Proxmox Development Environment"
echo 
echo "  ⚠️ EXPERIMENTAL"
echo "     Use at your own risk! Some tools may be missing or not work as expected."
echo 
echo "  • Go Version: $(go version | cut -d' ' -f3 | sed 's/^go//')"
echo "  • Terraform Version: $(terraform version -json | jq -r '.terraform_version')"
echo "  • OpenTofu Version: $(tofu version -json | jq -r '.terraform_version')"
echo "  • Working Directory: $(pwd)"
echo 
echo "════════════════════════════════════════════════════════════════════════════════════════════"
echo -e "\033[0m"

# Workaround for https://github.com/orgs/community/discussions/75161
unset GIT_COMMITTER_NAME
unset GIT_COMMITTER_EMAIL

cat <<EOF > ~/.terraformrc
provider_installation {

  dev_overrides {
      "bpg/proxmox" = "${GOPATH}/bin/"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
EOF
