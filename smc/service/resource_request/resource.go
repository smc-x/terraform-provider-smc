package request

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/terraform-community-providers/terraform-plugin-framework-utils/modifiers"

	"github.com/smc-x/terraform-provider-smc/utils/genid"
	"github.com/smc-x/terraform-provider-smc/utils/logging"
	"github.com/smc-x/terraform-provider-smc/utils/nats"
)

func _() resource.ResourceWithConfigure {
	return &request{}
}

// New is a helper function to simplify the provider implementation.
func New() resource.Resource {
	return &request{}
}

// request is the resource implementation.
type request struct {
	client *nats.Wrapper
}

// Metadata sets the resource type name.
func (r *request) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_service_request"
}

// model maps the resource schema data.
type model struct {
	ID types.String `tfsdk:"id"`

	Subj types.String  `tfsdk:"subj"`
	Time types.Float64 `tfsdk:"timeout"`

	Data types.String `tfsdk:"data"` // JSON encoded
	Resp types.String `tfsdk:"resp"` // JSON encoded
}

// GetSchema defines the schema for the resource.
func (r *request) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},

			"subj": {
				Type:     types.StringType,
				Required: true,
			},

			"timeout": {
				Type:     types.Float64Type,
				Optional: true,
				Computed: true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					modifiers.DefaultFloat(60),
				},
			},

			"data": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					modifiers.DefaultString("{}"),
				},
			},

			"resp": {
				Type:     types.StringType,
				Computed: true,
			},
		},
	}, nil
}

// Configure retrieves the shared client initialized by the provider.
func (r *request) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*nats.Wrapper)
}

// Create a new resource.
func (r *request) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	defer logging.RecoverDiags()

	// Retrieve values from plan
	var plan model
	diags := req.Plan.Get(ctx, &plan)
	logging.PanicIfDiags(diags, &resp.Diagnostics)

	subj := plan.Subj.ValueString()
	panicIfEmpty(subj, &resp.Diagnostics)

	timeout := time.Duration(plan.Time.ValueFloat64()) * time.Second
	if timeout < time.Second {
		timeout = time.Second
	}

	data := plan.Data.ValueString()

	// Generate an id for the resource
	id := genid.Short()

	// Invoke the remote method
	reply, err := r.client.Remote(subj+PatternCreate+id, []byte(data), timeout)
	logging.PanicIf(
		"invoke remote create method",
		err,
		&resp.Diagnostics,
	)
	if len(reply) == 0 || string(reply[0:1]) != "{" {
		logging.PanicIf(
			"finish the remote create method",
			errors.New(string(reply)),
			&resp.Diagnostics,
		)
	}

	// Fulfill the model
	plan.ID = types.StringValue(id)
	plan.Resp = types.StringValue(string(reply))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	logging.PanicIfDiags(diags, &resp.Diagnostics)
}

// Read refreshes the Terraform state with the latest data.
func (r *request) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	defer logging.RecoverDiags()

	// Get current state
	var state model
	diags := req.State.Get(ctx, &state)
	logging.PanicIfDiags(diags, &resp.Diagnostics)

	id := state.ID.ValueString()

	subj := state.Subj.ValueString()
	panicIfEmpty(subj, &resp.Diagnostics)

	timeout := time.Duration(state.Time.ValueFloat64()) * time.Second
	if timeout < time.Second {
		timeout = time.Second
	}

	data := state.Data.ValueString()

	// Invoke the remote method
	reply, err := r.client.Remote(subj+PatternRead+id, []byte(data), timeout)
	logging.PanicIf(
		"invoke remote read method",
		err,
		&resp.Diagnostics,
	)
	if len(reply) == 0 || string(reply[0:1]) != "{" {
		logging.PanicIf(
			"finish the remote read method",
			errors.New(string(reply)),
			&resp.Diagnostics,
		)
	}

	// Fulfill the model
	state.Resp = types.StringValue(string(reply))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	logging.PanicIfDiags(diags, &resp.Diagnostics)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *request) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	defer logging.RecoverDiags()

	// Retrieve values from plan
	var plan model
	diags := req.Plan.Get(ctx, &plan)
	logging.PanicIfDiags(diags, &resp.Diagnostics)

	id := plan.ID.ValueString()

	subj := plan.Subj.ValueString()
	panicIfEmpty(subj, &resp.Diagnostics)

	timeout := time.Duration(plan.Time.ValueFloat64()) * time.Second
	if timeout < time.Second {
		timeout = time.Second
	}

	data := plan.Data.ValueString()

	// Invoke the remote method
	reply, err := r.client.Remote(subj+PatternUpdate+id, []byte(data), timeout)
	logging.PanicIf(
		"invoke remote update method",
		err,
		&resp.Diagnostics,
	)
	if len(reply) == 0 || string(reply[0:1]) != "{" {
		logging.PanicIf(
			"finish the remote update method",
			errors.New(string(reply)),
			&resp.Diagnostics,
		)
	}

	// Fulfill the model
	plan.Resp = types.StringValue(string(reply))

	// Set updated state
	diags = resp.State.Set(ctx, plan)
	logging.PanicIfDiags(diags, &resp.Diagnostics)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *request) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	defer logging.RecoverDiags()

	// Retrieve values from state
	var state model
	diags := req.State.Get(ctx, &state)
	logging.PanicIfDiags(diags, &resp.Diagnostics)

	id := state.ID.ValueString()

	subj := state.Subj.ValueString()
	panicIfEmpty(subj, &resp.Diagnostics)

	timeout := time.Duration(state.Time.ValueFloat64()) * time.Second
	if timeout < time.Second {
		timeout = time.Second
	}

	data := state.Data.ValueString()

	// Invoke the remote method
	reply, err := r.client.Remote(subj+PatternDelete+id, []byte(data), timeout)
	logging.PanicIf(
		"invoke remote delete method",
		err,
		&resp.Diagnostics,
	)
	if len(reply) == 0 || string(reply[0:1]) != "{" {
		logging.PanicIf(
			"finish the remote delete method",
			errors.New(string(reply)),
			&resp.Diagnostics,
		)
	}

	// The state will be erased automatically
}

func panicIfEmpty(subj string, body *diag.Diagnostics) {
	logging.PanicIf(
		"validate subj",
		func() error {
			if subj == "" {
				return errors.New("subj must not be empty")
			}
			return nil
		}(),
		body,
	)
}

const (
	PatternCreate = ".create_"
	PatternRead   = ".readxx_"
	PatternUpdate = ".update_"
	PatternDelete = ".delete_"
)
