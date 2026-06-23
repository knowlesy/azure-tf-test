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

var _ resource.Resource = &KubernetesClusterResource{}
var _ resource.ResourceWithConfigure = &KubernetesClusterResource{}

func NewKubernetesClusterResource() resource.Resource {
	return &KubernetesClusterResource{}
}

type KubernetesClusterResource struct {
	client *mockclient.Client
}

type KubernetesClusterResourceModel struct {
	Id                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Location          types.String `tfsdk:"location"`
	ResourceGroupName types.String `tfsdk:"resource_group_name"`
	DnsPrefix         types.String `tfsdk:"dns_prefix"`
	Tags              types.Map    `tfsdk:"tags"`
	DefaultNodePool   types.List   `tfsdk:"default_node_pool"`
	Identity          types.List   `tfsdk:"identity"`
	NetworkProfile    types.List   `tfsdk:"network_profile"`
}

type DefaultNodePoolModel struct {
	Name         types.String `tfsdk:"name"`
	NodeCount    types.Int64  `tfsdk:"node_count"`
	VmSize       types.String `tfsdk:"vm_size"`
	VnetSubnetId types.String `tfsdk:"vnet_subnet_id"`
	Type         types.String `tfsdk:"type"`
}

type IdentityModel struct {
	Type types.String `tfsdk:"type"`
}

type NetworkProfileModel struct {
	NetworkPlugin types.String `tfsdk:"network_plugin"`
	NetworkPolicy types.String `tfsdk:"network_policy"`
	DnsServiceIp  types.String `tfsdk:"dns_service_ip"`
	ServiceCidr   types.String `tfsdk:"service_cidr"`
}

func (r *KubernetesClusterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_cluster"
}

func (r *KubernetesClusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"dns_prefix": schema.StringAttribute{
				Required: true,
			},
			"tags": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"default_node_pool": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
						"node_count": schema.Int64Attribute{
							Required: true,
						},
						"vm_size": schema.StringAttribute{
							Required: true,
						},
						"vnet_subnet_id": schema.StringAttribute{
							Optional: true,
						},
						"type": schema.StringAttribute{
							Optional: true,
						},
					},
				},
			},
			"identity": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
			"network_profile": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"network_plugin": schema.StringAttribute{
							Required: true,
						},
						"network_policy": schema.StringAttribute{
							Optional: true,
						},
						"dns_service_ip": schema.StringAttribute{
							Optional: true,
						},
						"service_cidr": schema.StringAttribute{
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func (r *KubernetesClusterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *KubernetesClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan KubernetesClusterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subID := "00000000-0000-0000-0000-000000000000"
	mockID := mockclient.GenerateResourceID(
		subID,
		plan.ResourceGroupName.ValueString(),
		"Microsoft.ContainerService",
		"managedClusters",
		plan.Name.ValueString(),
	)

	plan.Id = types.StringValue(mockID)

	props := map[string]interface{}{
		"location":            plan.Location.ValueString(),
		"resource_group_name": plan.ResourceGroupName.ValueString(),
		"default_node_pool":   plan.DefaultNodePool.Elements(),
		"identity":            plan.Identity.Elements(),
		"network_profile":     plan.NetworkProfile.Elements(),
	}

	err := r.client.Save("KubernetesCluster", mockID, props)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Kubernetes Cluster", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *KubernetesClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state KubernetesClusterResourceModel
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

func (r *KubernetesClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan KubernetesClusterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	props := map[string]interface{}{
		"location":            plan.Location.ValueString(),
		"resource_group_name": plan.ResourceGroupName.ValueString(),
		"default_node_pool":   plan.DefaultNodePool.Elements(),
		"identity":            plan.Identity.Elements(),
		"network_profile":     plan.NetworkProfile.Elements(),
	}

	err := r.client.Save("KubernetesCluster", plan.Id.ValueString(), props)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Kubernetes Cluster", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *KubernetesClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state KubernetesClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Kubernetes Cluster", err.Error())
		return
	}
}
