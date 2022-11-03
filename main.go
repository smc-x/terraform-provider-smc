package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/smc-x/terraform-provider-smc/smc"
)

func main() {
	providerserver.Serve(context.Background(), smc.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/smc-x/smc",
	})
}
