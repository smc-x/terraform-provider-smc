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

// NewLifecycleFragileResource is a helper function to simplify the provider implementation.
func NewLifecycleFragileResource() resource.Resource {
	return &lifecycleFragileResource{}
}

// lifecycleFragileResource is the resource implementation.
type lifecycleFragileResource struct{}

// Metadata sets the resource type name.
func (r *lifecycleFragileResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_lifecycle_fragile"
}

// lifecycleFragileResourceModel maps the resource schema data.
type lifecycleFragileResourceModel struct {
	Msg types.String `tfsdk:"msg"`
	ID  types.String `tfsdk:"id"`
}

// GetSchema defines the schema for the resource.
func (r *lifecycleFragileResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"msg": {
				Type:     types.StringType,
				Required: true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.RequiresReplace(),
				},
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
func (r *lifecycleFragileResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	tflog.Info(ctx, "======== create lifecycleFragileResource")

	// Retrieve values from plan
	var plan lifecycleFragileResourceModel
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
func (r *lifecycleFragileResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	tflog.Info(ctx, "======== read lifecycleFragileResource")

	// Get current state
	var state lifecycleFragileResourceModel
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
func (r *lifecycleFragileResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	tflog.Info(ctx, "======== update lifecycleFragileResource")

	// Retrieve values from plan
	var plan lifecycleFragileResourceModel
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
func (r *lifecycleFragileResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	tflog.Info(ctx, "======== delete lifecycleFragileResource")

	// Retrieve values from state
	var state lifecycleFragileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The state will be erased automatically
}
