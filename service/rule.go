package service

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/casbin/caswaf/object"
)

type UaRule struct {
	Rule  object.Rule
	check CheckRule
}

type CheckRule interface {
	checkRule(expressions []*object.Expression, req *http.Request) (string, string, error)
}

func (r *UaRule) checkRule(expressions []*object.Expression, req *http.Request) (string, string, error) {
	userAgent := req.UserAgent()
	for _, expression := range expressions {
		ua := expression.Value
		switch expression.Operator {
		case "contains":
			if strings.Contains(userAgent, ua) {
				return r.Rule.Action, r.Rule.Reason, nil
			}
		case "does not contain":
			if !strings.Contains(userAgent, ua) {
				return r.Rule.Action, r.Rule.Reason, nil
			}
		case "equals":
			if userAgent == ua {
				return r.Rule.Action, r.Rule.Reason, nil
			}
		case "does not equal":
			if strings.Compare(userAgent, ua) != 0 {
				return r.Rule.Action, r.Rule.Reason, nil
			}
		case "match":
			// regex match
			isMatched, err := regexp.MatchString(ua, userAgent)
			if err != nil {
				return "", "", err
			}
			if isMatched {
				return r.Rule.Action, r.Rule.Reason, nil
			}
		}
	}
	return "", "", nil
}
