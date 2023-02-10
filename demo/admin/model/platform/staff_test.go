package model

import (
	"fmt"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/demo/model/db/admin"
	adminsvr "github.com/xhigher/hzgo/server/admin"
	"github.com/xhigher/hzgo/types"
	"github.com/xhigher/hzgo/utils"
	"testing"
)

func TestStaff(t *testing.T) {
	uid := createUserid()
	password := randomPassword()
	password = "99637572"
	ct := utils.NowTime()
	fmt.Println("password=", password)
	staffInfo := &admin.StaffInfoModel{
		Uid:   uid,
		Username: "admin",
		Password: utils.MD5(fmt.Sprintf("%s-%d", utils.MD5(password), ct)),
		Nickname: "admin",
		Avatar:   "",
		Phone: "13528867472",
		Email: "xhigher@qq.com",
		Roles: types.StringArray{adminsvr.RoleMaintainer},
		Status:   consts.UserStatusActive,
		Ut: ct,
		Ct:       ct,
	}
	sql := "INSERT INTO `%s` (`uid`,`username`,`password`,`nickname`,`avatar`,`phone`,`email`,`roles`,`status`,`ut`,`ct`) VALUES "
	sql = sql + "('%s','%s','%s','%s','%s','%s','%s','%s',%d,%d,%d)"
	sql = fmt.Sprintf(sql, staffInfo.TableName(), staffInfo.Uid, staffInfo.Username, staffInfo.Password, staffInfo.Nickname, staffInfo.Avatar,
		staffInfo.Phone, staffInfo.Email, utils.JSONString(staffInfo.Roles), staffInfo.Status, staffInfo.Ut, staffInfo.Ct)
	fmt.Println(sql)

}
