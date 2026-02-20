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
	"strconv"
	"time"

	"github.com/xorm-io/core"
)

type Record struct {
	Id          int64  `xorm:"int notnull pk autoincr" json:"id"`
	Owner       string `xorm:"varchar(100) notnull" json:"owner"`
	CreatedTime string `xorm:"varchar(100) notnull" json:"createdTime"`

	Method    string `xorm:"varchar(100)" json:"method"`
	Host      string `xorm:"varchar(100)" json:"host"`
	Path      string `xorm:"varchar(100)" json:"path"`
	ClientIp  string `xorm:"varchar(100)" json:"clientIp"`
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

func AddRecord(record *Record) (bool, error) {
	affected, err := ormer.Engine.Insert(record)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
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
		return nil, err
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

type DataCount struct {
	Data  string `json:"data"`
	Count int64  `json:"count"`
}

// unixTimestampExpr returns a SQL expression that converts a varchar timestamp
// column to a Unix epoch integer, compatible with the configured DB driver.
func unixTimestampExpr(col string) string {
	switch ormer.driverName {
	case "postgres":
		return "CAST(EXTRACT(EPOCH FROM " + col + "::timestamp) AS bigint)"
	case "mssql":
		return "DATEDIFF(SECOND, '1970-01-01', CONVERT(datetime, " + col + "))"
	case "sqlite":
		return "strftime('%s', " + col + ")"
	default: // mysql
		return "UNIX_TIMESTAMP(" + col + ")"
	}
}

// dateFormatExpr returns a SQL expression that formats a varchar timestamp
// column using the given driver-appropriate format string.
func dateFormatExpr(col, format string) string {
	switch ormer.driverName {
	case "postgres":
		return "TO_CHAR(" + col + "::timestamp, '" + format + "')"
	case "mssql":
		return "FORMAT(CONVERT(datetime, " + col + "), '" + format + "')"
	case "sqlite":
		return "strftime('" + format + "', " + col + ")"
	default: // mysql
		return "DATE_FORMAT(" + col + ", '" + format + "')"
	}
}

func GetMetrics(dataType string, startAt time.Time, top int) (*[]DataCount, error) {
	var dataCounts []DataCount
	err := ormer.Engine.Table("record").
		Where(unixTimestampExpr("created_time")+" > ?", startAt.Unix()).
		Select(dataType + " as data, COUNT(*) as count").
		GroupBy("data").
		Desc("count").
		Limit(top).
		Find(&dataCounts)
	if err != nil {
		return nil, err
	}
	return &dataCounts, nil
}

func GetMetricsOverTime(startAt time.Time, timeType string) (*[]DataCount, error) {
	var dataCounts []DataCount
	createdTime := dateFormatExpr("created_time", timeType2Format(timeType))
	err := ormer.Engine.Table("record").
		Where(unixTimestampExpr("created_time")+" > ?", startAt.Unix()).
		GroupBy(createdTime).
		Select(createdTime + " as data, COUNT(*) as count").
		Asc("data").
		Find(&dataCounts)
	if err != nil {
		return nil, err
	}
	return &dataCounts, nil
}

// timeType2Format returns a DB-appropriate strftime/date-format pattern for
// the given time granularity. The format is chosen to match the active driver.
func timeType2Format(timeType string) string {
	isPostgres := ormer.driverName == "postgres"
	isMssql := ormer.driverName == "mssql"

	switch timeType {
	case "hour":
		if isPostgres {
			// Postgres uses to_char with its own format tokens.
			return "YYYY-MM-DD HH24"
		}
		if isMssql {
			// SQL Server FORMAT() uses .NET-style format strings.
			return "yyyy-MM-dd HH"
		}
		if ormer.driverName == "sqlite" {
			return "%Y-%m-%d %H"
		}
		return "%Y-%m-%d %H"
	case "day":
		if isPostgres {
			return "YYYY-MM-DD"
		}
		if isMssql {
			return "yyyy-MM-dd"
		}
		return "%Y-%m-%d"
	case "month":
		if isPostgres {
			return "YYYY-MM"
		}
		if isMssql {
			return "yyyy-MM"
		}
		return "%Y-%m"
	}

	if isPostgres {
		return "YYYY-MM-DD HH24"
	}
	if isMssql {
		return "yyyy-MM-dd HH"
	}
	return "%Y-%m-%d %H"
}

func GetRecordCount(owner, field, value string) (int64, error) {
	session := GetSession(owner, -1, -1, field, value, "", "")
	return session.Count(&Record{})
}

func GetPaginationRecords(owner string, offset, limit int, field, value, sortField, sortOrder string) ([]*Record, error) {
	records := []*Record{}
	session := GetSession(owner, offset, limit, field, value, sortField, sortOrder)
	err := session.Where("owner = ? or owner = ?", "admin", owner).Find(&records)
	if err != nil {
		return records, err
	}

	return records, nil
}
