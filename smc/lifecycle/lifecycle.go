package lifecycle

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"

	fragile "github.com/smc-x/terraform-provider-smc/smc/lifecycle/resource_fragile"
	normal "github.com/smc-x/terraform-provider-smc/smc/lifecycle/resource_normal"
)

func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		normal.New,
		fragile.New,
	}
}
