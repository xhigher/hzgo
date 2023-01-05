package user

import (
	"fmt"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/mysql"
	"github.com/xhigher/hzgo/utils"
	"gorm.io/gorm"
	"time"
)

func createUserid() string {
	return utils.IntToBase36(utils.NowTimeMicro() - 1040000000000000)
}

type CreateUserLogic struct {
	Username string
	Password string
	userInfo *UserInfo
	existed  bool
}

func (logic CreateUserLogic) Do() (userInfo *UserInfo, existed bool, err error) {
	err = DB().Transaction(logic.getTransaction)
	if err != nil {
		logger.Errorf("transaction error %v ", err)
		return
	}
	userInfo = logic.userInfo
	existed = logic.existed
	return
}

func (logic CreateUserLogic) getTransaction(tx *gorm.DB) (err error) {
	err = tx.Where("username = ?", logic.Username).First(&logic.userInfo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logic.userInfo = nil
			err = nil
		} else {
			return
		}
	}
	if logic.userInfo != nil {
		logic.existed = true
		return
	}
	ct := utils.NowTime()
	userid := createUserid()
	logic.userInfo = &UserInfo{
		Userid:   userid,
		Username: logic.Username,
		Password: utils.MD5(fmt.Sprintf("%s-%d", logic.Password, ct)),
		Nid:      "",
		Nickname: "",
		Avatar:   "",
		Sex:      0,
		Birthday: "",
		Inviter:  "",
		Status:   consts.UserStatusNormal,
		Ct:       ct,
	}
	err = tx.Create(logic.userInfo).Error
	if err != nil {
		if mysql.ErrDuplicateKey(err) {
			time.Sleep(time.Millisecond * 10)
			logic.userInfo.Userid = createUserid()
			err = tx.Create(logic.userInfo).Error
			if err != nil {
				logic.userInfo = nil
				return
			}
		} else {
			logic.userInfo = nil
		}
	}
	return
}

func GetUser(username string) (data *UserInfo, err error) {
	err = DB().Where("username = ?", username).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			data = nil
			err = nil
			return
		}
	}
	return
}

func CheckPassword(data *UserInfo, password string) bool {
	password = utils.MD5(fmt.Sprintf("%s-%d", password, data.Ct))
	if password == data.Password {
		return true
	}
	return false
}

func GetUserById(userid string) (data *UserInfo, err error) {
	err = DB().Where("userid = ?", userid).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			data = nil
			err = nil
			return
		}
	}
	return
}

func GetUserByIcode(icode string) (data *UserInfo, err error) {
	err = DB().Where("icode = ?", icode).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			data = nil
			err = nil
			return
		}
	}
	return
}
