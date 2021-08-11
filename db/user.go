package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"kka-zentao-server/common/message"
	"kka-zentao-server/utils"

	_ "github.com/go-sql-driver/mysql"
)

func LinkSql(fileName string) *sql.DB {
	sc := message.SqlConfig{}
	err := json.Unmarshal(utils.ReadConfig(fileName), &sc)
	if err != nil {
		fmt.Println("main.go | init() | json解析" + fileName + "文件失败")
		panic(err)
	}
	db, err := sql.Open("mysql", utils.StringStitching(sc.Account, ":", sc.Password, "@tcp(", sc.Ip, sc.Port, ")/", sc.Database))
	db.Ping()

	if err != nil {
		fmt.Println(fileName + "链接失败")
		panic(err)
	}
	// 最大连接数
	// db.SetMaxOpenConns(10)
	// 设置的执行完闲置的连接，这些就算是执行结束了sql语句还是会保留着的

	db.SetMaxIdleConns(sc.MaxIdleConns)
	return db
}

func ZenTaoInsert(db *sql.DB, sqlStr string) {
	ret, err := db.Exec(sqlStr)
	if err != nil {
		fmt.Printf("插入失败,err:%v\n", err)
		return
	}
	//如果是插入数据的操作，能够拿到插入数据的id
	id, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get id failed,err:%v\n", err)
		return
	}
	fmt.Println("id", id)
}

// 查询kkb用户数据
func KkbUserLookUp(db1 *sql.DB, db2 *sql.DB, sqlStr string) {

	rows, err := db1.Query(sqlStr)

	if err != nil {
		fmt.Printf("query failed, err: %v\n", err)
		return
	}
	// 非常重要：关闭rows释放持有的数据库链接
	defer rows.Close()

	// 循环读取结果集中的数据
	for rows.Next() {
		var u message.User

		err := rows.Scan(&u.UserId, &u.RealName, &u.Email, &u.MobileNumber)
		if err != nil {
			fmt.Printf("user.go | KkbUserLookUp | rows.scan failed, err:%v\n", err)
			return
		}
		u.Account = u.MobileNumber
		u.Id = u.UserId + 1000
		if u.Email != "" {
			u.Password = utils.FastMD5(u.Email)
		} else {
			u.Password = utils.FastMD5("123456")
		}
		u.Company = 0
		u.Role = "dev"

		zenTaoinsertStr := "INSERT IGNORE INTO zt_user(id,company,account,password,role,commiter,avatar,email,mobile,nature,analysis,strategy) values (?,?,?,?,?,?,?,?,?,?,?,?)"

		r, err := db2.Exec(zenTaoinsertStr, u.Id, u.Company, u.Account, u.Password, u.Role, 1, 1, u.Email, u.MobileNumber, 1, 1, 1)
		if err != nil {
			fmt.Printf("user.go | KkbUserLookUp | db2.Exec failed, err:%v\n", err)
		}

		fmt.Println(r.LastInsertId())

		//fmt.Printf("id:%d name:%s email:%s mobilePhone:%s\n", u.UserId, u.RealName, u.Email, u.MobileNumber)

	}
}
