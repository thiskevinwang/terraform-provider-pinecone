package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
}

type PineconeDescribeIndexResponse struct {
	Database struct {
		Name      string `json:"name"`
		Metric    string `json:"metric"`
		Dimension int    `json:"dimension"`
		Replicas  int    `json:"replicas"`
		Shards    int    `json:"shards"`
		Pods      int    `json:"pods"`
	} `json:"database"`
	Status struct {
		Waiting []interface{} `json:"waiting"`
		Crashed []interface{} `json:"crashed"`
		Host    string        `json:"host"`
		Port    int           `json:"port"`
		State   string        `json:"state"`
		Ready   bool          `json:"ready"`
	} `json:"status"`
}

// indexResourceModel maps the resource schema data.
type indexResourceModel struct {
	Database struct {
		Name      string `json:"name"`
		Metric    string `json:"metric"`
		Dimension int    `json:"dimension"`
		Replicas  int    `json:"replicas"`
		Shards    int    `json:"shards"`
		Pods      int    `json:"pods"`
	} `json:"database" tfsdk:"database"`
	Status struct {
		Waiting []interface{} `json:"waiting"`
		Crashed []interface{} `json:"crashed"`
		Host    string        `json:"host"`
		Port    int           `json:"port"`
		State   string        `json:"state"`
		Ready   bool          `json:"ready"`
	} `json:"status" tfsdk:"status"`
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
				Description: "The name of the index.",
				Required:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *indexResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
}

// Create a new resource.
func (r *indexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan indexResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *indexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

func (r *indexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *indexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

func (r *indexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
