package main

import (
  "fmt"
  "net/http"
  "github.com/gin-gonic/gin"
)

func main() {
  router := gin.Default()

  router.GET("/greet", Greeting)
  
  router.Run()
}

func Greeting(c *gin.Context) {
  fmt.Println("Greeting: ")
  c.IndentedJSON(http.StatusOK, "hi")
}
