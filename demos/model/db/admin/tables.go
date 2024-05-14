package admin

import (
	"fmt"
	"github.com/xhigher/hzgo/demos/model/db"
	"github.com/xhigher/hzgo/types"
	"gorm.io/gorm"
	"time"
)

func DB() *gorm.DB {
	return db.AdminDB()
}

type StaffInfoModel struct {
	Uid      string            `json:"uid" xorm:"uid" gorm:"column:uid"`
	Username string            `json:"username" xorm:"username" gorm:"column:username"`
	Password string            `json:"password" xorm:"password" gorm:"column:password"`
	Nickname string            `json:"nickname" xorm:"nickname" gorm:"column:nickname"`
	Avatar   string            `json:"avatar" xorm:"avatar" gorm:"column:avatar"`
	Email    string            `json:"email" xorm:"email" gorm:"column:email"`
	Phone    string            `json:"phone" xorm:"phone" gorm:"column:phone"`
	Roles    types.StringArray `json:"roles" xorm:"roles" gorm:"column:roles"`
	Status   int32             `json:"status" xorm:"status" gorm:"column:status"`
	Ct       int64             `json:"ct" xorm:"ct" gorm:"column:ct"`
	Ut       int64             `json:"ut" xorm:"ut" gorm:"column:ut"`
}

func (StaffInfoModel) TableName() string {
	return "staff_info"
}

type StaffTokenModel struct {
	Uid   string `json:"uid" gorm:"column:uid"`
	Token string `json:"token" gorm:"column:token"`
	Et    int64  `json:"et" gorm:"column:et"`
	It    int64  `json:"it" gorm:"column:it"`
	Ut    int64  `json:"ut" gorm:"column:ut"`
}

func (t *StaffTokenModel) TableName() string {
	return "staff_token"
}

type TraceLogModel struct {
	Id     int64  `json:"id" gorm:"column:id"`
	Module string `json:"module" gorm:"column:module"`
	Action string `json:"action" gorm:"column:action"`
	Params string `json:"params" gorm:"column:params"`
	Result string `json:"result" gorm:"column:result"`
	Roles  string `json:"roles" gorm:"column:roles"`
	Uid    string `json:"uid" gorm:"column:uid"`
	Ts     int64  `json:"ts" gorm:"column:ts"`
}

func (TraceLogModel) TableName() string {
	return fmt.Sprintf("trace_log_%d", time.Now().Year())
}
