package mysql

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/alphafast/asmt-fw/libs/domain/noti/model"
)

type MySqlNotiResult struct {
	ID        string `gorm:"column:id;primaryKey"`
	ReqID     string `gorm:"column:req_id"`
	IsSuccess bool   `gorm:"column:is_success"`
	Reason    string `gorm:"column:reason"`
}

func (m *MySqlNotiResult) TableName() string {
	return "noti_result"
}

func ToNotiResult(mysqlNotiResult *MySqlNotiResult) *model.NotiResult {
	return &model.NotiResult{
		ID:        mysqlNotiResult.ID,
		ReqID:     mysqlNotiResult.ReqID,
		IsSuccess: mysqlNotiResult.IsSuccess,
		Reason:    mysqlNotiResult.Reason,
	}
}

func ToMySqlNotiResult(notiResult *model.NotiResult) *MySqlNotiResult {
	return &MySqlNotiResult{
		ID:        notiResult.ID,
		ReqID:     notiResult.ReqID,
		IsSuccess: notiResult.IsSuccess,
		Reason:    notiResult.Reason,
	}
}

type MySqlNotiUserNotiChannels []model.NotiUserNotiChannel

func (m MySqlNotiUserNotiChannels) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *MySqlNotiUserNotiChannels) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), m)
}

type MySqlNotiUser struct {
	UserID   string                    `gorm:"column:id;primaryKey"`
	Channels MySqlNotiUserNotiChannels `gorm:"column:channels;type:json"`
}

func (m *MySqlNotiUser) TableName() string {
	return "noti_user"
}

func ToNotiUser(mysqlNotiUser *MySqlNotiUser) *model.NotiUser {
	return &model.NotiUser{
		UserID:   mysqlNotiUser.UserID,
		Channels: mysqlNotiUser.Channels,
	}
}

func ToMySqlNotiUser(notiUser *model.NotiUser) *MySqlNotiUser {
	return &MySqlNotiUser{
		UserID:   notiUser.UserID,
		Channels: notiUser.Channels,
	}
}
