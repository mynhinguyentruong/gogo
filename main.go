package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/mynhinguyentruong/gogo/schema"
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
  result := graphql.Do(graphql.Params{
    Schema: schema,
    RequestString: query,
  })
  
  if result.HasErrors() {
    fmt.Printf("Wrong result, unexpected error: %v", result.Errors)
  }

  return result
} 

func main() {
  router := gin.Default()

  router.GET("/greet", Greeting)
  router.GET("/list", func (c *gin.Context) {
    list := schema.InitTodoList()
    c.IndentedJSON(http.StatusOK, list)
  })
  
  router.Use(TokenAuthMiddleware)

  router.GET("/graphql", func (c *gin.Context) {
    // http://localhost:8080/graphql?query={todo(id:"a"){id, text}}
    result := executeQuery(c.Query("query"), schema.TodoSchema)

    if result.HasErrors() {
      c.IndentedJSON(http.StatusBadRequest, result.Errors)
    } else {
    c.IndentedJSON(http.StatusOK, result) 
    }

  })

  if err := router.Run(":8080"); err != nil {
    log.Fatalf("Couldnot run the server %v", err)
  }
}

func TokenAuthMiddleware(c *gin.Context) {
  token := c.Query("token")
  
  if token == "" {
    c.AbortWithStatusJSON(http.StatusForbidden, "Unauthorize")
  }

  c.Next()
}

func Greeting(c *gin.Context) {
  query := c.Request.URL.Query().Get("query")
  fmt.Println("Query: ", query)
  fmt.Println("Greeting: ")
  c.IndentedJSON(http.StatusOK, "hi")
}
