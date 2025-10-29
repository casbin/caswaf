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

package supervisor

import (
	"os"
	"testing"
)

func TestIsSupervisedProcess(t *testing.T) {
	// Test when not supervised
	os.Unsetenv(EnvSupervisorKey)
	if IsSupervisedProcess() {
		t.Error("Expected IsSupervisedProcess() to return false when env var is not set")
	}

	// Test when supervised
	os.Setenv(EnvSupervisorKey, "1")
	if !IsSupervisedProcess() {
		t.Error("Expected IsSupervisedProcess() to return true when env var is set to '1'")
	}

	// Test with invalid value
	os.Setenv(EnvSupervisorKey, "0")
	if IsSupervisedProcess() {
		t.Error("Expected IsSupervisedProcess() to return false when env var is set to '0'")
	}

	// Cleanup
	os.Unsetenv(EnvSupervisorKey)
}

func TestConstants(t *testing.T) {
	if EnvSupervisorKey != "CASWAF_SUPERVISED" {
		t.Errorf("Expected EnvSupervisorKey to be 'CASWAF_SUPERVISED', got '%s'", EnvSupervisorKey)
	}

	if MaxRestarts != 5 {
		t.Errorf("Expected MaxRestarts to be 5, got %d", MaxRestarts)
	}
}
