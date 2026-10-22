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
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/apache/casbin-gateway/agentusage"
	"github.com/apache/casbin-gateway/util"
)

// Window names. The short and the weekly window are what a plan runs out of;
// premium and chat are a month's allowance, which a subscription is not
// refused for when it is spent.
const (
	WindowShort   = "5h"
	WindowWeek    = "week"
	WindowPremium = "premium"
	WindowChat    = "chat"
)

// QuotaWindow is one limit of a subscription: what is left of it, and when it
// fills again. An empty ResetTime is a window that has not started.
type QuotaWindow struct {
	Name             string  `json:"name"`
	RemainingPercent float64 `json:"remainingPercent"`
	WindowMinutes    int     `json:"windowMinutes"`
	ResetTime        string  `json:"resetTime,omitempty"`
}

func (window *QuotaWindow) reset() time.Time {
	at, err := time.Parse(time.RFC3339, window.ResetTime)
	if err != nil {
		return time.Time{}
	}
	return at
}

// spent reports whether the plan refuses requests until this window resets.
func (window *QuotaWindow) spent(now time.Time) bool {
	if window.Name != WindowShort && window.Name != WindowWeek {
		return false
	}
	return window.RemainingPercent <= 0 && window.reset().After(now)
}

func windowName(minutes int, fallback string) string {
	switch {
	case minutes == 300:
		return WindowShort
	case minutes == 10080:
		return WindowWeek
	case minutes > 0 && minutes%60 == 0:
		return fmt.Sprintf("%dh", minutes/60)
	}
	return fallback
}

// usageWindows names the two windows a vendor's usage endpoint reports.
func usageWindows(usage *agentusage.Usage, primary string, secondary string) []*QuotaWindow {
	windows := []*QuotaWindow{}
	for _, item := range []struct {
		window *agentusage.Window
		name   string
	}{{usage.Primary, primary}, {usage.Secondary, secondary}} {
		if item.window == nil {
			continue
		}
		name := item.name
		if name == "" {
			name = windowName(item.window.WindowMinutes, WindowShort)
		}
		windows = append(windows, &QuotaWindow{
			Name:             name,
			RemainingPercent: item.window.RemainingPercent,
			WindowMinutes:    item.window.WindowMinutes,
			ResetTime:        item.window.ResetTime,
		})
	}
	return windows
}

// subscriptionQuota asks the vendor what the plan behind a held sign-in has
// left, with that sign-in's own token.
func subscriptionQuota(provider *Provider, quota *ProviderQuota) *ProviderQuota {
	vendor := vendorOfProvider(provider)
	if vendor == nil {
		quota.Error = "this provider is not signed in"
		return quota
	}
	if vendor.usage == nil {
		return quota
	}
	quota.Supported = true

	if err := EnsureSubscription(provider); err != nil {
		quota.Error = err.Error()
		return quota
	}
	subscription := subscriptionOf(provider)
	windows, plan, err := vendor.usage(subscription)
	if err != nil {
		quota.Error = err.Error()
		return quota
	}
	quota.Windows = windows
	quota.Plan = plan
	if quota.Plan == "" {
		quota.Plan = subscription.Plan
	}
	return quota
}

// limitBlocks holds the providers a plan refused, until the moment it said it
// would take requests again. It is what keeps a spent account out of the way
// for the rest of its window rather than for the breaker's few minutes.
var limitBlocks = struct {
	sync.Mutex
	until map[string]time.Time
}{until: map[string]time.Time{}}

// ObserveSubscriptionLimits reads the limits a subscription endpoint reports on
// every answer, which keeps what is left current without asking for it.
func ObserveSubscriptionLimits(provider *Provider, statusCode int, header http.Header) {
	vendor := vendorOfProvider(provider)
	if vendor == nil || vendor.limits == nil {
		return
	}

	now := time.Now()
	windows, plan := vendor.limits(header)
	if len(windows) > 0 {
		storeObservedWindows(provider.GetId(), windows, plan)
	}

	if statusCode != http.StatusTooManyRequests {
		return
	}
	until := time.Time{}
	for _, window := range windows {
		if window.spent(now) && window.reset().After(until) {
			until = window.reset()
		}
	}
	if until.IsZero() {
		if seconds, err := strconv.Atoi(strings.TrimSpace(header.Get("Retry-After"))); err == nil && seconds > 0 {
			until = now.Add(time.Duration(seconds) * time.Second)
		}
	}
	if until.IsZero() {
		return
	}
	limitBlocks.Lock()
	limitBlocks.until[provider.GetId()] = until
	limitBlocks.Unlock()
}

