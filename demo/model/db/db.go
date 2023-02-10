package db

import (
	"github.com/xhigher/hzgo/mysql"
	"gorm.io/gorm"
)

const (
	namePrefix = "hertz_"

	nameAdmin = namePrefix + "admin"

	nameMisc = namePrefix + "misc"
	nameUser = namePrefix + "user"
	nameActivity = namePrefix + "activity"
	nameStat = namePrefix + "stat"
)

func AdminDB() *gorm.DB {
	return mysql.DB(nameAdmin)
}

func MiscDB() *gorm.DB {
	return mysql.DB(nameMisc)
}

func MiscStandbyDB() *gorm.DB {
	return mysql.StandbyDB(nameMisc)
}

func UserDB() *gorm.DB {
	return mysql.DB(nameUser)
}

func UserStandbyDB() *gorm.DB {
	return mysql.StandbyDB(nameUser)
}

func ActivityDB() *gorm.DB {
	return mysql.DB(nameActivity)
}

func ActivityStandbyDB() *gorm.DB {
	return mysql.StandbyDB(nameActivity)
}

func StatDB() *gorm.DB {
	return mysql.DB(nameStat)
}

func StatStandbyDB() *gorm.DB {
	return mysql.StandbyDB(nameStat)
}