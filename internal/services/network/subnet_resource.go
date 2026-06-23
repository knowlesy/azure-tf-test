package network

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/peterknowles/terraform-provider-mockazurerm/internal/mockclient"
)

var _ resource.Resource = &SubnetResource{}
var _ resource.ResourceWithConfigure = &SubnetResource{}

func NewSubnetResource() resource.Resource {
	return &SubnetResource{}
}

type SubnetResource struct {
	client *mockclient.Client
}

type SubnetResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	ResourceGroupName  types.String `tfsdk:"resource_group_name"`
	VirtualNetworkName types.String `tfsdk:"virtual_network_name"`
	AddressPrefixes    types.List   `tfsdk:"address_prefixes"`
}

func (r *SubnetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subnet"
}

func (r *SubnetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"resource_group_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"virtual_network_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"address_prefixes": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
		},
	}
}

func (r *SubnetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SubnetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SubnetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subID := "00000000-0000-0000-0000-000000000000"
	vnetID := mockclient.GenerateResourceID(
		subID,
		plan.ResourceGroupName.ValueString(),
		"Microsoft.Network",
		"virtualNetworks",
		plan.VirtualNetworkName.ValueString(),
	)
	mockID := mockclient.GenerateSubnetID(vnetID, plan.Name.ValueString())

	plan.Id = types.StringValue(mockID)

	props := map[string]interface{}{
		"resource_group_name":  plan.ResourceGroupName.ValueString(),
		"virtual_network_name": plan.VirtualNetworkName.ValueString(),
		"address_prefixes":     plan.AddressPrefixes.Elements(),
	}

	err := r.client.Save("Subnet", mockID, props)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Subnet", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SubnetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SubnetResourceModel
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

func (r *SubnetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SubnetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	props := map[string]interface{}{
		"resource_group_name":  plan.ResourceGroupName.ValueString(),
		"virtual_network_name": plan.VirtualNetworkName.ValueString(),
		"address_prefixes":     plan.AddressPrefixes.Elements(),
	}

	err := r.client.Save("Subnet", plan.Id.ValueString(), props)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Subnet", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SubnetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SubnetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subnetID := state.Id.ValueString()

	// Strict dependency enforcement: Can't delete if a NetworkInterface exists pointing to this Subnet.
	dependentNICs := r.client.FindResourcesByTypeAndProperty("NetworkInterface", "subnet_id", subnetID)
	if len(dependentNICs) > 0 {
		resp.Diagnostics.AddError(
			"Subnet In Use",
			fmt.Sprintf("Subnet %s is in use by %d network interface(s) (e.g., %s) and cannot be deleted.",
				subnetID, len(dependentNICs), dependentNICs[0].ID),
		)
		return
	}

	err := r.client.Delete(subnetID)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Subnet", err.Error())
		return
	}
}
