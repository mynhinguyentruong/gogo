package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type GithubAccessTokenResponse struct{
  AccessToken string `json:"access_token"`
  Scope string `json:"scope"`
  TokenType string `json:"token_type"` 
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

