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
	"io"
	"kka-zentao-server/common/message"
	"kka-zentao-server/server/dboperate"
	"kka-zentao-server/server/network"
	"kka-zentao-server/server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db_zentao *sql.DB //连接池对象

func main() {
	db_zentao = dboperate.LinkSql("./config/config_zentao.json")
	defer func() {
		if db_zentao != nil {
			db_zentao.Close()
		}
	}()
	zts := dboperate.ZenTaoService{
		Db_zentao: db_zentao,
	}

	r := gin.Default()
	r.GET("/auth", network.ZenTaoAuthHandler) // 获取token

	userGroup := r.Group("/user")
	{

		userGroup.POST("/ZenTaoInsertUser", func(c *gin.Context) {
			var u1 []message.UserInfo
			err := json.Unmarshal(utils.ReadConfig("./config/Account.json"), &u1)
			if err != nil {
				c.JSON(http.StatusAccepted, gin.H{
					"msg": "打开./config/Account.json失败",
				})
				err = errors.New(fmt.Sprintln("main.go | userGroup.POST(\"/ZenTaoInsertUser\") | json.Unmarshal(utils.ReadConfig(\"./config/Account.json\"), &u1) failed , err: ", err))
				fmt.Println(err)
				return
			}
			var u message.KkbUser
			b, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusAccepted, gin.H{
					"msg": "请将数据json化传递",
				})
				err = errors.New(fmt.Sprintln("main.go | userGroup.POST(\"/ZenTaoInsertUser\") | io.ReadAll(c.Request.Body) failed , err: ", err))
				fmt.Println(err)
				return
			}

			err = json.Unmarshal(b, &u)
			if err != nil {
				c.JSON(http.StatusAccepted, gin.H{
					"msg": "json解析失败",
				})
				err = errors.New(fmt.Sprintln("main.go | userGroup.POST(\"/ZenTaoInsertUser\") | json.Unmarshal(b, &u) failed , err: ", err))
				fmt.Println(err)
				return
			}
			if u.Token == "" {
				c.JSON(http.StatusOK, gin.H{
					"code": 2003,
					"msg":  "请求体中token为空",
				})
				c.Abort()
				return
			}
			mc, err := utils.ParseToken(u.Token)
			fmt.Println(mc)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": 2005,
					"msg":  "无效的token",
				})

				return
			}

			exist := false
			for _, v := range u1 {
				if v.UserName == mc.UserName {
					exist = true
					break
				}

			}
			if !exist {
				c.JSON(http.StatusAccepted, gin.H{
					"username": mc.UserName,
					"msg":      "用户不存在，请重新获取token",
				})
				return
			}
			var replay string

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
