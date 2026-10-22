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

package object

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/apache/casbin-gateway/protocol"
	"github.com/apache/casbin-gateway/proxy"
	"github.com/apache/casbin-gateway/util"
)

// WarmupConfig starts a subscription's short window at set times of the day. A
// Claude or ChatGPT plan counts its five hours from the first request, so one
// sent early in the morning has the window reset by midday, and the afternoon
// starts on a full one.
type WarmupConfig struct {
	Enabled bool `json:"enabled"`
	// Times are local, "HH:MM".
	Times []string `json:"times"`
	// Model is what is asked for. Empty takes the model this provider last
	// relayed, or the vendor's smallest.
	Model string `json:"model"`
}

// WarmupResult is how one warm-up went.
type WarmupResult struct {
	Time    string `json:"time"`
	Ok      bool   `json:"ok"`
	Model   string `json:"model,omitempty"`
	Message string `json:"message"`
}

var warmupTime = regexp.MustCompile(`^([01][0-9]|2[0-3]):[0-5][0-9]$`)

// normalizeWarmup drops what a schedule cannot mean, and refuses a time that
// is not one.
func normalizeWarmup(provider *Provider) error {
	config := provider.Warmup
	if config == nil {
		return nil
	}
	if !UsesSubscription(provider) {
		provider.Warmup = nil
		return nil
	}

	times := []string{}
	seen := map[string]bool{}
	for _, value := range config.Times {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		if !warmupTime.MatchString(value) {
			return fmt.Errorf("the warm-up time %q is not a time of day written as HH:MM", value)
		}
		seen[value] = true
		times = append(times, value)
	}
	config.Times = times
	config.Model = strings.TrimSpace(config.Model)
	return nil
}

var warmups = struct {
	sync.Mutex
	last  map[string]*WarmupResult
	fired map[string]bool
}{last: map[string]*WarmupResult{}, fired: map[string]bool{}}

func lastWarmup(providerId string) *WarmupResult {
	warmups.Lock()
	defer warmups.Unlock()
	return warmups.last[providerId]
}

func rememberWarmup(providerId string, result *WarmupResult) *WarmupResult {
	warmups.Lock()
	warmups.last[providerId] = result
	warmups.Unlock()
	return result
}

// WarmUpProvider sends a subscription the smallest request its endpoint takes,
// which is what starts the plan's short window.
func WarmUpProvider(provider *Provider) *WarmupResult {
	result := &WarmupResult{Time: util.GetCurrentTime()}
	model, err := warmUp(provider)
	result.Model = model
	switch {
	case err != nil:
		result.Message = err.Error()
	case runningShortWindow(provider.GetId(), time.Now()) != "":
		result.Ok = true
		result.Message = "the window runs until " + runningShortWindow(provider.GetId(), time.Now())
	default:
		result.Ok = true
		result.Message = "the request was answered"
	}
	return rememberWarmup(provider.GetId(), result)
}

