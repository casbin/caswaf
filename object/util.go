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
	"net"
	"strings"

	"golang.org/x/net/publicsuffix"
)

func resolveDomainToIp(domain string) string {
	ips, err := net.LookupIP(domain)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return "(empty)"
		}

		fmt.Printf("resolveDomainToIp() error: %s\n", err.Error())
		return err.Error()
	}

	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			return ipv4.String()
		}
	}
	return "(empty)"
}

func getBaseDomain(domain string) string {
	// abc.com -> abc.com
	// abc.com.it -> abc.com.it
	// subdomain.abc.io -> abc.io
	// subdomain.abc.org.us -> abc.org.us
	baseDomain, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return baseDomain
}
