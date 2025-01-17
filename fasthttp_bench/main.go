package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	cache "github.com/patrickmn/go-cache"
	"github.com/valyala/fasthttp"
	_ "net/http/pprof"
)

var cacheInstance *cache.Cache
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
	cacheInstance = cache.New(time.Minute*30, time.Minute*20)
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		if string(ctx.Path()) == "/ping" {
			rand.Seed(time.Now().Unix())
			key := fmt.Sprintf("solitest:%d", rand.Int()%(1e3-1))
			_, found := cacheInstance.Get(key)
			if found {
				ctx.Success("", []byte("pong"))
				return
			}
			cacheInstance.Set(key, cacheValue, -1)
			ctx.Success("", []byte("pong"))
			return
		}
	}

	if err := fasthttp.ListenAndServe(":8080", requestHandler); err != nil {
		log.Fatalf("error in ListenAndServe: %v", err)
	}
}
