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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/apache/casbin-gateway/agentusage"
	"github.com/apache/casbin-gateway/proxy"
)

// A Copilot plan is signed in to with GitHub's device flow, under the client id
// the Copilot editor plugins use: the GitHub token it grants is exchanged for a
// short-lived Copilot token, which is what the requests carry.
const (
	copilotClientId      = "Iv1.b507a08c87ecfe98"
	copilotBaseUrl       = "https://api.githubcopilot.com"
	copilotTokenUrl      = "https://api.github.com/copilot_internal/v2/token"
	copilotUserUrl       = "https://api.github.com/user"
	copilotEditorVersion = "vscode/1.107.0"
	copilotPluginVersion = "copilot-chat/0.35.0"
)

func startDeviceLogin(vendor *SubscriptionVendor) (*SubscriptionLogin, error) {
	form := url.Values{}
	form.Set("client_id", vendor.clientId)
	form.Set("scope", vendor.scope)

	var device struct {
		DeviceCode      string `json:"device_code"`
		UserCode        string `json:"user_code"`
		VerificationUri string `json:"verification_uri"`
		ExpiresIn       int    `json:"expires_in"`
		Interval        int    `json:"interval"`
		Error           string `json:"error"`
		ErrorText       string `json:"error_description"`
	}
	if err := postGithubForm(vendor.deviceCodeUrl, form, &device); err != nil {
		return nil, err
	}
	if device.Error != "" {
		return nil, fmt.Errorf("github refused the sign-in: %s %s", device.Error, device.ErrorText)
	}
	if device.DeviceCode == "" || device.UserCode == "" {
		return nil, fmt.Errorf("github returned no code to sign in with")
	}

	id, err := randomString()
	if err != nil {
		return nil, err
	}
	session := &loginSession{
		SubscriptionLogin: SubscriptionLogin{
			Id:       id,
			Vendor:   vendor.Id,
			Url:      device.VerificationUri,
			UserCode: device.UserCode,
			Running:  true,
		},
		started: time.Now(),
	}

	cancelRunningLogins(vendor.Id)
	logins.Lock()
	forgetFinishedLoginsLocked()
	logins.byId[id] = session
	logins.Unlock()

	expires := time.Duration(device.ExpiresIn) * time.Second
	if expires <= 0 || expires > pendingTtl {
		expires = pendingTtl
	}
	go pollDeviceLogin(session, vendor, device.DeviceCode, device.Interval, time.Now().Add(expires))

	answer := session.snapshot()
	return &answer, nil
}

// pollDeviceLogin asks GitHub, at the pace it set, whether the code was entered.
func pollDeviceLogin(session *loginSession, vendor *SubscriptionVendor, deviceCode string, interval int, deadline time.Time) {
	if interval <= 0 {
		interval = 5
	}

	for time.Now().Before(deadline) {
		time.Sleep(time.Duration(interval) * time.Second)

		session.Lock()
		running := session.Running
		session.Unlock()
		if !running {
			return
		}

		form := url.Values{}
		form.Set("client_id", vendor.clientId)
		form.Set("device_code", deviceCode)
		form.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

		var answer struct {
			AccessToken string `json:"access_token"`
			Error       string `json:"error"`
			ErrorText   string `json:"error_description"`
		}
		if err := postGithubForm(vendor.tokenUrl, form, &answer); err != nil {
			continue
		}

		switch answer.Error {
		case "":
		case "authorization_pending":
			continue
		case "slow_down":
			interval += 5
			continue
		default:
			session.finish(nil, fmt.Errorf("%s %s", answer.Error, answer.ErrorText))
			return
		}
		if answer.AccessToken == "" {
			continue
		}

		granted, err := renewCopilot(&ProviderSubscription{Vendor: vendor.Id, RefreshToken: answer.AccessToken})
		if err == nil {
			describeCopilot(granted)
			models := copilotModels(granted)
			session.Lock()
			session.Models = models
			session.Unlock()
		}
		session.finish(granted, err)
		return
	}

	session.finish(nil, fmt.Errorf("the sign-in was not finished in time"))
}

func postGithubForm(endpoint string, form url.Values, target any) error {
	request, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: exchangeTimeout, Transport: proxy.Transport()}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(io.LimitReader(response.Body, 1024*1024))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("github answered %s: %s", response.Status, summarize(body))
	}
	return nil
}

