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
	"crypto"
	"fmt"
	"github.com/beego/beego"

	"github.com/casbin/caswaf/proxy"
	"github.com/casbin/lego/v4/acme"
	"github.com/casbin/lego/v4/certcrypto"
	"github.com/casbin/lego/v4/lego"
	"github.com/casbin/lego/v4/registration"
)

type Account struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

// GetEmail returns the email address for the account.
func (a *Account) GetEmail() string {
	return a.Email
}

// GetPrivateKey returns the private RSA account key.
func (a *Account) GetPrivateKey() crypto.PrivateKey {
	return a.key
}

// GetRegistration returns the server registration.
func (a *Account) GetRegistration() *registration.Resource {
	return a.Registration
}

func getLegoClientAndAccount(email string, privateKey string, devMode bool) (*lego.Client, *Account) {
	account := &Account{
		Email: email,
		key:   decodeEccKey(privateKey),
	}

	config := lego.NewConfig(account)
	if devMode {
		config.CADirURL = lego.LEDirectoryStaging
	} else {
		config.CADirURL = lego.LEDirectoryProduction
	}

	config.Certificate.KeyType = certcrypto.RSA2048
	config.HTTPClient = proxy.ProxyHttpClient

	client, err := lego.NewClient(config)
	if err != nil {
		panic(err)
	}

	return client, account
}

func getAcmeClient(email string, privateKey string, devMode bool) *lego.Client {
	// Create a user. New accounts need an email and private key to start.
	client, account := getLegoClientAndAccount(email, privateKey, devMode)

	// try to obtain an account based on the private key
	var err error
	account.Registration, err = client.Registration.ResolveAccountByKey()
	if err != nil {
		acmeError, ok := err.(*acme.ProblemDetails)
		if !ok {
			panic(err)
		}

		if acmeError.Type != "urn:ietf:params:acme:error:accountDoesNotExist" {
			panic(acmeError)
		}

		// Failed to get account, so create an account based on the private key.
		var err error
		account.Registration, err = client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
		if err != nil {
			panic(err)
		}
	}

	return client
}

func GetAcmeClient() (*lego.Client, error) {
	acmeEmail := beego.AppConfig.String("acmeEmail")
	acmePrivateKey := beego.AppConfig.String("acmePrivateKey")
	if acmeEmail == "" {
		return nil, fmt.Errorf("acmeEmail should not be empty")
	}
	if acmePrivateKey == "" {
		return nil, fmt.Errorf("acmePrivateKey should not be empty")
	}

	client := getAcmeClient(acmeEmail, acmePrivateKey, false)
	return client, nil
}
