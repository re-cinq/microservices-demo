// Copyright 2026 Google LLC
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

import "sync"

// ratingAggregate is the running tally of submitted star ratings for one
// product.
type ratingAggregate struct {
	sum   int64
	count int64
}

// productRatingStore holds shopper-submitted 1-5 star ratings in memory in the
// frontend. Ratings are not persisted (no datastore), so they reset whenever
// the frontend restarts.
type productRatingStore struct {
	mu  sync.Mutex
	agg map[string]*ratingAggregate
}

func newProductRatingStore() *productRatingStore {
	return &productRatingStore{agg: make(map[string]*ratingAggregate)}
}

// record folds a single 1-5 star submission into the product's aggregate.
func (s *productRatingStore) record(productID string, stars int32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a := s.agg[productID]
	if a == nil {
		a = &ratingAggregate{}
		s.agg[productID] = a
	}
	a.sum += int64(stars)
	a.count++
}

// get returns the average rating and the number of ratings for a product.
// A count of 0 means the product has no ratings yet.
func (s *productRatingStore) get(productID string) (avg float32, count int32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a := s.agg[productID]
	if a == nil || a.count == 0 {
		return 0, 0
	}
	return float32(a.sum) / float32(a.count), int32(a.count)
}

// productRatings is the process-wide in-memory rating store for the frontend.
var productRatings = newProductRatingStore()
