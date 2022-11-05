package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/smc-x/terraform-provider-smc/smc"
	"github.com/smc-x/terraform-provider-smc/smc/global"
)

func main() {
	defer global.Run()
	providerserver.Serve(context.Background(), smc.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/smc-x/smc",
	})
}
