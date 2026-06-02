# Technical Constraints (do not deviate)

- Use only the services already in this repo. Do NOT add new services.
- Do NOT introduce Elasticsearch, Solr, vector databases, or any new datastore.
- Use the existing in-memory product catalogue loaded from `productcatalogservice/products.json`. Filter in memory.
- Match the language of the service you're editing: frontend is Go, productcatalogservice is Go. Use the existing protobuf/gRPC patterns if you extend productcatalogservice.
- The CI deploys whatever lands on `attendee/<your-name>`. Do NOT add new infrastructure config, Helm charts, manifests, or environment variables.
- Stay inside this branch and this repo. Do NOT touch the build pipeline.
