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

package agentusage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

var (
	claudeUsageUrl  = "https://api.anthropic.com/api/oauth/usage"
	copilotUsageUrl = "https://api.github.com/copilot_internal/user"
)

// Claude reads what a Claude subscription has left, with the OAuth token the
// sign-in was granted. The token needs the user:profile scope.
func Claude(accessToken string) (*Usage, error) {
	request, err := http.NewRequest(http.MethodGet, claudeUsageUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Anthropic-Version", "2023-06-01")
	request.Header.Set("Anthropic-Beta", "oauth-2025-04-20")

	var answer struct {
		FiveHour *claudeWindow `json:"five_hour"`
		SevenDay *claudeWindow `json:"seven_day"`
	}
	if err := getUsage(request, "claude", &answer); err != nil {
		return nil, err
	}
	if answer.FiveHour == nil && answer.SevenDay == nil {
		return nil, errors.New("claude reported no limits for this account")
	}
	return &Usage{
		Primary:   answer.FiveHour.window(300),
		Secondary: answer.SevenDay.window(10080),
	}, nil
}

// claudeWindow reports what has been spent, in percent, and no reset at all for
// a window that has not started.
type claudeWindow struct {
	Utilization float64 `json:"utilization"`
	ResetsAt    *string `json:"resets_at"`
}

func (window *claudeWindow) window(minutes int) *Window {
	if window == nil {
		return nil
	}

	reset := ""
	if window.ResetsAt != nil {
		if at, err := time.Parse(time.RFC3339, *window.ResetsAt); err == nil {
			reset = at.UTC().Format(time.RFC3339)
		}
	}
	return &Window{
		RemainingPercent: clampPercent(100 - window.Utilization),
		WindowMinutes:    minutes,
		ResetTime:        reset,
	}
}

// Copilot reads what a GitHub Copilot plan has left this month, with the GitHub
// token the device sign-in was granted. Primary is the premium requests, and
// Secondary the chat allowance of the free plan; a limit the plan does not have
// is left out.
func Copilot(githubToken string) (*Usage, error) {
	request, err := http.NewRequest(http.MethodGet, copilotUsageUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "token "+githubToken)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("X-Github-Api-Version", "2025-04-01")

	var answer struct {
		Plan      string `json:"copilot_plan"`
		ResetDate string `json:"quota_reset_date"`
		Snapshots struct {
			Premium *copilotSnapshot `json:"premium_interactions"`
			Chat    *copilotSnapshot `json:"chat"`
		} `json:"quota_snapshots"`
	}
	if err := getUsage(request, "github", &answer); err != nil {
		return nil, err
	}

	reset := ""
	if at, err := time.Parse("2006-01-02", answer.ResetDate); err == nil {
		reset = at.UTC().Format(time.RFC3339)
	}
	return &Usage{
		Primary:   answer.Snapshots.Premium.window(reset),
		Secondary: answer.Snapshots.Chat.window(reset),
		Plan:      answer.Plan,
	}, nil
}

type copilotSnapshot struct {
	PercentRemaining float64 `json:"percent_remaining"`
	Unlimited        bool    `json:"unlimited"`
}

func (snapshot *copilotSnapshot) window(reset string) *Window {
	if snapshot == nil || snapshot.Unlimited {
		return nil
	}
	return &Window{RemainingPercent: clampPercent(snapshot.PercentRemaining), ResetTime: reset}
}

// getUsage performs the GET and decodes the answer into target.
func getUsage(request *http.Request, vendor string, target any) error {
	response, err := usageClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusForbidden {
		return fmt.Errorf("the stored sign-in is no longer accepted (%s), sign in again", response.Status)
	}
	if response.StatusCode >= 400 {
		return fmt.Errorf("%s answered %s", vendor, response.Status)
	}
	return json.NewDecoder(io.LimitReader(response.Body, maxUsageReply)).Decode(target)
}

func clampPercent(value float64) float64 {
	return math.Min(100, math.Max(0, math.Round(value*10)/10))
}
