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
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/casbin/caswaf/util"
)

func getOriginalName(name string) string {
	tokens := strings.Split(name, "_")
	if len(tokens) > 0 {
		return tokens[0]
	} else {
		return name
	}
}

func getRepoUrl(name string) string {
	if name == "casdoor" {
		return "https://github.com/casdoor/casdoor"
	} else {
		return fmt.Sprintf("https://github.com/casbin/%s", name)
	}
}

func ensureFileFolderExists(path string) {
	if !util.FileExist(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

func updateBatFile(name string) {
	fmt.Printf("Updating BAT file: [%s]\n", name)

	batPath := getBatPath(name)
	ensureFileFolderExists(filepath.Dir(batPath))
	content := fmt.Sprintf("cd %s\ngo run main.go", GetRepoPath(name))
	util.WriteStringToPath(content, batPath)
}

func updateShortcutFile(name string) {
	fmt.Printf("Updating shortcut file: [%s]\n", name)

	cmd := exec.Command("powershell", fmt.Sprintf("$s=(New-Object -COM WScript.Shell).CreateShortcut('%s');$s.TargetPath='%s';$s.Save()", getShortcutPath(name), getBatPath(name)))
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func startProcess(name string) {
	fmt.Printf("Starting process: [%s]\n", name)

	cmd := exec.Command("cmd", "/C", "start", "", getShortcutPath(name))
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
