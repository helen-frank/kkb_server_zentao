/*
	@file	main.go
	@author	helenfrank(helenfrank@protonmail.com)
	@date	2021-08-08 19:46:16
*/

package main

import (
	"database/sql"
	"kka-zentao-server/db"

	_ "github.com/go-sql-driver/mysql"
)

var db_zentao *sql.DB //连接池对象
var db_kkb *sql.DB

func initdb() {
	db_kkb = db.LinkSql("./config_kkb.json") //
	db_zentao = db.LinkSql("./config_zentao.json")
}

func main() {
	initdb()
	defer func() {
		if db_kkb != nil {
			db_kkb.Close()
		}
	}()

	defer func() {
		if db_zentao != nil {
			db_zentao.Close()
		}
	}()

	looUpStr := "SELECT suanke_student.user_id,suanke_student.real_name, suanke_user.email, suanke_user.mobile_number FROM suanke_user,suanke_student  where suanke_user.id=suanke_student.user_id  ORDER BY suanke_student.user_id limit 10;"
	db.KkbUserLookUp(db_kkb, db_zentao, looUpStr)

}
