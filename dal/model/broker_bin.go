package model

import "time"

// BrokerBin 存放 broker 节点的可执行程序。
type BrokerBin struct {
	ID        int64     `json:"id,string"  gorm:"column:id;primaryKey"`
	Name      string    `json:"name"       gorm:"column:name"`
	FileID    int64     `json:"-"          gorm:"column:file_id"`
	Size      int64     `json:"size"       gorm:"column:size"`
	Hash      string    `json:"hash"       gorm:"column:hash"`
	Goos      string    `json:"goos"       gorm:"column:goos"`
	Arch      string    `json:"arch"       gorm:"column:arch"`
	Semver    Semver    `json:"semver"     gorm:"column:semver"`
	Changelog string    `json:"changelog"  gorm:"column:changelog"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (BrokerBin) TableName() string {
	return "broker_bin"
}
