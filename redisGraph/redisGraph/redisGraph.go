package redisgraph

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	rg "github.com/redislabs/redisgraph-go"
)

func ConnecRedisGraph(userID int32) {
	// mapUserID := make(map[int32]int32)
	t := time.Second
	pool := &redis.Pool{
		MaxIdle:   2000,
		MaxActive: 2000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", "172.26.37.31:6388",
				redis.DialPassword("4JAsQpqCEBjmN5tr"),
				redis.DialConnectTimeout(t),
				redis.DialReadTimeout(t),
				redis.DialWriteTimeout(t),
			)
			return conn, err
		},
	}

	conn := pool.Get()

	graph := rg.GraphNew("kingtalk", conn)

	query := fmt.Sprintf(`MATCH (us {UserID:%d})-[i:Inbox]-(ud:User) RETURN i,ud,COUNT(us) AS data`, userID)

	result, _ := graph.Query(query)

	if result != nil {
		for result.Next() {
			record, _ := result.Record().Get("data")
			fmt.Printf("data := %v , userID :=%d  \n", record, userID)

			CountIB, _ := result.Record().Get("i")
			countIbox := CountIB.(*rg.Edge).Properties["Count"]
			fmt.Printf("count := %v , userID :=%d  \n", countIbox, userID)

			Participant, _ := result.Record().Get("ud")
			fmt.Printf("ParticipantID  := %v , userID :=%d  \n", Participant, userID)
		}

	}
}
