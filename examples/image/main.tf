terraform {
  required_providers {
    smc = {
      source = "registry.terraform.io/smc-x/smc"
    }
  }
  required_version = ">= 1.1.0"
}

provider "smc" {}

resource "smc_image" "ws" {
    workspace = "ws"
    owner = "me"
    base = "ubuntu:20.04"
}

output "ws_image" {
  value = smc_image.ws.id
}
