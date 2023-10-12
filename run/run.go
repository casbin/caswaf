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
	"path/filepath"
	"strings"

	"github.com/beego/beego"
	"github.com/casbin/caswaf/util"
)

func isTargetRepo(siteName string) bool {
	return strings.HasPrefix(siteName, "cc_") || strings.Count(siteName, "_") == 2
}

func wrapRepoError(function string, err error) (int, error) {
	return 0, fmt.Errorf("%s(): %s", function, err.Error())
}

func CreateRepo(siteName string, needStart bool, diff string, providerName string) (int, error) {
	path := GetRepoPath(siteName)
	if !util.FileExist(path) {
		originalName := getOriginalName(siteName)
		repoUrl := getRepoUrl(originalName)
		err := gitClone(repoUrl, path)
		if err != nil {
			return wrapRepoError("gitClone", err)
		}

		language := beego.AppConfig.String("language")
		if language == "en" {
			_, err = gitCreateDatabase(siteName)
			if err != nil {
				return wrapRepoError("gitCreateDatabase", err)
			}
		}

		needWebBuild := false
		if isTargetRepo(siteName) {
			index := getNameIndex(siteName)
			updateAppConfFile(siteName, index)
			if index == 0 {
				needWebBuild = true
			}
		} else {
			needWebBuild = true

			if diff != "" {
				err = gitApply(path, diff)
				if err != nil {
					return wrapRepoError("gitApply", err)
				}
			}
		}

		if needWebBuild {
			err = gitWebBuild(path)
			if err != nil {
				return wrapRepoError("gitWebBuild", err)
			}

			err = gitUploadCdn(providerName, siteName)
			if err != nil {
				return wrapRepoError("gitUploadCdn", err)
			}
		}

		batExisted, err := updateBatFile(siteName)
		if err != nil {
			return wrapRepoError("updateBatFile", err)
		}

		if !batExisted {
			err = updateShortcutFile(siteName)
			if err != nil {
				return wrapRepoError("updateShortcutFile", err)
			}
		}
	} else {
		affected, err := gitPull(path)
		if err != nil {
			return wrapRepoError("gitPull", err)
		}

		needWebBuild := false
		if affected {
			if isTargetRepo(siteName) {
				index := getNameIndex(siteName)
				if index == 0 {
					needWebBuild = true
				}
			} else {
				needWebBuild = true
			}
		} else {
			webIndex := filepath.Join(path, "web/build/index.html")
			if !util.FileExist(webIndex) {
				needWebBuild = true
			}

			if isTargetRepo(siteName) {
				index := getNameIndex(siteName)
				if index != 0 {
					needWebBuild = false
				}
			}
		}

		if needWebBuild {
			err = gitWebBuild(path)
			if err != nil {
				return wrapRepoError("gitWebBuild", err)
			}

			err = gitUploadCdn(providerName, siteName)
			if err != nil {
				return wrapRepoError("gitUploadCdn", err)
			}
		}

		batExisted, err := updateBatFile(siteName)
		if err != nil {
			return wrapRepoError("updateBatFile", err)
		}

		if !batExisted {
			err = updateShortcutFile(siteName)
			if err != nil {
				return wrapRepoError("updateShortcutFile", err)
			}
		}

		if affected {
			if !needStart {
				err = stopProcess(siteName)
				if err != nil {
					return wrapRepoError("stopProcess", err)
				}

				err = startProcess(siteName)
				if err != nil {
					return wrapRepoError("startProcess", err)
				}

				var pid int
				pid, err = getPid(siteName)
				if err != nil {
					return wrapRepoError("getPid", err)
				}

				return pid, nil
			}
		}
	}

	if needStart {
		err := startProcess(siteName)
		if err != nil {
			return wrapRepoError("startProcess", err)
		}

		pid, err := getPid(siteName)
		if err != nil {
			return wrapRepoError("getPid", err)
		}

		return pid, nil
	}

	return 0, nil
}
