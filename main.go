package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func executeQuery(query string, schema graphql.Schema) {
} 

func main() {
  router := gin.Default()

  router.GET("/greet", Greeting)

  router.Run()
}

func Greeting(c *gin.Context) {
  query := c.Request.URL.Query().Get("query")
  fmt.Println("Query: ", query)
  fmt.Println("Greeting: ")
  c.IndentedJSON(http.StatusOK, "hi")
}
