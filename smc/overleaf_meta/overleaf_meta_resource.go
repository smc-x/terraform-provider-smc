package meta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/teris-io/shortid"
)

// NewOverleafMetaResource is a helper function to simplify the provider implementation.
func NewOverleafMetaResource() resource.Resource {
	return &overleafMetaResource{}
}

// overleafMetaResource is the resource implementation.
type overleafMetaResource struct{}

// Metadata sets the resource type name.
func (r *overleafMetaResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_overleaf_meta"
}

// overleafMetaResourceModel maps the resource schema data.
type overleafMetaResourceModel struct {
	Owner types.String `tfsdk:"owner"`
	ID    types.String `tfsdk:"id"`
	Image types.String `tfsdk:"image"`
}

// GetSchema defines the schema for the resource.
func (r *overleafMetaResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"owner": {
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
			"image": {
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
func (r *overleafMetaResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	// Retrieve values from plan
	var plan overleafMetaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.StringValue(shortid.MustGenerate())
	plan.Image = types.StringValue(
		fmt.Sprintf("sharelatex-%s:%s", plan.Owner, shortid.MustGenerate()),
	)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *overleafMetaResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	// Do nothing
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *overleafMetaResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	// Do nothing
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *overleafMetaResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	// The state will be erased automatically
}
