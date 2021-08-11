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
	UserId       int    `json:"userId"`
	RealName     string `json:"realName"`
	Email        string `json:"email"`
	MobileNumber string `json:"mobileNumber"`
	Company      int    `json:"company"`
	Account      string `json:"account"`
	Password     string `json:"password"`
	Role         string `json:"role"`
	Id           int    `json:"id"`
}
