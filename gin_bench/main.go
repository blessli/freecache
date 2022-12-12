package main

import (
  "net/http"

  "github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()//New()
  r.Use(gin.Recovery())
  r.GET("/ping", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })
  r.Run(":8080")
}
// wrk压测http：./wrk -t 20 -c 5000 -d 20 --latency http://localhost:8080/ping