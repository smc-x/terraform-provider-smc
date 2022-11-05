package fragile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/smc-x/terraform-provider-smc/utils/genid"
	"github.com/smc-x/terraform-provider-smc/utils/logging"
)

// New is a helper function to simplify the provider implementation.
func New() resource.Resource {
	return &fragile{}
}

// fragile is the resource implementation.
type fragile struct{}

// Metadata sets the resource type name.
func (r *fragile) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_lifecycle_fragile"
}

// model maps the resource schema data.
type model struct {
	Msg types.String `tfsdk:"msg"`
	ID  types.String `tfsdk:"id"`
}

// GetSchema defines the schema for the resource.
func (r *fragile) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
func (r *fragile) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	defer logging.RecoverDiags()

	// Retrieve values from plan
	var plan model
	diags := req.Plan.Get(ctx, &plan)
	logging.PanicIfDiags(diags, resp.Diagnostics)

	plan.ID = types.StringValue(genid.Short())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	logging.PanicIfDiags(diags, resp.Diagnostics)
}

// Read refreshes the Terraform state with the latest data.
func (r *fragile) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	defer logging.RecoverDiags()

	// Get current state
	var state model
	diags := req.State.Get(ctx, &state)
	logging.PanicIfDiags(diags, resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	logging.PanicIfDiags(diags, resp.Diagnostics)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *fragile) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	defer logging.RecoverDiags()

	// Retrieve values from plan
	var plan model
	diags := req.Plan.Get(ctx, &plan)
	logging.PanicIfDiags(diags, resp.Diagnostics)

	// Set updated state
	diags = resp.State.Set(ctx, plan)
	logging.PanicIfDiags(diags, resp.Diagnostics)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *fragile) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	defer logging.RecoverDiags()

	// Retrieve values from state
	var state model
	diags := req.State.Get(ctx, &state)
	logging.PanicIfDiags(diags, resp.Diagnostics)

	// The state will be erased automatically
}
