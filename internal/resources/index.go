package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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

// indexResourceModel maps the resource schema data.
// - "github.com/hashicorp/terraform-plugin-framework/types"
type indexResourceModel struct {
	Id        types.String `tfsdk:"id"` // for TF
	Name      types.String `tfsdk:"name"`
	Dimension types.Int64  `tfsdk:"dimension"`
	Metric    types.String `tfsdk:"metric"`
	Replicas  types.Int64  `tfsdk:"replicas"`
	Pods      types.Int64  `tfsdk:"pods"`
	// Shards    types.Number  `tfsdk:"shards"`
	SourceCollection types.String `tfsdk:"source_collection"`
}

// Metadata returns the resource type name.
func (r *indexResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Debug(ctx, "indexResource.Metadata", map[string]any{"req": req, "resp": resp})

	resp.TypeName = req.ProviderTypeName + "_index"
}

// Schema defines the schema for the resource.
func (r *indexResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Debug(ctx, "indexResource.Schema", map[string]any{"req": req, "resp": resp})

	resp.Schema = schema.Schema{
		Description: "Manages an index.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Service generated identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
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
			"source_collection": schema.StringAttribute{
				Description: "The name of the collection to create an index from",
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *indexResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "indexResource.Configure", map[string]any{"req": req, "resp": resp})
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
	tflog.Debug(ctx, "indexResource.Create", map[string]any{"req": req, "resp": resp})
	var plan indexResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	name := plan.Name.ValueString()
	dimension := plan.Dimension.ValueInt64()
	metric := plan.Metric.ValueString()
	sourceCollection := plan.SourceCollection.ValueString()

	// Create new index
	response, err := r.client.CreateIndex(services.CreateIndexBodyParams{
		Name:             name,
		Dimension:        dimension,
		Metric:           metric,
		SourceCollection: sourceCollection,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create index",
			fmt.Sprintf("Failed to create index: %s", err),
		)
		return
	}
	// poll the describe index endpoint until the index is ready
	// Poll every n seconds
	ticker := time.NewTicker(10 * time.Second)
	shouldRetry := true
	for shouldRetry {
		diRes, err := r.client.DescribeIndex(name)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to poll index",
				fmt.Sprintf("Failed to describe index: %s", err),
			)
			return
		}

		if diRes.Status.State == "Ready" || diRes.Status.Ready {
			shouldRetry = false
		} else {
			<-ticker.C // keep polling
		}
	}

	// log the response
	tflog.Info(ctx, "CreateIndex OK: %s", map[string]any{"response": *response})

	plan.Id = types.StringValue(fmt.Sprintf("%s/%s", r.client.Environment, name))

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read resource information.
func (r *indexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "indexResource.Read", map[string]any{"req": req, "resp": resp})

	// Get current state
	// Read data from Terraform state
	var state indexResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get fresh state from Pinecone
	// Generate API request body from plan
	name := state.Name.ValueString()
	response, err := r.client.DescribeIndex(name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to describe index",
			err.Error(),
		)
		return
	}

	// log the response
	tflog.Info(ctx, "DescribeIndex OK", map[string]any{"response": *response})

	// Set refreshed state
	state.Name = types.StringValue(response.Database.Name)
	state.Dimension = types.Int64Value(response.Database.Dimension)
	state.Metric = types.StringValue(response.Database.Metric)
	state.Replicas = types.Int64Value(response.Database.Replicas)
	state.Pods = types.Int64Value(response.Database.Pods)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource information.
func (r *indexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "indexResource.Update", map[string]any{"req": req, "resp": resp})

	var plan indexResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	indexName := plan.Name.ValueString()

	confIdxResp, err := r.client.ConfigureIndex(indexName, &services.ConfigureIndexRequest{})
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update index",
			fmt.Sprintf("Failed to update index: %s", err),
		)
		return
	}

	// log the response
	tflog.Info(ctx, "ConfigureIndex OK", map[string]any{"response": *confIdxResp})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete resource information.
func (r *indexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "indexResource.Delete", map[string]any{"req": req, "resp": resp})

	var state indexResourceModel

	// Read Terraform plan data into the model
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	delIdxResp, err := r.client.DeleteIndex(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete index",
			fmt.Sprintf("Failed to delete index: %s", err),
		)
		return
	}

	// log the response
	tflog.Info(ctx, "DeleteIndex OK", map[string]any{"response": *delIdxResp})

}

func (r *indexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "indexResource.ImportState", map[string]any{"req": req, "resp": resp})

	// fetch a fresh index from pinecone
	// req.ID appears to be our only way to get the index name
	indexName := req.ID
	// note that what gets fetched from pinecone, based on purely
	// the index name, may differ from the rest of whatever is
	// specified in the resource stanza in HCL

	// Get fresh state from Pinecone
	response, err := r.client.DescribeIndex(indexName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to describe index",
			err.Error(),
		)
		return
	}

	// TODO(kevinwang) wait for the index to be in a ready state
	state := indexResourceModel{}
	state.Dimension = types.Int64Value(response.Database.Dimension)
	state.Metric = types.StringValue(response.Database.Metric)
	state.Replicas = types.Int64Value(response.Database.Replicas)
	state.Pods = types.Int64Value(response.Database.Pods)
	state.Name = types.StringValue(response.Database.Name)

	// Save data into Terraform state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
