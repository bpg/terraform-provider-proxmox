#!/usr/bin/bash

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
