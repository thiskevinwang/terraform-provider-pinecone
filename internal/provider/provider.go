package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	resources "github.com/thiskevinwang/terraform-provider-pinecone/internal/resources"
	services "github.com/thiskevinwang/terraform-provider-pinecone/internal/services"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &pineconeProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &pineconeProvider{
			version: version,
		}
	}
}

type pineconeProvider struct {
	version string
}

// pineconeProviderModel maps provider schema data to a Go type.
type pineconeProviderModel struct {
	// ex. uuid
	ApiKey types.String `tfsdk:"apikey"`
	// ex. us-west4-gcp-free
	Environment types.String `tfsdk:"environment"`
}

// Metadata returns the provider type name.
func (p *pineconeProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pinecone"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *pineconeProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Pinecone.io",
		Attributes: map[string]schema.Attribute{
			"apikey": schema.StringAttribute{
				Description: "...Or PINECONE_API_KEY",
				Optional:    false,
				Required:    true,
				Sensitive:   true,
			},
			"environment": schema.StringAttribute{
				Description: "...Or PINECONE_ENVIRONMENT",
				Optional:    false,
				Required:    true,
			},
		},
	}
}

func (p *pineconeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring client")

	// Retrieve provider data from configuration
	var config pineconeProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"Summary (ApiKey)",
			"Detail (ApiKey)",
		)
	}

	if config.Environment.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("environment"),
			"Summary (Environment)",
			"Detail (Environment)",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	apikey := os.Getenv("PINECONE_API_KEY")
	environment := os.Getenv("PINECONE_ENVIRONMENT")

	if !config.ApiKey.IsNull() {
		apikey = config.ApiKey.ValueString()
	}

	if !config.Environment.IsNull() {
		environment = config.Environment.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if apikey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"Missing (ApiKey)",
			"Detail (ApiKey)",
		)
	}

	if environment == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("environment"),
			"Missing (Environment)",
			"Detail (Environment)",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "pinecone_api_key", apikey)
	ctx = tflog.SetField(ctx, "pinecone_environment", environment)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "pinecone_api_key")

	tflog.Debug(ctx, "Creating client")

	// TODO(kevinwang): create pinecone client?
	client := services.Pinecone{
		ApiKey:      apikey,
		Environment: environment,
	}

	// TODO(kevinwang): Make the client available during DataSource and Resource type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *pineconeProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources defines the resources implemented in the provider.
func (p *pineconeProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewIndexResource,
	}
}
