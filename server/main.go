/*
	@file	main.go
	@author	helenfrank(helenfrank@protonmail.com)
	@date	2021-08-08 19:46:16
*/

package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"kka-zentao-server/common/message"
	"kka-zentao-server/server/dboperate"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db_zentao *sql.DB //连接池对象

func main() {
	db_zentao = dboperate.LinkSql("./config_zentao.json")
	defer func() {
		if db_zentao != nil {
			db_zentao.Close()
		}
	}()
	zts := dboperate.ZenTaoService{
		Db_zentao: db_zentao,
	}

	r := gin.Default()
	userGroup := r.Group("/user")
	{
		// userGroup.GET("/inquire", func(c *gin.Context) {})
		userGroup.POST("/ZenTaoInsertUser", func(c *gin.Context) {
			u := message.KkbUser{}
			var replay string
			b, err := c.GetRawData()
			if err != nil {
				c.JSON(http.StatusAccepted, gin.H{
					"err": err,
				})
				err = errors.New(fmt.Sprintln("main.go | userGroup.POST(\"/ZenTaoInsertUser\") | c.GetRawData failed , err: ", err))
				fmt.Println(err)
				return
			}

			err = json.Unmarshal(b, &u)
			if err != nil {
				c.JSON(http.StatusAccepted, gin.H{
					"err": err,
				})
				err = errors.New(fmt.Sprintln("main.go | userGroup.POST(\"/ZenTaoInsertUser\") | json.Unmarshal failed , err: ", err))
				fmt.Println(err)
				return
			}
			err = zts.ZenTaoInsertUser(u, &replay)
			if err != nil {
				c.JSON(http.StatusAccepted, gin.H{
					"err":    err,
					"replay": replay,
				})
				err = errors.New(fmt.Sprintln("main.go | userGroup.POST(\"/ZenTaoInsertUser\") | zts.ZenTaoInsertUser failed , err: ", err))
				fmt.Println(err)
				return
			}

			// 执行完毕
			c.JSON(http.StatusOK, gin.H{
				"replay": replay,
			})
		})
	}
	r.Run(":10227")
}
