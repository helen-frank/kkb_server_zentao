package dboperate

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"kka-zentao-server/common/message"
	"kka-zentao-server/server/utils"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

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
	fmt.Println("connect kkb sql success")
	return db
}

// 查询kkb用户数据
func KkbUserLookUp(db1 *sql.DB, sqlStr string) {
	var u message.User
	var replay string

	rows, err := db1.Query(sqlStr)

	if err != nil {
		fmt.Printf("query failed, err: %v\n", err)
		return
	}
	// 非常重要：关闭rows释放持有的数据库链接
	defer rows.Close()

	// 循环读取结果集中的数据
	for rows.Next() {
		conn, err := net.Dial("tcp", ":10227")
		if err != nil {
			err = errors.New(fmt.Sprintln("连接失败", err))
			fmt.Println(err)
			return
		}
		client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
		err = rows.Scan(&u.RealName, &u.MobileNumber, &u.Email)
		if err != nil {
			err = errors.New(fmt.Sprintln("user.go | KkbUserLookUp | rows.scan failed, err:\n", err))
			fmt.Println(err)
			return
		}
		u.Gender = "f"
		err = client.Call("zts.ZenTaoInsertUser", u, &replay)
		if err != nil {
			err = errors.New(fmt.Sprintln("user.go | KkbUserLookUp | client.Call -> zts.ZenTaoInsertUser 调用失败,err: ", err))
			fmt.Println(err)
			return
		}
		fmt.Print(replay)

	}
}
