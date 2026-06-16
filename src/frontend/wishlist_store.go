// Copyright 2018 Google LLC
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

// wishlistStore holds each session's saved product IDs in memory. It is
// process-local by design: per the project constitution this feature introduces
// no datastore, so saved products are scoped to the process and the shopper's
// session and are not shared across frontend replicas.
type wishlistStore struct {
	mu    sync.Mutex
	items map[string][]string // sessionID -> product IDs, most recent first
}

func newWishlistStore() *wishlistStore {
	return &wishlistStore{items: make(map[string][]string)}
}

// Add saves productID for the session. It is idempotent: saving a product that is
// already saved does not create a duplicate and does not change ordering. Empty
// session or product IDs are ignored.
func (s *wishlistStore) Add(sessionID, productID string) {
	if sessionID == "" || productID == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, id := range s.items[sessionID] {
		if id == productID {
			return
		}
	}
	// Most recent first, for display order.
	s.items[sessionID] = append([]string{productID}, s.items[sessionID]...)
}

// List returns the product IDs saved for the session, most recent first. The
// returned slice is a copy and safe for the caller to retain.
func (s *wishlistStore) List(sessionID string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	ids := s.items[sessionID]
	out := make([]string, len(ids))
	copy(out, ids)
	return out
}

// Contains reports whether productID is currently saved for the session.
func (s *wishlistStore) Contains(sessionID, productID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, id := range s.items[sessionID] {
		if id == productID {
			return true
		}
	}
	return false
}
