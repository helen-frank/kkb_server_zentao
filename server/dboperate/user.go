package dboperate

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"kka-zentao-server/common/message"
	"kka-zentao-server/server/utils"

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
	fmt.Println("connect zentao sql success")
	return db
}
