package dboperate

import (
	"database/sql"
	"fmt"
	"kka-zentao-server/common/message"
	"kka-zentao-server/server/utils"
)

type ZenTaoService struct {
	Db_zentao *sql.DB
}

func (zts *ZenTaoService) ZenTaoInsertUser(ku message.KkbUser, replay *string) error {
	u := message.User{
		RealName: ku.RealName,
		Account:  ku.Email,
		Password: utils.FastMD5(ku.MobileNumber),
		Company:  0,
		Role:     "dev",
		Commiter: nil,
		Avatar:   nil,
		Nature:   nil,
		Analysis: nil,
		Strategy: nil,
	}
	if ku.Gender == "男" || ku.Gender == "man" {
		u.Gender = "m"
	} else {
		u.Gender = "f"
	}

	zenTaoInsertUser := "INSERT IGNORE INTO zt_user(company,account,password,realname,gender,role,commiter,avatar,nature,analysis,strategy) values (?,?,?,?,?,?,?,?,?,?,?)"

	r, err := zts.Db_zentao.Exec(zenTaoInsertUser, u.Company, u.Account, u.Password, u.RealName, u.Gender, u.Role, u.Commiter, u.Avatar, u.Nature, u.Analysis, u.Strategy)
	if err != nil {
		fmt.Printf("user.go | KkbUserLookUp | db2.Exec -> zenTaoInsertUser failed, err:%v\n", err)
	}

	// 后面写个日志
	fmt.Println(r.LastInsertId())
	*replay = fmt.Sprintln(r.LastInsertId())
	return err
}
