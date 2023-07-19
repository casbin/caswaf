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
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/casbin/caswaf/object"
)

func ForwardHandler(targetURL string, writer http.ResponseWriter, request *http.Request) {
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
	fmt.Println("port :80 receive request:", r.URL.Path)

	site := object.GetSiteByDomain(r.Host)
	if site == nil {
		// cache miss
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//fmt.Println("site.Domain:", site.Domain)
	//fmt.Println("site.Host:", site.Host)

	// forward to target
	ForwardHandler(site.Host+r.URL.Path, w, r)
}

func Start() {
	fmt.Println("listening port 80")

	go func() {
		http.HandleFunc("/", handleRequest)

		err := http.ListenAndServe(":80", nil)
		if err != nil {
			panic(err)
		}
	}()
}
