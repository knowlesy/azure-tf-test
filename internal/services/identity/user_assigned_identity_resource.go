package identity

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/peterknowles/terraform-provider-mockazurerm/internal/mockclient"
)

var _ resource.Resource = &UserAssignedIdentityResource{}
var _ resource.ResourceWithConfigure = &UserAssignedIdentityResource{}

func NewUserAssignedIdentityResource() resource.Resource {
	return &UserAssignedIdentityResource{}
}

type UserAssignedIdentityResource struct {
	client *mockclient.Client
}

type UserAssignedIdentityResourceModel struct {
	Id                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Location          types.String `tfsdk:"location"`
	ResourceGroupName types.String `tfsdk:"resource_group_name"`
	ClientId          types.String `tfsdk:"client_id"`
	PrincipalId       types.String `tfsdk:"principal_id"`
	TenantId          types.String `tfsdk:"tenant_id"`
}

func (r *UserAssignedIdentityResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_assigned_identity"
}

func (r *UserAssignedIdentityResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"client_id": schema.StringAttribute{
				Computed: true,
			},
			"principal_id": schema.StringAttribute{
				Computed: true,
			},
			"tenant_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *UserAssignedIdentityResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserAssignedIdentityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserAssignedIdentityResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subID := "00000000-0000-0000-0000-000000000000"
	mockID := mockclient.GenerateResourceID(
		subID,
		plan.ResourceGroupName.ValueString(),
		"Microsoft.ManagedIdentity",
		"userAssignedIdentities",
		plan.Name.ValueString(),
	)

	plan.Id = types.StringValue(mockID)
	plan.ClientId = types.StringValue("mock-client-id-0000-0000")
	plan.PrincipalId = types.StringValue("mock-principal-id-0000-0000")
	plan.TenantId = types.StringValue("mock-tenant-id-0000-0000")

	props := map[string]interface{}{
		"location":            plan.Location.ValueString(),
		"resource_group_name": plan.ResourceGroupName.ValueString(),
	}

	err := r.client.Save("UserAssignedIdentity", mockID, props)
	if err != nil {
		resp.Diagnostics.AddError("Error creating User Assigned Identity", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserAssignedIdentityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserAssignedIdentityResourceModel
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

func (r *UserAssignedIdentityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UserAssignedIdentityResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read existing to keep computed properties
	var state UserAssignedIdentityResourceModel
	req.State.Get(ctx, &state)
	plan.ClientId = state.ClientId
	plan.PrincipalId = state.PrincipalId
	plan.TenantId = state.TenantId

	props := map[string]interface{}{
		"location":            plan.Location.ValueString(),
		"resource_group_name": plan.ResourceGroupName.ValueString(),
	}

	err := r.client.Save("UserAssignedIdentity", plan.Id.ValueString(), props)
	if err != nil {
		resp.Diagnostics.AddError("Error updating User Assigned Identity", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserAssignedIdentityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserAssignedIdentityResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting User Assigned Identity", err.Error())
		return
	}
}
