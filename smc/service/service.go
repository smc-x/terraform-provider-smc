package service

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"

	request "github.com/smc-x/terraform-provider-smc/smc/service/resource_request"
)

func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		request.New,
	}
}
