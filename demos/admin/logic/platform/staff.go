package platform

import (
	"github.com/xhigher/hzgo/bizerr"
	"github.com/xhigher/hzgo/consts"
	model "github.com/xhigher/hzgo/demo/admin/model/platform"
	"github.com/xhigher/hzgo/demo/model/db/admin"
	"github.com/xhigher/hzgo/types"
)

func CheckStaff(username, password string) (uid string, roles []string, be *bizerr.Error) {
	staffInfo, err := model.GetStaff(username)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	if staffInfo == nil {
		be = bizerr.NotFound("用户不存在")
		return
	}
	if !model.CheckPassword(staffInfo, password) {
		be = bizerr.PasswordWrong("")
		return
	}
	uid = staffInfo.Uid
	roles = staffInfo.Roles
	return
}

func GetStaff(uid string) (staffInfo *admin.StaffInfoModel, be *bizerr.Error) {
	staffInfo, err := model.GetStaffByUid(uid)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	if staffInfo == nil {
		be = bizerr.NotFound("用户不存在")
		return
	}
	staffInfo.Password = ""
	return
}

func GetStaffList(status, offset, limit int32) (total int64, staffList []*admin.StaffInfoModel, be *bizerr.Error) {
	total, staffList, err := model.GetStaffList(status, offset, limit)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	return
}

func CreateStaff(username, nickname, phone, email string) (be *bizerr.Error) {
	task := model.CreateStaffTask{
		Username: username,
		Nickname: nickname,
		Phone:    phone,
		Email:    email,
	}
	_, existed, err := task.Do()
	if err != nil {
		be = bizerr.New(err)
		return
	}
	if existed {
		be = bizerr.AlreadyExists("用户已存在")
		return
	}
	return
}

func UpdateStaffRoles(uid string, roles types.StringArray) (be *bizerr.Error) {
	staffInfo, err := model.GetStaffByUid(uid)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	if staffInfo == nil {
		be = bizerr.NotFound("用户不存在")
		return
	}

	err = model.UpdateStaffRoles(uid, roles)
	if err != nil {
		be = bizerr.New(err)
		return
	}

	if staffInfo.Status == consts.StatusOnline {
		CleanStaffToken(uid)
	}

	return
}

func ResetStaffPassword(uid string) (password string, be *bizerr.Error) {
	staffInfo, err := model.GetStaffByUid(uid)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	if staffInfo == nil {
		be = bizerr.NotFound("用户不存在")
		return
	}
	password, err = model.ResetStaffPassword(uid, staffInfo.Ct)
	if err != nil {
		be = bizerr.New(err)
		return
	}

	if staffInfo.Status == consts.StatusOnline {
		CleanStaffToken(uid)
	}

	return
}

func UpdateStaffStatus(uid string, status int32) (be *bizerr.Error) {
	staffInfo, err := model.GetStaffByUid(uid)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	if staffInfo == nil {
		be = bizerr.NotFound("用户不存在")
		return
	}
	if staffInfo.Status == status {
		return
	}

	err = model.UpdateStaffStatus(uid, status)
	if err != nil {
		be = bizerr.New(err)
		return
	}

	if status == consts.StatusOffline {
		CleanStaffToken(uid)
	}

	return
}

func CleanStaffToken(uid string) (be *bizerr.Error) {
	err := model.SaveToken(uid, "", 0, 0)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	return
}

func GetTraceLogs(uid, module string, offset, limit int32) (total int64, logs []*admin.TraceLogModel, be *bizerr.Error) {
	total, logs, err := model.GetTraceLogs(uid, module, offset, limit)
	if err != nil {
		be = bizerr.New(err)
		return
	}
	return
}
