package message

const Hextable = "0123456789abcdef"

type SqlConfig struct {
	Account      string `json:"account"`
	Password     string `json:"password"`
	Ip           string `json:"ip"`
	Port         string `json:"port"`
	Database     string `json:"database"`
	MaxIdleConns int    `json:"maxIdleConns"`
}

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
	// 暂时未知参数，设为 nil
	Commiter interface{} `json:"commiter"`
	Avatar   interface{} `json:"avatar"`
	Nature   interface{} `json:"nature"`
	Analysis interface{} `json:"analysis"`
	Strategy interface{} `json:"strategy"`
}

type KkbUser struct {
	RealName     string `json:"realName"`
	MobileNumber string `json:"mobileNumber"`
	Email        string `json:"email"`
	Gender       string `json:"gender"`
}
