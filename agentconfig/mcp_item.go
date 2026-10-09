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

package agentconfig

import (
	"crypto/sha256"
	"encoding/json"
	"strings"

	"github.com/apache/casbin-gateway/auditutil"
)

// mcpItems summarizes one file's entries for a list. The fields below are the
// ones every MCP host spells the same way; anything else an agent keeps in an
// entry is preserved on copy and shown in the detail view.
func mcpItems(agentId string, owner string, file string, entries map[string]map[string]any) []*Item {
	items := make([]*Item, 0, len(entries))
	for name, entry := range entries {
		item := &Item{
			AgentId:   agentId,
			Owner:     owner,
			Kind:      KindMcp,
			Name:      name,
			Path:      file,
			Command:   commandOf(entry),
			Url:       stringField(entry, "url", "serverUrl", "endpoint"),
			Managed:   name == ManagedEntryName,
			Transport: transportOf(entry),
			Digest:    entryDigest(entry),
		}
		items = append(items, item)
	}
	sortItems(items)
	return items
}

// mcpDetail shows one entry as it is written in the file, with credentials
// masked. Copying reads the file again, so what is migrated is the real entry
// rather than this redacted view of it.
func mcpDetail(agentId string, owner string, file string, name string, entry map[string]any) (*Detail, error) {
	redacted := auditutil.SanitizeValue("", entry)
	content, err := json.MarshalIndent(redacted, "", "  ")
	if err != nil {
		return nil, err
	}

	items := mcpItems(agentId, owner, file, map[string]map[string]any{name: entry})
	return &Detail{Item: items[0], Content: string(content)}, nil
}

// entryDigest identifies one server's definition, so the same name in two
// agents can be told apart when the two definitions differ. Marshalling sorts
// the keys, so the digest does not depend on how the file was written.
func entryDigest(entry map[string]any) string {
	raw, err := json.Marshal(entry)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(raw)
	return shortDigest(sum[:])
}

// transportOf reports how the agent reaches the server. An entry that names its
// own type is believed; otherwise a command means the server is spawned locally
// and a URL means it is reached over HTTP.
func transportOf(entry map[string]any) string {
	if declared := strings.ToLower(stringField(entry, "type", "transport")); declared != "" {
		return declared
	}
	if commandOf(entry) != "" {
		return "stdio"
	}
	if stringField(entry, "url", "serverUrl", "endpoint") != "" {
		return "http"
	}
	return ""
}

// commandOf renders the spawn command the way a shell would show it, so a list
// row says which binary runs rather than only that the server is local.
func commandOf(entry map[string]any) string {
	command := stringField(entry, "command")
	if command == "" {
		return ""
	}

	arguments, _ := entry["args"].([]any)
	parts := make([]string, 0, len(arguments)+1)
	parts = append(parts, command)
	redactNext := false
	for _, argument := range arguments {
		text, ok := argument.(string)
		if !ok {
			continue
		}
		if redactNext {
			redactNext = false
			parts = append(parts, "[REDACTED]")
			continue
		}

		flag, _, hasValue := strings.Cut(text, "=")
		if strings.HasPrefix(flag, "-") && auditutil.SensitiveKey(flag) {
			if hasValue {
				text = flag + "=[REDACTED]"
			} else {
				redactNext = true
			}
		}
		parts = append(parts, text)
	}
	return auditutil.SanitizeString(strings.Join(parts, " "))
}

func stringField(entry map[string]any, keys ...string) string {
	for _, key := range keys {
		if text, ok := entry[key].(string); ok && text != "" {
			return text
		}
	}
	return ""
}
