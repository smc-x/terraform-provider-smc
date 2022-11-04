package overleaf

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"

	meta "github.com/smc-x/terraform-provider-smc/smc/overleaf_meta"
)

func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		meta.NewOverleafMetaResource,
	}
}
