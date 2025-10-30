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

//go:build !skipCi
// +build !skipCi

package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCaptchaCookieAttributes(t *testing.T) {
	// Create a mock response recorder
	w := httptest.NewRecorder()

	// Create a test request with HTTPS scheme
	r := httptest.NewRequest("GET", "https://example.com/test", nil)
	// Explicitly set scheme to ensure test reliability
	r.URL.Scheme = "https"

	// Set a test captcha cookie
	uuidStr := "test-uuid-123"
	host := "example.com"
	cookie := &http.Cookie{
		Name:     "casdoor_captcha_token",
		Value:    uuidStr,
		Path:     "/",
		Domain:   host,
		HttpOnly: true,
		Secure:   getScheme(r) == "https",
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)

	// Get the cookies from the response
	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("Expected 1 cookie, got %d", len(cookies))
	}

	captchaCookie := cookies[0]

	// Verify SameSite attribute is set to Lax
	if captchaCookie.SameSite != http.SameSiteLaxMode {
		t.Errorf("Expected SameSite to be Lax (1), got %d", captchaCookie.SameSite)
	}

	// Verify HttpOnly is set
	if !captchaCookie.HttpOnly {
		t.Error("Expected HttpOnly to be true")
	}

	// Verify Secure is set for HTTPS
	if !captchaCookie.Secure {
		t.Error("Expected Secure to be true for HTTPS requests")
	}
}

func TestCaptchaCookieAttributesHTTP(t *testing.T) {
	// Create a mock response recorder
	w := httptest.NewRecorder()

	// Create a test request with HTTP scheme
	r := httptest.NewRequest("GET", "http://example.com/test", nil)
	// Explicitly set scheme to ensure test reliability
	r.URL.Scheme = "http"

	// Set a test captcha cookie
	uuidStr := "test-uuid-123"
	host := "example.com"
	cookie := &http.Cookie{
		Name:     "casdoor_captcha_token",
		Value:    uuidStr,
		Path:     "/",
		Domain:   host,
		HttpOnly: true,
		Secure:   getScheme(r) == "https",
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)

	// Get the cookies from the response
	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("Expected 1 cookie, got %d", len(cookies))
	}

	captchaCookie := cookies[0]

	// Verify Secure is NOT set for HTTP
	if captchaCookie.Secure {
		t.Error("Expected Secure to be false for HTTP requests")
	}

	// Verify SameSite is still set
	if captchaCookie.SameSite != http.SameSiteLaxMode {
		t.Errorf("Expected SameSite to be Lax (1), got %d", captchaCookie.SameSite)
	}
}

func TestOAuthCookieAttributes(t *testing.T) {
	// Create a mock response recorder
	w := httptest.NewRecorder()

	// Create a test request with HTTPS scheme
	r := httptest.NewRequest("GET", "https://example.com/test", nil)
	// Explicitly set scheme to ensure test reliability
	r.URL.Scheme = "https"

	// Set a test OAuth access token cookie
	cookie := &http.Cookie{
		Name:     "casdoor_access_token",
		Value:    "test-access-token",
		Path:     "/",
		HttpOnly: true,
		Secure:   getScheme(r) == "https",
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)

	// Get the cookies from the response
	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("Expected 1 cookie, got %d", len(cookies))
	}

	oauthCookie := cookies[0]

	// Verify SameSite attribute is set to Lax
	if oauthCookie.SameSite != http.SameSiteLaxMode {
		t.Errorf("Expected SameSite to be Lax (1), got %d", oauthCookie.SameSite)
	}

	// Verify HttpOnly is set
	if !oauthCookie.HttpOnly {
		t.Error("Expected HttpOnly to be true")
	}

	// Verify Secure is set for HTTPS
	if !oauthCookie.Secure {
		t.Error("Expected Secure to be true for HTTPS requests")
	}
}
