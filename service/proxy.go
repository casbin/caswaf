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
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/casbin/caswaf/object"
)

func forwardHandler(targetURL string, writer http.ResponseWriter, request *http.Request) {
	u, err := url.Parse(targetURL)
	if nil != err {
		log.Println(err)
		return
	}

	proxy := httputil.ReverseProxy{
		Director: func(request *http.Request) {
			request.URL = u
		},
	}

	proxy.ServeHTTP(writer, request)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	site := object.GetSiteByDomain(r.Host)
	if site == nil {
		// cache miss
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if site.SslMode == "HTTPS Only" {
		// This domain only supports https but receive http request, redirect to https
		if r.TLS == nil {
			safetyURL := fmt.Sprintf("https://%s%s", r.Host, r.URL.Path)

			http.Redirect(w, r, safetyURL, http.StatusMovedPermanently)
		}
	}

	forwardHandler(site.Host+r.URL.Path, w, r)
}

func getCertificateForDomain(domain string) (*tls.Certificate, error) {
	site := object.GetSiteByDomain(domain)
	tlsCert, certErr := tls.X509KeyPair([]byte(site.SslCertObj.Certificate), []byte(site.SslCertObj.PrivateKey))

	return &tlsCert, certErr
}

func Start() {
	http.HandleFunc("/", handleRequest)

	go func() {
		err := http.ListenAndServe(":80", nil)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		server := &http.Server{
			Addr:      ":443",
			TLSConfig: &tls.Config{},
		}

		// start https server and set how to get certificate
		server.TLSConfig.GetCertificate = func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			domain := info.ServerName
			cert, err := getCertificateForDomain(domain)

			if err != nil {
				return nil, err
			}

			return cert, nil
		}

		err := server.ListenAndServeTLS("", "")
		if err != nil {
			panic(err)
		}
	}()
}
