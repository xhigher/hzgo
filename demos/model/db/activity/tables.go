package activity

import (
	"github.com/xhigher/hzgo/demos/model/db"
	"gorm.io/gorm"
)

func DB() *gorm.DB {
	return db.ActivityDB()
}
