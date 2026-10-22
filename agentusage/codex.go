// Copyright 2026 The casbin Authors. All Rights Reserved.
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

// Package agentusage asks the vendor what the account behind a stored sign-in
// has left of its limits, using that account's own token.
package agentusage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/apache/casbin-gateway/proxy"
)

const (
	// codexOriginator names the client this token was granted to, which the
	// endpoint wants on every request made with it.
	codexOriginator = "codex_cli_rs"
	usageTimeout    = 15 * time.Second
	maxUsageReply   = 1 << 20
)

// Vars so a test can answer the usage call locally.
var (
	codexUsageUrl = "https://chatgpt.com/backend-api/wham/usage"
	usageClient   = &http.Client{Timeout: usageTimeout, Transport: proxy.Transport()}
)

// Window is one rate limit window: what is left of it, and when it resets.
type Window struct {
	RemainingPercent float64 `json:"remainingPercent"`
	WindowMinutes    int     `json:"windowMinutes"`
	ResetTime        string  `json:"resetTime,omitempty"`
}

// Usage is what one account has left. Primary is the short window - five hours
// on the Codex and Claude plans today - and Secondary the weekly one.
type Usage struct {
	Primary   *Window `json:"primary,omitempty"`
	Secondary *Window `json:"secondary,omitempty"`
	// ResetCredits is how many rate limit resets the account still holds. Each
	// one empties a window before its own time.
	ResetCredits int    `json:"resetCredits"`
	Plan         string `json:"plan,omitempty"`
}

// codexWindow is one window as the endpoint reports it: what has been spent
// rather than what is left, and the reset both as a moment and as a countdown.
type codexWindow struct {
	UsedPercent   float64 `json:"used_percent"`
	WindowSeconds int64   `json:"limit_window_seconds"`
	ResetAfter    int64   `json:"reset_after_seconds"`
	ResetAt       int64   `json:"reset_at"`
}

// Codex reads what one Codex credential's account has left. The credential is
// the whole auth.json the agent keeps; an API key has no account to ask about.
func Codex(credential string) (*Usage, error) {
	var auth struct {
		Tokens struct {
			AccessToken string `json:"access_token"`
			AccountId   string `json:"account_id"`
		} `json:"tokens"`
	}
	if err := json.Unmarshal([]byte(credential), &auth); err != nil {
		return nil, err
	}
	if auth.Tokens.AccessToken == "" {
		return nil, errors.New("this account holds no ChatGPT sign-in, so it has no plan limits to report")
	}
	return CodexToken(auth.Tokens.AccessToken, auth.Tokens.AccountId)
}

// CodexToken is the same reading for a ChatGPT access token held anywhere.
func CodexToken(accessToken string, accountId string) (*Usage, error) {
	request, err := http.NewRequest(http.MethodGet, codexUsageUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Originator", codexOriginator)
	request.Header.Set("User-Agent", codexOriginator)
	if accountId != "" {
		request.Header.Set("Chatgpt-Account-Id", accountId)
	}

	response, err := usageClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// A refused sign-in is answered with an HTML sign-in page, not with numbers.
	if response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("the stored sign-in is no longer accepted (%s), sign in again", response.Status)
	}
	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("chatgpt answered %s", response.Status)
	}

	var answer struct {
		PlanType  string `json:"plan_type"`
		RateLimit struct {
			Primary   *codexWindow `json:"primary_window"`
			Secondary *codexWindow `json:"secondary_window"`
		} `json:"rate_limit"`
		ResetCredits struct {
			AvailableCount int `json:"available_count"`
		} `json:"rate_limit_reset_credits"`
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, maxUsageReply)).Decode(&answer); err != nil {
		return nil, err
	}
	if answer.RateLimit.Primary == nil && answer.RateLimit.Secondary == nil {
		return nil, errors.New("chatgpt reported no limits for this account")
	}

	return &Usage{
		Primary:      codexWindowOf(answer.RateLimit.Primary),
		Secondary:    codexWindowOf(answer.RateLimit.Secondary),
		ResetCredits: answer.ResetCredits.AvailableCount,
		Plan:         answer.PlanType,
	}, nil
}

func codexWindowOf(window *codexWindow) *Window {
	if window == nil {
		return nil
	}

	// A window nothing has been spent in has not started: the reset it reports
	// is only where it would end if it started now, so none is given.
	idle := window.UsedPercent == 0 && window.WindowSeconds > 0 && window.ResetAfter >= window.WindowSeconds
	reset := ""
	switch {
	case idle:
	case window.ResetAt > 0:
		// The moment beats the countdown, which the local clock can skew.
		reset = time.Unix(window.ResetAt, 0).UTC().Format(time.RFC3339)
	case window.ResetAfter > 0:
		reset = time.Now().Add(time.Duration(window.ResetAfter) * time.Second).UTC().Format(time.RFC3339)
	}

	remaining := math.Round((100-window.UsedPercent)*10) / 10
	return &Window{
		RemainingPercent: math.Min(100, math.Max(0, remaining)),
		WindowMinutes:    int(window.WindowSeconds / 60),
		ResetTime:        reset,
	}
}
