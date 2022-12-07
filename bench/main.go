package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	_ "net/http/pprof"

	"github.com/coocood/freecache"
)
var goroutineNums = flag.Int("gn", 2, "goroutine nums")
func main() {
	flag.Parse()
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
 
	rand.Seed(time.Now().Unix())
	// lc := cache.New(time.Minute*5, time.Minute*2)
	lc := freecache.NewCache(256 * 1024 * 1024) // 256MB
	log.Printf("start at:%v", time.Now())
	aaaKey := "aaa:%d:buy:cnt"
	log.Println("set run over")
 
	for i := 0; i < *goroutineNums; i++ {
		go func(idx int) {
			for {
				key := fmt.Sprintf(aaaKey, rand.Int())
				newKey := fmt.Sprintf("%s:%d", key, rand.Int())
				v := rand.Int()
				lc.Set([]byte(newKey), []byte(strconv.Itoa(v)), 10*int(time.Millisecond))
			}
		}(i)
	}
 
	// 查看 go-cache 中 key 的数量
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				log.Printf("lc key count:%d", lc.EntryCount())
			}
		}
	}()
	select {}
}