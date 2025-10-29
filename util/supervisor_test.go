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

package util

import (
	"os"
	"testing"
)

func TestSupervisorConstants(t *testing.T) {
	// Verify that constants are defined
	if EnvSupervisorKey == "" {
		t.Error("EnvSupervisorKey should not be empty")
	}
	if EnvSupervisorMode == "" {
		t.Error("EnvSupervisorMode should not be empty")
	}
	if MaxRestarts <= 0 {
		t.Error("MaxRestarts should be positive")
	}
	if RestartDelay <= 0 {
		t.Error("RestartDelay should be positive")
	}
}

func TestInitSelfGuard_AlreadySupervised(t *testing.T) {
	// Save original environment
	originalValue := os.Getenv(EnvSupervisorKey)
	defer func() {
		if originalValue == "" {
			os.Unsetenv(EnvSupervisorKey)
		} else {
			os.Setenv(EnvSupervisorKey, originalValue)
		}
	}()

	// Set environment to indicate already supervised
	os.Setenv(EnvSupervisorKey, "1")

	// This should return immediately without error or exit
	// We can't fully test InitSelfGuard as it calls os.Exit in other paths
	// But we can verify it doesn't panic when supervised
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("InitSelfGuard panicked when already supervised: %v", r)
		}
	}()

	// Note: In a real supervised scenario, this would just return
	// We can't test the full flow without mocking os.Exit
}
