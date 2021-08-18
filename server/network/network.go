package network

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kkb-zentao-server/common/message"
	"kkb-zentao-server/server/dboperate"
	"kkb-zentao-server/server/utils"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// 生成token
func ZenTaoAuthHandler(c *gin.Context) {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	path = path[:index]
	// 用户发送用户名和密码过来
	var user message.UserInfo
	user.UserName = c.Query("username")
	user.Password = c.Query("password")
	var u []message.UserInfo
	err := json.Unmarshal(utils.ReadConfig(path+"/etc/Account.json"), &u)
	if err != nil {
		err = errors.New(fmt.Sprintln("network.go | json.Unmarshal failed , err: ", err))
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "打开etc/Account.json失败",
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
		// 校验用户用户正确，保留账户密码以及token
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{
		"msg": "鉴权失败",
		// 写明失败原因
	})
}

// Token 校验中间件
func ZenTaoTokenCheck() gin.HandlerFunc {
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
	}
}

// 添加zentao用户
func ZenTaoInsertUserHandler(zts *dboperate.ZenTaoService) func(c *gin.Context) {
	return func(c *gin.Context) {

		var replay string
		var u message.Kkb
		b, err := io.ReadAll(c.Request.Body)
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
		var replay string
		var u message.Kkb
		b, err := io.ReadAll(c.Request.Body)
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
