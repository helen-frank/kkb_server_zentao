package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kka-zentao-server/common/message"
	"kka-zentao-server/server/dboperate"
	"kka-zentao-server/server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ZenTaoAuthHandler(c *gin.Context) {
	// 用户发送用户名和密码过来
	var user message.UserInfo
	user.UserName = c.Query("username")
	user.Password = c.Query("password")
	var u []message.UserInfo
	err := json.Unmarshal(utils.ReadConfig("./config/Account.json"), &u)
	if err != nil {
		err = errors.New(fmt.Sprintln("network.go | json.Unmarshal failed , err: ", err))
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "打开./config/Account.json失败",
		})
		return
	}
	uMap := make(map[string]string)
	for _, v := range u {
		uMap[v.UserName] = v.Password
	}

	// 校验用户名和密码是否正确
	if uMap[user.UserName] == user.Password {
		// 生成Token
		tokenString, _ := utils.GenToken(user.UserName, user.Password)
		c.JSON(http.StatusOK, gin.H{
			"msg":   "success",
			"token": tokenString,
		})
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{
		"msg": "鉴权失败",
	})
}

// 添加zentao用户
func ZenTaoInsertUserHandler(zts *dboperate.ZenTaoService) func(c *gin.Context) {
	return func(c *gin.Context) {
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
		var u message.Kkb
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
			c.JSON(http.StatusAccepted, gin.H{
				"code": 2003,
				"msg":  "请求体中token为空",
			})
			c.Abort()
			return
		}
		mc, err := utils.ParseToken(u.Token)
		fmt.Println(mc)
		if err != nil {
			c.JSON(http.StatusAccepted, gin.H{
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
	}
}

// 添加用户到项目
func ZenTaoInsertUserProjectHandler(zts *dboperate.ZenTaoService) func(c *gin.Context) {
	return func(c *gin.Context) {
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
		var u message.Kkb
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
			c.JSON(http.StatusAccepted, gin.H{
				"code": 2003,
				"msg":  "请求体中token为空",
			})
			c.Abort()
			return
		}
		mc, err := utils.ParseToken(u.Token)
		fmt.Println(mc)
		if err != nil {
			c.JSON(http.StatusAccepted, gin.H{
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
		// 校验

		var replay string
		err = zts.ZenTaoInsertUserProject(u, &replay)
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
	}
}
