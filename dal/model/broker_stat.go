package model

import "time"

type BrokerStat struct {
	ID         int64     `json:"id,string"   gorm:"column:id;primaryKey"`
	Name       string    `json:"name"        gorm:"name"`
	MemUsed    uint64    `json:"mem_used"    gorm:"mem_used"`
	MemTotal   uint64    `json:"mem_total"   gorm:"mem_total"`
	CPUPercent float64   `json:"cpu_percent" gorm:"cpu_percent"`
	CreatedAt  time.Time `json:"created_at"  gorm:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"  gorm:"updated_at"`
}

func (BrokerStat) TableName() string {
	return "broker_stat"
}
