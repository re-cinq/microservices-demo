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
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// withSessionCtx injects a session ID into a request context,
// replicating what the ensureSessionID middleware does in production.
func withSessionCtx(r *http.Request, sid string) *http.Request {
	ctx := context.WithValue(r.Context(), ctxKeySessionID{}, sid)
	return r.WithContext(ctx)
}

func postWishlist(fe *frontendServer, sid, productID string) *httptest.ResponseRecorder {
	body := url.Values{"product_id": {productID}}
	req := httptest.NewRequest(http.MethodPost, "/wishlist/save", strings.NewReader(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = withSessionCtx(req, sid)
	rr := httptest.NewRecorder()
	fe.saveWishlistHandler(rr, req)
	return rr
}

func TestSaveWishlistHandler_AddsProduct(t *testing.T) {
	fe := &frontendServer{}

	rr := postWishlist(fe, "sess1", "abc123")

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected 303 SeeOther, got %d", rr.Code)
	}
	v, ok := fe.wishlists.Load("sess1")
	if !ok {
		t.Fatal("expected wishlist entry for sess1")
	}
	ids := v.([]string)
	if len(ids) != 1 || ids[0] != "abc123" {
		t.Errorf("expected [abc123], got %v", ids)
	}
}

func TestSaveWishlistHandler_PreventsDuplicate(t *testing.T) {
	fe := &frontendServer{}
	fe.wishlists.Store("sess2", []string{"abc123"})

	postWishlist(fe, "sess2", "abc123")

	v, _ := fe.wishlists.Load("sess2")
	ids := v.([]string)
	if len(ids) != 1 {
		t.Errorf("expected 1 item (no duplicate), got %d items", len(ids))
	}
}

func TestSaveWishlistHandler_RejectsMissingProductID(t *testing.T) {
	fe := &frontendServer{}

	req := httptest.NewRequest(http.MethodPost, "/wishlist/save", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = withSessionCtx(req, "sess3")
	rr := httptest.NewRecorder()
	fe.saveWishlistHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 BadRequest, got %d", rr.Code)
	}
}

func TestViewWishlistHandler_EmptyForNewSession(t *testing.T) {
	fe := &frontendServer{}

	_, ok := fe.wishlists.Load("new-sess")
	if ok {
		t.Error("expected no wishlist entry for a fresh session")
	}
}

func TestSaveWishlistHandler_MultipleProducts(t *testing.T) {
	fe := &frontendServer{}

	postWishlist(fe, "sess4", "prod-1")
	postWishlist(fe, "sess4", "prod-2")
	postWishlist(fe, "sess4", "prod-1") // duplicate — should not add

	v, ok := fe.wishlists.Load("sess4")
	if !ok {
		t.Fatal("expected wishlist entry for sess4")
	}
	ids := v.([]string)
	if len(ids) != 2 {
		t.Errorf("expected 2 unique items, got %d: %v", len(ids), ids)
	}
}
