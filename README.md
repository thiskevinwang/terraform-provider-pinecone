# terraform-provider-pinecone

A Terraform provider for managing resources on [Pinecone](https://pinecone.io/). For full provider documentation, refer to the [registry’s provider page](https://registry.terraform.io/providers/thiskevinwang/pinecone/latest/docs)

## Quickstart

```hcl
terraform {
  required_providers {
    pinecone = {
      source  = "thiskevinwang/pinecone"
      version = "0.1.2"
    }
  }
}

provider "pinecone" {
  apikey      = var.pinecone_api_key # optional; PINECONE_API_KEY
  environment = var.pinecone_environment # optional; PINECONE_ENVIRONMENT
}

data "pinecone_collection" "existing-collection" {
  name = "my-existing-collection"
}

resource "pinecone_index" "my-first-index" {
  name   = "index-to-be-managed"
  metric = "cosine"
  pods   = 1

  source_collection = data.pinecone_collection.existing-collection.name
  dimension         = data.pinecone_collection.existing-collection.dimension
}
```

## Development

Check out the [examples](./examples) directory for various examples that can be run locally.

During local development, make sure to set a `provider_override` in your Terraform configuration
file (`~/.terraformrc`) to point to the local binary.

```hcl
provider_installation {
  dev_overrides {
    # Adjust this path to point to the directory where the local provider
    # binary is present. It is likely your $GOPATH/bin directory.
    "thekevinwang.com/terraform-providers/pinecone" = "/Users/kevin/go/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

```hcl
terraform {
  required_providers {
    pinecone = {
      source = "thekevinwang.com/terraform-providers/pinecone"
    }
  }
}
```

## Testing

```console
cp .env.example .env
go test -v ./...
```

## Documenting

This project relies on `./tools/tools.go` to install [`tfplugindocs`](github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs).

```console
export GOBIN=$PWD/bin
export PATH=$GOBIN:$PATH
go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
which tfplugindocs
```

Run `tfplugindocs` to generate docs, and preview individual files at https://registry.terraform.io/tools/doc-preview

## Releasing

There are a few one time steps that have been done already and will not be covered in this README. See the following footnote for more information. [^release]

[^release]: Docs on publishing: https://developer.hashicorp.com/terraform/registry/providers/publishing#publishing-to-the-registry

To release a new version of the provider to the registry, a new GitHub release needs to be created.
Use the following steps to proceed.

1. Ideally you’re on `main`, and have a clean working tree.
2. Ensure the commit (aka HEAD) you're about to release is tagged.
   - `git tag v0.1.2`
   - `git push origin v0.1.2`
3. Run `goreleaser`: `GITHUB_TOKEN=$(gh auth token) goreleaser release --clean`
   - This will create a new GitHub release which should be detected by the Terraform registry shortly after.

> [!NOTE]
>
> The goreleaser config requires `GPG_FINGERPRINT` to be set. Maybe try the following (needs to be tested)
>
> ```bash
> export GPG_FINGERPRINT=$(gpg --list-keys --with-colons | awk -F: '$1 == "fpr" || $1 == "fp2" {print $10}' | head -n 1)
> export GITHUB_TOKEN=$(gh auth token)
> goreleaser release --clean
> ```
>
> - Ref: https://unix.stackexchange.com/a/743986

## Appendix

- https://thekevinwang.com/2023/10/05/build-publish-terraform-provider
- https://docs.pinecone.io/
- https://github.com/hashicorp/terraform-plugin-framework
- https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider
- https://github.com/hashicorp/terraform-provider-hashicups-pf
