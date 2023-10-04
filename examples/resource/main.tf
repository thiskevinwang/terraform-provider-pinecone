provider "pinecone" {
  apikey      = var.pinecone_api_key
  environment = var.pinecone_environment
}

resource "pinecone_index" "my-first-index" {
  name      = "my-little-pony"
  dimension = 1536
  metric    = "cosine"
}
