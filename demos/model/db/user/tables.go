package user

import (
	"github.com/xhigher/hzgo/demos/model/db"
	"gorm.io/gorm"
)

func DB() *gorm.DB {
	return db.UserDB()
}

type UserInfoModel struct {
	Userid   string `json:"userid" gorm:"column:userid"`
	Username string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
	Nid      string `json:"nid" gorm:"column:nid"`
	Nickname string `json:"nickname" gorm:"column:nickname"`
	Avatar   string `json:"avatar" gorm:"column:avatar"`
	Sex      int32  `json:"sex" gorm:"column:sex"`
	Birthday string `json:"birthday" gorm:"column:birthday"`
	Inviter  string `json:"inviter" gorm:"column:inviter"`
	Status   int32  `json:"status" gorm:"column:status"`
	Cd       string `json:"cd" xorm:"cd" gorm:"column:cd"`
	Ct       int64  `json:"ct" gorm:"column:ct"`
	Ut       int64  `json:"ut" gorm:"column:ut"`
}

func (t *UserInfoModel) TableName() string {
	return "user_info"
}

type UserTokenModel struct {
	Userid string `json:"userid" gorm:"column:userid"`
	Token  string `json:"token" gorm:"column:token"`
	Et     int64  `json:"et" gorm:"column:et"`
	It     int64  `json:"it" gorm:"column:it"`
	Ut     int64  `json:"ut" gorm:"column:ut"`
}

func (t *UserTokenModel) TableName() string {
	return "user_token"
}
