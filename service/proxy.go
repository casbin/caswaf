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
	"net/http/httputil"
	"net/url"

	"github.com/casbin/caswaf/object"
)

func forwardHandler(targetUrl string, writer http.ResponseWriter, request *http.Request) {
	target, err := url.Parse(targetUrl)

	if nil != err {
		panic(err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = func(r *http.Request) {
		r.URL = target

		if clientIP, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
				newXff := fmt.Sprintf("%s, %s", xff, clientIP)
				r.Header.Set("X-Forwarded-For", newXff)
			} else {
				r.Header.Set("X-Forwarded-For", clientIP)
			}
		}
	}

	proxy.ServeHTTP(writer, request)
}

func redirectToHttps(w http.ResponseWriter, r *http.Request) {
	safetyUrl := fmt.Sprintf("https://%s%s", r.Host, r.RequestURI)

	w.Header().Set("Location", safetyUrl)
	w.WriteHeader(http.StatusMovedPermanently)

	html := `
				<!DOCTYPE html>
				<html>
				<head>
					<title>301 Moved Permanently</title>
				</head>
				<body>
					<center>
						<h1>301 Moved Permanently to </h1>
					</center>
					<hr>
					<center>caswaf</center>
				</body>
				</html>
			`
	_, err := fmt.Fprint(w, html)
	if err != nil {
		return
	}
	return
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
			redirectToHttps(w, r)
		}
	}

	targetUrl := fmt.Sprintf("%s%s", site.Host, r.RequestURI)
	forwardHandler(targetUrl, w, r)
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
