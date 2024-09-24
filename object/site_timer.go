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
	"sync"
	"time"

	"github.com/casbin/caswaf/util"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

var (
	siteUpdateMap          = map[string]string{}
	lock                   = &sync.Mutex{}
	healthCheckTryTimesMap = map[string]int{}
)

func monitorSiteNodes() error {
	sites, err := GetGlobalSites()
	if err != nil {
		return err
	}

	for _, site := range sites {
		//updatedTime, ok := siteUpdateMap[site.GetId()]
		//if ok && updatedTime != "" && updatedTime == site.UpdatedTime {
		//	continue
		//}

		lock.Lock()
		err = site.checkNodes()
		lock.Unlock()
		if err != nil {
			return err
		}

		siteUpdateMap[site.GetId()] = site.UpdatedTime
	}

	return err
}

func monitorSiteCerts() error {
	sites, err := GetGlobalSites()
	if err != nil {
		return err
	}

	for _, site := range sites {
		//updatedTime, ok := siteUpdateMap[site.GetId()]
		//if ok && updatedTime != "" && updatedTime == site.UpdatedTime {
		//	continue
		//}

		lock.Lock()
		err = site.checkCerts()
		lock.Unlock()
		if err != nil {
			return err
		}

		siteUpdateMap[site.GetId()] = site.UpdatedTime
	}

	return err
}

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

func StartMonitorSitesLoop() {
	fmt.Printf("StartMonitorSitesLoop() Start!\n\n")
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[%s] Recovered from StartMonitorSitesLoop() panic: %v\n", util.GetCurrentTime(), r)
				StartMonitorSitesLoop()
			}
		}()

		for {
			err := refreshSiteMap()
			if err != nil {
				fmt.Println(err)
				continue
			}

			err = monitorSiteNodes()
			if err != nil {
				fmt.Println(err)
				continue
			}

			err = monitorSiteCerts()
			if err != nil {
				fmt.Println(err)
				continue
			}

			startHealthCheckLoop()

			time.Sleep(5 * time.Second)
		}
	}()
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
