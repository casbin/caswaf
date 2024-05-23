// Copyright 2024 The casbin Authors. All Rights Reserved.
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
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/casbin/caswaf/util"
	"github.com/xorm-io/core"
)

type Record struct {
	Id          int64  `xorm:"int notnull pk autoincr" json:"id"`
	Owner       string `xorm:"varchar(100) notnull" json:"owner"`
	CreatedTime string `xorm:"varchar(100) notnull" json:"createdTime"`

	Method    string `xorm:"varchar(100)" json:"method"`
	Host      string `xorm:"varchar(100)" json:"host"`
	Path      string `xorm:"varchar(100)" json:"path"`
	UserAgent string `xorm:"varchar(512)" json:"userAgent"`
}

func GetRecords(owner string) ([]*Record, error) {
	records := []*Record{}
	err := ormer.Engine.Asc("id").Asc("host").Find(&records, &Record{Owner: owner})
	if err != nil {
		return nil, err
	}

	return records, nil
}

func AddRecord(record *Record) (*Record, error) {
	_, err := ormer.Engine.Insert(record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func DeleteRecord(record *Record) (bool, error) {
	affected, err := ormer.Engine.ID(core.PK{record.Id}).Delete(&Record{})
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func UpdateRecord(owner string, id string, record *Record) (bool, error) {
	affected, err := ormer.Engine.ID(core.PK{record.Id}).AllCols().Update(record)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func GetRecord(owner string, id string) (*Record, error) {
	idNum, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Failed to transform id(string) to num: ", err)
	}
	record, err := getRecord(owner, int64(idNum))
	if err != nil {
		return nil, err
	}

	return record, nil
}

func getRecord(owner string, id int64) (*Record, error) {
	record := Record{Owner: owner, Id: id}
	existed, err := ormer.Engine.Get(&record)
	if err != nil {
		return nil, err
	}

	if existed {
		return &record, nil
	}
	return nil, nil
}

func LogRequest(r *http.Request) {
	if !strings.Contains(r.UserAgent(), "Uptime-Kuma") {
		fmt.Printf("handleRequest: %s\t%s\t%s\t%s\t%s\n", r.RemoteAddr, r.Method, r.Host, r.RequestURI, r.UserAgent())
		record := Record{
			Owner:       "admin",
			CreatedTime: util.GetCurrentFormattedTime(),
			Method:      r.Method,
			Host:        r.Host,
			Path:        r.RequestURI,
			UserAgent:   r.UserAgent(),
		}
		fmt.Println(util.GetCurrentTime())
		_, err := AddRecord(&record)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
