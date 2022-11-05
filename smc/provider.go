package smc

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/terraform-community-providers/terraform-plugin-framework-utils/modifiers"

	"github.com/smc-x/terraform-provider-smc/smc/global"
	"github.com/smc-x/terraform-provider-smc/smc/lifecycle"
	"github.com/smc-x/terraform-provider-smc/utils/logging"
	"github.com/smc-x/terraform-provider-smc/utils/nats"
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &smcProvider{}
}

// smcProvider is the provider implementation.
type smcProvider struct{}

// Metadata returns the provider type name.
func (p *smcProvider) Metadata(
	_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse,
) {
	resp.TypeName = "smc"
}

// model maps the resource schema data.
type model struct {
	Token      types.String `tfsdk:"token"`
	Endpoint   types.String `tfsdk:"endpoint"`
	SkipVerify types.Bool   `tfsdk:"skip_verify"`
}

// GetSchema defines the provider-level schema for configuration data.
func (p *smcProvider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"token": {
				Type:      types.StringType,
				Required:  true,
				Sensitive: true,
			},
			"endpoint": {
				Type:      types.StringType,
				Required:  true,
				Sensitive: true,
			},
			"skip_verify": {
				Type:     types.BoolType,
				Optional: true,
				Computed: true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					modifiers.DefaultBool(false),
				},
			},
		},
	}, nil
}

// Configure prepares an smc api client for data sources and resources.
func (p *smcProvider) Configure(
	ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse,
) {
	defer logging.RecoverDiags()

	// Retrieve and validate provider data from configuration and environment variables
	var config model
	diags := req.Config.Get(ctx, &config)
	logging.PanicIfDiags(diags, resp.Diagnostics)

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	token := config.Token.ValueString()
	endpoint := config.Endpoint.ValueString()
	skipVerify := config.SkipVerify.ValueBool()
	logging.PanicIf(
		"validate configurations",
		func() error {
			if token == "" {
				return errors.New("token must not be empty")
			}
			if endpoint == "" {
				return errors.New("endpoint must not be empty")
			}
			return nil
		}(),
		resp.Diagnostics,
	)

	// Create a new client using the configuration values
	client, cleanClient, err := nats.New(token, endpoint, skipVerify)
	logging.PanicIf("create a new client", err, resp.Diagnostics)
	global.Defer(cleanClient)

	// Make the smc client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *smcProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources defines the resources implemented in the provider.
func (p *smcProvider) Resources(_ context.Context) []func() resource.Resource {
	resources := lifecycle.Resources()
	return resources
}
