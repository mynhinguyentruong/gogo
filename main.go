package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"crypto/sha256"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/mynhinguyentruong/gogo/schema"
	"golang.org/x/crypto/bcrypt"
)

func encryptTokenToBase64String(access_token string) string {
  hash := sha256.New() 
  hash.Write([]byte(access_token))

  bs := hash.Sum(nil)
  hashedToken, err := bcrypt.GenerateFromPassword(bs, bcrypt.DefaultCost)
  if err != nil {
    panic(err)
  }

  strEncoded := base64.StdEncoding.EncodeToString(hashedToken)
 
  return strEncoded 
  // Next: save this to a DB with key access_token
}

func getGithubUser(access_token string) map[string]interface{} {
  client := &http.Client{}
  req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
  if err != nil {
    log.Fatalf("Error when create new request: %v", req)
  }

  req.Header.Add("Authorization", "Bearer " + access_token)

  resp, err := client.Do(req)
  if err != nil {
    log.Fatalf("Error when sending Get User request to github: %v", err)
  }

  var data map[string]interface{}
  dataBytes, err := io.ReadAll(resp.Body)
  if err != nil {
    fmt.Println("Error while reading response body to bytes", err)
  }
  resp.Body.Close()
  err = json.Unmarshal(dataBytes, &data)
  if err != nil {
    fmt.Println("Error while Unmarshal JSON to map: ", err)
  }

  log.Printf("Github user: %v", data)

  return data
}

func verifyAccessToken(str1, str2 string) bool {
  err := bcrypt.CompareHashAndPassword([]byte(str1), []byte(str2)) 
  
  return err == nil
}

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
  port := os.Getenv("PORT")
  if port == "" {
    log.Println("Cannot find PORT env, default to run on port 8080 instead")
    port = "8080"
  }
  router := gin.Default()
  
  router.Use(ExperimentalMiddleware)
  router.Use(CORSMiddleware())

  router.SetTrustedProxies([]string{"127.0.0.69"})

  router.GET("/test", func (c *gin.Context) {
    val := encryptTokenToBase64String("123")
    fmt.Println("Val: ", val)
    c.IndentedJSON(http.StatusOK, val)
  })

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

      if os.Getenv("github_clientsecret") == "" || os.Getenv("github_clientid") == "" {
        log.Fatalf("Provide github env\n github_clientid: %v \n github_clientsecret: %v", os.Getenv("github_clientid"), os.Getenv("github_clientsecret"))
      }

      gitURL := "https://github.com/login/oauth/access_token?" + params.Encode()

      resp, err := http.Get(gitURL)
      if err != nil {
        log.Fatalf("error in get req: %v", err)
        // c.AbortWithStatusJSON(500, err)
        return
      }

      data, err := io.ReadAll(resp.Body)
      resp.Body.Close()
      if err != nil {
        log.Fatal("Error in reading data: ", err)
        return
      }

      fmt.Println("Data: ", data)
      fmt.Println("Data: ", data)
      fmt.Println("Data: ", data)

      fmt.Println("Haha: ", string(data))
      fmt.Println("Haha: ", string(data))
      fmt.Println("Haha: ", string(data))

      access_token_response := string(data) 
      
      m, _ := url.ParseQuery(access_token_response)
      // type Values map[string][]string

      fmt.Println("m: ", m)
      fmt.Println("m: ", m)
      fmt.Println("m: ", m)
      //
      // //
      // // fmt.Println("Response: ", access_token_response)
      // // fmt.Println("Response: ", access_token_response)
      // // fmt.Println("Response: ", access_token_response)

      fmt.Println("Access token: ", m["access_token"])
      // fmt.Println("Access token: ", access_token_response.AccessToken)
      // fmt.Println("Access token: ", access_token_response.AccessToken)

      getGithubUser(m["access_token"][0])

      // Set Cookie to access_token_response
      // key: backend_auth.session_token
      // value: access_token_response.AccessToken

      // Hash
      // Salt
      // Save to db
      c.SetCookie("backend_auth.session", m["access_token"][0], 3600, "/", "localhost", true, true)
      c.SetSameSite(http.SameSiteLaxMode)

      c.IndentedJSON(http.StatusOK, m["access_token"][0])
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

  if err := router.Run(":"+port); err != nil {
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
