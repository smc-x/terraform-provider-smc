#!/bin/bash
cat > ~/.terraformrc <<EOF
provider_installation {

  dev_overrides {
    "registry.terraform.io/smc-x/smc" = "$(pwd)/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
EOF
