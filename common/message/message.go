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
	UserId       string `json:"userId"`
	RealName     string `json:"realName"`
	Email        string `json:"email"`
	MobileNumber string `json:"mobileNumber"`
}
