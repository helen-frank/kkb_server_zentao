package network

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kkb-zentao-server/common/message"
	"kkb-zentao-server/server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Token User 校验中间件
func ZenTaoUserTokenCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		var u1 []message.UserInfo
		err := json.Unmarshal(utils.ReadConfig("etc/Account.json"), &u1)
		if err != nil {
			c.JSON(http.StatusAccepted, gin.H{
				"msg": "打开etc/Account.json失败",
			})
			err = errors.New(fmt.Sprintln("network.go | ZenTaoTokenCheck() | json.Unmarshal(utils.ReadConfig(\"etc/Account.json\"), &u1) failed , err: ", err))

			fmt.Println(err)
			return
		}
		var u message.Kkb

		// 读取request body 后恢复request body
		b, err := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(b))

		if err != nil {
			c.JSON(http.StatusAccepted, gin.H{
				"msg": "请将数据json化传递,或者是io.ReadAll(c.Request.Body) failed",
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
		//fmt.Println(mc)
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
	}
}

// Token Project 校验中间件
func ZenTaoProjectTokenCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		var u1 []message.UserInfo
		err := json.Unmarshal(utils.ReadConfig("etc/Account.json"), &u1)
		if err != nil {
			c.JSON(http.StatusAccepted, gin.H{
				"msg": "打开etc/Account.json失败",
			})
			err = errors.New(fmt.Sprintln("network.go | ZenTaoTokenCheck() | json.Unmarshal(utils.ReadConfig(\"etc/Account.json\"), &u1) failed , err: ", err))

			fmt.Println(err)
			return
		}
		var u message.KkbProject

		// 读取request body 后恢复request body
		b, err := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(b))

		if err != nil {
			c.JSON(http.StatusAccepted, gin.H{
				"msg": "请将数据json化传递,或者是io.ReadAll(c.Request.Body) failed",
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
		//fmt.Println(mc)
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
	}
}
