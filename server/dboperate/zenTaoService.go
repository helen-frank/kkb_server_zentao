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

	zenTaoInsertUserStr := "INSERT IGNORE INTO zt_user(company,account,password,realname,gender,role,mobile,email,commiter,avatar,nature,analysis,strategy) values (?,?,?,?,?,?,?,?,?,?,?,?,?)"

	r, err := zts.Db_zentao.Exec(zenTaoInsertUserStr, u.Company, u.Account, u.Password, u.RealName, u.Gender, u.Role, u.MobileNumber, u.Email, u.Commiter, u.Avatar, u.Nature, u.Analysis, u.Strategy)
	if err != nil {
		fmt.Printf("user.go | KkbUserLookUp | db2.Exec -> zenTaoInsertUserStr failed, err: %v\n", err)
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

func (zts *ZenTaoService) ZenTaoInsertUserProject(k message.KkbProject, replay *string) (err error) {
	up := message.UserProject{
		Role:  "研发",
		Hours: 7.0,
		KkbProject: message.KkbProject{
			Root:    k.Root,
			Account: k.Account,
			Days:    k.Days,
		},
	}
	zenTaoInsertUserProjectStr := "INSERT IGNORE INTO zt_team(root,account,days,role,hours) values (?,?,?,?,?)"
	for _, v := range up.Account {
		r, err := zts.Db_zentao.Exec(zenTaoInsertUserProjectStr, up.Root, v, up.Days, up.Role, up.Hours)
		if err != nil {
			fmt.Printf("zenTaoService.go | ZenTaoInsertUserProject | zts.Db_zentao.Exec -> zenTaoInsertUserProjectStr failed, err: %v\n", err)
			return err
		}

		// 后面写个日志
		fmt.Print("添加到zt_team | ")
		fmt.Println(r.LastInsertId())
		*replay += fmt.Sprintln(r.LastInsertId())
	}
	return err
}
