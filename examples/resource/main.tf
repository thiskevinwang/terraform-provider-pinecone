# terraform plan -var-file=../.tfvars
terraform {
  required_providers {
    pinecone = {
      source = "thekevinwang.com/terraform-providers/pinecone"
    }
  }
}

variable "pinecone_api_key" {
  type      = string
  sensitive = true
}

variable "pinecone_environment" {
  type = string
}

provider "pinecone" {
  apikey      = var.pinecone_api_key
  environment = var.pinecone_environment
}

resource "pinecone_index" "my-first-index" {
  name = "test"
}
