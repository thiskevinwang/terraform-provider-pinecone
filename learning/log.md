reference https://github.com/hashicorp/terraform-provider-hashicups-pf

create main.go
create /internal/provider/provider.go

go mod init github.com/thiskevinwang/terraform-provider-pinecone

go mod tidy

nvim ~/.terraformrc

```hcl
provider_installation {

  dev_overrides {
      "thekevinwang.com/terraform-providers/pinecone" = "/Users/kevin/go/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

Install local project to go bin

go install .

cd examples/basic/main.tf

terraform plan

✅ No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against
your configuration and found no differences, so no changes
are needed.
