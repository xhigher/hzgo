package model

import (
	"github.com/xhigher/hzgo/demo/model/db/admin"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/utils"
	"gorm.io/gorm"
)

func AddTraceLog(module, action string, params, result interface{}, roles []string, uid string) (err error) {
	model := &admin.TraceLogModel{
		Module: module,
		Action: action,
		Params: utils.JSONString(params),
		Result: utils.JSONString(result),
		Roles:  utils.JSONString(roles),
		Uid:    uid,
		Ts:     utils.NowTime(),
	}
	err = admin.DB().Create(model).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		return
	}
	return
}

func GetTraceLogs(uid, module string, offset, limit int32) (total int64, data []*admin.TraceLogModel, err error) {
	tx := admin.DB().Model(admin.TraceLogModel{})
	if len(uid) > 0 {
		tx.Where("uid = ?", uid)
	}
	if len(module) > 0 {
		tx.Where("module = ?", module)
	}
	tx = tx.Session(&gorm.Session{})
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
