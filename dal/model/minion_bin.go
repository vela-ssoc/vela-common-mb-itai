package model

import "time"

// MinionBin minion 节点二进制发行版本记录表
type MinionBin struct {
	ID         int64     `json:"id,string"  gorm:"column:id;primaryKey"` // 表 ID
	FileID     int64     `json:"-"          gorm:"column:file_id"`       // 关联文件表 ID
	Goos       string    `json:"goos"       gorm:"column:goos"`          // 操作系统
	Arch       string    `json:"arch"       gorm:"column:arch"`          // 系统架构
	Name       string    `json:"name"       gorm:"column:name"`          // 文件名称
	Customized string    `json:"customized" gorm:"column:customized"`    // 定制版标记
	Unstable   bool      `json:"unstable"   gorm:"column:unstable"`      // 不稳定版本，内测版本
	Caution    string    `json:"caution"    gorm:"column:caution"`       // 注意事项
	Ability    string    `json:"ability"    gorm:"column:ability"`       // 功能说明
	Size       int64     `json:"size"       gorm:"column:size"`          // 文件大小
	Hash       string    `json:"hash"       gorm:"column:hash"`          // 文件哈希
	Semver     Semver    `json:"semver"     gorm:"column:semver"`        // 版本号
	Changelog  string    `json:"changelog"  gorm:"column:changelog"`     // 更新日志
	Weight     int64     `json:"-"          gorm:"column:weight"`        // 版本号权重，用于比较版本号大小
	Deprecated bool      `json:"deprecated" gorm:"column:deprecated"`    // 是否已过期
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// TableName implement gorm schema.Tabler
func (MinionBin) TableName() string {
	return "minion_bin"
}
