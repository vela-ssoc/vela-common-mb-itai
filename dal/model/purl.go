package model

type Purl struct {
	ID string `json:"id" gorm:"column:id;primaryKey"`
}

func (p Purl) TableName() string {
	return "purl"
}
