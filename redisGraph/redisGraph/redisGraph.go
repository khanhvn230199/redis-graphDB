package redisgraph

import (
	"fmt"
	"test/user"
	"time"

	"github.com/gomodule/redigo/redis"
	rg "github.com/redislabs/redisgraph-go"
	"gorm.io/gorm"
)

type Server struct {
	Db *gorm.DB
}

func ConnecRedisGraph() redis.Pool {
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
	return *pool

}

func QueryGraph(userID int32) []user.UsersData {

	conn := ConnecRedisGraph()

	graph := rg.GraphNew("kingtalk", conn.Get())

	query := fmt.Sprintf(`MATCH (us {UserID:%d})-[i:Inbox]-(ud:User) RETURN i, ud.UserID`, userID)

	result, _ := graph.Query(query)

	m := make(map[string]bool)

	list := []user.UsersData{}

	if result != nil {
		for result.Next() {
			CountIB, _ := result.Record().Get("i")
			countIbox := CountIB.(*rg.Edge).Properties["Count"]
			// fmt.Printf("count := %d   \n", countIbox.(int))
			if countIbox == nil {
				continue
			}
			c := countIbox.(int)
			if c == 0 {
				continue
			}

			Participant, _ := result.Record().Get("ud.UserID")
			// fmt.Printf("ParticipantID  := %v , userID :=%d  \n", Participant, userID)
			p := Participant.(int)

			key := genKey(int32(p), userID)

			if exist, _ := m[key]; !exist {
				list = append(list, user.UsersData{
					UserID:      userID,
					Participant: int32(p),
					CountInbox:  int32(c),
				})
				m[key] = true
			}
		}

		return list

	}
	return list

}

func genKey(p int32, u int32) string {
	if p > u {
		p, u = u, p
	}
	return fmt.Sprintf("%d-%d", p, u)
}

func DeleteArray(list []user.UsersData, p int32) []user.UsersData {
	for i := 0; i < len(list)-1; i++ {
		if list[i].Participant == p {
			list = append(list[:i+1], list[i+2:]...)
		}
	}
	return list
}

func (s *Server) CreateDB(u, p, c int32) {
	err := user.CreateUser(s.Db, u, p, c)
	if err != nil {
		fmt.Println(err)
	}
}

func (s *Server) CreateMultiDB(u []user.UsersData) {
	err := user.CreateMultiUser(s.Db, u)
	if err != nil {
		fmt.Println(err)
	}
}

func DeleteNode(userID int32) {
	conn := ConnecRedisGraph()

	graph := rg.GraphNew("kingtalk", conn.Get())

	query := fmt.Sprintf(`MATCH (us {UserID:%d} DELETE us`, userID)

	result, _ := graph.Query(query)

	if result != nil {
		fmt.Println(result)
	}
}

func CreateNode(userID int32) {
	conn := ConnecRedisGraph()

	graph := rg.GraphNew("kingtalk", conn.Get())

	query := fmt.Sprintf(`MATCH (us:User {UserID: %d}) RETURN us.User`, userID)

	result, _ := graph.Query(query)

	fmt.Println("ngoai vong for")

	// query := fmt.Sprintf("CREATE  (us:User {UserID:%d})", userID)
	for !result.Next() {

		fmt.Println("trong vong for")
		query2 := fmt.Sprintf(`CREATE (us:User {UserID: %d})`, userID)

		graph.Query(query2)

		fmt.Println(graph.Query(query2))
	}

}

func CreateEdge(c, p, u int32) {
	conn := ConnecRedisGraph()

	if c != 0 {
		graph := rg.GraphNew("kingtalk", conn.Get())
		query := fmt.Sprintf(`MATCH (us:User {UserID: %d}), (ud:User {UserID: %d})
							  MERGE (us)-[i:Inbox {Count:%d}]->(ud)`, u, p, c)

		result, _ := graph.Query(query)
		if result != nil {
			fmt.Println(result)
		}
	}
}
