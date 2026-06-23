package network

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/peterknowles/terraform-provider-mockazurerm/internal/mockclient"
)

var _ resource.Resource = &PrivateEndpointResource{}
var _ resource.ResourceWithConfigure = &PrivateEndpointResource{}

func NewPrivateEndpointResource() resource.Resource {
	return &PrivateEndpointResource{}
}

type PrivateEndpointResource struct {
	client *mockclient.Client
}

type PrivateEndpointResourceModel struct {
	Id                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Location          types.String `tfsdk:"location"`
	ResourceGroupName types.String `tfsdk:"resource_group_name"`
	SubnetId          types.String `tfsdk:"subnet_id"`
}

func (r *PrivateEndpointResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_endpoint"
}

func (r *PrivateEndpointResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"location": schema.StringAttribute{
				Required: true,
			},
			"resource_group_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *PrivateEndpointResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*mockclient.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", "")
		return
	}
	r.client = client
}

func (r *PrivateEndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PrivateEndpointResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subID := "00000000-0000-0000-0000-000000000000"
	mockID := mockclient.GenerateResourceID(
		subID,
		plan.ResourceGroupName.ValueString(),
		"Microsoft.Network",
		"privateEndpoints",
		plan.Name.ValueString(),
	)

	plan.Id = types.StringValue(mockID)

	props := map[string]interface{}{
		"location":            plan.Location.ValueString(),
		"resource_group_name": plan.ResourceGroupName.ValueString(),
		"subnet_id":           plan.SubnetId.ValueString(),
	}

	err := r.client.Save("PrivateEndpoint", mockID, props)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Private Endpoint", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PrivateEndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PrivateEndpointResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, exists := r.client.Read(state.Id.ValueString())
	if !exists {
		resp.State.RemoveResource(ctx)
		return
	}

	_ = res
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PrivateEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PrivateEndpointResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	props := map[string]interface{}{
		"location":            plan.Location.ValueString(),
		"resource_group_name": plan.ResourceGroupName.ValueString(),
		"subnet_id":           plan.SubnetId.ValueString(),
	}

	err := r.client.Save("PrivateEndpoint", plan.Id.ValueString(), props)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Private Endpoint", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PrivateEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PrivateEndpointResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Private Endpoint", err.Error())
		return
	}
}
