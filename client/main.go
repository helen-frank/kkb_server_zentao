/*
	@file	main.go
	@author	helenfrank(helenfrank@protonmail.com)
	@date	2021-08-12 13:33:32
*/

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kkb-zentao-server/client/dboperate"
	"net/http"
)

var db_kkb *sql.DB

type GetToken struct {
	Token string `json:"token"`
	Msg   string `json:"msg"`
}

func main() {
	db_kkb = dboperate.LinkSql("./dboperate/config_kkb.json")
	defer func() {
		if db_kkb != nil {
			db_kkb.Close()
		}
	}()

	getToken := "http://127.0.0.1:10227/auth?username=test1&password=010227"
	// getToken := "http://39.103.184.136:10227/auth?username=test1&password=010227"
	resp1, err := http.Get(getToken)
	if err != nil {
		panic(err)
	}
	defer resp1.Body.Close()
	var gettokens GetToken
	b, _ := ioutil.ReadAll(resp1.Body)
	json.Unmarshal(b, &gettokens)
	fmt.Println(gettokens)
	if resp1.StatusCode != http.StatusOK {
		return
	}

	//zentaouserinserturl := "http://127.0.0.1:10227/user/ZenTaoInsertUser"
	fmt.Println("KkbUserLookUp service start")
	//UserInsertStr := "SELECT suanke_student.real_name,suanke_user.mobile_number,suanke_user.email FROM suanke_user,suanke_student  where suanke_user.id=suanke_student.user_id  ORDER BY suanke_student.user_id limit 100;"
	//dboperate.KkbUserInsert(db_kkb, UserInsertStr, zentaouserinserturl, gettokens.Token)

	zentaouserinserturl := "http://127.0.0.1:10227/project/ZenTaoInsertUserProject"
	UserProjectInsertStr := "SELECT suanke_user.mobile_number FROM suanke_user,suanke_student  where suanke_user.id=suanke_student.user_id  ORDER BY suanke_student.user_id limit 100"
	dboperate.KkbUserProjectInsert(db_kkb, UserProjectInsertStr, zentaouserinserturl, gettokens.Token)
}
