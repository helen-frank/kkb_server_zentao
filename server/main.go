/*
	@file	main.go
	@author	helenfrank(helenfrank@protonmail.com)
	@date	2021-08-08 19:46:16
*/

package main

import (
	"database/sql"
	"kka-zentao-server/server/dboperate"
	"kka-zentao-server/server/network"

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
		userGroup.POST("/ZenTaoInsertUser", network.ZenTaoInsertUserHandler(&zts))
	}

	projectGroup := r.Group("/project")
	{
		projectGroup.POST("/ZenTaoInsertUserProject", network.ZenTaoInsertUserProjectHandler(&zts))
	}
	r.Run(":10227")
}
