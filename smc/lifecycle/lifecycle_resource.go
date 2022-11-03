package lifecycle

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/teris-io/shortid"
)

// NewLifecycleResource is a helper function to simplify the provider implementation.
func NewLifecycleResource() resource.Resource {
	return &lifecycleResource{}
}

// lifecycleResource is the resource implementation.
type lifecycleResource struct{}

// Metadata sets the resource type name.
func (r *lifecycleResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_lifecycle"
}

// lifecycleResourceModel maps the resource schema data.
type lifecycleResourceModel struct {
	Msg types.String `tfsdk:"msg"`
	ID  types.String `tfsdk:"id"`
}

// GetSchema defines the schema for the resource.
func (r *lifecycleResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"msg": {
				Type:     types.StringType,
				Required: true,
			},
			"id": {
				Type:     types.StringType,
				Computed: true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
		},
	}, nil
}

// Create a new resource.
func (r *lifecycleResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	tflog.Info(ctx, "======== create lifecycleResource")

	// Retrieve values from plan
	var plan lifecycleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.StringValue(shortid.MustGenerate())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *lifecycleResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	tflog.Info(ctx, "======== read lifecycleResource")

	// Get current state
	var state lifecycleResourceModel
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
func (r *lifecycleResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	tflog.Info(ctx, "======== update lifecycleResource")

	// Retrieve values from plan
	var plan lifecycleResourceModel
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
func (r *lifecycleResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	tflog.Info(ctx, "======== delete lifecycleResource")

	// Retrieve values from state
	var state lifecycleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The state will be erased automatically
}
