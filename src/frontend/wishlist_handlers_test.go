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
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

// withTestContext attaches the logger and session-ID values that the handlers
// expect to find on the request context (normally set by middleware).
func withTestContext(r *http.Request, session string) *http.Request {
	ctx := context.WithValue(r.Context(), ctxKeyLog{}, logrus.New().WithField("test", true))
	ctx = context.WithValue(ctx, ctxKeySessionID{}, session)
	return r.WithContext(ctx)
}

// Saving with no product_id must not error or touch the catalogue; it redirects
// back to the store and saves nothing. This is the portion of the save handler
// exercisable without the gRPC backends; the full save/view paths depend on the
// product-catalogue, currency, and cart services and are validated via
// quickstart.md.
func TestSaveProductHandler_EmptyProductID(t *testing.T) {
	fe := &frontendServer{wishlistStore: newWishlistStore()}
	req := httptest.NewRequest(http.MethodPost, "/wishlist", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = withTestContext(req, "sess-1")
	rr := httptest.NewRecorder()

	fe.saveProductHandler(rr, req)

	if rr.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusFound)
	}
	if got := rr.Header().Get("location"); got != baseUrl+"/" {
		t.Errorf("redirect location = %q, want %q", got, baseUrl+"/")
	}
	if got := fe.wishlistStore.List("sess-1"); len(got) != 0 {
		t.Errorf("nothing should be saved for empty product_id, got %v", got)
	}
}
