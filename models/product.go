package models

type Products struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"unique"`
	Price    string
	Category string
}
