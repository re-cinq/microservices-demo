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
	"testing"
)

func TestRecentlyViewed_MostRecentFirst(t *testing.T) {
	s := newRecentlyViewedStore()
	s.Record("sess", "A")
	s.Record("sess", "B")
	s.Record("sess", "C")

	got := s.List("sess", "", 4)
	want := []string{"C", "B", "A"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("List() = %v, want %v", got, want)
	}
}

func TestRecentlyViewed_CapAtMax(t *testing.T) {
	s := newRecentlyViewedStore()
	for _, id := range []string{"A", "B", "C", "D", "E", "F"} {
		s.Record("sess", id)
	}

	got := s.List("sess", "", 4)
	want := []string{"F", "E", "D", "C"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("List() = %v, want %v (must cap at 4)", got, want)
	}
}

func TestRecentlyViewed_ExcludesCurrentProduct(t *testing.T) {
	s := newRecentlyViewedStore()
	s.Record("sess", "A")
	s.Record("sess", "B")
	s.Record("sess", "C")

	got := s.List("sess", "C", 4)
	want := []string{"B", "A"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("List(exclude=C) = %v, want %v", got, want)
	}
}

func TestRecentlyViewed_ExcludesAllOccurrencesOfCurrent(t *testing.T) {
	// With no de-duplication (AIP-201), the current product may appear several
	// times; excluding it must remove every occurrence.
	s := newRecentlyViewedStore()
	s.Record("sess", "A")
	s.Record("sess", "B")
	s.Record("sess", "A")

	got := s.List("sess", "A", 4)
	want := []string{"B"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("List(exclude=A) = %v, want %v", got, want)
	}
}

func TestRecentlyViewed_DuplicatesPreserved(t *testing.T) {
	// This story does NOT de-duplicate (that is AIP-201). Re-viewing a product
	// records another entry.
	s := newRecentlyViewedStore()
	s.Record("sess", "A")
	s.Record("sess", "B")
	s.Record("sess", "A")

	got := s.List("sess", "", 4)
	want := []string{"A", "B", "A"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("List() = %v, want %v (duplicates expected without dedup)", got, want)
	}
}

func TestRecentlyViewed_PerSessionIsolation(t *testing.T) {
	s := newRecentlyViewedStore()
	s.Record("s1", "A")
	s.Record("s2", "B")

	if got := s.List("s1", "", 4); !reflect.DeepEqual(got, []string{"A"}) {
		t.Fatalf("s1 List() = %v, want [A]", got)
	}
	if got := s.List("s2", "", 4); !reflect.DeepEqual(got, []string{"B"}) {
		t.Fatalf("s2 List() = %v, want [B]", got)
	}
}

func TestRecentlyViewed_UnknownSessionEmpty(t *testing.T) {
	s := newRecentlyViewedStore()
	if got := s.List("nobody", "", 4); len(got) != 0 {
		t.Fatalf("List(unknown) = %v, want empty", got)
	}
}

func TestRecentlyViewed_IgnoresEmptyInputs(t *testing.T) {
	s := newRecentlyViewedStore()
	s.Record("", "A")
	s.Record("sess", "")
	if got := s.List("sess", "", 4); len(got) != 0 {
		t.Fatalf("List() = %v, want empty after no-op records", got)
	}
}
