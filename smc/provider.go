package smc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/smc-x/terraform-provider-smc/cli"
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

// GetSchema defines the provider-level schema for configuration data.
func (p *smcProvider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{},
	}, nil
}

// Configure prepares an smc api client for data sources and resources.
func (p *smcProvider) Configure(
	ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse,
) {
	tflog.Info(ctx, "configure smc client")

	// Retrieve and validate provider data from configuration and environment variables
	// TODO: we don't have provider-level configuration by far
	// var config smcProviderModel

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	// Create a new smc client using the configuration values
	client, err := cli.NewCli()
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to create smc api client",
			"client error: "+err.Error(),
		)
		return
	}

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
	return []func() resource.Resource{
		NewImageResource,
	}
}
