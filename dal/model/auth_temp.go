package model

import "time"

type AuthTemp struct {
	ID        int64     `json:"id"         gorm:"column:id;primaryKey"`
	UID       string    `json:"uid"        gorm:"column:uid"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

func (AuthTemp) TableName() string {
	return "auth_temp"
}

func (t AuthTemp) Expired(current time.Time, du time.Duration) bool {
	return t.CreatedAt.Add(du).Before(current)
}
