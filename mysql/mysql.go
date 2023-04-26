package mysql

import (
	basemysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
	"time"
	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/logger"
)

var (
	gormDBs map[string]*gorm.DB
)

func Init(configs []*config.MysqlConfig) {
	gormDBs = make(map[string]*gorm.DB)
	if len(configs) == 0 {
		logger.Warn("no mysql configs")
		return
	}
	for _, config := range configs {
		mysqlConfig := mysql.Config{
			DSN:                       config.Dsn(), // DSN data source name
			DefaultStringSize:         191,          // string 类型字段的默认长度
			DisableDatetimePrecision:  true,         // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,         // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,         // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false,        // 根据版本自动配置
		}

		if db, err := gorm.Open(mysql.New(mysqlConfig), gormOptions(config.DbName)); err == nil {
			sqlDB, _ := db.DB()
			if config.MaxIdleConns > 0 {
				sqlDB.SetMaxIdleConns(config.MaxIdleConns)
			}
			if config.MaxOpenConns > 0 {
				sqlDB.SetMaxOpenConns(config.MaxOpenConns)
			}
			sqlDB.SetConnMaxIdleTime(time.Minute * 20)

			dbName := config.DbName
			if config.Standby {
				dbName = dbName + ":standby"
			}
			gormDBs[dbName] = db
			logger.Infof("mysql success, db-name: %s, dsn: %s", dbName, config.Dsn())
		} else {
			logger.Errorf("mysql failed, err: %v", err)
			return
		}
	}
}

func gormOptions(dbName string) *gorm.Config {
	return &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   newGormLogger(dbName),
	}
}

func DB(dbName string) *gorm.DB {
	if db, ok := gormDBs[dbName]; ok {
		return db
	}
	return nil
}

func StandbyDB(dbName string) *gorm.DB {
	if db, ok := gormDBs[dbName+":standby"]; ok {
		return db
	}
	return nil
}

func ErrNotFound(err error) bool {
	if err == gorm.ErrRecordNotFound {
		return true
	}
	return false
}

func ErrDuplicateKey(err error) bool {
	if sqlError, ok := err.(*basemysql.MySQLError); ok && sqlError.Number == 1062 {
		return true
	}
	return false
}

type LogPrinter struct {
	dbName string
}

func (p LogPrinter) Printf(f string, v ...interface{}) {
	logger.Infof("db: "+p.dbName+", "+f, v)
}

func newGormLogger(dbName string) gormlog.Interface {
	return gormlog.New(LogPrinter{
		dbName,
	}, gormlog.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  gormlog.Warn,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	})
}
