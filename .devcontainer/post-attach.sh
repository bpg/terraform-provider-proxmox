#!/usr/bin/env bash
# Post-attach setup script for development container

set -e

# Color codes
GREEN='\033[1;32m'
CYAN='\033[1;36m'
RESET='\033[0m'

echo -e "${CYAN}🚀 Development Environment${RESET}"
echo "  Go:        $(go version | awk '{print $3}' | sed 's/^go//')"
echo "  Terraform: $(terraform version -json 2>/dev/null | jq -r '.terraform_version' 2>/dev/null || echo 'N/A')"
echo "  Node.js:   $(node --version 2>/dev/null || echo 'N/A')"
echo "  Docker:    $(docker --version 2>/dev/null | awk '{print $3}' | sed 's/,$//' || echo 'N/A')"
echo ""

# Unset Git committer variables
unset GIT_COMMITTER_NAME
unset GIT_COMMITTER_EMAIL

# Configure Terraform dev overrides
cat > ~/.terraformrc <<EOF
provider_installation {
  dev_overrides {
    "bpg/proxmox" = "${GOPATH}/bin/"
  }
  direct {}
}
EOF

# Configure OpenTofu if available
if command -v tofu &> /dev/null; then
	cat > ~/.tofurc <<EOF
provider_installation {
  dev_overrides {
    "bpg/proxmox" = "\${GOPATH}/bin/"
  }
  direct {}
}
EOF
fi

# Set environment variables
export TF_LOG="INFO"
export TF_LOG_PATH="${HOME}/.terraform/log"
mkdir -p ~/.terraform

# Verify Docker
if docker ps > /dev/null 2>&1; then
	echo -e "${GREEN}✓${RESET} Docker daemon accessible"
fi

# Verify Go tools
for tool in golangci-lint goimports air; do
	if command -v "$tool" &> /dev/null; then
		echo -e "${GREEN}✓${RESET} $tool installed"
	fi
done

echo ""
echo -e "${GREEN}✅ Ready to code!${RESET}"
