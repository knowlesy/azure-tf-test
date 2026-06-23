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

var _ resource.Resource = &SubnetNatGatewayAssociationResource{}
var _ resource.ResourceWithConfigure = &SubnetNatGatewayAssociationResource{}

func NewSubnetNatGatewayAssociationResource() resource.Resource {
	return &SubnetNatGatewayAssociationResource{}
}

type SubnetNatGatewayAssociationResource struct {
	client *mockclient.Client
}

type SubnetNatGatewayAssociationResourceModel struct {
	Id           types.String `tfsdk:"id"`
	SubnetId     types.String `tfsdk:"subnet_id"`
	NatGatewayId types.String `tfsdk:"nat_gateway_id"`
}

func (r *SubnetNatGatewayAssociationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subnet_nat_gateway_association"
}

func (r *SubnetNatGatewayAssociationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"nat_gateway_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *SubnetNatGatewayAssociationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SubnetNatGatewayAssociationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SubnetNatGatewayAssociationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mockID := plan.SubnetId.ValueString() + "/natGatewayAssociations/" + plan.NatGatewayId.ValueString()
	plan.Id = types.StringValue(mockID)

	props := map[string]interface{}{
		"subnet_id":      plan.SubnetId.ValueString(),
		"nat_gateway_id": plan.NatGatewayId.ValueString(),
	}

	err := r.client.Save("SubnetNatGatewayAssociation", mockID, props)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Association", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SubnetNatGatewayAssociationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SubnetNatGatewayAssociationResourceModel
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

func (r *SubnetNatGatewayAssociationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SubnetNatGatewayAssociationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	props := map[string]interface{}{
		"subnet_id":      plan.SubnetId.ValueString(),
		"nat_gateway_id": plan.NatGatewayId.ValueString(),
	}

	err := r.client.Save("SubnetNatGatewayAssociation", plan.Id.ValueString(), props)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Association", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SubnetNatGatewayAssociationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SubnetNatGatewayAssociationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Association", err.Error())
		return
	}
}
