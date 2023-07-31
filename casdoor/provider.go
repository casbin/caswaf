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

package casdoor

import (
	"github.com/astaxie/beego"
	"github.com/casbin/caswaf/util"
)

type Application struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`

	DisplayName         string   `xorm:"varchar(100)" json:"displayName"`
	Logo                string   `xorm:"varchar(200)" json:"logo"`
	HomepageUrl         string   `xorm:"varchar(100)" json:"homepageUrl"`
	Description         string   `xorm:"varchar(100)" json:"description"`
	Organization        string   `xorm:"varchar(100)" json:"organization"`
	Cert                string   `xorm:"varchar(100)" json:"cert"`
	EnablePassword      bool     `json:"enablePassword"`
	EnableSignUp        bool     `json:"enableSignUp"`
	EnableSigninSession bool     `json:"enableSigninSession"`
	EnableAutoSignin    bool     `json:"enableAutoSignin"`
	EnableCodeSignin    bool     `json:"enableCodeSignin"`
	EnableSamlCompress  bool     `json:"enableSamlCompress"`
	EnableWebAuthn      bool     `json:"enableWebAuthn"`
	EnableLinkWithEmail bool     `json:"enableLinkWithEmail"`
	OrgChoiceMode       string   `json:"orgChoiceMode"`
	SamlReplyUrl        string   `xorm:"varchar(100)" json:"samlReplyUrl"`
	GrantTypes          []string `xorm:"varchar(1000)" json:"grantTypes"`
	Tags                []string `xorm:"mediumtext" json:"tags"`

	ClientId             string   `xorm:"varchar(100)" json:"clientId"`
	ClientSecret         string   `xorm:"varchar(100)" json:"clientSecret"`
	RedirectUris         []string `xorm:"varchar(1000)" json:"redirectUris"`
	TokenFormat          string   `xorm:"varchar(100)" json:"tokenFormat"`
	ExpireInHours        int      `json:"expireInHours"`
	RefreshExpireInHours int      `json:"refreshExpireInHours"`
	SignupUrl            string   `xorm:"varchar(200)" json:"signupUrl"`
	SigninUrl            string   `xorm:"varchar(200)" json:"signinUrl"`
	ForgetUrl            string   `xorm:"varchar(200)" json:"forgetUrl"`
	AffiliationUrl       string   `xorm:"varchar(100)" json:"affiliationUrl"`
	TermsOfUse           string   `xorm:"varchar(100)" json:"termsOfUse"`
	SignupHtml           string   `xorm:"mediumtext" json:"signupHtml"`
	SigninHtml           string   `xorm:"mediumtext" json:"signinHtml"`
	FormCss              string   `xorm:"text" json:"formCss"`
	FormCssMobile        string   `xorm:"text" json:"formCssMobile"`
	FormOffset           int      `json:"formOffset"`
	FormSideHtml         string   `xorm:"mediumtext" json:"formSideHtml"`
	FormBackgroundUrl    string   `xorm:"varchar(200)" json:"formBackgroundUrl"`
}

func GetApplications(owner string) ([]*Application, error) {
	applications := []*Application{}
	err := adapter.Engine.Desc("created_time").Find(&applications, &Application{Owner: owner, Organization: beego.AppConfig.String("casdoorOrganization")})
	if err != nil {
		return applications, err
	}

	return applications, nil
}

func GetApplication(id string) (*Application, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	return getApplication(owner, name)
}

func getApplication(owner string, name string) (*Application, error) {
	if owner == "" || name == "" {
		return nil, nil
	}

	application := Application{Owner: owner, Name: name}
	existed, err := adapter.Engine.Get(&application)
	if err != nil {
		return &application, err
	}

	if existed {
		return &application, nil
	} else {
		return nil, nil
	}
}
