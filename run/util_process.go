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
	"os/exec"
	"regexp"
	"strings"

	"github.com/casbin/caswaf/util"
)

var reBatNames *regexp.Regexp

func init() {
	reBatNames = regexp.MustCompile(`\\Desktop\\run\\(.*?)\.bat`)
}

func parseBatName(s string) string {
	res := reBatNames.FindStringSubmatch(s)
	if res == nil {
		return ""
	}

	return res[1]
}

func getBatNamesFromOutput(output string) map[string]int {
	batNameMap := map[string]int{}

	output = strings.ReplaceAll(output, "\r", "")
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		tokens := strings.Split(line, " ")
		tokens2 := []string{}
		for _, token := range tokens {
			if token != "" {
				tokens2 = append(tokens2, token)
			}
		}

		if len(tokens2) < 5 || strings.ToLower(tokens2[0]) != `c:\windows\system32\cmd.exe` || tokens2[1] != "/c" {
			continue
		}

		batName := parseBatName(tokens2[2])
		processId := util.ParseInt(tokens2[len(tokens2)-1])
		batNameMap[batName] = processId
		//fmt.Printf("%s, %d\n", batName, processId)
	}

	return batNameMap
}

func getPid(name string) int {
	// wmic process where (name="cmd.exe") get CommandLine, ProcessID
	cmd := exec.Command("wmic", "process", "where", "name='cmd.exe'", "get", "CommandLine,ProcessID")
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	batNameMap := getBatNamesFromOutput(string(out))
	pid, ok := batNameMap[name]
	if ok {
		return pid
	} else {
		panic(fmt.Sprintf("getBatNamesFromOutput() error, name = %s, batNameMap = %v", name, batNameMap))
	}
}
