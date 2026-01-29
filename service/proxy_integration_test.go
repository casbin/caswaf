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

package service

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/casbin/caswaf/object"
)

// TestProxyWithCookieSecurity is an integration test that verifies the proxy
// correctly modifies cookies when security flags are enabled
func TestProxyWithCookieSecurity(t *testing.T) {
	// Create a test backend server that sets cookies
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a backend that sets various cookies without security flags
		w.Header().Add("Set-Cookie", "PHPSESSID=abc123; Path=/")
		w.Header().Add("Set-Cookie", "loginToken=xyz789; Path=/; Domain=example.com")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Backend response"))
	}))
	defer backendServer.Close()

	// Create a test request
	req := httptest.NewRequest("GET", "http://example.com/test", nil)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a mock site with cookie security enabled
	site := &object.Site{
		EnableCookieSecure:   true,
		EnableCookieHttpOnly: true,
		EnableCookieSameSite: true,
	}

	// Call forwardHandler with the test backend
	forwardHandler(backendServer.URL, rr, req, site)

	// Check the response
	resp := rr.Result()
	defer resp.Body.Close()

	// Verify status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Verify body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	if string(body) != "Backend response" {
		t.Errorf("Expected body 'Backend response', got '%s'", string(body))
	}

	// Verify cookies have security flags
	cookies := resp.Header["Set-Cookie"]
	if len(cookies) != 2 {
		t.Fatalf("Expected 2 cookies, got %d", len(cookies))
	}

	for i, cookie := range cookies {
		// Each cookie should have Secure, HttpOnly, and SameSite flags
		if !strings.Contains(cookie, "Secure") {
			t.Errorf("Cookie %d missing Secure flag: %s", i, cookie)
		}
		if !strings.Contains(cookie, "HttpOnly") {
			t.Errorf("Cookie %d missing HttpOnly flag: %s", i, cookie)
		}
		if !strings.Contains(cookie, "SameSite=Lax") {
			t.Errorf("Cookie %d missing SameSite flag: %s", i, cookie)
		}
	}

	// Verify specific cookies
	phpSessionCookie := cookies[0]
	if !strings.HasPrefix(phpSessionCookie, "PHPSESSID=abc123") {
		t.Errorf("Expected PHPSESSID cookie, got: %s", phpSessionCookie)
	}

	loginTokenCookie := cookies[1]
	if !strings.HasPrefix(loginTokenCookie, "loginToken=xyz789") {
		t.Errorf("Expected loginToken cookie, got: %s", loginTokenCookie)
	}
}

// TestProxyWithoutCookieSecurity verifies that cookies pass through unchanged
// when security flags are disabled
func TestProxyWithoutCookieSecurity(t *testing.T) {
	// Create a test backend server that sets cookies
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Set-Cookie", "sessionid=test123; Path=/")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer backendServer.Close()

	// Create a test request
	req := httptest.NewRequest("GET", "http://example.com/test", nil)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a mock site with cookie security disabled
	site := &object.Site{
		EnableCookieSecure:   false,
		EnableCookieHttpOnly: false,
		EnableCookieSameSite: false,
	}

	// Call forwardHandler with the test backend
	forwardHandler(backendServer.URL, rr, req, site)

	// Check the response
	resp := rr.Result()
	defer resp.Body.Close()

	// Verify cookies don't have security flags added
	cookies := resp.Header["Set-Cookie"]
	if len(cookies) != 1 {
		t.Fatalf("Expected 1 cookie, got %d", len(cookies))
	}

	cookie := cookies[0]
	// Cookie should be unchanged from backend
	if cookie != "sessionid=test123; Path=/" {
		t.Errorf("Expected cookie unchanged, got: %s", cookie)
	}
}