func warmUp(provider *Provider) (string, error) {
	vendor := vendorOfProvider(provider)
	if vendor == nil {
		return "", errors.New("this provider holds no sign-in to warm up")
	}
	if !vendor.warms {
		return "", fmt.Errorf("a %s plan has no window that a request starts", vendor.Label)
	}
	if err := EnsureSubscription(provider); err != nil {
		return "", err
	}

	model := ""
	if provider.Warmup != nil {
		model = provider.Warmup.Model
	}
	if model == "" {
		model = lastRelayedModel(provider.GetId())
	}
	if model == "" {
		model = vendor.warmModel
	}
	if model == "" {
		return "", errors.New("name the model to warm up with: this provider has not relayed a request yet")
	}

	name := ProviderProtocol(provider)
	upstream, err := protocol.UpstreamOf(name)
	if err != nil {
		return model, err
	}
	target, err := BuildProviderUrl(UpstreamBaseUrl(provider), name, upstream.Endpoint())
	if err != nil {
		return model, err
	}
	target = AppendQuery(target, SubscriptionQuery(provider))

	var body map[string]any
	switch name {
	case ProtocolAnthropic:
		body = map[string]any{
			"model":      model,
			"max_tokens": 1,
			"messages":   []any{map[string]any{"role": "user", "content": "hi"}},
		}
	case ProtocolResponses:
		body = map[string]any{
			"model": model,
			"input": []any{map[string]any{
				"role":    "user",
				"content": []any{map[string]any{"type": "input_text", "text": "hi"}},
			}},
			"store":  false,
			"stream": true,
		}
	default:
		body = map[string]any{
			"model":      model,
			"max_tokens": 1,
			"messages":   []any{map[string]any{"role": "user", "content": "hi"}},
		}
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		return model, err
	}
	encoded = ShapeSubscriptionBody(provider, encoded)

	request, err := http.NewRequest(http.MethodPost, target, bytes.NewReader(encoded))
	if err != nil {
		return model, err
	}
	request.Header.Set("Content-Type", "application/json")
	SetProviderAuth(request.Header, provider)

	client := &http.Client{Timeout: time.Minute, Transport: proxy.Transport()}
	response, err := client.Do(request)
	if err != nil {
		return model, err
	}
	defer response.Body.Close()
	answer, _ := io.ReadAll(io.LimitReader(response.Body, maxProbeBody))

	ObserveSubscriptionLimits(provider, response.StatusCode, response.Header)
	if response.StatusCode >= 300 {
		return model, fmt.Errorf("the vendor answered %s%s", response.Status, probeDetail(answer))
	}
	return model, nil
}

// lastRelayedModel is the model this provider last answered for, which is one
// its plan certainly serves.
func lastRelayedModel(providerId string) string {
	if ormer == nil || ormer.Engine == nil {
		return ""
	}
	record := &LlmRecord{}
	existed, err := ormer.Engine.Where("provider = ?", providerId).
		And("status >= 200 and status < 300").
		Desc("created_time").Get(record)
	if err != nil || !existed {
		return ""
	}
	// A request relayed as it arrived names its model only once.
	if record.UpstreamModel != "" {
		return record.UpstreamModel
	}
	return record.Model
}

// StartWarmupSchedule checks every half minute for a warm-up that is due.
func StartWarmupSchedule() {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for now := range ticker.C {
			runDueWarmups(now)
		}
	}()
}

func runDueWarmups(now time.Time) {
	providers, err := listEnabledProviders()
	if err != nil {
		return
	}
	decryptProviders(providers)

	clock := now.Format("15:04")
	day := now.Format("2006-01-02")
	for _, provider := range providers {
		config := provider.Warmup
		if !UsesSubscription(provider) || config == nil || !config.Enabled {
			continue
		}
		for _, at := range config.Times {
			if at != clock {
				continue
			}
			key := provider.GetId() + "|" + day + "|" + at
			warmups.Lock()
			fired := warmups.fired[key]
			warmups.fired[key] = true
			warmups.Unlock()
			if fired {
				continue
			}

			go warmUpIfIdle(provider)
		}
	}
	forgetOldWarmups(day)
}

// warmUpIfIdle warms a subscription up unless its short window is already
// running, which another request would not start again. What is left is asked
// of the vendor first when nothing is known, which costs none of the plan.
func warmUpIfIdle(provider *Provider) {
	RefreshProviderQuotas([]*Provider{provider}, false)
	if running := runningShortWindow(provider.GetId(), time.Now()); running != "" {
		rememberWarmup(provider.GetId(), &WarmupResult{
			Time:    util.GetCurrentTime(),
			Ok:      true,
			Message: "the window was already running, until " + running,
		})
		return
	}
	WarmUpProvider(provider)
}

func runningShortWindow(providerId string, now time.Time) string {
	for _, window := range cachedWindows(providerId) {
		if window.Name == WindowShort && window.reset().After(now) {
			return window.ResetTime
		}
	}
	return ""
}

func forgetOldWarmups(day string) {
	warmups.Lock()
	defer warmups.Unlock()
	for key := range warmups.fired {
		if !strings.Contains(key, "|"+day+"|") {
			delete(warmups.fired, key)
		}
	}
}
