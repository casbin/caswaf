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
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/beego/beego"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

//go:embed token_jwt_key.pem
var JwtPublicKey string

type JWKSResponse struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	X5C []string `json:"x5c"` // X.509 certificate chain
}

func InitCasdoorConfig() {
	casdoorEndpoint := beego.AppConfig.String("casdoorEndpoint")
	clientId := beego.AppConfig.String("clientId")
	clientSecret := beego.AppConfig.String("clientSecret")
	casdoorOrganization := beego.AppConfig.String("casdoorOrganization")
	casdoorApplication := beego.AppConfig.String("casdoorApplication")

	// Try to fetch certificate from Casdoor's standard JWKS endpoint
	jwtPublicKey := JwtPublicKey // default to embedded certificate

	jwksUrl := fmt.Sprintf("%s/.well-known/jwks", strings.TrimRight(casdoorEndpoint, "/"))
	beego.Info("Fetching JWT certificate from JWKS: ", jwksUrl)

	if cert, err := fetchCertificateFromJWKS(jwksUrl); err == nil && cert != "" {
		jwtPublicKey = cert
		beego.Info("Successfully loaded JWT certificate from JWKS endpoint")
	} else {
		beego.Warn("Failed to fetch certificate from JWKS (", err, "), using embedded certificate")
	}

	casdoorsdk.InitConfig(casdoorEndpoint, clientId, clientSecret, jwtPublicKey, casdoorOrganization, casdoorApplication)
}

func fetchCertificateFromJWKS(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var jwks JWKSResponse
	if err := json.Unmarshal(body, &jwks); err != nil {
		return "", err
	}

	if len(jwks.Keys) == 0 || len(jwks.Keys[0].X5C) == 0 {
		return "", fmt.Errorf("no certificates in JWKS")
	}

	// Extract the first certificate and format as PEM
	certData := jwks.Keys[0].X5C[0]
	pemCert := fmt.Sprintf("-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----\n", certData)

	return pemCert, nil
}
