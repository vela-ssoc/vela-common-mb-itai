package model

import "time"

type ThirdDown struct {
	ID        int64     `json:"id"`         // 数据库 ID
	Name      string    `json:"name"`       // 三方文件名
	BrokerID  int64     `json:"broker_id"`  // broker_id
	MinionID  int64     `json:"minion_id"`  // minion_id
	CreatedAt time.Time `json:"created_at"` // 下载时间
}

func (ThirdDown) TableName() string {
	return "third_down"
}
