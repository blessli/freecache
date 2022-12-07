package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/coocood/freecache"
	"github.com/panjf2000/gnet"
	"github.com/google/uuid"
)

var server *Server

type Server struct {
	cache *freecache.Cache
	cacheValue []byte
}

func NewServer(cacheSize int) (server *Server) {
	server = new(Server)
	server.cache = freecache.NewCache(cacheSize)
	for i:=0;i<1e3;i++ {
		server.cacheValue = append(server.cacheValue, []byte(uuid.New().String()))
	}
	log.Println("cache value size: ", len(server.cacheValue))
	return
}

type echoServer struct {
	*gnet.EventServer
}

func (es *echoServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	out = frame
	server.cache.Set(frame, frame, 100*int(time.Millisecond))
	return
}

func main() {
	server = NewServer(256 * 1024 * 1024)
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	echo := new(echoServer)
	log.Fatal(gnet.Serve(echo, "tcp://:18000", gnet.WithMulticore(false)))
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
