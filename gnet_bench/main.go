package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"

	"github.com/coocood/freecache"
	"github.com/google/uuid"
	"github.com/panjf2000/gnet"
)

var server *Server

type Server struct {
	cache      *freecache.Cache
	cacheValue []byte
}

func NewServer(cacheSize int) (server *Server) {
	server = new(Server)
	server.cache = freecache.NewCache(cacheSize)
	ss := ""
	for i := 0; i < 1e2; i++ {
		str := uuid.New().String()
		ss += str
	}
	server.cacheValue = []byte(ss)
	log.Println("cache value size: ", len(server.cacheValue)) // 3.5kb
	return
}

type echoServer struct {
	*gnet.EventServer
}

// 只有set接口
// func (es *echoServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
// 	out = frame
// 	err:=server.cache.Set(frame, server.cacheValue, 100*int(time.Millisecond))
// 	if err!=nil{
// 		log.Println("freecache set error: ",err)
// 	}
// 	return
// }

func (es *echoServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	out = frame
	v, err := server.cache.Get(frame)
	if err != nil&& err.Error()!=freecache.ErrNotFound.Error() {
		log.Println("freecache get error: ", err)
		return
	}
	if len(v) > 0 {
		log.Println(string(frame),"exists")
		return
	}
	err = server.cache.Set(frame, server.cacheValue, 100*int(time.Millisecond))
	if err != nil {
		log.Println("freecache set error: ", err)
	}
	return
}

func main() {
	server = NewServer(1024 * 1024 * 1024)
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()
	// 查看 go-cache 中 key 的数量
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				log.Printf("freecache key count:%d", server.cache.EntryCount())
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
