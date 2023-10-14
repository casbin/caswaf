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
	"path/filepath"
	"strings"

	"github.com/beego/beego"
	"github.com/casbin/caswaf/object"
	"github.com/casbin/caswaf/util"
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
			if xff := r.Header.Get("X-Forwarded-For"); xff != "" && xff != clientIP {
				newXff := fmt.Sprintf("%s, %s", xff, clientIP)
				//r.Header.Set("X-Forwarded-For", newXff)
				r.Header.Set("X-Real-Ip", newXff)
			} else {
				//r.Header.Set("X-Forwarded-For", clientIP)
				r.Header.Set("X-Real-Ip", clientIP)
			}
		}
	}

	proxy.ServeHTTP(writer, request)
}

func getHostNonWww(host string) string {
	res := ""
	tokens := strings.Split(host, ".")
	if len(tokens) > 2 && tokens[0] == "www" {
		res = strings.Join(tokens[1:], ".")
	}
	return res
}

func redirectToHttps(w http.ResponseWriter, r *http.Request) {
	targetUrl := fmt.Sprintf("https://%s", joinPath(r.Host, r.RequestURI))
	http.Redirect(w, r, targetUrl, http.StatusMovedPermanently)
}

func redirectToHost(w http.ResponseWriter, r *http.Request, host string) {
	protocol := "https"
	if r.TLS == nil {
		protocol = "http"
	}

	targetUrl := fmt.Sprintf("%s://%s", protocol, joinPath(host, r.RequestURI))
	http.Redirect(w, r, targetUrl, http.StatusMovedPermanently)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.UserAgent(), "Uptime-Kuma") {
		fmt.Printf("handleRequest: %s\t%s\t%s\t%s\t%s\n", r.RemoteAddr, r.Method, r.Host, r.RequestURI, r.UserAgent())
	}

	site := getSiteByDomainWithWww(r.Host)
	if site == nil {
		responseError(w, "CasWAF error: site not found for host: %s", r.Host)
		return
	}

	hostNonWww := getHostNonWww(r.Host)
	if hostNonWww != "" {
		redirectToHost(w, r, hostNonWww)
	}

	if site.Domain != r.Host && site.NeedRedirect {
		redirectToHost(w, r, site.Domain)
	}

	if site.Node == "" {
		site.Node = util.GetHostname()
		_, err := object.UpdateSiteNoRefresh(site.GetId(), site)
		responseError(w, "CasWAF error: UpdateSiteNoRefresh() error: %v", err)
		return
	}

	if strings.HasPrefix(r.RequestURI, "/.well-known/acme-challenge/") {
		challengeMap := site.GetChallengeMap()
		for token, keyAuth := range challengeMap {
			if r.RequestURI == fmt.Sprintf("/.well-known/acme-challenge/%s", token) {
				responseOk(w, keyAuth)
				return
			}
		}
	}

	if site.SslMode == "HTTPS Only" {
		// This domain only supports https but receive http request, redirect to https
		if r.TLS == nil {
			redirectToHttps(w, r)
		}
	}

	// oAuth proxy
	if site.CasdoorApplication != "" {
		// handle oAuth proxy
		cookie, err := r.Cookie("casdoor_access_token")
		if err != nil && err.Error() != "http: named cookie not present" {
			panic(err)
		}

		casdoorClient, err := getCasdoorClientFromSite(site)
		if err != nil {
			responseError(w, "CasWAF error: getCasdoorClientFromSite() error: %s", err.Error())
			return
		}

		if cookie == nil {
			// not logged in
			redirectToCasdoor(casdoorClient, w, r)
			return
		} else {
			_, err = casdoorClient.ParseJwtToken(cookie.Value)
			if err != nil {
				responseError(w, "CasWAF error: casdoorClient.ParseJwtToken() error: %s", err.Error())
				return
			}
		}
	}

	host := site.GetHost()
	if host == "" {
		responseError(w, "CasWAF error: targetUrl should not be empty for host: %s, site = %v", r.Host, site)
		return
	}

	if site.SslMode == "Static Folder" {
		var path string
		if r.RequestURI != "/" {
			path = filepath.Join(host, r.RequestURI)
		} else {
			path = filepath.Join(host, "/index.htm")
			if !util.FileExist(path) {
				path = filepath.Join(host, "/index.html")
				if !util.FileExist(path) {
					path = filepath.Join(host, r.RequestURI)
				}
			}
		}
		http.ServeFile(w, r, path)
	} else {
		targetUrl := joinPath(site.GetHost(), r.RequestURI)
		forwardHandler(targetUrl, w, r)
	}
}

func Start() {
	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/caswaf-handler", handleAuthCallback)

	gatewayEnabled, err := beego.AppConfig.Bool("gatewayEnabled")
	if err != nil {
		panic(err)
	}
	if !gatewayEnabled {
		fmt.Printf("CasWAF gateway not enabled (gatewayEnabled == \"false\")\n")
		return
	}

	go func() {
		fmt.Printf("CasWAF gateway running on: http://127.0.0.1:80\n")
		err := http.ListenAndServe(":80", nil)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		fmt.Printf("CasWAF gateway running on: https://127.0.0.1:443\n")
		server := &http.Server{
			Addr:      ":443",
			TLSConfig: &tls.Config{},
		}

		// start https server and set how to get certificate
		server.TLSConfig.GetCertificate = func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			domain := info.ServerName
			cert, err := getX509CertByDomain(domain)
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
