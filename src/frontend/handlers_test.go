// Copyright 2024 Google LLC
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
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// T009: table-driven tests for parseRecentlyViewed
func TestParseRecentlyViewed(t *testing.T) {
	tests := []struct {
		name   string
		cookie string // empty string = no cookie set
		want   []string
	}{
		{name: "no cookie", cookie: "", want: nil},
		{name: "empty value", cookie: " ", want: nil},
		{name: "single id", cookie: "OLJCESPC7Z", want: []string{"OLJCESPC7Z"}},
		{name: "two ids", cookie: "OLJCESPC7Z|66VCHSJNUP", want: []string{"OLJCESPC7Z", "66VCHSJNUP"}},
		{name: "five ids", cookie: "A|B|C|D|E", want: []string{"A", "B", "C", "D", "E"}},
		{name: "trailing pipe ignored", cookie: "A|B|", want: []string{"A", "B"}},
		{name: "leading pipe ignored", cookie: "|A|B", want: []string{"A", "B"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.cookie != "" {
				r.AddCookie(&http.Cookie{Name: cookieRecentlyViewed, Value: tt.cookie})
			}
			got := parseRecentlyViewed(r)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseRecentlyViewed() = %v, want %v", got, tt.want)
			}
		})
	}
}

// T010: table-driven tests for prependDedup (cookie dedup/prepend/truncate logic)
func TestPrependDedup(t *testing.T) {
	tests := []struct {
		name string
		ids  []string
		id   string
		want []string
	}{
		{
			name: "empty list, first view",
			ids:  nil,
			id:   "A",
			want: []string{"A"},
		},
		{
			name: "new product prepended",
			ids:  []string{"A", "B"},
			id:   "C",
			want: []string{"C", "A", "B"},
		},
		{
			name: "duplicate moves to front",
			ids:  []string{"A", "B", "C"},
			id:   "B",
			want: []string{"B", "A", "C"},
		},
		{
			name: "first item re-viewed stays at front without duplicate",
			ids:  []string{"A", "B", "C"},
			id:   "A",
			want: []string{"A", "B", "C"},
		},
		{
			name: "sixth product evicts oldest",
			ids:  []string{"A", "B", "C", "D", "E"},
			id:   "F",
			want: []string{"F", "A", "B", "C", "D"},
		},
		{
			name: "mid-list duplicate moves to front, oldest evicted after cap",
			ids:  []string{"A", "B", "C", "D", "E"},
			id:   "C",
			want: []string{"C", "A", "B", "D", "E"},
		},
		{
			name: "last item viewed again moves to front",
			ids:  []string{"A", "B", "C", "D", "E"},
			id:   "E",
			want: []string{"E", "A", "B", "C", "D"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := prependDedup(tt.ids, tt.id)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("prependDedup(%v, %q) = %v, want %v", tt.ids, tt.id, got, tt.want)
			}
		})
	}
}
