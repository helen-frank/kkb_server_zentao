package utils

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"kkb-zentao-server/common/message"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type MyClaims struct {
	ApiKey    string `json:"apiKey"`
	SecretKey string `json:"secretKey"`
	jwt.StandardClaims
}

func ObtainPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	path = path[:index]
	return path
}

// MD5字符串
func FastMD5(str string) string {
	src := md5.Sum([]byte(str))
	var dst = make([]byte, 32)
	j := 0
	for _, v := range src {
		dst[j] = message.Hextable[v>>4]
		dst[j+1] = message.Hextable[v&0x0f]
		j += 2
	}
	return string(dst)
}

// 读取config.json
func ReadConfig(filename string) []byte {
	fi, err := os.Open(filename)
	if err != nil {
		fmt.Println("utils | utils.go | readConfig() | 打开" + filename + "文件失败，请把文件放置于etc目录")
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		fmt.Println("utils | main.go | init() | ioutil.ReadAll(fi)失败")
		panic(err)
	}
	return fd
}

// 字符串拼接
func StringStitching(str ...string) string {
	var buf bytes.Buffer

	buf.Grow(100)
	for _, v := range str {
		buf.WriteString(v)
	}
	return buf.String()
}

// GenToken 生成JWT
func GenToken(apiKey, secretKey string) (string, error) {
	mc := MyClaims{
		ApiKey:    apiKey, // 自定义字段
		SecretKey: secretKey,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(message.TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "helen",                                            // 签发人
		},
	}

	// 使用指定的签名方法创建签名对象
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, mc)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return tokenClaims.SignedString(message.SecretKey)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	tokenClaims, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return message.SecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*MyClaims); ok && tokenClaims.Valid { // 校验token
			return claims, nil
		}
	}
	return nil, errors.New("invalid token")
}
