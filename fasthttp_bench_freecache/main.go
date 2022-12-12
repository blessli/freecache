package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/coocood/freecache"
	"github.com/valyala/fasthttp"
	_ "net/http/pprof"
)

var cacheInstance *freecache.Cache
var cacheValue []byte
func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	ss := ""
	for i := 0; i < 1e2; i++ {
		str := uuid.New().String()
		ss += str
	}
	cacheValue = []byte(ss)
	log.Println("cacheValue size: ",len(cacheValue))
	cacheInstance = freecache.NewCache(50 * 1024 * 1024)
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		if string(ctx.Path()) == "/ping" {
			rand.Seed(time.Now().Unix())
			key := fmt.Sprintf("solitest:%d", rand.Int()%(1e3-1))
			v, err := cacheInstance.Get([]byte(key))
			if err != nil&& err.Error()!=freecache.ErrNotFound.Error() {
				log.Println("freecache get error: ", err)
				ctx.Success("", []byte("pong"))
				return
			}
			if len(v) == 0{
				cacheInstance.Set([]byte(key), cacheValue, -1)
			}
			ctx.Success("", []byte("pong"))
			return
		}
	}

	if err := fasthttp.ListenAndServe(":8080", requestHandler); err != nil {
		log.Fatalf("error in ListenAndServe: %v", err)
	}
}
