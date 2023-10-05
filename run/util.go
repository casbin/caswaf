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
)

func GetSitePath(siteName string) string {
	appDir := beego.AppConfig.String("appDir")
	res := filepath.Join(appDir, siteName)
	return res
}

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
