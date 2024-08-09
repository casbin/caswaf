package rule

import (
	"net/http"
	"testing"

	"github.com/casbin/caswaf/object"
)

func TestIpRateRule_checkRule(t *testing.T) {
	type fields struct {
		ruleName string
	}
	type args struct {
		args []struct {
			expressions []*object.Expression
			req         *http.Request
		}
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []bool
		want1   []string
		want2   []string
		wantErr []bool
	}{
		{
			name: "Test 1",
			fields: fields{
				ruleName: "rule1",
			},
			args: args{
				args: []struct {
					expressions []*object.Expression
					req         *http.Request
				}{
					{
						expressions: []*object.Expression{
							{
								Operator: "1",
								Value:    "1",
							},
						},
						req: &http.Request{
							RemoteAddr: "127.0.0.1",
						},
					},
					{
						expressions: []*object.Expression{
							{
								Operator: "1",
								Value:    "1",
							},
						},
						req: &http.Request{
							RemoteAddr: "127.0.0.1",
						},
					},
				},
			},
			want:    []bool{false, true},
			want1:   []string{"", "Block"},
			want2:   []string{"", "Rate limit exceeded"},
			wantErr: []bool{false, false},
		},
		{
			name: "Test 2",
			fields: fields{
				ruleName: "rule2",
			},
			args: args{
				args: []struct {
					expressions []*object.Expression
					req         *http.Request
				}{
					{
						expressions: []*object.Expression{
							{
								Operator: "1",
								Value:    "1",
							},
						},
						req: &http.Request{
							RemoteAddr: "127.0.0.1",
						},
					},
					{
						expressions: []*object.Expression{
							{
								Operator: "10",
								Value:    "1",
							},
						},
						req: &http.Request{
							RemoteAddr: "127.0.0.1",
						},
					},
				},
			},
			want:    []bool{false, false},
			want1:   []string{"", ""},
			want2:   []string{"", ""},
			wantErr: []bool{false, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &IpRateRule{
				ruleName: tt.fields.ruleName,
			}
			for i, arg := range tt.args.args {
				got, got1, got2, err := r.checkRule(arg.expressions, arg.req)
				if (err != nil) != tt.wantErr[i] {
					t.Errorf("checkRule() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want[i] {
					t.Errorf("checkRule() got = %v, want %v", got, tt.want)
				}
				if got1 != tt.want1[i] {
					t.Errorf("checkRule() got1 = %v, want %v", got1, tt.want1)
				}
				if got2 != tt.want2[i] {
					t.Errorf("checkRule() got2 = %v, want %v", got2, tt.want2)
				}
			}
		})
	}
}
