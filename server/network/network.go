package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"kka-zentao-server/common/message"
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
