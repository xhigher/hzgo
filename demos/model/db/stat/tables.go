package stat

import (
	"github.com/xhigher/hzgo/demo/model/db"
	"gorm.io/gorm"
)

func DB() *gorm.DB {
	return db.StatDB()
}
