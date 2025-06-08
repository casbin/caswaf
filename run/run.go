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
	return strings.HasPrefix(siteName, "cc_") || strings.HasPrefix(siteName, "casibase_customer_") || strings.Count(siteName, "_") == 2
}

func isFrontendBaseDirEnabledRepo(siteName string) bool {
	return strings.HasPrefix(siteName, "casibase_") && !strings.HasSuffix(siteName, "keli")
}

func wrapRepoError(function string, path string, err error) (int, error) {
	return 0, fmt.Errorf("%s(): path = %s, %s", function, path, err.Error())
}

func CreateRepo(siteName string, needStart bool, diff string, providerName string, orgName string) (int, error) {
	path := GetRepoPath(siteName)
	if !util.FileExist(path) {
		originalName := getOriginalName(siteName)
		repoUrl := getRepoUrl(originalName)
		err := gitClone(repoUrl, path)
		if err != nil {
			return wrapRepoError("gitClone", path, err)
		}

		dbInstanceId := beego.AppConfig.String("dbInstanceId")
		if dbInstanceId == "" {
			_, err = gitCreateDatabase(siteName)
			if err != nil {
				return wrapRepoError("gitCreateDatabase", path, err)
			}
		} else {
			_, err = gitCreateDatabaseCloud(siteName)
			if err != nil {
				return wrapRepoError("gitCreateDatabaseCloud", path, err)
			}
		}

		needWebBuild := false
		if isTargetRepo(siteName) {
			index := getNameIndex(siteName)
			updateAppConfFile(siteName, index, orgName)
			if index == 0 {
				needWebBuild = true
			}
		} else {
			needWebBuild = true

			if diff != "" {
				err = gitApply(path, diff)
				if err != nil {
					return wrapRepoError("gitApply", path, err)
				}
			}
		}

		if needWebBuild && !isFrontendBaseDirEnabledRepo(siteName) {
			err = gitWebBuild(path)
			if err != nil {
				return wrapRepoError("gitWebBuild", path, err)
			}

			err = gitUploadCdn(providerName, siteName)
			if err != nil {
				return wrapRepoError("gitUploadCdn", path, err)
			}
		}

		batExisted, err := updateBatFile(siteName)
		if err != nil {
			return wrapRepoError("updateBatFile", path, err)
		}

		if !batExisted {
			err = updateShortcutFile(siteName)
			if err != nil {
				return wrapRepoError("updateShortcutFile", path, err)
			}
		}
	} else {
		affected, err := gitPull(path)
		if err != nil {
			return wrapRepoError("gitPull", path, err)
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

		if needWebBuild && !isFrontendBaseDirEnabledRepo(siteName) {
			err = gitWebBuild(path)
			if err != nil {
				return wrapRepoError("gitWebBuild", path, err)
			}

			err = gitUploadCdn(providerName, siteName)
			if err != nil {
				return wrapRepoError("gitUploadCdn", path, err)
			}
		}

		batExisted, err := updateBatFile(siteName)
		if err != nil {
			return wrapRepoError("updateBatFile", path, err)
		}

		if !batExisted {
			err = updateShortcutFile(siteName)
			if err != nil {
				return wrapRepoError("updateShortcutFile", path, err)
			}
		}

		if affected {
			if !needStart {
				if !strings.HasPrefix(siteName, "casdoor") && !strings.HasPrefix(siteName, "casibase") {
					err = stopProcess(siteName)
					if err != nil {
						return wrapRepoError("stopProcess", path, err)
					}
				}

				err = startProcess(siteName)
				if err != nil {
					return wrapRepoError("startProcess", path, err)
				}

				var pid int
				pid, err = getPid(siteName)
				if err != nil {
					return wrapRepoError("getPid", path, err)
				}

				return pid, nil
			}
		}
	}

	if needStart {
		err := startProcess(siteName)
		if err != nil {
			return wrapRepoError("startProcess", path, err)
		}

		pid, err := getPid(siteName)
		if err != nil {
			return wrapRepoError("getPid", path, err)
		}

		return pid, nil
	}

	return 0, nil
}
