package utils

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"kka-zentao-server/common/message"
	"os"
)

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
		fmt.Println("utils | utils.go | readConfig() | 打开config.json文件失败，请把config.json文件放置于该程序同级目录")
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
