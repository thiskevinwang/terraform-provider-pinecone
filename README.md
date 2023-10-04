# terraform-provider-pinecone

## Development

```hcl
terraform {
	required_providers {
		pinecone = {
			source ="thekevinwang.com/terraform-providers/pinecone"
		}
	}
}

provider "pinecone" {}
```

## Appendix

- https://docs.pinecone.io/
- https://github.com/hashicorp/terraform-plugin-framework
- https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider
- https://github.com/hashicorp/terraform-provider-hashicups-pf
