// Copyright 2023 The casbin Authors. All Rights Reserved.
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

package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego/context"
)

func TestBlockDebugEndpoints(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		shouldBlock    bool
	}{
		{
			name:           "Block /debug/vars",
			path:           "/debug/vars",
			expectedStatus: http.StatusNotFound,
			shouldBlock:    true,
		},
		{
			name:           "Block /debug/pprof",
			path:           "/debug/pprof",
			expectedStatus: http.StatusNotFound,
			shouldBlock:    true,
		},
		{
			name:        "Allow /api/signin",
			path:        "/api/signin",
			shouldBlock: false,
		},
		{
			name:        "Allow root path",
			path:        "/",
			shouldBlock: false,
		},
		{
			name:        "Allow /debugging (not a debug endpoint)",
			path:        "/debugging",
			shouldBlock: false,
		},
		{
			name:        "Allow /debug (without trailing slash)",
			path:        "/debug",
			shouldBlock: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock request
			r, _ := http.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			// Create a Beego context
			ctx := context.NewContext()
			ctx.Reset(w, r)

			// Call the filter
			BlockDebugEndpoints(ctx)

			// Check if the path was blocked
			if tt.shouldBlock {
				if w.Code != tt.expectedStatus {
					t.Errorf("Expected status %d for path %s, got %d", tt.expectedStatus, tt.path, w.Code)
				}
				body := w.Body.String()
				if body != "404 page not found" {
					t.Errorf("Expected '404 page not found' message, got %s", body)
				}
			} else {
				// For non-blocked paths, the filter should not write anything to the body
				if w.Body.Len() > 0 {
					t.Errorf("Filter should not write anything for path %s, but got %s", tt.path, w.Body.String())
				}
			}
		})
	}
}
