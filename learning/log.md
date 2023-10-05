## Implement Provider

https://github.com/thiskevinwang/terraform-provider-pinecone/pull/1

10-04-2023 12am (start)
10-04-2023 12-1am (stop; +1h; total elapsed 1h)

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

10-04-2023 7am (start)

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

https://x.com/thekevinwang/status/1709553196203434169?s=20

10-04-2023 7am-8:47am (stop; +1.5hr of work; total elapsed 2.5h)

10-04-2023 3:30pm (start) - early end to my last day at HashiCorp. Resume work on this provider.

Run `terraform apply`

│ Error: Value Conversion Error
│
│ with pinecone_index.my-first-index,
│ An unexpected error was encountered trying to convert tftypes.Value into
│ resources.indexResourceModel. This is always an error in the provider. Please
│ report the following to the provider developer:
│
│ mismatch between struct and object: Struct defines fields not found in object:
│ status and database. Object defines fields not found in struct: name.
╵

reference https://github.com/hashicorp/terraform-provider-scaffolding-framework/blob/main/internal/provider/example_resource.go#L35-L39

```go
type indexResourceModel struct {
	Name types.String `tfsdk:"name"`
}
```

│ Error: Missing Resource State After Create
│
│ with pinecone_index.my-first-index,
│ on main.tf line 24, in resource "pinecone_index" "my-first-index":
│ 24: resource "pinecone_index" "my-first-index" {
│
│ The Terraform Provider unexpectedly returned no resource state after having no
│ errors in the resource creation. This is always an issue in the Terraform
│ Provider and should be reported to the provider developers.
│
│ The resource may have been successfully created, but Terraform is not tracking
│ it. Applying the configuration again with no other action may result in duplicate
│ resource errors. Import the resource if the resource was actually created and
│ Terraform should be tracking it.

TF_LOG=trace terraform apply -var-file="../.tfvars"

`tflog.Error(ctx, "reading a resource")` output looks like

```
2023-10-04T15:51:32.287-0400 [ERROR] provider.terraform-provider-pinecone: reading a resource: @module=pinecone tf_resource_type=pinecone_index tf_req_id=95b0918e-f096-338a-19ca-5c232de3ca63 tf_rpc=ReadResource @caller=/Users/kevin/repos/terraform-provider-pinecone/internal/resources/index.go:95 tf_provider_addr=thekevinwang.com/terraform-providers/pinecone timestamp=2023-10-04T15:51:32.287-0400
```

Here, the https://github.com/hashicorp/terraform-provider-scaffolding-framework/blob/c7f8b736aec6b14daac8533176931af51a0df22a/internal/provider/example_resource.go#L122 repo starts to fall short in code examples for resource `Read` and `Create` methods.

... and https://github.com/hashicorp/terraform-provider-hashicups-pf/blob/c42733f24b8c4e0583d750e28c0490ad82a20972/internal/provider/order_resource.go#L205 is more helpful.

I probably should've started with implementing a data_source instead since that is Read-only which is more simple.

Fix dupe content here: https://developer.hashicorp.com/terraform/plugin/framework/migrating/attributes-blocks/default-values

- needs to show import statement from
  - "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
  - "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"

10-04-2023 3:30pm-4:34pm (stop; nap. +1hr; total elapsed 3.5h)

10-04-2023 5:43pm (start)

Play Dragonball Z on crunchyroll in the background

blew away my statefile because I was getting errors
about some value being null

Now..

│ Error: Value Conversion Error
│
│ with pinecone_index.my-first-index,
│ on main.tf line 27, in resource "pinecone_index" "my-first-index":
│ 27: metric = "cosine"
│
│ An unexpected error was encountered trying to convert tftypes.Value into
│ basetypes.SetType. This is always an error in the provider. Please report
│ the following to the provider developer:
│
│ cannot reflect tftypes.String into a struct, must be an object

APPLY WORKS

10-04-2023 5:43pm-7:23pm (stop; +1.5hr; total elapsed 5h)

10-04-2023 10pm (start)

Do more docs reading

implement more service methods. Pinecone API docs so good.

Do some whimsical diagramming ... it helps alot

10-04-2023 10pm-12:30am (stop; +2.5hr; total elapsed 7.5h)

Implemented update and destroy.

Next: acceptance tests?

10-05-2023 Flexible — I'm not going to track time anymore

Started working on tests and new-to-me golang module/import/testing isms
https://github.com/hashicorp/terraform-plugin-testing/issues/185

10-05-2023 12:52pm Got a basic CREATE-DESTROY acceptance test working

10-05-2023 12:58pm Merged: https://github.com/thiskevinwang/terraform-provider-pinecone/pull/2

Continue on to implement a `pinecone_collection` data-source
