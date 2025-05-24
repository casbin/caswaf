// Copyright 2025 The casbin Authors. All Rights Reserved.
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

package run

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/beego/beego"
)

var rdsClient *rds.Client

func InitRdsClient() {
	dbRegionId := beego.AppConfig.String("dbRegionId")
	dbAccessKeyId := beego.AppConfig.String("dbAccessKeyId")
	dbAccessKeySecret := beego.AppConfig.String("dbAccessKeySecret")

	if dbRegionId == "" || dbAccessKeyId == "" || dbAccessKeySecret == "" {
		return
	}

	config := sdk.NewConfig()
	credential := credentials.NewAccessKeyCredential(dbAccessKeyId, dbAccessKeySecret)

	var err error
	rdsClient, err = rds.NewClientWithOptions(dbRegionId, config, credential)
	if err != nil {
		panic(err)
	}
}

func gitCreateDatabaseCloud(name string) (bool, error) {
	fmt.Printf("gitCreateDatabaseCloud(): [%s]\n", name)

	dbInstanceId := beego.AppConfig.String("dbInstanceId")

	// https://help.aliyun.com/document_detail/26258.htm
	r := rds.CreateCreateDatabaseRequest()
	r.DBInstanceId = dbInstanceId
	r.DBName = name
	r.CharacterSetName = "utf8mb4"

	_, err := rdsClient.CreateDatabase(r)
	if err != nil {
		return false, err
	}

	err = addDatabaseUser(name)
	if err != nil {
		return false, err
	}

	return true, err
}

func addDatabaseUser(dbName string) error {
	dbInstanceId := beego.AppConfig.String("dbInstanceId")
	dbUser := beego.AppConfig.String("dbUser")

	// https://help.aliyun.com/document_detail/26266.html
	r := rds.CreateGrantAccountPrivilegeRequest()
	r.DBInstanceId = dbInstanceId
	r.AccountName = dbUser
	r.DBName = dbName
	r.AccountPrivilege = "ReadWrite"

	_, err := rdsClient.GrantAccountPrivilege(r)
	if err != nil {
		return err
	}

	return nil
}
