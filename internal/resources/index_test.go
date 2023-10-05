package resources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/thiskevinwang/terraform-provider-pinecone/internal/provider"

	"github.com/joho/godotenv"
)

var (
	providerConfig = `provider "pinecone" {
		# will use PINECONE_API_KEY
		# and PINECONE_ENVIRONMENT env vars
	}`
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"pinecone": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
)

// setup
func TestMain(m *testing.M) {
	fmt.Println(`
#######################
# 🧪 Setting up tests #
#######################
	`)

	// Load environment variables from a .env file
	if err := godotenv.Load("../../.env"); err != nil {
		panic(fmt.Sprintf("Error loading ../../.env file: %v", err))
	}

	// Sanity checks
	// fmt.Println(fmt.Sprintf("PINECONE_API_KEY: %s", os.Getenv("PINECONE_API_KEY")))
	// fmt.Println(fmt.Sprintf("PINECONE_ENVIRONMENT: %s", os.Getenv("PINECONE_ENVIRONMENT")))

	// Run the tests
	exitCode := m.Run()

	// Exit with the appropriate exit code
	os.Exit(exitCode)
}

// Note: this test requires a Pinecone account with a valid API key
// and will create and destroy REAL infrastructure.
func TestAccOrderResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `

resource "pinecone_index" "test" {
	name      = "acceptance-test"
	dimension = 1536
	metric    = "cosine"
	pods      = 1
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("pinecone_index.test", "dimension", "1536"),
					resource.TestCheckResourceAttr("pinecone_index.test", "metric", "cosine"),
					resource.TestCheckResourceAttr("pinecone_index.test", "name", "acceptance-test"),
					resource.TestCheckResourceAttr("pinecone_index.test", "pods", "1"),
					resource.TestCheckResourceAttr("pinecone_index.test", "replicas", "1"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("pinecone_index.test", "id"),
				),
			},
			// ImportState testing
			// TODO(kevinwang)
			// {
			// 	ResourceName:      "pinecone_index.test",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	// The last_updated attribute does not exist in the HashiCups
			// 	// API, therefore there is no value for it during import.
			// 	ImportStateVerifyIgnore: []string{"last_updated"},
			// },
			// Update and Read testing
			// TODO(kevinwang) - update doesn't work yet on the Pinecone free tier
			// 			{
			// 				Config: providerConfig + `

			// resource "pinecone_index" "test" {
			// 	name      = "acceptance-test"
			// 	dimension = 1536
			// 	metric    = "cosine"
			// 	pods      = 1
			// }
			// `,
			// 				Check: resource.ComposeAggregateTestCheckFunc(
			// 					// Verify attributes
			// 					resource.TestCheckResourceAttr("pinecone_index.test", "dimension", "1536"),
			// 					resource.TestCheckResourceAttr("pinecone_index.test", "metric", "cosine"),
			// 					resource.TestCheckResourceAttr("pinecone_index.test", "name", "acceptance-test"),
			// 					resource.TestCheckResourceAttr("pinecone_index.test", "pods", "1"),
			// 					resource.TestCheckResourceAttr("pinecone_index.test", "replicas", "1"),
			// 				),
			// 			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
