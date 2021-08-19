/*
	@file	main.go
	@author	helenfrank(helenfrank@protonmail.com)
	@date	2021-08-08 19:46:16
*/

package main

import (
	"database/sql"
	"fmt"
	"kkb-zentao-server/server/dboperate"
	"kkb-zentao-server/server/network"
	"kkb-zentao-server/server/utils"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/ini.v1"
)

var db_zentao *sql.DB //连接池对象

func main() {
	// file, _ := exec.LookPath(os.Args[0])
	// path, _ := filepath.Abs(file)
	// index := strings.LastIndex(path, string(os.PathSeparator))
	// path = path[:index]
	path := utils.ObtainPath()
	cfg, err := ini.Load(path + "/etc/my.cnf")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	Port := cfg.Section("server").Key("port").String()
	db_zentao = dboperate.LinkSql(path + "/etc/config_zentao.json")
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

	userGroup := r.Group("/user", network.ZenTaoUserTokenCheck())
	{
		userGroup.POST("/ZenTaoInsertUser", network.ZenTaoInsertUserHandler(&zts))
	}

	projectGroup := r.Group("/project", network.ZenTaoProjectTokenCheck())
	{
		projectGroup.POST("/ZenTaoInsertUserProject", network.ZenTaoInsertUserProjectHandler(&zts))
	}
	r.Run(":" + Port)
}
