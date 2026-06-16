// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"strings"
	"sync"
	"time"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/productcatalogservice/genproto"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// ratingAggregate is the in-memory tally of submitted star ratings for a
// single product. It is never persisted (no datastore), so all ratings reset
// when the service restarts.
type ratingAggregate struct {
	sum   int64
	count int64
}

type productCatalog struct {
	pb.UnimplementedProductCatalogServiceServer
	catalog pb.ListProductsResponse

	ratingsMu sync.Mutex
	ratings   map[string]*ratingAggregate
}

func (p *productCatalog) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (p *productCatalog) Watch(req *healthpb.HealthCheckRequest, ws healthpb.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "health check via Watch not implemented")
}

func (p *productCatalog) ListProducts(context.Context, *pb.Empty) (*pb.ListProductsResponse, error) {
	time.Sleep(extraLatency)

	catalog := p.parseCatalog()
	products := make([]*pb.Product, len(catalog))
	for i, product := range catalog {
		products[i] = p.withRating(product)
	}

	return &pb.ListProductsResponse{Products: products}, nil
}

func (p *productCatalog) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	time.Sleep(extraLatency)

	catalog := p.parseCatalog()
	for _, product := range catalog {
		if req.Id == product.Id {
			return p.withRating(product), nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "no product with ID %s", req.Id)
}

// RateProduct records a single 1-5 star submission for a product, folding it
// into that product's in-memory aggregate. Ratings are not persisted.
func (p *productCatalog) RateProduct(ctx context.Context, req *pb.RateProductRequest) (*pb.Empty, error) {
	time.Sleep(extraLatency)

	if req.Stars < 1 || req.Stars > 5 {
		return nil, status.Errorf(codes.InvalidArgument, "stars must be between 1 and 5, got %d", req.Stars)
	}

	// Confirm the product exists before recording a rating for it.
	found := false
	for _, product := range p.parseCatalog() {
		if product.Id == req.ProductId {
			found = true
			break
		}
	}
	if !found {
		return nil, status.Errorf(codes.NotFound, "no product with ID %s", req.ProductId)
	}

	p.ratingsMu.Lock()
	defer p.ratingsMu.Unlock()
	if p.ratings == nil {
		p.ratings = make(map[string]*ratingAggregate)
	}
	agg := p.ratings[req.ProductId]
	if agg == nil {
		agg = &ratingAggregate{}
		p.ratings[req.ProductId] = agg
	}
	agg.sum += int64(req.Stars)
	agg.count++

	return &pb.Empty{}, nil
}

// withRating returns a copy of the product annotated with its current average
// rating and rating count. A copy is used so the cached catalog objects are
// never mutated (which would race with concurrent reads).
func (p *productCatalog) withRating(product *pb.Product) *pb.Product {
	out := proto.Clone(product).(*pb.Product)

	p.ratingsMu.Lock()
	agg := p.ratings[product.Id]
	if agg != nil && agg.count > 0 {
		out.Rating = float32(agg.sum) / float32(agg.count)
		out.NumRatings = int32(agg.count)
	}
	p.ratingsMu.Unlock()

	return out
}

func (p *productCatalog) SearchProducts(ctx context.Context, req *pb.SearchProductsRequest) (*pb.SearchProductsResponse, error) {
	time.Sleep(extraLatency)

	var ps []*pb.Product
	for _, product := range p.parseCatalog() {
		if strings.Contains(strings.ToLower(product.Name), strings.ToLower(req.Query)) ||
			strings.Contains(strings.ToLower(product.Description), strings.ToLower(req.Query)) {
			ps = append(ps, product)
		}
	}

	return &pb.SearchProductsResponse{Results: ps}, nil
}

func (p *productCatalog) parseCatalog() []*pb.Product {
	if reloadCatalog || len(p.catalog.Products) == 0 {
		err := loadCatalog(&p.catalog)
		if err != nil {
			return []*pb.Product{}
		}
	}

	return p.catalog.Products
}
