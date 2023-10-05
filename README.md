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

## Testing

```console
cp .env.example .env
go test -v ./...
```

## Release

https://developer.hashicorp.com/terraform/registry/providers/publishing#publishing-to-the-registry

- Sign in to https://registry.terraform.io/publish/provider/github/thiskevinwang/terraform-provider-pinecone
- Ensure repo has `.goreleaser.yml`
- gpg --armor --export "[EMAIL]"
  - add this to registry signing keys
- gpg --armor --detach-sign
- git tag v0.1.1
- GITHUB_TOKEN=$(gh auth token) goreleaser release --clean

## Appendix

- https://docs.pinecone.io/
- https://github.com/hashicorp/terraform-plugin-framework
- https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider
- https://github.com/hashicorp/terraform-provider-hashicups-pf
