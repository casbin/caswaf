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
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/beego/beego"
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

func getNameIndex(name string) int {
	tokens := strings.Split(name, "_")
	if len(tokens) > 0 {
		return util.ParseInt(tokens[len(tokens)-1])
	} else {
		panic(fmt.Sprintf("getNameIndex() error, name = %s", name))
	}
}

func getRepoUrl(name string) string {
	if name == "casdoor" {
		return "https://github.com/casdoor/casdoor"
	} else {
		return fmt.Sprintf("https://github.com/casbin/%s", name)
	}
}

func getShortcut() string {
	res := "Shortcut"
	language := beego.AppConfig.String("language")
	if language != "en" {
		res = "快捷方式"
	}
	return res
}

func ensureFileFolderExists(path string) error {
	if !util.FileExist(path) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

func updateAppConfFile(name string, i int) {
	fmt.Printf("updateAppConfFile(): [%s]\n", name)

	confPath := getCodeAppConfPath(name)
	content := util.ReadStringFromPath(confPath)
	content = strings.ReplaceAll(content, "httpport = 8000", fmt.Sprintf("httpport = %d", 30000+i))
	content = strings.ReplaceAll(content, "123456", beego.AppConfig.String("dbPass"))
	content = strings.ReplaceAll(content, "dbName = casdoor", fmt.Sprintf("dbName = casdoor_customer_%d", i))
	content = strings.ReplaceAll(content, "redisEndpoint =", "redisEndpoint = \"localhost:6379\"")
	content = strings.ReplaceAll(content, "socks5Proxy = \"127.0.0.1:10808\"", "socks5Proxy =")
	util.WriteStringToPath(content, confPath)
}

func updateBatFile(name string) error {
	fmt.Printf("updateBatFile(): [%s]\n", name)

	batPath := getBatPath(name)
	err := ensureFileFolderExists(filepath.Dir(batPath))
	if err != nil {
		return err
	}
	content := fmt.Sprintf("cd %s\ngo run main.go", GetRepoPath(name))
	util.WriteStringToPath(content, batPath)
	return nil
}

func updateShortcutFile(name string) error {
	fmt.Printf("updateShortcutFile(): [%s]\n", name)

	cmd := exec.Command("powershell", fmt.Sprintf("$s=(New-Object -COM WScript.Shell).CreateShortcut('%s');$s.TargetPath='%s';$s.Save()", getShortcutPath(name), getBatPath(name)))
	return cmd.Run()
}

func startProcess(name string) error {
	fmt.Printf("startProcess(): [%s]\n", name)

	cmd := exec.Command("cmd", "/C", "start", "", getShortcutPath(name))
	return cmd.Run()
}

func stopProcess(name string) error {
	fmt.Printf("stopProcess(): [%s]\n", name)

	windowName := fmt.Sprintf("%s.bat - %s", name, getShortcut())
	// taskkill /IM "casdoor.bat - Shortcut" /F
	// taskkill /F /FI "WINDOWTITLE eq casdoor.bat - Shortcut" /T
	cmd := exec.Command("taskkill", "/F", "/FI", fmt.Sprintf("WINDOWTITLE eq %s", windowName), "/T")
	return cmd.Run()
}

func IsProcessActive(pid int) (bool, error) {
	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid))
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, err
	}

	output := out.String()
	res := strings.Contains(output, strconv.Itoa(pid))
	return res, nil
}
