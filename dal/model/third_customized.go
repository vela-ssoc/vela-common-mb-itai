package model

import "time"

type ThirdCustomized struct {
	ID        int64     `json:"id,string"  gorm:"column:id;primaryKey"`
	Name      string    `json:"name"       gorm:"column:name"`
	Icon      string    `json:"icon"       gorm:"column:icon"`
	Remark    string    `json:"remark"     gorm:"column:remark"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

func (ThirdCustomized) TableName() string {
	return "third_customized"
}
