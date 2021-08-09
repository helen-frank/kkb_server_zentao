package db

import (
	"database/sql"
	"fmt"
)

func insert(db *sql.DB) {
	sqlStr := `insert into user(name,age) values("加油呀",28)`
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
