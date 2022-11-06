terraform {
  required_providers {
    smc = {
      source = "registry.terraform.io/smc-x/smc"
    }
  }
}

variable "token" {}
variable "endpoint" {}
variable "skip_verify" {}

provider "smc" {
  token       = var.token
  endpoint    = var.endpoint
  skip_verify = var.skip_verify
}

resource "smc_service_request" "step1" {
  # id = {computed}

  subj    = "workers"
  timeout = 1

  # resp = {computed json}
}

locals {
  worker = jsondecode(smc_service_request.step1.resp).worker
}

resource "smc_service_request" "step2" {
  # id = {computed}

  subj = local.worker

  data = jsonencode({
    "msg" : "Hello, ${local.worker}"
  })

  # resp = {computed json}
}
