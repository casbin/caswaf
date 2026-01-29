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

package service

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/beego/beego"
	"github.com/casbin/caswaf/conf"
	"github.com/casbin/caswaf/object"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

func joinPath(a string, b string) string {
	if strings.HasSuffix(a, "/") && strings.HasPrefix(b, "/") {
		b = b[1:]
	} else if !strings.HasSuffix(a, "/") && !strings.HasPrefix(b, "/") {
		b = "/" + b
	}
	res := a + b
	return res
}

func isHostIp(host string) bool {
	hostWithoutPort := strings.Split(host, ":")[0]
	ip := net.ParseIP(hostWithoutPort)
	return ip != nil
}

func responseOk(w http.ResponseWriter, format string, a ...interface{}) {
	w.WriteHeader(http.StatusOK)

	msg := fmt.Sprintf(format, a...)
	fmt.Println(msg)
	_, err := fmt.Fprintf(w, msg)
	if err != nil {
		panic(err)
	}
}

func responseError(w http.ResponseWriter, format string, a ...interface{}) {
	w.WriteHeader(http.StatusInternalServerError)

	msg := fmt.Sprintf(format, a...)
	fmt.Println(msg)
	_, err := fmt.Fprintf(w, msg)
	if err != nil {
		panic(err)
	}
}

func responseErrorWithoutCode(w http.ResponseWriter, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println(msg)
	_, err := fmt.Fprintf(w, msg)
	if err != nil {
		panic(err)
	}
}

func getDomainWithoutPort(domain string) string {
	if !strings.Contains(domain, ":") {
		return domain
	}

	tokens := strings.SplitN(domain, ":", 2)
	if len(tokens) > 1 {
		return tokens[0]
	}
	return domain
}

func getSiteByDomainWithWww(domain string) *object.Site {
	hostNonWww := getHostNonWww(domain)
	if hostNonWww != "" {
		domain = hostNonWww
	}

	domainWithoutPort := getDomainWithoutPort(domain)

	site := object.GetSiteByDomain(domainWithoutPort)
	return site
}

func getX509CertByDomain(domain string) (*tls.Certificate, error) {
	cert, err := object.GetCertByDomain(domain)
	if err != nil {
		return nil, fmt.Errorf("getX509CertByDomain() error: %v, domain: [%s]", err, domain)
	}
	if cert == nil {
		return nil, fmt.Errorf("getX509CertByDomain() error: cert not found for domain: [%s]", domain)
	}

	tlsCert, certErr := tls.X509KeyPair([]byte(cert.Certificate), []byte(cert.PrivateKey))

	return &tlsCert, certErr
}

func getCasdoorClientFromSite(site *object.Site) (*casdoorsdk.Client, error) {
	if site.ApplicationObj == nil {
		return nil, fmt.Errorf("site.ApplicationObj is empty")
	}

	casdoorEndpoint := beego.AppConfig.String("casdoorEndpoint")
	if casdoorEndpoint == "http://localhost:8000" {
		casdoorEndpoint = "http://localhost:7001"
	}

	clientId := site.ApplicationObj.ClientId
	clientSecret := site.ApplicationObj.ClientSecret

	certificate := ""
	if site.ApplicationObj.CertObj != nil {
		certificate = site.ApplicationObj.CertObj.Certificate
	}

	res := casdoorsdk.NewClient(casdoorEndpoint, clientId, clientSecret, certificate, site.ApplicationObj.Organization, site.CasdoorApplication)
	return res, nil
}

func getScheme(r *http.Request) string {
	scheme := r.URL.Scheme
	if scheme == "" {
		scheme = "http"
	}
	return scheme
}

func getCasdoorEndpoint() string {
	endpoint := conf.GetConfigString("casdoorEndpoint")
	if endpoint == "http://localhost:8000" {
		endpoint = "http://localhost:7001"
	}
	return endpoint
}

// isAllowedOrigin checks if the given origin is in the allowed list
// Returns true if the origin is trusted, false otherwise
func isAllowedOrigin(origin string) bool {
	// Empty origin is not allowed
	if origin == "" {
		return false
	}

	// Hardcoded list of allowed origins for CORS with credentials
	// This should be configured based on your specific deployment
	allowedOrigins := []string{
		"http://localhost:7001",
		"http://localhost:17000",
		"https://localhost:7001",
		"https://localhost:17000",
	}

	for _, allowed := range allowedOrigins {
		if origin == allowed {
			return true
		}
	}

	return false
}

// setSecureCORSHeaders sets secure CORS headers on the response
// Only allows credentials for trusted origins, blocks all others
func setSecureCORSHeaders(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	// Only set CORS headers if there's an Origin header in the request
	if origin == "" {
		return
	}

	// Check if origin is allowed
	if isAllowedOrigin(origin) {
		// Only allow credentials for trusted origins
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length")
	}
	// If origin is not allowed, don't set any CORS headers
	// This prevents credential-bearing cross-origin requests
}
