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

package util

import (
	"testing"
)

func TestStopOldInstance(t *testing.T) {
	// Test with a port that's not in use
	// This should succeed without error
	err := StopOldInstance(54321)
	if err != nil {
		t.Errorf("StopOldInstance failed for unused port: %v", err)
	}
}

func TestGetPidByPort(t *testing.T) {
	// Test with a port that's not in use
	// Should return 0 and no error
	pid, err := getPidByPort(54321)
	if err != nil {
		t.Errorf("getPidByPort failed: %v", err)
	}
	if pid != 0 {
		t.Errorf("Expected pid to be 0 for unused port, got %d", pid)
	}
}
