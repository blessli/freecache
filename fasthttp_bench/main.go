package main

import (
	"log"

	"github.com/valyala/fasthttp"
)

func main() {
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		if string(ctx.Path()) == "/ping" {
			ctx.Success("",[]byte("pong"))
			return
		}
	}

	if err := fasthttp.ListenAndServe(":8080", requestHandler); err != nil {
		log.Fatalf("error in ListenAndServe: %v", err)
	}
}