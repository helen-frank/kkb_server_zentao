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
		path := utils.ObtainPath()
		err := json.Unmarshal(utils.ReadConfig(path+"/etc/Account.json"), &u1)
		if err != nil {
			err = errors.New(fmt.Sprintln("network.go | ZenTaoTokenCheck() | json.Unmarshal(utils.ReadConfig(\"etc/Account.json\"), &u1) failed , err: ", err))
			fmt.Println(err)
			c.JSON(http.StatusOK, gin.H{
				"msg": "打开etc/Account.json失败",
				"err": err,
			})
			return
		}
		var u message.Kkb

		// 读取request body 后恢复request body
		b, err := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(b))

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
		if u.Token == "" {
			c.JSON(http.StatusOK, gin.H{
				"err": "请求体中token为空",
			})
			c.Abort()
			return
		}
		mc, err := utils.ParseToken(u.Token)
		//fmt.Println(mc)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg": "无效的token",
				"err": err,
			})
			return
		}

		exist := false
		for _, v := range u1 {
			if v.SecretKey == mc.SecretKey {
				exist = true
				break
			}

		}
		if !exist {
			c.JSON(http.StatusOK, gin.H{
				"apiKey": mc.ApiKey,
				"err":    "apiKey不存在，请重新获取token",
			})
			return
		}
	}
}
