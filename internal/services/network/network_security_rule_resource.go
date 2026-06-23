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

var _ resource.Resource = &NetworkSecurityRuleResource{}
var _ resource.ResourceWithConfigure = &NetworkSecurityRuleResource{}

func NewNetworkSecurityRuleResource() resource.Resource {
	return &NetworkSecurityRuleResource{}
}

type NetworkSecurityRuleResource struct {
	client *mockclient.Client
}

type NetworkSecurityRuleResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	ResourceGroupName         types.String `tfsdk:"resource_group_name"`
	NetworkSecurityGroupName  types.String `tfsdk:"network_security_group_name"`
	Priority                  types.Int64  `tfsdk:"priority"`
	Direction                 types.String `tfsdk:"direction"`
	Access                    types.String `tfsdk:"access"`
	Protocol                  types.String `tfsdk:"protocol"`
	SourcePortRange           types.String `tfsdk:"source_port_range"`
	DestinationPortRange      types.String `tfsdk:"destination_port_range"`
	SourceAddressPrefix       types.String `tfsdk:"source_address_prefix"`
	DestinationAddressPrefix  types.String `tfsdk:"destination_address_prefix"`
}

func (r *NetworkSecurityRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_security_rule"
}

func (r *NetworkSecurityRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"network_security_group_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"priority": schema.Int64Attribute{
				Required: true,
			},
			"direction": schema.StringAttribute{
				Required: true,
			},
			"access": schema.StringAttribute{
				Required: true,
			},
			"protocol": schema.StringAttribute{
				Required: true,
			},
			"source_port_range": schema.StringAttribute{
				Required: true,
			},
			"destination_port_range": schema.StringAttribute{
				Required: true,
			},
			"source_address_prefix": schema.StringAttribute{
				Required: true,
			},
			"destination_address_prefix": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *NetworkSecurityRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworkSecurityRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NetworkSecurityRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subID := "00000000-0000-0000-0000-000000000000"
	nsgID := mockclient.GenerateResourceID(
		subID,
		plan.ResourceGroupName.ValueString(),
		"Microsoft.Network",
		"networkSecurityGroups",
		plan.NetworkSecurityGroupName.ValueString(),
	)
	mockID := nsgID + "/securityRules/" + plan.Name.ValueString()

	plan.Id = types.StringValue(mockID)

	props := map[string]interface{}{
		"resource_group_name":         plan.ResourceGroupName.ValueString(),
		"network_security_group_name": plan.NetworkSecurityGroupName.ValueString(),
	}

	err := r.client.Save("NetworkSecurityRule", mockID, props)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Network Security Rule", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NetworkSecurityRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NetworkSecurityRuleResourceModel
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

func (r *NetworkSecurityRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NetworkSecurityRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	props := map[string]interface{}{
		"resource_group_name":         plan.ResourceGroupName.ValueString(),
		"network_security_group_name": plan.NetworkSecurityGroupName.ValueString(),
	}

	err := r.client.Save("NetworkSecurityRule", plan.Id.ValueString(), props)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Network Security Rule", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NetworkSecurityRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state NetworkSecurityRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Network Security Rule", err.Error())
		return
	}
}
