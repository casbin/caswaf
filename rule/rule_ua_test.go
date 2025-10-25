package rule

import (
	"net/http"
	"testing"

	"github.com/casbin/caswaf/object"
)

func TestUaRule_checkRule(t *testing.T) {
	tests := []struct {
		name        string
		expressions []*object.Expression
		userAgent   string
		wantHit     bool
		wantAction  string
		wantReason  string
		wantErr     bool
	}{
		{
			name: "equals operator - match",
			expressions: []*object.Expression{
				{
					Name:     "Current User-Agent",
					Operator: "equals",
					Value:    "test-agent",
				},
			},
			userAgent:  "test-agent",
			wantHit:    true,
			wantAction: "",
			wantReason: "expression matched: \"test-agent equals test-agent\"",
			wantErr:    false,
		},
		{
			name: "equals operator - no match",
			expressions: []*object.Expression{
				{
					Name:     "Current User-Agent",
					Operator: "equals",
					Value:    "aaa",
				},
			},
			userAgent:  "Mozilla/5.0",
			wantHit:    false,
			wantAction: "",
			wantReason: "",
			wantErr:    false,
		},
		{
			name: "contains operator - match",
			expressions: []*object.Expression{
				{
					Name:     "Current User-Agent",
					Operator: "contains",
					Value:    "Mozilla",
				},
			},
			userAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			wantHit:    true,
			wantAction: "",
			wantReason: "expression matched: \"Mozilla/5.0 (Windows NT 10.0; Win64; x64) contains Mozilla\"",
			wantErr:    false,
		},
		{
			name: "contains operator - no match",
			expressions: []*object.Expression{
				{
					Name:     "Current User-Agent",
					Operator: "contains",
					Value:    "Chrome",
				},
			},
			userAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Safari/537.36",
			wantHit:    false,
			wantAction: "",
			wantReason: "",
			wantErr:    false,
		},
		{
			name: "does not equal operator - match",
			expressions: []*object.Expression{
				{
					Name:     "Current User-Agent",
					Operator: "does not equal",
					Value:    "bad-agent",
				},
			},
			userAgent:  "good-agent",
			wantHit:    true,
			wantAction: "",
			wantReason: "expression matched: \"good-agent does not equal bad-agent\"",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UaRule{}
			req := &http.Request{
				Header: http.Header{
					"User-Agent": []string{tt.userAgent},
				},
			}

			gotHit, gotAction, gotReason, err := r.checkRule(tt.expressions, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UaRule.checkRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHit != tt.wantHit {
				t.Errorf("UaRule.checkRule() gotHit = %v, want %v", gotHit, tt.wantHit)
			}
			if gotAction != tt.wantAction {
				t.Errorf("UaRule.checkRule() gotAction = %v, want %v", gotAction, tt.wantAction)
			}
			if gotReason != tt.wantReason {
				t.Errorf("UaRule.checkRule() gotReason = %v, want %v", gotReason, tt.wantReason)
			}
		})
	}
}
