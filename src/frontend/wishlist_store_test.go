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

import (
	"reflect"
	"sync"
	"testing"
)

func TestWishlistStore_AddAndList(t *testing.T) {
	s := newWishlistStore()
	s.Add("sess-1", "OLJCESPC7Z")
	s.Add("sess-1", "66VCHSJNUP")

	// Most recent first.
	got := s.List("sess-1")
	want := []string{"66VCHSJNUP", "OLJCESPC7Z"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("List = %v, want %v", got, want)
	}
}

func TestWishlistStore_AddIsIdempotent(t *testing.T) {
	s := newWishlistStore()
	s.Add("sess-1", "OLJCESPC7Z")
	s.Add("sess-1", "OLJCESPC7Z") // FR-007: no duplicate

	if got := s.List("sess-1"); len(got) != 1 {
		t.Fatalf("expected 1 saved product after duplicate Add, got %v", got)
	}
}

func TestWishlistStore_Contains(t *testing.T) {
	s := newWishlistStore()
	s.Add("sess-1", "OLJCESPC7Z")

	if !s.Contains("sess-1", "OLJCESPC7Z") {
		t.Errorf("Contains = false for a saved product, want true")
	}
	if s.Contains("sess-1", "NOTSAVED") {
		t.Errorf("Contains = true for an unsaved product, want false")
	}
}

func TestWishlistStore_Remove(t *testing.T) {
	s := newWishlistStore()
	s.Add("sess-1", "OLJCESPC7Z")
	s.Add("sess-1", "66VCHSJNUP")

	s.Remove("sess-1", "OLJCESPC7Z")
	if s.Contains("sess-1", "OLJCESPC7Z") {
		t.Errorf("product still saved after Remove")
	}
	// The other product is untouched.
	if !s.Contains("sess-1", "66VCHSJNUP") {
		t.Errorf("Remove deleted the wrong product")
	}
	// Removing something not saved is a no-op (no panic, no change).
	s.Remove("sess-1", "NOTSAVED")
	if got := s.List("sess-1"); len(got) != 1 {
		t.Errorf("expected 1 product after no-op Remove, got %v", got)
	}
}

func TestWishlistStore_SessionsAreIsolated(t *testing.T) {
	s := newWishlistStore()
	s.Add("sess-1", "OLJCESPC7Z")

	if s.Contains("sess-2", "OLJCESPC7Z") {
		t.Errorf("session 2 sees session 1's saved product; sessions must be isolated")
	}
	if got := s.List("sess-2"); len(got) != 0 {
		t.Errorf("List for empty session = %v, want empty", got)
	}
}

func TestWishlistStore_IgnoresEmptyIDs(t *testing.T) {
	s := newWishlistStore()
	s.Add("", "OLJCESPC7Z")
	s.Add("sess-1", "")

	if got := s.List("sess-1"); len(got) != 0 {
		t.Errorf("empty IDs should be ignored, got %v", got)
	}
}

func TestWishlistStore_ListReturnsCopy(t *testing.T) {
	s := newWishlistStore()
	s.Add("sess-1", "OLJCESPC7Z")

	got := s.List("sess-1")
	got[0] = "MUTATED"
	if s.Contains("sess-1", "MUTATED") {
		t.Errorf("mutating the returned slice corrupted the store; List must return a copy")
	}
}

func TestWishlistStore_ConcurrentAccess(t *testing.T) {
	s := newWishlistStore()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Add("sess-1", "OLJCESPC7Z")
			_ = s.List("sess-1")
			_ = s.Contains("sess-1", "OLJCESPC7Z")
		}()
	}
	wg.Wait()

	if got := s.List("sess-1"); len(got) != 1 {
		t.Fatalf("after concurrent idempotent Adds, expected 1 product, got %v", got)
	}
}
