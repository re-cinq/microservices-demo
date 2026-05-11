# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

**Online Boutique** is a polyglot microservices e-commerce demo app. Eleven services communicate exclusively over **gRPC** (protobuf definitions in `protos/demo.proto`). The frontend is the only HTTP-facing service; everything else is internal gRPC only.

## Build & Deploy

The primary development tool is **skaffold**, which builds all container images and deploys them to a Kubernetes cluster.

```sh
# Build and deploy to the currently active kubectl context (local cluster)
skaffold run

# Build, deploy, and watch for changes (hot-reload loop)
skaffold dev

# Deploy using Google Cloud Build instead of local Docker
skaffold run -p gcb --default-repo=us-docker.pkg.dev/PROJECT_ID/microservices-demo

# Tear down everything deployed by skaffold
skaffold delete
```

For local cluster options: Minikube (`minikube start --cpus=4 --memory 4096 --disk-size 32g`), Kind (`kind create cluster`), or Docker Desktop with Kubernetes enabled.

After deploying locally, forward the frontend:
```sh
kubectl port-forward deployment/frontend 8080:8080
```

## Running Tests

Each service has its own test command run from within the service directory:

| Language | Service(s) | Command |
|---|---|---|
| Go | `frontend`, `checkoutservice`, `productcatalogservice`, `shippingservice` | `go test ./...` |
| C# | `cartservice` | `dotnet test` (run from `src/cartservice/`) |
| Java | `adservice` | `./gradlew test` (run from `src/adservice/`) |
| Python | `emailservice`, `recommendationservice`, `loadgenerator` | `pytest` (if tests exist) |
| Node.js | `currencyservice`, `paymentservice` | `npm test` (run from service directory) |

To run a single Go test:
```sh
cd src/shippingservice && go test -run TestGetQuote ./...
```

## Architecture

```
Browser → frontend (Go, HTTP :8080)
             ↓ gRPC calls to all other services
    ┌────────────────────────────────────┐
    │  productcatalogservice (Go)        │  reads products.json
    │  currencyservice (Node.js)         │  highest QPS, currency conversion
    │  cartservice (C#)                  │  backed by Redis (default), Spanner, or AlloyDB
    │  recommendationservice (Python)    │
    │  checkoutservice (Go)              │  orchestrates: cart→payment→shipping→email
    │  paymentservice (Node.js)          │  mock charge
    │  shippingservice (Go)              │  mock shipping
    │  emailservice (Python)             │  mock email
    │  adservice (Java)                  │  context-based ads
    │  shoppingassistantservice (Python) │  Gemini-powered AI assistant via AlloyDB
    └────────────────────────────────────┘
```

**checkoutservice** is the critical orchestrator: it calls cartservice, productcatalogservice, currencyservice, shippingservice, paymentservice, and emailservice in sequence to complete an order.

**cartservice** supports three pluggable backends selected at runtime: `RedisCartStore` (default), `SpannerCartStore`, and `AlloyDBCartStore` — see `src/cartservice/src/cartstore/`.

**shoppingassistantservice** is an optional Gemini + LangChain service (Python/Flask) that uses AlloyDB with pgvector for product embeddings.

## Proto & Generated Code

The single source of truth for all gRPC interfaces is `protos/demo.proto`. Each service has its own pre-generated protobuf code (do not edit generated files directly). To regenerate, run the `genproto.sh` script inside the relevant service directory.

## Kubernetes & Kustomize

Base manifests live in `kubernetes-manifests/`. Optional features are toggled via kustomize components in `kustomize/components/`:

- `memorystore` — swap Redis for Google Cloud Memorystore
- `spanner` — swap cartservice backend to Spanner
- `alloydb` — swap cartservice backend to AlloyDB
- `service-mesh-istio` — add Istio sidecar configuration
- `cymbal-branding` — apply Cymbal brand theming
- `shopping-assistant` — deploy shoppingassistantservice
- `network-policies` — add Kubernetes NetworkPolicy resources
- `without-loadgenerator` — exclude the load generator

Enable a component by uncommenting it in `kustomize/kustomization.yaml`.

A Helm chart is also available in `helm-chart/` as an alternative to kustomize.

## Adding a New Microservice

1. Create `src/<servicename>/` with source code and a `Dockerfile`
2. Add the image to `skaffold.yaml` under `build.artifacts`
3. Create a kustomize component in `kustomize/components/<servicename>/` with Deployment + Service YAMLs
4. Reference the component in `kustomize/kustomization.yaml`
5. Add values to `helm-chart/values.yaml` and a template in `helm-chart/templates/`

See `docs/adding-new-microservice.md` for the full guide.

## Service Environment Variables

Services discover each other via environment variables set in the Kubernetes manifests. The frontend requires these to be set (panics on missing values):
- `PRODUCT_CATALOG_SERVICE_ADDR`, `CURRENCY_SERVICE_ADDR`, `CART_SERVICE_ADDR`
- `RECOMMENDATION_SERVICE_ADDR`, `CHECKOUT_SERVICE_ADDR`, `SHIPPING_SERVICE_ADDR`
- `AD_SERVICE_ADDR`, `SHOPPING_ASSISTANT_SERVICE_ADDR`

Optional: `ENABLE_TRACING=1` (OpenTelemetry/OTLP), `ENABLE_PROFILER=1` (Cloud Profiler), `COLLECTOR_SERVICE_ADDR`.

<!-- SPECKIT START -->
For additional context about technologies to be used, project structure,
shell commands, and other important information, read the current plan:
specs/001-browser-product-search/plan.md
<!-- SPECKIT END -->
