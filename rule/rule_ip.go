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

package rule

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/casbin/caswaf/object"
	"github.com/casbin/caswaf/util"
)

type IpRule struct{}

func (r *IpRule) checkRule(expressions []*object.Expression, req *http.Request) (bool, string, error) {
	clientIp := util.GetClientIp(req)
	for _, expression := range expressions {
		reason := fmt.Sprintf("expression matched: \"%s %s %s\"", clientIp, expression.Operator, expression.Value)
		ips := strings.Split(expression.Value, " ")
		op := expression.Operator == "is in"
		for _, ip := range ips {
			_, ipNet, err := net.ParseCIDR(ip)
			// use err to determine if ip is a CIDR
			if err == nil && ipNet.Contains(net.ParseIP(clientIp)) == op {
				return true, reason, nil
			}
			if err != nil && (ip == clientIp) == op {
				return true, reason, nil
			}
		}
	}
	return false, "", nil
}
