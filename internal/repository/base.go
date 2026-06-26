package repository

type Base struct {
	ID int `gorm:"primarykey;autoIncrement" json:"id"`
}
