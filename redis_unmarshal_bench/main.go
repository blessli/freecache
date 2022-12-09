package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"

	"github.com/coocood/freecache/utils"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/panjf2000/gnet"
)

var server *Server

type Server struct {
	redisClient *redis.Client
	cacheValue  []byte
	cacheString string
}
type Activitys []*Activity

type Activity struct {
	Field0 string `json:"field_0"`
	Field1 string `json:"field_1"`
	Field2 string `json:"field_2"`
	Field3 string `json:"field_3"`
	Field4 string `json:"field_4"`
	Field5 string `json:"field_5"`
	Field6 string `json:"field_6"`
	Field7 string `json:"field_7"`
	Field8 string `json:"field_8"`
	Field9 string `json:"field_9"`
}

func mockAct() *Activity {
	var act = &Activity{}
	for j := 0; j < 10; j++ {
		ss := ""
		for i := 0; i < 10; i++ {// 10-30组合就是109kb
			str := uuid.New().String()
			ss += str
		}
		switch j {
		case 0:
			act.Field0 = ss
		case 1:
			act.Field1 = ss
		case 2:
			act.Field2 = ss
		case 3:
			act.Field3 = ss
		case 4:
			act.Field4 = ss
		case 5:
			act.Field5 = ss
		case 6:
			act.Field6 = ss
		case 7:
			act.Field7 = ss
		case 8:
			act.Field8 = ss
		case 9:
			act.Field9 = ss
		}
	}
	return act
}
func NewServer(redisClient *redis.Client) (server *Server) {
	server = new(Server)
	server.redisClient = redisClient
	for e := 0; e < 1e3; e++ {
		acts := []*Activity{}
		for i := 0; i < 30; i++ {
			acts = append(acts, mockAct())
		}
		m, err := json.Marshal(acts)
		if err != nil {
			panic(err)
		}
		log.Println(len(string(m))) // 109kb
		server.redisClient.Set(fmt.Sprintf("solitest:%d", e), string(m), -1)
	}
	return
}

type echoServer struct {
	*gnet.EventServer
}

// redis 大key反序列化 cpu飙升问题复现
func (es *echoServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	out = frame
	// key := string(frame)
	rand.Seed(time.Now().Unix())
	key := fmt.Sprintf("solitest:%d", rand.Int()%(1e3-1))
	data, err := server.redisClient.Get(key).Result()
	if err != nil {
		log.Printf("redis get error: %s|%v\n", key, err)
		return
	}
	acts := []*Activity{}
	err = json.Unmarshal([]byte(data), &acts)
	if err != nil {
		log.Printf("react unmarshal error: %v\n", err)
		return
	}
	return
}

func init() {
	// 实例化RedisSingleObj结构体
	conn := &utils.RedisSingleObj{
		Redis_host: "0.0.0.0",
		Redis_port: 6379,
	}

	// 初始化连接 Single Redis 服务端
	err := conn.InitSingleRedis()
	if err != nil {
		panic(err)
	}
	server = NewServer(conn.Db)
}

func main() {

	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()
	// 查看 cache 中 key 的数量
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
			}
		}
	}()
	echo := new(echoServer)
	log.Fatal(gnet.Serve(echo, "tcp://:18000", gnet.WithMulticore(true)))
	exitOnSignal()
}

// exitOnSignal 监听退出信号
func exitOnSignal() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	select {
	case <-quit:
		// nop
	}
}
