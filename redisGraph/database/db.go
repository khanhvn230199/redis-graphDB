package database

import (
	"fmt"
	"test/user"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDb() *gorm.DB {
	Db = connectDB()
	return Db
}

func connectDB() *gorm.DB {
	var err error
	dsn := "root:123456789@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Printf("Error connecting to database : error=%v \n", err)
		return nil
	}

	db.AutoMigrate(&user.UsersData{})

	return db
}
