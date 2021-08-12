/*
	@file	main.go
	@author	helenfrank(helenfrank@protonmail.com)
	@date	2021-08-08 19:46:16
*/

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"kka-zentao-server/server/dboperate"
	"net/http"
	"net/rpc"

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

	r := gin.Default()
	userGroup := r.Group("/user")
	{
		// userGroup.GET("/inquire", func(c *gin.Context) {})
		userGroup.POST("/ZenTaoInsertUser", func(c *gin.Context) {
			b, err := c.GetRawData()

			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"err": err,
				})
				return
			}

		})
	}
	r.Run()

	rpc.RegisterName("zts", &dboperate.ZenTaoService{
		Db_zentao: db_zentao,
	})

	if err != nil {
		err = errors.New(fmt.Sprintln("main.go | net.Listen err:", err))
		panic(err)
	}

	fmt.Println("zentao server rpc start monitor Port 10227")
}
