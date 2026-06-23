package provider

import (
	"context"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/peterknowles/terraform-provider-mockazurerm/internal/mockclient"
	"github.com/peterknowles/terraform-provider-mockazurerm/internal/services/compute"
	"github.com/peterknowles/terraform-provider-mockazurerm/internal/services/foundational"
	"github.com/peterknowles/terraform-provider-mockazurerm/internal/services/identity"
	"github.com/peterknowles/terraform-provider-mockazurerm/internal/services/messaging"
	"github.com/peterknowles/terraform-provider-mockazurerm/internal/services/network"
	"github.com/peterknowles/terraform-provider-mockazurerm/internal/services/storage"
)

var _ provider.Provider = &mockAzureProvider{}

type mockAzureProvider struct {
	version string
}

type mockAzureProviderModel struct {
	SubscriptionId types.String `tfsdk:"subscription_id"`
	TenantId       types.String `tfsdk:"tenant_id"`
	ClientId       types.String `tfsdk:"client_id"`
	ClientSecret   types.String `tfsdk:"client_secret"`
	DbPath         types.String `tfsdk:"db_path"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &mockAzureProvider{
			version: version,
		}
	}
}

func (p *mockAzureProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "mockazurerm"
	resp.Version = p.version
}

func (p *mockAzureProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"subscription_id": schema.StringAttribute{
				Optional:    true,
				Description: "The Subscription ID which should be used.",
			},
			"tenant_id": schema.StringAttribute{
				Optional:    true,
				Description: "The Tenant ID which should be used.",
			},
			"client_id": schema.StringAttribute{
				Optional:    true,
				Description: "The Client ID which should be used.",
			},
			"client_secret": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"db_path": schema.StringAttribute{
				Optional:    true,
				Description: "Path to the local JSON mock DB. Defaults to .mockazurerm_db.json in the current dir.",
			},
		},
	}
}

func (p *mockAzureProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Mock Azure Provider")

	var data mockAzureProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dbPath := ".mockazurerm_db.json"
	if !data.DbPath.IsNull() && !data.DbPath.IsUnknown() {
		dbPath = data.DbPath.ValueString()
	}
	dbPath, _ = filepath.Abs(dbPath)

	client, err := mockclient.NewClient(dbPath)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to load Mock DB",
			"An unexpected error occurred when reading the mock db: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Mock Azure Provider configured")
}

func (p *mockAzureProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Foundational
		foundational.NewResourceGroupResource,
		foundational.NewKeyVaultResource,
		foundational.NewKeyVaultSecretResource,
		foundational.NewLogAnalyticsWorkspaceResource,

		// Network
		network.NewVirtualNetworkResource,
		network.NewSubnetResource,
		network.NewNetworkInterfaceResource,
		network.NewNetworkSecurityGroupResource,
		network.NewNetworkSecurityRuleResource,
		network.NewRouteTableResource,
		network.NewRouteResource,
		network.NewNatGatewayResource,
		network.NewSubnetNatGatewayAssociationResource,
		network.NewPrivateEndpointResource,
		network.NewPrivateLinkServiceResource,
		network.NewSubnetServiceEndpointStoragePolicyResource,

		// Compute
		compute.NewKubernetesClusterResource,
		compute.NewLinuxVirtualMachineResource,
		compute.NewWindowsVirtualMachineResource,
		compute.NewVirtualMachineExtensionResource,

		// Storage
		storage.NewStorageAccountResource,
		storage.NewStorageContainerResource,
		storage.NewStorageBlobResource,

		// Messaging
		messaging.NewServiceBusNamespaceResource,
		messaging.NewServiceBusQueueResource,
		messaging.NewEventHubNamespaceResource,
		messaging.NewEventHubResource,

		// Identity
		identity.NewUserAssignedIdentityResource,
		identity.NewServicePrincipalResource,
		identity.NewRoleAssignmentResource,
	}
}

func (p *mockAzureProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
