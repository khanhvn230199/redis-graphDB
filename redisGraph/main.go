package main

import (
	"fmt"
	redisgraph "redis/redisGraph"
	"redis/user"
)

func main() {
	users := user.ListUser()
	fmt.Printf("userID := %d\n", len(users))

	// for này cho chạy bằng cơm
	// for i := 0; i < len(users)/100000; i++ {
	// 	redisgraph.ConnecRedisGraph(users[i].UserID)
	// }

	redisgraph.ConnecRedisGraph(777132)

}
