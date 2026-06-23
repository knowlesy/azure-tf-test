package compute

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/peterknowles/terraform-provider-mockazurerm/internal/mockclient"
)

var _ resource.Resource = &WindowsVirtualMachineResource{}
var _ resource.ResourceWithConfigure = &WindowsVirtualMachineResource{}

func NewWindowsVirtualMachineResource() resource.Resource {
	return &WindowsVirtualMachineResource{}
}

type WindowsVirtualMachineResource struct {
	client *mockclient.Client
}

type WindowsVirtualMachineResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	ResourceGroupName   types.String `tfsdk:"resource_group_name"`
	Location            types.String `tfsdk:"location"`
	Size                types.String `tfsdk:"size"`
	AdminUsername       types.String `tfsdk:"admin_username"`
	NetworkInterfaceIds types.List   `tfsdk:"network_interface_ids"`
	Tags                types.Map    `tfsdk:"tags"`
	OsDisk              types.List   `tfsdk:"os_disk"`
	SourceImageReference types.List  `tfsdk:"source_image_reference"`
}

type WindowsOsDiskModel struct {
	Caching              types.String `tfsdk:"caching"`
	StorageAccountType   types.String `tfsdk:"storage_account_type"`
	DiskSizeGb           types.Int64  `tfsdk:"disk_size_gb"`
}

type WindowsSourceImageReferenceModel struct {
	Publisher types.String `tfsdk:"publisher"`
	Offer     types.String `tfsdk:"offer"`
	Sku       types.String `tfsdk:"sku"`
	Version   types.String `tfsdk:"version"`
}

func (r *WindowsVirtualMachineResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_windows_virtual_machine"
}

func (r *WindowsVirtualMachineResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"size": schema.StringAttribute{
				Required: true,
			},
			"admin_username": schema.StringAttribute{
				Required: true,
			},
			"network_interface_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"tags": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"os_disk": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"caching": schema.StringAttribute{
							Required: true,
						},
						"storage_account_type": schema.StringAttribute{
							Required: true,
						},
						"disk_size_gb": schema.Int64Attribute{
							Optional: true,
						},
					},
				},
			},
			"source_image_reference": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"publisher": schema.StringAttribute{
							Required: true,
						},
						"offer": schema.StringAttribute{
							Required: true,
						},
						"sku": schema.StringAttribute{
							Required: true,
						},
						"version": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
		},
	}
}

func (r *WindowsVirtualMachineResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WindowsVirtualMachineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WindowsVirtualMachineResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subID := "00000000-0000-0000-0000-000000000000"
	mockID := mockclient.GenerateResourceID(
		subID,
		plan.ResourceGroupName.ValueString(),
		"Microsoft.Compute",
		"virtualMachines",
		plan.Name.ValueString(),
	)

	plan.Id = types.StringValue(mockID)

	props := map[string]interface{}{
		"location":               plan.Location.ValueString(),
		"resource_group_name":    plan.ResourceGroupName.ValueString(),
		"os_disk":                plan.OsDisk.Elements(),
		"source_image_reference": plan.SourceImageReference.Elements(),
	}

	err := r.client.Save("WindowsVirtualMachine", mockID, props)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Windows Virtual Machine", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WindowsVirtualMachineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WindowsVirtualMachineResourceModel
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

func (r *WindowsVirtualMachineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WindowsVirtualMachineResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	props := map[string]interface{}{
		"location":               plan.Location.ValueString(),
		"resource_group_name":    plan.ResourceGroupName.ValueString(),
		"os_disk":                plan.OsDisk.Elements(),
		"source_image_reference": plan.SourceImageReference.Elements(),
	}

	err := r.client.Save("WindowsVirtualMachine", plan.Id.ValueString(), props)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Windows Virtual Machine", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WindowsVirtualMachineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WindowsVirtualMachineResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Windows Virtual Machine", err.Error())
		return
	}
}
