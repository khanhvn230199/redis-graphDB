package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type User struct {
	UserID int32 `json:"id"`
}

type Users struct {
	Users []User `json:"users"`
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
