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

// recentlyViewedStore keeps, per session, the product IDs a shopper has viewed,
// ordered most-recently-viewed first. It is in-memory only: no datastore is
// introduced (per the epic's technical constraints). It is safe for concurrent
// use across requests.
//
// De-duplication is intentionally NOT performed here; repeated views append
// repeated entries. Moving a revisited product to the front / removing
// duplicates is a separate story (AIP-201).
type recentlyViewedStore struct {
	mu        sync.Mutex
	bySession map[string][]string
}

func newRecentlyViewedStore() *recentlyViewedStore {
	return &recentlyViewedStore{bySession: make(map[string][]string)}
}

// Record prepends productID to the session's list, so the most recently viewed
// product is always at index 0.
func (s *recentlyViewedStore) Record(sessionID, productID string) {
	if sessionID == "" || productID == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bySession[sessionID] = append([]string{productID}, s.bySession[sessionID]...)
}

// List returns up to max product IDs for the session, most-recently-viewed
// first, skipping any entry equal to excludeProductID (used to keep the product
// currently being viewed out of its own strip). An unknown session yields an
// empty slice.
func (s *recentlyViewedStore) List(sessionID, excludeProductID string, limit int) []string {
	if limit <= 0 {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	ids := s.bySession[sessionID]
	out := make([]string, 0, limit)
	for _, id := range ids {
		if id == excludeProductID {
			continue
		}
		out = append(out, id)
		if len(out) == limit {
			break
		}
	}
	return out
}
