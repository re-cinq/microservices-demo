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
	"html"
	"html/template"
	"sort"
	"strings"
	"unicode"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/frontend/genproto"
)

// wordMatchTier returns the best (lowest) rank tier for a query match within text.
// baseOffset is 0 for product name, 3 for product description.
// Tiers: 1=name-word-start, 2=name-mid, 3=name-word-end,
//
//	4=desc-word-start, 5=desc-mid, 6=desc-word-end.
//
// Returns 7 (no match) if the query does not appear in text at all.
func wordMatchTier(text, query string, baseOffset int) int {
	lowerText := strings.ToLower(text)
	lowerQuery := strings.ToLower(query)
	if !strings.Contains(lowerText, lowerQuery) {
		return 7
	}

	best := 7
	words := strings.FieldsFunc(lowerText, func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsPunct(r)
	})
	for _, word := range words {
		if !strings.Contains(word, lowerQuery) {
			continue
		}
		var tier int
		if strings.HasPrefix(word, lowerQuery) {
			tier = baseOffset + 1 // word-start
		} else if strings.HasSuffix(word, lowerQuery) {
			tier = baseOffset + 3 // word-end
		} else {
			tier = baseOffset + 2 // mid-word
		}
		if tier < best {
			best = tier
		}
	}
	return best
}

// rankSearchResults stable-sorts products so name-match tiers precede
// description-match tiers, and word-start matches precede mid/end matches.
func rankSearchResults(products []*pb.Product, query string) []*pb.Product {
	type ranked struct {
		product *pb.Product
		tier    int
	}
	rs := make([]ranked, len(products))
	for i, p := range products {
		nameTier := wordMatchTier(p.GetName(), query, 0)
		descTier := wordMatchTier(p.GetDescription(), query, 3)
		tier := nameTier
		if descTier < tier {
			tier = descTier
		}
		rs[i] = ranked{p, tier}
	}
	sort.SliceStable(rs, func(i, j int) bool {
		return rs[i].tier < rs[j].tier
	})
	out := make([]*pb.Product, len(rs))
	for i, r := range rs {
		out[i] = r.product
	}
	return out
}

// highlightQuery wraps every case-insensitive occurrence of query in text
// with <mark>…</mark>, returning safe HTML. Input text is escaped before
// any mark tags are inserted.
func highlightQuery(text, query string) template.HTML {
	if query == "" {
		return template.HTML(html.EscapeString(text))
	}
	escaped := html.EscapeString(text)
	lowerEscaped := strings.ToLower(escaped)
	lowerQuery := strings.ToLower(html.EscapeString(query))
	if lowerQuery == "" || !strings.Contains(lowerEscaped, lowerQuery) {
		return template.HTML(escaped)
	}

	var b strings.Builder
	remaining := escaped
	lowerRemaining := lowerEscaped
	for {
		idx := strings.Index(lowerRemaining, lowerQuery)
		if idx < 0 {
			b.WriteString(remaining)
			break
		}
		b.WriteString(remaining[:idx])
		b.WriteString("<mark>")
		b.WriteString(remaining[idx : idx+len(lowerQuery)])
		b.WriteString("</mark>")
		remaining = remaining[idx+len(lowerQuery):]
		lowerRemaining = lowerRemaining[idx+len(lowerQuery):]
	}
	return template.HTML(b.String())
}
