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

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &VirtualNetworkResource{}
var _ resource.ResourceWithConfigure = &VirtualNetworkResource{}

func NewVirtualNetworkResource() resource.Resource {
	return &VirtualNetworkResource{}
}

type VirtualNetworkResource struct {
	client *mockclient.Client
}

type VirtualNetworkResourceModel struct {
	Id                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Location          types.String `tfsdk:"location"`
	ResourceGroupName types.String `tfsdk:"resource_group_name"`
	AddressSpace      types.List   `tfsdk:"address_space"`
	Tags              types.Map    `tfsdk:"tags"`
	DnsServers        types.List   `tfsdk:"dns_servers"`
}

func (r *VirtualNetworkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_network"
}

func (r *VirtualNetworkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"location": schema.StringAttribute{
				Required: true,
			},
			"address_space": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"tags": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"dns_servers": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (r *VirtualNetworkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*mockclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *mockclient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *VirtualNetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan VirtualNetworkResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subID := "00000000-0000-0000-0000-000000000000" // Mock subscription
	mockID := mockclient.GenerateResourceID(
		subID,
		plan.ResourceGroupName.ValueString(),
		"Microsoft.Network",
		"virtualNetworks",
		plan.Name.ValueString(),
	)

	plan.Id = types.StringValue(mockID)

	// Persist to mock DB
	props := map[string]interface{}{
		"location":            plan.Location.ValueString(),
		"resource_group_name": plan.ResourceGroupName.ValueString(),
		"address_space":       plan.AddressSpace.Elements(),
		"tags":                plan.Tags.Elements(),
		"dns_servers":         plan.DnsServers.Elements(),
	}
	err := r.client.Save("VirtualNetwork", mockID, props)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Virtual Network in Mock DB", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VirtualNetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VirtualNetworkResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, exists := r.client.Read(state.Id.ValueString())
	if !exists {
		// Resource doesn't exist in mock DB, remove from state
		resp.State.RemoveResource(ctx)
		return
	}

	// In a complete implementation we would map properties from 'res' back to 'state'
	// For this mock, we assume the state is accurate if it exists
	_ = res

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *VirtualNetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan VirtualNetworkResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mockID := plan.Id.ValueString()

	props := map[string]interface{}{
		"location":            plan.Location.ValueString(),
		"resource_group_name": plan.ResourceGroupName.ValueString(),
		"address_space":       plan.AddressSpace.Elements(),
		"tags":                plan.Tags.Elements(),
		"dns_servers":         plan.DnsServers.Elements(),
	}
	err := r.client.Save("VirtualNetwork", mockID, props)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Virtual Network in Mock DB", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VirtualNetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VirtualNetworkResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic check: Are there any subnets dependent on this VNet?
	subnets := r.client.FindResourcesByTypeAndProperty("Subnet", "virtual_network_name", state.Name.ValueString())
	if len(subnets) > 0 {
		resp.Diagnostics.AddError("Delete Blocked", fmt.Sprintf("Virtual Network %s cannot be deleted because it contains %d subnet(s).", state.Name.ValueString(), len(subnets)))
		return
	}

	err := r.client.Delete(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Virtual Network from Mock DB", err.Error())
		return
	}
}
