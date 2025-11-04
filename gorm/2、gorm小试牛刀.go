package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID    int
	Name  string
	Email string
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm_new_db"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	var userList []User
	db.Find(&userList)
	fmt.Println(userList)
}
