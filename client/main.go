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
	"strconv"
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

	getToken := "http://39.98.121.3:10227/auth?apiKey=3C4mjFIdvqfgzl2LTER095V732f8VNp5&secretKey=bD0SDba8Hdn8721EIgUUtzQ6JgHe6clu"
	// getToken := "http://127.0.0.1:10227/auth?apiKey=test1&secretKey=010227"
	// getToken := "http://39.103.184.136:10227/auth?apiKey=test1&secretKey=010227"
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

	fmt.Println("KkbUserLookUp service start")

	// 添加用户
	// zentaouserinserturl := "http://127.0.0.1:10227/user/ZenTaoInsertUser"
	// UserInsertStr := "SELECT suanke_student.real_name,suanke_user.mobile_number,suanke_user.email FROM suanke_user,suanke_student  where suanke_user.id=suanke_student.user_id  ORDER BY suanke_student.user_id limit 100;"
	// dboperate.KkbUserInsert(db_kkb, UserInsertStr, zentaouserinserturl, gettokens.Token)

	// 添加用户到项目
	// zentaouserinserturl := "http://127.0.0.1:10227/project/ZenTaoInsertUserProject"
	// UserProjectInsertStr := "SELECT suanke_user.mobile_number FROM suanke_user,suanke_student  where suanke_user.id=suanke_student.user_id  ORDER BY suanke_student.user_id limit 5"
	// dboperate.KkbUserProjectInsert(db_kkb, UserProjectInsertStr, zentaouserinserturl, gettokens.Token)

	zenTaoInsertUserUrl := "http://39.98.121.3:10227/user/ZenTaoInsertUser"
	// zenTaoInsertUserUrl := "http://127.0.0.1:10227/user/ZenTaoInsertUser"
	// zenTaoInsertUserProjectUrl := "http://39.98.121.3:10227/project/ZenTaoInsertUserProject"
	// zentao内项目id
	var planId int
	// 选择planId课的电话
	//UserProjectInsertStr := "SELECT suanke_student.real_name,suanke_user.mobile_number,suanke_user.email,suanke_plan.id,suanke_plan.name  FROM suanke_user,suanke_student_plan,suanke_student,suanke_plan  where suanke_user.id=suanke_student.user_id and suanke_student.id=suanke_student_plan.student_id and suanke_student_plan.is_available !=0  and  suanke_student_plan.plan_id=" + strconv.Itoa(planId) + " and suanke_plan.id = suanke_student_plan.plan_id ORDER BY suanke_student.user_id"
	for i := 1; i <= 64; i++ {
		planId = i
		UserInsertStr := "SELECT suanke_student.real_name,suanke_user.mobile_number,suanke_user.email,suanke_plan.id,suanke_plan.name  FROM suanke_user,suanke_student_plan,suanke_student,suanke_plan  where suanke_user.id=suanke_student.user_id and suanke_student.id=suanke_student_plan.student_id and suanke_student_plan.is_available !=0  and  suanke_student_plan.plan_id=" + strconv.Itoa(planId) + " and suanke_plan.id = suanke_student_plan.plan_id ORDER BY suanke_student.user_id"
		dboperate.KkbUserInsert(db_kkb, UserInsertStr, zenTaoInsertUserUrl, gettokens.Token)
		fmt.Println(UserInsertStr)
	}
}

// SELECT suanke_student.real_name,suanke_user.mobile_number,suanke_user.email,suanke_plan.id,suanke_plan.name  FROM suanke_user,suanke_student_plan,suanke_student,suanke_plan  where suanke_user.id=suanke_student.user_id and suanke_student.id=suanke_student_plan.student_id and suanke_student_plan.is_available !=0  and  suanke_student_plan.plan_id=64 and suanke_plan.id = suanke_student_plan.plan_id ORDER BY suanke_student.user_id
