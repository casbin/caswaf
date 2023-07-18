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
