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

import (
	"bytes"
	"html/template"
	"os"
	"strings"
	"testing"
)

// TestHomeRendersSearchMarkup locks in the DOM contract that
// static/js/search.js depends on (see specs/001-product-search/data-model.md).
// It does NOT exercise the full homeHandler — that would require stubbing all
// upstream gRPC clients. Instead it parses templates/home.html in isolation,
// substituting tiny harness templates for the header and footer it includes,
// and asserts the rendered body contains the expected attributes, IDs, and
// per-card data-product-name values.
func TestHomeRendersSearchMarkup(t *testing.T) {
	homeBytes, err := os.ReadFile("templates/home.html")
	if err != nil {
		t.Fatalf("read templates/home.html: %v", err)
	}

	tpl := template.New("").Funcs(template.FuncMap{
		"renderMoney":        func(_ interface{}) string { return "$0.00" },
		"renderCurrencyLogo": func(_ interface{}) string { return "" },
	})

	// Stub the two templates that home.html includes. We only care about
	// home.html's own markup, not the surrounding chrome.
	if _, err := tpl.New("header").Parse(`<!doctype html><html><body>`); err != nil {
		t.Fatalf("parse header stub: %v", err)
	}
	if _, err := tpl.New("footer").Parse(`</body></html>`); err != nil {
		t.Fatalf("parse footer stub: %v", err)
	}
	if _, err := tpl.Parse(string(homeBytes)); err != nil {
		t.Fatalf("parse home.html: %v", err)
	}

	type item struct {
		Id, Name, Picture string
	}
	type productView struct {
		Item  item
		Price interface{}
	}

	data := map[string]interface{}{
		"baseUrl":       "",
		"platform_css":  "",
		"platform_name": "",
		"products": []productView{
			{Item: item{Id: "WATCH1", Name: "Watch", Picture: "/w.jpg"}},
			{Item: item{Id: "SP1", Name: "Salt & Pepper Shakers", Picture: "/sp.jpg"}},
		},
	}

	var buf bytes.Buffer
	if err := tpl.ExecuteTemplate(&buf, "home", data); err != nil {
		t.Fatalf("execute home template: %v", err)
	}
	body := buf.String()

	wantSubstrings := []string{
		`id="search-input"`,
		`id="search-clear"`,
		`id="search-empty-state"`,
		`role="status"`,
		`aria-live="polite"`,
		`No products match your search.`,
		`data-product-name="Watch"`,
		// html/template escapes `&` in attribute context.
		`data-product-name="Salt &amp; Pepper Shakers"`,
		`/static/js/search.js`,
	}
	for _, want := range wantSubstrings {
		if !strings.Contains(body, want) {
			t.Errorf("rendered home body is missing %q", want)
		}
	}
}
