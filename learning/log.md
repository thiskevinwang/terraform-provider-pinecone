## Implement Provider

https://github.com/thiskevinwang/terraform-provider-pinecone/pull/1

~ 1 hour, done at 12-1am

---

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

---

## Implement Resource

https://github.com/thiskevinwang/terraform-provider-pinecone/pull/2

Start at 7am

---

get `.tfvars` working

create ./internal/resources/index.go

copy https://github.com/hashicorp/terraform-provider-scaffolding-framework/blob/main/internal/provider/example_resource.go

grok the structure. "symbols" view in GitHub gives a greate overview of what is implemented

- func New<ResourceName>
- type <ResourceName> struct {}
- type <ResourceName>Model struct {}
- func (r \*<ResourceName>) Metadata() {}
- func (r \*<ResourceName>) Schema() {}
- func (r \*<ResourceName>) Configure() {}
- func (r \*<ResourceName>) Create() {}
- func (r \*<ResourceName>) Read() {}
- func (r \*<ResourceName>) Update() {}
- func (r \*<ResourceName>) Delete() {}
- func (r \*<ResourceName>) ImportState() {}

Look at Pinecone API Docs

https://docs.pinecone.io/reference/create_index

201 String ("Created")
400
409
500

```
curl --request POST \
     --url https://controller.[[us-west4-gcp-free]].pinecone.io/databases \
     --header 'Api-Key: [[...]]' \
     --header 'accept: text/plain' \
     --header 'content-type: application/json' \
     --data '
{
  "metric": "cosine",
  "pods": 1,
  "replicas": 1,
  "pod_type": "p1.x1",
  "name": "[[..]]",
  "dimension": 1536
}
'
```

Task Copilot VSCode extension: convert the selected curl cmd to golang

...Take break at 7:40 - 7:50 to poop and mak coffee...

Go back to Pinecone API docs and select Go language instead of Shell.

- code sames don't include payload unless you _touch_ the inputs
- Google "golang map to json" because I don't like the stringified json in the code sample

https://stackoverflow.com/q/24652775

Implement `CreateIndex` method and `DescribeIndex` on my Pinecone service struct

## create ./internal/services/pinecone.go

Was tricky to implement nested object Schema

Luckily landed on https://developer.hashicorp.com/terraform/plugin/framework/handling-data/attributes/map-nested

- Google: terraform provider resource schema nested object
  - Result #2 https://developer.hashicorp.com/terraform/plugin/framework/handling-data/attributes/list-nested
    - Sidebar overview https://developer.hashicorp.com/terraform/plugin/framework/handling-data/attributes
      - Object Attribute https://developer.hashicorp.com/terraform/plugin/framework/handling-data/attributes/object
        - Nested attribute https://developer.hashicorp.com/terraform/plugin/framework/handling-data/attributes#nested-attribute-types
          - Map nested https://developer.hashicorp.com/terraform/plugin/framework/handling-data/attributes/map-nested

Oops, I misundertood the `Schema` as what is coming back from an external API, but it's actually
the definition behind the HCL stanza that is written.

I was getting and error with plan

│ Error: Invalid resource type
│
│ on main.tf line 24, in resource "pinecone_index" "my-first-index":
│ 24: resource "pinecone_index" "my-first-index" {
│
│ The provider thekevinwang.com/terraform-providers/pinecone
│ does not support resource type "pinecone_index".

Fix: import the newly created resource into provider.go > Resources()

Run go install . and try again

8:47am - take a break; 1.5 hours of work — https://x.com/thekevinwang/status/1709553196203434169?s=20
