package dboperate

import (
	"database/sql"
	"fmt"
	"kkb-zentao-server/common/message"
	"kkb-zentao-server/server/utils"
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

	zenTaoInsertUser := "INSERT IGNORE INTO zt_user(company,account,password,realname,gender,role,mobile,email,commiter,avatar,nature,analysis,strategy) values (?,?,?,?,?,?,?,?,?,?,?,?,?)"

	r, err := zts.Db_zentao.Exec(zenTaoInsertUser, u.Company, u.Account, u.Password, u.RealName, u.Gender, u.Role, u.MobileNumber, u.Email, u.Commiter, u.Avatar, u.Nature, u.Analysis, u.Strategy)
	if err != nil {
		fmt.Printf("user.go | KkbUserLookUp | db2.Exec -> zenTaoInsertUser failed, err: %v\n", err)
		return err
	}

	// 后面写个日志
	fmt.Print("添加到zt_user | ")
	fmt.Println(r.LastInsertId())
	*replay = fmt.Sprintln(r.LastInsertId())

	zenTaoInsertUserGroup := "INSERT IGNORE INTO zt_usergroup(`account`,`group`,`project`) values (?,?,?)"
	r, err = zts.Db_zentao.Exec(zenTaoInsertUserGroup, u.Account, 2, nil)

	if err != nil {
		fmt.Printf("user.go | KkbUserLookUp | zts.Db_zentao.Exec -> zenTaoInsertUserGroup failed, err: %v\n", err)
		return err
	}
	fmt.Print("添加到zt_usergroup | ")
	fmt.Println(r.LastInsertId())
	return err
}

func (zts *ZenTaoService) ZenTaoInsertUserProject(k message.Kkb, replay *string) error {
	up := message.UserProject{
		Role:  "研发",
		Hours: 7.0,
		Kkb: message.Kkb{
			Root:    k.Root,
			Account: k.Account,
			Days:    k.Days,
		},
	}
	zenTaoInsertUserProject := "INSERT IGNORE INTO zt_team(root,account,days,role,hours) values (?,?,?,?,?)"

	r, err := zts.Db_zentao.Exec(zenTaoInsertUserProject, up.Root, up.Account, up.Days, up.Role, up.Hours)
	if err != nil {
		fmt.Printf("zenTaoService.go | ZenTaoInsertUserProject | zts.Db_zentao.Exec -> zenTaoInsertUserProject failed, err: %v\n", err)
		return err
	}

	// 后面写个日志
	fmt.Print("添加到zt_team | ")
	fmt.Println(r.LastInsertId())
	*replay = fmt.Sprintln(r.LastInsertId())
	return err
}
