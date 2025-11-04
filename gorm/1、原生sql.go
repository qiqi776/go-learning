package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm_new_db"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("dsn格式错误", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("数据库连接失败", err)
	}
	fmt.Println(db)

	res, err := db.Exec(
		// 错误： "INSERT INTO users (name, age, email) VALUES (?, ?, ?)"
		// 修正：
		"INSERT INTO user (username, age, email) VALUES (?, ?, ?)",
		"王五", 25, "wangwu@email.com",
	)
	if err != nil {
		fmt.Printf("插入失败: %v\n", err)
		return
	}
	lastId, _ := res.LastInsertId()
	fmt.Printf("插入成功, ID: %d\n", lastId)

	rows, err := db.Query("select id, username from user")
	if err != nil {
		log.Fatal("查询失败", err)
	}
	defer rows.Close()

	fmt.Println("查询所有用户:")
	for rows.Next() {
		var id int
		var username string // 变量名也最好对应

		// 修正：
		err = rows.Scan(&id, &username)
		if err != nil {
			log.Printf("扫描单行数据失败: %v\n", err)
			continue
		}
		fmt.Printf("  ID: %d, Username: %s\n", id, username)
	}
	if err = rows.Err(); err != nil {
		log.Fatalf("rows 迭代出错: %v", err)
	}

	// 5. 查询单行 (修正：表名 user, 列名 username)
	var id int
	var username string // 变量名也最好对应

	// 错误： "SELECT id, name FROM users WHERE id = ?"
	// 修正：
	err = db.QueryRow("SELECT id, username FROM user WHERE id = ?", lastId).Scan(&id, &username)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("未找到指定 ID 的用户")
		} else {
			log.Fatalf("查询单行失败: %v", err)
		}
	} else {
		// 修正：打印值
		fmt.Printf("查询单个用户 (ID: %d): Username: %s\n", id, username)
	}
}
