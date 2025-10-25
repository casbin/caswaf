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

func (r *IpRule) checkRule(expressions []*object.Expression, req *http.Request) (*RuleResult, error) {
	clientIp := util.GetClientIp(req)
	netIp, err := parseIp(clientIp)
	if err != nil {
		return nil, err
	}
	for _, expression := range expressions {
		reason := fmt.Sprintf("expression matched: \"%s %s %s\"", clientIp, expression.Operator, expression.Value)
		ips := strings.Split(expression.Value, ",")
		for _, ip := range ips {
			if strings.Contains(ip, "/") {
				_, ipNet, err := net.ParseCIDR(ip)
				if err != nil {
					return nil, err
				}

				switch expression.Operator {
				case "is in":
					if ipNet.Contains(netIp) {
						return &RuleResult{Reason: reason}, nil
					}
				case "is not in":
					if !ipNet.Contains(netIp) {
						return &RuleResult{Reason: reason}, nil
					}
				default:
					return nil, fmt.Errorf("unknown operator: %s", expression.Operator)
				}
			} else if strings.ContainsAny(ip, ".:") {
				switch expression.Operator {
				case "is in":
					if ip == clientIp {
						return &RuleResult{Reason: reason}, nil
					}
				case "is not in":
					if ip != clientIp {
						return &RuleResult{Reason: reason}, nil
					}
				default:
					return nil, fmt.Errorf("unknown operator: %s", expression.Operator)
				}
			} else {
				return nil, fmt.Errorf("unknown IP or CIDR format: %s", ip)
			}
		}
	}
	return nil, nil
}

func parseIp(ipStr string) (net.IP, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("unknown IP or CIDR format: %s", ipStr)
	}
	return ip, nil
}
