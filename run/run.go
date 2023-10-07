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
	"path/filepath"
	"strings"

	"github.com/casbin/caswaf/util"
)

func CreateRepo(siteName string, needStart bool, diff string) (int, error) {
	path := GetRepoPath(siteName)
	if !util.FileExist(path) {
		originalName := getOriginalName(siteName)
		repoUrl := getRepoUrl(originalName)
		err := gitClone(repoUrl, path)
		if err != nil {
			return 0, err
		}

		_, err = gitCreateDatabase(siteName)
		if err != nil {
			return 0, err
		}

		if strings.HasPrefix(siteName, "cc_") || strings.Count(siteName, "_") == 2 {
			index := getNameIndex(siteName)
			updateAppConfFile(siteName, index)
			if index == 0 {
				err = gitWebBuild(path)
				if err != nil {
					return 0, err
				}
			}
		} else if diff != "" {
			err = gitApply(path, diff)
			if err != nil {
				return 0, err
			}

			err = gitWebBuild(path)
			if err != nil {
				return 0, err
			}
		}

		err = updateBatFile(siteName)
		if err != nil {
			return 0, err
		}

		err = updateShortcutFile(siteName)
		if err != nil {
			return 0, err
		}
	} else {
		affected, err := gitPull(path)
		if err != nil {
			return 0, err
		}
		if affected {
			err = gitWebBuild(path)
			if err != nil {
				return 0, err
			}

			if !needStart {
				err = stopProcess(siteName)
				if err != nil {
					return 0, err
				}

				err = startProcess(siteName)
				if err != nil {
					return 0, err
				}

				var pid int
				pid, err = getPid(siteName)
				if err != nil {
					return 0, err
				}

				return pid, nil
			}
		} else {
			webIndex := filepath.Join(path, "web/build/index.html")
			if !util.FileExist(webIndex) {
				if strings.HasPrefix(siteName, "cc_") || strings.Count(siteName, "_") == 2 {
					index := getNameIndex(siteName)
					if index == 0 {
						err = gitWebBuild(path)
						if err != nil {
							return 0, err
						}
					}
				} else {
					err = gitWebBuild(path)
					if err != nil {
						return 0, err
					}
				}
			}
		}
	}

	if needStart {
		err := startProcess(siteName)
		if err != nil {
			return 0, err
		}

		pid, err := getPid(siteName)
		if err != nil {
			return 0, err
		}

		return pid, nil
	}

	return 0, nil
}