// renewCopilot exchanges the GitHub token for a Copilot token, which lives for
// about half an hour and names the API host this account is served from.
func renewCopilot(subscription *ProviderSubscription) (*ProviderSubscription, error) {
	request, err := http.NewRequest(http.MethodGet, copilotTokenUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "token "+subscription.RefreshToken)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Editor-Version", copilotEditorVersion)
	request.Header.Set("Editor-Plugin-Version", copilotPluginVersion)

	client := &http.Client{Timeout: exchangeTimeout, Transport: proxy.Transport()}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(io.LimitReader(response.Body, 1024*1024))
	if err != nil {
		return nil, err
	}
	if response.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("github no longer accepts this sign-in, sign in again")
	}
	if response.StatusCode == http.StatusForbidden || response.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("this GitHub account has no Copilot plan: %s", summarize(body))
	}
	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("github answered %s: %s", response.Status, summarize(body))
	}

	var answer struct {
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expires_at"`
		Endpoints struct {
			Api string `json:"api"`
		} `json:"endpoints"`
	}
	if err := json.Unmarshal(body, &answer); err != nil || answer.Token == "" {
		return nil, fmt.Errorf("github returned no Copilot token: %s", summarize(body))
	}

	renewed := *subscription
	renewed.AccessToken = answer.Token
	renewed.Endpoint = strings.TrimRight(answer.Endpoints.Api, "/")
	renewed.ExpiresAt = ""
	if answer.ExpiresAt > 0 {
		renewed.ExpiresAt = time.Unix(answer.ExpiresAt, 0).UTC().Format(time.RFC3339)
	}
	return &renewed, nil
}

// describeCopilot fills in whose account this is and which plan it is on. What
// cannot be read is left empty rather than failing a sign-in that worked.
func describeCopilot(subscription *ProviderSubscription) {
	request, err := http.NewRequest(http.MethodGet, copilotUserUrl, nil)
	if err == nil {
		request.Header.Set("Authorization", "token "+subscription.RefreshToken)
		request.Header.Set("Accept", "application/json")
		client := &http.Client{Timeout: exchangeTimeout, Transport: proxy.Transport()}
		if response, err := client.Do(request); err == nil {
			var user struct {
				Login string `json:"login"`
				Email string `json:"email"`
			}
			if json.NewDecoder(io.LimitReader(response.Body, 1024*1024)).Decode(&user) == nil {
				fillClaim(&subscription.Account, user.Login)
			}
			response.Body.Close()
		}
	}

	if usage, err := agentusage.Copilot(subscription.RefreshToken); err == nil {
		fillClaim(&subscription.Plan, usage.Plan)
	}
}

// copilotModels lists the chat models the account may pick. A Copilot model is
// named its own way (claude-sonnet-4.5), so the list is what lets an agent's
// own model names be matched to one.
func copilotModels(subscription *ProviderSubscription) []string {
	base := subscription.Endpoint
	if base == "" {
		base = copilotBaseUrl
	}
	request, err := http.NewRequest(http.MethodGet, base+"/models", nil)
	if err != nil {
		return nil
	}
	applyCopilot(request.Header, subscription)

	client := &http.Client{Timeout: exchangeTimeout, Transport: proxy.Transport()}
	response, err := client.Do(request)
	if err != nil {
		return nil
	}
	defer response.Body.Close()
	if response.StatusCode >= 400 {
		return nil
	}

	var answer struct {
		Data []struct {
			Id           string `json:"id"`
			PickerOn     bool   `json:"model_picker_enabled"`
			Capabilities struct {
				Type string `json:"type"`
			} `json:"capabilities"`
		} `json:"data"`
	}
	if json.NewDecoder(io.LimitReader(response.Body, maxProbeBody)).Decode(&answer) != nil {
		return nil
	}

	models := []string{}
	seen := map[string]bool{}
	for _, model := range answer.Data {
		if !model.PickerOn || (model.Capabilities.Type != "" && model.Capabilities.Type != "chat") || seen[model.Id] {
			continue
		}
		seen[model.Id] = true
		models = append(models, model.Id)
	}
	return models
}

func applyCopilot(header http.Header, subscription *ProviderSubscription) {
	header.Set("Authorization", "Bearer "+subscription.AccessToken)
	header.Set("Copilot-Integration-Id", "vscode-chat")
	header.Set("Editor-Version", copilotEditorVersion)
	header.Set("Editor-Plugin-Version", copilotPluginVersion)
	header.Set("Openai-Intent", "conversation-edits")
}

// copilotRequestHeaders says who started the turn. A turn the agent took on its
// own - answering a tool result - is not a premium request the person made, and
// Copilot bills it that way only when told.
func copilotRequestHeaders(header http.Header, body []byte) {
	var request struct {
		Messages []struct {
			Role    string          `json:"role"`
			Content json.RawMessage `json:"content"`
		} `json:"messages"`
	}
	if json.Unmarshal(body, &request) != nil || len(request.Messages) == 0 {
		return
	}

	initiator := "agent"
	if last := request.Messages[len(request.Messages)-1]; last.Role == "user" {
		initiator = "user"
	}
	header.Set("X-Initiator", initiator)

	for _, message := range request.Messages {
		if strings.Contains(string(message.Content), `"image_url"`) {
			header.Set("Copilot-Vision-Request", "true")
			break
		}
	}
}
