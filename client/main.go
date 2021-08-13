/*
	@file	main.go
	@author	helenfrank(helenfrank@protonmail.com)
	@date	2021-08-12 13:33:32
*/

package main

import (
	"database/sql"
	"fmt"
	"kka-zentao-server/client/dboperate"
)

var db_kkb *sql.DB

func main() {
	db_kkb = dboperate.LinkSql("./config_kkb.json")
	defer func() {
		if db_kkb != nil {
			db_kkb.Close()
		}
	}()

	looUpStr := "SELECT suanke_student.real_name,suanke_user.mobile_number,suanke_user.email FROM suanke_user,suanke_student  where suanke_user.id=suanke_student.user_id  ORDER BY suanke_student.user_id limit 100;"
	zentaourl := "http://127.0.0.1:10227/user/ZenTaoInsertUser"
	fmt.Println("KkbUserLookUp service start")
	dboperate.KkbUserLookUp(db_kkb, looUpStr, zentaourl)
}
