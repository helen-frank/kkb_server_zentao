package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kkb-zentao-server/common/message"
	"kkb-zentao-server/server/dboperate"
	"kkb-zentao-server/server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 生成token
func ZenTaoAuthHandler(c *gin.Context) {
	path := utils.ObtainPath()
	// 用户发送ApiKey和SecretKey过来
	var user message.UserInfo
	user.ApiKey = c.Query("apiKey")
	user.SecretKey = c.Query("secretKey")
	var u []message.UserInfo
	err := json.Unmarshal(utils.ReadConfig(path+"/etc/Account.json"), &u)
	if err != nil {
		err = errors.New(fmt.Sprintln("network.go | json.Unmarshal failed , err: ", err))
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}
	uMap := make(map[string]string)
	for _, v := range u {
		uMap[v.ApiKey] = v.SecretKey
	}

	// 校验用户名和密码是否正确
	if uMap[user.ApiKey] == user.SecretKey {
		// 生成Token
		tokenString, _ := utils.GenToken(user.ApiKey, user.SecretKey)
		c.JSON(http.StatusOK, gin.H{
			"token": tokenString,
			"err":   nil,
		})
		// 校验用户用户正确，保留账户密码以及token
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{
		"err": "鉴权失败",
		// 写明失败原因
	})
}

// 添加zentao用户
func ZenTaoInsertUserHandler(zts *dboperate.ZenTaoService) func(c *gin.Context) {
	return func(c *gin.Context) {

		var replay string
		var u message.Kkb
		b, err := io.ReadAll(c.Request.Body)
		if err != nil {

			err = errors.New(fmt.Sprintln("main.go | userGroup.POST(\"/ZenTaoInsertUser\") | io.ReadAll(c.Request.Body) failed , err: ", err))
			fmt.Println(err)
			c.JSON(http.StatusOK, gin.H{
				"msg": "请将数据json化传递,或者是io.ReadAll(c.Request.Body) failed",
				"err": err,
			})
			return
		}
		err = json.Unmarshal(b, &u)
		if err != nil {
			err = errors.New(fmt.Sprintln("main.go | userGroup.POST(\"/ZenTaoInsertUser\") | json.Unmarshal(b, &u) failed , err: ", err))
			fmt.Println(err)
			c.JSON(http.StatusOK, gin.H{
				"msg": "json解析失败",
				"err": err,
			})
			return
		}
		err = zts.ZenTaoInsertUser(u, &replay)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"err":    err,
				"replay": replay,
			})
			err = errors.New(fmt.Sprintln("main.go | userGroup.POST(\"/ZenTaoInsertUser\") | zts.ZenTaoInsertUser failed , err: ", err))
			fmt.Println(err)
			return
		}

		// 执行完毕
		c.JSON(http.StatusOK, gin.H{
			"err":    err,
			"replay": replay,
		})
	}
}

// 添加用户到项目
func ZenTaoInsertUserProjectHandler(zts *dboperate.ZenTaoService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var replay string
		var u message.KkbProject
		b, err := io.ReadAll(c.Request.Body)
		if err != nil {

			err = errors.New(fmt.Sprintln("main.go | userGroup.POST(\"/ZenTaoInsertUser\") | io.ReadAll(c.Request.Body) failed , err: ", err))
			fmt.Println(err)
			c.JSON(http.StatusOK, gin.H{
				"msg": "请将数据json化传递,或者是io.ReadAll(c.Request.Body) failed",
				"err": err,
			})
			return
		}
		err = json.Unmarshal(b, &u)
		if err != nil {
			err = errors.New(fmt.Sprintln("main.go | userGroup.POST(\"/ZenTaoInsertUser\") | json.Unmarshal(b, &u) failed , err: ", err))
			fmt.Println(err)
			c.JSON(http.StatusOK, gin.H{
				"msg": "json解析失败",
				"err": err,
			})
			return
		}
		err = zts.ZenTaoInsertUserProject(u, &replay)
		if err != nil {
			err = errors.New(fmt.Sprintln("main.go | userGroup.POST(\"/ZenTaoInsertUser\") | zts.ZenTaoInsertUser failed , err: ", err))
			fmt.Println(err)
			c.JSON(http.StatusOK, gin.H{
				"err":    err,
				"replay": replay,
			})
			return
		}

		// 执行完毕
		c.JSON(http.StatusOK, gin.H{
			"err":    err,
			"replay": replay,
		})
	}
}
