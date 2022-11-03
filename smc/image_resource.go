package smc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/smc-x/terraform-provider-smc/cli"
)

func _() resource.ResourceWithConfigure {
	return &imageResource{}
}

// NewImageResource is a helper function to simplify the provider implementation.
func NewImageResource() resource.Resource {
	return &imageResource{}
}

// imageResource is the resource implementation.
type imageResource struct {
	client *cli.Cli
}

// Metadata sets the resource type name.
func (r *imageResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_image"
}

// imageResourceModel maps the resource schema data.
type imageResourceModel struct {
	Workspace types.String `tfsdk:"workspace"`
	Owner     types.String `tfsdk:"owner"`
	Base      types.String `tfsdk:"base"`
	ID        types.String `tfsdk:"id"`
}

// GetSchema defines the schema for the resource.
func (r *imageResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"workspace": {
				Type:     types.StringType,
				Required: true,
			},
			"owner": {
				Type:     types.StringType,
				Required: true,
			},
			"base": {
				Type:     types.StringType,
				Required: true,
			},
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
		},
	}, nil
}

// Create a new resource.
func (r *imageResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	tflog.Info(ctx, "create imageResource")

	// Retrieve values from plan
	var plan imageResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new tag for the image
	plan.ID = types.StringValue(
		plan.Base.ValueString() +
			"-" + plan.Owner.ValueString() + "-" + plan.Workspace.ValueString(),
	)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *imageResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	tflog.Info(ctx, "read imageResource")

	// Get current state
	var state imageResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *imageResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	tflog.Info(ctx, "update imageResource")

	// Retrieve values from plan
	var plan imageResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set updated state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *imageResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	tflog.Info(ctx, "delete imageResource")

	// Retrieve values from state
	var state imageResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing image
	// The state will be erased automatically
}

// Configure adds the provider configured client to the resource.
func (r *imageResource) Configure(
	ctx context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse,
) {
	tflog.Info(ctx, "configure imageResource")

	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*cli.Cli)
}
