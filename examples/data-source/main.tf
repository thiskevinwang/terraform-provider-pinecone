provider "pinecone" {
  apikey      = var.pinecone_api_key
  environment = var.pinecone_environment
}

data "pinecone_collection" "existing-collection" {
  name = "testindex"
}

resource "pinecone_index" "my-first-index" {
  name   = "testidx"
  metric = "cosine"
  pods   = 1

  source_collection = data.pinecone_collection.existing-collection.name
  # index and collection dimension must match
  dimension = data.pinecone_collection.existing-collection.dimension
}
