package model

import (
	"fmt"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/demos/model/db/user"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/mysql"
	"github.com/xhigher/hzgo/utils"
	"gorm.io/gorm"
	"time"
)

func createUserid() string {
	return utils.IntToBase36(utils.NowTimeMicro() - 1040000000000000)
}

type CreateUserTask struct {
	Username string
	Password string
	userInfo *user.UserInfoModel
	existed  bool
}

func (task *CreateUserTask) Do() (userInfo *user.UserInfoModel, existed bool, err error) {
	err = user.DB().Transaction(task.getTransaction)
	if err != nil {
		logger.Errorf("error %v ", err)
		return
	}
	userInfo = task.userInfo
	existed = task.existed
	return
}

func (task *CreateUserTask) getTransaction(tx *gorm.DB) (err error) {
	err = tx.First(&task.userInfo, "username = ?", task.Username).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			task.userInfo = nil
			err = nil
		} else {
			logger.Errorf("error: %v", err)
			return
		}
	}
	if task.userInfo != nil {
		task.existed = true
		return
	}
	ct := utils.NowTime()
	userid := createUserid()
	task.userInfo = &user.UserInfoModel{
		Userid:   userid,
		Username: task.Username,
		Password: utils.MD5(fmt.Sprintf("%s-%d", task.Password, ct)),
		Nid:      "",
		Nickname: "",
		Avatar:   "",
		Sex:      0,
		Birthday: "",
		Inviter:  "",
		Status:   consts.UserStatusActive,
		Ct:       ct,
	}
	err = tx.Create(task.userInfo).Error
	if err != nil {
		if mysql.ErrDuplicateKey(err) {
			time.Sleep(time.Millisecond * 10)
			task.userInfo.Userid = createUserid()
			err = tx.Create(task.userInfo).Error
			if err != nil {
				logger.Errorf("error: %v", err)
				task.userInfo = nil
				return
			}
		} else {
			logger.Errorf("error: %v", err)
			task.userInfo = nil
		}
	}
	return
}

func GetUser(username string) (data *user.UserInfoModel, err error) {
	err = user.DB().First(&data, "username = ?", username).Error
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

func CheckPassword(data *user.UserInfoModel, password string) bool {
	password = utils.MD5(fmt.Sprintf("%s-%d", password, data.Ct))
	if password == data.Password {
		return true
	}
	return false
}

func GetUserById(userid string) (data *user.UserInfoModel, err error) {
	err = user.DB().First(&data, "userid = ?", userid).Error
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
