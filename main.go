package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io/ioutil"
	"time"
)

type ConfigStruct struct {
	Url      string       `json:"url"`
	User     string       `json:"user"`
	Password string       `json:"password"`
	Tables   []string     `json:"tables"`
	Filter   FilterStruct `json:"filter"`
}

type FilterStruct struct {
	UseFilter bool     `json:"use_filter"`
	AllTables bool     `json:"all_tables"`
	Excludes  []string `json:"excludes"`
}

var ConfigInfo = make(map[string]ConfigStruct)
var MySQLDsn = "%s:%s@tcp(%s)/%s?charset=utf8&interpolateParams=true&parseTime=True&loc=Local"

func main() {
	configContent, err := ioutil.ReadFile("./db.json")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(configContent))
	if err := json.Unmarshal(configContent, &ConfigInfo); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", ConfigInfo["live-update"])
	for k, v := range ConfigInfo {
		eng := conn(k, v)
		for _,table:=range v.Tables{
			desc(eng, table)
		}
	}
}

func setDsn(dbName string, cs ConfigStruct) string {
	return fmt.Sprintf(MySQLDsn, cs.User, cs.Password, cs.Url, dbName)
}

func conn(dbName string, cs ConfigStruct) *xorm.Engine {
	eng, err := xorm.NewEngine("mysql", setDsn(dbName, cs))
	if err != nil {
		panic(err)
	}
	if err := eng.DB().Ping(); err != nil {
		panic(err)
	}
	// 打开调试模式
	eng.ShowSQL(true)
	eng.SetMaxOpenConns(10)
	eng.SetMaxIdleConns(10)
	eng.SetConnMaxLifetime(time.Second * time.Duration(1200))
	return eng
}

func desc(eng *xorm.Engine, table string) {
	sql := fmt.Sprintf("DESC `%s`", table)
	descInfo := make([]interface{},0)
	if err := eng.SQL(sql).Find(&descInfo); err != nil {
		panic(err)
	}
	fmt.Println(descInfo)

}
