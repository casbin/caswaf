package object

import (
	"fmt"
	"github.com/xorm-io/core"
	"strconv"
)

type Record struct {
	Id          int64  `xorm:"int notnull pk autoincr" json:"id"`
	Owner       string `xorm:"varchar(100) notnull" json:"owner"`
	CreatedTime string `xorm:"varchar(100) notnull" json:"createdTime"`

	Method     string `xorm:"varchar(100)" json:"method"`
	Host       string `xorm:"varchar(100)" json:"host"`
	RequestURI string `xorm:"varchar(100)" json:"requestURI"`
	UserAgent  string `xorm:"varchar(512)" json:"userAgent"`
}

func GetRecords(owner string) ([]*Record, error) {
	records := []*Record{}
	err := ormer.Engine.Asc("Id").Desc("created_time").Asc("Host").Asc("request_u_r_i").Find(&records, &Record{Owner: owner})
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
	//id_num, err := strconv.ParseInt(id, 10, 64)
	//if err != nil {
	//	fmt.Println("Failed to transform id(string) to num: ", err)
	//}

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
