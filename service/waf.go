// Copyright 2024 The casbin Authors. All Rights Reserved.
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

package service

import (
	"fmt"

	"github.com/casbin/caswaf/conf"
	"github.com/casbin/caswaf/object"
	"github.com/corazawaf/coraza/v3"
	"github.com/corazawaf/coraza/v3/types"
)

var wafs = map[*object.Site]*coraza.WAF{}
var sites []*object.Site

func createWaf(site *object.Site) {
	waf, err := coraza.NewWAF(
		coraza.NewWAFConfig().
			WithErrorCallback(logError).
			WithDirectives(conf.WafConf).
			WithDirectives(object.GetWafRulesByIds(site.Rules)),
	)
	if err != nil {
		fmt.Printf("createWaf(): %s\n", err.Error())
	}
	wafs[site] = &waf
	site.Waf = waf
}

func getWaf(site *object.Site) {
	if wafs[site] == nil {
		createWaf(site)
		sites = append(sites, site)
	} else {
		site.Waf = *wafs[site]
	}
}

func UpdateWafs() {
	for _, site := range sites {
		createWaf(site)
	}
}

func logError(error types.MatchedRule) {
	msg := error.ErrorLog()
	fmt.Printf("[WAFlogError][%s] %s\n", error.Rule().Severity(), msg)
}
