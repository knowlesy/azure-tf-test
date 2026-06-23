package storage

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/peterknowles/terraform-provider-mockazurerm/internal/mockclient"
)

var _ resource.Resource = &StorageBlobResource{}
var _ resource.ResourceWithConfigure = &StorageBlobResource{}

func NewStorageBlobResource() resource.Resource {
	return &StorageBlobResource{}
}

type StorageBlobResource struct {
	client *mockclient.Client
}

type StorageBlobResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	StorageAccountName   types.String `tfsdk:"storage_account_name"`
	StorageContainerName types.String `tfsdk:"storage_container_name"`
	Type                 types.String `tfsdk:"type"`
	Source               types.String `tfsdk:"source"`
}

func (r *StorageBlobResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_blob"
}

func (r *StorageBlobResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"storage_account_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"storage_container_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required: true,
			},
			"source": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (r *StorageBlobResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *StorageBlobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan StorageBlobResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mockID := "https://" + plan.StorageAccountName.ValueString() + ".blob.core.windows.net/" + plan.StorageContainerName.ValueString() + "/" + plan.Name.ValueString()
	plan.Id = types.StringValue(mockID)

	props := map[string]interface{}{
		"storage_account_name":   plan.StorageAccountName.ValueString(),
		"storage_container_name": plan.StorageContainerName.ValueString(),
	}

	err := r.client.Save("StorageBlob", mockID, props)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Storage Blob", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StorageBlobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state StorageBlobResourceModel
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

func (r *StorageBlobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan StorageBlobResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	props := map[string]interface{}{
		"storage_account_name":   plan.StorageAccountName.ValueString(),
		"storage_container_name": plan.StorageContainerName.ValueString(),
	}

	err := r.client.Save("StorageBlob", plan.Id.ValueString(), props)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Storage Blob", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StorageBlobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state StorageBlobResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Storage Blob", err.Error())
		return
	}
}
