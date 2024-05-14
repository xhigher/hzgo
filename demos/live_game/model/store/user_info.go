package store

type UserInfo struct {
	Id       string `json:"id"`
	Openid   string `json:"openid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Token    string `json:"token"`
	Sex      int    `json:"sex"`
	Level    int    `json:"level"`
	Skin     int    `json:"role"`
}
