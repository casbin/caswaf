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

package main

import (
	"github.com/beego/beego"
	"github.com/beego/beego/plugins/cors"
	_ "github.com/beego/beego/session/redis"
	"github.com/casbin/caswaf/casdoor"
	"github.com/casbin/caswaf/object"
	"github.com/casbin/caswaf/proxy"
	"github.com/casbin/caswaf/routers"
	"github.com/casbin/caswaf/run"
	"github.com/casbin/caswaf/service"
)

func main() {
	object.InitFlag()
	object.InitAdapter()
	object.CreateTables()
	casdoor.InitCasdoorConfig()
	proxy.InitHttpClient()
	object.InitSiteMap()
	object.InitRuleMap()
	run.InitAppMap()
	run.InitSelfStart()
	object.StartMonitorSitesLoop()

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	//beego.DelStaticPath("/static")
	beego.SetStaticPath("/static", "web/build/static")
	// https://studygolang.com/articles/2303
	beego.InsertFilter("/", beego.BeforeRouter, routers.TransparentStatic) // must has this for default page
	beego.InsertFilter("/*", beego.BeforeRouter, routers.TransparentStatic)
	beego.InsertFilter("/api/*", beego.BeforeRouter, routers.ApiFilter)

	if beego.AppConfig.String("redisEndpoint") == "" {
		beego.BConfig.WebConfig.Session.SessionProvider = "file"
		beego.BConfig.WebConfig.Session.SessionProviderConfig = "./tmp"
	} else {
		beego.BConfig.WebConfig.Session.SessionProvider = "redis"
		beego.BConfig.WebConfig.Session.SessionProviderConfig = beego.AppConfig.String("redisEndpoint")
	}
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600 * 24 * 365

	service.Start()

	beego.Run()
}
