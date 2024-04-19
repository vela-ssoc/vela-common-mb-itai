package model

import "time"

type KVAudit struct {
	ID        int64     `json:"id,string"  gorm:"column:id"`
	MinionID  int64     `json:"minion_id"  gorm:"column:minion_id"`
	Inet      string    `json:"inet"       gorm:"column:inet"`
	Bucket    string    `json:"bucket"     gorm:"column:bucket"`
	Key       string    `json:"key"        gorm:"column:key"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (KVAudit) TableName() string {
	return "kv_audit"
}
