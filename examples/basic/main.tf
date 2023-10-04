terraform {
  required_providers {
    pinecone = {
      source = "thekevinwang.com/terraform-providers/pinecone"
    }
  }
}

provider "pinecone" {
  apikey      = "1"
  environment = "dev"
}
