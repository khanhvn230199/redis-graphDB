package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	redisgraph "test/redisGraph"
	"test/user"

	"github.com/panjf2000/ants/v2"
)

const lengthFile = 100000

var saveFile chan []user.UsersData
var userCount int32

func saveFileFunc() {
	i := int32(1)
	lUserData := make([]user.UsersData, 0)
	for {
		lUserData = append(lUserData, <-saveFile...)

		if i%lengthFile == 0 || i == userCount {
			b, _ := json.Marshal(lUserData)
			ioutil.WriteFile(fmt.Sprintf("list-%d.json", i), b, 0777)
			lUserData = make([]user.UsersData, 0)
		}

		if i == userCount {
			panic("done")
		}
		fmt.Println(i)
		i++
	}
}

func demoFunc(i interface{}) {
	saveFile <- redisgraph.QueryGraph(i.(int32))
}

func main() {
	// fmt.Printf("running goroutines: %d\n", p.Running())
	defer ants.Release()

	users := user.ListUser()
	userCount = int32(len(users))
	saveFile = make(chan []user.UsersData, 1)
	// runTimes := 1000

	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
		demoFunc(i)
	})
	defer p.Release()
	// Submit tasks one by one.
	go func(users []user.User) {
		for _, u := range users {
			_ = p.Invoke(u.UserID)
		}
	}(users)

	saveFileFunc()

}
