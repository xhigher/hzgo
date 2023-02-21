package model

import (
	"fmt"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/demo/model/db/admin"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/mysql"
	"github.com/xhigher/hzgo/types"
	"github.com/xhigher/hzgo/utils"
	"gorm.io/gorm"
	"time"
)

func createUserid() string {
	return utils.IntToBase36(utils.NowTime() - 1050000000)
}

func randomPassword() string {
	return utils.RandNumberString(8)
}

type CreateStaffTask struct {
	Username string
	Nickname string
	Phone string
	Email string
	staffInfo *admin.StaffInfoModel
	existed  bool
}

func (task *CreateStaffTask) Do() (staffInfo *admin.StaffInfoModel, existed bool, err error) {
	err = admin.DB().Transaction(task.getTransaction)
	if err != nil {
		logger.Errorf("transaction error %v ", err)
		return
	}
	staffInfo = task.staffInfo
	existed = task.existed
	return
}

func (task *CreateStaffTask) getTransaction(tx *gorm.DB) (err error) {
	err = admin.DB().Where("username = ?", task.Username).First(&task.staffInfo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			task.staffInfo = nil
			err = nil
		} else {
			return
		}
	}
	if task.staffInfo != nil {
		task.existed = true
		return
	}
	ct := utils.NowTime()
	uid := createUserid()
	password := randomPassword()
	task.staffInfo = &admin.StaffInfoModel{
		Uid:   uid,
		Username: task.Username,
		Password: utils.MD5(fmt.Sprintf("%s-%d", utils.MD5(password), ct)),
		Nickname: task.Username,
		Avatar:   "",
		Phone: task.Phone,
		Email: task.Email,
		Roles: types.StringArray{},
		Status:   consts.UserStatusActive,
		Ut:       ct,
		Ct:       ct,
	}
	err = tx.Create(task.staffInfo).Error
	if err != nil {
		if mysql.ErrDuplicateKey(err) {
			time.Sleep(time.Second * 1)
			task.staffInfo.Uid = createUserid()
			err = tx.Create(task.staffInfo).Error
			if err != nil {
				task.staffInfo = nil
				return
			}
		} else {
			task.staffInfo = nil
		}
	}
	return
}

func GetStaff(username string) (data *admin.StaffInfoModel, err error) {
	err = admin.DB().Where("username = ?", username).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			data = nil
			err = nil
			return
		}
		logger.Errorf("error: %v", err)
	}
	return
}

func GetStaffByUid(uid string) (data *admin.StaffInfoModel, err error) {
	err = admin.DB().Where("uid = ?", uid).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			data = nil
			err = nil
			return
		}
		logger.Errorf("error: %v", err)
	}
	return
}

func GetStaffList(status, offset, limit int32) (total int64, data []*admin.StaffInfoModel, err error) {
	tx := admin.DB().Model(admin.StaffInfoModel{}).Where("status = ?", status).Session(&gorm.Session{})
	err = tx.Count(&total).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	if total == 0 {
		return
	}

	err = tx.Offset(int(offset)).Limit(int(limit)).Find(&data).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	return
}

func CheckPassword(data *admin.StaffInfoModel, password string) bool {
	password = utils.MD5(fmt.Sprintf("%s-%d", password, data.Ct))
	if password == data.Password {
		return true
	}
	return false
}

func UpdateStaffRoles(uid string, roles types.StringArray) (err error){
	updates := map[string]interface{}{
		"roles": roles,
		"ut": utils.NowTime(),
	}
	err = admin.DB().Model(admin.StaffInfoModel{}).Where("uid = ?", uid).Updates(&updates).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	return
}

func UpdateStaffStatus(uid string, status int32) (err error){
	updates := map[string]interface{}{
		"status": status,
		"ut": utils.NowTime(),
	}
	tx := admin.DB().Model(admin.StaffInfoModel{}).Where("uid = ?", uid)
	if status == consts.UserStatusActive {
		tx.Where("status=?", consts.UserStatusBlocked)
	}else if status == consts.UserStatusBlocked {
		tx.Where("status=?", consts.UserStatusActive)
	}else{
		err = fmt.Errorf("status[%d] error", status)
		return
	}
	err = tx.Updates(&updates).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	return
}

func ResetStaffPassword(uid string, ct int64) (password string, err error){
	password = randomPassword()
	updates := map[string]interface{}{
		"password": utils.MD5(fmt.Sprintf("%s-%d", utils.MD5(password), ct)),
		"ut": utils.NowTime(),
	}
	err = admin.DB().Model(admin.StaffInfoModel{}).Where("uid = ?", uid).Updates(&updates).Error
	if err != nil {
		password = ""
		logger.Errorf("error: %v", err)
		return
	}
	return
}
