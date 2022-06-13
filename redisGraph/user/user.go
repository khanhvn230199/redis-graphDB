package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"gorm.io/gorm"
)

type User struct {
	UserID int32 `json:"id"`
}

type Users struct {
	Users []User `json:"users"`
}

type UsersData struct {
	ID          int32 `json:"id"`
	UserID      int32 `json:"user"`
	Participant int32 `json:"participant"`
	CountInbox  int32 `json:"count"`
}

func CreateUser(db *gorm.DB, u, p, c int32) (err error) {
	err = db.Create(UsersData{UserID: u, Participant: p, CountInbox: c}).Error
	if err != nil {
		return err
	}
	return nil
}

func CreateMultiUser(db *gorm.DB, u []UsersData) (err error) {
	err = db.Create(u).Error
	if err != nil {
		return err
	}
	return nil
}

func ListUser() []User {
	jsonFile, err := os.Open("users.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var users Users
	json.Unmarshal(byteValue, &users)

	return users.Users
}
