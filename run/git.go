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
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/beego/beego"
	"xorm.io/xorm"
)

func gitClone(repoUrl string, path string) error {
	fmt.Printf("gitClone(): [%s]\n", path)

	cmd := exec.Command("git", "clone", repoUrl, path)
	return cmd.Run()
}

func GitDiff(path string) (string, error) {
	cmd := exec.Command("git", "diff")
	cmd.Dir = path

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	return out.String(), nil
}

func gitApply(path string, patch string) error {
	fmt.Printf("gitApply(): [%s]\n", path)

	tmpFile, err := ioutil.TempFile("", "patch")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(patch)
	if err != nil {
		return err
	}

	err = tmpFile.Close()
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "apply", tmpFile.Name())
	cmd.Dir = path
	return cmd.Run()
}

func gitGetLatestCommitHash(path string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func gitPull(path string) (bool, error) {
	oldHash, err := gitGetLatestCommitHash(path)
	if err != nil {
		return false, err
	}

	cmd := exec.Command("git", "pull", "--rebase", "--autostash")
	cmd.Dir = path
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}

	newHash, err := gitGetLatestCommitHash(path)
	if err != nil {
		return false, err
	}

	affected := oldHash != newHash

	if affected {
		fmt.Printf("gitPull(): [%s]\n", path)
		fmt.Printf("Output: %s\n", string(out))
		fmt.Printf("Affected: [%s] -> [%s]\n", oldHash, newHash)
	}

	return affected, nil
}

func runCmd(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	return cmd.Run()
}

func gitWebBuild(path string) error {
	webDir := filepath.Join(path, "web")
	fmt.Printf("gitWebBuild(): [%s]\n", webDir)

	err := runCmd(webDir, "yarn", "install")
	if err != nil {
		return err
	}

	return runCmd(webDir, "yarn", "build")
}

func gitCreateDatabase(name string) (bool, error) {
	fmt.Printf("gitCreateDatabase(): [%s]\n", name)

	driverName := "mysql"
	dataSourceName := fmt.Sprintf("root:%s@tcp(localhost:3306)/", beego.AppConfig.String("dbPass"))
	engine, err := xorm.NewEngine(driverName, dataSourceName)
	if err != nil {
		return false, err
	}

	cmd := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;", name)
	result, err := engine.Exec(cmd)
	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	engine.Close()

	return affected != 0, err
}
