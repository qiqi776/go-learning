package models

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type UserModel struct {
	ID        int
	Name      string `gorm:"size:16;not null;unique"`
	Age       int    `gorm:"default:18"`
	CreatedAt time.Time
}

func (u *UserModel) BeforeCreate(tx *gorm.DB) error {
	fmt.Println("before create user")
	return nil
}
