package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

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

type GithubAccessTokenResponse struct{
  AccessToken string `json:"access_token"`
  Scope string `json:"scope"`
  TokenType string `json:"token_type"` 
}

func main() {
  router := gin.Default()
  
  router.Use(ExperimentalMiddleware)
  router.Use(CORSMiddleware())

  router.SetTrustedProxies([]string{"127.0.0.69"})

  router.GET("/greet", Greeting)
  router.GET("/list", func (c *gin.Context) {
    list := schema.InitTodoList()
    c.IndentedJSON(http.StatusOK, list)
  })
  router.GET("/api/auth/callback", func (c *gin.Context) {
    method := c.Query("method")

    if method == "github_oauth" {
      code := c.Query("code")
      fmt.Println("code: ", code)
      
      // POST https://github.com/login/oauth/access_token?

      params := url.Values{}
      params.Add("client_id", os.Getenv("github_clientid"))
      params.Add("client_secret", os.Getenv("github_clientsecret"))
      params.Add("code", code)

      url := "https://github.com/login/oauth/access_token?"

      resp, err := http.Post(url + params.Encode(), "application/json", nil)
      if err != nil {
        c.AbortWithStatusJSON(500, err)
      }

      data, err := io.ReadAll(resp.Body)
      if err != nil {
        c.AbortWithStatusJSON(500, err)
      }

      var access_token_response GithubAccessTokenResponse

      err = json.Unmarshal(data, &access_token_response)
      if err != nil {
        c.AbortWithStatusJSON(500, err)
      }

      // Set Cookie to access_token_response
      // key: backend_auth.session_token
      // value: access_token_response.AccessToken

      // Hash
      // Salt
      // Save to db
      c.SetCookie("backend_auth.session", access_token_response.AccessToken, 3600, "/", "localhost", true, true)

      c.IndentedJSON(http.StatusOK, code)
    }

    if method == "" {
      c.AbortWithError(400, errors.New("unhandled method"))
    }
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

func ExperimentalMiddleware(c *gin.Context) {
  fmt.Println("ExperimentalMiddleware ran")

  c.Next()
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

func CORSMiddleware() gin.HandlerFunc {
  log.Println("CORSMiddleware ran")

    return func(c *gin.Context) {
      fmt.Println("The return func ran")
      c.Writer.Header().Set("Access-Control-Allow-Origin", "https://go-graphql-test123.fly.dev")
      c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
      c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
      c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

      if c.Request.Method == "OPTIONS" {
          c.AbortWithStatus(204)
          return
      }

      c.Next()
    }
}
