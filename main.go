/*
	@file	main.go
	@author	helenfrank(helenfrank@protonmail.com)
	@date	2021-08-08 19:46:16
*/

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"kka-zentao-server/common/message"
	"kka-zentao-server/utils"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB //连接池对象

func init() {
	sc := message.SqlConfig{}
	err := json.Unmarshal(utils.ReadConfig("./config.json"), &sc)
	if err != nil {
		fmt.Println("main.go | init() | json解析config.json文件失败")
		panic(err)
	}

	db, err = sql.Open("mysql", utils.StringStitching(sc.Account, ":", sc.Password, "@tcp(", sc.Ip, sc.Port, ")/", sc.Database))
	db.Ping()
	defer func() {
		if db != nil {
			db.Close()
		}
	}()
	if err != nil {
		fmt.Println("数据库链接失败")
		panic(err)
	}
	//设置数据库连接池的最大连接数
	db.SetMaxIdleConns(sc.MaxIdleConns)
}

func main() {

}
