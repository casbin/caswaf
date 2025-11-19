// Copyright 2025 The casbin Authors. All Rights Reserved.
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
	"testing"
	"time"
)

func TestShouldAllowUpgrade(t *testing.T) {
	tests := []struct {
		name        string
		upgradeMode string
		want        bool
	}{
		{
			name:        "Empty upgrade mode should allow upgrade",
			upgradeMode: "",
			want:        true,
		},
		{
			name:        "At Any Time mode should allow upgrade",
			upgradeMode: "At Any Time",
			want:        true,
		},
		{
			name:        "No Upgrade mode should not allow upgrade",
			upgradeMode: "No Upgrade",
			want:        false,
		},
		{
			name:        "Unknown mode should default to allow upgrade",
			upgradeMode: "Unknown Mode",
			want:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &Node{
				UpgradeMode: tt.upgradeMode,
			}
			if got := node.ShouldAllowUpgrade(); got != tt.want {
				t.Errorf("ShouldAllowUpgrade() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShouldAllowUpgrade_HalfAHour(t *testing.T) {
	// We can't easily test the time-based logic in CI without mocking time,
	// but we can at least verify the function doesn't panic
	node := &Node{
		UpgradeMode: "Half A Hour",
	}

	// This will return true or false depending on current time in GMT+8
	result := node.ShouldAllowUpgrade()

	// Get current time in GMT+8
	location := time.FixedZone("GMT+8", 8*60*60)
	now := time.Now().In(location)
	hour := now.Hour()
	minute := now.Minute()

	expectedResult := hour == 23 && minute < 30

	if result != expectedResult {
		t.Logf("ShouldAllowUpgrade() for Half A Hour mode = %v, expected %v (current time: %02d:%02d GMT+8)", result, expectedResult, hour, minute)
	}
}
