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
	"strings"

	"github.com/casbin/caswaf/util"
)

func CreateRepo(siteName string, needStart bool, diff string) {
	path := GetRepoPath(siteName)
	if !util.FileExist(path) {
		originalName := getOriginalName(siteName)
		repoUrl := getRepoUrl(originalName)
		gitClone(repoUrl, path)

		if strings.HasPrefix(siteName, "cc_") || strings.Count(siteName, "_") == 2 {
			index := getNameIndex(siteName)
			updateAppConfFile(siteName, index)
		} else if diff != "" {
			gitApply(path, diff)
		}
	} else {
		gitPull(path)
	}

	if needStart {
		updateBatFile(siteName)
		updateShortcutFile(siteName)
		startProcess(siteName)
	}
}
