package misc

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/xhigher/hzgo/demo/model/db"
	"gorm.io/gorm"
)

func DB() *gorm.DB {
	return db.MiscDB()
}

type ConfigInfoModel struct {
	Id      string `json:"id" gorm:"column:id"`
	Name    string `json:"name" gorm:"column:name"`
	Data    string `json:"data" gorm:"column:data"`
	Filters string `json:"filters" gorm:"column:filters"`
	Ut      int64  `json:"ut" gorm:"column:ut"`
}

func (t *ConfigInfoModel) TableName() string {
	return "config_info"
}

type ConfigItem struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

type ConfigData []ConfigItem

func (i *ConfigData) Scan(v interface{}) error {
	data, ok := v.([]uint8)
	if !ok {
		return fmt.Errorf("error v type not []uint8")
	}
	if len(data) == 0 {
		return nil
	}
	var d ConfigData
	err := json.Unmarshal(data, &d)
	if err != nil {
		return err
	}
	*i = d
	return nil
}

func (i ConfigData) Value() (driver.Value, error) {
	data, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type BannerInfoModel struct {
	Id     int64  `json:"id" gorm:"column:id"`
	Site   string `json:"site" gorm:"column:site"`
	Type   int32  `json:"type" gorm:"column:type"`
	Name   string `json:"name" gorm:"column:name"`
	Img    string `json:"img" gorm:"column:img"`
	Data   string `json:"data" gorm:"column:data"`
	Sn     int64  `json:"sn" gorm:"column:sn"`
	Status int32  `json:"status" gorm:"column:status"`
	Ct     int64  `json:"ct" gorm:"column:ct"`
	Ut     int64  `json:"ut" gorm:"column:ut"`
}

func (BannerInfoModel) TableName() string {
	return "banner_info"
}
