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

	"github.com/casbin/caswaf/util"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

var healthCheckTryTimesMap = map[string]int{}

func healthCheck(site *Site, domain string) error {
	var flag bool
	var log string
	switch site.SslMode {
	case "HTTPS Only":
		flag, log = pingUrl("https://" + domain)
	case "HTTP":
		flag, log = pingUrl("http://" + domain)
	case "HTTPS and HTTP":
		flag, log = pingUrl("https://" + domain)
		flagHttp, logHttp := pingUrl("http://" + domain)
		flag = flag || flagHttp
		log = log + logHttp
	}
	if !flag {
		fmt.Println(log)
		healthCheckTryTimesMap[domain]--
		if healthCheckTryTimesMap[domain] != 0 {
			return nil
		}
		log = fmt.Sprintf("CasWAF health check failed for domain %s, %s", domain, log)
		user, err := casdoorsdk.GetUser(site.Owner)
		if err != nil {
			fmt.Println(err)
		}
		err = casdoorsdk.SendEmail("CasWAF HealthCheck Check Alert", log, "CasWAF", user.Email)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		healthCheckTryTimesMap[domain] = GetSiteByDomain(domain).AlertTryTimes
	}
	return nil
}

func startHealthCheckLoop() {
	for _, domain := range healthCheckNeededDomains {
		domain := domain
		if _, ok := healthCheckTryTimesMap[domain]; ok {
			continue
		}
		healthCheckTryTimesMap[domain] = GetSiteByDomain(domain).AlertTryTimes
		go func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("[%s] Recovered from healthCheck() panic: %v\n", util.GetCurrentTime(), r)
				}
			}()
			for {
				site := GetSiteByDomain(domain)
				if site == nil || !site.EnableAlert {
					return
				}

				err := healthCheck(site, domain)
				if err != nil {
					fmt.Println(err)
				}
				time.Sleep(time.Duration(site.AlertInterval) * time.Second)
			}
		}()
	}
}