// storeObservedWindows puts what an answer reported in the quota cache, beside
// whatever the usage endpoint said about the windows the answer did not carry.
func storeObservedWindows(providerId string, windows []*QuotaWindow, plan string) {
	providerQuotaMutex.Lock()
	defer providerQuotaMutex.Unlock()

	quota := &ProviderQuota{Provider: providerId, Supported: true}
	if cached, ok := providerQuotaMap[providerId]; ok {
		copied := *cached.quota
		quota = &copied
	}
	merged := append([]*QuotaWindow{}, windows...)
	for _, existing := range quota.Windows {
		found := false
		for _, window := range windows {
			found = found || window.Name == existing.Name
		}
		if !found {
			merged = append(merged, existing)
		}
	}
	quota.Windows = merged
	quota.Supported = true
	quota.Error = ""
	if plan != "" {
		quota.Plan = plan
	}
	quota.Time = util.GetCurrentTime()
	providerQuotaMap[providerId] = &cachedQuota{quota: quota, time: time.Now()}
}

// IsSubscriptionSpent reports whether the plan behind a provider has nothing
// left until one of its windows resets. Like a provider in cooldown it is tried
// last rather than dropped.
func IsSubscriptionSpent(providerId string) bool {
	now := time.Now()
	limitBlocks.Lock()
	until, blocked := limitBlocks.until[providerId]
	if blocked && !until.After(now) {
		delete(limitBlocks.until, providerId)
		blocked = false
	}
	limitBlocks.Unlock()
	if blocked {
		return true
	}

	for _, window := range cachedWindows(providerId) {
		if window.spent(now) {
			return true
		}
	}
	return false
}

func cachedWindows(providerId string) []*QuotaWindow {
	providerQuotaMutex.Lock()
	defer providerQuotaMutex.Unlock()
	if cached, ok := providerQuotaMap[providerId]; ok {
		return cached.quota.Windows
	}
	return nil
}

// headerFloat reads a number a header carries, false when there is none.
func headerFloat(header http.Header, name string) (float64, bool) {
	value := strings.TrimSpace(header.Get(name))
	if value == "" {
		return 0, false
	}
	number, err := strconv.ParseFloat(value, 64)
	return number, err == nil
}

func unixTime(seconds float64) string {
	if seconds <= 0 {
		return ""
	}
	return time.Unix(int64(seconds), 0).UTC().Format(time.RFC3339)
}

// claudeLimits reads the windows Anthropic reports on every answer to a
// subscription token, as the fraction spent and the reset in Unix seconds.
func claudeLimits(header http.Header) ([]*QuotaWindow, string) {
	windows := []*QuotaWindow{}
	for _, item := range []struct {
		key     string
		name    string
		minutes int
	}{{"5h", WindowShort, 300}, {"7d", WindowWeek, 10080}} {
		prefix := "Anthropic-Ratelimit-Unified-" + item.key
		used, hasUsed := headerFloat(header, prefix+"-Utilization")
		reset, hasReset := headerFloat(header, prefix+"-Reset")
		if !hasUsed && !hasReset {
			continue
		}
		windows = append(windows, &QuotaWindow{
			Name:             item.name,
			RemainingPercent: clampPercent(100 - used*100),
			WindowMinutes:    item.minutes,
			ResetTime:        unixTime(reset),
		})
	}
	return windows, ""
}

// codexLimits reads the windows the Codex backend reports on every answer, in
// the same terms as its usage endpoint.
func codexLimits(header http.Header) ([]*QuotaWindow, string) {
	windows := []*QuotaWindow{}
	for _, item := range []struct {
		key      string
		fallback string
	}{{"primary", WindowShort}, {"secondary", WindowWeek}} {
		prefix := "X-Codex-" + item.key
		used, hasUsed := headerFloat(header, prefix+"-Used-Percent")
		if !hasUsed {
			continue
		}
		minutes, _ := headerFloat(header, prefix+"-Window-Minutes")
		resetAt, _ := headerFloat(header, prefix+"-Reset-At")
		resetAfter, _ := headerFloat(header, prefix+"-Reset-After-Seconds")

		reset := unixTime(resetAt)
		if used == 0 && minutes > 0 && resetAfter >= minutes*60 {
			reset = ""
		}
		windows = append(windows, &QuotaWindow{
			Name:             windowName(int(minutes), item.fallback),
			RemainingPercent: clampPercent(100 - used),
			WindowMinutes:    int(minutes),
			ResetTime:        reset,
		})
	}
	return windows, header.Get("X-Codex-Plan-Type")
}

func clampPercent(value float64) float64 {
	rounded := float64(int64(value*10+0.5)) / 10
	if value < 0 {
		rounded = 0
	}
	if rounded > 100 {
		return 100
	}
	return rounded
}
