package dboperate

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"kkb-zentao-server/common/message"
	"kkb-zentao-server/server/utils"
	"net/http"

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

// 查询kkb用户数据并插入到zentao
func KkbUserInsert(db1 *sql.DB, sqlStr, url, token string) {
	var u message.Kkb
	u.Token = token
	rows, err := db1.Query(sqlStr)

	if err != nil {
		fmt.Printf("query failed, err: %v\n", err)
		return
	}
	// 非常重要：关闭rows释放持有的数据库链接
	defer rows.Close()

	// 循环读取结果集中的数据
	for rows.Next() {
		if err != nil {
			err = errors.New(fmt.Sprintln("连接失败", err))
			fmt.Println(err)
			return
		}
		err = rows.Scan(&u.RealName, &u.MobileNumber, &u.Email)
		if err != nil {
			err = errors.New(fmt.Sprintln("user.go | KkbUserLookUp | rows.scan failed, err:\n", err))
			fmt.Println(err)
			return
		}
		u.Gender = "f"
		u.Account = u.MobileNumber
		u.Password = u.MobileNumber
		jsonStr, err := json.Marshal(u)
		if err != nil {
			err = errors.New(fmt.Sprintln("user.go | KkbUserLookUp | json.Marshal failed, err:\n", err))
			fmt.Println(err)
			return
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		if err != nil {
			err = errors.New(fmt.Sprintln("user.go | KkbUserLookUp | json.Marshal failed, err:\n", err))
			fmt.Println(err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			err = errors.New(fmt.Sprintln("user.go | KkbUserLookUp | client.Do failed, err:\n", err))
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		statuscode := resp.StatusCode
		//hea := resp.Header
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			err = errors.New(fmt.Sprintln("user.go | KkbUserLookUp | ioutil.ReadAll failed, err:\n", err))
			fmt.Println(err)
			return
		}
		fmt.Println(string(body))
		fmt.Println(statuscode)
	}
}

// 查询kkb用户数据并插入到zentao
func KkbUserProjectInsert(db1 *sql.DB, sqlStr, url, token string) {
	var u message.KkbProject
	var tempKkbAccout string
	u.Token = token
	rows, err := db1.Query(sqlStr)

	if err != nil {
		fmt.Printf("query failed, err: %v\n", err)
		return
	}
	// 非常重要：关闭rows释放持有的数据库链接
	defer rows.Close()
	t := 0
	u.Root = 2
	u.Days = 7
	// 循环读取结果集中的数据
	for rows.Next() {
		if err != nil {
			err = errors.New(fmt.Sprintln("连接失败", err))
			fmt.Println(err)
			return
		}
		err = rows.Scan(&tempKkbAccout)
		u.Account = append(u.Account, tempKkbAccout)
		if err != nil {
			err = errors.New(fmt.Sprintln("user.go | KkbUserLookUp | rows.scan failed, err:\n", err))
			fmt.Println(err)
			return
		}
		t++
	}
	jsonStr, err := json.Marshal(u)
	fmt.Println(string(jsonStr))
	if err != nil {
		err = errors.New(fmt.Sprintln("user.go | KkbUserLookUp | json.Marshal failed, err:\n", err))
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		err = errors.New(fmt.Sprintln("user.go | KkbUserLookUp | json.Marshal failed, err:\n", err))
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err = errors.New(fmt.Sprintln("user.go | KkbUserLookUp | client.Do failed, err:\n", err))
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	statuscode := resp.StatusCode
	//hea := resp.Header
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		err = errors.New(fmt.Sprintln("user.go | KkbUserLookUp | ioutil.ReadAll failed, err:\n", err))
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
	fmt.Println(statuscode)
}
