package models

import "gorm.io/gorm"

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string
	Email    string
}

// Migrate user table
func MigrateUsers(db *gorm.DB) {
	db.AutoMigrate(&User{})
}
