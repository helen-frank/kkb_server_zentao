package dboperate

import (
	"database/sql"
	"fmt"
	"kkb-zentao-server/common/message"
	"kkb-zentao-server/server/utils"
	"strconv"
)

type ZenTaoService struct {
	Db_zentao *sql.DB
}

func (zts *ZenTaoService) ZenTaoInsertUser(ku message.Kkb, replay *string) error {
	u := message.User{
		RealName:     ku.RealName,
		Account:      ku.Account,
		Password:     utils.FastMD5(ku.Password),
		Company:      0,
		Role:         "dev",
		MobileNumber: ku.MobileNumber,
		Email:        ku.Email,
		DeptId:       ku.DeptId,
		DeptName:     ku.DeptName,
		Commiter:     nil,
		Avatar:       nil,
		Nature:       nil,
		Analysis:     nil,
		Strategy:     nil,
	}
	if ku.Gender == "男" || ku.Gender == "man" {
		u.Gender = "m"
	} else {
		u.Gender = "f"
	}

	zenTaoInsertUserStr := "INSERT IGNORE INTO zt_user(company,account,password,realname,gender,role,mobile,email,dept,commiter,avatar,nature,analysis,strategy) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

	r, err := zts.Db_zentao.Exec(zenTaoInsertUserStr, u.Company, u.Account, u.Password, u.RealName, u.Gender, u.Role, u.MobileNumber, u.Email, u.DeptId, u.Commiter, u.Avatar, u.Nature, u.Analysis, u.Strategy)
	if err != nil {
		fmt.Printf("user.go | KkbUserLookUp | db2.Exec -> zenTaoInsertUserStr failed, err: %v\n", err)
		fmt.Println(u)
		return err
	}
	fmt.Println()
	// 后面写个日志
	fmt.Print("添加到zt_user | ")
	fmt.Println(r.LastInsertId())
	*replay = fmt.Sprintln(r.LastInsertId())

	zenTaoInsertUserGroup := "INSERT IGNORE INTO  zt_usergroup(`account`,`group`,`project`) values (?,?,?)"
	r, err = zts.Db_zentao.Exec(zenTaoInsertUserGroup, u.Account, 2, nil)

	if err != nil {
		fmt.Printf("user.go | KkbUserLookUp | zts.Db_zentao.Exec -> zenTaoInsertUserGroup failed, err: %v\n", err)
		return err
	}
	fmt.Print("添加到zt_usergroup | ")
	fmt.Println(r.LastInsertId())

	zenTaoInsertUserDept := "INSERT IGNORE INTO zt_dept(`id`,`name`,`path`) values (?,?,?)"
	r, err = zts.Db_zentao.Exec(zenTaoInsertUserDept, u.DeptId, u.DeptName, ","+strconv.Itoa(u.DeptId)+",")
	if err != nil {
		fmt.Printf("user.go | KkbUserLookUp | zts.Db_zentao.Exec -> zenTaoInsertUserDept failed, err: %v\n", err)
		return err
	}
	fmt.Print("添加到zt_dept | ")
	fmt.Println(r.LastInsertId())
	return err
}
