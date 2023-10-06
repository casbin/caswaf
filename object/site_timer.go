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

package object

import (
	"fmt"
	"time"
)

var siteUpdateMap = map[string]string{}

func monitorSites() {
	sites := GetGlobalSites()
	for _, site := range sites {
		updatedTime, ok := siteUpdateMap[site.GetId()]
		if ok && updatedTime != "" && updatedTime == site.UpdatedTime {
			continue
		}

		site.checkNodes()
		siteUpdateMap[site.GetId()] = site.UpdatedTime
	}
}

func StartMonitorSitesLoop() {
	fmt.Printf("StartMonitorSitesLoop() Start!\n\n")
	go func() {
		for {
			refreshSiteMap()
			monitorSites()
			time.Sleep(5 * time.Second)
		}
	}()

}
