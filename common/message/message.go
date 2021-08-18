package message

import (
	"time"
)

const Hextable = "0123456789abcdef"

// 定义JWT的过期时间，这里以2小时
const TokenExpireDuration = time.Hour * 2

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中

var SecretKey = []byte("划水的一天")

type SqlConfig struct {
	Account      string `json:"account"`
	Password     string `json:"password"`
	Ip           string `json:"ip"`
	Port         string `json:"port"`
	Database     string `json:"database"`
	MaxIdleConns int    `json:"maxIdleConns"`
}

// 开课吧添加用户
type Kkb struct {
	RealName     string `json:"realName"`
	MobileNumber string `json:"mobileNumber"`
	Email        string `json:"email"`
	Gender       string `json:"gender"`
	Token        string `json:"token"`
	Account      string `json:"account"`
	Password     string `json:"password"`
	Days         int    `json:"days"`
	Root         int    `json:"root"`
}

// ZenTao用户
type User struct {
	//UserId       int    `json:"userId"`
	//Id           int    `json:"id"`
	RealName     string `json:"realName"`
	Email        string `json:"email"`
	MobileNumber string `json:"mobileNumber"`
	Company      int    `json:"company"`
	Account      string `json:"account"`
	Password     string `json:"password"`
	Role         string `json:"role"`
	Gender       string `json:"gender"`
	Token        string `json:"token"`
	// 暂时未知参数，设为 nil
	Commiter interface{} `json:"commiter"`
	Avatar   interface{} `json:"avatar"`
	Nature   interface{} `json:"nature"`
	Analysis interface{} `json:"analysis"`
	Strategy interface{} `json:"strategy"`
}

// zentao项目
type UserProject struct {
	Kkb
	Role  string  `json:"role"`
	Hours float32 `json:"hours"`
}

// 获取token校验
type UserInfo struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}
