package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	services "github.com/thiskevinwang/terraform-provider-pinecone/internal/services"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &indexResource{}
	_ resource.ResourceWithConfigure   = &indexResource{}
	_ resource.ResourceWithImportState = &indexResource{}
)

func NewIndexResource() resource.Resource {
	return &indexResource{}
}

// indexResource is the resource implementation.
type indexResource struct {
	// this client is set by the provider
	client services.Pinecone
}

// {
//   "database": {
//     "name": "test",
//     "metric": "cosine",
//     "dimension": 1536,
//     "replicas": 1,
//     "shards": 1,
//     "pods": 1
//   },
//   "status": {
//     "waiting": [],
//     "crashed": [],
//     "host": "test-260030d.svc.us-west4-gcp-free.pinecone.io",
//     "port": 433,
//     "state": "Ready",
//     "ready": true
//   }
// }

// indexResourceModel maps the resource schema data.
// - "github.com/hashicorp/terraform-plugin-framework/types"
type indexResourceModel struct {
	Name      types.String `tfsdk:"name"`
	Dimension types.Int64  `tfsdk:"dimension"`
	Metric    types.String `tfsdk:"metric"`
	Replicas  types.Int64  `tfsdk:"replicas"`
	Pods      types.Int64  `tfsdk:"pods"`
	// Shards    types.Number  `tfsdk:"shards"`
}

// Metadata returns the resource type name.
func (r *indexResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_index"
}

// Schema defines the schema for the resource.
func (r *indexResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an index.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the index to be created. The maximum length is 45 characters.",
				Required:    true,
			},
			"dimension": schema.Int64Attribute{
				Description: "The dimensions of the vectors to be inserted in the index",
				Required:    true,
			},
			"metric": schema.StringAttribute{
				Description: "The distance metric to be used for similarity search. You can use 'euclidean', 'cosine', or 'dotproduct'.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("cosine"),
			},
			"replicas": schema.Int64Attribute{
				Description: "The number of replicas. Replicas duplicate your index. They provide higher availability and throughput.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1),
			},
			"pods": schema.Int64Attribute{
				Description: "The number of pods for the index to use,including replicas.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1),
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *indexResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

// Create a new resource.
func (r *indexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Error(ctx, "creating a resource")
	var plan indexResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	name := plan.Name.ValueString()
	dimension := plan.Dimension.ValueInt64()
	metric := plan.Metric.ValueString()

	// Create new index
	response, err := r.client.CreateIndex(name, dimension, metric)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create index",
			fmt.Sprintf("Failed to create index: %s", err),
		)
		return
	}

	// log the response
	tflog.Info(ctx, "CreateIndex OK: %s", map[string]any{"response": *response})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read resource information.
func (r *indexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read data from Terraform state
	var plan indexResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	name := plan.Name.ValueString()
	response, err := r.client.DescribeIndex(name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to describe index",
			fmt.Sprintf("Failed to describe index: %s", err),
		)
		return
	}

	// log the response
	tflog.Info(ctx, "DescribeIndex OK: %s", map[string]any{"response": *response})

	// Save data into Terraform plan
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *indexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data indexResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *indexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data indexResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *indexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
