package user

import (
	"github.com/xhigher/hzgo/mysql"
	"gorm.io/gorm"
)

const (
	dbName = "hertz_user"
)

func DB() *gorm.DB {
	return mysql.DB(dbName)
}

type UserInfo struct {
	Userid   string `json:"userid" column:"userid"`
	Username string `json:"username" column:"username"`
	Password string `json:"password" column:"password"`
	Nid      string `json:"nid" column:"nid"`
	Nickname string `json:"nickname" column:"nickname"`
	Avatar   string `json:"avatar" column:"avatar"`
	Sex      int32  `json:"sex" column:"sex"`
	Birthday string `json:"birthday" column:"birthday"`
	Inviter  string `json:"inviter" column:"inviter"`
	Status   int32  `json:"status" column:"status"`
	Ct       int64  `json:"ct" column:"ct"`
	Ut       int64  `json:"ut" column:"ut"`
}

func (u *UserInfo) TableName() string {
	return "user_info"
}

type UserToken struct {
	Userid string `json:"userid" column:"userid"`
	Token  string `json:"token" column:"token"`
	Et     int64  `json:"et" column:"et"`
	It     int64  `json:"it" gorm:"it"`
}

func (u *UserToken) TableName() string {
	return "user_token"
}
