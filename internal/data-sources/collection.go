package data_sources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	services "github.com/thiskevinwang/terraform-provider-pinecone/internal/services"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource = &CollectionDataSource{}
)

func NewCollectionDataSource() datasource.DataSource {
	return &CollectionDataSource{}
}

// CollectionDataSource defines the data source implementation.
type CollectionDataSource struct {
	client services.Pinecone
}

// CollectionDataSourceModel describes the data source data model.
type CollectionDataSourceModel struct {
	Name      types.String `tfsdk:"name"`
	Dimension types.Int64  `tfsdk:"dimension"`
	Id        types.String `tfsdk:"id"`
}

func (d *CollectionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collection"
}

func (d *CollectionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: `A Pinecone collection
- See [Understanding collections](https://docs.pinecone.io/docs/collections)
- See [API Docs](https://docs.pinecone.io/reference/describe_collection)
`,

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the collection",
				Required:            true,
			},
			"dimension": schema.Int64Attribute{
				MarkdownDescription: "The dimension of the collection",
				Required:            false,
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the datasource
func (d *CollectionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	// extract the client from the provider data
	client, ok := req.ProviderData.(services.Pinecone)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected pinecone.Pinecone, got: %T", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *CollectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CollectionDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from configuration
	name := data.Name.ValueString()
	response, err := d.client.DescribeCollection(name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to describe collection",
			fmt.Sprintf("Failed to describe collection: %s", err),
		)
		return
	}

	// log the response
	tflog.Info(ctx, "DescribeCollection OK", map[string]any{"respond": *response})

	data.Id = types.StringValue(fmt.Sprintf("datasource-pinecone_collection-%s/%s", d.client.Environment, name))
	data.Name = types.StringValue(response.Name)
	data.Dimension = types.Int64Value(response.Dimension)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
