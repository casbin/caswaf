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

package run

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
)

func gitClone(repoUrl string, path string) {
	cmd := exec.Command("git", "clone", repoUrl, path)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func GitDiff(path string) string {
	cmd := exec.Command("git", "diff")
	cmd.Dir = path

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	return out.String()
}

func gitApply(path string, patch string) {
	tmpFile, err := ioutil.TempFile("", "patch")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err = tmpFile.WriteString(patch); err != nil {
		panic(err)
	}
	if err = tmpFile.Close(); err != nil {
		panic(err)
	}

	cmd := exec.Command("git", "apply", tmpFile.Name())
	cmd.Dir = path
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
