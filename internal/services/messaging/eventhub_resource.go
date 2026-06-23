package messaging

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/peterknowles/terraform-provider-mockazurerm/internal/mockclient"
)

var _ resource.Resource = &EventHubResource{}
var _ resource.ResourceWithConfigure = &EventHubResource{}

func NewEventHubResource() resource.Resource {
	return &EventHubResource{}
}

type EventHubResource struct {
	client *mockclient.Client
}

type EventHubResourceModel struct {
	Id                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	NamespaceName     types.String `tfsdk:"namespace_name"`
	ResourceGroupName types.String `tfsdk:"resource_group_name"`
	PartitionCount    types.Int64  `tfsdk:"partition_count"`
	MessageRetention  types.Int64  `tfsdk:"message_retention"`
}

func (r *EventHubResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_eventhub"
}

func (r *EventHubResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"namespace_name": schema.StringAttribute{
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
			"partition_count": schema.Int64Attribute{
				Required: true,
			},
			"message_retention": schema.Int64Attribute{
				Required: true,
			},
		},
	}
}

func (r *EventHubResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EventHubResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EventHubResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subID := "00000000-0000-0000-0000-000000000000"
	namespaceID := mockclient.GenerateResourceID(
		subID,
		plan.ResourceGroupName.ValueString(),
		"Microsoft.EventHub",
		"namespaces",
		plan.NamespaceName.ValueString(),
	)
	mockID := namespaceID + "/eventhubs/" + plan.Name.ValueString()

	plan.Id = types.StringValue(mockID)

	props := map[string]interface{}{
		"namespace_name":      plan.NamespaceName.ValueString(),
		"resource_group_name": plan.ResourceGroupName.ValueString(),
	}

	err := r.client.Save("EventHub", mockID, props)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Event Hub", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EventHubResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EventHubResourceModel
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

func (r *EventHubResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EventHubResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	props := map[string]interface{}{
		"namespace_name":      plan.NamespaceName.ValueString(),
		"resource_group_name": plan.ResourceGroupName.ValueString(),
	}

	err := r.client.Save("EventHub", plan.Id.ValueString(), props)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Event Hub", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EventHubResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EventHubResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Event Hub", err.Error())
		return
	}
}
